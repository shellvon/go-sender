package wecombot

// MarkdownBuilder 提供了一个流式 API 来构建企业微信 Markdown 消息。
// 示例：
//
//	msg := wecombot.Markdown().Content("**Hello** _world_").Build()
type MarkdownBuilder struct {
	content string
	version MarkdownVersion
}

// Markdown 创建一个新的 MarkdownBuilder 实例。
//
// 返回值：*MarkdownBuilder - 新创建的 MarkdownBuilder 实例，用于构建 MarkdownMessage。
func Markdown() *MarkdownBuilder {
	return &MarkdownBuilder{version: MarkdownVersionLegacy}
}

// Content 设置 Markdown 消息内容（最大 4096 个字符）。
// 基于企业微信机器人 API 的 SendMarkdownParams
// 参考：https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
//
// 参数：c string - 要设置的 Markdown 内容。
//
// 返回值：*MarkdownBuilder - 返回 MarkdownBuilder 实例以支持链式调用。
func (b *MarkdownBuilder) Content(c string) *MarkdownBuilder {
	b.content = c
	return b
}

// Version 设置 Markdown 版本（"legacy" 或 "v2"）。
// 基于企业微信机器人 API 的 SendMarkdownParams
// 参考：https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
//   - 如果未提供或版本为空，则版本为 "legacy"。
//   - 如果提供的版本为 "v2"，则版本为 "v2"。
//
// 使用 markdown_message.go 中的常量 `MarkdownVersionV2`/`MarkdownVersionLegacy`。
// 更多详情请参见 https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
//
// 参数：v MarkdownVersion - 要设置的 Markdown 版本。
//
// 返回值：*MarkdownBuilder - 返回 MarkdownBuilder 实例以支持链式调用。
func (b *MarkdownBuilder) Version(v MarkdownVersion) *MarkdownBuilder {
	b.version = v
	return b
}

// Build 构建并返回一个准备发送的 MarkdownMessage 实例。
// 根据设置的版本，消息类型将为 TypeMarkdown 或 "markdown_v2"。
//
// 返回值：*MarkdownMessage - 基于 MarkdownBuilder 配置创建的 Markdown 消息实例。
func (b *MarkdownBuilder) Build() *MarkdownMessage {
	msgType := TypeMarkdown
	if b.version == MarkdownVersionV2 {
		msgType = MessageType("markdown_v2")
	}
	return &MarkdownMessage{
		BaseMessage: newBaseMessage(msgType),
		Markdown: MarkdownContent{
			Content: b.content,
			Version: b.version,
		},
	}
}
