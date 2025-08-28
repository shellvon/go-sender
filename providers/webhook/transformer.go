package webhook

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

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

// buildEndpointURL constructs the base URL with endpoint's query parameters.
func (wt *webhookTransformer) buildEndpointURL(endpoint *Endpoint) (string, error) {
	if len(endpoint.QueryParams) == 0 {
		return endpoint.URL, nil
	}

	// Parse the endpoint URL
	parsedURL, err := url.Parse(endpoint.URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	// Add endpoint's query parameters
	query := parsedURL.Query()
	for key, value := range endpoint.QueryParams {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

// transform builds the HTTPRequestSpec for a webhook message.
func (wt *webhookTransformer) transform(
	_ context.Context,
	msg *Message,
	endpoint *Endpoint,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// First, build the base URL with endpoint's query parameters
	baseURL, err := wt.buildEndpointURL(endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build endpoint URL: %w", err)
	}

	// Then apply message-level path params and query params
	url := baseURL
	if len(msg.PathParams) > 0 || len(msg.QueryParams) > 0 {
		builtURL, err := msg.buildURL(baseURL)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build URL: %w", err)
		}
		url = builtURL
	}

	// Merge headers using http.Header for proper case-insensitive handling
	httpHeaders := make(http.Header)

	// Add endpoint headers first
	for k, v := range endpoint.Headers {
		httpHeaders.Set(k, v)
	}

	// Add message headers (will override endpoint headers if same name, case-insensitive)
	for k, v := range msg.Headers {
		httpHeaders.Set(k, v)
	}

	// Set default Content-Type if not present
	if httpHeaders.Get("Content-Type") == "" {
		httpHeaders.Set("Content-Type", "application/json")
	}

	// Convert back to map[string]string for HTTPRequestSpec
	headers := make(map[string]string, len(httpHeaders))
	for k, v := range httpHeaders {
		if len(v) > 0 {
			headers[k] = v[0] // Take the first value
		}
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
	return reqSpec, core.NewSendResultHandler(endpoint.ResponseConfig), nil
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
