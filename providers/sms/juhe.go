package sms

// @ProviderName: Juhe / 聚合数据
// @Website: https://www.juhe.cn
// @APIDoc: https://www.juhe.cn/docs
//
// # JuheProvider implements SMSProviderInterface for Juhe SMS
//
// 官方文档:
//   - 短信API文档: https://www.juhe.cn/docs/api/id/54
//   - 国内短信API: https://www.juhe.cn/docs/api/id/54
//   - 国际短信API: https://www.juhe.cn/docs/api/id/357
//   - 视频短信API: https://www.juhe.cn/docs/api/id/363
//
// 能力说明:
//   - 国内短信：仅支持单发，需模板ID。
//   - 国际短信：仅支持单发，需模板ID，需带区号。
//   - 视频短信（彩信）：仅支持单发，仅移动号码。
//   - 语音：暂不支持。
//
// 注意：短信仅支持单发，不支持群发。
import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shellvon/go-sender/utils"
)

const (
	juheEndpoint = "v.juhe.cn"
)

type JuheProvider struct {
	config SMSProvider
}

// NewJuheProvider creates a new Juhe SMS provider
func NewJuheProvider(config SMSProvider) *JuheProvider {
	return &JuheProvider{config: config}
}

// Send sends an SMS/MMS message via Juhe
func (provider *JuheProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	switch msg.Type {
	case MMS:
		return provider.sendMMS(ctx, msg)
	default:
		if msg.IsIntl() {
			return provider.sendIntlSMS(ctx, msg)
		}
		return provider.sendSMS(ctx, msg)
	}
}

// sendSMS sends domestic SMS message via Juhe API
// 国内短信API: https://www.juhe.cn/docs/api/id/54
func (provider *JuheProvider) sendSMS(ctx context.Context, msg *Message) error {
	endpoint := provider.config.GetEndpoint(false, juheEndpoint)
	url := "http://" + endpoint + "/sms/send"
	params := map[string]string{
		"mobile": msg.Mobiles[0],
		"tpl_id": msg.TemplateID,
		"key":    provider.config.AppID,
		"vars":   utils.ToJSONString(msg.TemplateParams),
	}
	// 扩展参数
	if ext := msg.GetExtraStringOrDefault("ext", ""); ext != "" {
		params["ext"] = ext
	}

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Data:    params,
	})
	if err != nil {
		return fmt.Errorf("juhe SMS request failed: %w", err)
	}
	return parseJuheResponse(resp)
}

// sendIntlSMS sends international SMS message via Juhe API
// 国际短信API: https://www.juhe.cn/docs/api/id/357
func (provider *JuheProvider) sendIntlSMS(ctx context.Context, msg *Message) error {
	endpoint := provider.config.GetEndpoint(true, juheEndpoint)
	url := "http://" + endpoint + "/smsInternational/send"
	params := map[string]string{
		"key":      provider.config.AppID,
		"mobile":   msg.Mobiles[0],
		"areaNum":  fmt.Sprintf("%d", msg.RegionCode),
		"tplId":    msg.TemplateID,
		"tplValue": buildTemplateValue(msg.TemplateParams),
	}
	// tplValue	否	string	如果您的模板里面有变量则需要提交此参数,如:#code#=123456,参数需要urlencode

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Data:    params,
	})
	if err != nil {
		return fmt.Errorf("juhe intl SMS request failed: %w", err)
	}
	return parseJuheResponse(resp)
}

// sendMMS sends MMS message via Juhe API
// 聚合数据视频短信API文档: https://www.juhe.cn/docs/api/id/363
func (provider *JuheProvider) sendMMS(ctx context.Context, msg *Message) error {
	endpoint := provider.config.GetEndpoint(false, juheEndpoint)
	url := "http://" + endpoint + "/caixinv2/send"

	if len(msg.Mobiles) != 1 {
		return fmt.Errorf("juhe MMS only supports single mobile per request")
	}
	// 模板ID是必需的
	if msg.TemplateID == "" {
		return fmt.Errorf("juhe MMS requires template ID")
	}

	// 聚合数据视频短信API参数
	params := map[string]string{
		"mobile": msg.Mobiles[0],
		"key":    provider.config.AppID,
		"tplId":  msg.TemplateID,
	}

	// 模板参数处理 - 使用vars参数
	if len(msg.TemplateParams) > 0 {
		// 聚合数据要求vars参数为JSON格式
		varsJSON, err := json.Marshal(msg.TemplateParams)
		if err != nil {
			return fmt.Errorf("failed to marshal template params: %w", err)
		}
		params["vars"] = string(varsJSON)
	}

	// 扩展参数
	if ext := msg.GetExtraStringOrDefault("ext", ""); ext != "" {
		params["ext"] = ext
	}

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Data:    params,
	})
	if err != nil {
		return fmt.Errorf("juhe MMS request failed: %w", err)
	}
	return parseJuheResponse(resp)
}

func parseJuheResponse(resp []byte) error {
	var result struct {
		ErrorCode int    `json:"error_code"`
		Reason    string `json:"reason"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse juhe response: %w", err)
	}
	if result.ErrorCode != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.ErrorCode),
			Message:  result.Reason,
			Provider: string(ProviderTypeJuhe),
		}
	}
	return nil
}

func (p *JuheProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信
	capabilities.SMS.Domestic = NewRegionCapability(
		true, false,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，仅支持单发，需模板ID",
	)
	// 国际短信
	capabilities.SMS.International = NewRegionCapability(
		true, false,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际短信，仅支持单发，需模板ID，需带区号",
	)
	capabilities.SMS.Limits.MaxBatchSize = 1
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "未知"
	capabilities.SMS.Limits.DailyLimit = "未知"
	// 彩信能力
	capabilities.MMS.Domestic = NewRegionCapability(
		true, false,
		[]MessageType{MMS},
		[]MessageCategory{CategoryNotification, CategoryPromotion},
		`支持国内视频短信（彩信），仅支持单发，仅移动号码。状态码为0计费，其余不计费，目前只能发送移动号码。视频大小限制：2048K，展现形式：平铺，发送速度：300条/s。（游戏、保健品等）频次：1条/日,3条/周；（会员营销）频次：1条/日，量级动态调整。屏蔽：北京、沈阳，发送时间：早九晚六。详见: https://www.juhe.cn/docs/api/id/363`,
	)
	capabilities.MMS.Limits.MaxBatchSize = 1
	capabilities.MMS.Limits.MaxContentLen = 2048 * 1024 // 2MB
	capabilities.MMS.Limits.RateLimit = "300条/s"
	capabilities.MMS.Limits.DailyLimit = "1条/日"
	// 国际彩信不支持
	capabilities.MMS.International = NewRegionCapability(false, false, nil, nil, "不支持国际彩信")
	// 语音能力
	capabilities.Voice.Domestic = NewRegionCapability(false, false, nil, nil, "不支持国内语音")
	capabilities.Voice.International = NewRegionCapability(false, false, nil, nil, "不支持国际语音")
	return capabilities
}

func (p *JuheProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}

func (p *JuheProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	default:
		return Limits{}
	}
}

func (p *JuheProvider) GetName() string {
	return p.config.Name
}

func (p *JuheProvider) GetType() string {
	return string(p.config.Type)
}

func (p *JuheProvider) IsEnabled() bool {
	return !p.config.Disabled
}

func (p *JuheProvider) GetWeight() int {
	return p.config.GetWeight()
}

func (p *JuheProvider) CheckConfigured() error {
	if p.config.AppID == "" {
		return fmt.Errorf("juhe SMS provider requires AppID (key)")
	}
	return nil
}
