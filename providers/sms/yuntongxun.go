package sms

// @ProviderName: Yuntongxun / 云讯通
// @Website: https://www.yuntongxun.com
// @APIDoc: https://www.yuntongxun.com/developer-center
//
// # YuntongxunProvider implements SMSProviderInterface for Yuntongxun (容联云通讯) SMS
//
// 官方文档:
//   - 国内短信: https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
//   - 国际短信: https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d
//
// 云讯通支持国内和国际短信，使用不同的 API 版本和参数格式
// 国内短信使用模板发送，国际短信使用内容发送
// 认证方式：使用 accountSid + accountToken + datetime 生成签名
import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shellvon/go-sender/utils"
)

const (
	// 云讯通 API 配置
	yuntongxunEndpoint      = "app.cloopen.com:8883"
	yuntongxunSDKVersion    = "2013-12-26" // 国内短信 API 版本
	yuntongxunSDKVersionInt = "v2"         // 国际短信 API 版本
	yuntongxunSuccessCode   = "000000"
)

type YuntongxunProvider struct {
	config SMSProvider
}

// NewYuntongxunProvider creates a new Yuntongxun SMS provider
func NewYuntongxunProvider(config SMSProvider) *YuntongxunProvider {
	return &YuntongxunProvider{
		config: config,
	}
}

// Send sends an SMS message via Yuntongxun
func (provider *YuntongxunProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	// 根据消息类型选择不同的 API
	switch msg.Type {
	case SMSText:
		return provider.sendSMS(ctx, msg)
	case Voice:
		if msg.Category == CategoryVerification {
			return NewUnsupportedMessageTypeError(string(ProviderTypeYuntongxun), msg.Type.String(), msg.Category.String())
		}
		return provider.sendVoice(ctx, msg)
	default:
		return NewUnsupportedMessageTypeError(string(ProviderTypeYuntongxun), msg.Type.String(), "")
	}
}

// sendSMS sends SMS message via Yuntongxun API
// https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
func (provider *YuntongxunProvider) sendSMS(ctx context.Context, msg *Message) error {
	// 判断是否为国际短信
	if msg.IsIntl() {
		return provider.sendIntlSMS(ctx, msg)
	}
	return provider.sendDomesticSMS(ctx, msg)
}

// sendDomesticSMS sends domestic SMS message via Yuntongxun API
// https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
func (provider *YuntongxunProvider) sendDomesticSMS(ctx context.Context, msg *Message) error {
	datetime := time.Now().Format("20060102150405")
	accountSid := provider.config.Channel
	accountToken := provider.config.AppSecret

	// 生成签名
	sig := strings.ToUpper(utils.MD5Hex(accountSid + accountToken + datetime))

	// 构建请求头
	headers := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json;charset=utf-8",
		"Authorization": utils.Base64Encode(fmt.Sprintf("%s:%s", accountSid, datetime)),
	}

	// 构建请求体
	data := map[string]interface{}{
		"to":         strings.Join(msg.Mobiles, ","),
		"appId":      provider.config.AppID,
		"templateId": msg.TemplateID,
		"datas":      msg.ParamsOrder,
	}

	// 构建完整URL
	endpoint := provider.config.GetEndpoint(false, yuntongxunEndpoint)
	url := fmt.Sprintf("https://%s/%s/Accounts/%s/SMS/TemplateSMS?sig=%s",
		endpoint, yuntongxunSDKVersion, accountSid, sig)

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: headers,
		JSON:    data,
	})
	if err != nil {
		return fmt.Errorf("yuntongxun domestic SMS request failed: %w", err)
	}

	return provider.parseSMSResponse(resp)
}

// sendIntlSMS sends international SMS message via Yuntongxun API
// https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d
func (provider *YuntongxunProvider) sendIntlSMS(ctx context.Context, msg *Message) error {
	datetime := time.Now().Format("20060102150405")
	accountSid := provider.config.Channel
	accountToken := provider.config.AppSecret

	// 生成签名
	sig := strings.ToUpper(utils.MD5Hex(accountSid + accountToken + datetime))

	// 构建请求头
	headers := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json;charset=utf-8",
		"Authorization": utils.Base64Encode(fmt.Sprintf("%s:%s", accountSid, datetime)),
	}

	// 构建请求体
	data := map[string]interface{}{
		"mobile":  strings.Join(msg.Mobiles, ","),
		"content": utils.AddSignature(msg.Content, msg.SignName),
		"appId":   provider.config.AppID,
	}

	// 构建完整URL
	endpoint := provider.config.GetEndpoint(true, yuntongxunEndpoint)
	url := fmt.Sprintf("https://%s/%s/account/%s/international/send?sig=%s",
		endpoint, yuntongxunSDKVersionInt, accountSid, sig)

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: headers,
		JSON:    data,
	})
	if err != nil {
		return fmt.Errorf("yuntongxun international SMS request failed: %w", err)
	}

	return provider.parseSMSResponse(resp)
}

