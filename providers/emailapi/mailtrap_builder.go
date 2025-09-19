package emailapi

// MailtrapBuilder provides a chainable interface for building Mailtrap messages.
// It follows the Mailtrap API specification and supports all available features
// including templates, custom variables, category tracking, and attachments.
//
// For template emails, use [TemplateUUID] and [TemplateVariables].
// For regular emails, use [Subject], [Text]/[HTML], and other content methods.
type MailtrapBuilder struct {
	BaseBuilder[*MailtrapBuilder]
}

// NewMailtrapBuilder creates a new MailtrapBuilder.
func NewMailtrapBuilder() *MailtrapBuilder {
	builder := &MailtrapBuilder{}
	builder.self = builder
	return builder
}

// Category sets the category for message tracking (Mailtrap specific).
//
// See also: [CustomVariables].
func (b *MailtrapBuilder) Category(category string) *MailtrapBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[mailtrapCategoryKey] = category
	return b
}

// CustomVariables sets custom variables for tracking (Mailtrap specific).
// This follows the OpenAPI specification which uses 'custom_variables' field.
//
// See also: [Category].
func (b *MailtrapBuilder) CustomVariables(variables map[string]interface{}) *MailtrapBuilder {
	if b.extras == nil {
		b.extras = make(map[string]interface{})
	}
	b.extras[mailtrapCustomVariablesKey] = variables
	return b
}

// CustomArgs is deprecated, use [CustomVariables] instead.
// Kept for backward compatibility.
//
// Deprecated: Use [CustomVariables] instead.
func (b *MailtrapBuilder) CustomArgs(args map[string]interface{}) *MailtrapBuilder {
	return b.CustomVariables(args)
}

// TemplateUUID sets the template UUID for template-based emails.
//
// See also: [BuildTemplateEmail], [EmailTypeTemplate].
func (b *MailtrapBuilder) TemplateUUID(templateUUID string) *MailtrapBuilder {
	b.templateID = templateUUID
	return b
}

// TemplateVariables sets the template variables for template-based emails.
//
// See also: [BuildTemplateEmail], [EmailTypeTemplate], [TemplateUUID].
func (b *MailtrapBuilder) TemplateVariables(variables map[string]interface{}) *MailtrapBuilder {
	b.templateData = variables
	return b
}

// Build constructs a *Message for Mailtrap.
func (b *MailtrapBuilder) Build() *Message {
	return b.BuildMessage(string(SubProviderMailtrap))
}

// Builder convenience methods for different email types

// BuildTextEmail creates a text-only email message.
// The resulting message will have EmailType EmailTypeText.
//
// See also: [EmailTypeText], [GetEmailType].
func (b *MailtrapBuilder) BuildTextEmail() *Message {
	msg := b.Build()
	// Ensure it's a valid text email type by removing HTML if present
	msg.HTML = ""
	return msg
}

// BuildHTMLEmail creates an HTML-only email message.
// The resulting message will have EmailType EmailTypeHTML.
//
// See also: [EmailTypeHTML], [GetEmailType].
func (b *MailtrapBuilder) BuildHTMLEmail() *Message {
	msg := b.Build()
	// Ensure it's a valid HTML email type by removing text if present
	msg.Text = ""
	return msg
}

// BuildTextAndHTMLEmail creates an email with both text and HTML content.
// The resulting message will have EmailType EmailTypeTextAndHTML.
//
// See also: [EmailTypeTextAndHTML], [GetEmailType].
func (b *MailtrapBuilder) BuildTextAndHTMLEmail() *Message {
	return b.Build() // Both text and HTML should be present
}

// BuildTemplateEmail creates a template-based email message.
// The resulting message will have EmailType EmailTypeTemplate.
//
// Note: subject, text, and html fields are cleared as they should not be used with templates.
//
// See also: [EmailTypeTemplate], [GetEmailType], [TemplateUUID], [TemplateVariables].
func (b *MailtrapBuilder) BuildTemplateEmail() *Message {
	msg := b.Build()
	// Clear content fields as they're forbidden with templates
	msg.Subject = ""
	msg.Text = ""
	msg.HTML = ""
	return msg
}
