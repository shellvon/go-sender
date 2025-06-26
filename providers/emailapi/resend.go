package emailapi

// @ProviderName: resend
// @Website: https://resend.com/
// @APIDoc: https://resend.com/docs/api-reference/emails/send-batch-emails
//
// To use Resend, you need to create an API key in the Resend dashboard.
// The API key must be sent in the Authorization header as 'Bearer <API_KEY>'.
//
// This provider supports batch sending (up to 50 emails per request).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/utils"
)

const (
	ProviderResend = "resend"
	resendAPIURL   = "https://api.resend.com/emails/batch"
)

// ResendProvider implements APIClient for Resend batch email API.
type ResendProvider struct {
	account AccountConfig
}

var _ APIClient = (*ResendProvider)(nil)

func init() {
	RegisterProvider(ProviderResend, NewResendProvider)
}

func NewResendProvider(account AccountConfig) APIClient {
	return &ResendProvider{account: account}
}

// Send sends a batch of emails via Resend API.
func (p *ResendProvider) Send(ctx context.Context, msg *Message) error {
	if p.account.APIKey == "" {
		return errors.New("resend: APIKey is required in provider config")
	}
	if len(msg.To) == 0 {
		return errors.New("resend: To field is required (at least one recipient)")
	}
	if len(msg.To) > 50 {
		return errors.New("resend: To field is limited to 50 recipients")
	}
	if msg.From == "" {
		// To include a friendly name, use the format "Your Name <sender@domain.com>".
		return errors.New("resend: From field is required")
	}
	if msg.Subject == "" {
		return errors.New("resend: Subject is required")
	}
	if msg.HTML == "" && msg.Text == "" {
		return errors.New("resend: at least one of HTML or Text must be provided")
	}

	params := map[string]interface{}{
		"from":     utils.DefaultStringIfEmpty(msg.From, p.account.From),
		"to":       msg.To,
		"subject":  msg.Subject,
		"html":     msg.HTML,
		"text":     msg.Text,
		"reply_to": msg.ReplyTo,
		"cc":       msg.Cc,
		"bcc":      msg.Bcc,
		"headers":  msg.Headers,
	}

	options := utils.RequestOptions{
		Method: http.MethodPost,
		Headers: map[string]string{
			// To authenticate you need to add an Authorization header with the contents of the header being Bearer re_xxxxxxxxx where re_xxxxxxxxx is your API Key.
			"Authorization": "Bearer " + p.account.APIKey,
			"Content-Type":  "application/json",
			// Add an idempotency key to prevent duplicated emails.
			//  - Should be unique per API request
			//  - Idempotency keys expire after 24 hours
			//  - Have a maximum length of 256 characters
			"Idempotency-Key": msg.MsgID(),
		},
		JSON: params,
	}

	respBody, statusCode, err := utils.DoRequest(ctx, resendAPIURL, options)
	if err != nil {
		return fmt.Errorf("resend: request failed: %w", err)
	}
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("resend: API error: %s", string(respBody))
	}

	// Parse response to check for data field
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("resend: failed to parse response: %w", err)
	}
	if len(resp.Data) == 0 {
		return fmt.Errorf("resend: send failed, no data returned: %s", string(respBody))
	}
	return nil
}

func (p *ResendProvider) Name() string {
	return ProviderResend
}
