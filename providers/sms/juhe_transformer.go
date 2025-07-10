package sms

import (
	"context"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/core"
)

// juheTransformer implements HTTPRequestTransformer for Juhe SMS.
// It supports sending text message, voice message, and mms message.
//
// Reference:
//   - Official Website: https://www.juhe.cn
//   - API Docs: https://www.juhe.cn/docs/api/id/54
//   - SMS API(Domestic): https://www.juhe.cn/docs/api/id/54
//   - SMS API(International): https://www.juhe.cn/docs/api/id/357
//   - MMS API: https://www.juhe.cn/docs/api/id/363

const (
	juheDefaultBaseURI = "https//v.juhe.cn"

	// API路径常量.
	juheDomesticAPIPath = "/sms/send"
	juheIntlAPIPath     = "/smsInternational/send"
	juheMMSAPIPath      = "/caixinv2/send"
)

type juheTransformer struct {
	*BaseTransformer
}

func newJuheTransformer() *juheTransformer {
	transformer := &juheTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(core.ProviderTypeSMS),
		&core.ResponseHandlerConfig{
			BodyType:  core.BodyTypeJSON,
			CheckBody: true,
			Path:      "error_code",
			Expect:    "0",
			Mode:      core.MatchEq,
		},
		nil,
		WithSMSHandler(transformer.transformSMS),
		WithMMSHandler(transformer.transformMMS),
	)
	return transformer
}

func init() {
	RegisterTransformer(string(SubProviderJuhe), newJuheTransformer())
}

// transformSMS 处理短信（国内/国际）
//   - 国内短信 API: https://www.juhe.cn/docs/api/id/54
//   - 国际短信 API: https://www.juhe.cn/docs/api/id/357
func (t *juheTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	params := url.Values{}
	params.Set("mobile", msg.Mobiles[0])
	params.Set("tpl_id", msg.TemplateID)
	params.Set("tpl_value", t.buildTemplateValue(msg.TemplateParams))
	params.Set("key", account.APIKey)

	// 根据是否为国际短信设置不同的参数和URL
	var apiPath string
	if msg.IsIntl() {
		// 国际短信特有参数
		params.Set("areaNum", strconv.Itoa(msg.RegionCode))
		apiPath = juheIntlAPIPath
	} else {
		// 国内短信特有参数
		params.Set("ext", msg.Extend)
		apiPath = juheDomesticAPIPath
	}

	body := []byte(params.Encode())
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      juheDefaultBaseURI + apiPath,
		Body:     body,
		BodyType: core.BodyTypeForm,
	}
	return reqSpec, nil, nil
}

// transformMMS 处理彩信/视频短信
//   - API: https://www.juhe.cn/docs/api/id/363
func (t *juheTransformer) transformMMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	params := url.Values{}
	params.Set("mobile", strings.Join(msg.Mobiles, ","))
	params.Set("tpl_id", msg.TemplateID)
	params.Set("key", account.APIKey)

	body := []byte(params.Encode())
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      juheDefaultBaseURI + juheMMSAPIPath,
		Body:     body,
		BodyType: core.BodyTypeForm,
	}
	return reqSpec, nil, nil
}

// 如果您的模板里面有变量则需要提交此参数,如:#code#=123456,参数需要urlencode
// -> https://www.juhe.cn/docs/api/id/357
func (t *juheTransformer) buildTemplateValue(params map[string]string) string {
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
