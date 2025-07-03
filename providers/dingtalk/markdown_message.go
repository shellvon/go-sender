package dingtalk

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

const maxMarkdownTextLength = 2048

// MarkdownContent represents the markdown content and mentions for a DingTalk message.
type MarkdownContent struct {
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

	Markdown MarkdownContent `json:"markdown"`
}

// NewMarkdownMessage creates a new MarkdownMessage with required content and applies optional configurations.
func NewMarkdownMessage(title, text string) *MarkdownMessage {
	return Markdown().Title(title).Text(text).Build()
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
		return core.NewParamError(fmt.Sprintf("markdown text exceeds %d bytes", maxMarkdownTextLength))
	}

	return nil
}
