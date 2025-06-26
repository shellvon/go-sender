package sms

// @ProviderName: Aliyun / 阿里云
// @Website: https://www.aliyun.com
// @APIDoc: https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
//
// # AliyunProvider implements SMSProviderInterface for Aliyun SMS
//
// 官方文档:
// - 短信模板即具体发送的短信内容，模板类型支持验证码、通知短信和推广短信。模板由模板变量和模板内容构成，您需要遵守模板内容规范和变量规范。
// - SMS API(国内/国外/单发/群发): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
//
// 阿里云支持多媒体彩信:
//   - 多媒体彩信: https://help.aliyun.com/zh/sms/user-guide/what-is-multimedia-sms
//   - 定价: 卡片短信默认定价是0.2元/条,数字短信默认定价为0.4元/条
//     https://help.aliyun.com/zh/sms/user-guide/multimedia-sms-pricing
//
// 阿里云支持语音短信:
//   - 语音API文档: https://help.aliyun.com/zh/dyvms/user-guide/voice-notification-overview
//   - 语音验证码API: https://help.aliyun.com/zh/dyvms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - 语音通知API: https://help.aliyun.com/zh/dyvms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
//   - 语音短信支持国内单发，支持验证码和通知类型，需开通语音服务。
// 对于目前，使用语音发送短信时，当发送验证码时，会使用 TTS 接口，当发送通知时，会使用 Voice 接口。
import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/shellvon/go-sender/utils"
)

type AliyunProvider struct {
	config SMSProvider
}

// NewAliyunProvider creates a new Aliyun SMS provider
func NewAliyunProvider(config SMSProvider) *AliyunProvider {
	return &AliyunProvider{
		config: config,
	}
}

// Send sends an SMS message via Aliyun
func (provider *AliyunProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	switch msg.Type {
	case MMS:
		return provider.sendCardSMS(ctx, msg)
	case Voice:
		return provider.sendVoice(ctx, msg)
	default:
		return provider.sendSMS(ctx, msg)
	}
}

// formatPhoneNumber 格式化手机号，添加国家代码前缀 适配阿里云的短信发送API
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
func formatPhoneNumber(mobile string, regionCode int) string {
	if regionCode == 0 || regionCode == 86 {
		return mobile
	}
	// 国际/港澳台消息：国际区号+号码，例如 852000012****。
	return fmt.Sprintf("%d%s", regionCode, mobile)
}

// GetCapabilities returns Aliyun's capabilities
func (p *AliyunProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 阿里云 SMS 能力配置
	capabilities.SMS.International = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification},
		"支持国际验证码和通知短信",
	)
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内验证码、通知和营销短信",
	)
	capabilities.SMS.Limits.MaxBatchSize = 1000
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "5000 QPS"
	capabilities.SMS.Limits.DailyLimit = ""
	// 阿里云支持彩信（卡片短信）
	capabilities.MMS.International = NewRegionCapability(
		false, false,
		[]MessageType{},
		[]MessageCategory{},
		"暂不支持国际彩信",
	)
	capabilities.MMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{MMS},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内彩信（卡片短信），需要申请开通",
	)
	capabilities.MMS.Limits.MaxBatchSize = 1000
	capabilities.MMS.Limits.MaxContentLen = 500
	capabilities.MMS.Limits.RateLimit = "1000 QPS"
	capabilities.MMS.Limits.DailyLimit = ""
	// 阿里云语音能力
	capabilities.Voice.Domestic = NewRegionCapability(
		true, false, // 只支持国内，且不支持群呼
		[]MessageType{Voice},
		[]MessageCategory{CategoryVerification, CategoryNotification},
		"支持国内语音验证码和语音通知，仅支持单号码",
	)
	capabilities.Voice.International = NewRegionCapability(
		false, false,
		[]MessageType{},
		[]MessageCategory{},
		"暂不支持国际语音",
	)
	return capabilities
}

// CheckCapability checks if a specific capability is supported
func (p *AliyunProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}

