package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Huawei Cloud / 华为云
// @Website: https://www.huaweicloud.com
// @APIDoc: https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 官方文档:
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 仅支持国内/国际短信，不支持彩信和语音。
//
// transformer 仅支持 text（模板短信）类型。

const (
	huaweiEndpoint     = "https://api.rtc.huaweicloud.com:10443"
	huaweiEndpointPath = "/sms/batchSendSms/v1"
	huaweiSuccessCode  = "000000"
	huaweiTimeout      = 30 * time.Second
)

// huaweiTransformer implements HTTPRequestTransformer for Huawei Cloud SMS.
type huaweiTransformer struct{}

// init automatically registers the Huawei transformer.
func init() {
	RegisterTransformer(string(SubProviderHuawei), &huaweiTransformer{})
}

// CanTransform checks if this transformer can handle the given message.
func (t *huaweiTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return false
	}
	return smsMsg.SubProvider == string(SubProviderHuawei)
}

// Transform converts a Huawei SMS message to HTTP request specification
//
// - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
func (t *huaweiTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for Huawei: %T", msg)
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}
	return t.transformSMS(smsMsg, account)
}

// validateMessage validates the message for Huawei.
func (t *huaweiTransformer) validateMessage(msg *Message) error {
	if msg.TemplateID == "" {
		return errors.New("templateId is required for Huawei SMS")
	}
	if len(msg.Mobiles) == 0 {
		return errors.New("at least one mobile number is required")
	}
	if msg.SignName == "" {
		return errors.New("sign name is required for Huawei SMS")
	}
	return nil
}

// transformSMS transforms SMS message to HTTP request
//
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 如果"to"参数携带的号码中包含除数字和+之外的其他字符，则无法向该参数携带的所有号码发送短信。如果"to"参数携带的所有号码只包含数字和+，
// 但部分号码不符合号码规则要求，则在响应消息中会通过状态码标识发送失败的号码，不影响其他正常号码的短信发送。号码之间以英文逗号分隔，
// 每个号码最大长度为21位，最多允许携带500个号码。如果携带超过500个号码，则全部号码都会发送失败。
func (t *huaweiTransformer) transformSMS(
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 格式化手机号
	formattedMobiles := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		formattedMobiles[i] = t.formatHuaweiPhoneNumber(mobile, msg.RegionCode)
	}
	paras := ""
	if len(msg.ParamsOrder) > 0 {
		b, _ := json.Marshal(msg.ParamsOrder)
		paras = string(b)
	}
	// 构建参数，注释与每个参数紧密对应
	params := url.Values{}
	params.Set("from", msg.GetExtraStringOrDefault(huaweiFrom, account.From))
	params.Set("to", strings.Join(formattedMobiles, ","))
	params.Set("templateId", msg.TemplateID)
	params.Set("templateParas", paras)
	// 签名名称，必须是已审核通过的，与模板类型一致的签名名称。
	// 仅中国大陆短信可携带此参数。
	// 仅在templateId指定的模板类型为通用模板时生效且必填，用于指定在通用模板短信内容前面补充的签名
	if msg.IsDomestic() {
		params.Set("signature", msg.SignName)
	}
	// 客户的回调地址，用于接收短信状态报告，如：http://my.com/receiveSMSReport。
	// 如果设置了该字段，则该消息的状态报告将通过"接收状态报告"接口直接通知客户。
	// 如果未设置该字段，则短信平台收到运营商短信中心返回的状态报告不会推送给客户，该状态报告将在短信平台中保存1个小时，超时后系统会自动删除。
	// 回调地址推荐使用域名。
	if msg.CallbackURL != "" {
		params.Set("statusCallback", msg.CallbackURL)
	}
	// 扩展参数，在状态报告中会原样返回。
	// 不允许赋空值，不允许携带以下字符："{","}"（即大括号）。
	if msg.Extend != "" {
		params.Set("extend", msg.Extend)
	} else if extend := msg.GetExtraStringOrDefault(huaweiExtend, ""); extend != "" {
		params.Set("extend", extend)
	}

	body := []byte(params.Encode())

	appKey := account.Key
	appSecret := account.Secret

	endpoint := t.getHuaweiEndpoint(account.Endpoint)
	if msg.IsIntl() {
		endpoint = t.getHuaweiEndpoint(account.IntlEndpoint)
	}
	headers := t.buildHeaders(appKey, appSecret)
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      endpoint,
		Headers:  headers,
		Body:     body,
		BodyType: "form",
	}
	return reqSpec, t.handleHuaweiResponse, nil
}

// formatHuaweiPhoneNumber formats phone number for Huawei Cloud API
// 华为云要求：标准号码格式为：+{国家码}{地区码}{终端号码}.
func (t *huaweiTransformer) formatHuaweiPhoneNumber(mobile string, regionCode int) string {
	if regionCode == 0 {
		regionCode = 86
	}
	return fmt.Sprintf("+%d%s", regionCode, mobile)
}

// getHuaweiEndpoint returns the full endpoint URL.
func (t *huaweiTransformer) getHuaweiEndpoint(endpoint string) string {
	if endpoint == "" {
		endpoint = huaweiEndpoint
	}
	return strings.TrimRight(endpoint, "/") + huaweiEndpointPath
}

// buildHeaders 构建华为云短信请求头.
func (t *huaweiTransformer) buildHeaders(appKey, appSecret string) map[string]string {
	return map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "WSSE realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\"",
		"X-WSSE":        t.buildHuaweiWsseHeader(appKey, appSecret),
	}
}

// buildHuaweiWsseHeader 构建 X-WSSE 认证头.
func (t *huaweiTransformer) buildHuaweiWsseHeader(appKey, appSecret string) string {
	now := time.Now().UTC().Format(time.RFC3339)
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	passwordDigest := utils.Base64EncodeBytes(utils.SHA256Sum([]byte(nonce + now + appSecret)))
	return fmt.Sprintf(
		"UsernameToken Username=\"%s\",PasswordDigest=\"%s\",Nonce=\"%s\",Created=\"%s\"",
		appKey, passwordDigest, nonce, now,
	)
}

// handleHuaweiResponse 处理华为云短信 API 响应.
func (t *huaweiTransformer) handleHuaweiResponse(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}
	var result struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse huawei response: %w", err)
	}
	if result.Code != huaweiSuccessCode {
		return &Error{
			Code:     result.Code,
			Message:  result.Description,
			Provider: string(SubProviderHuawei),
		}
	}
	return nil
}
