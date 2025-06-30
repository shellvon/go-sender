package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

type larkTransformer struct{}

// CanTransform 判断是否为 Lark 消息.
func (t *larkTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeLark
}

// Transform 构造 Lark HTTPRequestSpec
//   - API: https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
//
// 参数:
//   - ctx: 上下文
//   - msg: Lark 消息体
//   - account: 账号配置
//
// 返回:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *larkTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	larkMsg, ok := msg.(Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for lark transformer: %T", msg)
	}

	// 构造 webhook URL
	var webhookURL string
	var timestamp, sign string
	if account.Webhook != "" {
		webhookURL = account.Webhook
		if account.Key != "" {
			timestamp = strconv.FormatInt(time.Now().Unix(), 10)
			sign = utils.HMACSHA256Base64(account.Key, timestamp+"\n"+account.Key)
		}
	} else {
		webhookURL = fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", account.Key)
	}

	// 构造 payload
	payload := map[string]interface{}{
		"msg_type": larkMsg.GetMsgType(),
		"content":  larkMsg,
	}
	if timestamp != "" && sign != "" {
		payload["timestamp"] = timestamp
		payload["sign"] = sign
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal lark payload: %w", err)
	}

	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      webhookURL,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     body,
		BodyType: "json",
	}
	return reqSpec, t.handleLarkResponse, nil
}

// handleLarkResponse 处理 Lark API 响应.
func (t *larkTransformer) handleLarkResponse(statusCode int, body []byte) error {
	if statusCode != http.StatusOK {
		return fmt.Errorf("lark API returned non-OK status: %d", statusCode)
	}
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse lark response: %w", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("lark error: code=%d, msg=%s", result.Code, result.Message)
	}
	return nil
}

// newLarkTransformer 创建 Lark 的 transformer 实例
// 返回实现 core.HTTPTransformer[*core.Account] 的 larkTransformer.
func newLarkTransformer() core.HTTPTransformer[*core.Account] {
	return &larkTransformer{}
}
