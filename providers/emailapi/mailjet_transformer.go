package emailapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// mailjetTransformer implements HTTPRequestTransformer for Mailjet.
// API Reference: https://dev.mailjet.com/email/guides/send-api-v31/#send-a-basic-email

func init() {
	RegisterTransformer(string(SubProviderMailjet), newMailjetTransformer())
}

const (
	mailjetAPIPath       = "https://api.mailjet.com/v3.1/send"
	mailjetMaxRecipients = 50
)

// MailjetEmail represents an email in Mailjet format.
type MailjetEmail struct {
	Email string `json:"Email"`
	Name  string `json:"Name,omitempty"`
}

// MailjetAttachment represents an attachment in Mailjet format.
type MailjetAttachment struct {
	ContentType   string `json:"ContentType"`
	Filename      string `json:"Filename"`
	Base64Content string `json:"Base64Content"`
}

// MailjetMessage represents a message in Mailjet format.
type MailjetMessage struct {
	From                   MailjetEmail           `json:"From"`
	To                     []MailjetEmail         `json:"To"`
	Cc                     []MailjetEmail         `json:"Cc,omitempty"`
	Bcc                    []MailjetEmail         `json:"Bcc,omitempty"`
	Subject                string                 `json:"Subject"`
	TextPart               string                 `json:"TextPart,omitempty"`
	HTMLPart               string                 `json:"HTMLPart,omitempty"`
	TemplateID             int                    `json:"TemplateID,omitempty"`
	TemplateLanguage       bool                   `json:"TemplateLanguage,omitempty"`
	Variables              map[string]interface{} `json:"Variables,omitempty"`
	Attachments            []MailjetAttachment    `json:"Attachments,omitempty"`
	InlinedAttachments     []MailjetAttachment    `json:"InlinedAttachments,omitempty"`
	Headers                map[string]string      `json:"Headers,omitempty"`
	CustomID               string                 `json:"CustomID,omitempty"`
	EventPayload           string                 `json:"EventPayload,omitempty"`
	CustomCampaign         string                 `json:"CustomCampaign,omitempty"`
	DeduplicateCampaign    bool                   `json:"DeduplicateCampaign,omitempty"`
	URLTags                string                 `json:"URLTags,omitempty"`
	TemplateErrorReporting MailjetEmail           `json:"TemplateErrorReporting,omitempty"`
	TemplateErrorDeliver   bool                   `json:"TemplateErrorDeliver,omitempty"`
}

// MailjetRequest represents the complete Mailjet API request.
type MailjetRequest struct {
	Messages    []MailjetMessage `json:"Messages"`
	SandBoxMode bool             `json:"SandBoxMode,omitempty"`
}

// mailjetTransformer implements Mailjet logic via BaseHTTPTransformer.
type mailjetTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func newMailjetTransformer() *mailjetTransformer {
	mt := &mailjetTransformer{}
	mt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeEmailAPI,
		string(SubProviderMailjet),
		&core.ResponseHandlerConfig{},
		mt.transform,
		transformer.AddBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			return mt.validate(msg, account)
		}),
	)
	return mt
}

// transform handles Mailjet message → HTTPRequestSpec.
func (mt *mailjetTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// Build request body
	requestBody := mt.buildRequestBody(msg, account)

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Mailjet request body: %w", err)
	}

	// Basic authentication with API key and secret
	auth := base64.StdEncoding.EncodeToString([]byte(account.APIKey + ":" + account.APISecret))

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + auth,
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      mailjetAPIPath,
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

