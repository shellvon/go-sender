package wecomapp

// MarkdownBuilder 提供流畅的API来构造企业微信应用markdown消息
//
// 使用示例:
//
//	msg := wecomapp.Markdown().
//	         Content("# Hello World\n\nThis is **bold** text").
//	         ToUser("user1|user2").
//	         AgentID("1000001").
//	         Build()
//
// 这遵循与其他provider相同的构建器风格模式以保持一致性.
type MarkdownBuilder struct {
	content                string
	toUser                 string
	toParty                string
	toTag                  string
	agentID                string
	safe                   int
	enableIDTrans          int
	enableDuplicateCheck   int
	duplicateCheckInterval int
}

// Markdown 创建新的MarkdownBuilder实例.
func Markdown() *MarkdownBuilder {
	return &MarkdownBuilder{}
}

// Content 设置markdown消息内容（最大4096字节UTF-8）.
func (b *MarkdownBuilder) Content(content string) *MarkdownBuilder {
	b.content = content
	return b
}

// ToUser 设置发送给的用户ID，用"|"分隔。使用"@all"发送给所有用户.
func (b *MarkdownBuilder) ToUser(toUser string) *MarkdownBuilder {
	b.toUser = toUser
	return b
}

// ToParty 设置发送给的部门ID，用"|"分隔.
func (b *MarkdownBuilder) ToParty(toParty string) *MarkdownBuilder {
	b.toParty = toParty
	return b
}

// ToTag 设置发送给的标签ID，用"|"分隔.
func (b *MarkdownBuilder) ToTag(toTag string) *MarkdownBuilder {
	b.toTag = toTag
	return b
}

// AgentID 设置应用ID（必需）.
func (b *MarkdownBuilder) AgentID(agentID string) *MarkdownBuilder {
	b.agentID = agentID
	return b
}

// Safe 设置是否启用安全模式（0：否，1：是）.
func (b *MarkdownBuilder) Safe(safe int) *MarkdownBuilder {
	b.safe = safe
	return b
}

// EnableIDTrans 设置是否启用ID转换（0：否，1：是）.
func (b *MarkdownBuilder) EnableIDTrans(enable int) *MarkdownBuilder {
	b.enableIDTrans = enable
	return b
}

// EnableDuplicateCheck 设置是否启用重复消息检查（0：否，1：是）.
func (b *MarkdownBuilder) EnableDuplicateCheck(enable int) *MarkdownBuilder {
	b.enableDuplicateCheck = enable
	return b
}

// DuplicateCheckInterval 设置重复检查间隔（秒，最大4小时）.
func (b *MarkdownBuilder) DuplicateCheckInterval(interval int) *MarkdownBuilder {
	b.duplicateCheckInterval = interval
	return b
}

// Build 组装一个可发送的*MarkdownMessage.
func (b *MarkdownBuilder) Build() *MarkdownMessage {
	return &MarkdownMessage{
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
			MsgType: TypeMarkdown,
		},
		Markdown: MarkdownMessageContent{
			Content: b.content,
		},
	}
}
