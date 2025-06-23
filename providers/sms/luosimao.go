package sms

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Luosimao SMS implementation
//
// API Documentation: https://luosimao.com/docs/api/
//
// 螺丝帽短信服务特点：
// - 支持单个发送和批量发送
// - 使用HTTP Basic Auth认证
// - 支持JSON和XML返回格式
// - 批量发送限制：单次最多10万个号码
// - 短信内容需要包含签名信息（格式：内容【签名】）
//
// 接口地址：
// - 单个发送: http://sms-api.luosimao.com/v1/send.json
// - 批量发送: http://sms-api.luosimao.com/v1/send_batch.json
// - 账户查询: http://sms-api.luosimao.com/v1/status.json
func sendLuosimaoSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	// 构建短信内容
	content := msg.Content
	if msg.TemplateCode != "" {
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	// 检查是否有签名
	if msg.SignName != "" && !strings.Contains(content, "【") {
		content = content + "【" + msg.SignName + "】"
	}

	// 确定发送方式：批量还是单个
	if len(msg.Mobiles) == 1 {
		return sendLuosimaoSingle(ctx, provider, msg.Mobiles[0], content)
	} else {
		return sendLuosimaoBatch(ctx, provider, msg.Mobiles, content)
	}
}

// sendLuosimaoSingle sends a single SMS
//
// 单个发送接口：
//   - URL: http://sms-api.luosimao.com/v1/send.json
//   - Method: POST
//   - Content-Type: application/x-www-form-urlencoded
//   - Auth: Basic Auth (api:key-{api_key})
//
// 请求参数：
//   - mobile: 目标手机号码
//   - message: 短信内容（需包含签名）
func sendLuosimaoSingle(ctx context.Context, provider *SMSProvider, mobile, content string) error {
	requestBody := map[string]string{
		"mobile":  mobile,
		"message": content,
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("api:key-"+provider.AppSecret))

	resp, _, err := utils.DoRequest(ctx, "http://sms-api.luosimao.com/v1/send.json", utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Data: requestBody,
	})

	if err != nil {
		return fmt.Errorf("luosimao SMS request failed: %w", err)
	}

	return parseLuosimaoResponse(resp)
}

// sendLuosimaoBatch sends batch SMS
//
// 批量发送接口：
//   - URL: http://sms-api.luosimao.com/v1/send_batch.json
//   - Method: POST
//   - Content-Type: application/x-www-form-urlencoded
//   - Auth: Basic Auth (api:key-{api_key})
//
// 请求参数：
//   - mobile_list: 目标手机号码列表（逗号分隔）
//   - message: 短信内容（需包含签名）
//   - time: 定时发送时间（可选）
//
// 限制：
//   - 单次提交控制在10万个号码以内
//   - 批量接口专门用于大量号码的内容群发，不建议发送验证码等有时效性要求的内容
func sendLuosimaoBatch(ctx context.Context, provider *SMSProvider, mobiles []string, content string) error {
	if len(mobiles) > 100000 {
		return fmt.Errorf("luosimao batch SMS limit exceeded: max 100,000 numbers per request")
	}

	mobileList := strings.Join(mobiles, ",")
	requestBody := map[string]string{
		"mobile_list": mobileList,
		"message":     content,
	}

	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if t, ok := metadata["time"].(string); ok && t != "" {
			requestBody["time"] = t // 螺丝帽批量接口支持定时发送，格式为"YYYY-MM-DD HH:MM:SS"
		}
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("api:key-"+provider.AppSecret))

	resp, _, err := utils.DoRequest(ctx, "http://sms-api.luosimao.com/v1/send_batch.json", utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Data: requestBody,
	})

	if err != nil {
		return fmt.Errorf("luosimao batch SMS request failed: %w", err)
	}

	return parseLuosimaoResponse(resp)
}

// 响应格式：
//
//	{
//	  "error": 0,           // 错误码，0表示成功
//	  "msg": "ok",          // 错误描述
//	  "batch_id": "...",    // 批次号（批量发送时返回）
//	  "hit": "..."          // 敏感词（error为-31时返回）
//	}
func parseLuosimaoResponse(resp []byte) error {
	var result struct {
		Error   int    `json:"error"`
		Msg     string `json:"msg"`
		BatchID string `json:"batch_id"`
		Hit     string `json:"hit"`
	}

	err := json.Unmarshal(resp, &result)
	if err != nil {
		return fmt.Errorf("failed to parse luosimao response: %w", err)
	}

	if result.Error != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.Error),
			Message:  result.Msg,
			Provider: string(ProviderTypeLuosimao),
		}
	}

	return nil
}
