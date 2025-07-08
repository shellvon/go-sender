package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// tencentTransformer implements HTTPRequestTransformer for Tencent SMS.
// It supports sending text message and voice message.
//
// Reference:
//   - Official Website: https://cloud.tencent.com
//   - API Docs: https://cloud.tencent.com/document/product/382/55981
//   - SMS API: https://cloud.tencent.com/document/product/382/55981
//   - Voice API: https://cloud.tencent.com/document/product/1128/51559
//   - Signature Algorithm: https://github.com/TencentCloud/signature-process-demo/blob/main/services/sms/signature-v3/golang/demo.go

const (
	tencentAPIDomain         = "tencentcloudapi.com"
	tencentSMSAPIVersion     = "2021-01-11"
	tencentVoiceAPIVersion   = "2020-09-02"
	tencentSmsAction         = "SendSms"
	tencentVoiceAction       = "SendCodeVoice"
	tencentVoiceNotifyAction = "SendVoice"
)

type tencentTransformer struct {
	*BaseTransformer
}

func newTencentTransformer() *tencentTransformer {
	transformer := &tencentTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(SubProviderTencent),
		nil,
		HTTPOptions{
			AddBeforeHook(transformer.beforeSend),
		},
		WithSMSHandler(transformer.transformSMS),
		WithVoiceHandler(transformer.transformVoice),
	)
	return transformer
}

func init() {
	RegisterTransformer(string(SubProviderTencent), newTencentTransformer())
}

func (t *tencentTransformer) beforeSend(ctx context.Context, msg *Message, account *Account) error {
	if err := t.applyTencentDefaults(ctx, msg, account); err != nil {
		return err
	}
	if err := t.validateMessage(msg); err != nil {
		return err
	}
	return nil
}

func (t *tencentTransformer) applyTencentDefaults(_ context.Context, msg *Message, account *Account) error {
	// Apply common defaults first
	msg.ApplyCommonDefaults(account)

	// Apply Aliyun-specific defaults
	region := utils.FirstNonEmpty(
		// 优先使用消息中的 region，其次使用 account 中的 region，最后使用默认值
		msg.GetExtraStringOrDefaultEmpty(tencentRegionKey),
		account.Region,
		aliyunDefaultRegion,
	)
	msg.Extras[aliyunRegionKey] = region

	if msg.Extras[tencentSmsSdkAppIDKey] == "" {
		msg.Extras[tencentSmsSdkAppIDKey] = account.AppID
	}
	if msg.Extras[tencentVoiceSdkAppIDKey] == "" {
		msg.Extras[tencentVoiceSdkAppIDKey] = account.AppID
	}
	return nil
}

func (t *tencentTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return NewProviderError(string(SubProviderTencent), "MISSING_PARAM", "mobiles is required")
	}
	if msg.TemplateID == "" {
		return NewProviderError(string(SubProviderTencent), "MISSING_PARAM", "templateID is required")
	}
	if msg.Type == SMSText && msg.SignName == "" && !msg.IsIntl() {
		return NewProviderError(string(SubProviderTencent), "MISSING_SIGNATURE", "domestic sms requires sign name")
	}
	if msg.Type == Voice && len(msg.Mobiles) != 1 {
		return NewProviderError(
			string(SubProviderTencent),
			"INVALID_MOBILE_NUMBER",
			"voice sms only supports single mobile",
		)
	}
	return nil
}

