package telegram

import "github.com/shellvon/go-sender/core"

// DiceMessage represents a dice message for Telegram
// Based on SendDiceParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddice
type DiceMessage struct {
	BaseMessage

	// Emoji on which the dice throw animation is based. Must be one of "🎲", "🎯", "🏀", "⚽", "🎳", or "🎰". Defaults to "🎲"
	Emoji DiceEmoji `json:"emoji,omitempty"`
}

// NewDiceMessage creates a new DiceMessage instance.
// Based on SendDiceParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddice
//   - Only chat_id is required.
func NewDiceMessage(chatID string) *DiceMessage {
	return &DiceMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeDice,
			ChatID:  chatID,
		},
	}
}

func (m *DiceMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Emoji != "" && !m.Emoji.IsValid() {
		return core.NewParamError("invalid dice emoji: " + string(m.Emoji))
	}
	return nil
}
