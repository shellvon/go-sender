package sms

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	smsbaoDefaultBaseURI = "https://api.smsbao.com"
	maxMobilesPerRequest = 99
)

// smsbaoTransformer implements HTTPRequestTransformer for Smsbao
// 统一风格实现

type smsbaoTransformer struct {
	*BaseTransformer
}

func newSmsbaoTransformer() *smsbaoTransformer {
	transformer := &smsbaoTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(core.ProviderTypeSMS),
		string(SubProviderSmsbao),
		nil,
		WithSMSHandler(transformer.transformSMS),
		WithVoiceHandler(transformer.transformVoice),
	)
	return transformer
}

// init 自动注册 Smsbao transformer.
func init() {
	RegisterTransformer(string(SubProviderSmsbao), newSmsbaoTransformer())
}

// transformSMS transforms SMS message to HTTP request
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//   - 国际短信: https://www.smsbao.com/openapi/299.html
func (t *smsbaoTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if account == nil || account.APIKey == "" || account.APISecret == "" {
		return nil, nil, NewProviderError(
			string(ProviderTypeSmsbao),
			"AUTH_ERROR",
			"smsbao account Key(username) and Secret(password) are required",
		)
	}
	mobiles := strings.Join(msg.Mobiles, ",")
	content := utils.AddSignature(msg.Content, msg.SignName)

	// 构建查询参数
	queryParams := url.Values{
		"u": {account.APIKey},
		"p": {utils.MD5Hex(account.APISecret)},
		"m": {mobiles},
		"c": {content},
	}
	if prod := msg.GetExtraStringOrDefault(smsbaoProductIDKey, ""); prod != "" {
		queryParams.Set("g", prod)
	}

	var apiPath string
	if msg.IsIntl() {
		apiPath = "/wsms"
	} else {
		apiPath = "/sms"
	}

	return &core.HTTPRequestSpec{
		Method:      http.MethodGet,
		URL:         fmt.Sprintf("%s%s", smsbaoDefaultBaseURI, apiPath),
		QueryParams: queryParams,
	}, t.handleSMSBaoResponse, nil
}

// transformVoice transforms voice message to HTTP request
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
//
// 能力说明:
//   - 语音验证码：仅支持国内、仅验证码类型、仅单号码。
func (t *smsbaoTransformer) transformVoice(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 检查语音短信的限制
	if msg.IsIntl() {
		return nil, nil, NewProviderError(
			string(ProviderTypeSmsbao),
			"UNSUPPORTED_COUNTRY",
			"voice sms only supports domestic mobile",
		)
	}
	if len(msg.Mobiles) != 1 {
		return nil, nil, NewProviderError(
			string(ProviderTypeSmsbao),
			"INVALID_MOBILE_NUMBER",
			fmt.Sprintf("smsbao voice only supports single mobile, got %d", len(msg.Mobiles)),
		)
	}
	if len(msg.Mobiles[0]) != 11 || msg.Mobiles[0][0] != '1' {
		return nil, nil, NewProviderError(
			string(ProviderTypeSmsbao),
			"INVALID_MOBILE_FORMAT",
			"only support domestic mobile for voice sms",
		)
	}

	if account == nil || account.APIKey == "" || account.APISecret == "" {
		return nil, nil, NewProviderError(
			string(ProviderTypeSmsbao),
			"AUTH_ERROR",
			"smsbao account Key(username) and Secret(password) are required",
		)
	}

	// 构建查询参数
	queryParams := url.Values{
		"u": {account.APIKey},
		"p": {utils.MD5Hex(account.APISecret)},
		"m": {msg.Mobiles[0]},
		"c": {msg.Content},
	}

	return &core.HTTPRequestSpec{
		Method:      http.MethodGet,
		URL:         fmt.Sprintf("%s/voice", smsbaoDefaultBaseURI),
		QueryParams: queryParams,
	}, t.handleSMSBaoResponse, nil
}

// handleSMSBaoResponse 处理短信宝 API 响应.
func (t *smsbaoTransformer) handleSMSBaoResponse(resp *http.Response) error {
	body, _, err := utils.ReadAndClose(resp)
	if err != nil {
		return NewProviderError(string(ProviderTypeSmsbao), "READ_ERROR", err.Error())
	}
	var smsBaoErrorMap = map[string]string{
		"30": "password error",
		"40": "account does not exist",
		"41": "insufficient balance",
		"42": "account expired",
		"43": "IP address restriction",
		"50": "content contains sensitive words",
		"51": "incorrect mobile number",
	}
	code := string(body)
	if code != "0" {
		return NewProviderError(string(ProviderTypeSmsbao), code, smsBaoErrorMap[code])
	}
	return nil
}
