package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
)

// @ProviderName: EmailJS
// @Website: https://www.emailjs.com/
// @APIDoc: https://www.emailjs.com/docs/rest-api/send/
//
// You need to activate API requests through [Account:Security](https://dashboard.emailjs.com/admin/account/security) for non-browser applications.

// init automatically registers the EmailJS transformer.
func init() {
	RegisterTransformer(string(SubProviderEmailJS), newEmailJSTransformer())
}

const (
	defaultEmailjsAPIPath   = "https://api.emailjs.com/api/v1.0/email/send"
	defaultEmailjsServiceID = "default_service"
)

// emailJSTransformer implements HTTPRequestTransformer for EmailJS.
type emailJSTransformer struct{}

// newEmailJSTransformer creates a new EmailJS transformer.
func newEmailJSTransformer() core.HTTPTransformer[*Account] {
	return &emailJSTransformer{}
}

// CanTransform checks if this transformer can handle the given message.
func (t *emailJSTransformer) CanTransform(message core.Message) bool {
	if emailMsg, ok := message.(*Message); ok {
		return ok && emailMsg.SubProvider == string(SubProviderEmailJS)
	}
	return false
}

// Transform converts an EmailJS message to HTTP request specification.
func (t *emailJSTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	emailMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for EmailJS: %T", msg)
	}

	if err := t.validateMessage(emailMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}

	return t.transformEmail(ctx, emailMsg, account)
}

// validateMessage validates the message for EmailJS.
func (t *emailJSTransformer) validateMessage(msg *Message) error {
	if msg.TemplateID == "" {
		return errors.New("template_id is required for EmailJS")
	}
	return nil
}

// transformEmail transforms email message to HTTP request.
func (t *emailJSTransformer) transformEmail(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// Get required parameters
	serviceID := msg.From
	userID := account.APIKey
	accessToken := account.APISecret

	if serviceID == "" {
		serviceID = defaultEmailjsServiceID
	}
	if userID == "" {
		return nil, nil, errors.New("EmailJS: user_id (APIKey) is required")
	}
	if accessToken == "" {
		return nil, nil, errors.New("EmailJS: accessToken (APISecret) is required")
	}

	// Prepare template data with smart field merging
	templateData := t.prepareTemplateData(msg)
	// Build request parameters
	params := map[string]interface{}{
		"service_id":      serviceID,
		"template_id":     msg.TemplateID,
		"user_id":         userID,
		"accessToken":     accessToken,
		"template_params": templateData,
	}

	// Convert to JSON
	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal EmailJS request body: %w", err)
	}
	return &core.HTTPRequestSpec{
			Method:   http.MethodPost,
			URL:      defaultEmailjsAPIPath,
			Body:     bodyData,
			BodyType: core.BodyTypeJSON,
		}, core.NewResponseHandler(&core.ResponseHandlerConfig{
			ResponseType:   core.BodyTypeText,
			SuccessPattern: "OK",
		}), nil
}

// prepareTemplateData intelligently merges message fields with template data
// EmailJS is special - recipients, sender, subject can come from template variables
// We prioritize user's template data, only adding message fields if they don't exist in template.
func (t *emailJSTransformer) prepareTemplateData(msg *Message) map[string]interface{} {
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
