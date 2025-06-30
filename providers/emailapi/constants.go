package emailapi

// Internal constants for EmailAPI message extra fields
// These constants are used internally by transformers and should not be exposed to users

// EmailJS specific constants.
const (
	// emailjsServiceID is the key for EmailJS service ID.
	emailjsServiceID = "service_id"

	// emailjsTemplateID is the key for EmailJS template ID.
	emailjsTemplateID = "template_id"

	// emailjsUserID is the key for EmailJS user ID.
	emailjsUserID = "user_id"

	// emailjsAccessToken is the key for EmailJS access token.
	emailjsAccessToken = "accessToken"

	// emailjsTemplateParams is the key for EmailJS template parameters.
	emailjsTemplateParams = "template_params"
)

// Resend specific constants.
const (
	// resendTags is the key for Resend tags.
	resendTags = "tags"
)
