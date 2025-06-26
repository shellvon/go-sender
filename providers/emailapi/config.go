package emailapi

// Config holds configuration for the emailapi provider and all supported API email services.
type Config struct {
	// ProviderType: e.g., "mailgun", "mailjet", "resend", etc.
	ProviderType string `json:"provider_type"`
	// Accounts: support multiple API accounts for load balancing, failover, etc.
	Accounts []AccountConfig `json:"accounts"`
	// Strategy: selection strategy (round_robin, weighted, etc.)
	Strategy string `json:"strategy"`
	// Disabled: whether this provider is disabled
	Disabled bool `json:"disabled"`
}

// AccountConfig holds credentials and settings for a single API account.
type AccountConfig struct {
	Name      string                 `json:"name"`
	APIKey    string                 `json:"api_key"`
	APISecret string                 `json:"api_secret,omitempty"`
	Domain    string                 `json:"domain,omitempty"` // For providers like Mailgun
	Region    string                 `json:"region,omitempty"`
	From      string                 `json:"from"`
	ReplyTo   string                 `json:"reply_to,omitempty"`
	Weight    int                    `json:"weight,omitempty"`
	Disabled  bool                   `json:"disabled,omitempty"`
	Extras    map[string]interface{} `json:"extras,omitempty"` // Provider-specific params
	// Add more fields as needed for specific providers
}
