package dingtalk

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

const maxTextContentLength = 2048

// TextContent represents the text content and mentions for a DingTalk message.
type TextContent struct {
	// Content of the text message
	Content string `json:"content"`
	// List of user IDs to mention in the group chat (@member)
	AtMobiles []string `json:"atMobiles,omitempty"`
	// AtUserIDs specifies the user IDs to @mention in the text message
	AtUserIDs []string `json:"atUserIds,omitempty"`
	// Whether to mention everyone
	IsAtAll bool `json:"isAtAll,omitempty"`
}

// TextMessage represents a text message for DingTalk.
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access
type TextMessage struct {
	BaseMessage

	Text TextContent `json:"text"`
}

// NewTextMessage creates a new TextMessage with required content.
func NewTextMessage(content string) *TextMessage {
	return Text().Content(content).Build()
}

// Validate validates the TextMessage to ensure it meets DingTalk API requirements.
func (m *TextMessage) Validate() error {
	if m.Text.Content == "" {
		return core.NewParamError("text content cannot be empty")
	}
	// DingTalk text message content has a maximum length of 2048 bytes
	if len(m.Text.Content) > maxTextContentLength {
		return core.NewParamError(fmt.Sprintf("text content exceeds %d bytes", maxTextContentLength))
	}

	return nil
}
