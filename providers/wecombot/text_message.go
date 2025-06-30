package wecombot

import (
	"github.com/shellvon/go-sender/core"
)

const maxTextContentLength = 2048

// Text represents the text content and mentions for a WeCom message.
type Text struct {
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

	Text Text `json:"text"`
}

// NewTextMessage creates a new TextMessage with required content and applies optional configurations.
func NewTextMessage(content string, opts ...TextMessageOption) *TextMessage {
	msg := &TextMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeText,
		},
		Text: Text{
			Content: content,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
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

	// Validate that at least one of MentionedList or MentionedMobileList is used if not empty,
	// and that "@all" is not combined with specific users/mobiles.
	if len(m.Text.MentionedList) > 0 && containsAll(m.Text.MentionedList) && len(m.Text.MentionedList) > 1 {
		return core.NewParamError("cannot combine '@all' with specific user IDs in mentioned_list")
	}
	if len(m.Text.MentionedMobileList) > 0 && containsAll(m.Text.MentionedMobileList) &&
		len(m.Text.MentionedMobileList) > 1 {
		return core.NewParamError("cannot combine '@all' with specific mobile numbers in mentioned_mobile_list")
	}

	return nil
}

// containsAll is a helper to check if a slice contains the "@all" string.
func containsAll(list []string) bool {
	for _, item := range list {
		if item == "@all" {
			return true
		}
	}
	return false
}

// TextMessageOption defines a function type for configuring TextMessage.
type TextMessageOption func(*TextMessage)

// WithMentionedList sets the MentionedList for TextMessage.
func WithMentionedList(list []string) TextMessageOption {
	return func(m *TextMessage) {
		m.Text.MentionedList = list
	}
}

// WithMentionedMobileList sets the MentionedMobileList for TextMessage.
func WithMentionedMobileList(list []string) TextMessageOption {
	return func(m *TextMessage) {
		m.Text.MentionedMobileList = list
	}
}
