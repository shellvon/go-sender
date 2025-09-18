package emailapi

import (
	"time"

	"github.com/shellvon/go-sender/core"
)

// SubProviderType represents the type of EmailAPI sub-provider.
type SubProviderType string

const (
	SubProviderEmailJS    SubProviderType = "emailjs"
	SubProviderResend     SubProviderType = "resend"
	SubProviderMailerSend SubProviderType = "mailersend"
	SubProviderMailtrap   SubProviderType = "mailtrap"
	SubProviderBrevo      SubProviderType = "brevo"
	SubProviderMailgun    SubProviderType = "mailgun"
	SubProviderMailjet    SubProviderType = "mailjet"
	// 可以继续添加其他 EmailAPI 提供商.
)

// EmailType represents the type of email based on content characteristics.
type EmailType int

const (
	// EmailTypeText indicates email with text content only.
	//
	// See also: GetEmailType
	EmailTypeText EmailType = iota

	// EmailTypeHTML indicates email with HTML content only.
	//
	// See also: GetEmailType
	EmailTypeHTML

	// EmailTypeTextAndHTML indicates email with both text and HTML content.
	//
	// See also: GetEmailType
	EmailTypeTextAndHTML

	// EmailTypeTemplate indicates email using template-based content.
	//
	// See also: GetEmailType, TemplateID, TemplateData
	EmailTypeTemplate
)

// String returns a human-readable representation of the email type.
func (t EmailType) String() string {
	switch t {
	case EmailTypeText:
		return "Text"
	case EmailTypeHTML:
		return "HTML"
	case EmailTypeTextAndHTML:
		return "TextAndHTML"
	case EmailTypeTemplate:
		return "Template"
	default:
		return "Unknown"
	}
}

// Message represents a unified email message for API-based providers.
//
// The message structure supports different email types:
// - Text: contains text content only
// - HTML: contains HTML content only
// - TextAndHTML: contains both text and HTML content
// - Template: uses template-based content generation
//
// Use GetEmailType() to determine the message content type.
type Message struct {
	*core.BaseMessage
	*core.WithExtraFields // Add extra fields support for provider-specific configurations

	// Sub-provider specification (required for EmailAPI due to multiple platforms)
	SubProvider string `json:"sub_provider"` // 子提供商类型（emailjs, resend等）

	// Basic email fields
	To      []string `json:"to"`                 // 收件人
	Cc      []string `json:"cc,omitempty"`       // 抄送
	Bcc     []string `json:"bcc,omitempty"`      // 密送
	Subject string   `json:"subject"`            // 邮件主题
	From    string   `json:"from,omitempty"`     // 发件人
	ReplyTo []string `json:"reply_to,omitempty"` // 回复地址（支持多个）

	// Content fields
	Text string `json:"text,omitempty"` // 纯文本内容
	HTML string `json:"html,omitempty"` // HTML内容

	// Template related fields
	TemplateID   string                 `json:"template_id,omitempty"`   // 模板ID
	TemplateData map[string]interface{} `json:"template_data,omitempty"` // 模板数据（通用格式）
	// Note: For providers like MailerSend that support per-recipient personalization,
	// TemplateData can be structured as: key=email_address, value=personalization_data

	// Additional fields
	Attachments []Attachment      `json:"attachments,omitempty"` // 附件
	Headers     map[string]string `json:"headers,omitempty"`     // 自定义头部

	// Common callback field (unified interface)
	CallbackURL string `json:"callback_url,omitempty"` // 统一回调地址 - 各平台内部适配

	// ScheduledAt is the time to send the message.
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"` // 统一发送时间 - 各平台内部适配
}

// Attachment represents an email attachment.
// For Mailtrap, it supports content_id for inline attachments as per OpenAPI spec.
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
	// ContentID is used for inline attachments (when disposition is "inline")
	// This allows the attachment to be referenced within the email body
	ContentID string `json:"content_id,omitempty"`
	// Disposition specifies how the attachment should be displayed
	Disposition string `json:"disposition,omitempty"` // "attachment" or "inline"
}

// Compile-time assertions: Message implements core.Message, core.SubProviderAware, and core.Validatable.
var (
	_ core.Message          = (*Message)(nil)
	_ core.SubProviderAware = (*Message)(nil)
	_ core.Validatable      = (*Message)(nil)
)

// GetSubProvider Implements the SubProviderAware interface.
// Returns the sub-provider type.
func (m *Message) GetSubProvider() string {
	return m.SubProvider
}

// Validate checks if the Message is valid.
func (m *Message) Validate() error {
	if m.SubProvider == "" {
		return core.NewParamError("sub_provider must be specified for EmailAPI messages (e.g., emailjs, resend)")
	}
	return nil
}

// GetEmailType determines the email type based on content characteristics.
//
// NewMessage creates a new EmailAPI message with the specified sub-provider.
func NewMessage(subProvider string) *Message {
	return &Message{
		BaseMessage:     core.NewBaseMessage(core.ProviderTypeEmailAPI),
		WithExtraFields: core.NewWithExtraFields(),
		SubProvider:     subProvider,
	}
}

// See also: EmailTypeText, EmailTypeHTML, EmailTypeTextAndHTML, EmailTypeTemplate
func (m *Message) GetEmailType() EmailType {
	hasText := m.Text != ""
	hasHTML := m.HTML != ""
	hasTemplate := m.TemplateID != ""

	if hasTemplate {
		return EmailTypeTemplate
	}

	if hasText && hasHTML {
		return EmailTypeTextAndHTML
	}

	if hasHTML {
		return EmailTypeHTML
	}

	if hasText {
		return EmailTypeText
	}
	return EmailTypeText
}
