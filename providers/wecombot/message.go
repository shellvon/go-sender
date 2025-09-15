package wecombot

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// 企业微信群机器人消息，与应用消息不同（https://developer.work.weixin.qq.com/document/path/90236）
// 群机器人基于群聊，参数更简单：
// https://developer.work.weixin.qq.com/document/path/91770

// MessageType 定义企业微信支持的消息类型。
type MessageType string

const (
	// TypeText 表示文本消息类型。
	TypeText MessageType = providers.MsgTypeText
	// TypeMarkdown 表示 Markdown 消息类型。
	TypeMarkdown MessageType = providers.MsgTypeMarkdown
	// TypeImage 表示图片消息类型。
	TypeImage MessageType = providers.MsgTypeImage
	// TypeNews 表示新闻消息类型。
	TypeNews MessageType = providers.MsgTypeNews
	// TypeTemplateCard 表示模板卡片消息类型。
	TypeTemplateCard MessageType = providers.MsgTypeTemplateCard
	// TypeVoice 表示语音消息类型。
	TypeVoice MessageType = providers.MsgTypeVoice
	// TypeFile 表示文件消息类型。
	TypeFile MessageType = providers.MsgTypeFile
)

// Message 接口定义。
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// BaseMessage 是消息的基础结构。
type BaseMessage struct {
	core.DefaultMessage

	MsgType MessageType `json:"msgtype"`
}

// GetMsgType 实现 Message 接口，返回消息类型。
// 返回值：MessageType - 当前消息的类型。
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// ProviderType 返回提供者类型，表示企业微信群机器人。
// 返回值：core.ProviderType - 提供者类型，固定为 ProviderTypeWecombot。
func (m *BaseMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}
