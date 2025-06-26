package sms

// @ProviderName: Tencent / 腾讯云
// @Website: https://cloud.tencent.com
// @APIDoc: https://cloud.tencent.com/document/product/382/55981
//
// # TencentProvider implements SMSProviderInterface for Tencent Cloud SMS
//
// 官方文档:
//   - 短信API文档: https://cloud.tencent.com/document/product/382/55981
//   - 国内短信: 支持验证码、通知类短信和营销短信
//   - 国际/港澳台短信: 支持验证码、通知类短信和营销短信
//   - 语音短信 验证码： https://cloud.tencent.com/document/product/1128/51559
//   -  语音短信-通知  https://cloud.tencent.com/document/product/1128/51558
//
// 能力说明:
//   - 国内短信：支持单发和群发，需模板ID，需签名。
//   - 国际/港澳台短信：支持单发和群发，需模板ID，签名可选。
//   - 彩信/语音：暂不支持。
//   - 彩信: WIP
//
// 注意：支持国内短信与国际/港澳台短信，默认接口请求频率限制：3000次/秒。
import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/shellvon/go-sender/utils"
)

const (
	tencentSMSEndpoint = "sms.tencentcloudapi.com"
	tencentSMSVersion  = "2021-01-11"
	tencentSMSAction   = "SendSms"
)

type TencentProvider struct {
	config SMSProvider
}

// NewTencentProvider creates a new Tencent Cloud SMS provider
func NewTencentProvider(config SMSProvider) *TencentProvider {
	return &TencentProvider{config: config}
}

// Send sends an SMS message via Tencent Cloud
func (provider *TencentProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	switch msg.Type {
	case MMS:
		return NewUnsupportedMessageTypeError(string(ProviderTypeTencent), msg.Type.String(), "")
	case Voice:
		return provider.sendVoice(ctx, msg)
	default:
		return provider.sendSMS(ctx, msg)
	}
}

// formatTencentPhone 格式化手机号，始终+开头，国内强制+86，国际为+regionCode
func formatTencentPhone(mobile string, regionCode int) string {
	if regionCode == 0 {
		regionCode = 86
	}
	return fmt.Sprintf("+%d%s", regionCode, mobile)
}

// buildTencentHeaders 构造腾讯云API请求头并签名
func buildTencentHeaders(endpoint, action, version, region, appSecret string, bodyData []byte, timestamp int64, credentialScope string, appSecretForCredential string) map[string]string {
	signature := calculateTencentSignature(appSecret, bodyData, timestamp, credentialScope)
	return map[string]string{
		"Content-Type":    "application/json",
		"Host":            endpoint,
		"X-TC-Action":     action,
		"X-TC-Version":    version,
		"X-TC-Timestamp":  fmt.Sprintf("%d", timestamp),
		"X-TC-Region":     region,
		"Authorization":   signature,
		"X-TC-Credential": fmt.Sprintf("%s/%s", appSecretForCredential, credentialScope),
	}
}

// calculateTencentSignature 计算腾讯云API签名（原calculateSignature逻辑）
func calculateTencentSignature(secretKey string, payload []byte, timestamp int64, credentialScope string) string {
	h := hmac.New(sha256.New, []byte("TC3"+secretKey))
	date := time.Unix(timestamp, 0).UTC().Format("20060102")
	h.Write([]byte(date))
	dateKey := h.Sum(nil)

	h = hmac.New(sha256.New, dateKey)
	h.Write([]byte("sms"))
	dateServiceKey := h.Sum(nil)

	h = hmac.New(sha256.New, dateServiceKey)
	h.Write([]byte("tc3_request"))
	signingKey := h.Sum(nil)

	h = hmac.New(sha256.New, signingKey)
	h.Write(payload)
	signature := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=content-type;host, Signature=%s",
		secretKey, credentialScope, signature)
}

