//nolint:depguard // intentional use of math/rand for compatibility or legacy reasons
package sms

import (
	"context"
	"encoding/hex"
	"encoding/json"
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
	RegisterTransformer(string(SubProviderAliyun), &aliyunTransformer{})
}

const (
	aliyunDefaultSmsEndpoint   = "dysmsapi.aliyuncs.com"
	aliyunDefaultVoiceEndpoint = "dyvmsapi.aliyuncs.com"
)

// aliyunTransformer implements HTTPTransformer[*core.Account] for Aliyun SMS.
type aliyunTransformer struct{}

// CanTransform checks if this transformer can handle the given message.
func (t *aliyunTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return false
	}
	return smsMsg.SubProvider == string(SubProviderAliyun)
}

// Transform converts an Aliyun SMS message to HTTP request specification.
func (t *aliyunTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for Aliyun: %T", msg)
	}

	// Validate message before processing
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}

	switch smsMsg.Type {
	case SMSText:
		return t.transformTextMessage(smsMsg, account)
	case Voice:
		return t.transformVoiceMessage(smsMsg, account)
	case MMS:
		return t.transformMMSMessage(smsMsg, account)
	default:
		return nil, nil, fmt.Errorf("unsupported message type: %v", smsMsg.Type)
	}
}

// validateMessage validates the message based on its type.
func (t *aliyunTransformer) validateMessage(msg *Message) error {
	switch msg.Type {
	case SMSText:
		return t.validateTextMessage(msg)
	case Voice:
		return t.validateVoiceMessage(msg)
	case MMS:
		return t.validateMMSMessage(msg)
	default:
		return fmt.Errorf("unsupported message type: %s", msg.Type)
	}
}

// validateTextMessage validates text message options.
func (t *aliyunTransformer) validateTextMessage(msg *Message) error {
	// Check for voice-only options
	voiceOnlyOptions := []string{"Volume", "PlayTimes", "CalledShowNumber", "Speed", "OutId"}
	for _, opt := range voiceOnlyOptions {
		if msg.Extras != nil && msg.Extras[opt] != nil {
			return fmt.Errorf("option %s is only applicable to voice messages", opt)
		}
	}

	// Validate template code format if provided
	if msg.TemplateID != "" && !strings.HasPrefix(msg.TemplateID, "SMS_") {
		return errors.New("aliyun template code must start with 'SMS_'")
	}

	// Check required fields for domestic SMS
	if msg.SignName == "" && msg.IsDomestic() {
		return errors.New("aliyun sign name is required for domestic SMS")
	}

	return nil
}

// validateVoiceMessage validates voice message options.
func (t *aliyunTransformer) validateVoiceMessage(msg *Message) error {
	// Check for text-only options
	textOnlyOptions := []string{"SignName"}
	for _, opt := range textOnlyOptions {
		if msg.Extras != nil && msg.Extras[opt] != nil {
			return fmt.Errorf("option %s is not applicable to voice messages", opt)
		}
	}

	// Validate template code format if provided
	if msg.TemplateID != "" && !strings.HasPrefix(msg.TemplateID, "TTS_") {
		return errors.New("Aliyun voice template code must start with 'TTS_'")
	}

	return nil
}

// validateMMSMessage validates MMS message options.
func (t *aliyunTransformer) validateMMSMessage(msg *Message) error {
	// Check for voice-only options
	voiceOnlyOptions := []string{"Volume", "PlayTimes", "CalledShowNumber", "Speed"}
	for _, opt := range voiceOnlyOptions {
		if msg.Extras != nil && msg.Extras[opt] != nil {
			return fmt.Errorf("option %s is only applicable to voice messages", opt)
		}
	}

	// Check for text-only options
	textOnlyOptions := []string{"SignName"}
	for _, opt := range textOnlyOptions {
		if msg.Extras != nil && msg.Extras[opt] != nil {
			return fmt.Errorf("option %s is not applicable to MMS messages", opt)
		}
	}

	// Validate template code format if provided
	if msg.TemplateID != "" && !strings.HasPrefix(msg.TemplateID, "CARD_") {
		return errors.New("Aliyun MMS template code must start with 'CARD_'")
	}

	// Check required fields for MMS
	if msg.SignName == "" {
		return errors.New("sign name is required for Aliyun MMS messages")
	}

	// Validate fallback type if provided
	if fallbackType := msg.GetExtraStringOrDefault(aliyunFallbackTypeKey, ""); fallbackType != "" {
		//   - SMS：不支持卡片短信的号码，回落文本短信。
		//   - DIGITALSMS：不支持卡片短信的号码，回落数字短信。
		//   - NONE：不需要回落。
		validFallbackTypes := []string{"SMS", "DIGITALSMS", "NONE"}
		isValid := false
		for _, validType := range validFallbackTypes {
			if fallbackType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid fallback type: %s, must be one of %v", fallbackType, validFallbackTypes)
		}
	}

	return nil
}

