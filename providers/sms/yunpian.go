package sms

// @ProviderName: Yunpian / 云片
// @Website: https://www.yunpian.com
// @APIDoc: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
//
// # YunpianProvider implements SMSProviderInterface for Yunpian SMS
//
// 官方文档:
// - 短信API文档: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
// - 国内短信API: https://www.yunpian.com/official/document/sms/zh_CN/domestic_single_send
// - 国际短信API: https://www.yunpian.com/official/document/sms/zh_CN/international_single_send
//
// 云片支持超级短信（彩信）:
//   - 超级短信API: https://www.yunpian.com/official/document/sms/zh_CN/super_sms_send
//   - 超级短信功能需要联系客服开通后使用
//   - 因运营商夜间防骚扰限制，建议在早9:00至晚18:00间发送
//
// 云片支持语音短信:
//   - 语音API文档: https://www.yunpian.com/official/document/sms/zh_CN/voice_send
//   - 语音验证码支持国内和国际，仅支持单发验证码
//   - 频率限制：30秒1次，1小时3次，24小时5次
import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

type YunpianProvider struct {
	config SMSProvider
}

// NewYunpianProvider creates a new Yunpian SMS provider
func NewYunpianProvider(config SMSProvider) *YunpianProvider {
	return &YunpianProvider{config: config}
}

const (
	yunpianDomain                = "yunpian.com"
	yunpianSMSService            = "sms"
	yunpianVMSService            = "vsms"
	yunpianVoiceService          = "voice"
	yunpianSingleSendPath        = "/v2/sms/single_send.json"
	yunpianBatchSendPath         = "/v2/sms/batch_send.json"
	yunpianTplSingleSendPath     = "/v2/sms/tpl_single_send.json"
	yunpianTplBatchSendPath      = "/v2/sms/tpl_batch_send.json"
	yunpianSuperMMSBatchSendPath = "/v2/vsms/tpl_batch_send.json"
	yunpianVoiceSendPath         = "/v2/voice/send.json"
)

// yunpianEndpoint 统一生成云片 API endpoint，支持可选自定义域名（如 us.yunpian.com）
func yunpianEndpoint(service, path string, domainOverride ...string) string {
	domain := fmt.Sprintf("%s.%s", service, yunpianDomain)
	if len(domainOverride) > 0 && domainOverride[0] != "" {
		domain = domainOverride[0]
	}
	return fmt.Sprintf("https://%s%s", domain, path)
}

// Send sends an SMS/Voice/MMS message via Yunpian
func (provider *YunpianProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	if msg.Type == MMS {
		return provider.sendMMS(ctx, msg)
	}
	if msg.Type == Voice && msg.Category == CategoryVerification {
		return provider.sendVoice(ctx, msg)
	}
	if msg.IsIntl() {
		if msg.HasMultipleRecipients() {
			return fmt.Errorf("yunpian international SMS only supports single send")
		}
		if msg.TemplateID != "" {
			return fmt.Errorf("yunpian international SMS does not support template")
		}
		return provider.sendIntlSMS(ctx, msg)
	}
	// 国内短信
	if msg.TemplateID != "" {
		if msg.HasMultipleRecipients() {
			return provider.sendTplBatchSMS(ctx, msg)
		}
		return provider.sendTplSMS(ctx, msg)
	} else {
		if msg.HasMultipleRecipients() {
			return provider.sendBatchSMS(ctx, msg)
		}
		return provider.sendSMS(ctx, msg)
	}
}

