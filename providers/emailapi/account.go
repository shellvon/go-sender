package emailapi

import (
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// Account represents a single Email API service account (Mailgun, Resend, EmailJS etc.).
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey, APISecret, AppID (e.g., Mailgun domain) (from core.BaseAccount)
//   - Extra: Region, Callback, From (API service-specific configuration)
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount

	// Region is the API service region.
	Region string `json:"region,omitempty"   yaml:"region,omitempty"`
	// Callback is the callback URL for webhooks.
	Callback string `json:"callback,omitempty" yaml:"callback,omitempty"`
	// From is the default "From" address for emails.
	From string `json:"from,omitempty"     yaml:"from,omitempty"`
}

// AccountOption represents a function that modifies Email API Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new Email API account with the given configuration and options.
//
// The subType parameter specifies the Email API provider (e.g., "resend", "sendgrid").
//
// Example:
//
//	account := emailapi.NewAccount("resend", "your-api-key",
//	    emailapi.Name("resend-main"),
//	    emailapi.Weight(3),
//	    emailapi.WithRegion("us-east-1"),
//	    emailapi.WithCallback("https://example.com/webhook"))
func NewAccount(subType, apiKey string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeEmailAPI,
		"emailapi-default",
		subType,
		core.Credentials{
			APIKey: apiKey,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, subType string) string {
			if subType != "" {
				return fmt.Sprintf("%s-default", subType)
			}
			return defaultName
		},
		opts...,
	)
}

// Validate checks if the account is valid.
// It ensures that the subType is set for Email API providers.
func (a *Account) Validate() error {
	if a.SubType == "" {
		return errors.New("subType is required for Email API provider")
	}
	return a.BaseAccount.Validate()
}

// EmailAPI-specific account options

// WithRegion sets the email API service region.
func WithRegion(region string) AccountOption {
	return func(account *Account) {
		account.Region = region
	}
}

// WithCallback sets the callback URL for webhooks.
func WithCallback(callback string) AccountOption {
	return func(account *Account) {
		account.Callback = callback
	}
}

// WithFrom sets the default "From" address for emails.
func WithFrom(from string) AccountOption {
	return func(account *Account) {
		account.From = from
	}
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: emailapi.Name("test") instead of core.WithName[*emailapi.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
