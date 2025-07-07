package serverchan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/shellvon/go-sender/core"
)

// serverchanTransformer implements core.HTTPTransformer[*Account] for ServerChan messages.
// See API documentation: https://sct.ftqq.com/sendkey

type serverchanTransformer struct{}

func newServerChanTransformer() core.HTTPTransformer[*Account] {
	return &serverchanTransformer{}
}

func (t *serverchanTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeServerChan
}

// Transform constructs a ServerChan HTTPRequestSpec.
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
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}
	return reqSpec, core.NewResponseHandler(&core.ResponseHandlerConfig{
		SuccessField:      "code",
		SuccessValue:      "0",
		ErrorCodeField:    "code",
		ErrorMessageField: "message",
	}), nil
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
