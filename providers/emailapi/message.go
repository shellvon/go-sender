package emailapi

import (
	"github.com/shellvon/go-sender/core"
)

// SubProviderType represents the type of EmailAPI sub-provider
type SubProviderType string

const (
	SubProviderEmailJS SubProviderType = "emailjs"
	SubProviderResend  SubProviderType = "resend"
	// 可以继续添加其他 EmailAPI 提供商
)

// Message represents a unified email message for API-based providers
type Message struct {
	core.DefaultMessage

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
	TemplateData map[string]interface{} `json:"template_data,omitempty"` // 模板数据

	// Additional fields
	Attachments []Attachment      `json:"attachments,omitempty"` // 附件
	Headers     map[string]string `json:"headers,omitempty"`     // 自定义头部

	// Common callback field (unified interface)
	CallbackURL string `json:"callback_url,omitempty"` // 统一回调地址 - 各平台内部适配

	// Extensions for platform-specific parameters
	Extras map[string]interface{} `json:"extras"` // 扩展字段（平台特定参数）
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
}

var (
	_ core.Message = (*Message)(nil)
)

// ProviderType returns the provider type for this message
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeEmailAPI
}

// Validate checks if the Message is valid
func (m *Message) Validate() error {
	if len(m.To) == 0 {
		return core.NewParamError("to recipients cannot be empty")
	}
	if m.SubProvider == "" {
		return core.NewParamError("sub_provider must be specified for EmailAPI messages (e.g., emailjs, resend)")
	}
	if m.Subject == "" && m.TemplateID == "" {
		return core.NewParamError("either subject or template_id must be specified")
	}
	return nil
}

func (m *Message) MsgID() string {
	return m.DefaultMessage.MsgID()
}

// GetSubProvider 实现 SubProviderMessage 接口
func (m *Message) GetSubProvider() string {
	return m.SubProvider
}
