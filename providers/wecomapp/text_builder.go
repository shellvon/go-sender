package wecomapp

// TextBuilder 提供流畅的API来构造企业微信应用文本消息
//
// 使用示例:
//
//	msg := wecomapp.Text().
//	         Content("Hello World").
//	         ToUser("user1|user2").
//	         Build()
//
// 注意: AgentID会在发送过程中从账号配置自动设置
//
// 这遵循与其他provider相同的构建器风格模式以保持一致性
type TextBuilder struct {
	content                string
	toUser                 string
	toParty                string
	toTag                  string
	safe                   int
	enableIDTrans          int
	enableDuplicateCheck   int
	duplicateCheckInterval int
}

// Text 创建新的TextBuilder实例
func Text() *TextBuilder {
	return &TextBuilder{}
}

// Content 设置文本消息内容
func (b *TextBuilder) Content(content string) *TextBuilder {
	b.content = content
	return b
}

// ToUser 设置发送给的用户ID，用"|"分隔。使用"@all"发送给所有用户
func (b *TextBuilder) ToUser(toUser string) *TextBuilder {
	b.toUser = toUser
	return b
}

// ToParty 设置发送给的部门ID，用"|"分隔
func (b *TextBuilder) ToParty(toParty string) *TextBuilder {
	b.toParty = toParty
	return b
}

// ToTag 设置发送给的标签ID，用"|"分隔
func (b *TextBuilder) ToTag(toTag string) *TextBuilder {
	b.toTag = toTag
	return b
}

// Safe 设置是否启用安全模式（0：否，1：是）
func (b *TextBuilder) Safe(safe int) *TextBuilder {
	b.safe = safe
	return b
}

// EnableIDTrans 设置是否启用ID转换（0：否，1：是）
func (b *TextBuilder) EnableIDTrans(enable int) *TextBuilder {
	b.enableIDTrans = enable
	return b
}

// EnableDuplicateCheck 设置是否启用重复消息检查（0：否，1：是）
func (b *TextBuilder) EnableDuplicateCheck(enable int) *TextBuilder {
	b.enableDuplicateCheck = enable
	return b
}

// DuplicateCheckInterval 设置重复检查间隔（秒，最大4小时）
func (b *TextBuilder) DuplicateCheckInterval(interval int) *TextBuilder {
	b.duplicateCheckInterval = interval
	return b
}

// Build 组装一个可发送的*TextMessage
func (b *TextBuilder) Build() *TextMessage {
	return &TextMessage{
		BaseMessage: BaseMessage{
			CommonFields: CommonFields{
				ToUser:                 b.toUser,
				ToParty:                b.toParty,
				ToTag:                  b.toTag,
				Safe:                   b.safe,
				EnableIDTrans:          b.enableIDTrans,
				EnableDuplicateCheck:   b.enableDuplicateCheck,
				DuplicateCheckInterval: b.duplicateCheckInterval,
			},
			MsgType: TypeText,
		},
		Text: TextMessageContent{
			Content: b.content,
		},
	}
}
