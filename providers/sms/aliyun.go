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

// sendAliyunSMS 发送阿里云中国内地短信（支持单个/批量）
// 文档: https://next.api.aliyun.com/api/Dysmsapi/2017-05-25/SendSms
func sendAliyunSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	params := map[string]string{
		"RegionId":         "cn-hangzhou",
		"AccessKeyId":      provider.AppID,
		"Format":           "JSON",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"PhoneNumbers":     strings.Join(msg.Mobiles, ","),
		"SignName":         msg.SignName,
		"TemplateCode":     msg.TemplateCode,
		"TemplateParam":    toJSONString(msg.TemplateParams),
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

	// 构造POST表单
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
