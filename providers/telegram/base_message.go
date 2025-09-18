package telegram

import "github.com/shellvon/go-sender/core"

// BaseMessage represents a base message with common fields for all Telegram messages.
type BaseMessage struct {
	*core.BaseMessage

	MsgType MessageType `json:"msgtype"`

	// Unique identifier of the business connection on behalf of which the message will be sent
	BusinessConnectionID string `json:"business_connection_id,omitempty"`

	// Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	ChatID string `json:"chat_id"`

	// Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	MessageThreadID int `json:"message_thread_id,omitempty"`

	// Sends the message silently. Users will receive a notification with no sound.
	DisableNotification bool `json:"disable_notification,omitempty"`

	// Protects the contents of the sent message from forwarding
	ProtectContent bool `json:"protect_content,omitempty"`

	// Pass True to allow up to 1000 messages per second, ignoring broadcasting limits for a fee of 0.1 Telegram Stars per message.
	// The relevant Stars will be withdrawn from the bot's balance.
	// By default, all bots are able to broadcast up to 30 messages per second to their users.
	// Developers can increase this limit by enabling Paid Broadcasts in @Botfather - allowing their bot to broadcast up to 1000 messages per second.
	AllowPaidBroadcast bool `json:"allow_paid_broadcast,omitempty"`

	// Unique identifier of the message effect to be added to the message; for private chats only
	MessageEffectID string `json:"message_effect_id,omitempty"`

	// Description of the message to reply to
	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`

	// Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove a reply keyboard or to force a reply from the user
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

// Compile-time assertion: BaseMessage implements Message interface.
var (
	_ core.Message = (*BaseMessage)(nil)
)

// GetMsgType Implements the Message interface.
// Returns the message type.
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// MediaMessage represents common fields for media messages.
type MediaMessage struct {
	BaseMessage

	// Caption for the media, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`

	// ParseMode - Mode for parsing entities in the message text. see https://core.telegram.org/bots/api#formatting-options
	//
	// Can be:
	//  - "HTML": Use HTML-style formatting (<b>bold</b>, <i>italic</i>, etc.)
	//  - "Markdown": This is a legacy mode, retained for backward compatibility. To use this mode, pass Markdown in the parse_mode field, Use Markdown-style formatting (*bold*, _italic_, etc.)
	//  - "MarkdownV2": Use MarkdownV2-style formatting (more strict)
	ParseMode ParseMode `json:"parse_mode,omitempty"`

	// A JSON-serialized list of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	// Pass True if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
}