// transformTextMessage transforms text SMS message to HTTP request.
func (t *aliyunTransformer) transformTextMessage(
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	return t.transformSMS(msg, account)
}

// transformVoiceMessage transforms voice message to HTTP request.
func (t *aliyunTransformer) transformVoiceMessage(
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	return t.transformVoice(msg, account)
}

// transformMMSMessage transforms MMS message to HTTP request.
func (t *aliyunTransformer) transformMMSMessage(
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	return t.transformCardSMS(msg, account)
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
	Account *core.Account
	Headers map[string]string
}

// signAliyunRequest signs the Aliyun request
// https://help.aliyun.com/zh/sdk/product-overview/v3-request-structure-and-signature
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
	// 过滤 canonical headers
	filtered := make(map[string]string)
	for k, v := range headers {
		lowerKey := strings.ToLower(k)
		if lowerKey == "host" || lowerKey == "content-type" || strings.HasPrefix(lowerKey, "x-acs-") {
			filtered[lowerKey] = strings.TrimSpace(v)
		}
	}
	// 升序排列
	canonKeys := make([]string, 0, len(filtered))
	for k := range filtered {
		canonKeys = append(canonKeys, k)
	}
	sort.Strings(canonKeys)
	// 拼接 canonicalHeaders 和 signedHeaders
	canonicalHeaders := ""
	signedHeaders := ""
	for _, k := range canonKeys {
		canonicalHeaders += k + ":" + filtered[k] + "\n"
		signedHeaders += k + ";"
	}
	signedHeaders = strings.TrimSuffix(signedHeaders, ";")
	canonicalQueryString := ""
	keys := make([]string, 0, len(params.Query))
	for k := range params.Query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := params.Query[k]
		canonicalQueryString += t.percentCode(url.QueryEscape(k)) + "=" + t.percentCode(url.QueryEscape(v)) + "&"
	}
	canonicalQueryString = strings.TrimSuffix(canonicalQueryString, "&")
	canonicalRequest := params.Method + "\n" + params.Path + "\n" + canonicalQueryString + "\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + hashedRequestPayload
	hashedCanonicalRequest := utils.SHA256Hex([]byte(canonicalRequest))
	stringToSign := algorithm + "\n" + hashedCanonicalRequest
	signature := hex.EncodeToString(utils.HMACSHA256([]byte(params.Account.Secret), []byte(stringToSign)))
	authorization := algorithm + " Credential=" + params.Account.Key + ",SignedHeaders=" + signedHeaders + ",Signature=" + signature
	headers["Authorization"] = authorization
	return headers
}

// transformSMS transforms SMS message to HTTP request
// SMS API(国内/国外/单发/群发): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
// 短信模板即具体发送的短信内容，模板类型支持验证码、通知短信和推广短信。模板由模板变量和模板内容构成，您需要遵守模板内容规范和变量规范。
func (t *aliyunTransformer) transformSMS(
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.getEndpoint(msg.IsIntl(), account)
	phones := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		phones[i] = t.formatPhoneNumber(mobile, msg.RegionCode)
	}

	params := map[string]string{
		"PhoneNumbers":  strings.Join(phones, ","),
		"SignName":      msg.SignName,
		"TemplateCode":  msg.TemplateID,
		"TemplateParam": utils.ToJSONString(msg.TemplateParams),
	}

	// 使用新的Message结构体字段
	if msg.Extend != "" {
		params["SmsUpExtendCode"] = msg.Extend
	}
	if msg.UID != "" {
		params["OutId"] = msg.UID
	}

	if v := msg.GetExtraStringOrDefault(aliyunSmsUpExtendCodeKey, ""); v != "" {
		params["SmsUpExtendCode"] = v
	}
	if v := msg.GetExtraStringOrDefault(aliyunOutIDKey, ""); v != "" {
		params["OutId"] = v
	}

	reqSpec := &core.HTTPRequestSpec{
		Method: "POST",
		URL:    "https://" + endpoint + "/",
		Headers: t.signAliyunRequest(aliyunSignParams{
			Host:    endpoint,
			Method:  "POST",
			Path:    "/",
			Query:   params,
			Body:    nil,
			Action:  "SendSms",
			Version: "2017-05-25",
			Account: account,
		}),
		QueryParams: params,
		BodyType:    "raw",
	}

	return reqSpec, t.handleAliyunResponse, nil
}

