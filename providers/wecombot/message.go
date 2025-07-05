package wecombot

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// WeCom group robot message, distinct from app messages (https://developer.work.weixin.qq.com/document/path/90236)
// The robot is based on group chats with simpler and more simple parameters:
// https://developer.work.weixin.qq.com/document/path/91770

// MessageType defines the message types supported by WeCom.
type MessageType string

const (
	// TypeText represents text message type.
	TypeText MessageType = providers.MsgTypeText
	// TypeMarkdown represents markdown message type.
	TypeMarkdown MessageType = providers.MsgTypeMarkdown
	// TypeImage represents image message type.
	TypeImage MessageType = providers.MsgTypeImage
	// TypeNews represents news message type.
	TypeNews MessageType = providers.MsgTypeNews
	// TypeTemplateCard represents template card message type.
	TypeTemplateCard MessageType = providers.MsgTypeTemplateCard
	// TypeVoice represents voice message type.
	TypeVoice MessageType = providers.MsgTypeVoice
	// TypeFile represents file message type.
	TypeFile MessageType = providers.MsgTypeFile
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

// ProviderType implementations for new message types.
func (m *VoiceMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}

func (m *FileMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}
