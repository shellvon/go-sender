package wecombot

import (
	"github.com/shellvon/go-sender/core"
)

// Account 表示一个企业微信机器人账户。
// 它遵循三层设计：账户元数据 + 凭证 + 额外字段
//   - 账户元数据：名称、权重、禁用状态（来自 core.BaseAccount）
//   - 凭证：APIKey（机器人密钥）、APISecret（可选）、AppID（可选）（来自 core.BaseAccount）
//   - 额外字段：企业微信机器人不需要额外的字段
type Account struct {
	core.BaseAccount
}

// AccountOption 表示一个修改企业微信机器人账户配置的函数。
type AccountOption func(*Account)

// NewAccount 创建一个新的企业微信机器人账户，使用给定的机器人密钥和选项。
//
// 示例：
//
//	account := wecombot.NewAccount("your-bot-key",
//	    wecombot.Name("wecom-main"),
//	    wecombot.Weight(2))
func NewAccount(botKey string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeWecombot,
		"wecombot-default",
		"", // 企业微信机器人不使用子类型
		core.Credentials{
			APIKey: botKey,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, subType string) string {
			return defaultName // 企业微信机器人使用固定的默认名称
		},
		opts...,
	)
}

// 重新导出的核心账户选项，以提供更简洁的 API
// 这些选项提供了便捷的别名：wecombot.Name("test") 而不是 core.WithName[*wecombot.Account]("test")。
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
