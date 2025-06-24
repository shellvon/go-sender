package sms

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// API endpoint paths
const (
	// 国际短信-模版单发
	intlTemplateSingle = "/internationalsms/xsend" // https://www.mysubmail.com/documents/87QTB2
	// 国际短信-批量群发
	intlBatch = "internationalsms/batchsend" // https://www.mysubmail.com/documents/yD46O
	// 国际短信-单发
	intlSingle = "/internationalsms/send" // https://www.mysubmail.com/documents/3UQA3
	// 估计短信-模版一对多(没有找到批量的API)
	intlTemplateBatch = "internationalsms/multixsend" // https://www.mysubmail.com/documents/B70hy
	// 国内短信-模版单发
	domesticTemplateSingle = "/sms/xsend" // https://www.mysubmail.com/documents/OOVyh
	// 国内短信-模版-群发
	domesticTemplateBatch = "/sms/multixsend" // https://www.mysubmail.com/documents/G5KBR
	// 国内短信-单发
	domesticSingle = "/sms/send" // https://www.mysubmail.com/documents/FppOR3
	// 国内短信-群发
	domesticBatch = "/sms/multisend" // https://www.mysubmail.com/documents/AzD4Z4
)

// SubmailResponse represents the common response structure from SUBMAIL API
// API Documentation: https://www.mysubmail.com/documents/FppOR3
type SubmailResponse struct {
	Status string `json:"status"`  // Request status: success/error
	SendID string `json:"send_id"` // Unique send ID
	Fee    int    `json:"fee"`     // Billing count
	Code   string `json:"code"`    // Error code (when status is error)
	Msg    string `json:"msg"`     // Error message (when status is error)
}

// isSubmailInternational checks if the mobile number is international (not China)
func isSubmailInternational(mobile string) bool {
	cleanMobile := mobile
	if strings.HasPrefix(cleanMobile, "+") {
		cleanMobile = cleanMobile[1:]
	} else if strings.HasPrefix(cleanMobile, "00") {
		cleanMobile = cleanMobile[2:]
	}
	// +86/86/1xxxxxxxxxx 都视为国内
	if strings.HasPrefix(cleanMobile, "86") && len(cleanMobile) == 13 && cleanMobile[2] == '1' {
		return false
	}
	if len(cleanMobile) == 11 && cleanMobile[0] == '1' {
		return false
	}
	return strings.HasPrefix(mobile, "+") || strings.HasPrefix(mobile, "00")
}

// isSubmailIntlSMS checks if any of the mobile numbers are international
func isSubmailIntlSMS(mobiles []string) bool {
	for _, mobile := range mobiles {
		if isSubmailInternational(mobile) {
			return true
		}
	}
	return false
}

// buildSubmailEndpoint builds the API endpoint based on message type and recipient number
// API Documentation: https://www.mysubmail.com/documents/FppOR3
func buildSubmailEndpoint(provider *SMSProvider, msg *Message) string {
	baseDomain := "api-v4.mysubmail.com"
	if provider.Endpoint != "" {
		baseDomain = provider.Endpoint
	}

	isIntl := isSubmailIntlSMS(msg.Mobiles)
	isTemplate := msg.TemplateCode != ""

	// Path mappings based on provider type, template usage, and recipient count
	pathMap := map[bool]map[bool]string{
		true: { // International
			true:  intlTemplateSingle, // Template, single
			false: intlSingle,         // Non-template, single
		},
		false: { // Domestic
			true:  domesticTemplateSingle, // Template, single
			false: domesticSingle,         // Non-template, single
		},
	}

	// Update for batch send if multiple recipients
	if len(msg.Mobiles) > 1 {
		pathMap[true][true] = intlTemplateBatch      // International, template, batch
		pathMap[true][false] = intlBatch             // International, non-template, batch
		pathMap[false][true] = domesticTemplateBatch // Domestic, template, batch
		pathMap[false][false] = domesticBatch        // Domestic, non-template, batch
	}

	return fmt.Sprintf("https://%s%s", baseDomain, pathMap[isIntl][isTemplate])
}

