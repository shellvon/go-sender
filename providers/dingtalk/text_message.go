package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// Text represents the text content and mentions for a DingTalk message.
type Text struct {
	// Content of the text message
	Content string `json:"content"`
	// List of user IDs to mention in the group chat (@member)
	AtMobiles []string `json:"atMobiles,omitempty"`
	// List of user IDs to mention in the group chat (@member)
	AtUserIds []string `json:"atUserIds,omitempty"`
	// Whether to mention everyone
	IsAtAll bool `json:"isAtAll,omitempty"`
}

// TextMessage represents a text message for DingTalk.
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access
type TextMessage struct {
	BaseMessage
	Text Text `json:"text"`
}

// Validate validates the TextMessage to ensure it meets DingTalk API requirements.
func (m *TextMessage) Validate() error {
	if m.Text.Content == "" {
		return core.NewParamError("text content cannot be empty")
	}
	// DingTalk text message content has a maximum length of 2048 bytes
	if len(m.Text.Content) > 2048 {
		return core.NewParamError("text content exceeds 2048 bytes")
	}

	return nil
}

// TextMessageOption defines a function type for configuring TextMessage.
type TextMessageOption func(*TextMessage)

// WithAtMobiles sets the AtMobiles for TextMessage.
func WithAtMobiles(mobiles []string) TextMessageOption {
	return func(m *TextMessage) {
		m.Text.AtMobiles = mobiles
	}
}

// WithAtUserIds sets the AtUserIds for TextMessage.
func WithAtUserIds(userIds []string) TextMessageOption {
	return func(m *TextMessage) {
		m.Text.AtUserIds = userIds
	}
}

// WithIsAtAll sets the IsAtAll for TextMessage.
func WithIsAtAll(isAtAll bool) TextMessageOption {
	return func(m *TextMessage) {
		m.Text.IsAtAll = isAtAll
	}
}

// NewTextMessage creates a new TextMessage with required content and applies optional configurations.
func NewTextMessage(content string, opts ...TextMessageOption) *TextMessage {
	msg := &TextMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeText,
		},
		Text: Text{
			Content: content,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}
