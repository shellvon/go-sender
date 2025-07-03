package wecombot

// TextBuilder provides a fluent API to construct WeCom text messages.
// Example usage:
//   msg := wecombot.Text().
//            Content("Hello").
//            MentionUsers([]string{"@all"}).
//            Build()
//
// This aligns WeCom Bot with the same builder-style ergonomics used by other providers
// (e.g. Telegram) so that callers only need to learn one invocation pattern.
//
// Note: The builder is completely optional – existing NewTextMessage + functional-option
// constructors continue to work for backward-compatibility.

type TextBuilder struct {
	content             string
	mentionedList       []string
	mentionedMobileList []string
}

// Text creates a new TextBuilder.
func Text() *TextBuilder {
	return &TextBuilder{}
}

// Content sets the required message content (raw text, max 2048 bytes UTF-8).
func (b *TextBuilder) Content(c string) *TextBuilder {
	b.content = c
	return b
}

// MentionUsers sets the list of user IDs ("mentioned_list").
// Use []string{"@all"} to mention everyone.
func (b *TextBuilder) MentionUsers(users []string) *TextBuilder {
	b.mentionedList = users
	return b
}

// MentionMobiles sets the list of mobile numbers ("mentioned_mobile_list").
// Use []string{"@all"} to mention everyone.
func (b *TextBuilder) MentionMobiles(mobiles []string) *TextBuilder {
	b.mentionedMobileList = mobiles
	return b
}

// Build assembles a ready-to-send *TextMessage value.
func (b *TextBuilder) Build() *TextMessage {
	return &TextMessage{
		BaseMessage: BaseMessage{MsgType: TypeText},
		Text: TextContent{
			Content:             b.content,
			MentionedList:       b.mentionedList,
			MentionedMobileList: b.mentionedMobileList,
		},
	}
}
