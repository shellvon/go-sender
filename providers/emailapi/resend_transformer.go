package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Resend
// @Website: https://resend.com/
// @APIDoc: https://resend.com/docs/api-reference/emails
//
// To use Resend, you need to create an API key in the Resend dashboard.
// The API key must be sent in the Authorization header as 'Bearer <API_KEY>'.
//
// This provider supports batch sending (up to 50 emails per request).

// init automatically registers the Resend transformer.
func init() {
	RegisterTransformer(string(SubProviderResend), newResendTransformer())
}

const (
	resendDefaultAPIPath  = "https://api.resend.com/emails"
	maxRecipientsPerBatch = 50
)

// resendTransformer implements HTTPRequestTransformer for Resend.
type resendTransformer struct{}

// newResendTransformer creates a new Resend transformer.
func newResendTransformer() core.HTTPTransformer[*Account] {
	return &resendTransformer{}
}

// CanTransform checks if this transformer can handle the given message.
func (t *resendTransformer) CanTransform(message core.Message) bool {
	if emailMsg, ok := message.(*Message); ok {
		return ok && emailMsg.SubProvider == string(SubProviderResend)
	}
	return false
}

// Transform converts a Resend message to HTTP request specification.
func (t *resendTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	emailMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for Resend: %T", msg)
	}

	if err := t.validateMessage(emailMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}

	return t.transformEmail(ctx, emailMsg, account)
}

// validateMessage validates the message for Resend.
func (t *resendTransformer) validateMessage(msg *Message) error {
	if len(msg.To) == 0 {
		return errors.New("to recipients cannot be empty")
	}
	if len(msg.To) > maxRecipientsPerBatch {
		return fmt.Errorf("to recipients are limited to %d recipients", maxRecipientsPerBatch)
	}
	if msg.From == "" {
		return errors.New("from is required for Resend")
	}
	if msg.Subject == "" {
		return errors.New("subject is required for Resend")
	}
	if msg.HTML == "" && msg.Text == "" {
		return errors.New("at least one of HTML or Text must be provided for Resend")
	}
	return nil
}

// transformEmail transforms email message to HTTP request.
func (t *resendTransformer) transformEmail(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// Build request parameters
	params := map[string]interface{}{
		"from":        msg.From,
		"to":          msg.To,
		"subject":     msg.Subject,
		"html":        msg.HTML,
		"text":        msg.Text,
		"reply_to":    msg.ReplyTo,
		"cc":          msg.Cc,
		"bcc":         msg.Bcc,
		"headers":     msg.Headers,
		"attachments": msg.Attachments,
		"tags":        msg.Extras[resendTagsKey],
	}

	if msg.ScheduledAt != nil {
		params["scheduled_at"] = msg.ScheduledAt.Format(time.RFC3339)
	}

	// Convert to JSON
	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Resend request body: %w", err)
	}

	// Build headers
	headers := map[string]string{
		"Authorization":   "Bearer " + account.APISecret,
		"Idempotency-Key": msg.MsgID(),
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      resendDefaultAPIPath,
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, t.handleResendResponse, nil
}

// handleResendResponse handles Resend API response.
func (t *resendTransformer) handleResendResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	// Parse response to check for data field
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("failed to parse Resend response: %w", err)
	}

	// Check if data field exists and is not empty
	if len(resp.Data) == 0 {
		return errors.New("resend response contains no data")
	}

	return nil
}
