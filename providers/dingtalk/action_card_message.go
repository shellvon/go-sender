package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// ActionCard represents the action card content for a DingTalk message.
type ActionCard struct {
	// Title of the action card
	Title string `json:"title"`
	// Content of the action card
	Text string `json:"text"`
	// Button text
	BtnOrientation string `json:"btnOrientation,omitempty"`
	// Single button (for single action card)
	SingleTitle string `json:"singleTitle,omitempty"`
	SingleURL   string `json:"singleURL,omitempty"`
	// Multiple buttons (for multiple action card)
	Btns []ActionCardButton `json:"btns,omitempty"`
}

// ActionCardButton represents a button in action card.
type ActionCardButton struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

// ActionCardMessage represents an action card message for DingTalk.
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access
type ActionCardMessage struct {
	BaseMessage
	ActionCard ActionCard `json:"actionCard"`
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

	return nil
}

// ActionCardMessageOption defines a function type for configuring ActionCardMessage.
type ActionCardMessageOption func(*ActionCardMessage)

// WithBtnOrientation sets the BtnOrientation for ActionCardMessage.
func WithBtnOrientation(orientation string) ActionCardMessageOption {
	return func(m *ActionCardMessage) {
		m.ActionCard.BtnOrientation = orientation
	}
}

// WithSingleButton sets the single button for ActionCardMessage.
func WithSingleButton(title, url string) ActionCardMessageOption {
	return func(m *ActionCardMessage) {
		m.ActionCard.SingleTitle = title
		m.ActionCard.SingleURL = url
	}
}

// WithMultipleButtons sets the multiple buttons for ActionCardMessage.
func WithMultipleButtons(btns []ActionCardButton) ActionCardMessageOption {
	return func(m *ActionCardMessage) {
		m.ActionCard.Btns = btns
	}
}

// NewActionCardMessage creates a new ActionCardMessage with required content and applies optional configurations.
func NewActionCardMessage(title, text string, opts ...ActionCardMessageOption) *ActionCardMessage {
	msg := &ActionCardMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeActionCard,
		},
		ActionCard: ActionCard{
			Title: title,
			Text:  text,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}
