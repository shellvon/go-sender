package wecombot

// TextBuilder 提供了一个流式 API 来构建企业微信文本消息。
// 示例用法：
//   msg := wecombot.Text().
//            Content("Hello").
//            MentionUsers([]string{"@all"}).
//            Build()
//
// 这使企业微信机器人的构建器风格与其它提供者（如 Telegram）保持一致，
// 使用户只需学习一种调用模式。
//
// 注意：构建器完全是可选的 - 现有的 NewTextMessage 和功能选项构造函数
// 继续有效，以保持向后兼容性。

type TextBuilder struct {
	content             string
	mentionedList       []string
	mentionedMobileList []string
}

// Text 创建一个新的 TextBuilder 实例。
// 返回值：*TextBuilder - 新创建的 TextBuilder 实例，用于构建 TextMessage。
func Text() *TextBuilder {
	return &TextBuilder{}
}

// Content 设置消息内容（原始文本，最大 2048 字节，需为 UTF-8 编码）。
// 参数：c string - 要设置的文本内容。
// 返回值：*TextBuilder - 返回 TextBuilder 实例以支持链式调用。
func (b *TextBuilder) Content(c string) *TextBuilder {
	b.content = c
	return b
}

// MentionUsers 设置用户 ID 列表（"mentioned_list"）。
// 使用 []string{"@all"} 可提及所有人。
// 参数：users []string - 要@提及的用户 ID 列表。
// 返回值：*TextBuilder - 返回 TextBuilder 实例以支持链式调用。
func (b *TextBuilder) MentionUsers(users []string) *TextBuilder {
	b.mentionedList = users
	return b
}

// MentionMobiles 设置手机号码列表（"mentioned_mobile_list"）。
// 使用 []string{"@all"} 可提及所有人。
// 参数：mobiles []string - 要@提及的手机号码列表。
// 返回值：*TextBuilder - 返回 TextBuilder 实例以支持链式调用。
func (b *TextBuilder) MentionMobiles(mobiles []string) *TextBuilder {
	b.mentionedMobileList = mobiles
	return b
}

// Build 构建并返回一个准备发送的 TextMessage 实例。
// 返回值：*TextMessage - 基于 TextBuilder 配置创建的文本消息实例。
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
