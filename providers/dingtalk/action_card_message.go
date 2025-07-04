package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// ActionCardContent represents the action card content for a DingTalk message.
type ActionCardContent struct {
	// Title of the action card
	Title string `json:"title"`
	// Content of the action card, 如果需要实现 @ 功能 ，在 text 内容中添加 @ 用户的 userId。比如 @manager7675
	Text string `json:"text"`
	// 按钮排列方式，0：竖向排列，1：横向排列
	//   - 0: vertical (default)
	//   - 1: horizontal
	BtnOrientation string `json:"btnOrientation,omitempty"`

	// Single button (for single action card
	SingleTitle string `json:"singleTitle,omitempty"`
	SingleURL   string `json:"singleURL,omitempty"`

	// Multiple buttons (for multiple action card)
	Btns []ActionCardButton `json:"btns,omitempty"`
}

// ActionCardButton represents a button in action card.
type ActionCardButton struct {
	// 按钮标题
	Title string `json:"title"`
	// 按钮点击链接
	ActionURL string `json:"actionURL"`
}

// ActionCardMessage represents an action card message for DingTalk.
// Reference:
//   - https://open.dingtalk.com/document/robots/custom-robot-access
//   - https://open.dingtalk.com/document/orgapp/custom-bot-send-message-type
type ActionCardMessage struct {
	BaseMessage

	ActionCard ActionCardContent `json:"actionCard"`
}

// NewActionCardMessage creates a new ActionCardMessage.
func NewActionCardMessage(title, text string) *ActionCardMessage {
	return ActionCard().Title(title).Text(text).Build()
}

// Validate validates the ActionCardMessage to ensure it meets DingTalk API requirements.
func (m *ActionCardMessage) Validate() error {
	if m.ActionCard.Title == "" {
		return core.NewParamError("action card title cannot be empty")
	}
	if m.ActionCard.Text == "" {
		return core.NewParamError("action card text cannot be empty")
	}

	// Check if it's single action card or multiple action card
	hasSingle := m.ActionCard.SingleTitle != "" && m.ActionCard.SingleURL != ""
	hasMultiple := len(m.ActionCard.Btns) > 0

	if !hasSingle && !hasMultiple {
		return core.NewParamError("action card must have either single button or multiple buttons")
	}

	if hasSingle && hasMultiple {
		return core.NewParamError("action card cannot have both single button and multiple buttons")
	}

	if m.ActionCard.BtnOrientation != "0" && m.ActionCard.BtnOrientation != "1" {
		return core.NewParamError("action card button orientation must be 0 or 1")
	}

	return nil
}
