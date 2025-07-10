package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// emailJSTransformer implements HTTPRequestTransformer for EmailJS.
// You need to activate API requests through [Account:Security](https://dashboard.emailjs.com/admin/account/security)
//
// Reference:
//   - Official Website: https://www.emailjs.com/
//   - API Docs: https://www.emailjs.com/docs/rest-api/send/

func init() {
	RegisterTransformer(string(SubProviderEmailJS), newEmailJSTransformer())
}

const (
	defaultEmailjsAPIPath   = "https://api.emailjs.com/api/v1.0/email/send"
	defaultEmailjsServiceID = "default_service"
)

// emailJSTransformer utilises BaseHTTPTransformer for EmailJS.
type emailJSTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func newEmailJSTransformer() *emailJSTransformer {
	et := &emailJSTransformer{}

	cfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeText,
		CheckBody: true,
		Path:      "",
		Expect:    "OK",
		Mode:      core.MatchEq,
	}

	et.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeEmailAPI,
		string(SubProviderEmailJS),
		cfg,
		et.transform,
		transformer.AddBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			return et.validate(msg, account)
		}),
	)
	return et
}

// transform builds HTTPRequestSpec for EmailJS.
func (et *emailJSTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// Get required parameters
	serviceID := msg.From
	if serviceID == "" {
		serviceID = defaultEmailjsServiceID
	}
	userID := account.APIKey
	accessToken := account.APISecret
	templateData := et.prepareTemplateData(msg)

	params := map[string]interface{}{
		"service_id":      serviceID,
		"template_id":     msg.TemplateID,
		"user_id":         userID,
		"accessToken":     accessToken,
		"template_params": templateData,
	}

	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal EmailJS request body: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      defaultEmailjsAPIPath,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

func (et *emailJSTransformer) validate(msg *Message, account *Account) error {
	userID := account.APIKey
	accessToken := account.APISecret
	if userID == "" {
		return errors.New("EmailJS: user_id (APIKey) is required")
	}
	if accessToken == "" {
		return errors.New("EmailJS: accessToken (APISecret) is required")
	}
	if msg.TemplateID == "" {
		return errors.New("template_id is required for EmailJS")
	}
	return nil
}

// prepareTemplateData intelligently merges message fields with template data.
func (et *emailJSTransformer) prepareTemplateData(msg *Message) map[string]interface{} {
	// Start with user's template data if provided
	templateData := make(map[string]interface{})
	if msg.TemplateData != nil {
		for k, v := range msg.TemplateData {
			templateData[k] = v
		}
	}

	// Email field keys that EmailJS commonly expects
	emailFields := map[string]interface{}{
		"to":       msg.To,
		"cc":       msg.Cc,
		"bcc":      msg.Bcc,
		"subject":  msg.Subject,
		"from":     msg.From,
		"reply_to": msg.ReplyTo,
		"text":     msg.Text,
		"html":     msg.HTML,
	}

	// Only add fields that don't already exist in user's template data
	for fieldKey, fieldValue := range emailFields {
		// Skip empty values
		if isEmptyValue(fieldValue) {
			continue
		}

		// Check if this field already exists in user's template data
		if _, exists := templateData[fieldKey]; !exists {
			templateData[fieldKey] = fieldValue
		}
	}

	return templateData
}

// isEmptyValue checks if a value is considered empty for template data.
func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case []string:
		return len(val) == 0
	case []interface{}:
		return len(val) == 0
	default:
		return false
	}
}