// transformSMS transforms SMS message to HTTP request
//   - 国内短信: 支持验证码、通知类短信和营销短信
//   - 国际/港澳台短信: 支持验证码、通知类短信和营销短信
//   - API文档: https://cloud.tencent.com/document/product/382/55981
func (t *tencentTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 格式化手机号
	phoneNumbers := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		phoneNumbers[i] = t.formatTencentPhone(mobile, msg.RegionCode)
	}

	params := map[string]interface{}{
		"PhoneNumberSet": phoneNumbers,
		"SmsSdkAppId":    msg.GetExtraStringOrDefaultEmpty(tencentSmsSdkAppIDKey),
		"TemplateId":     msg.TemplateID,
		"SignName":       msg.SignName,
		"ExtendCode":     msg.Extend,
		"SenderId":       msg.GetExtraStringOrDefaultEmpty(tencentSenderIDKey),
	}

	if len(msg.ParamsOrder) > 0 {
		params["TemplateParamSet"] = msg.ParamsOrder
	}

	region := msg.GetExtraStringOrDefault(tencentRegionKey, tencentDefaultRegion)

	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, NewProviderError(
			string(SubProviderTencent),
			"JSON_MARSHAL_ERROR",
			fmt.Sprintf("failed to marshal tencent request body: %v", err),
		)
	}

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

	endpoint := fmt.Sprintf("%s.%s", "sms", tencentAPIDomain)
	headers := t.buildTencentHeaders(tencentHeaderParams{
		Endpoint:  endpoint,
		Action:    tencentSmsAction,
		Version:   tencentSMSAPIVersion,
		Region:    region,
		SecretID:  account.APIKey,
		SecretKey: account.APISecret,
		BodyData:  bodyData,
		Timestamp: timestamp,
		Date:      date,
		Service:   "sms",
	})

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      fmt.Sprintf("https://%s", endpoint),
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, t.handleTencentResponse, nil
}

// transformVoice transforms voice message to HTTP request
//   - 语音验证码API: https://cloud.tencent.com/document/product/1128/51559
//   - 语音通知API: https://cloud.tencent.com/document/product/1128/51558
//
// 当短信为验证码类型时，使用语音验证码API，否则使用语音通知API.
func (t *tencentTransformer) transformVoice(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 腾讯云语音只支持单发
	calledNumber := t.formatTencentPhone(msg.Mobiles[0], msg.RegionCode)

	params := map[string]interface{}{
		"CalledNumber": calledNumber,
		// 注意是Appid，不是AppId
		"VoiceSdkAppid": msg.GetExtraStringOrDefaultEmpty(tencentVoiceSdkAppIDKey),
	}

	var action string
	if msg.Category == CategoryVerification {
		action = tencentVoiceAction
		// 验证码，仅支持填写数字，实际播报语音时，会自动在数字前补充语音文本"您的验证码是"。示例值：8253
		params["CodeMessage"] = msg.Content
	} else {
		action = tencentVoiceNotifyAction
		params["VoiceId"] = msg.TemplateID
		if len(msg.ParamsOrder) > 0 {
			params["TemplateParamSet"] = msg.ParamsOrder
		}
	}
	if playTimes := msg.GetExtraIntOrDefault(tencentPlayTimesKey, tencentDefaultPlayTimes); playTimes != 0 {
		params["PlayTimes"] = playTimes
	}

	voiceRegion := utils.FirstNonEmpty(
		msg.GetExtraStringOrDefault(tencentRegionKey, ""),
		account.Region,
		tencentDefaultRegion,
	)

	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, NewProviderError(
			string(SubProviderTencent),
			"JSON_MARSHAL_ERROR",
			fmt.Sprintf("failed to marshal tencent voice request body: %v", err),
		)
	}

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

	endpoint := fmt.Sprintf("%s.%s", "vms", tencentAPIDomain)
	headers := t.buildTencentHeaders(tencentHeaderParams{
		Endpoint:  endpoint,
		Action:    action,
		Version:   tencentVoiceAPIVersion,
		Region:    voiceRegion,
		SecretID:  account.APIKey,
		SecretKey: account.APISecret,
		BodyData:  bodyData,
		Timestamp: timestamp,
		Date:      date,
		Service:   "vms",
	})

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      fmt.Sprintf("https://%s", endpoint),
		Headers:  headers,
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, t.handleTencentResponse, nil
}

