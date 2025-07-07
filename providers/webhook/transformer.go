package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
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

// webhookTransformer implements core.HTTPTransformer[*Endpoint].
type webhookTransformer struct{}

// Ensure webhookTransformer implements core.HTTPTransformer[*Endpoint].
var _ core.HTTPTransformer[*Endpoint] = (*webhookTransformer)(nil)

func newWebhookTransformer() core.HTTPTransformer[*Endpoint] {
	return &webhookTransformer{}
}

func (t *webhookTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeWebhook
}

// Transform constructs a Webhook HTTPRequestSpec.
func (t *webhookTransformer) Transform(
	_ context.Context,
	msg core.Message,
	endpoint *Endpoint,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	whMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for webhook transformer: %T", msg)
	}

	// Build URL with PathParams and QueryParams
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

	method := endpoint.Method
	if whMsg.Method != "" {
		method = whMsg.Method
	}
	if method == "" {
		method = http.MethodPost
	}
	reqSpec := &core.HTTPRequestSpec{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    whMsg.Body,
	}
	return reqSpec, core.NewResponseHandler(endpoint.ResponseConfig), nil
}
