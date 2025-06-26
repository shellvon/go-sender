package emailapi

// ProviderType defines the supported email API providers.
type ProviderType string

const (
	ProviderTypeEmailJS ProviderType = "emailjs"
	ProviderTypeResend  ProviderType = "resend"
)
