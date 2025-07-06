package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// RequestTransformer defines the interface for transforming webhook messages to HTTP requests.
type RequestTransformer interface {
	Transform(
		ctx context.Context,
		msg core.Message,
		endpoint *Endpoint,
	) (*core.HTTPRequestSpec, core.ResponseHandler, error)
	CanTransform(msg core.Message) bool
}

// webhookTransformer 实现 core.HTTPTransformer[*Endpoint].
type webhookTransformer struct{}

// 确保 webhookTransformer 实现了 core.HTTPTransformer[*Endpoint].
var _ core.HTTPTransformer[*Endpoint] = (*webhookTransformer)(nil)

// newWebhookTransformer constructs a new webhookTransformer.
func newWebhookTransformer() core.HTTPTransformer[*Endpoint] {
	return &webhookTransformer{}
}

// CanTransform 判断是否为 Webhook 消息.
func (t *webhookTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeWebhook
}

// Transform 构造 Webhook HTTPRequestSpec
//
// Parameters:
//   - ctx: 上下文
//   - msg: Webhook 消息体
//   - endpoint: webhook endpoint 配置
//
// Returns:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *webhookTransformer) Transform(
	_ context.Context,
	msg core.Message,
	endpoint *Endpoint,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	whMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for webhook transformer: %T", msg)
	}
	// 构造URL（支持PathParams/QueryParams）
	url := endpoint.URL
	if len(whMsg.PathParams) > 0 || len(whMsg.QueryParams) > 0 {
		builtURL, err := whMsg.buildURL(endpoint.URL)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build URL: %w", err)
		}
		url = builtURL
	}
	// 合并Headers
	headers := make(map[string]string)
	for k, v := range endpoint.Headers {
		headers[k] = v
	}
	for k, v := range whMsg.Headers {
		headers[k] = v
	}
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}
	// 选择HTTP方法
	method := endpoint.Method
	if whMsg.Method != "" {
		method = whMsg.Method
	}
	if method == "" {
		method = http.MethodPost
	}
	// Body已经是[]byte，直接使用
	reqSpec := &core.HTTPRequestSpec{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    whMsg.Body,
	}
	return reqSpec, t.buildResponseHandler(endpoint), nil
}

// buildResponseHandler 构造响应处理器，支持多种响应校验方式.
func (t *webhookTransformer) buildResponseHandler(endpoint *Endpoint) core.ResponseHandler {
	return func(statusCode int, body []byte) error {
		cfg := endpoint.ResponseConfig
		if cfg == nil || !cfg.ValidateResponse {
			return t.handleDefaultResponse(statusCode, body)
		}
		switch cfg.ResponseType {
		case core.BodyTypeJSON:
			return t.handleJSONResponse(cfg, statusCode, body)
		case core.BodyTypeText:
			return t.handleTextResponse(cfg, statusCode, body)
		case core.BodyTypeNone:
			return nil
		case core.BodyTypeXML:
			fallthrough
		case core.BodyTypeRaw:
			fallthrough
		case core.BodyTypeForm:
			// TODO: Implement form response handling if needed, or add a comment if not supported.
			// For now, just fall through or return an error if not supported.
			return errors.New("BodyTypeForm is not supported yet")
		default:
			return t.handleDefaultResponse(statusCode, body)
		}
	}
}

func (t *webhookTransformer) handleJSONResponse(cfg *ResponseConfig, _ int, body []byte) error {
	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("failed to parse webhook response: %w", err)
	}
	if cfg.SuccessField != "" {
		if v, ok := resp[cfg.SuccessField]; ok {
			if fmt.Sprintf("%v", v) != cfg.SuccessValue {
				errMsg := t.extractErrorMessage(cfg, resp)
				return fmt.Errorf("webhook returned failure: %s", errMsg)
			}
		}
	}
	return nil
}

func (t *webhookTransformer) handleTextResponse(cfg *ResponseConfig, _ int, body []byte) error {
	respText := string(body)
	if cfg.ErrorPattern != "" {
		matched, err := regexp.MatchString(cfg.ErrorPattern, respText)
		if err != nil {
			return fmt.Errorf("invalid error pattern: %w", err)
		}
		if matched {
			return fmt.Errorf("webhook returned error response: %s", respText)
		}
	}
	if cfg.SuccessPattern != "" {
		matched, err := regexp.MatchString(cfg.SuccessPattern, respText)
		if err != nil {
			return fmt.Errorf("invalid success pattern: %w", err)
		}
		if !matched {
			return fmt.Errorf("webhook response does not match success pattern: %s", respText)
		}
	}
	return nil
}

func (t *webhookTransformer) handleDefaultResponse(statusCode int, _ []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("webhook API returned non-2xx status: %d", statusCode)
	}
	return nil
}

// extractErrorMessage 从 JSON 响应中提取错误信息.
func (t *webhookTransformer) extractErrorMessage(cfg *ResponseConfig, resp map[string]interface{}) string {
	if cfg.ErrorField != "" {
		if v, ok := resp[cfg.ErrorField]; ok {
			return fmt.Sprintf("%v", v)
		}
	}
	if cfg.MessageField != "" {
		if v, ok := resp[cfg.MessageField]; ok {
			return fmt.Sprintf("%v", v)
		}
	}
	return "unknown error"
}
