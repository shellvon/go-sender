package lark

import "github.com/shellvon/go-sender/core"

// Lark/Feishu group robot message
// Reference: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN

// MessageType defines the message types supported by Lark.
type MessageType string

const (
	// TypeText represents text message type.
	TypeText MessageType = "text"
	// TypePost represents post message type (rich text).
	TypePost MessageType = "post"
	// TypeShareChat represents share chat message type.
	TypeShareChat MessageType = "share_chat"
	// TypeShareUser represents share user message type.
	TypeShareUser MessageType = "share_user"
	// TypeImage represents image message type.
	TypeImage MessageType = "image"
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

// WithLarkMsgID 设置 Lark 消息的 msg_id（可选，用于幂等、追踪等场景）
//
// 示例：WithLarkMsgID("custom-id-123").
func WithLarkMsgID(id string) func(m *BaseMessage) {
	return func(m *BaseMessage) {
		m.MsgID = id
	}
}
