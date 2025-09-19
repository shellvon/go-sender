package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	brevoBatchIDKey         = "batchId"
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

	bt.setSender(body, msg, account)
	bt.setRecipients(body, msg)
	bt.setContent(body, msg)
	bt.setTemplate(body, msg)
	bt.setAttachments(body, msg)
	bt.setExtras(body, msg)

	return body
}

// setSender configures the sender information.
func (bt *brevoTransformer) setSender(body map[string]interface{}, msg *Message, account *Account) {
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
}

// setRecipients configures all recipient addresses.
func (bt *brevoTransformer) setRecipients(body map[string]interface{}, msg *Message) {
	if len(msg.To) > 0 {
		body["to"] = parseEmailAddresses(msg.To)
	}
	if len(msg.Cc) > 0 {
		body["cc"] = parseEmailAddresses(msg.Cc)
	}
	if len(msg.Bcc) > 0 {
		body["bcc"] = parseEmailAddresses(msg.Bcc)
	}
	if len(msg.ReplyTo) > 0 {
		body["replyTo"] = parseEmailAddress(msg.ReplyTo[0])
	}
}

// setContent configures the email content.
func (bt *brevoTransformer) setContent(body map[string]interface{}, msg *Message) {
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
}

// setTemplate configures template-related fields.
func (bt *brevoTransformer) setTemplate(body map[string]interface{}, msg *Message) {
	if msg.TemplateID != "" {
		templateIDInt := 0
		if id, err := strconv.Atoi(msg.TemplateID); err == nil {
			templateIDInt = id
		}
		body["templateId"] = templateIDInt

		if msg.TemplateData != nil {
			body["params"] = msg.TemplateData
		}
	}
}

// setAttachments configures email attachments.
func (bt *brevoTransformer) setAttachments(body map[string]interface{}, msg *Message) {
	if len(msg.Attachments) == 0 {
		return
	}

	attachments := make([]BrevoAttachment, len(msg.Attachments))
	for i, att := range msg.Attachments {
		attachments[i] = BrevoAttachment{
			Content: string(att.Content),
			Name:    att.Filename,
		}
		contentOrURL := string(att.Content)
		if strings.HasPrefix(contentOrURL, "http://") || strings.HasPrefix(contentOrURL, "https://") {
			attachments[i].URL = contentOrURL
			attachments[i].Content = ""
		}
	}
	body["attachment"] = attachments
}

// setExtras configures additional fields from extras and scheduling.
func (bt *brevoTransformer) setExtras(body map[string]interface{}, msg *Message) {
	if len(msg.Headers) > 0 {
		body["headers"] = msg.Headers
	}
	if tags, ok := msg.Extras[brevoTagsKey]; ok {
		body["tags"] = tags
	}
	if msg.ScheduledAt != nil {
		body["scheduledAt"] = msg.ScheduledAt.Format(time.RFC3339)
	}
	if batchID, ok := msg.Extras[brevoBatchIDKey]; ok {
		body["batchId"] = batchID
	}
	if messageVersions, ok := msg.Extras[brevoMessageVersionsKey]; ok {
		body["messageVersions"] = messageVersions
	}
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
