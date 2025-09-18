package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// brevoTransformer implements HTTPRequestTransformer for Brevo.
// API Reference: https://developers.brevo.com/reference/sendtransacemail

func init() {
	RegisterTransformer(string(SubProviderBrevo), newBrevoTransformer())
}

const (
	brevoAPIPath            = "https://api.brevo.com/v3/smtp/email"
	brevoMaxRecipients      = 50
	brevoTagsKey            = "tags"
	brevoBatchIdKey         = "batchId"
	brevoMessageVersionsKey = "messageVersions"
)

// brevoTransformer implements Brevo logic via BaseHTTPTransformer.
type brevoTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func newBrevoTransformer() *brevoTransformer {
	bt := &brevoTransformer{}
	bt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeEmailAPI,
		string(SubProviderBrevo),
		&core.ResponseHandlerConfig{},
		bt.transform,
		transformer.AddBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			return bt.validate(msg, account)
		}),
	)
	return bt
}

// Note: Using unified EmailAddress from email_utils.go instead of provider-specific types

// BrevoAttachment represents an attachment in Brevo format.
type BrevoAttachment struct {
	URL     string `json:"url"`     // URL of the attachment
	Content string `json:"content"` // Base64 encoded content
	Name    string `json:"name"`    // Filename
}

// transform handles Brevo message → HTTPRequestSpec.
func (bt *brevoTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// Prepare request body
	requestBody := bt.buildRequestBody(msg, account)

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Brevo request body: %w", err)
	}

	headers := map[string]string{
		"api-key": account.APIKey,
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      brevoAPIPath,
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

// buildRequestBody constructs the Brevo API request body.
func (bt *brevoTransformer) buildRequestBody(msg *Message, account *Account) map[string]interface{} {
	body := make(map[string]interface{})

	// Sender (required) - use account.From as default if msg.From is empty
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	// Mandatory if templateId is not passed.
	// Pass name (optional) and email or id of sender from which emails will be sent.
	// name will be ignored if passed along with sender id.
	// For example,{"name":"Mary from MyShop", "email":"no-reply@myshop.com"}
	// {"id":2}
	if fromAddr != "" {
		if !strings.Contains(fromAddr, "@") {
			body["sender"] = map[string]string{
				"id": fromAddr,
			}
		} else {
			body["sender"] = parseEmailAddress(fromAddr)
		}
	}

	// To addresses (required)
	if len(msg.To) > 0 {
		toAddrs := parseEmailAddresses(msg.To)
		body["to"] = toAddrs
	}

	// CC addresses
	if len(msg.Cc) > 0 {
		ccAddrs := parseEmailAddresses(msg.Cc)
		body["cc"] = ccAddrs
	}

	// BCC addresses
	if len(msg.Bcc) > 0 {
		bccAddrs := parseEmailAddresses(msg.Bcc)
		body["bcc"] = bccAddrs
	}

	// Reply-to address (Brevo supports single reply-to)
	if len(msg.ReplyTo) > 0 {
		replyToAddr := parseEmailAddress(msg.ReplyTo[0]) // Use first reply-to address
		body["replyTo"] = replyToAddr
	}

	// Subject
	if msg.Subject != "" {
		body["subject"] = msg.Subject
	}

	// Preheader (preview text)
	if preheader, ok := msg.Extras["preheader"]; ok {
		body["preheader"] = preheader
	}

	// Content
	if msg.HTML != "" {
		body["htmlContent"] = msg.HTML
	}
	if msg.Text != "" {
		body["textContent"] = msg.Text
	}

	// Template support
	if msg.TemplateID != "" {
		templateIDInt := 0
		if id, err := parseTemplateID(msg.TemplateID); err == nil {
			templateIDInt = id
		}
		body["templateId"] = templateIDInt

		// Template parameters - direct mapping from TemplateData
		if msg.TemplateData != nil {
			body["params"] = msg.TemplateData
		}
	}

	// Attachments
	if len(msg.Attachments) > 0 {
		attachments := make([]BrevoAttachment, len(msg.Attachments))
		for i, att := range msg.Attachments {
			attachments[i] = BrevoAttachment{
				Content: string(att.Content),
				Name:    att.Filename,
			}
			contentOrURL := string(att.Content)
			// Pass the absolute URL (no local file) or the base64 content of the attachment along with the attachment name.
			// Mandatory if attachment content is passed.
			if strings.HasPrefix(contentOrURL, "http://") || strings.HasPrefix(contentOrURL, "https://") {
				attachments[i].URL = contentOrURL
				attachments[i].Content = ""
			}
		}
		body["attachment"] = attachments
	}

	// Custom headers (use the standard headers from message)
	if len(msg.Headers) > 0 {
		body["headers"] = msg.Headers
	}

	// Tags from extras
	if tags, ok := msg.Extras[brevoTagsKey]; ok {
		body["tags"] = tags
	}

	// Scheduled sending
	if msg.ScheduledAt != nil {
		body["scheduledAt"] = msg.ScheduledAt.Format(time.RFC3339)
	}

	// Batch ID for grouping emails
	if batchId, ok := msg.Extras[brevoBatchIdKey]; ok {
		body["batchId"] = batchId
	}

	// Message versions
	if messageVersions, ok := msg.Extras[brevoMessageVersionsKey]; ok {
		body["messageVersions"] = messageVersions
	}

	return body
}

// validate checks if the message and account are valid for Brevo.
func (bt *brevoTransformer) validate(msg *Message, account *Account) error {
	if account.APIKey == "" {
		return errors.New("APIKey is required for Brevo")
	}

	if len(msg.To) == 0 {
		return errors.New("to recipients cannot be empty")
	}

	if len(msg.To) > brevoMaxRecipients {
		return fmt.Errorf("to recipients are limited to %d recipients", brevoMaxRecipients)
	}

	// From field validation - check both message and account default
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr == "" {
		return errors.New("sender is required for Brevo")
	}

	// Either subject or template is required
	if msg.Subject == "" && msg.TemplateID == "" {
		return errors.New("either subject or templateId is required for Brevo")
	}

	// Either content or template is required
	if msg.HTML == "" && msg.Text == "" && msg.TemplateID == "" {
		return errors.New("either content (htmlContent/textContent) or templateId is required for Brevo")
	}

	return nil
}

// parseTemplateID attempts to parse template ID as integer.
func parseTemplateID(templateID string) (int, error) {
	if templateID == "" {
		return 0, errors.New("template ID is empty")
	}

	// Simple integer parsing
	result := 0
	for _, char := range templateID {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("invalid template ID: %s", templateID)
		}
		result = result*10 + int(char-'0')
	}

	return result, nil
}
