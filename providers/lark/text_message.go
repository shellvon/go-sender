package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// TextMessage represents a text message for Lark/Feishu
type TextMessage struct {
	BaseMessage
	Content TextContent `json:"content"`
}

// TextContent represents the content of a text message
type TextContent struct {
	Text string `json:"text"`
}

// NewTextMessage creates a new text message
func NewTextMessage(text string) *TextMessage {
	return &TextMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeText,
		},
		Content: TextContent{
			Text: text,
		},
	}
}

// GetMsgType returns the message type
func (m *TextMessage) GetMsgType() MessageType {
	return TypeText
}

// ProviderType returns the provider type
func (m *TextMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the text message
func (m *TextMessage) Validate() error {
	if m.Content.Text == "" {
		return errors.New("text content cannot be empty")
	}
	return nil
}