// formatTencentPhone 格式化手机号，始终+开头，国内强制+86，国际为+regionCode.
func (t *tencentTransformer) formatTencentPhone(mobile string, regionCode int) string {
	if regionCode == 0 {
		regionCode = 86
	}
	return fmt.Sprintf("+%d%s", regionCode, mobile)
}

// tencentHeaderParams 腾讯云API请求头参数.
type tencentHeaderParams struct {
	Endpoint  string
	Action    string
	Version   string
	Region    string
	SecretID  string
	SecretKey string
	Service   string
	BodyData  []byte
	Timestamp int64
	Date      string
}

// buildTencentHeaders 构造腾讯云API请求头并签名.
// SMS API: https://cloud.tencent.com/document/product/382/55981
// 参考: https://github.com/TencentCloud/signature-process-demo/blob/main/services/sms/signature-v3/golang/demo.go
func (t *tencentTransformer) buildTencentHeaders(p tencentHeaderParams) map[string]string {
	// ----- Step1: Canonical Request -----
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("content-type:application/json; charset=utf-8\nhost:%s\n", p.Endpoint)
	signedHeaders := "content-type;host"
	hashedPayload := t.sha256Hex(p.BodyData)
	canonicalRequest := strings.Join([]string{
		httpRequestMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedPayload,
	}, "\n")

	// ----- Step2: String to sign -----
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", p.Date, p.Service)
	stringToSign := fmt.Sprintf("TC3-HMAC-SHA256\n%d\n%s\n%s",
		p.Timestamp,
		credentialScope,
		t.sha256Hex([]byte(canonicalRequest)),
	)

	// ----- Step3: Signature -----
	secretDate := t.hmacSha256([]byte("TC3"+p.SecretKey), []byte(p.Date))
	secretService := t.hmacSha256(secretDate, []byte(p.Service))
	secretSigning := t.hmacSha256(secretService, []byte("tc3_request"))
	signature := hex.EncodeToString(t.hmacSha256(secretSigning, []byte(stringToSign)))

	authorization := fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		p.SecretID,
		credentialScope,
		signedHeaders,
		signature,
	)

	return map[string]string{
		"Content-Type":   "application/json; charset=utf-8",
		"Host":           p.Endpoint,
		"X-TC-Action":    p.Action,
		"X-TC-Version":   p.Version,
		"X-TC-Timestamp": strconv.FormatInt(p.Timestamp, 10),
		"X-TC-Region":    p.Region,
		"Authorization":  authorization,
	}
}

func (t *tencentTransformer) sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func (t *tencentTransformer) hmacSha256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// handleTencentResponse 处理腾讯云API响应.
func (t *tencentTransformer) handleTencentResponse(resp *http.Response) error {
	body, _, err := utils.ReadAndClose(resp)
	subProvider := string(SubProviderSmsbao)
	if err != nil {
		return NewProviderError(subProvider, "READ_ERROR", err.Error())
	}

	// Tencent 返回有两种结构：
	// 1. 成功/失败明细在 SendStatusSet 数组里
	// 2. 整体失败时，只有 Error 字段

	var response struct {
		Response struct {
			Error *struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error,omitempty"`

			SendStatusSet []struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"SendStatusSet,omitempty"`
		} `json:"Response"`
	}

	if decodeErr := json.Unmarshal(body, &response); decodeErr != nil {
		return NewProviderError(subProvider, "PARSE_ERROR", decodeErr.Error())
	}

	if response.Response.Error != nil {
		return &Error{
			Code:     response.Response.Error.Code,
			Message:  response.Response.Error.Message,
			Provider: subProvider,
		}
	}

	if len(response.Response.SendStatusSet) == 0 {
		return NewProviderError(
			subProvider,
			"NO_STATUS_SET",
			"tencent API returned success but no SendStatusSet",
		)
	}

	for _, status := range response.Response.SendStatusSet {
		if status.Code != "OK" {
			return &Error{
				Code:     status.Code,
				Message:  status.Message,
				Provider: subProvider,
			}
		}
	}
	return nil
}
