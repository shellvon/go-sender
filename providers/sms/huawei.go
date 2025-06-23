package sms

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shellvon/go-sender/utils"
)

// 华为云短信服务实现
// 仅支持模板短信发送，API文档：https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
const (
	huaweiEndpointHost = "https://api.rtc.huaweicloud.com:10443"
	huaweiEndpointURI  = "/sms/batchSendSms/v1"
	huaweiSuccessCode  = "000000"
)

// sendHuaweiSMS 华为云短信发送实现
func sendHuaweiSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	appKey := provider.AppID
	appSecret := provider.AppSecret
	channel := provider.Channel // 通道号
	statusCallback := provider.Callback

	templateId := msg.TemplateCode

	paras := "[]"
	if len(msg.TemplateParamsArray) > 0 {
		b, _ := json.Marshal(msg.TemplateParamsArray)
		paras = string(b)
	}
	params := map[string]string{
		"from":          channel,
		"to":            strings.Join(msg.Mobiles, ","),
		"templateId":    templateId,
		"templateParas": paras,
		// 签名名称，必须是已审核通过的，与模板类型一致的签名名称。
		"signature": provider.SignName,
	}
	if statusCallback != "" {
		params["statusCallback"] = statusCallback
	}

	headers := getHuaweiHeaders(appKey, appSecret)
	endpoint := getHuaweiEndpoint(provider)

	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method:  "POST",
		Headers: headers,
		Data:    params,
	})
	if err != nil {
		return fmt.Errorf("huawei SMS request failed: %w", err)
	}
	return parseHuaweiResponse(resp)
}

func getHuaweiEndpoint(provider *SMSProvider) string {
	endpoint := huaweiEndpointHost
	if provider.Endpoint != "" {
		endpoint = provider.Endpoint
	}
	return strings.TrimRight(endpoint, "/") + huaweiEndpointURI
}

func getHuaweiHeaders(appKey, appSecret string) map[string]string {
	return map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "WSSE realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\"",
		"X-WSSE":        buildHuaweiWsseHeader(appKey, appSecret),
	}
}

func buildHuaweiWsseHeader(appKey, appSecret string) string {
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	h := sha256.New()
	h.Write([]byte(nonce + now + appSecret))
	passwordDigest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf(
		"UsernameToken Username=\"%s\",PasswordDigest=\"%s\",Nonce=\"%s\",Created=\"%s\"",
		appKey, passwordDigest, nonce, now,
	)
}

func parseHuaweiResponse(resp []byte) error {
	var result struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse huawei response: %w", err)
	}
	if result.Code != huaweiSuccessCode {
		return &SMSError{
			Code:     result.Code,
			Message:  result.Description,
			Provider: string(ProviderTypeHuawei),
		}
	}
	return nil
}
