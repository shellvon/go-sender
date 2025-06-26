package sms

// @ProviderName: Submail / 赛邮
// @Website: https://www.mysubmail.com/
// @APIDoc: https://www.mysubmail.com/documents
//
// 官方文档:
//   国内短信:
//     - 国内短信单发 https://www.mysubmail.com/documents/FppOR3
//     - 国内短信群发 https://www.mysubmail.com/documents/AzD4Z4
//     - 国内短信模版单发 https://www.mysubmail.com/documents/OOVyh
//     - 国内短信模版群发 https://www.mysubmail.com/documents/G5KBR
//   国际短信:
//     - 国际短信单发 https://www.mysubmail.com/documents/3UQA3
//     - 国际短信模版单发 https://www.mysubmail.com/documents/87QTB2
//     - 国际短信批量群发 https://www.mysubmail.com/documents/yD46O
//     - 国际短信模版一对多 https://www.mysubmail.com/documents/B70hy
//  彩信:
//     - 单发(不知道是否支持国际): https://www.mysubmail.com/documents/N6ktR
//  语音:
//     - 单发(不知道是否支持国际): https://www.mysubmail.com/documents/meE3C1
//     - 模版单发(不知道是否支持国际): https://www.mysubmail.com/documents/KbG03
//     - 模版群发(不知道是否支持国际): https://www.mysubmail.com/documents/FkgkM2

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

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
	// 国际短信-模版一对多(没有找到批量的API)
	intlTemplateBatch = "internationalsms/multixsend" // https://www.mysubmail.com/documents/B70hy
	// 国内短信-模版单发
	domesticTemplateSingle = "/sms/xsend" // https://www.mysubmail.com/documents/OOVyh
	// 国内短信-模版-群发
	domesticTemplateBatch = "/sms/multixsend" // https://www.mysubmail.com/documents/G5KBR
	// 国内短信-单发
	domesticSingle = "/sms/send" // https://www.mysubmail.com/documents/FppOR3
	// 国内短信-群发
	domesticBatch = "/sms/multisend" // https://www.mysubmail.com/documents/AzD4Z4
	// 彩信-单发
	mmsSingle = "/mms/send" // https://www.mysubmail.com/documents/N6ktR
	// 语音-单发
	voiceSingle = "/voice/send" // https://www.mysubmail.com/documents/meE3C1
	// 语音-模版单发
	voiceTemplateSingle = "/voice/xsend" // https://www.mysubmail.com/documents/KbG03
	// 语音-模版群发
	voiceTemplateBatch   = "/voice/multixsend" // https://www.mysubmail.com/documents/FkgkM2
	submailDefaultDomain = "api-v4.mysubmail.com"
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

