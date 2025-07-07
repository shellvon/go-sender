package dingtalk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// dingTalkTransformer implements core.HTTPTransformer[*Account] for DingTalk
// It is stateless and does not hold any config.
type dingTalkTransformer struct{}

// newDingTalkTransformer creates a new DingTalk transformer (stateless).
func newDingTalkTransformer() core.HTTPTransformer[*Account] {
	return &dingTalkTransformer{}
}

// CanTransform checks if this transformer can handle the given message.
func (t *dingTalkTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeDingtalk
}

// Transform converts a DingTalk message to HTTP request specification.
//   - https://open.dingtalk.com/document/orgapp/custom-bot-send-message-type
//   - https://open.dingtalk.com/document/orgapp/customize-robot-security-settings#title-7fs-kgs-36x
func (t *dingTalkTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	dingMsg, ok := msg.(Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for DingTalk: %T", msg)
	}

	if err := dingMsg.Validate(); err != nil {
		return nil, nil, err
	}

	if account == nil {
		return nil, nil, errors.New("no account provided")
	}

	// Build webhook URL with signature if secret is provided
	// https://open.dingtalk.com/document/orgapp/customize-robot-security-settings#title-7fs-kgs-36x
	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", account.APIKey)
	if account.APISecret != "" {
		timestamp := time.Now().UnixMilli()
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, account.APISecret)
		sign := utils.HMACSHA256Base64(account.APISecret, stringToSign)
		webhookURL = fmt.Sprintf("%s&timestamp=%d&sign=%s", webhookURL, timestamp, sign)
	}

	// Convert message to JSON
	body, err := json.Marshal(dingMsg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Build request
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      webhookURL,
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}

	return reqSpec, core.NewResponseHandler(&core.ResponseHandlerConfig{
		SuccessField:      "errcode",
		SuccessValue:      "0",
		ErrorCodeField:    "errcode",
		ErrorMessageField: "errmsg",
	}), nil
}
