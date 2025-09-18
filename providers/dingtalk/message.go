package dingtalk

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// DingTalk group robot message
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access

// MessageType defines the message types supported by DingTalk.
//   - Supported message types: https://open.dingtalk.com/document/orgapp/custom-bot-send-message-type
type MessageType string

const (
	// TypeText represents text message type.
	TypeText MessageType = providers.MsgTypeText
	// TypeMarkdown represents markdown message type.
	TypeMarkdown MessageType = providers.MsgTypeMarkdown
	// TypeLink represents link message type.
	TypeLink MessageType = providers.MsgTypeLink
	// TypeActionCard represents action card message type.
	TypeActionCard MessageType = providers.MsgTypeActionCard
	// TypeFeedCard represents feed card message type.
	TypeFeedCard MessageType = providers.MsgTypeFeedCard
)

// Message interface definition.
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// BaseMessage is the base message structure.
type BaseMessage struct {
	*core.BaseMessage
	MsgType MessageType `json:"msgtype"`
}

func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// Compile-time assertion: BaseMessage implements Message interface.
var (
	_ core.Message = (*BaseMessage)(nil)
)

// newBaseMessage Creates a new BaseMessage instance and sets the provider type to DingTalk.
func newBaseMessage(msgType MessageType) BaseMessage {
	return BaseMessage{
		BaseMessage: core.NewBaseMessage(core.ProviderTypeDingtalk),
		MsgType:     msgType,
	}
}