// buildRequestBody constructs the Mailjet API request body.
func (mt *mailjetTransformer) buildRequestBody(msg *Message, account *Account) MailjetRequest {
	message := MailjetMessage{}

	// From (required)
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr != "" {
		parsedFrom := parseEmailAddress(fromAddr)
		message.From = MailjetEmail{
			Email: parsedFrom.Email,
			Name:  parsedFrom.Name,
		}
	}

	// To (required)
	if len(msg.To) > 0 {
		toEmails := make([]MailjetEmail, len(msg.To))
		for i, addr := range msg.To {
			parsed := parseEmailAddress(addr)
			toEmails[i] = MailjetEmail{
				Email: parsed.Email,
				Name:  parsed.Name,
			}
		}
		message.To = toEmails
	}

	// CC
	if len(msg.Cc) > 0 {
		ccEmails := make([]MailjetEmail, len(msg.Cc))
		for i, addr := range msg.Cc {
			parsed := parseEmailAddress(addr)
			ccEmails[i] = MailjetEmail{
				Email: parsed.Email,
				Name:  parsed.Name,
			}
		}
		message.Cc = ccEmails
	}

	// BCC
	if len(msg.Bcc) > 0 {
		bccEmails := make([]MailjetEmail, len(msg.Bcc))
		for i, addr := range msg.Bcc {
			parsed := parseEmailAddress(addr)
			bccEmails[i] = MailjetEmail{
				Email: parsed.Email,
				Name:  parsed.Name,
			}
		}
		message.Bcc = bccEmails
	}

	// Subject
	if msg.Subject != "" {
		message.Subject = msg.Subject
	}

	// Content
	if msg.Text != "" {
		message.TextPart = msg.Text
	}
	if msg.HTML != "" {
		message.HTMLPart = msg.HTML
	}

	// Template support
	if msg.TemplateID != "" {
		if templateID, err := parseTemplateIDInt(msg.TemplateID); err == nil {
			message.TemplateID = templateID
			message.TemplateLanguage = true

			// Template variables
			if msg.TemplateData != nil {
				message.Variables = msg.TemplateData
			}
		}
	}

	// Attachments
	if len(msg.Attachments) > 0 {
		attachments := make([]MailjetAttachment, 0)
		inlinedAttachments := make([]MailjetAttachment, 0)

		for _, att := range msg.Attachments {
			attachment := MailjetAttachment{
				ContentType:   att.ContentType,
				Filename:      att.Filename,
				Base64Content: base64.StdEncoding.EncodeToString(att.Content),
			}

			if att.Disposition == "inline" {
				inlinedAttachments = append(inlinedAttachments, attachment)
			} else {
				attachments = append(attachments, attachment)
			}
		}

		if len(attachments) > 0 {
			message.Attachments = attachments
		}
		if len(inlinedAttachments) > 0 {
			message.InlinedAttachments = inlinedAttachments
		}
	}

	// Custom headers
	if len(msg.Headers) > 0 {
		message.Headers = msg.Headers
	}

	// Custom ID from extras
	if customID, ok := msg.Extras["custom_id"]; ok {
		if id, ok := customID.(string); ok {
			message.CustomID = id
		}
	}

	// Campaign tracking
	if campaign, ok := msg.Extras["campaign"]; ok {
		if c, ok := campaign.(string); ok {
			message.CustomCampaign = c
		}
	}

	// URL Tags
	if urlTags, ok := msg.Extras["url_tags"]; ok {
		if tags, ok := urlTags.(string); ok {
			message.URLTags = tags
		}
	}

	// Sandbox mode
	sandboxMode := false
	if sandbox, ok := msg.Extras["sandbox"]; ok {
		if sb, ok := sandbox.(bool); ok {
			sandboxMode = sb
		}
	}

	return MailjetRequest{
		Messages:    []MailjetMessage{message},
		SandBoxMode: sandboxMode,
	}
}

// validate checks if the message and account are valid for Mailjet.
func (mt *mailjetTransformer) validate(msg *Message, account *Account) error {
	if account.APIKey == "" {
		return errors.New("APIKey is required for Mailjet")
	}

	if account.APISecret == "" {
		return errors.New("APISecret is required for Mailjet")
	}

	if len(msg.To) == 0 {
		return errors.New("to recipients cannot be empty")
	}

	if len(msg.To) > mailjetMaxRecipients {
		return fmt.Errorf("to recipients are limited to %d recipients", mailjetMaxRecipients)
	}

	// From field validation
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr == "" {
		return errors.New("from is required for Mailjet")
	}

	// Subject validation
	if msg.Subject == "" {
		return errors.New("subject is required for Mailjet")
	}

	// Content validation
	if msg.Text == "" && msg.HTML == "" && msg.TemplateID == "" {
		return errors.New("either text, html content or template is required for Mailjet")
	}

	return nil
}

// parseTemplateIDInt parses template ID as integer for Mailjet.
func parseTemplateIDInt(templateID string) (int, error) {
	if templateID == "" {
		return 0, errors.New("template ID is empty")
	}

	result := 0
	for _, char := range templateID {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("invalid template ID: %s", templateID)
		}
		result = result*10 + int(char-'0')
	}

	return result, nil
}
