package sms

// @ProviderName: Cl253 (Chuanglan) / 创蓝253
// @Website: https://www.253.com
// @APIDoc: https://www.253.com/api
//
// 官方文档:
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
//
// CL253 支持能力:
//   - 国内短信：支持验证码、通知、营销，单发/群发，签名自动拼接，需遵守工信部规范。
//   - 国际短信：支持验证码、通知、营销，仅单发，需带国际区号，内容需以签名开头。
//   - 彩信：暂不支持。
//   - 语音短信：暂不支持。
//   - 其它能力详见官方文档。
//
// 验证码(YZM开头)、通知(N开头)、会员营销（M开头）这三类账号没有限制，但M账号22点之后成功率会受到较大影响，会出现拦截
// 对于国际短信，API请求地址不同，但仍然可以用来发送+86的国内短信，不过本实现的时候通过判断regionCode=86的时候已经移动到国内发送API了，则可能不一定符合用户预期
// 创蓝在发国内短信时，签名和营销短信的结尾是拼接在里面的
// 中文括号是代表短信签名。内容长度支持1～3500个字符（含变量）。用营销账号提交短信时最末尾需带上退订语"拒收请回复R"不支持小写r，否则营销短信将进入人工审核
// 本实现也不会关心营销短信的结尾是否需要增加[拒收请回复R], 因为这个可能会随着工信部的规定而调整，需要业务方调用的时候自己拼接。但签名本实现会自动增加。
import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

// CL253 默认 endpoint 常量，放在 init 后便于查阅
const (
	cl253DomesticEndpoint = "smssh1.253.com"
	cl253IntlEndpoint     = "intapi.253.com"
)

// Cl253Provider implements SMSProviderInterface for CL253 SMS
type Cl253Provider struct {
	config SMSProvider
}

// NewCl253Provider creates a new CL253 SMS provider
func NewCl253Provider(config SMSProvider) *Cl253Provider {
	return &Cl253Provider{
		config: config,
	}
}

// Send sends an SMS message via CL253
func (provider *Cl253Provider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	if msg.IsIntl() {
		return provider.sendIntlSMS(ctx, msg)
	}
	return provider.sendDomesticSMS(ctx, msg)
}

// sendDomesticSMS 发送国内短信（单发/群发同一API）
// https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
func (provider *Cl253Provider) sendDomesticSMS(ctx context.Context, msg *Message) error {
	endpoint := provider.config.GetEndpoint(false, cl253DomesticEndpoint)
	url := "https://" + endpoint + "/msg/v1/send/json"
	params := map[string]interface{}{
		"msg":   utils.AddSignature(msg.Content, msg.SignName),
		"phone": strings.Join(msg.Mobiles, ","),
	}
	if report := msg.GetExtraStringOrDefault("report", ""); report != "" {
		params["report"] = report
	}
	if callbackUrl := msg.GetExtraStringOrDefault("callbackUrl", provider.config.Callback); callbackUrl != "" {
		params["callbackUrl"] = callbackUrl
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if sendtime := msg.GetExtraStringOrDefault("sendtime", ""); sendtime != "" {
		params["sendtime"] = sendtime
	}
	if extend := msg.GetExtraStringOrDefault("extend", ""); extend != "" {
		params["extend"] = extend
	}
	return provider.doCl253Request(ctx, url, params)
}

// sendIntlSMS 发送国际短信（只支持单发）
// https://doc.chuanglan.com/document/O58743GF76M7754H
func (provider *Cl253Provider) sendIntlSMS(ctx context.Context, msg *Message) error {
	if msg.HasMultipleRecipients() {
		return fmt.Errorf("cl253 international SMS only supports single send")
	}
	endpoint := provider.config.GetEndpoint(true, cl253IntlEndpoint)
	url := "https://" + endpoint + "/send/sms"
	params := map[string]interface{}{
		"mobile": msg.Mobiles[0],
		"msg":    utils.AddSignature(msg.Content, msg.SignName),
	}
	if senderId := msg.GetExtraStringOrDefault("senderId", ""); senderId != "" {
		params["senderId"] = senderId
	}
	if templateId := msg.GetExtraStringOrDefault("templateId", ""); templateId != "" {
		params["templateId"] = templateId
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if tdFlag := msg.GetExtraStringOrDefault("tdFlag", ""); tdFlag != "" {
		params["tdFlag"] = tdFlag
	}
	return provider.doCl253Request(ctx, url, params)
}

// doCl253Request 统一处理 CL253 请求、签名、响应解析
func (provider *Cl253Provider) doCl253Request(ctx context.Context, url string, params map[string]interface{}) error {
	params["account"] = provider.config.AppID
	params["password"] = provider.config.AppSecret
	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method: "POST",
		JSON:   params,
	})
	if err != nil {
		return fmt.Errorf("CL253 API request failed: %w", err)
	}
	var result struct {
		Code     string `json:"code"`
		MsgId    string `json:"msgId"`
		RespTime string `json:"time"`
		ErrorMsg string `json:"errorMsg"`
		Message  string `json:"message"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse CL253 response: %w", err)
	}
	if result.Code != "0" {
		return &SMSError{
			Code:     result.Code,
			Message:  result.ErrorMsg + result.Message,
			Provider: string(ProviderTypeCl253),
		}
	}
	return nil
}

// GetCapabilities returns CL253's capabilities
func (p *Cl253Provider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// CL253 短信能力配置
	capabilities.SMS.International = NewRegionCapability(
		true, false, // 支持国际短信，但不支持群发
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际验证码、通知和营销短信，不支持群发",
	)
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内验证码、通知和营销短信",
	)
	capabilities.SMS.Limits.MaxBatchSize = 1000
	capabilities.SMS.Limits.MaxContentLen = 2000 // 国际短信限制 2000 字符
	capabilities.SMS.Limits.RateLimit = "无限制"
	capabilities.SMS.Limits.DailyLimit = ""
	// CL253 不支持彩信
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
	// CL253 不支持语音短信
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
func (p *Cl253Provider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}

// GetLimits returns CL253's limits for a specific message type
func (p *Cl253Provider) GetLimits(msgType MessageType) Limits {
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
func (p *Cl253Provider) GetName() string {
	return p.config.Name
}

// GetType returns the provider type
func (p *Cl253Provider) GetType() string {
	return string(p.config.Type)
}

// IsEnabled returns if the provider is enabled
func (p *Cl253Provider) IsEnabled() bool {
	return !p.config.Disabled
}

// GetWeight returns the provider weight
func (p *Cl253Provider) GetWeight() int {
	return p.config.GetWeight()
}
func (p *Cl253Provider) CheckConfigured() error {
	if p.config.AppID == "" || p.config.AppSecret == "" {
		return fmt.Errorf("cl253 SMS provider requires AppID and AppSecret")
	}
	return nil
}
