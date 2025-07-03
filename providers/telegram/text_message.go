package telegram

import (
	"strings"

	"github.com/shellvon/go-sender/core"
)

// TextMessage represents a text message for Telegram
// Based on SendMessageParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendmessage
type TextMessage struct {
	BaseMessage

	// Text - Text of the message to be sent, 1-4096 characters after entities parsing.
	// Required. The actual message content.
	Text string `json:"text"`

	// ParseMode - Mode for parsing entities in the message text. see https://core.telegram.org/bots/api#formatting-options
	//
	// Can be:
	//  - "HTML": Use HTML-style formatting (<b>bold</b>, <i>italic</i>, etc.)
	//  - "Markdown": This is a legacy mode, retained for backward compatibility. To use this mode, pass Markdown in the parse_mode field, Use Markdown-style formatting (*bold*, _italic_, etc.)
	//  - "MarkdownV2": Use MarkdownV2-style formatting (more strict)
	ParseMode string `json:"parse_mode,omitempty"`

	// Entities - A JSON-serialized list of special entities that appear in message text, which can be specified instead of parse_mode.
	// Optional. Allows for more precise control over text formatting and special entities.
	Entities []MessageEntity `json:"entities,omitempty"`

	// LinkPreviewOptions - Link preview generation options for the message.
	// Optional. Controls how link previews are generated and displayed.
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

// NewTextMessage creates a new TextMessage instance.
// Based on SendMessageParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendmessage
//   - Only chat_id and text are required.
func NewTextMessage(chatID string, text string) *TextMessage {
	return Text().Chat(chatID).Text(text).Build()
}

// Validate checks if the message is valid.
func (m *TextMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if strings.TrimSpace(m.Text) == "" {
		return core.NewParamError("text cannot be empty")
	}
	return nil
}
