package wecomapp

// TextCardBuilder 提供流畅的API来构造企业微信应用文本卡片消息
//
// 使用示例:
//
//	msg := wecomapp.TextCard().
//	         Title("System Alert").
//	         Description("Server CPU usage is above 90%. Please check immediately.").
//	         URL("https://monitor.example.com").
//	         BtnTxt("View Details").
//	         ToUser("user1|user2").
//	         Build()
//
// 注意: AgentID会在发送过程中从账号配置自动设置
//
// 这遵循与其他provider相同的构建器风格模式以保持一致性
type TextCardBuilder struct {
	title                  string
	description            string
	url                    string
	btnTxt                 string
	toUser                 string
	toParty                string
	toTag                  string
	safe                   int
	enableIDTrans          int
	enableDuplicateCheck   int
	duplicateCheckInterval int
}

// TextCard 创建新的TextCardBuilder实例
func TextCard() *TextCardBuilder {
	return &TextCardBuilder{}
}

// Title 设置文本卡片标题（最大128字节）
func (b *TextCardBuilder) Title(title string) *TextCardBuilder {
	b.title = title
	return b
}

// Description 设置文本卡片描述（最大512字节）
func (b *TextCardBuilder) Description(description string) *TextCardBuilder {
	b.description = description
	return b
}

// URL 设置点击卡片时跳转的URL（最大2048字节）
func (b *TextCardBuilder) URL(url string) *TextCardBuilder {
	b.url = url
	return b
}

// BtnTxt 设置按钮文本。如果未指定，默认为"详情"
func (b *TextCardBuilder) BtnTxt(btnTxt string) *TextCardBuilder {
	b.btnTxt = btnTxt
	return b
}

// ToUser 设置发送给的用户ID，用"|"分隔。使用"@all"发送给所有用户
func (b *TextCardBuilder) ToUser(toUser string) *TextCardBuilder {
	b.toUser = toUser
	return b
}

// ToParty 设置发送给的部门ID，用"|"分隔
func (b *TextCardBuilder) ToParty(toParty string) *TextCardBuilder {
	b.toParty = toParty
	return b
}

// ToTag 设置发送给的标签ID，用"|"分隔
func (b *TextCardBuilder) ToTag(toTag string) *TextCardBuilder {
	b.toTag = toTag
	return b
}

// Safe 设置是否启用安全模式（0：否，1：是）
func (b *TextCardBuilder) Safe(safe int) *TextCardBuilder {
	b.safe = safe
	return b
}

// EnableIDTrans 设置是否启用ID转换（0：否，1：是）
func (b *TextCardBuilder) EnableIDTrans(enable int) *TextCardBuilder {
	b.enableIDTrans = enable
	return b
}

// EnableDuplicateCheck sets whether to enable duplicate message check (0: no, 1: yes).
func (b *TextCardBuilder) EnableDuplicateCheck(enable int) *TextCardBuilder {
	b.enableDuplicateCheck = enable
	return b
}

// DuplicateCheckInterval sets the duplicate check interval in seconds (max 4 hours).
func (b *TextCardBuilder) DuplicateCheckInterval(interval int) *TextCardBuilder {
	b.duplicateCheckInterval = interval
	return b
}

// Build assembles a ready-to-send *TextCardMessage.
func (b *TextCardBuilder) Build() *TextCardMessage {
	return &TextCardMessage{
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
			MsgType: TypeTextCard,
		},
		TextCard: TextCardMessageContent{
			Title:       b.title,
			Description: b.description,
			URL:         b.url,
			BtnTxt:      b.btnTxt,
		},
	}
}
