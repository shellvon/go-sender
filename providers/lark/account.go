package lark

import (
	"github.com/shellvon/go-sender/core"
)

// Account represents a Lark bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (webhook URL), APISecret (signing secret), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for Lark bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}

// AccountOption represents a function that modifies Lark Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new Lark bot account with the given webhook URL and options.
//
// Example:
//
//	account := lark.NewAccount("https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
//	    lark.Name("lark-main"),
//	    lark.Weight(2))
func NewAccount(webhookURL string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeLark,
		"lark-default",
		"", // Lark doesn't use subType
		core.Credentials{
			APIKey: webhookURL,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, _ string) string {
			return defaultName // Lark uses fixed default name
		},
		opts...,
	)
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: lark.Name("test") instead of core.WithName[*lark.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
