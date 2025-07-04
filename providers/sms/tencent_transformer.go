package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Tencent / 腾讯云
// @Website: https://cloud.tencent.com
// @APIDoc: https://cloud.tencent.com/document/product/382/55981
//
// 官方文档:
//   - 短信API: https://cloud.tencent.com/document/product/382/55981
//   - 语音API: https://cloud.tencent.com/document/product/1128/51559
//
// transformer 支持 text（普通短信）和 voice（语音短信）类型。

const (
	tencentAPIDomain         = "tencentcloudapi.com"
	tencentSMSAPIVersion     = "2021-01-11"
	tencentVoiceAPIVersion   = "2020-09-02"
	tencentSmsAction         = "SendSms"
	tencentVoiceAction       = "SendCodeVoice"
	tencentVoiceNotifyAction = "SendVoice"
)

type tencentTransformer struct{}

func init() {
	RegisterTransformer(string(SubProviderTencent), &tencentTransformer{})
}

func (t *tencentTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	return ok && smsMsg.SubProvider == string(SubProviderTencent)
}

func (t *tencentTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, errors.New("invalid message type for tencentTransformer")
	}

	// Apply Tencent-specific defaults
	t.applyTencentDefaults(smsMsg, account)

	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, err
	}

	switch smsMsg.Type {
	case SMSText:
		return t.transformSMS(smsMsg, account)
	case Voice:
		return t.transformVoice(smsMsg, account)
	case MMS:
		return nil, nil, fmt.Errorf("unsupported message type: %v", smsMsg.Type)
	default:
		return nil, nil, fmt.Errorf("unsupported message type: %v", smsMsg.Type)
	}
}

func (t *tencentTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("mobiles is required")
	}
	if msg.TemplateID == "" {
		return errors.New("templateID is required")
	}
	if msg.Type == SMSText && msg.SignName == "" && !msg.IsIntl() {
		return errors.New("domestic sms requires sign name")
	}
	if msg.Type == Voice && len(msg.Mobiles) != 1 {
		return errors.New("voice sms only supports single mobile")
	}
	return nil
}

