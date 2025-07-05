package emailapi

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// Account represents a single Email API service account (Mailgun, Resend, EmailJS etc.).
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey, APISecret, AppID (e.g., Mailgun domain) (from core.BaseAccount)
//   - Extra: Region, Callback (optional defaults for API services)
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount

	// Region is the API service region.
	Region string `json:"region,omitempty"   yaml:"region,omitempty"`
	// Callback is the callback URL for webhooks.
	Callback string `json:"callback,omitempty" yaml:"callback,omitempty"`
}

// Validate checks if the account is valid.
// It ensures that the subType is set for Email API providers.
func (a *Account) Validate() error {
	if a.SubType == "" {
		return errors.New("subType is required for Email API provider")
	}
	return a.BaseAccount.Validate()
}
