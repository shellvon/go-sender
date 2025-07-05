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

// BaseMessage is the base message structure
// MsgId 可选，允许用户自定义消息ID（用于幂等、追踪等场景）.
type BaseMessage struct {
	core.DefaultMessage

	MsgType MessageType `json:"msg_type"`
	MsgID   string      `json:"msg_id,omitempty"`
}

// GetMsgType implements the Message interface.
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// ProviderType returns the provider type.
func (m *BaseMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}