// sendVoice sends voice notification via Yuntongxun API
// https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (provider *YuntongxunProvider) sendVoice(ctx context.Context, msg *Message) error {
	// 只支持国内
	if msg.IsIntl() {
		return NewUnsupportedInternationalError(string(ProviderTypeYuntongxun), "voice call")
	}

	datetime := time.Now().Format("20060102150405")
	accountSid := provider.config.Channel
	accountToken := provider.config.AppSecret

	// 生成签名
	sig := strings.ToUpper(utils.MD5Hex(accountSid + accountToken + datetime))

	// 构建请求体
	body := map[string]interface{}{
		"to":        strings.Join(msg.Mobiles, ","),
		"appId":     provider.config.AppID,
		"mediaTxt":  msg.Content,
		"playTimes": msg.GetExtraStringOrDefault("playTimes", "3"),
	}

	// 可选参数
	if v := msg.GetExtraStringOrDefault("mediaName", ""); v != "" {
		body["mediaName"] = v
	}
	if v := msg.GetExtraStringOrDefault("displayNum", ""); v != "" {
		body["displayNum"] = v
	}
	if v := msg.GetExtraStringOrDefault("respUrl", provider.config.Callback); v != "" {
		body["respUrl"] = v
	}
	if v := msg.GetExtraStringOrDefault("userData", ""); v != "" {
		body["userData"] = v
	}
	if v := msg.GetExtraStringOrDefault("maxCallTime", ""); v != "" {
		body["maxCallTime"] = v
	}
	if v := msg.GetExtraStringOrDefault("speed", ""); v != "" {
		body["speed"] = v
	}
	if v := msg.GetExtraStringOrDefault("volume", ""); v != "" {
		body["volume"] = v
	}
	if v := msg.GetExtraStringOrDefault("pitch", ""); v != "" {
		body["pitch"] = v
	}
	if v := msg.GetExtraStringOrDefault("bgsound", ""); v != "" {
		body["bgsound"] = v
	}

	bodyData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal voice body: %w", err)
	}

	// 构建完整URL
	endpoint := provider.config.GetEndpoint(false, yuntongxunEndpoint)
	url := fmt.Sprintf("https://%s/%s/Accounts/%s/Calls/LandingCalls?sig=%s",
		endpoint, yuntongxunSDKVersion, accountSid, sig)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": utils.Base64Encode(fmt.Sprintf("%s:%s", accountSid, datetime)),
	}

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: headers,
		Raw:     bodyData,
	})
	if err != nil {
		return fmt.Errorf("yuntongxun voice call request failed: %w", err)
	}

	return provider.parseSMSResponse(resp)
}

// parseSMSResponse parses Yuntongxun SMS response
func (provider *YuntongxunProvider) parseSMSResponse(resp []byte) error {
	var result struct {
		StatusCode string `json:"statusCode"`
		StatusMsg  string `json:"statusMsg"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse yuntongxun response: %w", err)
	}
	if result.StatusCode != yuntongxunSuccessCode {
		return &SMSError{
			Code:     result.StatusCode,
			Message:  result.StatusMsg,
			Provider: string(ProviderTypeYuntongxun),
		}
	}
	return nil
}

// GetCapabilities returns Yuntongxun's capabilities
func (p *YuntongxunProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 云讯通 SMS 能力配置
	capabilities.SMS.International = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际验证码、通知和营销短信",
	)
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内验证码、通知和营销短信",
	)
	capabilities.SMS.Limits.MaxBatchSize = 1000
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "无限制"
	capabilities.SMS.Limits.DailyLimit = ""
	// 云讯通不支持彩信
	capabilities.MMS.International = NewRegionCapability(
		false, false,
		[]MessageType{},
		[]MessageCategory{},
		"不支持国际彩信",
	)
	capabilities.MMS.Domestic = NewRegionCapability(
		false, false,
		[]MessageType{},
		[]MessageCategory{},
		"不支持国内彩信",
	)
	// 云讯通不支持语音短信
	capabilities.Voice.International = NewRegionCapability(
		false, false,
		[]MessageType{},
		[]MessageCategory{},
		"不支持国际语音短信",
	)
	capabilities.Voice.Domestic = NewRegionCapability(
		false, false,
		[]MessageType{},
		[]MessageCategory{},
		"不支持国内语音短信",
	)
	return capabilities
}

// CheckCapability checks if a specific capability is supported
func (p *YuntongxunProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}

// GetLimits returns Yuntongxun's limits for a specific message type
func (p *YuntongxunProvider) GetLimits(msgType MessageType) Limits {
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
func (p *YuntongxunProvider) GetName() string {
	return p.config.Name
}

// GetType returns the provider type
func (p *YuntongxunProvider) GetType() string {
	return string(p.config.Type)
}

// IsEnabled returns if the provider is enabled
func (p *YuntongxunProvider) IsEnabled() bool {
	return !p.config.Disabled
}

// GetWeight returns the provider weight
func (p *YuntongxunProvider) GetWeight() int {
	return p.config.GetWeight()
}

func (p *YuntongxunProvider) CheckConfigured() error {
	if p.config.Channel == "" || p.config.AppSecret == "" || p.config.AppID == "" {
		return fmt.Errorf("yuntongxun SMS provider requires Channel (accountSid), AppSecret (accountToken), and AppID")
	}
	return nil
}