// GetLimits returns Aliyun's limits for a specific message type
func (p *AliyunProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	case Voice:
		return capabilities.Voice.GetLimits()
	case MMS:
		return capabilities.MMS.GetLimits()
	default:
		return capabilities.SMS.GetLimits() // 默认返回 SMS 的限制
	}
}

// GetName returns the provider name
func (p *AliyunProvider) GetName() string {
	return p.config.Name
}

// GetType returns the provider type
func (p *AliyunProvider) GetType() string {
	return string(p.config.Type)
}

// IsEnabled returns if the provider is enabled
func (p *AliyunProvider) IsEnabled() bool {
	return !p.config.Disabled
}

// GetWeight returns the provider weight
func (p *AliyunProvider) GetWeight() int {
	return p.config.GetWeight()
}

const (
	// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-endpoint
	aliyunDefaultSmsEndpoint   = "dysmsapi.aliyuncs.com"
	aliyunDefaultVoiceEndpoint = "dyvmsapi.aliyuncs.com"
)

type aliyunRequest struct {
	method      string
	host        string
	path        string
	query       map[string]string
	body        []byte
	contentType string
	action      string
	version     string
}

// doAliyunSMSRequest 统一处理阿里云短信/卡片/语音请求
func (provider *AliyunProvider) doAliyunSMSRequest(ctx context.Context, req aliyunRequest) error {
	headers := provider.signAliyunRequest(req)
	// 构建请求URL
	urlStr := "https://" + req.host + req.path
	q := url.Values{}
	keys := make([]string, 0, len(req.query))
	for k := range req.query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		q.Set(k, req.query[k])
	}
	urlStr += "?" + q.Encode()
	resp, _, err := utils.DoRequest(ctx, urlStr, utils.RequestOptions{
		Method:  req.method,
		Headers: headers,
		Raw:     req.body,
	})
	if err != nil {
		return fmt.Errorf("aliyun SMS request failed: %w", err)
	}
	// 解析响应
	var result struct {
		Message   string `json:"Message"`
		RequestId string `json:"RequestId"`
		BizId     string `json:"BizId"`
		Code      string `json:"Code"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse aliyun response: %w", err)
	}
	if result.Code != "OK" {
		return &SMSError{
			Code:     result.Code,
			Message:  result.Message,
			Provider: string(ProviderTypeAliyun),
		}
	}
	return nil
}

// signAliyunRequest 计算阿里云 V3 签名，返回完整请求头
func (provider *AliyunProvider) signAliyunRequest(req aliyunRequest) map[string]string {
	const algorithm = "ACS3-HMAC-SHA256"
	// 生成通用请求头
	xAcsSignatureNonce := fmt.Sprintf("%x", rand.New(rand.NewSource(time.Now().UnixNano())).Int63())
	xAcsDate := time.Now().UTC().Format(time.RFC3339)
	// 计算请求体哈希
	var hashedRequestPayload string
	if len(req.body) > 0 {
		hashedRequestPayload = utils.SHA256Hex(req.body)
	} else {
		hashedRequestPayload = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // 空字符串的SHA256
	}
	headers := map[string]string{
		"host":                  req.host,
		"x-acs-action":          req.action,
		"x-acs-version":         req.version,
		"x-acs-date":            xAcsDate,
		"x-acs-signature-nonce": xAcsSignatureNonce,
		"x-acs-content-sha256":  hashedRequestPayload,
	}
	if req.contentType != "" {
		headers["content-type"] = req.contentType
	}
	// 步骤 1：拼接规范请求串
	canonicalQueryString := ""
	keys := make([]string, 0, len(req.query))
	for k := range req.query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := req.query[k]
		canonicalQueryString += percentCode(url.QueryEscape(k)) + "=" + percentCode(url.QueryEscape(v)) + "&"
	}
	canonicalQueryString = strings.TrimSuffix(canonicalQueryString, "&")
	canonicalHeaders := ""
	signedHeaders := ""
	headerKeys := make([]string, 0, len(headers))
	for k := range headers {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)
	for _, k := range headerKeys {
		lowerKey := strings.ToLower(k)
		canonicalHeaders += lowerKey + ":" + headers[k] + "\n"
		signedHeaders += lowerKey + ";"
	}
	signedHeaders = strings.TrimSuffix(signedHeaders, ";")
	canonicalRequest := req.method + "\n" + req.path + "\n" + canonicalQueryString + "\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + hashedRequestPayload
	hashedCanonicalRequest := utils.SHA256Hex([]byte(canonicalRequest))
	stringToSign := algorithm + "\n" + hashedCanonicalRequest
	signature := hex.EncodeToString(utils.HMACSHA256([]byte(provider.config.AppSecret), []byte(stringToSign)))
	authorization := algorithm + " Credential=" + provider.config.AppID + ",SignedHeaders=" + signedHeaders + ",Signature=" + signature
	headers["Authorization"] = authorization
	return headers
}

// percentCode 按阿里云要求编码
func percentCode(str string) string {
	// 阿里云要求: 空格->%20, *->%2A, ~不变, +->%20
	encoded := url.QueryEscape(str)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

// sendCardSMS 发送卡片短信（彩信）
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (provider *AliyunProvider) sendCardSMS(ctx context.Context, msg *Message) error {
	if msg.IsIntl() {
		return NewUnsupportedInternationalError(string(ProviderTypeAliyun), "card SMS")
	}
	endpoint := provider.config.GetEndpoint(false, aliyunDefaultSmsEndpoint)
	cardObjects := make([]map[string]interface{}, 0, len(msg.Mobiles))
	for _, mobile := range msg.Mobiles {
		cardObject := map[string]interface{}{
			"mobile": mobile,
		}
		if len(msg.TemplateParams) > 0 {
			cardObject["dyncParams"] = utils.ToJSONString(msg.TemplateParams)
		}
		if msg.Extras != nil {
			if mobileParams, ok := msg.Extras["mobileParams"].(map[string]interface{}); ok {
				if mobileParam, exists := mobileParams[mobile]; exists {
					if paramMap, isMap := mobileParam.(map[string]string); isMap {
						cardObject["dyncParams"] = utils.ToJSONString(paramMap)
					}
				}
			}
		}
		if msg.Extras != nil {
			if mobileUrls, ok := msg.Extras["mobileUrls"].(map[string]interface{}); ok {
				if customUrl, exists := mobileUrls[mobile]; exists {
					if urlStr, isString := customUrl.(string); isString && urlStr != "" {
						cardObject["customUrl"] = urlStr
					}
				}
			}
		}
		if _, hasCustomUrl := cardObject["customUrl"]; !hasCustomUrl {
			if customUrl := msg.GetExtraStringOrDefault("customUrl", ""); customUrl != "" {
				cardObject["customUrl"] = customUrl
			}
		}
		cardObjects = append(cardObjects, cardObject)
	}
	params := map[string]string{
		"SignName":         msg.SignName,
		"CardTemplateCode": msg.TemplateID,
		"FallbackType":     msg.GetExtraStringOrDefault("FallbackType", "SMS"),
	}
	if smsTemplateCode := msg.GetExtraStringOrDefault("SmsTemplateCode", ""); smsTemplateCode != "" {
		params["SmsTemplateCode"] = smsTemplateCode
	}
	if len(msg.TemplateParams) > 0 {
		params["SmsTemplateParam"] = utils.ToJSONString(msg.TemplateParams)
	}
	if smsUpExtendCode := msg.GetExtraStringOrDefault("SmsUpExtendCode", ""); smsUpExtendCode != "" {
		params["SmsUpExtendCode"] = smsUpExtendCode
	}
	if templateCode := msg.GetExtraStringOrDefault("TemplateCode", ""); templateCode != "" {
		params["TemplateCode"] = templateCode
	}
	if len(msg.TemplateParams) > 0 {
		params["TemplateParam"] = utils.ToJSONString(msg.TemplateParams)
	}
	cardObjectsJSON, _ := json.Marshal(cardObjects)
	return provider.doAliyunSMSRequest(ctx, aliyunRequest{
		method:      "POST",
		host:        endpoint,
		path:        "/",
		query:       params,
		body:        cardObjectsJSON,
		contentType: "application/json",
		action:      "SendCardSms",
		version:     "2017-05-25",
	})
}

// sendVoice 发送语音（只支持国内，且只支持单号码）
// 根据消息类型自动选择：
// - 验证码类消息：使用TTS（文本转语音）API
// - 通知类消息：使用Voice（语音文件）API
func (provider *AliyunProvider) sendVoice(ctx context.Context, msg *Message) error {
	// 只支持国内
	if msg.IsIntl() {
		return NewUnsupportedInternationalError(string(ProviderTypeAliyun), "voice call")
	}
	// 只支持单号码
	if msg.HasMultipleRecipients() {
		return fmt.Errorf("aliyun voice only supports single number, got %d", len(msg.Mobiles))
	}

	params := map[string]string{
		"CalledNumber": msg.Mobiles[0],
	}

	// 根据消息类型选择API和参数
	isVerification := msg.Category == CategoryVerification
	if isVerification {
		// 验证码类消息：使用TTS API
		params["TtsCode"] = msg.TemplateID
		if len(msg.TemplateParams) > 0 {
			params["TtsParam"] = utils.ToJSONString(msg.TemplateParams)
		}
	} else {
		// 通知类消息：使用Voice API
		params["VoiceCode"] = msg.TemplateID
	}

	// 通用参数
	if v := msg.GetExtraStringOrDefault("CalledShowNumber", ""); v != "" {
		params["CalledShowNumber"] = v
	}
	if v := msg.GetExtraStringOrDefault("PlayTimes", "1"); v != "" {
		params["PlayTimes"] = v
	}
	if v := msg.GetExtraStringOrDefault("Volume", "100"); v != "" {
		params["Volume"] = v
	}
	if v := msg.GetExtraStringOrDefault("Speed", "0"); v != "" {
		params["Speed"] = v
	}
	if v := msg.GetExtraStringOrDefault("OutId", ""); v != "" {
		params["OutId"] = v
	}

	// 根据消息类型设置不同的action
	action := "SingleCallByVoice"
	if isVerification {
		action = "SingleCallByTts"
	}

	return provider.doAliyunSMSRequest(ctx, aliyunRequest{
		method:      "POST",
		host:        aliyunDefaultVoiceEndpoint,
		path:        "/",
		query:       params,
		body:        nil,
		contentType: "",
		action:      action,
		version:     "2017-05-25",
	})
}

// sendSMS 统一处理阿里云国内/国际/单发/群发短信
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
func (provider *AliyunProvider) sendSMS(ctx context.Context, msg *Message) error {
	endpoint := provider.config.GetEndpoint(msg.IsIntl(), aliyunDefaultSmsEndpoint)
	phones := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		phones[i] = formatPhoneNumber(mobile, msg.RegionCode)
	}
	params := map[string]string{
		"PhoneNumbers":  strings.Join(phones, ","),
		"SignName":      msg.SignName,
		"TemplateCode":  msg.TemplateID,
		"TemplateParam": utils.ToJSONString(msg.TemplateParams),
	}
	if v := msg.GetExtraStringOrDefault("SmsUpExtendCode", ""); v != "" {
		params["SmsUpExtendCode"] = v
	}
	if v := msg.GetExtraStringOrDefault("OutId", ""); v != "" {
		params["OutId"] = v
	}
	return provider.doAliyunSMSRequest(ctx, aliyunRequest{
		method:      "POST",
		host:        endpoint,
		path:        "/",
		query:       params,
		body:        nil,
		contentType: "",
		action:      "SendSms",
		version:     "2017-05-25",
	})
}
func (p *AliyunProvider) CheckConfigured() error {
	if p.config.AppID == "" || p.config.AppSecret == "" {
		return fmt.Errorf("aliyun SMS provider requires AppID and AppSecret")
	}
	return nil
}
