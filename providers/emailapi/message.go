package emailapi

import (
	"github.com/shellvon/go-sender/core"
)

// Message represents a unified email message for API-based providers.
type Message struct {
	core.DefaultMessage
	To           []string               `json:"to"`
	Cc           []string               `json:"cc,omitempty"`
	Bcc          []string               `json:"bcc,omitempty"`
	Subject      string                 `json:"subject"`
	Text         string                 `json:"text,omitempty"`
	HTML         string                 `json:"html,omitempty"`
	From         string                 `json:"from,omitempty"`
	ReplyTo      string                 `json:"reply_to,omitempty"`
	Attachments  []Attachment           `json:"attachments,omitempty"`
	Headers      map[string]string      `json:"headers,omitempty"`
	TemplateID   string                 `json:"template_id,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Extras       map[string]interface{} `json:"extras,omitempty"`
}

// Attachment represents an email attachment.
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
}

// ProviderType returns the provider type for this message
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeEmailAPI
}

func (m *Message) Validate() error {
	return nil
}

func (m *Message) MsgID() string {
	return m.DefaultMessage.MsgID()
}
