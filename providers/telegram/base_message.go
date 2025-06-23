package telegram

import "github.com/shellvon/go-sender/core"

// BaseMessage represents a base message with common fields for all Telegram messages.
type BaseMessage struct {
	core.DefaultMessage
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
	AllowPaidBroadcast bool `json:"allow_paid_broadcast,omitempty"`

	// Unique identifier of the message effect to be added to the message; for private chats only
	MessageEffectID string `json:"message_effect_id,omitempty"`

	// Description of the message to reply to
	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`

	// Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove a reply keyboard or to force a reply from the user
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// MessageWithBase is an interface for messages that embed BaseMessage.
type MessageWithBase interface {
	GetBase() *BaseMessage
}

// MessageOption defines an option for a message.
type MessageOption func(MessageWithBase)

// WithSilent sets the disable_notification field.
func WithSilent(silent bool) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().DisableNotification = silent
	}
}

// WithProtectContent sets the protect_content field.
func WithProtectContent(protect bool) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().ProtectContent = protect
	}
}

// WithAllowPaidBroadcast sets the allow_paid_broadcast field.
func WithAllowPaidBroadcast(allow bool) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().AllowPaidBroadcast = allow
	}
}

// WithMessageEffectID sets the message_effect_id field.
func WithMessageEffectID(effectID string) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().MessageEffectID = effectID
	}
}

// WithReplyParameters sets the reply_parameters field.
func WithReplyParameters(params *ReplyParameters) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().ReplyParameters = params
	}
}

// WithReplyMarkup sets the reply_markup field.
func WithReplyMarkup(markup ReplyMarkup) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().ReplyMarkup = markup
	}
}

// WithBusinessConnectionID sets the business_connection_id field.
func WithBusinessConnectionID(id string) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().BusinessConnectionID = id
	}
}

// WithMessageThreadID sets the message_thread_id field.
func WithMessageThreadID(threadID int) MessageOption {
	return func(m MessageWithBase) {
		m.GetBase().MessageThreadID = threadID
	}
}

// MediaMessage represents common fields for media messages
type MediaMessage struct {
	BaseMessage

	// Caption for the media, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`

	// Mode for parsing entities in the caption. See formatting options for more details on supported modes.
	// Options: "HTML", "Markdown", "MarkdownV2"
	ParseMode string `json:"parse_mode,omitempty"`

	// A JSON-serialized list of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	// Pass True if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
}

func (m *MediaMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

// WithCaption sets the caption for the media message
// Caption can be 0-1024 characters after entities parsing
// Use with ParseMode to format the caption (HTML, Markdown, MarkdownV2)
func WithCaption(caption string) MessageOption {
	return func(m MessageWithBase) {
		// Try to access MediaMessage fields through type assertions
		if p, ok := m.(*MediaMessage); ok {
			p.Caption = caption
		} else if p, ok := m.(*PhotoMessage); ok {
			p.Caption = caption
		} else if p, ok := m.(*AudioMessage); ok {
			p.Caption = caption
		} else if p, ok := m.(*DocumentMessage); ok {
			p.Caption = caption
		} else if p, ok := m.(*VideoMessage); ok {
			p.Caption = caption
		} else if p, ok := m.(*AnimationMessage); ok {
			p.Caption = caption
		} else if p, ok := m.(*VoiceMessage); ok {
			p.Caption = caption
		}
	}
}

// WithParseMode sets the parse mode for the media caption
// Supported modes: "HTML", "Markdown", "MarkdownV2"
// This enables formatting in the caption text
func WithParseMode(mode string) MessageOption {
	return func(m MessageWithBase) {
		// Try to access MediaMessage fields through type assertions
		if p, ok := m.(*MediaMessage); ok {
			p.ParseMode = mode
		} else if p, ok := m.(*PhotoMessage); ok {
			p.ParseMode = mode
		} else if p, ok := m.(*AudioMessage); ok {
			p.ParseMode = mode
		} else if p, ok := m.(*DocumentMessage); ok {
			p.ParseMode = mode
		} else if p, ok := m.(*VideoMessage); ok {
			p.ParseMode = mode
		} else if p, ok := m.(*AnimationMessage); ok {
			p.ParseMode = mode
		} else if p, ok := m.(*VoiceMessage); ok {
			p.ParseMode = mode
		} else if p, ok := m.(*TextMessage); ok {
			p.ParseMode = mode
		}
	}
}

// WithCaptionEntities sets the entities for the media caption
// A JSON-serialized list of special entities that appear in the caption
func WithCaptionEntities(entities []MessageEntity) MessageOption {
	return func(m MessageWithBase) {
		// Try to access MediaMessage fields through type assertions
		if p, ok := m.(*MediaMessage); ok {
			p.CaptionEntities = entities
		} else if p, ok := m.(*AudioMessage); ok {
			p.CaptionEntities = entities
		} else if p, ok := m.(*DocumentMessage); ok {
			p.CaptionEntities = entities
		} else if p, ok := m.(*VideoMessage); ok {
			p.CaptionEntities = entities
		} else if p, ok := m.(*AnimationMessage); ok {
			p.CaptionEntities = entities
		} else if p, ok := m.(*VoiceMessage); ok {
			p.CaptionEntities = entities
		} else if p, ok := m.(*TextMessage); ok {
			p.Entities = entities
		}
	}
}

// WithShowCaptionAboveMedia sets whether the caption should be shown above the media
// By default, captions appear below the media
func WithShowCaptionAboveMedia(show bool) MessageOption {
	return func(m MessageWithBase) {
		// Try to access MediaMessage fields through type assertions
		if p, ok := m.(*MediaMessage); ok {
			p.ShowCaptionAboveMedia = show
		} else if p, ok := m.(*PhotoMessage); ok {
			p.ShowCaptionAboveMedia = show
		} else if p, ok := m.(*VideoMessage); ok {
			p.ShowCaptionAboveMedia = show
		} else if p, ok := m.(*AnimationMessage); ok {
			p.ShowCaptionAboveMedia = show
		}
	}
}
