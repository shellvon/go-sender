package wecomapp

import (
	"github.com/shellvon/go-sender/core"
)

// 无需长度常量 - 企业微信会在需要时自动截断

// MarkdownMessageContent 代表企业微信应用API的markdown内容.
type MarkdownMessageContent struct {
	// Content markdown消息的内容
	Content string `json:"content"`
}

// MarkdownMessage 代表企业微信应用的markdown消息
// 更多详情请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/90372#markdown%E6%B6%88%E6%81%AF
type MarkdownMessage struct {
	BaseMessage

	Markdown MarkdownMessageContent `json:"markdown"`
}

// NewMarkdownMessage 创建新的MarkdownMessage.
func NewMarkdownMessage(content string) *MarkdownMessage {
	return &MarkdownMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeMarkdown,
		},
		Markdown: MarkdownMessageContent{
			Content: content,
		},
	}
}

// Validate 验证MarkdownMessage以确保满足企业微信API要求.
func (m *MarkdownMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	if m.Markdown.Content == "" {
		return core.NewParamError("markdown content cannot be empty")
	}

	// 无需长度验证 - 企业微信会在需要时自动截断

	if m.AgentID == "" {
		return core.NewParamError("agentid is required")
	}

	// 必须指定至少一个目标
	if m.ToUser == "" && m.ToParty == "" && m.ToTag == "" {
		return core.NewParamError("at least one of touser, toparty, or totag must be specified")
	}

	return nil
}
