package telegram

// DiceBuilder constructs Telegram dice messages (random animation with score).
// Example:
//   msg := telegram.Dice().
//            Chat("123").
//            Emoji("🏀").
//            Build()

type DiceBuilder struct {
	*baseBuilder[*DiceBuilder]

	emoji string
}

// Dice returns a new DiceBuilder.
func Dice() *DiceBuilder {
	b := &DiceBuilder{}
	b.baseBuilder = &baseBuilder[*DiceBuilder]{self: b}
	return b
}

// Emoji sets the dice emoji (🎲, 🎯, 🏀, ⚽, 🎳, or 🎰).
// Based on SendDiceParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddice
// Emoji on which the dice throw animation is based.
//
// Must be one of
//   - 🎲 (default)
//   - 🎯
//   - 🏀
//   - ⚽
//   - 🎳
//   - 🎰
func (b *DiceBuilder) Emoji(e string) *DiceBuilder {
	b.emoji = e
	return b
}

// Build assembles the *DiceMessage.
func (b *DiceBuilder) Build() *DiceMessage {
	msg := &DiceMessage{
		BaseMessage: b.baseBuilder.toBaseMessage(TypeDice),
		Emoji:       b.emoji,
	}
	return msg
}
