package telegram

import (
	"github.com/shellvon/go-sender/core"
)

// Account represents a Telegram bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (bot token), APISecret (optional), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for Telegram bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}

// AccountOption represents a function that modifies Telegram Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new Telegram bot account with the given bot token and options.
//
// Example:
//
//	account := telegram.NewAccount("your-bot-token",
//	    telegram.Name("main-bot"),
//	    telegram.Weight(2))
func NewAccount(botToken string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeTelegram,
		"telegram-default",
		"", // Telegram doesn't use subType
		core.Credentials{
			APIKey: botToken,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, subType string) string {
			return defaultName // Telegram uses fixed default name
		},
		opts...,
	)
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: telegram.Name("test") instead of core.WithName[*telegram.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