// 国内短信-单条发送接口
// API文档: https://www.yunpian.com/official/document/sms/zh_cn/domestic_single_send
//   - URL：https://sms.yunpian.com/v2/sms/single_send.json
//   - 注意：海外服务器地址 us.yunpian.com
//   - 访问方式：POST
func (provider *YunpianProvider) sendSMS(ctx context.Context, msg *Message) error {
	// 支持海外服务器地址，用户可在 config 配置覆盖
	defaultEndpoint := yunpianEndpoint(yunpianSMSService, yunpianSingleSendPath)
	endpoint := provider.config.GetEndpoint(false, defaultEndpoint)
	params := map[string]string{
		"apikey": provider.config.AppSecret,
		"mobile": msg.Mobiles[0],
		"text":   utils.AddSignature(msg.Content, msg.SignName),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", provider.config.Callback); cb != "" {
		params["callback_url"] = cb
	}
	if reg, ok := msg.GetExtraBool("register"); ok {
		params["register"] = fmt.Sprintf("%v", reg)
	}
	if stat, ok := msg.GetExtraBool("mobile_stat"); ok {
		params["mobile_stat"] = fmt.Sprintf("%v", stat)
	}
	return sendYunpianRequest(ctx, endpoint, params, provider.config.Type)
}

// 国内短信-批量发送接口
// API文档: https://www.yunpian.com/official/document/sms/zh_cn/domestic_batch_send
//   - URL：https://sms.yunpian.com/v2/sms/batch_send.json
//   - 注意：海外服务器地址 us.yunpian.com
//   - 访问方式：POST
func (provider *YunpianProvider) sendBatchSMS(ctx context.Context, msg *Message) error {
	defaultEndpoint := yunpianEndpoint(yunpianSMSService, yunpianBatchSendPath)
	endpoint := provider.config.GetEndpoint(false, defaultEndpoint)
	params := map[string]string{
		"apikey": provider.config.AppSecret,
		"mobile": strings.Join(msg.Mobiles, ","),
		"text":   utils.AddSignature(msg.Content, msg.SignName),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", provider.config.Callback); cb != "" {
		params["callback_url"] = cb
	}
	if stat, ok := msg.GetExtraBool("mobile_stat"); ok {
		params["mobile_stat"] = fmt.Sprintf("%v", stat)
	}
	return sendYunpianRequest(ctx, endpoint, params, provider.config.Type)
}

// 国内短信-指定模板单发
// API文档: https://www.yunpian.com/official/document/sms/zh_CN/domestic_tpl_single_send
//   - URL：https://sms.yunpian.com/v2/sms/tpl_single_send.json
//   - 注意：海外服务器地址 us.yunpian.com
//   - 访问方式：POST
func (provider *YunpianProvider) sendTplSMS(ctx context.Context, msg *Message) error {
	defaultEndpoint := yunpianEndpoint(yunpianSMSService, yunpianTplSingleSendPath)
	endpoint := provider.config.GetEndpoint(false, defaultEndpoint)
	params := map[string]string{
		"apikey":    provider.config.AppSecret,
		"mobile":    msg.Mobiles[0],
		"tpl_id":    msg.TemplateID,
		"tpl_value": buildTemplateValue(msg.TemplateParams),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", provider.config.Callback); cb != "" {
		params["callback_url"] = cb
	}
	return sendYunpianRequest(ctx, endpoint, params, provider.config.Type)
}

// 国内短信-指定模板群发
// API文档: https://www.yunpian.com/official/document/sms/zh_CN/domestic_tpl_batch_send
//   - URL：https://sms.yunpian.com/v2/sms/tpl_batch_send.json
//   - 注意：海外服务器地址 us.yunpian.com
//   - 访问方式：POST
func (provider *YunpianProvider) sendTplBatchSMS(ctx context.Context, msg *Message) error {
	defaultEndpoint := yunpianEndpoint(yunpianSMSService, yunpianTplBatchSendPath)
	endpoint := provider.config.GetEndpoint(false, defaultEndpoint)
	params := map[string]string{
		"apikey":    provider.config.AppSecret,
		"mobile":    strings.Join(msg.Mobiles, ","),
		"tpl_id":    msg.TemplateID,
		"tpl_value": buildTemplateValue(msg.TemplateParams),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", provider.config.Callback); cb != "" {
		params["callback_url"] = cb
	}
	return sendYunpianRequest(ctx, endpoint, params, provider.config.Type)
}

