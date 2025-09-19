package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// mailerSendTransformer implements HTTPRequestTransformer for MailerSend.
// MailerSend API requires a Bearer token for authentication.
// API Reference: https://developers.mailersend.com/api/v1/email.html#send-an-email

func init() {
	RegisterTransformer(string(SubProviderMailerSend), newMailerSendTransformer())
}

const (
	mailerSendAPIPath       = "https://api.mailersend.com/v1/email"
	mailerSendMaxRecipients = 50

	mailerSendTagsKey            = "tags"
	mailerSendSettingsKey        = "settings"
	mailerSendInReplyToKey       = "in_reply_to"
	mailerSendListUnsubscribeKey = "list_unsubscribe"
)

// mailerSendTransformer implements MailerSend logic via BaseHTTPTransformer.
type mailerSendTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func newMailerSendTransformer() *mailerSendTransformer {
	mt := &mailerSendTransformer{}
	mt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeEmailAPI,
		string(SubProviderMailerSend),
		&core.ResponseHandlerConfig{},
		mt.transform,
		transformer.AddBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			return mt.validate(msg, account)
		}),
	)
	return mt
}

// MailerSendAttachment represents an attachment in MailerSend format.
type MailerSendAttachment struct {
	Content     string `json:"content"` // Base64 encoded content
	Filename    string `json:"filename"`
	Disposition string `json:"disposition"`  // "attachment" or "inline"
	ID          string `json:"id,omitempty"` // For inline attachments
}

// MailerSendPersonalization represents personalization data for specific recipients.
type MailerSendPersonalization struct {
	Email string                 `json:"email"`
	Data  map[string]interface{} `json:"data"`
}

// MailerSendHeader represents a custom header in MailerSend format.
type MailerSendHeader struct {
	Name  string `json:"name"`  // Header name
	Value string `json:"value"` // Header value
}

// MailerSendSettings represents email settings.
type MailerSendSettings struct {
	TrackClicks      bool `json:"track_clicks,omitempty"`
	TrackOpens       bool `json:"track_opens,omitempty"`
	TrackContent     bool `json:"track_content,omitempty"`
	TrackUnsubscribe bool `json:"track_unsubscribe,omitempty"`
}

// transform handles MailerSend message → HTTPRequestSpec.
func (mt *mailerSendTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// Prepare request body
	requestBody := mt.buildRequestBody(msg, account)

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal MailerSend request body: %w", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + account.APIKey,
		"Content-Type":  "application/json",
	}

	// Add idempotency key if message ID is available
	if msgID := msg.MsgID(); msgID != "" {
		headers["X-Message-Id"] = msgID
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      mailerSendAPIPath,
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

// buildRequestBody constructs the MailerSend API request body.
func (mt *mailerSendTransformer) buildRequestBody(msg *Message, account *Account) map[string]interface{} {
	body := make(map[string]interface{})

	mt.setSender(body, msg, account)
	mt.setRecipients(body, msg)
	mt.setContent(body, msg)
	mt.setAttachments(body, msg)
	mt.setTemplate(body, msg)
	mt.setPersonalization(body, msg)
	mt.setHeaders(body, msg)
	mt.setExtras(body, msg)

	return body
}

// setSender configures the sender information.
func (mt *mailerSendTransformer) setSender(body map[string]interface{}, msg *Message, account *Account) {
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr != "" {
		body["from"] = parseEmailAddress(fromAddr)
	}
}

// setRecipients configures all recipient addresses.
func (mt *mailerSendTransformer) setRecipients(body map[string]interface{}, msg *Message) {
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
		body["reply_to"] = parseEmailAddress(msg.ReplyTo[0])
	}
}

// setContent configures the email content.
func (mt *mailerSendTransformer) setContent(body map[string]interface{}, msg *Message) {
	if msg.Subject != "" {
		body["subject"] = msg.Subject
	}
	if msg.Text != "" {
		body["text"] = msg.Text
	}
	if msg.HTML != "" {
		body["html"] = msg.HTML
	}
}

