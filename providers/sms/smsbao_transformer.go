package sms

import (
	"context"
	//nolint:gosec // compatibility with legacy system, not for security
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// smsbaoTransformer implements HTTPRequestTransformer for Smsbao SMS.
// It supports sending text message and voice message.
//
// Reference:
//   - Official Website: https://www.smsbao.com
//   - API Docs: https://www.smsbao.com/openapi
//   - SMS API(Domestic): https://www.smsbao.com/openapi/213.html
//   - SMS API(International): https://www.smsbao.com/openapi/299.html
//   - Voice API: https://www.smsbao.com/openapi/214.html

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
		string(SubProviderSmsbao),
		nil,
		HTTPOptions(nil),
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
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	if account == nil || account.APIKey == "" || account.APISecret == "" {
		return nil, nil, NewProviderError(
			string(SubProviderSmsbao),
			"AUTH_ERROR",
			"smsbao account Key(username) and Secret(password) are required",
		)
	}
	mobiles := strings.Join(msg.Mobiles, ",")
	content := utils.AddSignature(msg.Content, msg.SignName)

	// 构建查询参数
	queryParams := url.Values{
		"u": {account.APIKey},
		"p": {utils.HashHex(md5.New, []byte(account.APISecret))},
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
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	subProvider := string(SubProviderSmsbao)
	// 检查语音短信的限制
	if msg.IsIntl() {
		return nil, nil, NewProviderError(
			subProvider,
			"UNSUPPORTED_COUNTRY",
			"voice sms only supports domestic mobile",
		)
	}
	if len(msg.Mobiles) != 1 {
		return nil, nil, NewProviderError(
			subProvider,
			"INVALID_MOBILE_NUMBER",
			fmt.Sprintf("smsbao voice only supports single mobile, got %d", len(msg.Mobiles)),
		)
	}
	if len(msg.Mobiles[0]) != 11 || msg.Mobiles[0][0] != '1' {
		return nil, nil, NewProviderError(
			subProvider,
			"INVALID_MOBILE_FORMAT",
			"only support domestic mobile for voice sms",
		)
	}

	if account == nil || account.APIKey == "" || account.APISecret == "" {
		return nil, nil, NewProviderError(
			subProvider,
			"AUTH_ERROR",
			"smsbao account Key(username) and Secret(password) are required",
		)
	}

	// 构建查询参数
	queryParams := url.Values{
		"u": {account.APIKey},
		"p": {utils.HashHex(md5.New, []byte(account.APISecret))},
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
func (t *smsbaoTransformer) handleSMSBaoResponse(result *core.SendResult) error {
	subProvider := string(SubProviderSmsbao)
	var smsBaoErrorMap = map[string]string{
		"30": "password error",
		"40": "account does not exist",
		"41": "insufficient balance",
		"42": "account expired",
		"43": "IP address restriction",
		"50": "content contains sensitive words",
		"51": "incorrect mobile number",
	}
	code := string(result.Body)
	if code != "0" {
		return NewProviderError(subProvider, code, smsBaoErrorMap[code])
	}
	return nil
}
