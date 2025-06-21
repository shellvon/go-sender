package wecombot

import "github.com/shellvon/go-sender/core"

// Markdown represents the markdown content for a WeCom message.
type Markdown struct {
	// Content of the markdown message. Maximum length is 4096 bytes, and it must be UTF-8 encoded.
	Content string `json:"content"`
}

// MarkdownMessage represents a markdown message for WeCom.
// For more details, refer to the WeCom API documentation:
// https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
type MarkdownMessage struct {
	BaseMessage
	Markdown Markdown `json:"markdown"`
}

// Validate validates the MarkdownMessage to ensure it meets WeCom API requirements.
func (m *MarkdownMessage) Validate() error {
	if m.Markdown.Content == "" {
		return core.NewParamError("markdown content cannot be empty")
	}
	if len([]rune(m.Markdown.Content)) > 4096 {
		return core.NewParamError("markdown content exceeds 4096 characters")
	}
	return nil
}

// MarkdownMessageOption defines a function type for configuring MarkdownMessage.
type MarkdownMessageOption func(*MarkdownMessage)

// NewMarkdownMessage creates a new MarkdownMessage with the required content and applies optional configurations.
func NewMarkdownMessage(content string, opts ...MarkdownMessageOption) *MarkdownMessage {
	msg := &MarkdownMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeMarkdown,
		},
		Markdown: Markdown{
			Content: content,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}