// setAttachments configures email attachments.
func (mt *mailerSendTransformer) setAttachments(body map[string]interface{}, msg *Message) {
	if len(msg.Attachments) == 0 {
		return
	}

	attachments := make([]MailerSendAttachment, len(msg.Attachments))
	for i, att := range msg.Attachments {
		disposition := att.Disposition
		if disposition == "" {
			disposition = "attachment"
		}

		attachments[i] = MailerSendAttachment{
			Content:     string(att.Content),
			Filename:    att.Filename,
			Disposition: disposition,
			ID:          att.ContentID,
		}
	}
	body["attachments"] = attachments
}

// setTemplate configures template support.
func (mt *mailerSendTransformer) setTemplate(body map[string]interface{}, msg *Message) {
	if msg.TemplateID != "" {
		body["template_id"] = msg.TemplateID
	}
}

// setPersonalization configures per-recipient customization.
func (mt *mailerSendTransformer) setPersonalization(body map[string]interface{}, msg *Message) {
	if len(msg.TemplateData) == 0 {
		return
	}

	personalization := make([]MailerSendPersonalization, 0, len(msg.TemplateData))
	for email, data := range msg.TemplateData {
		if !strings.Contains(email, "@") {
			continue
		}
		if dataMap, ok := data.(map[string]interface{}); ok {
			personalization = append(personalization, MailerSendPersonalization{
				Email: email,
				Data:  dataMap,
			})
		}
	}

	if len(personalization) > 0 {
		body["personalization"] = personalization
	}
}

// setHeaders configures custom headers.
func (mt *mailerSendTransformer) setHeaders(body map[string]interface{}, msg *Message) {
	if len(msg.Headers) == 0 {
		return
	}

	headers := make([]MailerSendHeader, 0, len(msg.Headers))
	for name, value := range msg.Headers {
		headers = append(headers, MailerSendHeader{
			Name:  name,
			Value: value,
		})
	}
	body["headers"] = headers
}

// setExtras configures additional fields from extras and scheduling.
func (mt *mailerSendTransformer) setExtras(body map[string]interface{}, msg *Message) {
	if msg.ScheduledAt != nil {
		body["send_at"] = msg.ScheduledAt.Unix()
	}

	// Tags from extras
	if tags, ok := msg.Extras[mailerSendTagsKey]; ok {
		body["tags"] = tags
	}

	// Email settings from extras (track clicks, opens, etc.)
	if settings, ok := msg.Extras[mailerSendSettingsKey]; ok {
		body["settings"] = settings
	}

	// In-reply-to reference
	if inReplyTo, ok := msg.Extras[mailerSendInReplyToKey]; ok {
		body["in_reply_to"] = inReplyTo
	}

	// List-Unsubscribe header (RFC 2369)
	if listUnsubscribe, ok := msg.Extras[mailerSendListUnsubscribeKey]; ok {
		body["list_unsubscribe"] = listUnsubscribe
	}
}

// validate checks if the message and account are valid for MailerSend.
func (mt *mailerSendTransformer) validate(msg *Message, account *Account) error {
	if account.APIKey == "" {
		return errors.New("APIKey is required for MailerSend")
	}

	if len(msg.To) == 0 {
		return errors.New("to recipients cannot be empty")
	}

	if len(msg.To) > mailerSendMaxRecipients {
		return fmt.Errorf("to recipients are limited to %d recipients", mailerSendMaxRecipients)
	}

	// From field validation according to MailerSend API docs:
	// - Required for regular emails
	// - Optional for template emails (template can have default sender)
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}

	if msg.TemplateID == "" && fromAddr == "" {
		return errors.New("from is required for non-template emails")
	}

	// For template emails, content fields are usually not needed
	// For non-template emails, either subject or content is required
	if msg.TemplateID == "" {
		if msg.Subject == "" {
			return errors.New("subject is required for non-template emails")
		}
		if msg.HTML == "" && msg.Text == "" {
			return errors.New("either HTML or text content is required for non-template emails")
		}
	}

	return nil
}
