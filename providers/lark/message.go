package lark

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// Lark/Feishu group robot message
// Reference: https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot#5a997364

// MessageType defines the message types supported by Lark.
type MessageType string

const (
	// TypeText represents text message type.
	TypeText MessageType = providers.MsgTypeText
	// TypePost represents post message type (rich text).
	TypePost MessageType = "post"
	// TypeShareChat represents share chat message type.
	TypeShareChat MessageType = "share_chat"
	// TypeImage represents image message type.
	TypeImage MessageType = providers.MsgTypeImage
	// TypeInteractive represents interactive message type (card).
	TypeInteractive MessageType = "interactive"
)

// Message interface definition.
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// BaseMessage is the base message structure.
type BaseMessage struct {
	*core.BaseMessage

	MsgType MessageType `json:"msg_type"`
}

// Compile-time assertion: BaseMessage implements Message interface.
var (
	_ core.Message = (*BaseMessage)(nil)
)

// newBaseMessage Creates a new BaseMessage instance and sets the provider type to Lark.
func newBaseMessage(msgType MessageType) BaseMessage {
	return BaseMessage{
		BaseMessage: core.NewBaseMessage(core.ProviderTypeLark),
		MsgType:     msgType,
	}
}

// GetMsgType Implements the Message interface.
// Returns the message type.
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}
