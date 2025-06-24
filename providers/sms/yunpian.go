package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// YunpianResponse represents the common response structure from Yunpian API
// API Documentation: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
type YunpianResponse struct {
	Code int             `json:"code"`           // Response code, 0 means success
	Msg  string          `json:"msg"`            // Response message
	Data json.RawMessage `json:"data,omitempty"` // Response data
}

const (
	singleTemplatePath = "/v2/sms/tpl_single_send.json"
	batchTemplatePath  = "/v2/sms/tpl_batch_send.json"
	singlePath         = "/v2/sms/single_send.json"
	batchPath          = "/v2/sms/batch_send.json"
)

// buildYunpianEndpoint builds the API endpoint based on message type and provider configuration
// API Documentation: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
func buildYunpianEndpoint(provider *SMSProvider, msg *Message) string {
	// Determine domain - allow custom domain override via Endpoint field
	domain := "sms.yunpian.com" // Default domestic domain
	if provider.Endpoint != "" {
		domain = provider.Endpoint
	}

	pathMap := map[bool]string{
		// https://www.yunpian.com/official/document/sms/zh_CN/domestic_tpl_single_send
		true: singleTemplatePath, // Template, single
		// https://www.yunpian.com/official/document/sms/zh_CN/domestic_single_send
		false: singlePath, // Non-template, single
	}
	if len(msg.Mobiles) > 1 {
		// https://www.yunpian.com/official/document/sms/zh_CN/domestic_tpl_batch_send
		pathMap[true] = batchTemplatePath // Template, batch
		// https://www.yunpian.com/official/document/sms/zh_CN/domestic_batch_send
		pathMap[false] = batchPath // Non-template, batch
	}

	// Select path based on whether TemplateCode is set
	path := pathMap[msg.TemplateCode != ""]

	return fmt.Sprintf("https://%s%s", domain, path)
}

// buildYunpianParams builds the request parameters based on message type
// API Documentation: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
func buildYunpianParams(provider *SMSProvider, msg *Message, ctx context.Context) map[string]string {
	params := map[string]string{
		"apikey": provider.AppSecret,             // API key from Yunpian dashboard
		"mobile": strings.Join(msg.Mobiles, ","), // Mobile numbers (single or comma-separated)
	}

	isIntl := false
	for _, m := range msg.Mobiles {
		if utils.IsInternationalMobile(m) {
			isIntl = true
			break
		}
	}

	if msg.TemplateCode != "" {
		// Template SMS parameters
		params["tpl_id"] = msg.TemplateCode                          // Template ID from Yunpian dashboard
		params["tpl_value"] = buildTemplateValue(msg.TemplateParams) // Template variables in format #key#=value
	} else {
		// Non-template SMS parameters
		content := msg.Content
		// Add sign name only for domestic non-template SMS
		if !isIntl && msg.SignName != "" && !strings.Contains(content, msg.SignName) {
			content = "【" + msg.SignName + "】" + content
		}
		params["text"] = content // SMS content (signature only for domestic SMS)
	}

	// Add callback URL if configured
	if provider.Callback != "" {
		params["callback_url"] = provider.Callback // Callback URL for delivery status
	}

	// Add user ID if provided in metadata
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if uid, ok := metadata["uid"].(string); ok && uid != "" {
			params["uid"] = uid // User ID for tracking
		}
	}

	return params
}

// sendYunpianSMS sends SMS via Yunpian API (both domestic and international)
// Supports single send, batch send, template single send, and template batch send
// API Documentation: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
func sendYunpianSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	isIntl := false
	for _, m := range msg.Mobiles {
		if utils.IsInternationalMobile(m) {
			isIntl = true
			break
		}
	}
	// International SMS only supports single send
	if isIntl && len(msg.Mobiles) > 1 {
		return fmt.Errorf("yunpian international SMS only supports single send, got %d mobiles", len(msg.Mobiles))
	}

	endpoint := buildYunpianEndpoint(provider, msg)
	params := buildYunpianParams(provider, msg, ctx)

	return sendYunpianRequest(ctx, endpoint, params, provider.Type)
}

// sendYunpianRequest sends the actual HTTP request to Yunpian API
func sendYunpianRequest(ctx context.Context, endpoint string, params map[string]string, providerType ProviderType) error {
	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Raw: []byte(buildFormData(params)),
	})
	if err != nil {
		return fmt.Errorf("yunpian SMS request failed: %w", err)
	}

	var result YunpianResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse yunpian response: %w", err)
	}

	if result.Code != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.Code),
			Message:  result.Msg,
			Provider: string(providerType),
		}
	}

	return nil
}

// buildTemplateValue converts template parameters to Yunpian format with proper URL encoding
// Yunpian template format: urlencode("#key#") + "=" + urlencode("value") + "&" + urlencode("#key2#") + "=" + urlencode("value2")
// API Documentation: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
func buildTemplateValue(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var pairs []string
	for key, value := range params {
		// Format: urlencode("#key#") + "=" + urlencode("value")
		encodedKey := url.QueryEscape("#" + key + "#")
		encodedValue := url.QueryEscape(value)
		pairs = append(pairs, encodedKey+"="+encodedValue)
	}
	sort.Strings(pairs) // Sort for consistent ordering
	return strings.Join(pairs, "&")
}

// buildFormData builds URL-encoded form data from parameters
func buildFormData(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	return values.Encode()
}
