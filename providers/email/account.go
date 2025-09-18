package email

import (
	"github.com/shellvon/go-sender/core"
)

// Account represents a single SMTP account configuration.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (username), APISecret (password) (from core.BaseAccount)
//   - Extra: Host, Port, From (SMTP-specific configuration)
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount

	Host string `json:"host"` // SMTP server host
	Port int    `json:"port"` // SMTP port (25/465/587)
	From string `json:"from"` // Default "From" address
}

// AccountOption represents a function that modifies Email Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new SMTP email account with the given configuration and options.
//
// Example:
//
//	account := email.NewAccount("smtp.gmail.com", 587, "user@gmail.com", "password",
//	    email.Name("gmail-main"),
//	    email.Weight(2),
//	    email.WithFrom("noreply@myapp.com"))
func NewAccount(host string, port int, username, password string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeEmail,
		"email-default",
		"", // Email provider doesn't use subType
		core.Credentials{
			APIKey:    username,
			APISecret: password,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
				Host:        host,
				Port:        port,
				From:        username,
			}
		},
		func(defaultName, subType string) string {
			return defaultName // Email uses fixed default name
		},
		opts...,
	)
}

func (a *Account) Username() string { return a.GetCredentials().APIKey }
func (a *Account) Password() string { return a.GetCredentials().APISecret }

// Email-specific account options

// WithFrom sets the default "From" address for emails.
func WithFrom(from string) AccountOption {
	return func(account *Account) {
		account.From = from
	}
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: email.Name("test") instead of core.WithName[*email.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
