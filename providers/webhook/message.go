package webhook

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/shellvon/go-sender/core"
)

// Message represents a webhook message
type Message struct {
	core.DefaultMessage
	Body        []byte            `json:"body"`                   // Request body
	Headers     map[string]string `json:"headers,omitempty"`      // Additional headers to send with the request
	Method      string            `json:"method,omitempty"`       // HTTP method (overrides endpoint method)
	PathParams  map[string]string `json:"path_params,omitempty"`  // Path variables to replace in URL
	QueryParams map[string]string `json:"query_params,omitempty"` // Query parameters to add to URL
}

// NewMessage creates a new webhook message
func NewMessage(body []byte, opts ...MessageOption) *Message {
	msg := &Message{
		Body: body,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// MessageOption is a function option for webhook messages
type MessageOption func(*Message)

// WithHeaders adds headers to the webhook message
func WithHeaders(headers map[string]string) MessageOption {
	return func(m *Message) {
		if m.Headers == nil {
			m.Headers = make(map[string]string)
		}
		for k, v := range headers {
			m.Headers[k] = v
		}
	}
}

// WithMethod sets the HTTP method for this message
func WithMethod(method string) MessageOption {
	return func(m *Message) {
		m.Method = method
	}
}

// WithPathParams sets path parameters for URL template replacement
func WithPathParams(params map[string]string) MessageOption {
	return func(m *Message) {
		if m.PathParams == nil {
			m.PathParams = make(map[string]string)
		}
		for k, v := range params {
			m.PathParams[k] = v
		}
	}
}

// WithQueryParams sets query parameters for the request
func WithQueryParams(params map[string]string) MessageOption {
	return func(m *Message) {
		if m.QueryParams == nil {
			m.QueryParams = make(map[string]string)
		}
		for k, v := range params {
			m.QueryParams[k] = v
		}
	}
}

// buildURL constructs the final URL by replacing path variables and adding query parameters
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

// Validate validates the webhook message
func (m *Message) Validate() error {
	return nil
}

// ProviderType returns the provider type for webhook messages
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeWebhook
}
