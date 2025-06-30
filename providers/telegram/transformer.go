package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
)

// telegramTransformer 实现 providers.HTTPTransformer[*core.Account].
type telegramTransformer struct{}

// CanTransform 判断是否为 Telegram 消息.
func (t *telegramTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeTelegram
}

// Transform 构造 Telegram HTTPRequestSpec
// 参数:
//   - ctx: 上下文
//   - msg: Telegram 消息体
//   - account: 账号配置
//
// 返回:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *telegramTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	tgMsg, ok := msg.(Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for telegram transformer: %T", msg)
	}
	apiURL := buildAPIURL(account.Key)
	body, err := json.Marshal(tgMsg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal telegram payload: %w", err)
	}
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      apiURL,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     body,
		BodyType: "json",
	}
	return reqSpec, handleTelegramResponse, nil
}

// buildAPIURL 构造 Telegram API 地址.
func buildAPIURL(token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
}

// handleTelegramResponse 处理 Telegram API 响应.
func handleTelegramResponse(statusCode int, body []byte) error {
	if statusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned non-OK status: %d", statusCode)
	}
	var result struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code,omitempty"`
		Description string `json:"description,omitempty"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse telegram response: %w", err)
	}
	if !result.OK {
		if result.ErrorCode != 0 {
			return fmt.Errorf("telegram API error %d: %s", result.ErrorCode, result.Description)
		}
		return fmt.Errorf("telegram API request failed: %s", result.Description)
	}
	return nil
}
