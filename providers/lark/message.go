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
type BaseMessage struct {
	core.DefaultMessage
	MsgType MessageType `json:"msg_type"`
}

// GetMsgType implements the Message interface
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}