// transformSMS transforms SMS message to HTTP request
//   - 国内短信: 支持验证码、通知类短信和营销短信
//   - 国际/港澳台短信: 支持验证码、通知类短信和营销短信
//   - API文档: https://cloud.tencent.com/document/product/382/55981
func (t *tencentTransformer) transformSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 格式化手机号
	phoneNumbers := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		phoneNumbers[i] = t.formatTencentPhone(mobile, msg.RegionCode)
	}

	params := map[string]interface{}{
		"PhoneNumberSet":   phoneNumbers,
		"SmsSdkAppId":      msg.GetExtraStringOrDefaultEmpty(tencentSmsSdkAppIDKey),
		"TemplateId":       msg.TemplateID,
		"SignName":         msg.SignName,
		"TemplateParamSet": msg.ParamsOrder,
		"ExtendCode":       msg.Extend,
		"SenderId":         msg.GetExtraStringOrDefaultEmpty(tencentSenderIDKey),
	}

	region := msg.GetExtraStringOrDefault(tencentRegionKey, tencentDefaultRegion)
	requestBody := map[string]interface{}{
		"Action":  tencentSmsAction,
		"Version": tencentSMSAPIVersion,
		"Region":  region,
		"Request": params,
	}

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal tencent request body: %w", err)
	}

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

	endpoint := fmt.Sprintf("%s.%s", "sms", tencentAPIDomain)
	headers := t.buildTencentHeaders(tencentHeaderParams{
		Endpoint:  endpoint,
		Action:    tencentSmsAction,
		Version:   tencentSMSAPIVersion,
		Region:    region,
		AppSecret: account.APISecret,
		BodyData:  bodyData,
		Timestamp: timestamp,
		Date:      date,
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
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 腾讯云语音只支持单发
	calledNumber := t.formatTencentPhone(msg.Mobiles[0], msg.RegionCode)

	params := map[string]interface{}{
		"CalledNumber":  calledNumber,
		"VoiceSdkAppId": msg.GetExtraStringOrDefaultEmpty(tencentVoiceSdkAppIDKey),
	}

	var action string
	if msg.Category == CategoryVerification {
		action = tencentVoiceAction
		// 验证码，仅支持填写数字，实际播报语音时，会自动在数字前补充语音文本"您的验证码是"。示例值：8253
		params["CodeMessage"] = msg.Content
	} else {
		action = tencentVoiceNotifyAction
		params["VoiceId"] = msg.TemplateID
		params["TemplateParamSet"] = msg.ParamsOrder
	}
	if playTimes := msg.GetExtraIntOrDefault(tencentPlayTimesKey, tencentDefaultPlayTimes); playTimes != 0 {
		params["PlayTimes"] = playTimes
	}

	voiceRegion := utils.FirstNonEmpty(
		msg.GetExtraStringOrDefault(tencentRegionKey, ""),
		account.Region,
		tencentDefaultRegion,
	)
	requestBody := map[string]interface{}{
		"Action":  action,
		"Version": tencentVoiceAPIVersion,
		"Region":  voiceRegion,
		"Request": params,
	}

	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal tencent voice request body: %w", err)
	}

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

	endpoint := fmt.Sprintf("%s.%s", "vms", tencentAPIDomain)
	headers := t.buildTencentHeaders(tencentHeaderParams{
		Endpoint:  endpoint,
		Action:    action,
		Version:   tencentVoiceAPIVersion,
		Region:    voiceRegion,
		AppSecret: account.APISecret,
		BodyData:  bodyData,
		Timestamp: timestamp,
		Date:      date,
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
	AppSecret string
	BodyData  []byte
	Timestamp int64
	Date      string
}

// buildTencentHeaders 构造腾讯云API请求头并签名.
func (t *tencentTransformer) buildTencentHeaders(params tencentHeaderParams) map[string]string {
	credentialScope := fmt.Sprintf("%s/sms/tc3_request", params.Date)
	signature := t.calculateTencentSignature(params.AppSecret, params.BodyData, params.Timestamp, credentialScope)

	return map[string]string{
		"Content-Type":    "application/json",
		"Host":            params.Endpoint,
		"X-TC-Action":     params.Action,
		"X-TC-Version":    params.Version,
		"X-TC-Timestamp":  strconv.FormatInt(params.Timestamp, 10),
		"X-TC-Region":     params.Region,
		"Authorization":   signature,
		"X-TC-Credential": fmt.Sprintf("%s/%s", params.AppSecret, credentialScope),
	}
}

// calculateTencentSignature 计算腾讯云API签名.
func (t *tencentTransformer) calculateTencentSignature(
	secretKey string,
	payload []byte,
	timestamp int64,
	credentialScope string,
) string {
	h := hmac.New(sha256.New, []byte("TC3"+secretKey))
	date := time.Unix(timestamp, 0).UTC().Format("20060102")
	h.Write([]byte(date))
	dateKey := h.Sum(nil)

	h = hmac.New(sha256.New, dateKey)
	h.Write([]byte("sms"))
	dateServiceKey := h.Sum(nil)

	h = hmac.New(sha256.New, dateServiceKey)
	h.Write([]byte("tc3_request"))
	signingKey := h.Sum(nil)

	h = hmac.New(sha256.New, signingKey)
	h.Write(payload)
	signature := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=content-type;host, Signature=%s",
		secretKey, credentialScope, signature)
}

// handleTencentResponse 处理腾讯云API响应.
func (t *tencentTransformer) handleTencentResponse(_ int, body []byte) error {
	var result struct {
		Response struct {
			SendStatusSet []struct {
				SerialNo       string `json:"SerialNo"`
				PhoneNumber    string `json:"PhoneNumber"`
				Fee            int    `json:"Fee"`
				SessionContext string `json:"SessionContext"`
				Code           string `json:"Code"`
				Message        string `json:"Message"`
				IsoCode        string `json:"IsoCode"`
			} `json:"SendStatusSet"`
		} `json:"Response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to unmarshal tencent response: %w", err)
	}

	if len(result.Response.SendStatusSet) == 0 {
		return errors.New("no send status found in response")
	}

	for _, status := range result.Response.SendStatusSet {
		if status.Code != "OK" {
			return fmt.Errorf("tencent sms send failed: %s - %s", status.Code, status.Message)
		}
	}

	return nil
}

func (t *tencentTransformer) applyTencentDefaults(msg *Message, account *Account) {
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
}
