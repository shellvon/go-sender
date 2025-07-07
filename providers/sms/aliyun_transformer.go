//nolint:depguard // intentional use of math/rand for compatibility or legacy reasons
package sms

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Aliyun / 阿里云
// @Website: https://www.aliyun.com
//
// 官方文档:
// - 短信模板即具体发送的短信内容，模板类型支持验证码、通知短信和推广短信。模板由模板变量和模板内容构成，您需要遵守模板内容规范和变量规范。
// - SMS API(国内/国外/单发/群发): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
//
// 阿里云支持多媒体彩信:
//   - 多媒体彩信: https://help.aliyun.com/zh/sms/user-guide/what-is-multimedia-sms
//   - 定价: 卡片短信默认定价是0.2元/条,数字短信默认定价为0.4元/条
//     https://help.aliyun.com/zh/sms/user-guide/multimedia-sms-pricing
//
// 阿里云支持语音短信:
//   - 语音API文档: https://help.aliyun.com/zh/vms/getting-started/through-the-api-or-sdk-using-voice-notification-or-audio-captcha
//   - 语音验证码API: https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - 语音通知API: https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
//   - 语音短信支持国内单发，支持验证码和通知类型，需开通语音服务。
// 对于目前，使用语音发送短信时，当发送验证码时，会使用 TTS 接口，当发送通知时，会使用 Voice 接口。

// init automatically registers the Aliyun transformer.
func init() {
	RegisterTransformer(string(SubProviderAliyun), newAliyunTransformer())
}

const (
	aliyunDefaultSmsEndpoint      = "dysmsapi.aliyuncs.com"
	aliyunDefaultVoiceEndpoint    = "dyvmsapi.aliyuncs.com"
	aliyunDefaultRegion           = "cn-hangzhou"
	aliyunActionSendSms           = "SendSms"
	aliyunActionSingleCallByVoice = "SingleCallByVoice"
	aliyunActionSingleCallByTts   = "SingleCallByTts"
	aliyunDefaultAPIVersion       = "2017-05-25"
)

// aliyunTransformer implements HTTPTransformer[*Account] for Aliyun SMS.
type aliyunTransformer struct {
	*BaseTransformer
}

func newAliyunTransformer() *aliyunTransformer {
	transformer := &aliyunTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(core.ProviderTypeSMS),
		string(SubProviderCl253),
		&core.ResponseHandlerConfig{
			SuccessField:      "Code",
			SuccessValue:      "OK",
			ErrorCodeField:    "Code",
			ErrorMessageField: "Message",
			ErrorField:        "Code",
			MessageField:      "Message",
			ResponseType:      core.BodyTypeJSON,
			ValidateResponse:  true,
		},
		WithBeforeHook(func(_ context.Context, msg *Message, account *Account) error {
			transformer.applyAliyunDefaults(msg, account)
			return nil
		}),
		WithSMSHandler(transformer.transformSMS),
		WithVoiceHandler(transformer.transformVoice),
		WithMMSHandler(transformer.transformCardSMS),
	)
	return transformer
}

// https://help.aliyun.com/zh/sdk/product-overview/v3-request-structure-and-signature
type aliyunSignParams struct {
	Host    string
	Method  string
	Path    string
	Query   map[string]string
	Body    []byte
	Action  string
	Version string
	Account *Account
	Headers map[string]string
}

// signAliyunRequest signs the Aliyun request
//   - https://help.aliyun.com/zh/sdk/product-overview/v3-request-structure-and-signature
func (t *aliyunTransformer) signAliyunRequest(params aliyunSignParams) map[string]string {
	const algorithm = "ACS3-HMAC-SHA256"
	//nolint:gosec // Reason: not used for security, only for client nonce generation
	xAcsSignatureNonce := fmt.Sprintf("%x", rand.Int63())
	xAcsDate := time.Now().UTC().Format(time.RFC3339)
	// 计算请求体哈希
	var hashedRequestPayload string
	if len(params.Body) > 0 {
		hashedRequestPayload = utils.SHA256Hex(params.Body)
	} else {
		hashedRequestPayload = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // 空字符串的SHA256
	}
	// 内部生成的 headers
	headers := map[string]string{
		"host":                  params.Host,
		"x-acs-action":          params.Action,
		"x-acs-version":         params.Version,
		"x-acs-date":            xAcsDate,
		"x-acs-signature-nonce": xAcsSignatureNonce,
		"x-acs-content-sha256":  hashedRequestPayload,
	}
	// 合并用户自定义 header，优先用户自定义
	for k, v := range params.Headers {
		headers[strings.ToLower(k)] = v
	}
	if _, ok := headers["content-type"]; !ok {
		headers["content-type"] = "application/json"
	}

	flatQuery := make(map[string]string)
	for k, v := range params.Query {
		t.flattenParams(k, v, flatQuery)
	}

	canonicalQueryString := t.buildCanonicalQueryString(flatQuery)
	canonicalHeaders, signedHeaders := t.buildCanonicalHeaders(headers)
	canonicalRequest := strings.Join([]string{
		params.Method,
		params.Path,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedRequestPayload,
	}, "\n")

	hashedCanonicalRequest := utils.SHA256Hex([]byte(canonicalRequest))
	stringToSign := algorithm + "\n" + hashedCanonicalRequest

	signature := utils.HMACSHA256Hex(params.Account.APISecret, stringToSign)
	authorization := algorithm + " Credential=" + params.Account.APIKey +
		",SignedHeaders=" + signedHeaders + ",Signature=" + signature
	headers["Authorization"] = authorization
	return headers
}

