package wecombot

import (
	"github.com/shellvon/go-sender/core"
)

const maxTextContentLength = 2048

// TextContent represents the text content and mentions for a WeCom message.
type TextContent struct {
	// Content of the text message. Maximum length is 2048 bytes and must be UTF-8 encoded.
	Content string `json:"content"`
	// List of user IDs to mention in the group chat (@member). Use "@all" to mention everyone.
	// If user IDs are not available, use MentionedMobileList.
	MentionedList []string `json:"mentioned_list"`
	// List of mobile numbers to mention in the group chat (@member). Use "@all" to mention everyone.
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

// TextMessage represents a text message for WeCom.
// For more details, refer to the WeCom API documentation:
// https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E6%9C%AC%E7%B1%BB%E5%9E%8B
type TextMessage struct {
	BaseMessage

	Text TextContent `json:"text"`
}

// NewTextMessage creates a new TextMessage.
func NewTextMessage(content string) *TextMessage {
	return Text().Content(content).Build()
}

// Validate validates the TextMessage to ensure it meets WeCom API requirements.
func (m *TextMessage) Validate() error {
	if m.Text.Content == "" {
		return core.NewParamError("text content cannot be empty")
	}
	// WeCom text message content has a maximum length of 2048 bytes and must be UTF-8 encoded.
	if len(m.Text.Content) > maxTextContentLength {
		return core.NewParamError("text content exceeds 2048 bytes")
	}
	return nil
}
