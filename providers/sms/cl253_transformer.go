package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

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
//   - 彩信/语音短信：暂不支持。
//
// 签名和营销短信的结尾是拼接在内容里的，签名本实现会自动增加。

// init automatically registers the CL253 transformer
func init() {
	RegisterTransformer(string(SubProviderCl253), &cl253Transformer{})
}

const (
	cl253DomesticEndpoint = "smssh1.253.com"
	cl253IntlEndpoint     = "intapi.253.com"
)

// cl253Transformer implements HTTPRequestTransformer for CL253 SMS
type cl253Transformer struct{}

// CanTransform checks if this transformer can handle the given message
func (t *cl253Transformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return false
	}
	return smsMsg.SubProvider == string(SubProviderCl253)
}

// Transform converts a CL253 SMS message to HTTP request specification
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
func (t *cl253Transformer) Transform(ctx context.Context, msg core.Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for CL253: %T", msg)
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}
	if smsMsg.IsIntl() {
		return t.transformIntlSMS(ctx, smsMsg, account)
	}
	return t.transformDomesticSMS(ctx, smsMsg, account)
}

// validateMessage validates the message for CL253
func (t *cl253Transformer) validateMessage(msg *Message) error {
	// 国内短信必须有签名
	if msg.SignName == "" && utils.HasSignature(msg.Content) && msg.IsDomestic() {
		return errors.New("sign name is required for CL253 SMS")
	}
	if len(msg.Mobiles) == 0 {
		return errors.New("at least one mobile number is required")
	}
	if msg.Content == "" {
		return errors.New("content is required for CL253 SMS")
	}
	if msg.IsIntl() && len(msg.Mobiles) > 1 {
		return errors.New("CL253 international SMS only supports single recipient")
	}
	return nil
}

// transformDomesticSMS transforms domestic SMS message to HTTP request
//
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
func (t *cl253Transformer) transformDomesticSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	url := "https://" + cl253DomesticEndpoint + "/msg/v1/send/json"
	params := map[string]interface{}{
		"account":  account.Key,
		"password": account.Secret,
		"msg":      utils.AddSignature(msg.Content, msg.SignName),
		"phone":    strings.Join(msg.Mobiles, ","),
	}

	// 处理统一的接口字段 - 适配到CL253特定的key
	if report := msg.GetExtraStringOrDefault(cl253Report, ""); report != "" {
		params["report"] = report
	}
	if msg.CallbackURL != "" {
		params["callbackUrl"] = msg.CallbackURL // CL253使用 callbackUrl
	}
	if msg.UID != "" {
		params["uid"] = msg.UID // CL253使用 uid
	}
	if msg.Extend != "" {
		params["extend"] = msg.Extend // CL253使用 extend
	}

	// 处理发送时间 - CL253使用 sendtime 字段
	if msg.ScheduledAt != nil {
		// CL253使用 yyyyMMddHHmm 格式
		params["sendtime"] = msg.ScheduledAt.Format("200601021504")
	}

	body, _ := json.Marshal(params)
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      url,
		Headers:  t.buildHeaders(nil),
		Body:     body,
		BodyType: "json",
		Timeout:  30 * time.Second,
	}
	return reqSpec, t.handleCl253Response, nil
}

// transformIntlSMS transforms international SMS message to HTTP request
//
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
func (t *cl253Transformer) transformIntlSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	url := "https://" + cl253IntlEndpoint + "/send/sms"

	// 手机号码，格式(区号+手机号码)，例如：8615800000000，其中 86 为中国的区号， 区号前不使用 00 开头,15800000000 为接收短信的真实手机号码。5-20 位
	params := map[string]interface{}{
		"account":  account.Key,
		"password": account.Secret,
		"mobile":   fmt.Sprintf("%d%s", msg.RegionCode, msg.Mobiles[0]),
		"msg":      utils.AddSignature(msg.Content, msg.SignName),
	}

	// 处理统一的接口字段 - 适配到CL253国际短信特定的key
	if senderId := msg.GetExtraStringOrDefault(cl253SenderID, ""); senderId != "" {
		params["senderId"] = senderId
	}
	if templateId := msg.GetExtraStringOrDefault(cl253TemplateID, ""); templateId != "" {
		params["templateId"] = templateId
	}
	if msg.UID != "" {
		params["uid"] = msg.UID // CL253使用 uid
	}
	if tdFlag := msg.GetExtraStringOrDefault(cl253TDFlag, ""); tdFlag != "" {
		params["tdFlag"] = tdFlag
	}

	body, _ := json.Marshal(params)
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      url,
		Headers:  t.buildHeaders(nil),
		Body:     body,
		BodyType: "json",
		Timeout:  30 * time.Second,
	}
	return reqSpec, t.handleCl253Response, nil
}

// buildHeaders 构建请求头，支持用户自定义 header 合并，默认加 content-type
func (t *cl253Transformer) buildHeaders(userHeaders map[string]string) map[string]string {
	headers := map[string]string{
		"content-type": "application/json",
	}
	for k, v := range userHeaders {
		headers[strings.ToLower(k)] = v
	}
	return headers
}

// handleCl253Response 处理 CL253 API 响应
func (t *cl253Transformer) handleCl253Response(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}
	var result struct {
		Code     string `json:"code"`
		MsgId    string `json:"msgId"`
		RespTime string `json:"time"`
		ErrorMsg string `json:"errorMsg"`
		Message  string `json:"message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse CL253 response: %w", err)
	}
	if result.Code != "0" {
		return &SMSError{
			Code:     result.Code,
			Message:  result.ErrorMsg + result.Message,
			Provider: string(SubProviderCl253),
		}
	}
	return nil
}
