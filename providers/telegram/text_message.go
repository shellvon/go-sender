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

	// ParseMode - Mode for parsing entities in the message text. See formatting options for more details.
	// Optional. Can be:
	// - "HTML": Use HTML-style formatting (<b>bold</b>, <i>italic</i>, etc.)
	// - "Markdown": Use Markdown-style formatting (*bold*, _italic_, etc.)
	// - "MarkdownV2": Use MarkdownV2-style formatting (more strict)
	ParseMode string `json:"parse_mode,omitempty"`

	// Entities - A JSON-serialized list of special entities that appear in message text, which can be specified instead of parse_mode.
	// Optional. Allows for more precise control over text formatting and special entities.
	Entities []MessageEntity `json:"entities,omitempty"`

	// LinkPreviewOptions - Link preview generation options for the message.
	// Optional. Controls how link previews are generated and displayed.
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`

	// ReplyToMessageID - ID of the message to reply to. Deprecated: Use ReplyParameters instead.
	// Optional. This field is kept for backward compatibility.
	ReplyToMessageID int `json:"reply_to_message_id,omitempty"`

	// DisableWebPreview - Disables link previews for links in the message. Deprecated: Use LinkPreviewOptions instead.
	// Optional. This field is kept for backward compatibility.
	DisableWebPreview bool `json:"disable_web_page_preview,omitempty"`
}

// NewTextMessage creates a new TextMessage instance.
func NewTextMessage(chatID string, text string, opts ...interface{}) *TextMessage {
	return NewTextMessageWithBuilder(chatID, text, opts...)
}

// NewTextMessageWithBuilder creates a new TextMessage with builder-style options.
func NewTextMessageWithBuilder(chatID string, text string, opts ...interface{}) *TextMessage {
	msg := &TextMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeText,
			ChatID:  chatID,
		},
		Text: text,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case TextMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}

func (m *TextMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *TextMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *TextMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if strings.TrimSpace(m.Text) == "" {
		return core.NewParamError("text cannot be empty")
	}
	return nil
}

type TextMessageOption func(*TextMessage)

// WithDisableWebPreview sets disable_web_page_preview for the message.
// disable_web_page_preview: Optional. Disables link previews for links in the message.
func WithDisableWebPreview(disable bool) TextMessageOption {
	return func(m *TextMessage) { m.DisableWebPreview = disable }
}

// WithReplyTo sets reply_to_message_id for the message (deprecated, use WithReplyParameters).
// reply_to_message_id: Optional. ID of the message to reply to. Deprecated: Use reply_parameters instead.
func WithReplyTo(replyTo int) TextMessageOption {
	return func(m *TextMessage) { m.ReplyToMessageID = replyTo }
}

// WithEntities sets entities for the message.
// entities: Optional. A JSON-serialized list of special entities that appear in message text.
func WithEntities(entities []MessageEntity) TextMessageOption {
	return func(m *TextMessage) { m.Entities = entities }
}

// WithLinkPreviewOptions sets link_preview_options for the message.
// link_preview_options: Optional. Link preview generation options for the message.
func WithLinkPreviewOptions(options *LinkPreviewOptions) TextMessageOption {
	return func(m *TextMessage) { m.LinkPreviewOptions = options }
}
