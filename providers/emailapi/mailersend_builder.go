package emailapi

import "time"

// MailerSendBuilder provides a chainable interface for building MailerSend messages.
type MailerSendBuilder struct {
	BaseBuilder[*MailerSendBuilder]
}

// NewMailerSendBuilder creates a new MailerSendBuilder.
func NewMailerSendBuilder() *MailerSendBuilder {
	builder := &MailerSendBuilder{}
	builder.self = builder
	return builder
}

// Tags sets the tags for message tracking (MailerSend specific).
func (b *MailerSendBuilder) Tags(tags []string) *MailerSendBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[mailerSendTagsKey] = tags
	return b
}

// Personalize adds personalization data for a specific recipient email.
// This allows per-recipient customization of template variables.
// The data is stored in TemplateData with the email as the key.
//
// Example:
//
//	builder.Personalize("john@example.com", map[string]interface{}{
//	    "name": "John",
//	    "company": "ACME Corp",
//	})
//
// See also: GetEmailType, EmailTypeTemplate, TemplateData
func (b *MailerSendBuilder) Personalize(email string, data map[string]interface{}) *MailerSendBuilder {
	if b.templateData == nil {
		b.templateData = make(map[string]interface{})
	}

	// Store personalization data directly in TemplateData
	// Key = email address, Value = personalization data
	b.templateData[email] = data

	return b
}

// Settings configures email tracking settings.
//
// Example:
//
//	builder.Settings(map[string]bool{
//	    "track_clicks": true,
//	    "track_opens": true,
//	    "track_content": false,
//	})
func (b *MailerSendBuilder) Settings(settings map[string]bool) *MailerSendBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[mailerSendSettingsKey] = settings
	return b
}

// InReplyTo sets the In-Reply-To header for threading emails.
// Valid email address as per RFC 2821.
func (b *MailerSendBuilder) InReplyTo(inReplyTo string) *MailerSendBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[mailerSendInReplyToKey] = inReplyTo
	return b
}

// ListUnsubscribe sets a single value that complies with RFC 8058.
// This header provides recipients with a way to unsubscribe from mailing lists.
//
// Please note that this feature is available to Professional and Enterprise accounts only.
//
// Example:
//
//	builder.ListUnsubscribe("<mailto:unsubscribe@example.com>")
//	builder.ListUnsubscribe("<https://example.com/unsubscribe>")
func (b *MailerSendBuilder) ListUnsubscribe(unsubscribeInfo string) *MailerSendBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[mailerSendListUnsubscribeKey] = unsubscribeInfo
	return b
}

// SendAt sets the scheduled send time (MailerSend specific).
func (b *MailerSendBuilder) SendAt(sendAt time.Time) *MailerSendBuilder {
	b.scheduledAt = &sendAt
	return b
}

// Build constructs a *Message for MailerSend.
func (b *MailerSendBuilder) Build() *Message {
	return b.BuildMessage(string(SubProviderMailerSend))
}
