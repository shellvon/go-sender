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
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// sendAliyunSMS 发送阿里云短信（支持国内和国际）
// 文档: https://next.api.aliyun.com/api/Dysmsapi/2017-05-25/SendSms
// 国际短信文档: https://help.aliyun.com/document_detail/108084.html
func sendAliyunSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	isIntl := false
	for _, m := range msg.Mobiles {
		if utils.IsInternationalMobile(m) {
			isIntl = true
			break
		}
	}

	if isIntl && len(msg.Mobiles) > 1 {
		return fmt.Errorf("aliyun international SMS only supports single send, got %d mobiles", len(msg.Mobiles))
	}

	params := map[string]string{
		"AccessKeyId":      provider.AppID,
		"Format":           "JSON",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignName":         msg.SignName,
		"TemplateCode":     msg.TemplateCode,
		"TemplateParam":    toJSONString(msg.TemplateParams),
	}

	if isIntl {
		params["RegionId"] = "ap-southeast-1"
		params["Action"] = "SendSmsIntl"
		params["Version"] = "2018-05-01"
		params["To"] = msg.Mobiles[0]
	} else {
		params["RegionId"] = "cn-hangzhou"
		params["Action"] = "SendSms"
		params["Version"] = "2017-05-25"
		params["PhoneNumbers"] = strings.Join(msg.Mobiles, ",")
	}

	// 透传 OutId
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if outId, ok := metadata["OutId"].(string); ok && outId != "" {
			params["OutId"] = outId
		}
	}

	// 生成签名
	signature := aliyunSignForPost(params, provider.AppSecret)
	params["Signature"] = signature

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	resp, _, err := utils.DoRequest(ctx, "https://dysmsapi.aliyuncs.com/", utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Raw: []byte(form.Encode()),
	})
	if err != nil {
		return fmt.Errorf("aliyun SMS request failed: %w", err)
	}

	var result struct {
		Message   string `json:"Message"`
		RequestId string `json:"RequestId"`
		BizId     string `json:"BizId"`
		Code      string `json:"Code"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse aliyun response: %w", err)
	}

	if result.Code != "OK" {
		return &SMSError{
			Code:     result.Code,
			Message:  result.Message,
			Provider: string(ProviderTypeAliyun),
		}
	}

	return nil
}

func toJSONString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// aliyunSignForPost 生成POST方式的阿里云签名
func aliyunSignForPost(params map[string]string, accessKeySecret string) string {
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
	stringToSign := "POST&%2F&" + percentEncode(canonicalizedQueryString)
	mac := hmac.New(sha1.New, []byte(accessKeySecret+"&"))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// percentEncode 按阿里云要求编码
func percentEncode(str string) string {
	// 阿里云要求: 空格->%20, *->%2A, ~不变, +->%20
	encoded := url.QueryEscape(str)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}