// transformSMS transforms SMS message to HTTP request
//   - SMS API(国内/国外/单发/群发): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
//
// 短信模板即具体发送的短信内容，模板类型支持验证码、通知短信和推广短信。模板由模板变量和模板内容构成，您需要遵守模板内容规范和变量规范。
func (t *aliyunTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	phones := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		phones[i] = t.formatPhoneNumber(mobile, msg.RegionCode)
	}

	params := map[string]string{
		"PhoneNumbers":    strings.Join(phones, ","),
		"SignName":        msg.SignName,
		"TemplateCode":    msg.TemplateID,
		"TemplateParam":   utils.ToJSONString(msg.TemplateParams),
		"SmsUpExtendCode": msg.Extend,
		"OutId":           msg.GetExtraStringOrDefaultEmpty(aliyunOutIDKey),
	}
	urlVals := url.Values{}
	for k, v := range params {
		urlVals.Set(k, v)
	}

	endpoint := t.getEndpointByRegion(msg.Type, msg.GetExtraStringOrDefaultEmpty(aliyunRegionKey))

	reqSpec := &core.HTTPRequestSpec{
		Method: http.MethodPost,
		URL:    "https://" + endpoint + "/",
		Headers: t.signAliyunRequest(aliyunSignParams{
			Host:    endpoint,
			Method:  http.MethodPost,
			Path:    "/",
			Query:   params,
			Action:  aliyunActionSendSms,
			Version: aliyunDefaultAPIVersion,
			Account: account,
		}),
		QueryParams: urlVals,
	}

	return reqSpec, nil, nil
}

// getEndpointByRegion returns the correct endpoint for the given message type and region.
// For voice, always use the default voice endpoint.
// For SMS, use region-specific endpoint if available, otherwise use the default.
func (t *aliyunTransformer) getEndpointByRegion(msgType MessageType, region string) string {
	if msgType == Voice {
		return aliyunDefaultVoiceEndpoint
	}
	if ep, ok := t.regionSmsEndpoints()[region]; ok {
		return ep
	}
	return aliyunDefaultSmsEndpoint
}

// transformCardSMS transforms card SMS message to HTTP request
// TODO: API documentation not understood, temporarily not implemented
//   - 文档地址: https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (t *aliyunTransformer) transformCardSMS(
	_ context.Context,
	_ *Message,
	_ *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	return nil, nil, errors.New("not implemented")
}

// transformVoice transforms voice message to HTTP request
//   - 语音验证码API: https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - 语音通知API: https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
func (t *aliyunTransformer) transformVoice(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 只支持国内
	if msg.IsIntl() {
		return nil, nil, NewUnsupportedInternationalError(string(ProviderTypeAliyun), "voice call")
	}
	// 只支持单号码
	if msg.HasMultipleRecipients() {
		return nil, nil, fmt.Errorf("aliyun voice only supports single number, got %d", len(msg.Mobiles))
	}

	params2 := map[string]string{
		"CalledNumber":     msg.Mobiles[0],
		"CalledShowNumber": msg.GetExtraStringOrDefaultEmpty(aliyunCalledShowNumberKey),
		"PlayTimes":        strconv.Itoa(msg.GetExtraIntOrDefault(aliyunPlayTimesKey, 1)),
		"Volume":           strconv.Itoa(msg.GetExtraIntOrDefault(aliyunVolumeKey, aliyunDefaultVolume)),
		"Speed":            strconv.Itoa(msg.GetExtraIntOrDefault(aliyunSpeedKey, 0)),
		"OutId":            msg.GetExtraStringOrDefaultEmpty(aliyunOutIDKey),
	}
	urlVals2 := url.Values{}
	for k, v := range params2 {
		urlVals2.Set(k, v)
	}

	// 根据消息类型选择API和参数
	isVerification := msg.Category == CategoryVerification
	if isVerification {
		// 验证码类消息：使用TTS API
		params2["TtsCode"] = msg.TemplateID
		if len(msg.TemplateParams) > 0 {
			params2["TtsParam"] = utils.ToJSONString(msg.TemplateParams)
		}
	} else {
		// 通知类消息：使用Voice API
		params2["VoiceCode"] = msg.TemplateID
	}

	// 根据消息类型设置不同的action
	action := aliyunActionSingleCallByVoice
	if isVerification {
		action = aliyunActionSingleCallByTts
	}

	endpoint := t.getEndpointByRegion(msg.Type, msg.GetExtraStringOrDefaultEmpty(aliyunRegionKey))

	reqSpec := &core.HTTPRequestSpec{
		Method: http.MethodPost,
		URL:    "https://" + endpoint + "/",
		Headers: t.signAliyunRequest(aliyunSignParams{
			Host:    endpoint,
			Method:  http.MethodPost,
			Path:    "/",
			Query:   params2,
			Body:    nil,
			Action:  action,
			Version: aliyunDefaultAPIVersion,
			Account: account,
		}),
		QueryParams: urlVals2,
		BodyType:    core.BodyTypeRaw,
	}

	return reqSpec, nil, nil
}

