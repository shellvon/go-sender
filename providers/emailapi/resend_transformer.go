package emailapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// resendTransformer implements HTTPRequestTransformer for Resend.
// To use Resend, you need to create an API key in the Resend dashboard.
// The API key must be sent in the Authorization header as 'Bearer <API_KEY>'.
// This provider supports batch sending (up to 50 emails per request).
// Reference:
//   - Official Website: https://resend.com/
//   - API Docs: https://resend.com/docs/api-reference/emails

func init() {
	RegisterTransformer(string(SubProviderResend), newResendTransformer())
}

const (
	resendDefaultAPIPath  = "https://api.resend.com/emails"
	maxRecipientsPerBatch = 50
)

// resendTransformer implements Resend logic via BaseHTTPTransformer.
type resendTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func newResendTransformer() *resendTransformer {
	rt := &resendTransformer{}
	rt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeEmailAPI,
		string(SubProviderResend),
		&core.ResponseHandlerConfig{},
		rt.transform,
		transformer.AddBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			return rt.validate(msg, account)
		}),
	)
	return rt
}

// transform handles Resend message → HTTPRequestSpec.
func (rt *resendTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
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

	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Resend request body: %w", err)
	}

	headers := map[string]string{
		"Authorization":   "Bearer " + account.APIKey,
		"Idempotency-Key": msg.MsgID(),
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      resendDefaultAPIPath,
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

func (rt *resendTransformer) validate(msg *Message, account *Account) error {
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
	if account.APIKey == "" {
		return errors.New("APIKey is required for Resend")
	}
	return nil
}
