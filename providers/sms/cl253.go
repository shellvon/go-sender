package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// 当前文件仅实现了创蓝（CL253）国内短信协议，国际短信（如新加坡/上海节点、加签、特殊请求头等）尚未实现。
// https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
func sendCl253SMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	params := map[string]interface{}{
		"account":  provider.AppID,
		"password": provider.AppSecret,
		// 此处填写审核通过的签名和模板，中文括号是代表短信签名。内容长度支持1～3500个字符（含变量）。用营销账号提交短信时最末尾需带上退订语"拒收请回复R"不支持小写r，否则营销短信将进入人工审核。
		"msg":   msg.Content,
		"phone": strings.Join(msg.Mobiles, ","),
	}

	// 透传可选参数（report, callbackUrl, uid, sendtime, extend）
	if metadata := core.GetSendMetadataFromCtx(ctx); metadata != nil {
		// 如您需要状态回执，则需要传"true",不传默认为"false"，则无法获取状态回执。
		if v, ok := metadata["report"].(string); ok && v != "" {
			params["report"] = v
		}
		// 状态回执的回调地址，请传入完整带http开头的地址，不传默认为空，请勿传入空格，否则会造成地址推送错误。地址可通过接口入参传入，也可在控制台手动配置，可查看控制台操作指引。
		if v, ok := metadata["callbackUrl"].(string); ok && v != "" {
			params["callbackUrl"] = v
		}
		// 自定义参数，如订单号或短信发送记录流水号，最大支持64 位，状态回执会回传，不传默认为空。
		if v, ok := metadata["uid"].(string); ok && v != "" {
			params["uid"] = v
		}
		// 定时发送时间，*≤*当前时间则立即发送；只能定时 7 天内，格式为 yyyyMMddHHmm ，不传默认为空，立即发送。
		if v, ok := metadata["sendtime"].(string); ok && v != "" {
			params["sendtime"] = v
		}
		// 下发短信号码扩展码，用于匹配上行回复，上行报告会回传。一般5位以内(只支持传数字)，不传默认为空。
		if v, ok := metadata["extend"].(string); ok && v != "" {
			params["extend"] = v
		}
	}

	resp, _, err := utils.DoRequest(ctx, "https://smssh1.253.com/msg/v1/send/json", utils.RequestOptions{
		Method: "POST",
		JSON:   params,
	})
	if err != nil {
		return fmt.Errorf("cl253 request failed: %w", err)
	}

	// 5. 解析响应
	var result struct {
		Code     string `json:"code"`
		MsgId    string `json:"msgId"`
		RespTime string `json:"time"`
		ErrorMsg string `json:"errorMsg"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse cl253 response: %w", err)
	}
	if result.Code != "0" {
		return &SMSError{
			Code:     result.Code,
			Message:  result.ErrorMsg,
			Provider: string(ProviderTypeCl253),
		}
	}
	return nil
}
