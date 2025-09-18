package dingtalk

// FeedCardMsgBuilder provides a fluent API to construct DingTalk feed-card messages.
//
// Example:
//
//	msg := dingtalk.FeedCard().
//	         AddLink("title", "text", "https://example.com", "https://example.com/logo.png").
//	         Build()
//
// The underlying FeedCard message can contain multiple independent links.
type FeedCardMsgBuilder struct {
	links []FeedCardLink
}

// FeedCard creates a new FeedCardMsgBuilder instance.
func FeedCard() *FeedCardMsgBuilder { return &FeedCardMsgBuilder{} }

// Links sets the whole slice at once (overwrites previous additions).
func (b *FeedCardMsgBuilder) Links(links []FeedCardLink) *FeedCardMsgBuilder {
	b.links = links
	return b
}

// AddLink appends a single link entry.
func (b *FeedCardMsgBuilder) AddLink(title, messageURL, picURL string) *FeedCardMsgBuilder {
	b.links = append(b.links, FeedCardLink{
		Title:      title,
		MessageURL: messageURL,
		PicURL:     picURL,
	})
	return b
}

// Build assembles a *FeedCardMessage.
func (b *FeedCardMsgBuilder) Build() *FeedCardMessage {
	return &FeedCardMessage{
		BaseMessage: newBaseMessage(TypeFeedCard),
		FeedCard: FeedCardContent{
			Links: b.links,
		},
	}
}
