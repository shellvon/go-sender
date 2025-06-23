package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

// http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms
const (
	uccEndpointTemplate = "http://open2.ucpaas.com/sms-server/%s"
)

func sendUcpSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	var api string
	if len(msg.Mobiles) == 1 {
		api = "variablesms"
	} else if len(msg.Mobiles) > 1 {
		api = "templatesms"
	} else {
		return fmt.Errorf("ucp: mobiles cannot be empty")
	}
	endpoint := fmt.Sprintf(uccEndpointTemplate, api)

	params := map[string]interface{}{
		"clientid":   provider.AppID,
		"password":   provider.AppSecret,
		"templateid": msg.TemplateCode,
		"mobile":     strings.Join(msg.Mobiles, ","),
		// 模板中的替换参数，如该模板不存在参数则无需传该参数或者参数为空，如果有多个参数则需要写在同一个字符串中，以分号分隔 （如：“a;b;c”），参数中不能含有特殊符号“【】”和“,”
		"param": strings.Join(msg.TemplateParamsArray, ";"),
	}

	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/json"},
		JSON:    params,
	})
	if err != nil {
		return fmt.Errorf("ucp SMS request failed: %w", err)
	}
	return parseUcpResponse(resp)
}

func parseUcpResponse(resp []byte) error {
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse ucp response: %w", err)
	}
	if result.Code != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.Code),
			Message:  result.Msg,
			Provider: string(ProviderTypeUcp),
		}
	}
	return nil
}
