package lark

import "github.com/shellvon/go-sender/core"

// Lark/Feishu group robot message
// Reference: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN

// MessageType defines the message types supported by Lark
type MessageType string

const (
	// Text message type
	TypeText MessageType = "text"
	// Post message type (rich text)
	TypePost MessageType = "post"
	// Share chat message type
	TypeShareChat MessageType = "share_chat"
	// Share user message type
	TypeShareUser MessageType = "share_user"
	// Image message type
	TypeImage MessageType = "image"
	// Interactive message type (card)
	TypeInteractive MessageType = "interactive"
)

// Message interface definition
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// BaseMessage is the base message structure
// MsgId 可选，允许用户自定义消息ID（用于幂等、追踪等场景）
type BaseMessage struct {
	core.DefaultMessage
	MsgType MessageType `json:"msg_type"`
	MsgId   string      `json:"msg_id,omitempty"`
}

// GetMsgType implements the Message interface
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// WithLarkMsgId 设置 Lark 消息的 msg_id（可选，用于幂等、追踪等场景）
//
// 示例：WithLarkMsgId("custom-id-123")
func WithLarkMsgId(id string) func(m *BaseMessage) {
	return func(m *BaseMessage) {
		m.MsgId = id
	}
}
