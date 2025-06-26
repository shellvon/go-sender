package sms

// @ProviderName: Smsbao / 短信宝
// @Website: https://www.smsbao.com
// @APIDoc: https://www.smsbao.com/openapi
//
// # SmsbaoProvider implements SMSProviderInterface for Smsbao SMS
//
// 官方文档:
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//   - 国际短信: https://www.smsbao.com/openapi/299.html
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
//
// 能力说明:
//   - 国内短信：支持单发和群发，最多99个号码/次。
//   - 国际短信：支持单发和群发，最多99个号码/次。
//   - 语音验证码：仅支持单个号码，仅支持验证码类型。
//
// 注意：语音验证码仅支持国内手机号码，不支持国际号码。
import (
	"context"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

const smsbaoEndpoint = "api.smsbao.com"

type SmsbaoProvider struct {
	config SMSProvider
}

// NewSmsbaoProvider creates a new Smsbao SMS provider
func NewSmsbaoProvider(config SMSProvider) *SmsbaoProvider {
	return &SmsbaoProvider{config: config}
}

// Send sends an SMS message via Smsbao
func (provider *SmsbaoProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	switch msg.Type {
	case SMSText:
		return provider.sendSMS(ctx, msg)
	case Voice:
		return provider.sendVoice(ctx, msg)
	default:
		return NewUnsupportedMessageTypeError(string(ProviderTypeSmsbao), msg.Type.String(), msg.Category.String())
	}
}

// sendSMS sends SMS message via Smsbao API
func (provider *SmsbaoProvider) sendSMS(ctx context.Context, msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return fmt.Errorf("smsbao: mobiles cannot be empty")
	}
	if len(msg.Mobiles) > 99 {
		return fmt.Errorf("smsbao: at most 99 mobiles per request")
	}

	params := map[string]interface{}{
		"u": provider.config.AppID,
		"p": provider.config.AppSecret, // 需为MD5后的密码或ApiKey
		"m": strings.Join(msg.Mobiles, ","),
		"c": msg.Content,
	}
	if provider.config.Channel != "" {
		params["g"] = provider.config.Channel
	}

	isIntl := msg.IsIntl()

	// 根据是否为国际短信选择不同的API
	var apiPath string
	if isIntl {
		apiPath = "/wsms"
	} else {
		apiPath = "/sms"
	}

	endpoint := provider.config.GetEndpoint(isIntl, smsbaoEndpoint)
	url := "https://" + endpoint + apiPath

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method: "GET",
		Query:  params,
	})
	if err != nil {
		return fmt.Errorf("smsbao SMS request failed: %w", err)
	}
	return parseSmsbaoResponse(resp)
}

// sendVoice sends voice verification code via Smsbao API
// https://www.smsbao.com/openapi/214.html
// 同一被叫号码30秒内最多发送 1 条。
func (provider *SmsbaoProvider) sendVoice(ctx context.Context, msg *Message) error {
	// 只支持国内
	if msg.IsIntl() {
		return NewUnsupportedInternationalError(string(ProviderTypeSmsbao), "voice call")
	}
	// 只支持验证码
	if msg.Category != CategoryVerification {
		return NewUnsupportedMessageTypeError(string(ProviderTypeSmsbao), msg.Type.String(), msg.Category.String())
	}
	// 只支持单号码
	if msg.HasMultipleRecipients() {
		return fmt.Errorf("smsbao voice only supports single number, got %d", len(msg.Mobiles))
	}

	params := map[string]interface{}{
		"u": provider.config.AppID,
		"p": provider.config.AppSecret, // 需为MD5后的密码或ApiKey
		"m": msg.Mobiles[0],
		"c": msg.Content,
	}
	if provider.config.Channel != "" {
		params["g"] = provider.config.Channel
	}

	endpoint := provider.config.GetEndpoint(false, smsbaoEndpoint)
	url := "http://" + endpoint + "/voice"

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method: "GET",
		Query:  params,
	})
	if err != nil {
		return fmt.Errorf("smsbao voice request failed: %w", err)
	}
	return parseSmsbaoResponse(resp)
}

// parseSmsbaoResponse 解析短信宝响应
func parseSmsbaoResponse(resp []byte) error {
	result := strings.TrimSpace(string(resp))
	if result == "0" {
		return nil
	}
	return &SMSError{
		Code:     result,
		Message:  smsbaoErrorMsg(result),
		Provider: string(ProviderTypeSmsbao),
	}
}

// smsbaoErrorMsg 错误码转中文
func smsbaoErrorMsg(code string) string {
	switch code {
	case "0":
		return "短信发送成功"
	case "-1":
		return "参数不全"
	case "-2":
		return "服务器空间不支持,请确认支持curl或者fsocket，联系您的空间商解决或者更换空间！"
	case "30":
		return "密码错误"
	case "40":
		return "账号不存在"
	case "41":
		return "余额不足"
	case "42":
		return "帐户已过期"
	case "43":
		return "IP地址限制"
	case "50":
		return "内容含有敏感词"
	case "51":
		return "手机号码不正确"
	default:
		return "未知错误"
	}
}
func (p *SmsbaoProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，单发/群发，最多99个号码/次",
	)
	// 国际短信支持单发/群发
	capabilities.SMS.International = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际短信，单发/群发，最多99个号码/次",
	)
	capabilities.SMS.Limits.MaxBatchSize = 99
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "未知"
	capabilities.SMS.Limits.DailyLimit = "未知"
	// 彩信不支持
	capabilities.MMS.Domestic = NewRegionCapability(false, false, nil, nil, "不支持国内彩信")
	capabilities.MMS.International = NewRegionCapability(false, false, nil, nil, "不支持国际彩信")
	// 语音验证码支持（仅国内，仅验证码，仅单号码）
	capabilities.Voice.Domestic = NewRegionCapability(
		true, false, // 支持单发，不支持群发
		[]MessageType{Voice},
		[]MessageCategory{CategoryVerification}, // 仅支持验证码
		"支持国内语音验证码，仅支持单号码",
	)
	capabilities.Voice.International = NewRegionCapability(false, false, nil, nil, "不支持国际语音")
	return capabilities
}
func (p *SmsbaoProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}
func (p *SmsbaoProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	case Voice:
		return capabilities.Voice.GetLimits()
	default:
		return Limits{}
	}
}
func (p *SmsbaoProvider) GetName() string {
	return p.config.Name
}
func (p *SmsbaoProvider) GetType() string {
	return string(p.config.Type)
}
func (p *SmsbaoProvider) IsEnabled() bool {
	return !p.config.Disabled
}
func (p *SmsbaoProvider) GetWeight() int {
	return p.config.GetWeight()
}
func (p *SmsbaoProvider) CheckConfigured() error {
	if p.config.AppID == "" || p.config.AppSecret == "" {
		return fmt.Errorf("smsbao provider requires AppID and AppSecret")
	}
	return nil
}
