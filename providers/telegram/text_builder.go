package telegram

// TextBuilder provides chainable construction of Telegram text messages.
// Usage:
//   msg := telegram.Text().
//            Chat("123456").
//            Text("Hello World").
//            ParseMode("HTML").
//            Build()
//
// The builder ensures only Text-specific fields are exposed, giving compile-time
// safety and IDE auto-completion similar to SMS builders.

type TextBuilder struct {
	*baseBuilder[*TextBuilder]

	text               string
	parseMode          string
	entities           []MessageEntity
	linkPreviewOptions *LinkPreviewOptions
}

// Text creates a new TextBuilder.
func Text() *TextBuilder {
	b := &TextBuilder{}
	b.baseBuilder = &baseBuilder[*TextBuilder]{self: b}
	return b
}

// Text sets the message body.
// Based on SendMessageParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendmessage
// Text of the message to be sent, 1-4096 characters after entities parsing.
func (b *TextBuilder) Text(t string) *TextBuilder {
	b.text = t
	return b
}

// ParseMode - Mode for parsing entities in the message text. see https://core.telegram.org/bots/api#formatting-options
//
// Can be:
//   - "HTML": Use HTML-style formatting (<b>bold</b>, <i>italic</i>, etc.)
//   - "Markdown": This is a legacy mode, retained for backward compatibility. To use this mode, pass Markdown in the parse_mode field, Use Markdown-style formatting (*bold*, _italic_, etc.)
//   - "MarkdownV2": Use MarkdownV2-style formatting (more strict)
func (b *TextBuilder) ParseMode(mode string) *TextBuilder {
	b.parseMode = mode
	return b
}

// Entities sets entities, A JSON-serialized list of special entities that appear in message text, which can be specified instead of parse_mode.
func (b *TextBuilder) Entities(ents []MessageEntity) *TextBuilder {
	b.entities = ents
	return b
}

// LinkPreview sets link_preview_options.
// Based on SendMessageParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendmessage
// Link preview generation options for the message.
func (b *TextBuilder) LinkPreview(opts *LinkPreviewOptions) *TextBuilder {
	b.linkPreviewOptions = opts
	return b
}

// Build assembles *TextMessage ready for sending.
func (b *TextBuilder) Build() *TextMessage {
	msg := &TextMessage{
		BaseMessage:        b.baseBuilder.toBaseMessage(TypeText),
		Text:               b.text,
		ParseMode:          b.parseMode,
		Entities:           b.entities,
		LinkPreviewOptions: b.linkPreviewOptions,
	}
	return msg
}
