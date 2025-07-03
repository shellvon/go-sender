package webhook

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/shellvon/go-sender/core"
)

// Message represents a webhook message.
type Message struct {
	core.DefaultMessage

	Body        []byte            `json:"body"`                   // Request body
	Headers     map[string]string `json:"headers,omitempty"`      // Additional headers to send with the request
	Method      string            `json:"method,omitempty"`       // HTTP method (overrides endpoint method)
	PathParams  map[string]string `json:"path_params,omitempty"`  // Path variables to replace in URL
	QueryParams map[string]string `json:"query_params,omitempty"` // Query parameters to add to URL
}

// buildURL constructs the final URL by replacing path variables and adding query parameters.
func (m *Message) buildURL(baseURL string) (string, error) {
	// Replace path variables in the URL
	urlStr := baseURL
	for key, value := range m.PathParams {
		placeholder := fmt.Sprintf("{%s}", key)
		if !strings.Contains(urlStr, placeholder) {
			return "", fmt.Errorf("path parameter '%s' not found in URL template: %s", key, baseURL)
		}
		urlStr = strings.ReplaceAll(urlStr, placeholder, value)
	}

	// Parse the URL to add query parameters
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	query := parsedURL.Query()
	for key, value := range m.QueryParams {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

// Validate validates the webhook message.
// Method can be empty and may be set by the sender/provider if not specified.
// If provider config method is empty, it will be set to http.MethodPost.
// so this method will always return nil.
func (m *Message) Validate() error {
	return nil
}

// ProviderType returns the provider type for webhook messages.
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeWebhook
}

// WebhookBuilder provides a chainable API for constructing webhook messages.
type WebhookBuilder struct {
	body        []byte
	headers     map[string]string
	method      string
	pathParams  map[string]string
	queryParams map[string]string
}

// Webhook creates a new WebhookBuilder.
func Webhook() *WebhookBuilder {
	return &WebhookBuilder{
		headers:     make(map[string]string),
		pathParams:  make(map[string]string),
		queryParams: make(map[string]string),
	}
}

// Body sets the request body.
func (b *WebhookBuilder) Body(body []byte) *WebhookBuilder {
	b.body = body
	return b
}

// Method sets the HTTP method. Should use http.MethodXXX constants.
func (b *WebhookBuilder) Method(method string) *WebhookBuilder {
	b.method = method
	return b
}

// Header sets a single header key-value pair.
func (b *WebhookBuilder) Header(key, value string) *WebhookBuilder {
	b.headers[key] = value
	return b
}

// Headers sets multiple headers at once.
func (b *WebhookBuilder) Headers(headers map[string]string) *WebhookBuilder {
	for k, v := range headers {
		b.headers[k] = v
	}
	return b
}

// PathParam sets a single path parameter.
func (b *WebhookBuilder) PathParam(key, value string) *WebhookBuilder {
	b.pathParams[key] = value
	return b
}

// PathParams sets multiple path parameters at once.
func (b *WebhookBuilder) PathParams(params map[string]string) *WebhookBuilder {
	for k, v := range params {
		b.pathParams[k] = v
	}
	return b
}

// Query sets a single query parameter.
func (b *WebhookBuilder) Query(key, value string) *WebhookBuilder {
	b.queryParams[key] = value
	return b
}

// Queries sets multiple query parameters at once.
func (b *WebhookBuilder) Queries(params map[string]string) *WebhookBuilder {
	for k, v := range params {
		b.queryParams[k] = v
	}
	return b
}

// Build creates the webhook Message instance.
func (b *WebhookBuilder) Build() *Message {
	return &Message{
		Body:        b.body,
		Headers:     b.headers,
		Method:      b.method,
		PathParams:  b.pathParams,
		QueryParams: b.queryParams,
	}
}
