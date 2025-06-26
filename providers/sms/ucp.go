package sms

// @ProviderName: UCP / 云之讯
// @Website: https://www.ucpaas.com
// @APIDoc: http://docs.ucpaas.com
//
// # UcpProvider implements SMSProviderInterface for UCP SMS
//
// 官方文档:
//   - 短信API文档: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:about_sms
//   - 国内外指定模版单发: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms
//   - 国内外指定模版群发: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms_batch
//
// 能力说明:
//   - 国内短信：支持单发和群发，需模板ID。
//   - 国际短信：支持单发和群发，需模板ID。
//   - 彩信/语音：暂不支持。
//
// 注意：支持国内外手机号码，需模板ID。
import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

// http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms
const (
	ucpEndpoint = "open.ucpaas.com"
)

type UcpProvider struct {
	config SMSProvider
}

// NewUcpProvider creates a new UCP SMS provider
func NewUcpProvider(config SMSProvider) *UcpProvider {
	return &UcpProvider{config: config}
}

// Send sends an SMS message via UCP
func (provider *UcpProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	if msg.Type != SMSText {
		return NewUnsupportedMessageTypeError(string(ProviderTypeUcp), msg.Type.String(), msg.Category.String())
	}
	return provider.sendSMS(ctx, msg)
}

// sendSMS sends SMS message via UCP API
func (provider *UcpProvider) sendSMS(ctx context.Context, msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return fmt.Errorf("ucp: mobiles cannot be empty")
	}

	// 根据手机号数量选择API
	var apiPath string
	if msg.HasMultipleRecipients() {
		apiPath = "templatesms"
	} else {
		apiPath = "variablesms"
	}

	params := map[string]interface{}{
		"clientid":   provider.config.AppID,
		"password":   provider.config.AppSecret,
		"templateid": msg.TemplateID,
		"mobile":     strings.Join(msg.Mobiles, ","),
	}

	// 模板参数处理
	if len(msg.ParamsOrder) > 0 {
		// 模板中的替换参数，如该模板不存在参数则无需传该参数或者参数为空，如果有多个参数则需要写在同一个字符串中，以分号分隔 （如："a;b;c"），参数中不能含有特殊符号"【】"和","
		params["param"] = strings.Join(msg.ParamsOrder, ";")
	}

	// 可选参数：uid（用户自定义ID，用于回调时识别）
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}

	endpoint := provider.config.GetEndpoint(msg.IsIntl(), ucpEndpoint)
	url := "https://" + endpoint + "/sms-server/" + apiPath

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/json"},
		JSON:    params,
	})
	if err != nil {
		return fmt.Errorf("ucp SMS request failed: %w", err)
	}
	return provider.parseUcpResponse(resp)
}

// parseUcpResponse parses UCP SMS response
func (provider *UcpProvider) parseUcpResponse(resp []byte) error {
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse ucp response: %w", err)
	}
	if result.Code != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.Code),
			Message:  result.Msg,
			Provider: string(ProviderTypeUcp),
		}
	}
	return nil
}

func (p *UcpProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，单发/群发，需模板ID",
	)
	// 国际短信支持单发/群发
	capabilities.SMS.International = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际短信，单发/群发，需模板ID",
	)
	capabilities.SMS.Limits.MaxBatchSize = 1000
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "未知"
	capabilities.SMS.Limits.DailyLimit = "未知"
	// 彩信、语音均不支持
	capabilities.MMS.Domestic = NewRegionCapability(false, false, nil, nil, "不支持国内彩信")
	capabilities.MMS.International = NewRegionCapability(false, false, nil, nil, "不支持国际彩信")
	capabilities.Voice.Domestic = NewRegionCapability(false, false, nil, nil, "不支持国内语音")
	capabilities.Voice.International = NewRegionCapability(false, false, nil, nil, "不支持国际语音")
	return capabilities
}

func (p *UcpProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}

func (p *UcpProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	default:
		return Limits{}
	}
}

func (p *UcpProvider) GetName() string {
	return p.config.Name
}

func (p *UcpProvider) GetType() string {
	return string(p.config.Type)
}

func (p *UcpProvider) IsEnabled() bool {
	return !p.config.Disabled
}

func (p *UcpProvider) GetWeight() int {
	return p.config.GetWeight()
}

func (p *UcpProvider) CheckConfigured() error {
	if p.config.AppID == "" || p.config.AppSecret == "" {
		return fmt.Errorf("ucp provider requires AppID and AppSecret")
	}
	return nil
}
