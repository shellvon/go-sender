package sms

// @ProviderName: Volcengine / 火山引擎
// @Website: https://www.volcengine.com
// @APIDoc: https://www.volcengine.com/docs/63933
//
// # VolcProvider implements SMSProviderInterface for Volcengine SMS
//
// 官方文档:
//   - 短信API文档: https://www.volcengine.com/docs/63933
//   - 短信发送API: https://www.volcengine.com/docs/6361/67380
//   - 签名认证文档: https://www.volcengine.com/docs/6361/1205061
//
// 能力说明:
//   - 国内短信：支持单发和群发，需模板ID。
//   - 国际短信：暂不支持。
//   - 彩信/语音：暂不支持。
//
// 注意：仅支持国内手机号码，不支持国际号码。
import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shellvon/go-sender/utils"
)

// VolcResponse represents the response from Volcengine SMS API
// API Documentation: https://www.volcengine.com/docs/6361/67380
type VolcResponse struct {
	ResponseMetadata struct {
		RequestId string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Error     *struct {
			CodeN   int    `json:"CodeN"`
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error,omitempty"`
	} `json:"ResponseMetadata"`
	Result struct {
		MessageID string `json:"MessageID"`
	} `json:"Result"`
}

const (
	volcEndpoint = "sms.volcengineapi.com"
)

type VolcProvider struct {
	config SMSProvider
}

// NewVolcProvider creates a new Volcengine SMS provider
func NewVolcProvider(config SMSProvider) *VolcProvider {
	return &VolcProvider{config: config}
}

// Send sends an SMS message via Volcengine
func (provider *VolcProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	if msg.Type != SMSText {
		return NewUnsupportedMessageTypeError(string(ProviderTypeVolc), msg.Type.String(), msg.Category.String())
	}
	return provider.sendSMS(ctx, msg)
}

// sendSMS sends SMS message via Volcengine API
func (provider *VolcProvider) sendSMS(ctx context.Context, msg *Message) error {
	if msg.IsIntl() {
		return NewUnsupportedInternationalError(string(ProviderTypeSmsbao), "sendSMS")
	}
	body := map[string]interface{}{
		"SmsAccount":   provider.config.AppID,
		"Sign":         msg.SignName,
		"TemplateID":   msg.TemplateID,
		"PhoneNumbers": strings.Join(msg.Mobiles, ","),
	}
	if len(msg.TemplateParams) > 0 {
		body["TemplateParam"] = utils.ToJSONString(msg.TemplateParams)
	}
	if tag, ok := msg.GetExtraString("tag"); ok && tag != "" {
		body["Tag"] = tag
	}
	bodyJSON, _ := json.Marshal(body)
	headers := buildVolcHeaders(&provider.config, bodyJSON)

	endpoint := provider.config.GetEndpoint(false, volcEndpoint)
	url := "https://" + endpoint + "/?Action=SendSms&Version=2020-01-01"

	resp, _, err := utils.DoRequest(ctx, url, utils.RequestOptions{
		Method:  http.MethodPost,
		Headers: headers,
		Raw:     bodyJSON,
	})
	if err != nil {
		return fmt.Errorf("volcengine SMS request failed: %w", err)
	}
	var result VolcResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse volcengine response: %w", err)
	}
	if result.ResponseMetadata.Error != nil {
		return &SMSError{
			Code:     result.ResponseMetadata.Error.Code,
			Message:  result.ResponseMetadata.Error.Message,
			Provider: string(ProviderTypeVolc),
		}
	}
	return nil
}

// buildVolcHeaders 构建火山引擎TOP网关签名头
// 签名文档: https://www.volcengine.com/docs/6361/1205061
func buildVolcHeaders(provider *SMSProvider, body []byte) map[string]string {
	ak := provider.AppID
	sk := provider.AppSecret
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	xDate := timestamp[:10]
	headers := map[string]string{
		"Content-Type": "application/json;charset=utf-8",
		// TODO: ?? 这里需要根据 endpoint 来设置 Host吧？
		"Host":   "sms.volcengineapi.com",
		"X-Date": timestamp,
	}
	canonicalHeaders := "content-type:application/json;charset=utf-8\nhost:sms.volcengineapi.com\nx-date:" + timestamp + "\n"
	canonicalRequest := "POST\n/\n\n" + canonicalHeaders + "\n" + string(body)
	stringToSign := "HMAC-SHA256\n" + xDate + "\n" + utils.SHA256Hex([]byte(canonicalRequest))
	signingKey := utils.HMACSHA256([]byte(sk), []byte(xDate))
	signature := utils.HMACSHA256(signingKey, []byte(stringToSign))
	signatureBase64 := utils.Base64EncodeBytes(signature)
	authHeader := "HMAC-SHA256 Credential=" + ak + ", SignedHeaders=content-type;host;x-date, Signature=" + signatureBase64
	headers["Authorization"] = authHeader
	return headers
}
func hashSHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
func (p *VolcProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内短信支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，单发/群发，不支持国际号码",
	)
	// 国际短信不支持
	capabilities.SMS.International = NewRegionCapability(false, false, nil, nil, "不支持国际短信")
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
func (p *VolcProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}
func (p *VolcProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	default:
		return Limits{}
	}
}
func (p *VolcProvider) GetName() string {
	return p.config.Name
}
func (p *VolcProvider) GetType() string {
	return string(p.config.Type)
}
func (p *VolcProvider) IsEnabled() bool {
	return !p.config.Disabled
}
func (p *VolcProvider) GetWeight() int {
	return p.config.GetWeight()
}
func (p *VolcProvider) CheckConfigured() error {
	if p.config.AppID == "" || p.config.AppSecret == "" {
		return fmt.Errorf("volcengine provider requires AppID and AppSecret")
	}
	return nil
}
