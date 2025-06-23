package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// sendAliyunIntlSMS 发送阿里云国际短信（SendMessageWithTemplate）
func sendAliyunIntlSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	if len(msg.Mobiles) == 1 {
		return sendAliyunIntlSingleSMS(ctx, provider, msg)
	}
	return sendAliyunIntlBatchSMS(ctx, provider, msg)
}

// 单发实现
func sendAliyunIntlSingleSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	params := map[string]string{
		"RegionId":         "ap-southeast-1",
		"AccessKeyId":      provider.AppID,
		"Format":           "JSON",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Version":          "2018-05-01",
		"Action":           "SendMessageWithTemplate",
		"To":               msg.Mobiles[0],
		"From":             msg.SignName,
		"TemplateCode":     msg.TemplateCode,
		"TemplateParam":    toJSONString(msg.TemplateParams),
	}
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if outId, ok := metadata["OutId"].(string); ok && outId != "" {
			params["OutId"] = outId
		}
	}
	// 签名
	signature := aliyunIntlSign(params, provider.AppSecret)
	params["Signature"] = signature
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	resp, _, err := utils.DoRequest(ctx, "https://dysmsapi.ap-southeast-1.aliyuncs.com/", utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Raw: []byte(form.Encode()),
	})
	if err != nil {
		return fmt.Errorf("aliyun intl SMS request failed: %w", err)
	}
	var result struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseDescription string `json:"ResponseDescription"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse aliyun intl response: %w", err)
	}
	if result.ResponseCode != "OK" {
		return &SMSError{
			Code:     result.ResponseCode,
			Message:  result.ResponseDescription,
			Provider: string(ProviderTypeAliyunGlobe),
		}
	}
	return nil
}

// 批量实现
func sendAliyunIntlBatchSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	params := map[string]string{
		"RegionId":         "ap-southeast-1",
		"AccessKeyId":      provider.AppID,
		"Format":           "JSON",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Version":          "2018-05-01",
		"Action":           "BatchSendMessageToGlobe",
	}
	// To: 号码数组，Message: 内容数组（如无一对一内容则全部用msg.Content）
	toArr := msg.Mobiles
	msgArr := make([]string, len(toArr))
	for i := range msgArr {
		msgArr[i] = msg.Content
	}
	params["To"] = toJSONStringArray(toArr)
	params["Message"] = toJSONStringArray(msgArr)
	// 透传 Type, TaskId, ValidityPeriod, ChannelId
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if v, ok := metadata["Type"].(string); ok && v != "" {
			params["Type"] = v
		}
		if v, ok := metadata["TaskId"].(string); ok && v != "" {
			params["TaskId"] = v
		}
		if v, ok := metadata["ValidityPeriod"].(int); ok && v > 0 {
			params["ValidityPeriod"] = strconv.Itoa(v)
		}
		if v, ok := metadata["ChannelId"].(string); ok && v != "" {
			params["ChannelId"] = v
		}
	}
	// 签名
	signature := aliyunIntlSign(params, provider.AppSecret)
	params["Signature"] = signature
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	resp, _, err := utils.DoRequest(ctx, "https://dysmsapi.ap-southeast-1.aliyuncs.com/", utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Raw: []byte(form.Encode()),
	})
	if err != nil {
		return fmt.Errorf("aliyun intl batch SMS request failed: %w", err)
	}
	var result struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseDescription string `json:"ResponseDescription"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse aliyun intl batch response: %w", err)
	}
	if result.ResponseCode != "OK" {
		return &SMSError{
			Code:     result.ResponseCode,
			Message:  result.ResponseDescription,
			Provider: string(ProviderTypeAliyunGlobe),
		}
	}
	return nil
}

// toJSONStringArray 辅助函数，将字符串数组转为JSON数组字符串
func toJSONStringArray(arr []string) string {
	b, _ := json.Marshal(arr)
	return string(b)
}

// aliyunIntlSign 生成阿里云国际短信签名（GET方式）
func aliyunIntlSign(params map[string]string, accessKeySecret string) string {
	// 1. 参数排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var canonicalizedQueryString string
	for i, k := range keys {
		if i > 0 {
			canonicalizedQueryString += "&"
		}
		canonicalizedQueryString += percentEncode(k) + "=" + percentEncode(params[k])
	}
	stringToSign := "GET&%2F&" + percentEncode(canonicalizedQueryString)
	mac := hmac.New(sha1.New, []byte(accessKeySecret+"&"))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