// 超级短信（彩信）发送接口
// API文档: https://www.yunpian.com/official/document/sms/zh_CN/super_sms_send
//   - URL：https://vsms.yunpian.com/v2/vsms/tpl_batch_send.json
//   - 访问方式：POST
//   - 超级短信功能需要联系客服开通后使用。
//   - 因运营商夜间防骚扰限制，建议您在早9:00至晚18:00间发送，超出此时间的发送可能会失败。
//   - 仅支持国内手机号，tpl_id 必填，支持单发/群发。
func (provider *YunpianProvider) sendMMS(ctx context.Context, msg *Message) error {
	defaultEndpoint := yunpianEndpoint(yunpianVMSService, yunpianSuperMMSBatchSendPath)
	endpoint := provider.config.GetEndpoint(false, defaultEndpoint)
	if msg.IsIntl() {
		return fmt.Errorf("yunpian super sms (MMS) only supports domestic numbers")
	}
	if msg.TemplateID == "" {
		return fmt.Errorf("yunpian super sms (MMS) requires tpl_id (template ID)")
	}
	if len(msg.Mobiles) == 0 {
		return fmt.Errorf("yunpian super sms (MMS) requires at least one mobile number")
	}
	params := map[string]string{
		"apikey": provider.config.AppSecret,
		"tpl_id": msg.TemplateID,
		"mobile": strings.Join(msg.Mobiles, ","),
	}
	if len(msg.TemplateParams) > 0 {
		// 超级短信要求 tpl_param_json，内容为 JSON 字符串
		jsonBytes, err := json.Marshal(msg.TemplateParams)
		if err != nil {
			return fmt.Errorf("failed to marshal tpl_param_json: %w", err)
		}
		params["tpl_param_json"] = string(jsonBytes)
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", provider.config.Callback); cb != "" {
		params["callback_url"] = cb
	}
	return sendYunpianRequest(ctx, endpoint, params, provider.config.Type)
}

// 国际短信-单条发送接口
// API文档: https://www.yunpian.com/official/document/sms/zh_CN/intl_single_send
//   - URL：https://sms.yunpian.com/v2/sms/single_send.json
//   - 注意：海外服务器地址 us.yunpian.com
//   - 访问方式：POST
//   - 国际短信不支持批量发送，每次只能发送一条
//   - 国际号码需包含国际地区前缀号码（如 +93701234567）
//   - 发送内容需与已审核的短信模板相匹配
//   - 支持 https 与 http，两者均可，建议使用 https
func (provider *YunpianProvider) sendIntlSMS(ctx context.Context, msg *Message) error {
	// 默认用海外域名 us.yunpian.com
	defaultEndpoint := yunpianEndpoint(yunpianSMSService, yunpianSingleSendPath)
	endpoint := provider.config.GetEndpoint(true, defaultEndpoint)
	params := map[string]string{
		"apikey": provider.config.AppSecret,
		"mobile": msg.Mobiles[0],
		"text":   msg.Content,
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", provider.config.Callback); cb != "" {
		params["callback_url"] = cb
	}
	if reg, ok := msg.GetExtraBool("register"); ok {
		params["register"] = fmt.Sprintf("%v", reg)
	}
	return sendYunpianRequest(ctx, endpoint, params, provider.config.Type)
}

