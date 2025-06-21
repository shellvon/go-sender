package wecombot

import "github.com/shellvon/go-sender/core"

// WeCom group robot message, distinct from app messages (https://developer.work.weixin.qq.com/document/path/90236)
// The robot is based on group chats with simpler and more simple parameters:
// https://developer.work.weixin.qq.com/document/path/91770

// MessageType defines the message types supported by WeCom
type MessageType string

const (
	// Text message type
	TypeText MessageType = "text"
	// Markdown message type
	TypeMarkdown MessageType = "markdown"
	// Image message type
	TypeImage MessageType = "image"
	// News message type
	TypeNews MessageType = "news"
	// Template card message type
	TypeTemplateCard MessageType = "template_card"
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
	return core.ProviderTypeWecombot
}

func (m *MarkdownMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}

func (m *ImageMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}

func (m *NewsMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}

func (m *TemplateCardMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}
