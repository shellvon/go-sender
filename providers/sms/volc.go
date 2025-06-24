package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// VolcResponse represents the response from Volcengine SMS API
// API Documentation: https://www.volcengine.com/docs/6361/67380
type VolcResponse struct {
	ResponseMetadata struct {
		RequestId string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Error     *struct {
			CodeN   int    `json:"CodeN"`
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error,omitempty"`
	} `json:"ResponseMetadata"`
	Result struct {
		MessageID string `json:"MessageID"`
	} `json:"Result"`
}

const (
	volcEndpoint = "https://sms.volcengineapi.com/?Action=SendSms&Version=2020-01-01"
)

// sendVolcSMS 发送火山引擎短信（仅支持中国大陆号码）
// API文档: https://www.volcengine.com/docs/6361/67380
func sendVolcSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	// 检查是否有国际号码
	for _, m := range msg.Mobiles {
		if utils.IsInternationalMobile(m) {
			return fmt.Errorf("volcengine SMS does not support international numbers: %s", m)
		}
	}

	body := map[string]interface{}{
		"SmsAccount":   provider.AppID,
		"Sign":         provider.SignName,
		"TemplateID":   msg.TemplateCode,
		"PhoneNumbers": strings.Join(msg.Mobiles, ","),
	}
	if len(msg.TemplateParams) > 0 {
		b, _ := json.Marshal(msg.TemplateParams)
		body["TemplateParam"] = string(b)
	}
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if tag, ok := metadata["tag"].(string); ok && tag != "" {
			body["Tag"] = tag
		}
	}

	bodyJSON, _ := json.Marshal(body)

	headers := buildVolcHeaders(provider, bodyJSON)

	resp, _, err := utils.DoRequest(ctx, volcEndpoint, utils.RequestOptions{
		Method:  http.MethodPost,
		Headers: headers,
		Raw:     bodyJSON,
	})
	if err != nil {
		return fmt.Errorf("volcengine SMS request failed: %w", err)
	}

	var result VolcResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse volcengine response: %w", err)
	}

	if result.ResponseMetadata.Error != nil {
		return &SMSError{
			Code:     result.ResponseMetadata.Error.Code,
			Message:  result.ResponseMetadata.Error.Message,
			Provider: "volcengine",
		}
	}

	return nil
}

// buildVolcHeaders 构建火山引擎TOP网关签名头
// 签名文档: https://www.volcengine.com/docs/6361/1205061
func buildVolcHeaders(provider *SMSProvider, body []byte) map[string]string {
	ak := provider.AppID
	sk := provider.AppSecret
	// service := "volcSMS"
	// region := "cn-north-1"
	// version := "2020-01-01"
	// action := "SendSms"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	xDate := timestamp[:10]

	headers := map[string]string{
		"Content-Type": "application/json;charset=utf-8",
		"Host":         "sms.volcengineapi.com",
		"X-Date":       timestamp,
	}

	// 构建签名字符串
	canonicalHeaders := "content-type:application/json;charset=utf-8\nhost:sms.volcengineapi.com\nx-date:" + timestamp + "\n"
	canonicalRequest := "POST\n/\n\n" + canonicalHeaders + "\n" + string(body)
	stringToSign := "HMAC-SHA256\n" + xDate + "\n" + hashSHA256([]byte(canonicalRequest))
	signingKey := hmacSHA256([]byte(sk), []byte(xDate))
	signature := hmacSHA256(signingKey, []byte(stringToSign))
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	authHeader := "HMAC-SHA256 Credential=" + ak + ", SignedHeaders=content-type;host;x-date, Signature=" + signatureBase64
	headers["Authorization"] = authHeader

	return headers
}

func hashSHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
