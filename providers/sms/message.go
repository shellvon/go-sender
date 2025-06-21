package sms

import (
	"github.com/shellvon/go-sender/core"
)

// Message represents an SMS message
type Message struct {
	core.DefaultMessage
	Mobile         string            `json:"mobile"`          // Mobile phone number
	Content        string            `json:"content"`         // SMS content (for non-template SMS)
	TemplateCode   string            `json:"template_code"`   // Template code (for template SMS)
	TemplateParams map[string]string `json:"template_params"` // Template parameters
	SignName       string            `json:"sign_name"`       // SMS signature name
}

var (
	_ core.Message = (*Message)(nil)
)

// ProviderType returns the provider type for this message
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeSMS
}

// Validate checks if the Message is valid
func (m *Message) Validate() error {
	if m.Mobile == "" {
		return core.NewParamError("mobile cannot be empty")
	}

	// Either content or template_code must be provided
	if m.Content == "" && m.TemplateCode == "" {
		return core.NewParamError("either content or template_code must be provided")
	}

	// If template_code is provided, template_params should also be provided
	if m.TemplateCode != "" && len(m.TemplateParams) == 0 {
		return core.NewParamError("template_params must be provided when using template_code")
	}

	return nil
}

// MessageOption defines a function type for configuring Message
type MessageOption func(*Message)

// WithMobile sets the mobile phone number
func WithMobile(mobile string) MessageOption {
	return func(m *Message) {
		m.Mobile = mobile
	}
}

// WithContent sets the SMS content
func WithContent(content string) MessageOption {
	return func(m *Message) {
		m.Content = content
	}
}

// WithTemplateCode sets the template code
func WithTemplateCode(templateCode string) MessageOption {
	return func(m *Message) {
		m.TemplateCode = templateCode
	}
}

// WithTemplateParams sets the template parameters
func WithTemplateParams(params map[string]string) MessageOption {
	return func(m *Message) {
		m.TemplateParams = params
	}
}

// WithSignName sets the SMS signature name
func WithSignName(signName string) MessageOption {
	return func(m *Message) {
		m.SignName = signName
	}
}

// NewMessage creates a new Message with required fields and optional configurations
func NewMessage(mobile string, opts ...MessageOption) *Message {
	m := &Message{
		Mobile: mobile,
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(m)
	}

	return m
}
