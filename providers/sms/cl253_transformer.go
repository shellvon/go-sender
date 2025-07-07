package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

// init automatically registers the CL253 transformer.
func init() {
	RegisterTransformer(string(SubProviderCl253), newCL253Transformer())
}

const (
	cl253DomesticEndpoint      = "smssh1.253.com"
	cl253IntlSingaporeEndpoint = "sgap.253.com"
	cl253IntlShanghaiEndpoint  = "intapi.253.com"
	cl253DefaultRegion         = "sh" // 默认使用上海节点

	// API路径常量.
	cl253DomesticAPIPath = "/msg/v1/send/json"
	cl253IntlAPIPath     = "/send/sms"
)

// cl253Transformer implements HTTPRequestTransformer for CL253 SMS.
type cl253Transformer struct {
	*BaseTransformer
}

func newCL253Transformer() *cl253Transformer {
	transformer := &cl253Transformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(core.ProviderTypeSMS),
		string(SubProviderCl253),
		&core.ResponseHandlerConfig{
			SuccessField:      "code",
			SuccessValue:      "0",
			ErrorCodeField:    "code",
			ErrorMessageField: "message",
			ErrorField:        "errorMsg",
			MessageField:      "message",
			ResponseType:      core.BodyTypeJSON,
			ValidateResponse:  true,
		},
		WithBeforeHook(func(_ context.Context, msg *Message, _ *Account) error {
			return transformer.validateMessage(msg)
		}),
		WithSMSHandler(transformer.transformSMS),
	)
	return transformer
}

// CanTransform checks if this transformer can handle the given message.
// Inherited CanTransform from *BaseTransformer.

// Transform converts a CL253 SMS message to HTTP request specification
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
// Transform method is inherited from *BaseTransformer; no need to redeclare.

// validateMessage validates the message for CL253.
func (t *cl253Transformer) validateMessage(msg *Message) error {
	// 国内短信必须有签名
	if msg.SignName == "" && utils.HasSignature(msg.Content) && msg.IsDomestic() {
		return NewProviderError(string(SubProviderCl253), "MISSING_SIGNATURE", "sign name is required for CL253 SMS")
	}
	if len(msg.Mobiles) == 0 {
		return NewProviderError(string(SubProviderCl253), "MISSING_MOBILE", "at least one mobile number is required")
	}
	if msg.Content == "" {
		return NewProviderError(string(SubProviderCl253), "MISSING_CONTENT", "content is required for CL253 SMS")
	}
	if msg.IsIntl() && len(msg.Mobiles) > 1 {
		return NewProviderError(
			string(SubProviderCl253),
			"UNSUPPORTED_RECIPIENTS",
			"CL253 international SMS only supports single recipient",
		)
	}
	return nil
}

// transformSMS transforms SMS message to HTTP request
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
func (t *cl253Transformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 构建基础参数
	params := map[string]interface{}{
		"account":     account.APIKey,
		"password":    account.APISecret,
		"msg":         utils.AddSignature(msg.Content, msg.SignName),
		"report":      msg.GetExtraStringOrDefault(cl253ReportKey, ""),
		"callbackUrl": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
		"uid":         msg.UID,
		"extend":      msg.Extend,
	}

	// 添加国际短信特有参数
	if msg.IsIntl() {
		params["tdFlag"] = msg.GetExtraIntOrDefault(cl253TDFlagKey, 0)
		params["templateId"] = msg.TemplateID
		params["senderId"] = msg.GetExtraStringOrDefault(cl253SenderIDKey, "")
		params["mobile"] = fmt.Sprintf("%d%s", msg.RegionCode, msg.Mobiles[0])
	} else {
		params["phone"] = strings.Join(msg.Mobiles, ",")
	}

	// 处理发送时间 - CL253使用 sendtime 字段
	if msg.ScheduledAt != nil {
		// CL253使用 yyyyMMddHHmm 格式
		params["sendtime"] = msg.ScheduledAt.Format("200601021504")
	}

	body, _ := json.Marshal(params)
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      t.buildRequestURI(msg, account),
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}

	return reqSpec, nil, nil
}

// buildRequestURI 根据区域和是否国际获取最终的短信请求API.
func (t *cl253Transformer) buildRequestURI(msg *Message, account *Account) string {
	if msg.IsIntl() {
		// 国际短信
		ep := cl253IntlShanghaiEndpoint
		if account.Region != cl253DefaultRegion {
			ep = cl253IntlSingaporeEndpoint
		}
		return fmt.Sprintf("https://%s%s", ep, cl253IntlAPIPath)
	}
	// 国内短信
	return fmt.Sprintf("https://%s%s", cl253DomesticEndpoint, cl253DomesticAPIPath)
}