// buildSubmailEndpoint builds the API endpoint based on message type and recipient number
func buildSubmailEndpoint(provider *SMSProvider, msg *Message) string {
	baseDomain := provider.GetEndpoint(msg.IsIntl(), submailDefaultDomain)
	isIntl := msg.IsIntl()
	isTemplate := msg.TemplateID != ""
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
	if isIntl && msg.SignName != "" {
		params["sender"] = msg.SignName
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
	if msg.TemplateID != "" {
		params["project"] = msg.TemplateID // Template project ID from SUBMAIL dashboard
		if !isIntlTemplateBatch && len(msg.TemplateParams) > 0 {
			params["vars"] = utils.ToJSONString(msg.TemplateParams)
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
	if tag, ok := msg.GetExtraString("tag"); ok && tag != "" {
		params["tag"] = tag // Custom tag for tracking (32 chars max)
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
		return utils.SHA1Hex(stringToSign + provider.AppSecret)
	case "normal":
		return provider.AppSecret
	default: // md5
		return utils.MD5Hex(stringToSign + provider.AppSecret)
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
		Data: params,
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

type SubmailProvider struct {
	config SMSProvider
}

func NewSubmailProvider(config SMSProvider) *SubmailProvider {
	return &SubmailProvider{config: config}
}
func (provider *SubmailProvider) Send(ctx context.Context, msg *Message) error {
	if provider.config.AppID == "" || provider.config.AppSecret == "" {
		return fmt.Errorf("submail provider requires AppID and AppSecret")
	}
	if err := msg.Validate(); err != nil {
		return err
	}
	if err := provider.CheckCapability(msg); err != nil {
		return err
	}

	switch msg.Type {
	case SMSText:
		return provider.sendSMS(ctx, msg)
	case Voice:
		return provider.sendVoice(ctx, msg)
	case MMS:
		return provider.sendMMS(ctx, msg)
	default:
		return NewUnsupportedMessageTypeError(string(ProviderTypeSubmail), msg.Type.String(), msg.Category.String())
	}
}
func (provider *SubmailProvider) sendSMS(ctx context.Context, msg *Message) error {
	endpoint := buildSubmailEndpoint(&provider.config, msg)
	params := buildSubmailParams(ctx, endpoint, &provider.config, msg)
	return sendSubmailRequest(ctx, endpoint, params, ProviderTypeSubmail)
}

// sendVoice sends voice message via SUBMAIL API
// API Documentation: https://www.mysubmail.com/documents/meE3C1
func (provider *SubmailProvider) sendVoice(ctx context.Context, msg *Message) error {
	baseDomain := provider.config.GetEndpoint(false, submailDefaultDomain)

	var apiPath string
	if msg.TemplateID != "" {
		if len(msg.Mobiles) > 1 {
			apiPath = voiceTemplateBatch // /voice/multixsend
		} else {
			apiPath = voiceTemplateSingle // /voice/xsend
		}
	} else {
		apiPath = voiceSingle // /voice/send
	}

	endpoint := fmt.Sprintf("https://%s%s", baseDomain, apiPath)
	var params map[string]string
	if apiPath == voiceTemplateBatch {
		// 构建multi参数
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
		params = map[string]string{
			"appid":     provider.config.AppID,
			"multi":     string(multiJSON),
			"project":   msg.TemplateID,
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		}
		if tag, ok := msg.GetExtraString("tag"); ok && tag != "" {
			params["tag"] = tag
		}
		params["signature"] = calculateSubmailSignature(&provider.config, params)
	} else {
		params = buildSubmailVoiceParams(ctx, &provider.config, msg)
	}
	return sendSubmailRequest(ctx, endpoint, params, ProviderTypeSubmail)
}

// sendMMS sends MMS via SUBMAIL API
// API Documentation: https://www.mysubmail.com/documents/N6ktR
func (provider *SubmailProvider) sendMMS(ctx context.Context, msg *Message) error {
	// 彩信仅支持单发
	if len(msg.Mobiles) > 1 {
		return fmt.Errorf("submail MMS only supports single send")
	}

	baseDomain := provider.config.GetEndpoint(false, submailDefaultDomain)

	endpoint := fmt.Sprintf("https://%s%s", baseDomain, mmsSingle)
	params := buildSubmailMMSParams(ctx, &provider.config, msg)
	return sendSubmailRequest(ctx, endpoint, params, ProviderTypeSubmail)
}

func (p *SubmailProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，单发/群发，支持模板和非模板，单次最多10000条",
	)
	// 国际短信支持单发/群发
	capabilities.SMS.International = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际短信，单发/群发，支持模板和非模板，单次最多1000条",
	)
	capabilities.SMS.Limits.MaxBatchSize = 10000
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "未知"
	capabilities.SMS.Limits.DailyLimit = "未知"

	// 彩信支持单发
	capabilities.MMS.Domestic = NewRegionCapability(
		true, false,
		[]MessageType{MMS},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内彩信，仅单发，支持模板",
	)
	capabilities.MMS.International = NewRegionCapability(
		false, false, nil, nil, "不支持国际彩信",
	)
	capabilities.MMS.Limits.MaxBatchSize = 1
	capabilities.MMS.Limits.MaxContentLen = 500
	capabilities.MMS.Limits.RateLimit = "未知"
	capabilities.MMS.Limits.DailyLimit = "未知"

	// 语音支持单发
	capabilities.Voice.Domestic = NewRegionCapability(
		true, false,
		[]MessageType{Voice},
		[]MessageCategory{CategoryVerification, CategoryNotification},
		"支持国内语音，仅单发，支持模板和非模板",
	)
	capabilities.Voice.International = NewRegionCapability(
		false, false, nil, nil, "不支持国际语音",
	)
	capabilities.Voice.Limits.MaxBatchSize = 1
	capabilities.Voice.Limits.MaxContentLen = 500
	capabilities.Voice.Limits.RateLimit = "未知"
	capabilities.Voice.Limits.DailyLimit = "未知"

	return capabilities
}

// CheckCapability checks if a specific capability is supported
func (p *SubmailProvider) CheckCapability(msg *Message) error {
	// 语音和彩信只允许单发
	if (msg.Type == Voice || msg.Type == MMS) && msg.HasMultipleRecipients() {
		return fmt.Errorf("submail %s only supports single send", msg.Type.String())
	}
	return DefaultCheckCapability(p, msg)
}

func (p *SubmailProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	case Voice:
		return capabilities.Voice.GetLimits()
	case MMS:
		return capabilities.MMS.GetLimits()
	default:
		return Limits{}
	}
}

func (p *SubmailProvider) GetName() string {
	return p.config.Name
}
func (p *SubmailProvider) GetType() string {
	return string(p.config.Type)
}
func (p *SubmailProvider) IsEnabled() bool {
	return !p.config.Disabled
}
func (p *SubmailProvider) GetWeight() int {
	return p.config.GetWeight()
}
func (p *SubmailProvider) CheckConfigured() error {
	if p.config.AppID == "" || p.config.AppSecret == "" {
		return fmt.Errorf("submail provider requires AppID and AppSecret")
	}
	return nil
}

// buildSubmailVoiceParams builds voice request parameters
func buildSubmailVoiceParams(ctx context.Context, provider *SMSProvider, msg *Message) map[string]string {
	params := map[string]string{
		"appid": provider.AppID,
		"to":    msg.Mobiles[0], // 语音仅支持单发
	}

	if msg.TemplateID != "" {
		// 模板语音
		params["project"] = msg.TemplateID
		if len(msg.TemplateParams) > 0 {
			params["vars"] = utils.ToJSONString(msg.TemplateParams)
		}
	} else {
		// 非模板语音
		params["content"] = msg.Content
	}

	// Add tag if provided in metadata
	if tag, ok := msg.GetExtraString("tag"); ok && tag != "" {
		params["tag"] = tag
	}

	// Add timestamp for signature calculation
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	params["timestamp"] = timestamp
	// Add signature
	params["signature"] = calculateSubmailSignature(provider, params)
	return params
}

// buildSubmailMMSParams builds MMS request parameters
func buildSubmailMMSParams(ctx context.Context, provider *SMSProvider, msg *Message) map[string]string {
	params := map[string]string{
		"appid": provider.AppID,
		"to":    msg.Mobiles[0], // 彩信仅支持单发
	}

	// 彩信需要模板
	if msg.TemplateID != "" {
		params["project"] = msg.TemplateID
		if len(msg.TemplateParams) > 0 {
			params["vars"] = utils.ToJSONString(msg.TemplateParams)
		}
	} else {
		// 如果没有模板，使用content作为文本内容
		params["content"] = msg.Content
	}

	// Add tag if provided in metadata
	if tag, ok := msg.GetExtraString("tag"); ok && tag != "" {
		params["tag"] = tag
	}

	// Add timestamp for signature calculation
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	params["timestamp"] = timestamp
	// Add signature
	params["signature"] = calculateSubmailSignature(provider, params)
	return params
}
