package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// Account represents a DingTalk bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (access token), APISecret (signing secret), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for DingTalk bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}

// AccountOption represents a function that modifies DingTalk Account configuration.
type AccountOption func(*Account)

// NewAccount creates a new DingTalk bot account with the given webhook URL and options.
//
// Example:
//
//	account := dingtalk.NewAccount("https://oapi.dingtalk.com/robot/send?access_token=xxx",
//	    dingtalk.Name("main-bot"),
//	    dingtalk.Weight(2))
func NewAccount(webhookURL string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeDingtalk,
		"dingtalk-default",
		"", // DingTalk doesn't use subType
		core.Credentials{
			APIKey: webhookURL,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, subType string) string {
			return defaultName // DingTalk uses fixed default name
		},
		opts...,
	)
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: dingtalk.Name("test") instead of core.WithName[*dingtalk.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
