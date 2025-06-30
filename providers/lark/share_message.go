package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// ShareChatMessage represents a share chat message for Lark/Feishu.
type ShareChatMessage struct {
	BaseMessage

	Content ShareChatContent `json:"content"`
}

// ShareChatContent represents the content of a share chat message.
type ShareChatContent struct {
	ShareChat ShareChat `json:"share_chat"`
}

// ShareChat represents the share chat structure.
type ShareChat struct {
	ChatID string `json:"chat_id"`
}

// NewShareChatMessage creates a new share chat message.
func NewShareChatMessage(chatID string) *ShareChatMessage {
	return &ShareChatMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeShareChat,
		},
		Content: ShareChatContent{
			ShareChat: ShareChat{
				ChatID: chatID,
			},
		},
	}
}

// GetMsgType returns the message type.
func (m *ShareChatMessage) GetMsgType() MessageType {
	return TypeShareChat
}

// ProviderType returns the provider type.
func (m *ShareChatMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the share chat message.
func (m *ShareChatMessage) Validate() error {
	if m.Content.ShareChat.ChatID == "" {
		return errors.New("chat_id cannot be empty")
	}
	return nil
}

// ShareUserMessage represents a share user message for Lark/Feishu.
type ShareUserMessage struct {
	BaseMessage

	Content ShareUserContent `json:"content"`
}

// ShareUserContent represents the content of a share user message.
type ShareUserContent struct {
	ShareUser ShareUser `json:"share_user"`
}

// ShareUser represents the share user structure.
type ShareUser struct {
	UserID string `json:"user_id"`
}

// NewShareUserMessage creates a new share user message.
func NewShareUserMessage(userID string) *ShareUserMessage {
	return &ShareUserMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeShareUser,
		},
		Content: ShareUserContent{
			ShareUser: ShareUser{
				UserID: userID,
			},
		},
	}
}

// GetMsgType returns the message type.
func (m *ShareUserMessage) GetMsgType() MessageType {
	return TypeShareUser
}

// ProviderType returns the provider type.
func (m *ShareUserMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the share user message.
func (m *ShareUserMessage) Validate() error {
	if m.Content.ShareUser.UserID == "" {
		return errors.New("user_id cannot be empty")
	}
	return nil
}
