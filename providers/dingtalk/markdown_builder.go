package dingtalk

// MarkdownBuilder provides a fluent API to construct DingTalk markdown messages.
//
// Example:
//
//	msg := dingtalk.Markdown().
//	         Title("Alarm").
//	         Text("**service** down").
//	         AtAll().
//	         Build()
//
// This builder is additive; existing NewMarkdownMessage constructor is kept
// for compatibility.
type MarkdownBuilder struct {
	title     string
	text      string
	atMobiles []string
	atUserIDs []string
	isAtAll   bool
}

// Markdown creates a new MarkdownBuilder instance.
func Markdown() *MarkdownBuilder { return &MarkdownBuilder{} }

// Title sets the markdown title.
func (b *MarkdownBuilder) Title(t string) *MarkdownBuilder { b.title = t; return b }

// Text sets the markdown body text.
func (b *MarkdownBuilder) Text(markdown string) *MarkdownBuilder {
	b.text = markdown
	return b
}

// AtMobiles sets mobile numbers to mention.
func (b *MarkdownBuilder) AtMobiles(mobiles []string) *MarkdownBuilder {
	b.atMobiles = mobiles
	return b
}

// AtUserIDs sets user IDs to mention.
func (b *MarkdownBuilder) AtUserIDs(ids []string) *MarkdownBuilder {
	b.atUserIDs = ids
	return b
}

// AtAll marks the message to mention everyone.
func (b *MarkdownBuilder) AtAll() *MarkdownBuilder { b.isAtAll = true; return b }

// Build assembles a *MarkdownMessage.
func (b *MarkdownBuilder) Build() *MarkdownMessage {
	return &MarkdownMessage{
		BaseMessage: BaseMessage{MsgType: TypeMarkdown},
		Markdown: MarkdownContent{
			Title:     b.title,
			Text:      b.text,
			AtMobiles: b.atMobiles,
			AtUserIDs: b.atUserIDs,
			IsAtAll:   b.isAtAll,
		},
	}
}
