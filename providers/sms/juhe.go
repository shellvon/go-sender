package sms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// https://www.juhe.cn/docs/api/id/54
// sendJuheSMS 聚合短信发送实现（仅支持单发）
func sendJuheSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	if len(msg.Mobiles) > 1 {
		return fmt.Errorf("juhe does not support batch sending, only single mobile is allowed")
	}
	if len(msg.Mobiles) != 1 {
		return fmt.Errorf("juhe only supports single mobile per request")
	}
	mobile := msg.Mobiles[0]

	vars := ""
	if len(msg.TemplateParams) > 0 {
		b, _ := json.Marshal(msg.TemplateParams)
		vars = string(b)
	}

	params := map[string]string{
		"mobile": mobile,
		"tpl_id": msg.TemplateCode,
		"vars":   vars,
		"dtype":  "json",
		"key":    provider.AppID,
	}

	// 支持扩展码 ext
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		if ext, ok := metadata["ext"].(string); ok && ext != "" {
			params["ext"] = ext
		}
	}

	resp, _, err := utils.DoRequest(ctx, "http://v.juhe.cn/sms/send", utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Data: params,
	})
	if err != nil {
		return fmt.Errorf("juhe SMS request failed: %w", err)
	}
	return parseJuheResponse(resp)
}

// parseJuheResponse 解析聚合短信响应
func parseJuheResponse(resp []byte) error {
	var result struct {
		ErrorCode int    `json:"error_code"`
		Reason    string `json:"reason"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse juhe response: %w", err)
	}
	if result.ErrorCode != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.ErrorCode),
			Message:  result.Reason,
			Provider: string(ProviderTypeJuhe),
		}
	}
	return nil
}
