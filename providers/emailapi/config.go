package emailapi

import "github.com/shellvon/go-sender/core"

// Config holds configuration for the emailapi provider and all supported API email services.
type Config struct {
	// BaseConfig contains common configuration like strategy and disabled flag
	core.BaseConfig
	// Accounts contains multiple API accounts for load balancing, failover, etc.
	// Each account can be configured for different sub-providers (emailjs, resend, etc.)
	Accounts []Account `json:"accounts"`
}

// Account holds credentials and settings for a single API account.
type Account struct {
	// Name is the unique identifier for this account, used for account selection and logging
	Name string `json:"name"`
	// Type specifies the sub-provider type (e.g., "emailjs", "resend")
	Type SubProviderType `json:"type"`
	// APIKey is the authentication key for the email service API
	APIKey string `json:"api_key"`
	// APISecret is the secret key for API authentication (optional for some providers)
	APISecret string `json:"api_secret,omitempty"`
	// Domain is the custom domain for sending emails (e.g., for Mailgun)
	Domain string `json:"domain,omitempty"`
	// Region specifies the service region for multi-region providers
	Region string `json:"region,omitempty"`
	// From is the default sender email address
	From string `json:"from"`
	// ReplyTo is the reply-to email address for responses
	ReplyTo string `json:"reply_to,omitempty"`
	// Weight determines the load balancing weight for this account (higher weight = more traffic)
	Weight int `json:"weight,omitempty"`
	// Disabled indicates whether this account is disabled and should not be used
	Disabled bool `json:"disabled,omitempty"`
	// Extras contains provider-specific configuration parameters
	Extras map[string]interface{} `json:"extras,omitempty"`
}

// IsConfigured checks if the EmailAPI configuration is valid and ready to use.
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Accounts) > 0
}