// sendSMS sends SMS message via Tencent Cloud API
// API文档: https://cloud.tencent.com/document/product/382/55981
func (provider *TencentProvider) sendSMS(ctx context.Context, msg *Message) error {
	endpoint := provider.config.GetEndpoint(msg.IsIntl(), tencentSMSEndpoint)
	url := "https://" + endpoint

	// 格式化手机号
	phoneNumbers := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		phoneNumbers[i] = formatTencentPhone(mobile, msg.RegionCode)
	}

	params := map[string]interface{}{
		"PhoneNumberSet": phoneNumbers,
		"SmsSdkAppId":    provider.config.AppID,
		"TemplateId":     msg.TemplateID,
	}

	// 签名处理
	if msg.IsIntl() {
		if msg.SignName != "" {
			params["SignName"] = msg.SignName
		}
	} else {
		if msg.SignName == "" {
			return fmt.Errorf("tencent domestic SMS requires sign name")
		}
		params["SignName"] = msg.SignName
	}

	// 模板参数处理
	if len(msg.TemplateParams) > 0 {
		templateParams := make([]string, 0, len(msg.TemplateParams))
		if len(msg.ParamsOrder) > 0 {
			for _, key := range msg.ParamsOrder {
				if value, exists := msg.TemplateParams[key]; exists {
					templateParams = append(templateParams, value)
				}
			}
		} else {
			keys := make([]string, 0, len(msg.TemplateParams))
			for key := range msg.TemplateParams {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				templateParams = append(templateParams, msg.TemplateParams[key])
			}
		}
		params["TemplateParamSet"] = templateParams
	}

	if extendCode := msg.GetExtraStringOrDefault("ExtendCode", ""); extendCode != "" {
		params["ExtendCode"] = extendCode
	}
	if senderId := msg.GetExtraStringOrDefault("SenderId", ""); senderId != "" {
		params["SenderId"] = senderId
	}

	requestBody := map[string]interface{}{
		"Action":  tencentSMSAction,
		"Version": tencentSMSVersion,
		"Region":  msg.GetExtraStringOrDefault("Region", provider.config.Channel),
		"Request": params,
	}

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal tencent request body: %w", err)
	}

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/sms/tc3_request", date)

	headers := buildTencentHeaders(endpoint, tencentSMSAction, tencentSMSVersion, msg.GetExtraStringOrDefault("Region", "ap-guangzhou"), provider.config.AppSecret, bodyData, timestamp, credentialScope, provider.config.AppSecret)

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: headers,
		JSON:    requestBody,
	})
	if err != nil {
		return fmt.Errorf("tencent SMS request failed: %w", err)
	}
	return provider.parseResponse(resp)
}

// parseResponse 解析腾讯云API响应
func (provider *TencentProvider) parseResponse(resp []byte) error {
	var result struct {
		Response struct {
			SendStatusSet []struct {
				SerialNo       string `json:"SerialNo"`
				PhoneNumber    string `json:"PhoneNumber"`
				Fee            int    `json:"Fee"`
				SessionContext string `json:"SessionContext"`
				Code           string `json:"Code"`
				Message        string `json:"Message"`
				IsoCode        string `json:"IsoCode"`
			} `json:"SendStatusSet"`
			RequestId string `json:"RequestId"`
		} `json:"Response"`
		Error *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error,omitempty"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse tencent response: %w", err)
	}

	// 检查是否有错误
	if result.Error != nil {
		return &SMSError{
			Code:     result.Error.Code,
			Message:  result.Error.Message,
			Provider: string(ProviderTypeTencent),
		}
	}

	// 检查发送状态
	for _, status := range result.Response.SendStatusSet {
		if status.Code != "Ok" {
			return &SMSError{
				Code:     status.Code,
				Message:  status.Message,
				Provider: string(ProviderTypeTencent),
			}
		}
	}

	return nil
}

