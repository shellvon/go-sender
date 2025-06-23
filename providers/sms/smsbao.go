package sms

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

const smsbaoAPIURL = "http://api.smsbao.com/sms"

// https://www.smsbao.com/openapi/213.html
// sendSmsbaoSMS 短信宝短信发送实现
func sendSmsbaoSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return fmt.Errorf("smsbao: mobiles cannot be empty")
	}
	if len(msg.Mobiles) > 99 {
		return fmt.Errorf("smsbao: at most 99 mobiles per request")
	}

	content := msg.Content
	if msg.TemplateCode != "" {
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	// 短信宝要求内容为UTF-8 URL编码
	content = url.QueryEscape(content)

	params := map[string]interface{}{
		"u": provider.AppID,
		"p": provider.AppSecret, // 需为MD5后的密码或ApiKey
		"m": strings.Join(msg.Mobiles, ","),
		"c": content,
	}
	// 支持专用通道产品ID
	if provider.Channel != "" {
		params["g"] = provider.Channel
	}

	resp, _, err := utils.DoRequest(ctx, smsbaoAPIURL, utils.RequestOptions{
		Method: "GET",
		Query:  params,
	})
	if err != nil {
		return fmt.Errorf("smsbao SMS request failed: %w", err)
	}
	return parseSmsbaoResponse(resp)
}

// parseSmsbaoResponse 解析短信宝响应
func parseSmsbaoResponse(resp []byte) error {
	result := strings.TrimSpace(string(resp))
	if result == "0" {
		return nil
	}
	return &SMSError{
		Code:     result,
		Message:  smsbaoErrorMsg(result),
		Provider: string(ProviderTypeSmsbao),
	}
}

// smsbaoErrorMsg 错误码转中文
func smsbaoErrorMsg(code string) string {
	switch code {
	case "0":
		return "短信发送成功"
	case "-1":
		return "参数不全"
	case "-2":
		return "服务器空间不支持,请确认支持curl或者fsocket，联系您的空间商解决或者更换空间！"
	case "30":
		return "密码错误"
	case "40":
		return "账号不存在"
	case "41":
		return "余额不足"
	case "42":
		return "帐户已过期"
	case "43":
		return "IP地址限制"
	case "50":
		return "内容含有敏感词"
	case "51":
		return "手机号码不正确"
	default:
		return "未知错误"
	}
}
