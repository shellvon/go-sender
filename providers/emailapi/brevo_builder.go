package emailapi

import (
	"strconv"
	"time"
)

// BrevoBuilder provides a chainable interface for building Brevo messages.
type BrevoBuilder struct {
	BaseBuilder[*BrevoBuilder]
}

// NewBrevoBuilder creates a new BrevoBuilder.
func NewBrevoBuilder() *BrevoBuilder {
	builder := &BrevoBuilder{}
	builder.self = builder
	return builder
}

// Tags sets the tags for message tracking (Brevo specific).
func (b *BrevoBuilder) Tags(tags []string) *BrevoBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[brevoTagsKey] = tags
	return b
}

// Params sets template parameters (Brevo specific).
func (b *BrevoBuilder) Params(params map[string]interface{}) *BrevoBuilder {
	b.templateData = params
	return b
}

// TemplateIDInt sets the template ID as integer (Brevo specific).
func (b *BrevoBuilder) TemplateIDInt(templateID int) *BrevoBuilder {
	b.templateID = strconv.Itoa(templateID)
	return b
}

// Sender sets the sender email address (alias for From).
func (b *BrevoBuilder) Sender(sender string) *BrevoBuilder {
	b.from = sender

	return b
}

// ScheduledAt sets the scheduled send time (Brevo specific).
func (b *BrevoBuilder) ScheduledAt(scheduledAt time.Time) *BrevoBuilder {
	b.scheduledAt = &scheduledAt
	return b
}

// HTMLContent sets the HTML content (alias for HTML).
func (b *BrevoBuilder) HTMLContent(htmlContent string) *BrevoBuilder {
	b.html = htmlContent
	return b
}

// TextContent sets the text content (alias for Text).
func (b *BrevoBuilder) TextContent(textContent string) *BrevoBuilder {
	b.text = textContent
	return b
}

// Preheader sets the email preheader (preview text).
// A short summary that appears next to the subject line in the recipient’s inbox. This preview text gives recipients a quick idea of what the email is about before they open it.
// Docs: https://developers.brevo.com/reference/sendtransacemail
func (b *BrevoBuilder) Preheader(preheader string) *BrevoBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras["preheader"] = preheader
	return b
}

// BatchId sets the batch ID for grouping emails (Brevo specific).
// This allows you to group related emails together for tracking purposes.
func (b *BrevoBuilder) BatchId(batchId string) *BrevoBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[brevoBatchIdKey] = batchId
	return b
}

// MessageVersions sets message versions.
// Docs: https://developers.brevo.com/reference/sendtransacemail
//   - You can customize and send out multiple versions of a mail.
//   - templateId can be customized only if global parameter contains templateId.
//   - htmlContent and textContent can be customized only if any of the two, htmlContent or textContent, is present in global parameters.
//   - Some global parameters such as to(mandatory), bcc, cc, replyTo, subject can also be customized specific to each version.
//   - Total number of recipients in one API request must not exceed 2000. However, you can still pass upto 99 recipients maximum in one message version.
//   - The size of individual params in all the messageVersions shall not exceed 100 KB limit and that of cumulative params shall not exceed 1000 KB.
//   - You can follow this step-by-step guide on how to use messageVersions to batch send emails
//   - https://developers.brevo.com/docs/batch-send-transactional-emails
func (b *BrevoBuilder) MessageVersions(versions []map[string]interface{}) *BrevoBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[brevoMessageVersionsKey] = versions
	return b
}

// Build constructs a *Message for Brevo.
func (b *BrevoBuilder) Build() *Message {
	return b.BuildMessage(string(SubProviderBrevo))
}
