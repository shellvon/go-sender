package telegram

import "github.com/shellvon/go-sender/core"

// DiceMessage represents a dice message for Telegram
// Based on SendDiceParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddice
type DiceMessage struct {
	BaseMessage

	// Emoji on which the dice throw animation is based. Must be one of "🎲", "🎯", "🏀", "⚽", "🎳", or "🎰". Defaults to "🎲"
	Emoji string `json:"emoji,omitempty"`
}

func (m *DiceMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *DiceMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *DiceMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	return nil
}

type DiceMessageOption func(*DiceMessage)

// WithDiceEmoji sets the emoji for the dice animation
// Must be one of "🎲", "🎯", "🏀", "⚽", "🎳", or "🎰". Defaults to "🎲"
func WithDiceEmoji(emoji string) DiceMessageOption {
	return func(m *DiceMessage) { m.Emoji = emoji }
}

// NewDiceMessage creates a new DiceMessage
func NewDiceMessage(chatID string, opts ...interface{}) *DiceMessage {
	msg := &DiceMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeDice,
			ChatID:  chatID,
		},
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case DiceMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}
