package wecomapp

import (
	"errors"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// MessageType 定义企业微信应用支持的消息类型.
type MessageType string

const (
	// TypeText 代表文本消息类型.
	TypeText MessageType = providers.MsgTypeText
	// TypeImage 代表图片消息类型.
	TypeImage MessageType = providers.MsgTypeImage
	// TypeVoice 代表语音消息类型.
	TypeVoice MessageType = providers.MsgTypeVoice
	// TypeVideo 代表视频消息类型.
	TypeVideo MessageType = providers.MsgTypeVideo
	// TypeFile 代表文件消息类型.
	TypeFile MessageType = providers.MsgTypeFile
	// TypeNews 代表图文消息类型.
	TypeNews MessageType = providers.MsgTypeNews
	// TypeMarkdown 代表markdown消息类型.
	TypeMarkdown MessageType = providers.MsgTypeMarkdown
	// TypeTemplateCard 代表模板卡片消息类型.
	TypeTemplateCard MessageType = providers.MsgTypeTemplateCard
	// TypeTextCard 代表文本卡片消息类型.
	TypeTextCard MessageType = providers.MsgTypeTextCard
	// TypeMPNews 代表mpnews消息类型.
	TypeMPNews MessageType = providers.MsgTypeMPNews
	// TypeMiniprogramNotice 代表小程序通知消息类型.
	TypeMiniprogramNotice MessageType = providers.MsgTypeMiniprogramNotice
)

// MediaType 定义企业微信应用上传的媒体类型.
type MediaType string

const (
	// MediaTypeImage 代表上传的图片媒体类型.
	MediaTypeImage MediaType = "image"
	// MediaTypeVoice 代表上传的语音媒体类型.
	MediaTypeVoice MediaType = "voice"
	// MediaTypeVideo 代表上传的视频媒体类型.
	MediaTypeVideo MediaType = "video"
	// MediaTypeFile 代表上传的文件媒体类型.
	MediaTypeFile MediaType = "file"
)

// Message 企业微信应用消息的接口定义.
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// agentIDSetter 需要自动设置AgentID的消息接口
// 这是provider使用的内部接口.
type agentIDSetter interface {
	setAgentID(agentID string)
}

// CommonFields 包含所有企业微信应用消息共享的通用字段.
type CommonFields struct {
	// ToUser 指定发送给的用户ID，用"|"分隔。使用"@all"发送给所有用户
	ToUser string `json:"touser,omitempty"`
	// ToParty 指定发送给的部门ID，用"|"分隔
	ToParty string `json:"toparty,omitempty"`
	// ToTag 指定发送给的标签ID，用"|"分隔
	ToTag string `json:"totag,omitempty"`
	// AgentID 应用ID（从账号配置自动设置）
	AgentID string `json:"agentid"`
	// Safe 是否启用安全模式（0：否，1：是）
	Safe int `json:"safe,omitempty"`
	// EnableIDTrans 是否启用ID转换（0：否，1：是）。默认为0
	EnableIDTrans int `json:"enable_id_trans,omitempty"`
	// EnableDuplicateCheck 是否启用重复消息检查（0：否，1：是）。默认为0
	EnableDuplicateCheck int `json:"enable_duplicate_check,omitempty"`
	// DuplicateCheckInterval 重复检查间隔（秒）。默认1800秒，最大4小时
	DuplicateCheckInterval int `json:"duplicate_check_interval,omitempty"`
}

// ValidateCommonFields 验证通用字段以确保满足企业微信API要求.
func (c *CommonFields) ValidateCommonFields() error {
	// 至少必须指定一个目标
	if c.ToUser == "" && c.ToParty == "" && c.ToTag == "" {
		return core.NewParamError("at least one of touser, toparty, or totag must be specified")
	}

	return nil
}

// BaseMessage 企业微信应用的基础消息结构.
type BaseMessage struct {
	core.DefaultMessage
	CommonFields

	MsgType MessageType `json:"msgtype"`
}

// GetMsgType 实现Message接口.
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// ProviderType 返回企业微信应用的provider类型.
func (m *BaseMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeWecomApp
}

// setAgentID 设置消息的agent ID（仅内部使用）.
func (m *BaseMessage) setAgentID(agentID string) {
	m.AgentID = agentID
}

// Validate 对基础消息执行基本验证.
func (m *BaseMessage) Validate() error {
	if m.MsgType == "" {
		return errors.New("message type is required")
	}
	return m.ValidateCommonFields()
}
