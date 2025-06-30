package sms

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Smsbao / 短信宝
// @Website: https://www.smsbao.com
// @APIDoc: https://www.smsbao.com/openapi
//
// 官方文档:
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//   - 国际短信: https://www.smsbao.com/openapi/299.html
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
//
// 能力说明:
//   - 国内短信：支持单发和群发，最多99个号码/次。
//   - 国际短信：支持单发和群发，最多99个号码/次。
//   - 语音验证码：仅支持国内、仅验证码类型、仅单号码。
//
// transformer 支持文本短信（国内/国际）和语音验证码。

const (
	smsbaoDefaultEndpoint = "api.smsbao.com"
)

// smsbaoTransformer implements HTTPRequestTransformer for Smsbao
// 统一风格实现

type smsbaoTransformer struct{}

// newSMSBaoTransformer creates a new Smsbao transformer
func newSMSBaoTransformer() core.HTTPTransformer[*core.Account] {
	return &smsbaoTransformer{}
}

// init 自动注册 Smsbao transformer
func init() {
	RegisterTransformer(string(SubProviderSmsbao), &smsbaoTransformer{})
}

// CanTransform 判断是否为短信宝消息
func (t *smsbaoTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	return ok && smsMsg.SubProvider == string(SubProviderSmsbao)
}

// Transform 构造短信宝 HTTPRequestSpec
func (t *smsbaoTransformer) Transform(ctx context.Context, msg core.Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for smsbao: %T", msg)
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, err
	}
	switch smsMsg.Type {
	case SMSText:
		return t.transformTextSMS(ctx, smsMsg, account)
	case Voice:
		return t.transformVoiceSMS(ctx, smsMsg, account)
	default:
		return nil, nil, fmt.Errorf("unsupported smsbao message type: %s", smsMsg.Type)
	}
}

// validateMessage 校验参数
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//   - 国际短信: https://www.smsbao.com/openapi/299.html
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
func (t *smsbaoTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("mobiles is required")
	}
	if msg.Content == "" {
		return errors.New("content is required")
	}
	if len(msg.Mobiles) > 99 {
		return errors.New("smsbao: at most 99 mobiles per request")
	}
	if msg.Type == Voice {
		if msg.IsIntl() {
			return errors.New("voice sms only supports domestic mobile")
		}
		if len(msg.Mobiles) != 1 {
			return errors.New("voice sms only supports single domestic mobile")
		}
		if len(msg.Mobiles[0]) != 11 || msg.Mobiles[0][0] != '1' {
			return errors.New("only support domestic mobile for voice sms")
		}
	}
	return nil
}

// getBaseDomain 获取基础域名，优先使用account配置的endpoint
func (t *smsbaoTransformer) getBaseDomain(account *core.Account, isIntl bool) string {
	if isIntl && account.IntlEndpoint != "" {
		return account.IntlEndpoint
	}
	if account.Endpoint != "" {
		return account.Endpoint
	}
	return smsbaoDefaultEndpoint
}

// transformTextSMS 构造文本短信 HTTP 请求
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//   - 国际短信: https://www.smsbao.com/openapi/299.html
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
func (t *smsbaoTransformer) transformTextSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if account == nil || account.Key == "" || account.Secret == "" {
		return nil, nil, errors.New("smsbao account Key(username) and Secret(password) are required")
	}
	mobiles := strings.Join(msg.Mobiles, ",")
	content := utils.AddSignature(msg.Content, msg.SignName)

	// 构建查询参数
	queryParams := map[string]string{
		"u": account.Key,
		"p": utils.MD5Hex(account.Secret),
		"m": mobiles,
		"c": content,
	}

	// 当客户使用专用通道产品时，需要指定产品ID
	// 产品ID可在短信宝后台或联系客服获得,不填则默认使用通用短信产品
	// 文档
	//   - 国内短信: https://www.smsbao.com/openapi/213.html
	//
	// 对于短信宝而言，TemplateID 和 ProductID 是同一个概念
	// 且可以从账号的配置中from字段获取
	if msg.IsDomestic() {
		productID := utils.DefaultStringIfEmpty(msg.TemplateID, account.From)
		if productID != "" {
			queryParams["g"] = productID
		}
	}

	var apiPath string
	if msg.IsIntl() {
		apiPath = "/wsms"
	} else {
		apiPath = "/sms"
	}
	baseDomain := t.getBaseDomain(account, msg.IsIntl())

	return &core.HTTPRequestSpec{
		Method:      "GET",
		URL:         "https://" + baseDomain + apiPath,
		QueryParams: queryParams,
		Timeout:     10 * time.Second,
	}, handleSMSBaoResponse, nil
}

// transformVoiceSMS 构造语音验证码 HTTP 请求
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
//
// 能力说明:
//   - 语音验证码：仅支持国内、仅验证码类型、仅单号码。
func (t *smsbaoTransformer) transformVoiceSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if account == nil || account.Key == "" || account.Secret == "" {
		return nil, nil, errors.New("smsbao account Key(username) and Secret(password) are required")
	}

	// 构建查询参数
	queryParams := map[string]string{
		"u": account.Key,
		"p": utils.MD5Hex(account.Secret),
		"m": msg.Mobiles[0],
		"c": msg.Content,
	}

	baseDomain := t.getBaseDomain(account, false) // 语音仅支持国内
	return &core.HTTPRequestSpec{
		Method:      "GET",
		URL:         "http://" + baseDomain + "/voice",
		QueryParams: queryParams,
		Timeout:     10 * time.Second,
	}, handleSMSBaoResponse, nil
}

// handleSMSBaoResponse 处理短信宝 API 响应
func handleSMSBaoResponse(statusCode int, body []byte) error {
	if statusCode != 200 {
		return fmt.Errorf("smsbao API returned non-OK status: %d", statusCode)
	}

	resp := string(body)
	if resp == "0" {
		return nil // 成功
	}

	// 错误码映射表
	errorMessages := map[string]string{
		"30": "password error",
		"40": "account does not exist",
		"41": "insufficient balance",
		"42": "account expired",
		"43": "IP address restriction",
		"50": "content contains sensitive words",
		"51": "incorrect mobile number",
	}

	message, exists := errorMessages[resp]
	if !exists {
		message = fmt.Sprintf("smsbao unknown error: %s", resp)
	}

	return &SMSError{
		Code:     resp,
		Message:  message,
		Provider: string(SubProviderSmsbao),
	}
}
