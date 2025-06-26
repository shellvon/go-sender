package sms

// @ProviderName: Luosimao / 螺丝帽
// @Website: https://luosimao.com
// @APIDoc: https://luosimao.com/docs
//
// # LuosimaoProvider implements SMSProviderInterface for Luosimao SMS
//
// 官方文档:
//   - 短信API文档: https://luosimao.com/docs
//   - 短信发送API: https://luosimao.com/docs
//
// 能力说明:
//   - 国内短信：支持单发和群发，内容需包含签名，批量最多10万条/次。
//   - 国际短信：暂不支持。
//   - 彩信/语音：暂不支持。
//
// 注意：仅支持国内手机号码，不支持国际号码。
import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

const (
	luosimaoDomain        = "luosimao.com"
	luosimaoSMSService    = "sms-api"
	luosimaoVoiceService  = "voice-api"
	luosimaoSMSSendPath   = "/v1/send.json"
	luosimaoSMSBatchPath  = "/v1/send_batch.json"
	luosimaoVoiceSendPath = "/v1/verify.json"
	luosimaoSmsURL        = "https://sms-api.luosimao.com/v1/send.json"
	luosimaoVoiceURL      = "https://voice-api.luosimao.com/v1/verify.json"
)

type LuosimaoProvider struct {
	config SMSProvider
}

// NewLuosimaoProvider creates a new Luosimao SMS provider
func NewLuosimaoProvider(config SMSProvider) *LuosimaoProvider {
	return &LuosimaoProvider{config: config}
}

// Send sends an SMS message via Luosimao
func (provider *LuosimaoProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	switch msg.Type {
	case Voice:
		if msg.Category == CategoryVerification {
			return provider.sendVoice(ctx, msg)
		}
		return fmt.Errorf("luosimao voice only supports verification category")
	default:
		return provider.sendSMS(ctx, msg)
	}
}

// sendSMS sends SMS message via Luosimao API
func (provider *LuosimaoProvider) sendSMS(ctx context.Context, msg *Message) error {
	content := utils.AddSignature(msg.Content, msg.SignName)
	if msg.Category == CategoryVerification && len(msg.Mobiles) == 1 {
		return sendLuosimaoSingle(ctx, &provider.config, msg.Mobiles[0], content)
	}
	return sendLuosimaoBatch(ctx, &provider.config, msg.Mobiles, content, msg)
}

// luosimaoRequestURI 统一生成 API 请求地址，支持可选自定义域名
func luosimaoRequestURI(service, path string, override ...string) string {
	domain := service + "." + luosimaoDomain
	if len(override) > 0 && override[0] != "" {
		domain = override[0]
	}
	return "http://" + domain + path
}

// sendLuosimaoSingle sends a single SMS
//
// 单个发送接口：
//   - URL: http://sms-api.luosimao.com/v1/send.json
//   - Method: POST
//   - Content-Type: application/x-www-form-urlencoded
//   - Auth: Basic Auth (api:key-{api_key})
//
// 请求参数：
//   - mobile: 目标手机号码
//   - message: 短信内容（需包含签名）
func sendLuosimaoSingle(ctx context.Context, provider *SMSProvider, mobile, content string) error {
	defaultRequestURI := luosimaoRequestURI(luosimaoSMSService, luosimaoSMSSendPath)
	endpoint := provider.GetEndpoint(false, defaultRequestURI)
	requestBody := map[string]string{
		"mobile":  mobile,
		"message": content,
	}
	authHeader := "Basic " + utils.Base64EncodeBytes([]byte("api:key-"+provider.AppSecret))
	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Data: requestBody,
	})
	if err != nil {
		return fmt.Errorf("luosimao SMS request failed: %w", err)
	}
	return parseLuosimaoResponse(resp)
}

