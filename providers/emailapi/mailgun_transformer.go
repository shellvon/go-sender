package emailapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// mailgunTransformer implements HTTPRequestTransformer for Mailgun.
// API Reference: https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun/messages/post-v3--domain-name--messages

func init() {
	RegisterTransformer(string(SubProviderMailgun), newMailgunTransformer())
}

const (
	mailgunMaxRecipients = 1000 // Per API documentation
	mailgunTagsKey       = "tags"
)

// mailgunTransformer implements Mailgun logic via BaseHTTPTransformer.
type mailgunTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func newMailgunTransformer() *mailgunTransformer {
	mt := &mailgunTransformer{}
	mt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeEmailAPI,
		string(SubProviderMailgun),
		&core.ResponseHandlerConfig{},
		mt.transform,
		transformer.AddBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			return mt.validate(msg, account)
		}),
	)
	return mt
}

// transform handles Mailgun message → HTTPRequestSpec.
func (mt *mailgunTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// Build request body as form data
	formData := mt.buildFormData(msg, account)

	// Construct API URL with domain
	// Check if domain is overridden in message extras first
	domain := account.Region
	if overrideDomain, ok := msg.Extras["domain"]; ok {
		if domainStr, ok := overrideDomain.(string); ok && domainStr != "" {
			domain = domainStr
		}
	}
	if domain == "" {
		return nil, nil, errors.New("domain is required for Mailgun (set in Region field or message extras)")
	}

	// Determine base URL based on domain or region
	baseURL := "https://api.mailgun.net"
	if strings.Contains(domain, ".eu") {
		baseURL = "https://api.eu.mailgun.net"
	}

	apiURL := fmt.Sprintf("%s/v3/%s/messages", baseURL, domain)

	// Basic auth with api key
	auth := "api:" + account.APIKey

	headers := map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(auth)),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      apiURL,
		Headers:  headers,
		Body:     encodeFormData(formData),
		BodyType: core.BodyTypeForm,
	}, nil, nil
}

// buildFormData constructs form data for Mailgun API.
func (mt *mailgunTransformer) buildFormData(msg *Message, account *Account) map[string]interface{} {
	data := make(map[string]interface{})

	// From (required)
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr != "" {
		data["from"] = fromAddr
	}

	// To (required)
	if len(msg.To) > 0 {
		data["to"] = strings.Join(msg.To, ",")
	}

	// CC
	if len(msg.Cc) > 0 {
		data["cc"] = strings.Join(msg.Cc, ",")
	}

	// BCC
	if len(msg.Bcc) > 0 {
		data["bcc"] = strings.Join(msg.Bcc, ",")
	}

	// Subject
	if msg.Subject != "" {
		data["subject"] = msg.Subject
	}

	// Content
	if msg.Text != "" {
		data["text"] = msg.Text
	}
	if msg.HTML != "" {
		data["html"] = msg.HTML
	}

	// Template support
	if msg.TemplateID != "" {
		data["template"] = msg.TemplateID

		// Template variables (t:variables) - for template substitution
		if msg.TemplateData != nil {
			// Mailgun expects t:variables as JSON string
			if jsonBytes, err := json.Marshal(msg.TemplateData); err == nil {
				data["t:variables"] = string(jsonBytes)
			}
		}
	} else {
		// If not using template, use TemplateData as custom variables (v:)
		if msg.TemplateData != nil {
			for key, value := range msg.TemplateData {
				data["v:"+key] = value
			}
		}
	}

	// Custom headers
	if len(msg.Headers) > 0 {
		for name, value := range msg.Headers {
			data["h:"+name] = value
		}
	}

	// Tags - Mailgun accepts multiple o:tag params
	if tags, ok := msg.Extras[mailgunTagsKey]; ok {
		if tagSlice, ok := tags.([]string); ok {
			// Store as slice for proper handling in form encoding
			data["o:tag"] = tagSlice
		}
	}

	// Scheduled delivery
	if msg.ScheduledAt != nil {
		// Mailgun expects RFC-2822 format
		data["o:deliverytime"] = msg.ScheduledAt.Format("Mon, 02 Jan 2006 15:04:05 -0700")
	}

	// DKIM
	if dkim, ok := msg.Extras["dkim"]; ok {
		data["o:dkim"] = dkim
	}

	// Tracking
	if tracking, ok := msg.Extras["tracking"]; ok {
		data["o:tracking"] = tracking
	}

	// Attachments
	for i, att := range msg.Attachments {
		if att.Disposition == "inline" {
			data[fmt.Sprintf("inline[%d]", i)] = att.Content
		} else {
			data[fmt.Sprintf("attachment[%d]", i)] = att.Content
		}
	}

	return data
}

// validate checks if the message and account are valid for Mailgun.
func (mt *mailgunTransformer) validate(msg *Message, account *Account) error {
	if account.APIKey == "" {
		return errors.New("APIKey is required for Mailgun")
	}

	// Check domain from account or message extras
	domain := account.Region
	if overrideDomain, ok := msg.Extras["domain"]; ok {
		if domainStr, ok := overrideDomain.(string); ok && domainStr != "" {
			domain = domainStr
		}
	}
	if domain == "" {
		return errors.New("domain is required for Mailgun (set in Region field or message extras)")
	}

	if len(msg.To) == 0 {
		return errors.New("to recipients cannot be empty")
	}

	if len(msg.To) > mailgunMaxRecipients {
		return fmt.Errorf("to recipients are limited to %d recipients", mailgunMaxRecipients)
	}

	// From field validation
	fromAddr := msg.From
	if fromAddr == "" && account.From != "" {
		fromAddr = account.From
	}
	if fromAddr == "" {
		return errors.New("from is required for Mailgun")
	}

	// Content validation
	if msg.Text == "" && msg.HTML == "" && msg.TemplateID == "" {
		return errors.New("either text, html content or template is required for Mailgun")
	}

	return nil
}

// Helper functions.
func encodeFormData(data map[string]interface{}) []byte {
	// Use standard library for URL encoding
	values := url.Values{}

	for key, value := range data {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case []byte:
			values.Add(key, string(v))
		case []string:
			// Handle multiple values for the same key (like o:tag)
			for _, item := range v {
				values.Add(key, item)
			}
		default:
			values.Add(key, fmt.Sprintf("%v", v))
		}
	}

	return []byte(values.Encode())
}
