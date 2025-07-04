package emailapi

// EmailjsMessageBuilder provides a builder for EmailJS messages, supporting both standard fields and template params.
type EmailjsMessageBuilder struct {
	BaseBuilder[*EmailjsMessageBuilder]

	serviceID string
}

// Emailjs returns a new EmailjsMessageBuilder.
func Emailjs() *EmailjsMessageBuilder {
	b := &EmailjsMessageBuilder{}
	b.self = b
	return b
}

// ServiceID sets the service ID for the EmailJS message.
// This is the same as setting the From field.
//   - If not set, the default service ID will be used. set to reserved keyword `default_service`.
//   - If set, the From field will be set to the service ID.
//     that means serviceId priority is higher than From field.
//
// Reference: https://www.emailjs.com/docs/rest-api/send/
func (b *EmailjsMessageBuilder) ServiceID(serviceID string) *EmailjsMessageBuilder {
	b.serviceID = serviceID
	return b
}

// Build returns the constructed *Message, merging standard fields into template params if not set.
func (b *EmailjsMessageBuilder) Build() *Message {
	msg := b.BuildMessage(string(SubProviderEmailJS))
	// serviceId priority is higher than From field.
	if b.serviceID != "" {
		msg.From = b.serviceID
	}
	return msg
}