// sendLuosimaoBatch sends batch SMS
//
// 批量发送接口：
//   - URL: http://sms-api.luosimao.com/v1/send_batch.json
//   - Method: POST
//   - Content-Type: application/x-www-form-urlencoded
//   - Auth: Basic Auth (api:key-{api_key})
//
// 请求参数：
//   - mobile_list: 目标手机号码列表（逗号分隔）
//   - message: 短信内容（需包含签名）
//   - time: 定时发送时间（可选）
//
// 限制：
//   - 单次提交控制在10万个号码以内
//   - 批量接口专门用于大量号码的内容群发，不建议发送验证码等有时效性要求的内容
func sendLuosimaoBatch(ctx context.Context, provider *SMSProvider, mobiles []string, content string, msg *Message) error {
	defaultRequestURI := luosimaoRequestURI(luosimaoSMSService, luosimaoSMSBatchPath)
	endpoint := provider.GetEndpoint(false, defaultRequestURI)
	if len(mobiles) > 100000 {
		return fmt.Errorf("luosimao batch SMS limit exceeded: max 100,000 numbers per request")
	}
	mobileList := strings.Join(mobiles, ",")
	requestBody := map[string]string{
		"mobile_list": mobileList,
		"message":     content,
	}
	if t, ok := msg.GetExtraString("time"); ok && t != "" {
		requestBody["time"] = t // 螺丝帽批量接口支持定时发送，格式为"YYYY-MM-DD HH:MM:SS"
	}
	authHeader := "Basic " + utils.Base64EncodeBytes([]byte("api:key-"+provider.AppSecret))
	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Data: requestBody,
	})
	if err != nil {
		return fmt.Errorf("luosimao batch SMS request failed: %w", err)
	}
	return parseLuosimaoResponse(resp)
}

// 响应格式：
//
//	{
//	  "error": 0,           // 错误码，0表示成功
//	  "msg": "ok",          // 错误描述
//	  "batch_id": "...",    // 批次号（批量发送时返回）
//	  "hit": "..."          // 敏感词（error为-31时返回）
//	}
func parseLuosimaoResponse(resp []byte) error {
	var result struct {
		Error   int    `json:"error"`
		Msg     string `json:"msg"`
		BatchID string `json:"batch_id"`
		Hit     string `json:"hit"`
	}
	err := json.Unmarshal(resp, &result)
	if err != nil {
		return fmt.Errorf("failed to parse luosimao response: %w", err)
	}
	if result.Error != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.Error),
			Message:  result.Msg,
			Provider: string(ProviderTypeLuosimao),
		}
	}
	return nil
}
func (p *LuosimaoProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，单发/群发，内容需包含签名，批量最多10万条/次",
	)
	// 国际短信不支持
	capabilities.SMS.International = NewRegionCapability(false, false, nil, nil, "不支持国际短信")
	capabilities.SMS.Limits.MaxBatchSize = 100000
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
func (p *LuosimaoProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}
func (p *LuosimaoProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	default:
		return Limits{}
	}
}
func (p *LuosimaoProvider) GetName() string {
	return p.config.Name
}
func (p *LuosimaoProvider) GetType() string {
	return string(p.config.Type)
}
func (p *LuosimaoProvider) IsEnabled() bool {
	return !p.config.Disabled
}
func (p *LuosimaoProvider) GetWeight() int {
	return p.config.GetWeight()
}
func (p *LuosimaoProvider) CheckConfigured() error {
	if p.config.AppSecret == "" {
		return fmt.Errorf("luosimao SMS provider requires AppSecret (api_key)")
	}
	return nil
}

// sendVoice sends voice message via Luosimao API
func (provider *LuosimaoProvider) sendVoice(ctx context.Context, msg *Message) error {
	if len(msg.Mobiles) != 1 {
		return fmt.Errorf("luosimao voice only supports single mobile per request")
	}
	code := msg.Content
	if code == "" {
		return fmt.Errorf("luosimao voice requires code in Content field")
	}
	authHeader := "Basic " + utils.Base64EncodeBytes([]byte("api:key-"+provider.config.AppSecret))
	defaultRequestURI := luosimaoRequestURI(luosimaoVoiceService, luosimaoVoiceSendPath)
	endpoint := provider.config.GetEndpoint(true, defaultRequestURI)
	params := map[string]string{
		"mobile": msg.Mobiles[0],
		"code":   code,
	}
	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Data: params,
	})
	if err != nil {
		return fmt.Errorf("luosimao voice request failed: %w", err)
	}
	return parseLuosimaoResponse(resp)
}
