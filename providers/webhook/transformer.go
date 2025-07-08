package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// Webhook is a message provider for Webhook.
// It supports sending messages to a webhook URL.
type webhookTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Endpoint]
}

// Ensure interface compliance.
var _ core.HTTPTransformer[*Endpoint] = (*webhookTransformer)(nil)

// transform builds the HTTPRequestSpec for a webhook message.
func (wt *webhookTransformer) transform(
	_ context.Context,
	msg *Message,
	endpoint *Endpoint,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// Build URL with PathParams and QueryParams
	url := endpoint.URL
	if len(msg.PathParams) > 0 || len(msg.QueryParams) > 0 {
		builtURL, err := msg.buildURL(endpoint.URL)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build URL: %w", err)
		}
		url = builtURL
	}

	// Merge headers (endpoint first, then message overrides)
	headers := make(map[string]string, len(endpoint.Headers)+len(msg.Headers))
	for k, v := range endpoint.Headers {
		headers[k] = v
	}
	for k, v := range msg.Headers {
		headers[k] = v
	}
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	method := endpoint.Method
	if msg.Method != "" {
		method = msg.Method
	}
	if method == "" {
		method = http.MethodPost
	}

	reqSpec := &core.HTTPRequestSpec{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    msg.Body,
	}
	return reqSpec, core.NewResponseHandler(endpoint.ResponseConfig), nil
}

func newWebhookTransformer() core.HTTPTransformer[*Endpoint] {
	wt := &webhookTransformer{}
	wt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeWebhook,
		"",
		nil,
		wt.transform,
	)
	return wt
}
