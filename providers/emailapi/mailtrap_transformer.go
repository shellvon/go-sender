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

// mailtrapTransformer implements HTTPRequestTransformer for Mailtrap.
// API Reference: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/send-email
func init() {
	RegisterTransformer(string(SubProviderMailtrap), newMailtrapTransformer())
}

const (
	mailtrapAPIPath            = "https://send.api.mailtrap.io/api/send"
	mailtrapMaxRecipients      = 1000
	mailtrapCategoryKey        = "category"
	mailtrapCustomVariablesKey = "custom_variables"
)

// mailtrapTransformer implements Mailtrap logic via BaseHTTPTransformer.
type mailtrapTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func newMailtrapTransformer() *mailtrapTransformer {
	mt := &mailtrapTransformer{}
	mt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeEmailAPI,
		string(SubProviderMailtrap),
		&core.ResponseHandlerConfig{},
		mt.transform,
		transformer.AddBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			return mt.validate(msg, account)
		}),
	)
	return mt
}

// MailtrapAttachment represents an attachment in Mailtrap format.
type MailtrapAttachment struct {
	Content     string `json:"content"`               // Base64 encoded content (required)
	Filename    string `json:"filename"`              // Filename (required)
	Type        string `json:"type,omitempty"`        // MIME type
	Disposition string `json:"disposition,omitempty"` // "attachment" or "inline", defaults to "attachment"
	ContentID   string `json:"content_id,omitempty"`  // Content ID for inline attachments
}

