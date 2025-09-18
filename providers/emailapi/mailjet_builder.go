package emailapi

import "strconv"

// MailjetBuilder provides a chainable interface for building Mailjet messages.
type MailjetBuilder struct {
	BaseBuilder[*MailjetBuilder]
}

// NewMailjetBuilder creates a new MailjetBuilder.
func NewMailjetBuilder() *MailjetBuilder {
	builder := &MailjetBuilder{}
	builder.self = builder
	return builder
}

// TemplateIDInt sets the template ID as integer (Mailjet specific).
func (b *MailjetBuilder) TemplateIDInt(templateID int) *MailjetBuilder {
	b.templateID = strconv.Itoa(templateID)
	return b
}

// Variables sets template variables (Mailjet specific).
func (b *MailjetBuilder) Variables(variables map[string]interface{}) *MailjetBuilder {
	b.templateData = variables
	return b
}

// CustomID sets a custom ID for tracking (Mailjet specific).
func (b *MailjetBuilder) CustomID(customID string) *MailjetBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras["custom_id"] = customID
	return b
}

// Campaign sets the custom campaign name (Mailjet specific).
func (b *MailjetBuilder) Campaign(campaign string) *MailjetBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras["campaign"] = campaign
	return b
}

// URLTags sets URL tags for tracking (Mailjet specific).
func (b *MailjetBuilder) URLTags(urlTags string) *MailjetBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras["url_tags"] = urlTags
	return b
}

// SandboxMode enables or disables sandbox mode (Mailjet specific).
func (b *MailjetBuilder) SandboxMode(enabled bool) *MailjetBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras["sandbox"] = enabled
	return b
}

// TextPart sets the text content (alias for Text).
func (b *MailjetBuilder) TextPart(text string) *MailjetBuilder {
	b.text = text
	return b
}

// HTMLPart sets the HTML content (alias for HTML).
func (b *MailjetBuilder) HTMLPart(html string) *MailjetBuilder {
	b.html = html
	return b
}

// Build constructs a *Message for Mailjet.
func (b *MailjetBuilder) Build() *Message {
	return b.BuildMessage(string(SubProviderMailjet))
}
