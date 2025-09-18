package emailapi

import "time"

// MailgunBuilder provides a chainable interface for building Mailgun messages.
type MailgunBuilder struct {
	BaseBuilder[*MailgunBuilder]
}

// NewMailgunBuilder creates a new MailgunBuilder.
func NewMailgunBuilder() *MailgunBuilder {
	builder := &MailgunBuilder{}
	builder.self = builder
	return builder
}

// Tags sets the tags for message tracking (Mailgun specific).
func (b *MailgunBuilder) Tags(tags []string) *MailgunBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[mailgunTagsKey] = tags
	return b
}

// DeliveryTime sets the scheduled delivery time (Mailgun specific).
// Time should be in RFC-2822 format.
func (b *MailgunBuilder) DeliveryTime(deliveryTime time.Time) *MailgunBuilder {
	b.scheduledAt = &deliveryTime
	return b
}

// DKIM enables or disables DKIM signatures (Mailgun specific).
func (b *MailgunBuilder) DKIM(enabled bool) *MailgunBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	if enabled {
		b.extras["dkim"] = "yes"
	} else {
		b.extras["dkim"] = "no"
	}
	return b
}

// Tracking enables or disables tracking (Mailgun specific).
func (b *MailgunBuilder) Tracking(enabled bool) *MailgunBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	if enabled {
		b.extras["tracking"] = "yes"
	} else {
		b.extras["tracking"] = "no"
	}
	return b
}

// Variables sets custom variables for tracking and webhooks (Mailgun specific).
// Note: For template variables, use TemplateParams() instead.
// These v: variables are included in webhook events and can be used for tracking.
func (b *MailgunBuilder) Variables(variables map[string]interface{}) *MailgunBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras["variables"] = variables
	return b
}

// Domain overrides the Mailgun domain for this specific message.
// This allows using a different domain than the one configured in the account.
// If not set, the domain from account.Region will be used.
func (b *MailgunBuilder) Domain(domain string) *MailgunBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras["domain"] = domain
	return b
}

// Build constructs a *Message for Mailgun.
func (b *MailgunBuilder) Build() *Message {
	return b.BuildMessage(string(SubProviderMailgun))
}
