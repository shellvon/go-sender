package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	emailjsDefaultEndpoint = "api.emailjs.com"
	emailjsDefaultPath     = "/api/v1.0/email/send"
)

// emailJSTransformer implements HTTPRequestTransformer for EmailJS.
type emailJSTransformer struct{}

// newEmailJSTransformer creates a new EmailJS transformer.
func newEmailJSTransformer() core.HTTPTransformer[*core.Account] {
	return &emailJSTransformer{}
}

// CanTransform checks if this transformer can handle the given message.
func (t *emailJSTransformer) CanTransform(msg core.Message) bool {
	emailMsg, ok := msg.(*Message)
	if !ok {
		return false
	}
	return emailMsg.SubProvider == string(SubProviderEmailJS)
}

// Transform converts an EmailJS message to HTTP request specification.
func (t *emailJSTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *core.Account,
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
	if len(msg.To) == 0 {
		return errors.New("at least one recipient is required")
	}
	if msg.From == "" {
		return errors.New("from is required for EmailJS")
	}
	return nil
}

// transformEmail transforms email message to HTTP request.
func (t *emailJSTransformer) transformEmail(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// Get required parameters
	serviceID := msg.From
	userID := msg.GetExtraStringOrDefault(emailjsUserID, account.Key)
	accessToken := msg.GetExtraStringOrDefault(emailjsAccessToken, account.Secret)

	if serviceID == "" {
		return nil, nil, errors.New("EmailJS: service_id (From) is required")
	}
	if userID == "" {
		return nil, nil, errors.New("EmailJS: user_id (APIKey) is required")
	}
	if accessToken == "" {
		return nil, nil, errors.New("EmailJS: accessToken (APISecret) is required")
	}

	// Build request parameters
	params := map[string]interface{}{
		emailjsServiceID:   serviceID,
		emailjsTemplateID:  msg.TemplateID,
		emailjsUserID:      userID,
		emailjsAccessToken: accessToken,
	}

	// Prepare template data with smart field merging
	templateData := t.prepareTemplateData(msg)

	// Add template parameters
	params[emailjsTemplateParams] = templateData

	// Convert to JSON
	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal EmailJS request body: %w", err)
	}

	// Build headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Build endpoint URL
	endpoint := t.getEndpoint(account)

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      endpoint,
		Headers:  headers,
		Body:     bodyData,
		BodyType: "json",
	}, t.handleEmailJSResponse, nil
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

// getEndpoint returns the appropriate endpoint URL.
func (t *emailJSTransformer) getEndpoint(account *core.Account) string {
	switch {
	case account.Endpoint != "":
		return "https://" + account.Endpoint + emailjsDefaultPath
	case account.IntlEndpoint != "":
		return "https://" + account.IntlEndpoint + emailjsDefaultPath
	default:
		return "https://" + emailjsDefaultEndpoint + emailjsDefaultPath
	}
}

// handleEmailJSResponse handles EmailJS API response.
func (t *emailJSTransformer) handleEmailJSResponse(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	// EmailJS returns "OK" on success
	if strings.Contains(string(body), "OK") {
		return nil
	}

	return fmt.Errorf("EmailJS send failed: %s", string(body))
}
