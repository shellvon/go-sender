package wecombot

import (
	"github.com/shellvon/go-sender/core"
)

const maxTextContentLength = 2048

// TextContent 表示企业微信消息的文本内容和@提及信息。
type TextContent struct {
	// 文本消息的内容。最大长度为 2048 字节，且必须是 UTF-8 编码。
	Content string `json:"content"`
	// 群聊中@提及的用户 ID 列表（@member）。使用 "@all" 可提及所有人。
	// 如果用户 ID 不可用，请使用 MentionedMobileList。
	MentionedList []string `json:"mentioned_list"`
	// 群聊中@提及的手机号码列表（@member）。使用 "@all" 可提及所有人。
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

// TextMessage 表示企业微信的文本消息。
// 更多详情，请参考企业微信 API 文档：
// https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E6%9C%AC%E7%B1%BB%E5%9E%8B
type TextMessage struct {
	BaseMessage

	Text TextContent `json:"text"`
}

// NewTextMessage 创建一个新的 TextMessage 实例。
// 参数：content string - 文本消息内容。
// 返回值：*TextMessage - 新创建的文本消息实例。
func NewTextMessage(content string) *TextMessage {
	return Text().Content(content).Build()
}

// Validate 验证 TextMessage 是否满足企业微信 API 的要求。
// 该方法检查文本内容是否为空以及是否超过 2048 字节的限制。
// 返回值：error - 如果验证失败，返回具体的参数错误；否则返回 nil。
func (m *TextMessage) Validate() error {
	if m.Text.Content == "" {
		return core.NewParamError("文本内容不能为空")
	}
	// 企业微信文本消息内容最大长度为 2048 字节，且必须是 UTF-8 编码。
	if len(m.Text.Content) > maxTextContentLength {
		return core.NewParamError("文本内容超过 2048 字节")
	}
	return nil
}