// 语音验证码发送接口
// API文档: https://www.yunpian.com/official/document/sms/zh_CN/voicecode_list
//   - URL：https://voice.yunpian.com/v2/voice/send.json
//   - 访问方式：POST
//   - endpoint 固定，不支持自定义
//   - 仅支持单个手机号，国际号码需包含国际区号（如+93701234567），国内号码直接传号码
//   - 参数：apikey, mobile, code, language(国际), uid, callback_url
//   - 国内(+86)号码每天22:00至次日7:30禁止发送
//   - 频率限制：同一手机号30秒1次，1小时3次，24小时5次
func (provider *YunpianProvider) sendVoice(ctx context.Context, msg *Message) error {
	defaultEndpoint := yunpianEndpoint(yunpianVoiceService, yunpianVoiceSendPath)
	endpoint := provider.config.GetEndpoint(msg.IsIntl(), defaultEndpoint)
	var mobile string
	if msg.IsIntl() {
		mobile = fmt.Sprintf("+%d%s", msg.RegionCode, msg.Mobiles[0])
	} else {
		mobile = msg.Mobiles[0]
	}
	params := map[string]string{
		"apikey": provider.config.AppSecret,
		"mobile": mobile,
		"code":   msg.Content, // 4~6位数字验证码
	}
	if msg.IsIntl() {
		// 国际号码 language 可选
		if lang := msg.GetExtraStringOrDefault("language", ""); lang != "" {
			params["language"] = lang
		}
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", provider.config.Callback); cb != "" {
		params["callback_url"] = cb
	}
	return sendYunpianRequest(ctx, endpoint, params, provider.config.Type)
}

// YunpianResponse represents the common response structure from Yunpian API
// API Documentation: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
type YunpianResponse struct {
	Code int             `json:"code"`           // Response code, 0 means success
	Msg  string          `json:"msg"`            // Response message
	Data json.RawMessage `json:"data,omitempty"` // Response data
}

// sendYunpianRequest sends the actual HTTP request to Yunpian API
func sendYunpianRequest(ctx context.Context, endpoint string, params map[string]string, providerType ProviderType) error {
	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Data: params,
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

func (p *YunpianProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，单发/群发，支持模板和非模板",
	)
	// 国际短信仅支持单发
	capabilities.SMS.International = NewRegionCapability(
		true, false,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际短信，仅单发，支持模板和非模板",
	)
	capabilities.SMS.Limits.MaxBatchSize = 1000
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "未知"
	capabilities.SMS.Limits.DailyLimit = "未知"
	// 彩信（超级短信）
	capabilities.MMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{MMS},
		[]MessageCategory{CategoryNotification, CategoryPromotion},
		"支持超级短信（彩信），单发/群发，需模板，支持图片/视频等多媒体内容",
	)
	capabilities.MMS.International = NewRegionCapability(
		false, false, nil, nil, "不支持国际彩信",
	)
	// 语音短信
	capabilities.Voice.Domestic = NewRegionCapability(
		true, false,
		[]MessageType{Voice},
		[]MessageCategory{CategoryVerification},
		"仅支持单发验证码，30秒1次，1小时3次，24小时5次",
	)
	capabilities.Voice.International = NewRegionCapability(
		true, false,
		[]MessageType{Voice},
		[]MessageCategory{CategoryVerification},
		"仅支持国际单发验证码，30秒1次，1小时3次，24小时5次",
	)
	capabilities.Voice.Limits.RateLimit = "30秒1次，1小时3次，24小时5次"
	return capabilities
}
func (p *YunpianProvider) CheckCapability(msg *Message) error {
	// 国际短信只允许单发
	if msg.IsIntl() && msg.HasMultipleRecipients() {
		return fmt.Errorf("yunpian international SMS only supports single send")
	}
	return DefaultCheckCapability(p, msg)
}
func (p *YunpianProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	default:
		return Limits{}
	}
}
func (p *YunpianProvider) GetName() string {
	return p.config.Name
}
func (p *YunpianProvider) GetType() string {
	return string(p.config.Type)
}
func (p *YunpianProvider) IsEnabled() bool {
	return !p.config.Disabled
}
func (p *YunpianProvider) GetWeight() int {
	return p.config.GetWeight()
}
func (p *YunpianProvider) CheckConfigured() error {
	if p.config.AppSecret == "" {
		return fmt.Errorf("yunpian provider requires AppSecret (apikey)")
	}
	return nil
}
