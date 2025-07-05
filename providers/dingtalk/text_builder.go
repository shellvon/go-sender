package dingtalk

// TextBuilder provides a fluent API to construct DingTalk text messages.
//
// Example:
//
//	msg := dingtalk.Text().
//	         Content("hello").
//	         AtMobiles([]string{"***REMOVED***"}).
//	         Build()
//
// The builder is optional – existing NewTextMessage constructor remains for
// backward-compatibility.
type TextBuilder struct {
	content   string
	atMobiles []string
	atUserIDs []string
	isAtAll   bool
}

// Text creates a new TextBuilder instance.
func Text() *TextBuilder { return &TextBuilder{} }

// Content sets the required text content (max 2048 bytes).
func (b *TextBuilder) Content(c string) *TextBuilder {
	b.content = c
	return b
}

// AtMobiles sets the list of mobile numbers to mention ("@member").
func (b *TextBuilder) AtMobiles(mobiles []string) *TextBuilder {
	b.atMobiles = mobiles
	return b
}

// AtUserIDs sets the list of user IDs to mention.
func (b *TextBuilder) AtUserIDs(ids []string) *TextBuilder {
	b.atUserIDs = ids
	return b
}

// AtAll marks the message to mention everyone (@all).
func (b *TextBuilder) AtAll() *TextBuilder {
	b.isAtAll = true
	return b
}

// Build assembles a *TextMessage ready for sending.
func (b *TextBuilder) Build() *TextMessage {
	return &TextMessage{
		BaseMessage: BaseMessage{MsgType: TypeText},
		Text: TextContent{
			Content:   b.content,
			AtMobiles: b.atMobiles,
			AtUserIDs: b.atUserIDs,
			IsAtAll:   b.isAtAll,
		},
	}
}
