package sms

import (
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// Account represents a single SMS account configuration.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey, APISecret (from core.BaseAccount)
//   - Extra: SignName, Region, Callback (SMS-specific configuration)
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount

	// SMS-specific configuration
	SignName string `json:"sign_name,omitempty"` // SMS signature/sign name
	Region   string `json:"region,omitempty"`    // SMS service region (e.g. cn-hangzhou)
	Callback string `json:"callback,omitempty"`  // Callback URL for delivery reports
}

// Username returns the API key for authentication.
func (a *Account) Username() string { return a.GetCredentials().APIKey }

// Password returns the API secret for authentication.
func (a *Account) Password() string { return a.GetCredentials().APISecret }

// SubProvider returns the SMS sub-provider name (e.g., "aliyun", "tencent").
func (a *Account) SubProvider() string { return a.AccountMeta.SubType }

// Validate validates the SMS account configuration.
func (a *Account) Validate() error {
	// First run the base validation
	if err := a.BaseAccount.Validate(); err != nil {
		return err
	}

	// SMS-specific validation: SubType is required
	if a.AccountMeta.SubType == "" {
		return errors.New("subType is required for SMS accounts")
	}

	return nil
}

// AccountOption represents a function that modifies SMS Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new SMS account with the given configuration and options.
//
// The subType parameter specifies the SMS provider (e.g., "aliyun", "tencent").
// Uses pure generics for type safety without reflection.
//
// Example:
//
//	account := sms.NewAccount("aliyun", "key", "secret",
//	    sms.Name("aliyun-main"),           // Clean API through re-export
//	    sms.Weight(3),                     // Clean API through re-export
//	    sms.WithSignName("MyApp"),         // SMS-specific option
//	    sms.WithRegion("cn-hangzhou"))     // SMS-specific option
func NewAccount(subType, apiKey, apiSecret string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeSMS,
		"sms-default",
		subType,
		core.Credentials{
			APIKey:    apiKey,
			APISecret: apiSecret,
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

// SMS-specific account options

// WithSignName sets the default SMS signature for this account.
func WithSignName(signName string) AccountOption {
	return func(a *Account) {
		a.SignName = signName
	}
}

// WithRegion sets the SMS service region.
func WithRegion(region string) AccountOption {
	return func(a *Account) {
		a.Region = region
	}
}

// WithCallback sets the callback URL for delivery reports.
func WithCallback(callback string) AccountOption {
	return func(a *Account) {
		a.Callback = callback
	}
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: sms.Name("test") instead of core.WithName[*sms.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
