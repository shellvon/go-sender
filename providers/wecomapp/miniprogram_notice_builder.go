package wecomapp

// MiniprogramNoticeBuilder 提供流畅的API来构造企业微信应用小程序通知消息
//
// 使用示例:
//
//	msg := wecomapp.MiniprogramNotice().
//	         AppID("wx123456789").
//	         Title("Order Status Update").
//	         Description("Your order has been shipped").
//	         Page("pages/order/detail?id=123").
//	         EmphasisFirst("Status", "Shipped").
//	         EmphasisSecond("Tracking", "SF1234567890").
//	         AddContentItem("Order ID", "ORD123456").
//	         AddContentItem("Estimated Delivery", "2024-01-15").
//	         ToUser("user1|user2").
//	         Build()
//
// 注意: AgentID会在发送过程中从账号配置自动设置
//
// 这遵循与其他provider相同的构建器风格模式以保持一致性
type MiniprogramNoticeBuilder struct {
	appID                  string
	page                   string
	title                  string
	description            string
	emphasisFirstItem      *MiniprogramNoticeEmphasisFirstItem
	emphasisSecondItem     *MiniprogramNoticeEmphasisSecondItem
	contentItems           []*MiniprogramNoticeContentItem
	toUser                 string
	toParty                string
	toTag                  string
	safe                   int
	enableIDTrans          int
	enableDuplicateCheck   int
	duplicateCheckInterval int
}

// MiniprogramNotice 创建新的MiniprogramNoticeBuilder实例
func MiniprogramNotice() *MiniprogramNoticeBuilder {
	return &MiniprogramNoticeBuilder{}
}

// AppID sets the mini-program app ID (max 32 bytes).
func (b *MiniprogramNoticeBuilder) AppID(appID string) *MiniprogramNoticeBuilder {
	b.appID = appID
	return b
}

// Page sets the mini-program page path (max 128 bytes).
func (b *MiniprogramNoticeBuilder) Page(page string) *MiniprogramNoticeBuilder {
	b.page = page
	return b
}

// Title sets the mini-program notice title (max 64 bytes).
func (b *MiniprogramNoticeBuilder) Title(title string) *MiniprogramNoticeBuilder {
	b.title = title
	return b
}

// Description sets the mini-program notice description (max 600 bytes).
func (b *MiniprogramNoticeBuilder) Description(description string) *MiniprogramNoticeBuilder {
	b.description = description
	return b
}

// EmphasisFirst 设置第一个强调项（值最大10字节）
func (b *MiniprogramNoticeBuilder) EmphasisFirst(key, value string) *MiniprogramNoticeBuilder {
	b.emphasisFirstItem = &MiniprogramNoticeEmphasisFirstItem{
		Key:   key,
		Value: value,
	}
	return b
}

// EmphasisSecond 设置第二个强调项（值最大30字节）
func (b *MiniprogramNoticeBuilder) EmphasisSecond(key, value string) *MiniprogramNoticeBuilder {
	b.emphasisSecondItem = &MiniprogramNoticeEmphasisSecondItem{
		Key:   key,
		Value: value,
	}
	return b
}

// AddContentItem adds a content item to the notice (value max 200 bytes).
func (b *MiniprogramNoticeBuilder) AddContentItem(key, value string) *MiniprogramNoticeBuilder {
	b.contentItems = append(b.contentItems, &MiniprogramNoticeContentItem{
		Key:   key,
		Value: value,
	})
	return b
}

// ContentItems sets the complete content items slice (overwrites previous items).
func (b *MiniprogramNoticeBuilder) ContentItems(items []*MiniprogramNoticeContentItem) *MiniprogramNoticeBuilder {
	b.contentItems = items
	return b
}

// ToUser sets the user IDs to send to, separated by "|". Use "@all" to send to all users.
func (b *MiniprogramNoticeBuilder) ToUser(toUser string) *MiniprogramNoticeBuilder {
	b.toUser = toUser
	return b
}

// ToParty sets the department IDs to send to, separated by "|".
func (b *MiniprogramNoticeBuilder) ToParty(toParty string) *MiniprogramNoticeBuilder {
	b.toParty = toParty
	return b
}

// ToTag sets the tag IDs to send to, separated by "|".
func (b *MiniprogramNoticeBuilder) ToTag(toTag string) *MiniprogramNoticeBuilder {
	b.toTag = toTag
	return b
}

// Safe sets whether to enable safe mode (0: no, 1: yes).
func (b *MiniprogramNoticeBuilder) Safe(safe int) *MiniprogramNoticeBuilder {
	b.safe = safe
	return b
}

// EnableIDTrans sets whether to enable ID translation (0: no, 1: yes).
func (b *MiniprogramNoticeBuilder) EnableIDTrans(enable int) *MiniprogramNoticeBuilder {
	b.enableIDTrans = enable
	return b
}

// EnableDuplicateCheck sets whether to enable duplicate message check (0: no, 1: yes).
func (b *MiniprogramNoticeBuilder) EnableDuplicateCheck(enable int) *MiniprogramNoticeBuilder {
	b.enableDuplicateCheck = enable
	return b
}

// DuplicateCheckInterval sets the duplicate check interval in seconds (max 4 hours).
func (b *MiniprogramNoticeBuilder) DuplicateCheckInterval(interval int) *MiniprogramNoticeBuilder {
	b.duplicateCheckInterval = interval
	return b
}

// Build assembles a ready-to-send *MiniprogramNoticeMessage.
func (b *MiniprogramNoticeBuilder) Build() *MiniprogramNoticeMessage {
	return &MiniprogramNoticeMessage{
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
			MsgType: TypeMiniprogramNotice,
		},
		MiniprogramNotice: MiniprogramNoticeMessageContent{
			AppID:              b.appID,
			Page:               b.page,
			Title:              b.title,
			Description:        b.description,
			EmphasisFirstItem:  b.emphasisFirstItem,
			EmphasisSecondItem: b.emphasisSecondItem,
			ContentItem:        b.contentItems,
		},
	}
}
