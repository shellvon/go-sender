package dingtalk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shellvon/go-sender/core"
)

// dingTalkTransformer implements core.HTTPTransformer[*core.Account] for DingTalk
// It is stateless and does not hold any config.
type dingTalkTransformer struct{}

// newDingTalkTransformer creates a new DingTalk transformer (stateless)
func newDingTalkTransformer() core.HTTPTransformer[*core.Account] {
	return &dingTalkTransformer{}
}

// CanTransform checks if this transformer can handle the given message
func (t *dingTalkTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeDingtalk
}

// Transform converts a DingTalk message to HTTP request specification
func (t *dingTalkTransformer) Transform(ctx context.Context, msg core.Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	dingMsg, ok := msg.(Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for DingTalk: %T", msg)
	}

	if err := dingMsg.Validate(); err != nil {
		return nil, nil, err
	}

	if account == nil {
		return nil, nil, fmt.Errorf("no account provided")
	}

	// Build webhook URL
	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", account.Key)

	// Prepare the request payload
	payload := map[string]interface{}{
		"msgtype":                    string(dingMsg.GetMsgType()),
		string(dingMsg.GetMsgType()): dingMsg,
	}

	// Convert payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build request
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      webhookURL,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     body,
		BodyType: "json",
		Timeout:  30 * time.Second,
	}

	return reqSpec, t.handleDingTalkResponse, nil
}

// handleDingTalkResponse handles DingTalk API response
func (t *dingTalkTransformer) handleDingTalkResponse(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for DingTalk specific response format
	if errCode, ok := response["errcode"].(float64); ok {
		if errCode == 0 {
			return nil
		}
		errMsg := "unknown error"
		if msg, ok := response["errmsg"].(string); ok {
			errMsg = msg
		}
		return fmt.Errorf("dingtalk API error: code=%v, message=%s", errCode, errMsg)
	}

	return nil
}