// formatPhoneNumber formats phone number for Aliyun SMS API
// Aliyun requires: domestic numbers without country code, international numbers with country code
// 接收短信的手机号码。手机号码格式：
//   - 国内短信：+/+86/0086/86 或无任何前缀的手机号码，例如 1390000****。
//   - 国际/港澳台消息：国际区号+号码，例如 852000012****。
//   - 接收测试短信的手机号：必须先在控制台绑定测试手机号后才可以发送。
func (t *aliyunTransformer) formatPhoneNumber(mobile string, regionCode int) string {
	if regionCode == 0 || regionCode == 86 {
		return mobile
	}
	// 国际/港澳台消息：国际区号+号码，例如 852000012****。
	return fmt.Sprintf("%d%s", regionCode, mobile)
}

// percentCode encodes strings according to Aliyun's requirements.
func (t *aliyunTransformer) percentCode(str string) string {
	// 阿里云要求: 空格->%20, *->%2A, ~不变, +->%20
	encoded := url.QueryEscape(str)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

// flattenParams 平铺参数（支持数组/嵌套对象，key.1/key.2）.
func (t *aliyunTransformer) flattenParams(prefix string, value interface{}, out map[string]string) {
	switch v := value.(type) {
	case []interface{}:
		for i, item := range v {
			t.flattenParams(prefix+fmt.Sprintf(".%d", i+1), item, out)
		}
	case map[string]interface{}:
		for k, item := range v {
			t.flattenParams(prefix+"."+k, item, out)
		}
	default:
		key := strings.TrimPrefix(prefix, ".")
		out[key] = fmt.Sprintf("%v", v)
	}
}

// buildCanonicalQueryString 构造规范化 query string.
func (t *aliyunTransformer) buildCanonicalQueryString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		ek := t.percentCode(k)
		ev := t.percentCode(params[k])
		parts = append(parts, ek+"="+ev)
	}
	return strings.Join(parts, "&")
}

// buildCanonicalHeaders 构造规范化 header 字符串.
func (t *aliyunTransformer) buildCanonicalHeaders(headers map[string]string) (string, string) {
	keys := make([]string, 0, len(headers))
	for k := range headers {
		lk := strings.ToLower(k)
		if lk == "host" || strings.HasPrefix(lk, "x-acs-") || lk == "content-type" {
			keys = append(keys, lk)
		}
	}
	sort.Strings(keys)
	var canonicalHeaders, signedHeaders []string
	for _, k := range keys {
		canonicalHeaders = append(canonicalHeaders, k+":"+strings.TrimSpace(headers[k]))
		signedHeaders = append(signedHeaders, k)
	}
	return strings.Join(canonicalHeaders, "\n") + "\n", strings.Join(signedHeaders, ";")
}

func (t *aliyunTransformer) regionSmsEndpoints() map[string]string {
	return map[string]string{
		"ap-southeast-1": "dysmsapi.ap-southeast-1.aliyuncs.com",
		"cn-hangzhou":    "dysmsapi.aliyuncs.com",
		"cn-shanghai":    "dysmsapi.aliyuncs.com",
		"cn-shenzhen":    "dysmsapi.aliyuncs.com",
		"cn-beijing":     "dysmsapi.aliyuncs.com",
		"cn-hongkong":    "dysmsapi.aliyuncs.com",
	}
}

// applyAliyunDefaults applies Aliyun-specific defaults to the message.
func (t *aliyunTransformer) applyAliyunDefaults(msg *Message, account *Account) {
	// Apply common defaults first
	msg.ApplyCommonDefaults(account)

	// Apply Aliyun-specific defaults
	region := utils.FirstNonEmpty(
		// 优先使用消息中的 region，其次使用 account 中的 region，最后使用默认值
		msg.GetExtraStringOrDefault(aliyunRegionKey, ""),
		account.Region,
		aliyunDefaultRegion,
	)
	msg.Extras[aliyunRegionKey] = region
}
