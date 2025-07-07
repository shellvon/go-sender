package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// TextMessage represents a text message for Lark/Feishu.
type TextMessage struct {
	*core.DefaultMessage // Embed DefaultMessage

	BaseMessage *BaseMessage

	Content TextContent `json:"content"`
}

// TextContent represents the content of a text message.
type TextContent struct {
	Text string `json:"text"`
}

// TextBuilder provides a fluent API to construct Lark text messages.
type TextBuilder struct {
	text string
}

// Text creates a new TextBuilder instance.
func Text() *TextBuilder { return &TextBuilder{} }

// Content sets the text content.
func (b *TextBuilder) Content(text string) *TextBuilder {
	b.text = text
	return b
}

// Build assembles a *TextMessage.
func (b *TextBuilder) Build() *TextMessage {
	return &TextMessage{
		DefaultMessage: &core.DefaultMessage{}, // Correctly initialize embedded struct pointer
		BaseMessage:    &BaseMessage{MsgType: TypeText},
		Content:        TextContent{Text: b.text},
	}
}

// NewTextMessage creates a new text message.
func NewTextMessage(text string) *TextMessage {
	return Text().Content(text).Build()
}

// GetMsgType returns the message type.
func (m *TextMessage) GetMsgType() MessageType {
	return TypeText
}

// ProviderType returns the provider type.
func (m *TextMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the text message.
func (m *TextMessage) Validate() error {
	if m.Content.Text == "" {
		return errors.New("text content cannot be empty")
	}
	return nil
}
