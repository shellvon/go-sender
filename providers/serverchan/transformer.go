package serverchan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// serverchanTransformer 实现 core.HTTPTransformer[*Account]，用于 ServerChan 消息
// 详见API文档：https://sct.ftqq.com/sendkey

type serverchanTransformer struct{}

// newServerChanTransformer creates a new serverchanTransformer instance.
func newServerChanTransformer() core.HTTPTransformer[*Account] {
	return &serverchanTransformer{}
}

// CanTransform 判断是否为 ServerChan 消息.
func (t *serverchanTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeServerChan
}

// Transform 构造 ServerChan HTTPRequestSpec
// 参数:
//   - _: 上下文
//   - msg: ServerChan 消息体
//   - account: 账号配置
//
// 返回:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *serverchanTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	scMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for serverchan transformer: %T", msg)
	}
	apiURL := t.buildAPIURL(account.APIKey)
	body, err := json.Marshal(scMsg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal serverchan payload: %w", err)
	}
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      apiURL,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}
	return reqSpec, t.handleServerChanResponse, nil
}

// buildAPIURL 构造 ServerChan API 地址
// 支持企业版和专业版的 key 格式
// https://github.com/easychen/serverchan-demo/blob/481eccec5c3a6da71d6e744207a145d31a92b5a5/golang/main.go#L57C2-L69C3
func (t *serverchanTransformer) buildAPIURL(key string) string {
	// Check if key starts with 'sctp' (enterprise/professional version)
	if strings.HasPrefix(key, "sctp") {
		// Extract number from sctp{num}t format
		re := regexp.MustCompile(`^sctp(\d+)t`)
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 {
			num := matches[1]
			return fmt.Sprintf("https://%s.push.ft07.com/send/%s.send", num, key)
		}
	}
	// Standard Turbo version
	return fmt.Sprintf("https://sctapi.ftqq.com/%s.send", key)
}

// handleServerChanResponse 处理 ServerChan API 响应.
func (t *serverchanTransformer) handleServerChanResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("serverchan API returned non-OK status: %d", statusCode)
	}
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("serverchan error: code=%d, msg=%s", result.Code, result.Message)
	}
	return nil
}
