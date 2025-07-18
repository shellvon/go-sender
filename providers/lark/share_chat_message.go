package lark

// ShareChatMessage represents a share chat message for Lark/Feishu.
type ShareChatMessage struct {
	BaseMessage

	Content ShareChatContent `json:"content"`
}

// NewShareChatMessage creates a new ShareChatMessage.
// chatID is the id of the chat to share.
// See https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot#897b5321
func NewShareChatMessage(chatID string) *ShareChatMessage {
	return ShareChat().ChatID(chatID).Build()
}

func (m *ShareChatMessage) GetMsgType() MessageType {
	return TypeShareChat
}

// ShareChatContent represents the content of a share chat message.
type ShareChatContent struct {
	ChatID string `json:"share_chat_id"`
}

// ShareChatMsgBuilder provides a fluent API to construct Lark share chat messages.
type ShareChatMsgBuilder struct {
	chatID string
}

// ShareChat creates a new ShareChatMsgBuilder instance.
func ShareChat() *ShareChatMsgBuilder { return &ShareChatMsgBuilder{} }

// ChatID sets the chat ID.
func (b *ShareChatMsgBuilder) ChatID(id string) *ShareChatMsgBuilder {
	b.chatID = id
	return b
}

// Build assembles a *ShareChatMessage.
func (b *ShareChatMsgBuilder) Build() *ShareChatMessage {
	return &ShareChatMessage{
		BaseMessage: BaseMessage{MsgType: TypeShareChat},
		Content:     ShareChatContent{ChatID: b.chatID},
	}
}
