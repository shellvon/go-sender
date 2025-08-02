package wecombot

import (
	"github.com/shellvon/go-sender/core"
)

// Account represents a WeCom Bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (bot key), APISecret (optional), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for WeCom Bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}

// AccountOption represents a function that modifies WeCom Bot Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new WeCom Bot account with the given bot key and options.
//
// Example:
//
//	account := wecombot.NewAccount("your-bot-key",
//	    wecombot.Name("wecom-main"),
//	    wecombot.Weight(2))
func NewAccount(botKey string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeWecombot,
		"wecombot-default",
		"", // WeCom Bot doesn't use subType
		core.Credentials{
			APIKey: botKey,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, subType string) string {
			return defaultName // WeCom Bot uses fixed default name
		},
		opts...,
	)
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: wecombot.Name("test") instead of core.WithName[*wecombot.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