// transformCardSMS transforms card SMS message to HTTP request
//
//   - 文档地址: https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (t *aliyunTransformer) transformCardSMS(
	_ *Message,
	_ *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	return nil, nil, errors.New("not implemented")
}

// transformVoice transforms voice message to HTTP request
//   - 语音验证码API: https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - 语音通知API: https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
func (t *aliyunTransformer) transformVoice(
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 只支持国内
	if msg.IsIntl() {
		return nil, nil, NewUnsupportedInternationalError(string(ProviderTypeAliyun), "voice call")
	}
	// 只支持单号码
	if msg.HasMultipleRecipients() {
		return nil, nil, fmt.Errorf("aliyun voice only supports single number, got %d", len(msg.Mobiles))
	}

	params := map[string]string{
		"CalledNumber":     msg.Mobiles[0],
		"CalledShowNumber": msg.GetExtraStringOrDefault(aliyunCalledShowNumberKey, ""),
		"PlayTimes":        strconv.Itoa(msg.GetExtraIntOrDefault(aliyunPlayTimesKey, 1)),
		"Volume":           strconv.Itoa(msg.GetExtraIntOrDefault(aliyunVolumeKey, aliyunDefaultVolume)),
		"Speed":            strconv.Itoa(msg.GetExtraIntOrDefault(aliyunSpeedKey, 0)),
		"OutId":            msg.GetExtraStringOrDefault(aliyunOutIDKey, ""),
	}

	// 根据消息类型选择API和参数
	isVerification := msg.Category == CategoryVerification
	if isVerification {
		// 验证码类消息：使用TTS API
		params["TtsCode"] = msg.TemplateID
		if len(msg.TemplateParams) > 0 {
			params["TtsParam"] = utils.ToJSONString(msg.TemplateParams)
		}
	} else {
		// 通知类消息：使用Voice API
		params["VoiceCode"] = msg.TemplateID
	}

	// 根据消息类型设置不同的action
	action := "SingleCallByVoice"
	if isVerification {
		action = "SingleCallByTts"
	}

	reqSpec := &core.HTTPRequestSpec{
		Method: "POST",
		URL:    "https://" + aliyunDefaultVoiceEndpoint + "/",
		Headers: t.signAliyunRequest(aliyunSignParams{
			Host:    aliyunDefaultVoiceEndpoint,
			Method:  "POST",
			Path:    "/",
			Query:   params,
			Body:    nil,
			Action:  action,
			Version: "2017-05-25",
			Account: account,
		}),
		QueryParams: params,
		BodyType:    "raw",
	}

	return reqSpec, t.handleAliyunResponse, nil
}

// getEndpoint returns the appropriate endpoint based on international flag and account.
func (t *aliyunTransformer) getEndpoint(isIntl bool, account *core.Account) string {
	// Use default endpoints
	if isIntl && account.IntlEndpoint != "" {
		return account.IntlEndpoint
	}
	if account.Endpoint != "" {
		return account.Endpoint
	}
	return aliyunDefaultSmsEndpoint
}

// handleAliyunResponse handles Aliyun API response.
func (t *aliyunTransformer) handleAliyunResponse(statusCode int, body []byte) error {
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for Aliyun specific response format
	if code, ok := response["Code"].(string); ok {
		if code == "OK" {
			return nil
		}
		if msg, okMsg := response["Message"].(string); okMsg {
			return errors.New(msg)
		}
		return &Error{
			Code:     code,
			Message:  "unknown error",
			Provider: string(SubProviderAliyun),
		}
	}

	return nil
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
