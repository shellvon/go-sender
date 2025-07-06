package emailapi

const resendTagsKey = "tags"

// ResendTag represents a tag for Resend email messages.
type ResendTag struct {
	Name  string
	Value string
}

// ResendMessageBuilder provides a builder for Resend email messages.
type ResendMessageBuilder struct {
	BaseBuilder[*ResendMessageBuilder]

	tags []ResendTag
}

// Resend returns a new ResendMessageBuilder.
func Resend() *ResendMessageBuilder {
	b := &ResendMessageBuilder{}
	b.self = b
	return b
}

// Tag adds a single tag to the Resend message.
func (b *ResendMessageBuilder) Tag(name, value string) *ResendMessageBuilder {
	b.tags = append(b.tags, ResendTag{Name: name, Value: value})
	return b
}

// Tags sets all tags for the Resend message.
func (b *ResendMessageBuilder) Tags(tags []ResendTag) *ResendMessageBuilder {
	b.tags = tags
	return b
}

// Build returns the constructed *Message for Resend, with tags in Extras["tags"] as []ResendTag.
func (b *ResendMessageBuilder) Build() *Message {
	msg := b.BuildMessage(string(SubProviderResend))
	extra := map[string]interface{}{}
	if b.tags != nil {
		extra[resendTagsKey] = b.tags
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}
