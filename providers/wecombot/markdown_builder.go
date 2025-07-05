package wecombot

// MarkdownBuilder provides a fluent API to construct WeCom markdown messages.
// Example:
//
//	msg := wecombot.Markdown().Content("**Hello** _world_").Build()
type MarkdownBuilder struct {
	content string
	version string
}

// Markdown creates a new MarkdownBuilder.
func Markdown() *MarkdownBuilder {
	return &MarkdownBuilder{version: MarkdownVersionLegacy}
}

// Content sets the markdown content (max 4096 characters).
// Based on SendMarkdownParams from WeCom Bot API
// https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
func (b *MarkdownBuilder) Content(c string) *MarkdownBuilder {
	b.content = c
	return b
}

// Version sets markdown version ("legacy" or "v2")
// Based on SendMarkdownParams from WeCom Bot API
// https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
//   - version is "legacy" if not provided or empty.
//   - version is "v2" if provided version is "v2".
//
// using constant `MarkdownVersionV2`/`MarkdownVersionLegacy` from markdown_message.go.
func (b *MarkdownBuilder) Version(v string) *MarkdownBuilder {
	b.version = v
	return b
}

// Build assembles a *MarkdownMessage ready to send.
func (b *MarkdownBuilder) Build() *MarkdownMessage {
	msgType := TypeMarkdown
	if b.version == MarkdownVersionV2 {
		msgType = MessageType("markdown_v2")
	}
	return &MarkdownMessage{
		BaseMessage: BaseMessage{MsgType: msgType},
		Markdown: MarkdownContent{
			Content: b.content,
			Version: b.version,
		},
	}
}
