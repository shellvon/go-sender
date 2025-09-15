package wecomapp

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// Account 代表企业微信应用账号配置
// 遵循三层设计：AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (来自 core.BaseAccount)
//   - Credentials: APIKey (corpid), APISecret (corpsecret), AppID (agentid) (来自 core.BaseAccount)
//   - Extra: 企业微信应用不需要额外字段
//
// 它嵌入了core.BaseAccount，因此自动满足core.BasicAccount
// 和core.Selectable接口.
type Account struct {
	core.BaseAccount
}

// AccountOption 代表修改企业微信应用Account配置的函数.
type AccountOption func(*Account)

// NewAccount 使用给定的凭据和选项创建新的企业微信应用账号
//
// 示例:
//
//	account := wecomapp.NewAccount("your_corp_id", "your_corp_secret", "your_agent_id",
//	    wecomapp.Name("wecom-main"),
//	    wecomapp.Weight(2))
func NewAccount(corpID, corpSecret, agentID string, opts ...AccountOption) *Account {
	return core.CreateAccount(
		core.ProviderTypeWecomApp,
		"wecomapp-default",
		"", // 企业微信应用不使用subType
		core.Credentials{
			APIKey:    corpID,
			APISecret: corpSecret,
			AppID:     agentID,
		},
		func(baseAccount core.BaseAccount) *Account {
			return &Account{
				BaseAccount: baseAccount,
			}
		},
		func(defaultName, subType string) string {
			return defaultName // 企业微信应用使用固定的默认名称
		},
		opts...,
	)
}

// CorpID 返回企业ID.
func (a *Account) CorpID() string {
	return a.Credentials.APIKey
}

// CorpSecret 返回应用密钥.
func (a *Account) CorpSecret() string {
	return a.Credentials.APISecret
}

// AgentID 返回应用ID.
func (a *Account) AgentID() string {
	return a.Credentials.AppID
}

// Validate 验证企业微信应用账号配置
// 确保提供了所有必需的企业微信凭据.
func (a *Account) Validate() error {
	// 首先运行基础验证
	if err := a.BaseAccount.Validate(); err != nil {
		return err
	}

	// 企业微信应用特定的验证
	if a.CorpID() == "" {
		return errors.New("corpid is required for WeChat Work Application")
	}

	if a.CorpSecret() == "" {
		return errors.New("corpsecret is required for WeChat Work Application")
	}

	if a.AgentID() == "" {
		return errors.New("agentid is required for WeChat Work Application")
	}

	return nil
}

// Re-exported core account options for cleaner API
// These provide convenient aliases: wecomapp.Name("test") instead of core.WithName[*wecomapp.Account]("test").
var (
	Name     = core.WithName[*Account]
	Weight   = core.WithWeight[*Account]
	Disabled = core.WithDisabled[*Account]
	AppID    = core.WithAppID[*Account]
)
