package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// LinkContent represents the link content for a DingTalk message.
type LinkContent struct {
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

	Link LinkContent `json:"link"`
}

// NewLinkMessage creates a new LinkMessage with required content.
func NewLinkMessage(title, text, messageURL string) *LinkMessage {
	return Link().Title(title).Text(text).MessageURL(messageURL).Build()
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
