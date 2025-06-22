package webhook

import (
	"github.com/shellvon/go-sender/core"
)

// Message represents a webhook message
type Message struct {
	core.DefaultMessage
	Body    []byte            `json:"body"`              // Request body
	Headers map[string]string `json:"headers,omitempty"` // Additional headers to send with the request
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

// Validate validates the webhook message
func (m *Message) Validate() error {
	if len(m.Body) == 0 {
		return core.NewParamError("webhook message body cannot be empty")
	}
	return nil
}

// ProviderType returns the provider type for webhook messages
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeWebhook
}
