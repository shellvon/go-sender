package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// Link represents the link content for a DingTalk message.
type Link struct {
	// Title of the link message
	Title string `json:"title"`
	// Content of the link message
	Text string `json:"text"`
	// Link URL
	MessageURL string `json:"messageUrl"`
	// Picture URL (optional)
	PicURL string `json:"picUrl,omitempty"`
}

// LinkMessage represents a link message for DingTalk.
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access
type LinkMessage struct {
	BaseMessage
	Link Link `json:"link"`
}

// Validate validates the LinkMessage to ensure it meets DingTalk API requirements.
func (m *LinkMessage) Validate() error {
	if m.Link.Title == "" {
		return core.NewParamError("link title cannot be empty")
	}
	if m.Link.Text == "" {
		return core.NewParamError("link text cannot be empty")
	}
	if m.Link.MessageURL == "" {
		return core.NewParamError("link messageUrl cannot be empty")
	}

	return nil
}

// LinkMessageOption defines a function type for configuring LinkMessage.
type LinkMessageOption func(*LinkMessage)

// WithLinkPicURL sets the PicURL for LinkMessage.
func WithLinkPicURL(picURL string) LinkMessageOption {
	return func(m *LinkMessage) {
		m.Link.PicURL = picURL
	}
}

// NewLinkMessage creates a new LinkMessage with required content and applies optional configurations.
func NewLinkMessage(title, text, messageURL string, opts ...LinkMessageOption) *LinkMessage {
	msg := &LinkMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeLink,
		},
		Link: Link{
			Title:      title,
			Text:       text,
			MessageURL: messageURL,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}
