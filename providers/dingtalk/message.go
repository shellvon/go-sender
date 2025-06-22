package dingtalk

import "github.com/shellvon/go-sender/core"

// DingTalk group robot message
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access

// MessageType defines the message types supported by DingTalk
type MessageType string

const (
	// Text message type
	TypeText MessageType = "text"
	// Markdown message type
	TypeMarkdown MessageType = "markdown"
	// Link message type
	TypeLink MessageType = "link"
	// ActionCard message type
	TypeActionCard MessageType = "actionCard"
	// FeedCard message type
	TypeFeedCard MessageType = "feedCard"
)

// Message interface definition
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// BaseMessage is the base message structure
type BaseMessage struct {
	core.DefaultMessage
	MsgType MessageType `json:"msgtype"`
}

// GetMsgType implements the Message interface
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

func (m *TextMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeDingtalk
}

func (m *MarkdownMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeDingtalk
}

func (m *LinkMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeDingtalk
}

func (m *ActionCardMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeDingtalk
}

func (m *FeedCardMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeDingtalk
}
