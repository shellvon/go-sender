package serverchan

import (
	"github.com/shellvon/go-sender/core"
)

// Account represents a ServerChan account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (send key), APISecret (optional), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for ServerChan
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}

// AccountOption represents a function that modifies ServerChan Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new ServerChan account with the given SendKey and options.
//
// Example:
//
//	account := serverchan.NewAccount("your-send-key",
//	    serverchan.Name("serverchan-main"),
//	    serverchan.Weight(2))
func NewAccount(sendKey string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeServerChan,
		"serverchan-default",
		"", // ServerChan doesn't use subType
		core.Credentials{
			APIKey: sendKey,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, subType string) string {
			return defaultName // ServerChan uses fixed default name
		},
		opts...,
	)
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: serverchan.Name("test") instead of core.WithName[*serverchan.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
