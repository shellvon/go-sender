package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

const maxMarkdownTextLength = 2048

// Markdown represents the markdown content and mentions for a DingTalk message.
type Markdown struct {
	// Title of the markdown message
	Title string `json:"title"`
	// Content of the markdown message
	Text string `json:"text"`
	// List of user IDs to mention in the group chat (@member)
	AtMobiles []string `json:"atMobiles,omitempty"`
	// AtUserIDs specifies the user IDs to @mention in the markdown message
	AtUserIDs []string `json:"atUserIds,omitempty"`
	// Whether to mention everyone
	IsAtAll bool `json:"isAtAll,omitempty"`
}

// MarkdownMessage represents a markdown message for DingTalk.
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access
type MarkdownMessage struct {
	BaseMessage

	Markdown Markdown `json:"markdown"`
}

// NewMarkdownMessage creates a new MarkdownMessage with required content and applies optional configurations.
func NewMarkdownMessage(title, text string, opts ...MarkdownMessageOption) *MarkdownMessage {
	msg := &MarkdownMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeMarkdown,
		},
		Markdown: Markdown{
			Title: title,
			Text:  text,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// Validate validates the MarkdownMessage to ensure it meets DingTalk API requirements.
func (m *MarkdownMessage) Validate() error {
	if m.Markdown.Title == "" {
		return core.NewParamError("markdown title cannot be empty")
	}
	if m.Markdown.Text == "" {
		return core.NewParamError("markdown text cannot be empty")
	}
	// DingTalk markdown message content has a maximum length of 2048 bytes
	if len(m.Markdown.Text) > maxMarkdownTextLength {
		return core.NewParamError("markdown text exceeds 2048 bytes")
	}

	return nil
}

// MarkdownMessageOption defines a function type for configuring MarkdownMessage.
type MarkdownMessageOption func(*MarkdownMessage)

// WithMarkdownAtMobiles sets the AtMobiles for MarkdownMessage.
func WithMarkdownAtMobiles(mobiles []string) MarkdownMessageOption {
	return func(m *MarkdownMessage) {
		m.Markdown.AtMobiles = mobiles
	}
}

// WithMarkdownAtUserIDs sets the user IDs to @mention in the markdown message.
func WithMarkdownAtUserIDs(userIDs []string) MarkdownMessageOption {
	return func(m *MarkdownMessage) {
		m.Markdown.AtUserIDs = userIDs
	}
}

// WithMarkdownIsAtAll sets the IsAtAll for MarkdownMessage.
func WithMarkdownIsAtAll(isAtAll bool) MarkdownMessageOption {
	return func(m *MarkdownMessage) {
		m.Markdown.IsAtAll = isAtAll
	}
}
