package dingtalk

import "github.com/shellvon/go-sender/core"

// DingTalk group robot message
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access

// MessageType defines the message types supported by DingTalk.
type MessageType string

const (
	// TypeText represents text message type.
	TypeText = "text"
	// TypeMarkdown represents markdown message type.
	TypeMarkdown = "markdown"
	// TypeLink represents link message type.
	TypeLink = "link"
	// TypeActionCard represents action card message type.
	TypeActionCard = "actionCard"
	// TypeFeedCard represents feed card message type.
	TypeFeedCard = "feedCard"
)

// Message interface definition.
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// BaseMessage is the base message structure.
type BaseMessage struct {
	core.DefaultMessage

	MsgType MessageType `json:"msgtype"`
}

// GetMsgType implements the Message interface.
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
