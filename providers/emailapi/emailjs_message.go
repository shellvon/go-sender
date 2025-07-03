package emailapi

// NewEmailjsMessage creates a new EmailJS message with required parameters.
// This function is a shortcut for Emailjs().ServiceID(serviceID).TemplateID(templateID).TemplateParams(templateParams).Build().
//
// Reference:
//   - https://www.emailjs.com/docs/rest-api/send/
func NewEmailjsMessage(serviceID, templateID string, templateParams map[string]interface{}) *Message {
	return Emailjs().
		ServiceID(serviceID).
		TemplateID(templateID).
		TemplateParams(templateParams).
		Build()
}
