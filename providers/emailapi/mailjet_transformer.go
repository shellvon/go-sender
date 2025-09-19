package emailapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

// MailjetAttachment represents an attachment in Mailjet format.
type MailjetAttachment struct {
	ContentType   string `json:"ContentType"`
	Filename      string `json:"Filename"`
	Base64Content string `json:"Base64Content"`
}

// MailjetMessage represents a message in Mailjet format.
type MailjetMessage struct {
	From                   EmailAddress           `json:"From"`
	To                     []EmailAddress         `json:"To"`
	Cc                     []EmailAddress         `json:"Cc,omitempty"`
	Bcc                    []EmailAddress         `json:"Bcc,omitempty"`
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
	TemplateErrorReporting EmailAddress           `json:"TemplateErrorReporting,omitempty"`
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

	mt.setSender(&message, msg, account)
	mt.setRecipients(&message, msg)
	mt.setContent(&message, msg)
	mt.setTemplate(&message, msg)
	mt.setAttachments(&message, msg)
	mt.setHeaders(&message, msg)
	mt.setExtras(&message, msg)

	sandboxMode := mt.getSandboxMode(msg)
	return MailjetRequest{
		Messages:    []MailjetMessage{message},
		SandBoxMode: sandboxMode,
	}
}

// setSender configures the sender information.
func (mt *mailjetTransformer) setSender(message *MailjetMessage, msg *Message, account *Account) {
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr != "" {
		message.From = parseEmailAddress(fromAddr)
	}
}

// setRecipients configures all recipient addresses.
func (mt *mailjetTransformer) setRecipients(message *MailjetMessage, msg *Message) {
	if len(msg.To) > 0 {
		message.To = parseEmailAddresses(msg.To)
	}
	if len(msg.Cc) > 0 {
		message.Cc = parseEmailAddresses(msg.Cc)
	}
	if len(msg.Bcc) > 0 {
		message.Bcc = parseEmailAddresses(msg.Bcc)
	}
}

// setContent configures the email content.
func (mt *mailjetTransformer) setContent(message *MailjetMessage, msg *Message) {
	if msg.Subject != "" {
		message.Subject = msg.Subject
	}
	if msg.Text != "" {
		message.TextPart = msg.Text
	}
	if msg.HTML != "" {
		message.HTMLPart = msg.HTML
	}
}

// setTemplate configures template support.
func (mt *mailjetTransformer) setTemplate(message *MailjetMessage, msg *Message) {
	if msg.TemplateID == "" {
		return
	}

	if templateID, err := strconv.Atoi(msg.TemplateID); err == nil {
		message.TemplateID = templateID
		message.TemplateLanguage = true
		if msg.TemplateData != nil {
			message.Variables = msg.TemplateData
		}
	}
}

// setAttachments configures email attachments.
func (mt *mailjetTransformer) setAttachments(message *MailjetMessage, msg *Message) {
	if len(msg.Attachments) == 0 {
		return
	}

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

// setHeaders configures custom headers.
func (mt *mailjetTransformer) setHeaders(message *MailjetMessage, msg *Message) {
	if len(msg.Headers) > 0 {
		message.Headers = msg.Headers
	}
}

// setExtras configures additional fields from extras.
func (mt *mailjetTransformer) setExtras(message *MailjetMessage, msg *Message) {
	if customID, ok := msg.Extras["custom_id"]; ok {
		if id, idOk := customID.(string); idOk {
			message.CustomID = id
		}
	}

	if campaign, ok := msg.Extras["campaign"]; ok {
		if c, campaignOk := campaign.(string); campaignOk {
			message.CustomCampaign = c
		}
	}

	if urlTags, ok := msg.Extras["url_tags"]; ok {
		if tags, tagsOk := urlTags.(string); tagsOk {
			message.URLTags = tags
		}
	}
}

// getSandboxMode extracts sandbox mode from extras.
func (mt *mailjetTransformer) getSandboxMode(msg *Message) bool {
	if sandbox, ok := msg.Extras["sandbox"]; ok {
		if sb, sbOk := sandbox.(bool); sbOk {
			return sb
		}
	}
	return false
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
