package wecombot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shellvon/go-sender/core"
)

// wecombotTransformer implements core.HTTPTransformer[*core.Account] for WeCom Bot
type wecombotTransformer struct{}

// newWecombotTransformer creates a new WeCom Bot transformer (stateless)
func newWecombotTransformer() core.HTTPTransformer[*core.Account] {
	return &wecombotTransformer{}
}

// CanTransform checks if this transformer can handle the given message
func (t *wecombotTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeWecombot
}

// Transform converts a WeCom Bot message to HTTP request specification
func (t *wecombotTransformer) Transform(ctx context.Context, msg core.Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	wecomMsg, ok := msg.(Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for wecombot transformer: %T", msg)
	}
	if account == nil {
		return nil, nil, errors.New("no account provided")
	}

	// Build webhook URL
	webhookURL := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", account.Key)

	// Marshal message to JSON
	body, err := json.Marshal(wecomMsg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal wecombot message: %w", err)
	}

	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      webhookURL,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     body,
		BodyType: "json",
		Timeout:  30 * time.Second,
	}

	return reqSpec, t.handleWecombotResponse, nil
}

// handleWecombotResponse handles WeCom Bot API response
func (t *wecombotTransformer) handleWecombotResponse(statusCode int, body []byte) error {
	if statusCode != 200 {
		return fmt.Errorf("wecom API returned non-OK status: %d", statusCode)
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("wecom error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	return nil
}