// buildSubmailParams builds the request parameters based on message type
// API Documentation: https://www.mysubmail.com/documents/FppOR3
func buildSubmailParams(ctx context.Context, endpoint string, provider *SMSProvider, msg *Message) map[string]string {
	params := map[string]string{
		"appid": provider.AppID, // App ID from SUBMAIL dashboard
	}

	isIntl := strings.Contains(endpoint, "/internationalsms/")
	isIntlTemplateBatch := strings.Contains(endpoint, intlTemplateBatch)

	// 国际短信 sender 字段（可选）
	if isIntl && provider.SignName != "" {
		params["sender"] = provider.SignName
	}

	if isIntlTemplateBatch {
		// 国际批量模版发送 multi 参数
		multi := make([]map[string]interface{}, 0, len(msg.Mobiles))
		for _, mobile := range msg.Mobiles {
			item := map[string]interface{}{
				"to": mobile,
			}
			if len(msg.TemplateParams) > 0 {
				item["vars"] = msg.TemplateParams
			}
			multi = append(multi, item)
		}
		multiJSON, _ := json.Marshal(multi)
		params["multi"] = string(multiJSON)
	} else {
		params["to"] = strings.Join(msg.Mobiles, ",")
	}

	if msg.TemplateCode != "" {
		params["project"] = msg.TemplateCode // Template project ID from SUBMAIL dashboard
		if !isIntlTemplateBatch && len(msg.TemplateParams) > 0 {
			params["vars"] = toJSONString(msg.TemplateParams)
		}
	} else {
		// 非模版短信
		content := msg.Content
		if msg.SignName != "" && !strings.Contains(content, msg.SignName) {
			content = "【" + msg.SignName + "】" + content
		}
		params["content"] = content
	}

	// Add tag if provided in metadata
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if tag, ok := metadata["tag"].(string); ok && tag != "" {
			params["tag"] = tag // Custom tag for tracking (32 chars max)
		}
	}

	// Add timestamp for signature calculation
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	params["timestamp"] = timestamp

	// Add signature
	params["signature"] = calculateSubmailSignature(provider, params)

	return params
}

// calculateSubmailSignature calculates the signature for SUBMAIL API
// API Documentation: https://www.mysubmail.com/documents/VBcbe
func calculateSubmailSignature(provider *SMSProvider, params map[string]string) string {
	// Get sign type from provider config, default to "md5"
	signType := "md5"
	if provider.Channel != "" {
		signType = provider.Channel
	}

	// Build string to sign
	var keys []string
	for k := range params {
		if k != "signature" && k != "sign_type" && k != "sign_version" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, k+"="+params[k])
	}
	stringToSign := strings.Join(pairs, "&")

	// Calculate signature based on sign type
	switch signType {
	case "sha1":
		hash := sha1.Sum([]byte(stringToSign + provider.AppSecret))
		return fmt.Sprintf("%x", hash)
	case "normal":
		return provider.AppSecret
	default: // md5
		hash := md5.Sum([]byte(stringToSign + provider.AppSecret))
		return fmt.Sprintf("%x", hash)
	}
}

// sendSubmailSMS sends SMS via SUBMAIL API (domestic or international)
// API Documentation: https://www.mysubmail.com/documents/FppOR3
func sendSubmailSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	// 国际短信单次最大上限1000，国内10000，具体可根据业务需要加限制
	endpoint := buildSubmailEndpoint(provider, msg)
	params := buildSubmailParams(ctx, endpoint, provider, msg)

	return sendSubmailRequest(ctx, endpoint, params, ProviderTypeSubmail)
}

// sendSubmailRequest sends the actual HTTP request to SUBMAIL API
func sendSubmailRequest(ctx context.Context, endpoint string, params map[string]string, providerType ProviderType) error {
	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Raw: []byte(buildFormData(params)),
	})
	if err != nil {
		return fmt.Errorf("submail SMS request failed: %w", err)
	}

	var result SubmailResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse submail response: %w", err)
	}

	if result.Status != "success" {
		return &SMSError{
			Code:     result.Code,
			Message:  result.Msg,
			Provider: string(providerType),
		}
	}

	return nil
}
