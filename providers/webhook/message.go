package webhook

import (
	"github.com/shellvon/go-sender/core"
)

// Message represents a webhook message
type Message struct {
	core.DefaultMessage
	EndpointName string            `json:"endpoint_name"` // Specify which endpoint to use
	Body         interface{}       `json:"body"`          // Request body
	Headers      map[string]string `json:"headers"`       // Dynamic request headers
	QueryParams  map[string]string `json:"query_params"`  // Dynamic query parameters
	PathVars     map[string]string `json:"path_vars"`     // URL path variables for replacement
}

var (
	_ core.Message = (*Message)(nil)
)

// ProviderType returns the provider type for this message.
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeWebhook
}

// Validate checks if the Message is valid
func (m *Message) Validate() error {
	if m.EndpointName == "" {
		return core.NewParamError("endpoint_name cannot be empty")
	}
	if m.Body == nil {
		return core.NewParamError("body cannot be nil")
	}
	return nil
}

// MessageOption defines a function type for configuring Message
type MessageOption func(*Message)

// WithEndpointName sets the endpoint name for Message
func WithEndpointName(endpointName string) MessageOption {
	return func(m *Message) {
		m.EndpointName = endpointName
	}
}

// WithHeaders sets the dynamic headers for Message
func WithHeaders(headers map[string]string) MessageOption {
	return func(m *Message) {
		m.Headers = headers
	}
}

// WithQueryParams sets the dynamic query parameters for Message
func WithQueryParams(queryParams map[string]string) MessageOption {
	return func(m *Message) {
		m.QueryParams = queryParams
	}
}

// WithPathVars sets the path variables for Message
func WithPathVars(pathVars map[string]string) MessageOption {
	return func(m *Message) {
		m.PathVars = pathVars
	}
}

// NewMessage creates a new Message with required fields and optional configurations
func NewMessage(endpointName string, body interface{}, opts ...MessageOption) *Message {
	m := &Message{
		EndpointName: endpointName,
		Body:         body,
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(m)
	}

	return m
}
