package wecomapp

import (
	"github.com/shellvon/go-sender/core"
)

// 无需长度常量 - 企业微信会在需要时自动截断

// TextCardMessageContent 代表企业微信应用API的文本卡片内容
type TextCardMessageContent struct {
	// Title 文本卡片的标题。最大长度128字节
	Title string `json:"title"`
	// Description 文本卡片的描述。最大长度512字节
	Description string `json:"description"`
	// URL 点击卡片时跳转的URL。最大长度2048字节
	URL string `json:"url"`
	// BtnTxt 按钮文本。如果未指定，默认为"详情"
	BtnTxt string `json:"btntxt,omitempty"`
}

// TextCardMessage 代表企业微信应用的文本卡片消息
// 更多详情请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/90372#%E6%96%87%E6%9C%AC%E5%8D%A1%E7%89%87%E6%B6%88%E6%81%AF
type TextCardMessage struct {
	BaseMessage

	TextCard TextCardMessageContent `json:"textcard"`
}

// NewTextCardMessage 创建新的TextCardMessage
func NewTextCardMessage(title, description, url string) *TextCardMessage {
	return &TextCardMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeTextCard,
		},
		TextCard: TextCardMessageContent{
			Title:       title,
			Description: description,
			URL:         url,
		},
	}
}

// Validate 验证TextCardMessage以确保满足企业微信API要求
func (m *TextCardMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	if m.TextCard.Title == "" {
		return core.NewParamError("textcard title cannot be empty")
	}

	if m.TextCard.Description == "" {
		return core.NewParamError("textcard description cannot be empty")
	}

	if m.TextCard.URL == "" {
		return core.NewParamError("textcard URL cannot be empty")
	}

	// No length validation needed - WeChat Work will auto-truncate if needed

	return nil
}
