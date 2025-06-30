package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// FeedCard represents the feed card content for a DingTalk message.
type FeedCard struct {
	// List of feed card links
	Links []FeedCardLink `json:"links"`
}

// FeedCardLink represents a link in feed card.
type FeedCardLink struct {
	// Title of the link
	Title string `json:"title"`
	// Content of the link
	Text string `json:"text"`
	// Link URL
	MessageURL string `json:"messageURL"`
	// Picture URL
	PicURL string `json:"picURL"`
}

// FeedCardMessage represents a feed card message for DingTalk.
// Reference: https://open.dingtalk.com/document/robots/custom-robot-access
type FeedCardMessage struct {
	BaseMessage

	FeedCard FeedCard `json:"feedCard"`
}

// NewFeedCardMessage creates a new FeedCardMessage with required content and applies optional configurations.
func NewFeedCardMessage(links []FeedCardLink, opts ...FeedCardMessageOption) *FeedCardMessage {
	msg := &FeedCardMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeFeedCard,
		},
		FeedCard: FeedCard{
			Links: links,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// Validate validates the FeedCardMessage to ensure it meets DingTalk API requirements.
func (m *FeedCardMessage) Validate() error {
	if len(m.FeedCard.Links) == 0 {
		return core.NewParamError("feed card must have at least one link")
	}

	for _, link := range m.FeedCard.Links {
		if link.Title == "" {
			return core.NewParamError("feed card link title cannot be empty")
		}
		if link.Text == "" {
			return core.NewParamError("feed card link text cannot be empty")
		}
		if link.MessageURL == "" {
			return core.NewParamError("feed card link messageURL cannot be empty")
		}
		if link.PicURL == "" {
			return core.NewParamError("feed card link picURL cannot be empty")
		}
	}

	return nil
}

// FeedCardMessageOption defines a function type for configuring FeedCardMessage.
type FeedCardMessageOption func(*FeedCardMessage)

// WithFeedCardLinks sets the Links for FeedCardMessage.
func WithFeedCardLinks(links []FeedCardLink) FeedCardMessageOption {
	return func(m *FeedCardMessage) {
		m.FeedCard.Links = links
	}
}