// sendVoice sends voice message via Tencent Cloud API
// 语音验证码API: https://cloud.tencent.com/document/product/1128/51559
// 语音通知API: https://cloud.tencent.com/document/product/1128/51558
func (provider *TencentProvider) sendVoice(ctx context.Context, msg *Message) error {
	var action, templateKey string
	const voiceEndpoint = "vms.tencentcloudapi.com"
	const voiceVersion = "2020-09-02"
	if msg.Category == CategoryVerification {
		action = "SendTtsVoice"
		templateKey = "TemplateId"
	} else {
		action = "SendVoice"
		templateKey = "VoiceId"
	}

	endpoint := provider.config.GetEndpoint(msg.IsIntl(), voiceEndpoint)
	url := "https://" + endpoint

	// 腾讯云语音只支持单发
	if len(msg.Mobiles) != 1 {
		return fmt.Errorf("tencent voice only supports single mobile per request")
	}
	calledNumber := formatTencentPhone(msg.Mobiles[0], msg.RegionCode)

	params := map[string]interface{}{
		"CalledNumber":  calledNumber,
		"VoiceSdkAppId": provider.config.AppID,
	}
	if msg.TemplateID != "" {
		params[templateKey] = msg.TemplateID
	}
	if len(msg.ParamsOrder) > 0 {
		params["TemplateParamSet"] = msg.ParamsOrder
	} else if len(msg.TemplateParams) > 0 {
		keys := make([]string, 0, len(msg.TemplateParams))
		for key := range msg.TemplateParams {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		templateParams := make([]string, 0, len(keys))
		for _, key := range keys {
			templateParams = append(templateParams, msg.TemplateParams[key])
		}
		params["TemplateParamSet"] = templateParams
	}
	if playTimes := msg.GetExtraStringOrDefault("PlayTimes", "2"); playTimes != "" {
		params["PlayTimes"] = playTimes
	}
	if sessionContext := msg.GetExtraStringOrDefault("SessionContext", ""); sessionContext != "" {
		params["SessionContext"] = sessionContext
	}
	if callerId := msg.GetExtraStringOrDefault("CallerId", ""); callerId != "" {
		params["CallerId"] = callerId
	}

	requestBody := map[string]interface{}{
		"Action":  action,
		"Version": voiceVersion,
		"Region":  msg.GetExtraStringOrDefault("Region", provider.config.Channel),
		"Request": params,
	}

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal tencent voice request body: %w", err)
	}

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/sms/tc3_request", date)

	headers := buildTencentHeaders(endpoint, action, voiceVersion, msg.GetExtraStringOrDefault("Region", "ap-guangzhou"), provider.config.AppSecret, bodyData, timestamp, credentialScope, provider.config.AppSecret)

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: headers,
		JSON:    requestBody,
	})
	if err != nil {
		return fmt.Errorf("tencent voice request failed: %w", err)
	}
	return provider.parseResponse(resp)
}

func (p *TencentProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，单发/群发，需模板ID，需签名，默认频率限制3000次/秒",
	)
	// 国际/港澳台短信支持单发/群发
	capabilities.SMS.International = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际/港澳台短信，单发/群发，需模板ID，签名可选，默认频率限制3000次/秒",
	)
	capabilities.SMS.Limits.MaxBatchSize = 200
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "3000次/秒"
	capabilities.SMS.Limits.DailyLimit = "无限制"
	// 彩信不支持
	capabilities.MMS.Domestic = NewRegionCapability(false, false, nil, nil, "不支持国内彩信")
	capabilities.MMS.International = NewRegionCapability(false, false, nil, nil, "不支持国际彩信")
	// 语音短信能力
	capabilities.Voice.Domestic = NewRegionCapability(
		true, false,
		[]MessageType{Voice},
		[]MessageCategory{CategoryVerification, CategoryNotification},
		"支持国内语音短信，仅支持单发，支持验证码和通知语音，需模板ID。详见：https://cloud.tencent.com/document/product/1128/51559",
	)
	capabilities.Voice.International = NewRegionCapability(
		true, false,
		[]MessageType{Voice},
		[]MessageCategory{CategoryVerification, CategoryNotification},
		"支持国际/港澳台语音短信，仅支持单发，支持验证码和通知语音，需模板ID。详见：https://cloud.tencent.com/document/product/1128/51559",
	)
	capabilities.Voice.Limits.MaxBatchSize = 1
	capabilities.Voice.Limits.RateLimit = "未知"
	capabilities.Voice.Limits.DailyLimit = "未知"
	return capabilities
}

func (p *TencentProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}

func (p *TencentProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	default:
		return Limits{}
	}
}

func (p *TencentProvider) GetName() string {
	return p.config.Name
}

func (p *TencentProvider) GetType() string {
	return string(p.config.Type)
}

func (p *TencentProvider) IsEnabled() bool {
	return !p.config.Disabled
}

func (p *TencentProvider) GetWeight() int {
	return p.config.GetWeight()
}

func (p *TencentProvider) CheckConfigured() error {
	if p.config.AppID == "" {
		return fmt.Errorf("tencent SMS provider requires AppID (SmsSdkAppId)")
	}
	if p.config.AppSecret == "" {
		return fmt.Errorf("tencent SMS provider requires AppSecret (SecretKey)")
	}
	return nil
}
