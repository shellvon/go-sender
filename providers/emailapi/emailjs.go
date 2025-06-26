package emailapi

// @ProviderName: emailjs
// @Website: https://www.emailjs.com/
// @APIDoc: https://www.emailjs.com/docs/rest-api/send/
//
// You need to activate API requests through [Account:Security](https://dashboard.emailjs.com/admin/account/security) for non-browser applications.
//

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

const (
	ProviderEmailJS     = "emailjs"
	emailjsAPIURL       = "https://api.emailjs.com/api/v1.0/email/send"
	paramServiceID      = "service_id"
	paramTemplateID     = "template_id"
	paramUserID         = "user_id"
	paramTemplateParams = "template_params"
	paramAccessToken    = "accessToken"
)

// EmailJSProvider implements APIClient for EmailJS email API.
type EmailJSProvider struct {
	account AccountConfig
}

var _ APIClient = (*EmailJSProvider)(nil)

func init() {
	RegisterProvider(ProviderEmailJS, NewEmailJSProvider)
}

func NewEmailJSProvider(account AccountConfig) APIClient {
	return &EmailJSProvider{account: account}
}

// Send sends an email via EmailJS API using provider config for static values and msg for dynamic values.
func (p *EmailJSProvider) Send(ctx context.Context, msg *Message) error {
	serviceID := msg.From
	if serviceID == "" {
		serviceID = msg.GetExtraStringOrDefault(paramServiceID, p.account.From)
	}
	userID := msg.GetExtraStringOrDefault(paramUserID, p.account.APIKey)
	accessToken := msg.GetExtraStringOrDefault(paramAccessToken, p.account.APISecret)

	if serviceID == "" {
		return errors.New("emailjs: service_id (From) is required in provider config or msg.Extras")
	}
	if userID == "" {
		return errors.New("emailjs: user_id (APIKey) is required in provider config or msg.Extras")
	}
	if accessToken != "" {
		return errors.New("emailjs: accessToken is required in provider config or msg.TemplateID")
	}
	if msg.TemplateID == "" {
		return errors.New("emailjs: template_id is required in msg.TemplateID")
	}

	if len(msg.To) != 0 {
		return errors.New("emailjs: to is not supported")
	}
	if len(msg.Attachments) != 0 {
		return errors.New("emailjs: attachments is not supported")
	}

	params := map[string]interface{}{
		paramServiceID:   serviceID,
		paramTemplateID:  msg.TemplateID,
		paramUserID:      userID,
		paramAccessToken: accessToken,
	}
	if msg.TemplateData != nil {
		params[paramTemplateParams] = msg.TemplateData
	}

	options := utils.RequestOptions{
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/json"},
		JSON:    params,
	}

	requestURI := emailjsAPIURL
	if p.account.Domain != "" {
		requestURI = p.account.Domain
	}
	respBody, statusCode, err := utils.DoRequest(ctx, requestURI, options)
	if err != nil {
		return fmt.Errorf("emailjs: request failed: %w", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("emailjs: API error: %s", string(respBody))
	}
	if strings.Contains(string(respBody), "OK") {
		return nil
	}
	return fmt.Errorf("emailjs: send failed: %s", string(respBody))
}

func (p *EmailJSProvider) Name() string {
	return ProviderEmailJS
}