// transform handles Mailtrap message → HTTPRequestSpec.
func (mt *mailtrapTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// Prepare request body
	requestBody := mt.buildRequestBody(msg, account)

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Mailtrap request body: %w", err)
	}

	// Mailtrap supports two authentication methods as per OpenAPI spec:
	// 1. Api-Token header (preferred for API keys)
	// 2. Bearer token (for JWT tokens)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Determine authentication method based on APIKey format
	if account.APIKey != "" {
		// If it looks like a JWT (contains dots), use Bearer auth
		if strings.Contains(account.APIKey, ".") {
			headers["Authorization"] = "Bearer " + account.APIKey
		} else {
			// Otherwise use Api-Token header (standard API key)
			headers["Api-Token"] = account.APIKey
		}
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      mailtrapAPIPath,
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

// buildRequestBody constructs the Mailtrap API request body.
func (mt *mailtrapTransformer) buildRequestBody(msg *Message, account *Account) map[string]interface{} {
	body := make(map[string]interface{})

	// From address (required) - use account.From as default if msg.From is empty
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr != "" {
		parsedAddr := parseEmailAddress(fromAddr)
		body["from"] = parsedAddr
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

	// Reply-to address (Mailtrap supports only single reply-to address)
	if len(msg.ReplyTo) > 0 {
		// Use only the first reply-to address as Mailtrap API supports only one
		replyToAddr := parseEmailAddress(msg.ReplyTo[0])
		body["reply_to"] = replyToAddr
	}

	// Subject
	if msg.Subject != "" {
		body["subject"] = msg.Subject
	}

	// Content
	if msg.Text != "" {
		body["text"] = msg.Text
	}
	if msg.HTML != "" {
		body["html"] = msg.HTML
	}

	// Template support
	if msg.TemplateID != "" {
		body["template_uuid"] = msg.TemplateID

		// Template variables
		if msg.TemplateData != nil {
			body["template_variables"] = msg.TemplateData
		}
	}

	// Attachments - directly use the unified attachment structure
	if len(msg.Attachments) > 0 {
		attachments := make([]MailtrapAttachment, len(msg.Attachments))
		for i, att := range msg.Attachments {
			disposition := att.Disposition
			if disposition == "" {
				disposition = "attachment" // Default as per Mailtrap API
			}

			attachments[i] = MailtrapAttachment{
				Content:     string(att.Content), // Assuming base64 encoded
				Filename:    att.Filename,
				Type:        att.ContentType,
				Disposition: disposition,
				ContentID:   att.ContentID,
			}
		}
		body["attachments"] = attachments
	}

	// Custom headers
	if len(msg.Headers) > 0 {
		body["headers"] = msg.Headers
	}

	// Category from extras
	if category, ok := msg.Extras[mailtrapCategoryKey]; ok {
		body["category"] = category
	}

	// Custom variables from extras (following OpenAPI spec)
	if customVariables, ok := msg.Extras[mailtrapCustomVariablesKey]; ok {
		body["custom_variables"] = customVariables
	}

	return body
}

// validate checks if the message and account are valid for Mailtrap.
// Follows OpenAPI specification requirements and constraints.
// Uses the new email type detection for more precise validation.
func (mt *mailtrapTransformer) validate(msg *Message, account *Account) error {
	if account.APIKey == "" {
		return errors.New("APIKey is required for Mailtrap")
	}

	// Validate recipients according to OpenAPI spec
	// At least one recipient (to, cc, or bcc) is required
	if len(msg.To) == 0 && len(msg.Cc) == 0 && len(msg.Bcc) == 0 {
		return errors.New("at least one recipient (to, cc, or bcc) is required")
	}

	// Validate individual recipient limits (maxItems: 1000 each as per OpenAPI)
	if len(msg.To) > mailtrapMaxRecipients {
		return fmt.Errorf("to recipients are limited to %d recipients", mailtrapMaxRecipients)
	}
	if len(msg.Cc) > mailtrapMaxRecipients {
		return fmt.Errorf("cc recipients are limited to %d recipients", mailtrapMaxRecipients)
	}
	if len(msg.Bcc) > mailtrapMaxRecipients {
		return fmt.Errorf("bcc recipients are limited to %d recipients", mailtrapMaxRecipients)
	}

	// Mailtrap-specific validations

	// Mailtrap supports only single reply-to address
	if len(msg.ReplyTo) > 1 {
		return fmt.Errorf("provider Mailtrap supports only one reply-to address, got %d", len(msg.ReplyTo))
	}

	// Validate message content according to Mailtrap OpenAPI specification
	if err := mt.validateMailtrapSpecific(msg, account); err != nil {
		return fmt.Errorf("mailtrap validation failed: %w", err)
	}

	// Validate attachments if present
	for i, attachment := range msg.Attachments {
		if attachment.Filename == "" {
			return fmt.Errorf("attachment %d: filename is required", i)
		}
		if len(attachment.Content) == 0 {
			return fmt.Errorf("attachment %d: content is required", i)
		}
	}

	return nil
}

// validateMailtrapSpecific checks if the message satisfies Mailtrap OpenAPI requirements.
func (mt *mailtrapTransformer) validateMailtrapSpecific(msg *Message, account *Account) error {
	emailType := msg.GetEmailType()

	// Check common requirements for all types
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr == "" {
		return fmt.Errorf("from is required for %s email type", emailType.String())
	}

	// Check type-specific requirements according to Mailtrap OpenAPI schemas
	switch emailType {
	case EmailTypeText:
		// EmailWithText schema: requires from, subject, text
		if msg.Subject == "" {
			return fmt.Errorf("subject is required for %s email type", emailType.String())
		}
		if msg.Text == "" {
			return fmt.Errorf("text content is required for %s email type", emailType.String())
		}

	case EmailTypeHTML:
		// EmailWithHtml schema: requires from, subject, html
		if msg.Subject == "" {
			return fmt.Errorf("subject is required for %s email type", emailType.String())
		}
		if msg.HTML == "" {
			return fmt.Errorf("html content is required for %s email type", emailType.String())
		}

	case EmailTypeTextAndHTML:
		// EmailWithTextAndHtml schema: requires from, subject, text, html
		if msg.Subject == "" {
			return fmt.Errorf("subject is required for %s email type", emailType.String())
		}
		if msg.Text == "" {
			return fmt.Errorf("text content is required for %s email type", emailType.String())
		}
		if msg.HTML == "" {
			return fmt.Errorf("html content is required for %s email type", emailType.String())
		}

	case EmailTypeTemplate:
		// EmailFromTemplate schema: requires from, template_uuid
		if msg.TemplateID == "" {
			return fmt.Errorf("template_uuid is required for %s email type", emailType.String())
		}
		// Mailtrap OpenAPI: template emails forbid subject, text, html
		if msg.Subject != "" || msg.Text != "" || msg.HTML != "" {
			return fmt.Errorf("subject, text, and html are forbidden for %s email type", emailType.String())
		}
	}

	return nil
}
