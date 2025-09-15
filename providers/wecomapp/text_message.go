package wecomapp

import (
	"github.com/shellvon/go-sender/core"
)

// 无需长度常量 - 企业微信会在需要时自动截断

// TextMessageContent 代表企业微信应用API的文本内容.
type TextMessageContent struct {
	// Content 文本消息的内容
	Content string `json:"content"`
}

// TextMessage 代表企业微信应用的文本消息
// 更多详情请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/90372#%E6%96%87%E6%9C%AC%E6%B6%88%E6%81%AF
type TextMessage struct {
	BaseMessage

	Text TextMessageContent `json:"text"`
}

// NewTextMessage 创建新的TextMessage.
func NewTextMessage(content string) *TextMessage {
	return &TextMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeText,
		},
		Text: TextMessageContent{
			Content: content,
		},
	}
}

// Validate 验证TextMessage以确保满足企业微信API要求.
func (m *TextMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	if m.Text.Content == "" {
		return core.NewParamError("text content cannot be empty")
	}

	// 无需长度验证 - 企业微信会在需要时自动截断

	return nil
}
