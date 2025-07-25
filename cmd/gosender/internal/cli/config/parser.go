package config

import (
	"encoding/json"
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/dingtalk"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/serverchan"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/webhook"
	"github.com/shellvon/go-sender/providers/wecombot"
)

// AccountCreator 是创建账户对象的函数类型
type AccountCreator func() core.Selectable

// accountRegistry 存储所有已注册的账户创建器
var accountRegistry = make(map[core.ProviderType]AccountCreator)

// RegisterAccountType 注册一个账户类型的创建器
func RegisterAccountType(providerType core.ProviderType, creator AccountCreator) {
	accountRegistry[providerType] = creator
}

// 初始化内置的账户类型
func init() {
	// 注册所有内置账户类型
	RegisterAccountType(core.ProviderTypeSMS, func() core.Selectable { return &sms.Account{} })
	RegisterAccountType(core.ProviderTypeEmail, func() core.Selectable { return &email.Account{} })
	RegisterAccountType(core.ProviderTypeDingtalk, func() core.Selectable { return &dingtalk.Account{} })
	RegisterAccountType(core.ProviderTypeWebhook, func() core.Selectable { return &webhook.Endpoint{} })
	RegisterAccountType(core.ProviderTypeWecombot, func() core.Selectable { return &wecombot.Account{} })
	RegisterAccountType(core.ProviderTypeServerChan, func() core.Selectable { return &serverchan.Account{} })
}

// AccountParser handles parsing of account configurations
type AccountParser struct{}

// NewAccountParser creates a new account parser
func NewAccountParser() *AccountParser {
	return &AccountParser{}
}

// ParseAccount parses a single account from raw map
func (p *AccountParser) ParseAccount(providerType core.ProviderType, raw map[string]interface{}) (core.Selectable, error) {
	// 从注册表获取账户创建器
	creator, ok := accountRegistry[providerType]
	if !ok {
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}

	// 创建账户实例
	result := creator()

	// Convert map to JSON bytes and unmarshal to the struct
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal account to JSON: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	return result, nil
}

// ParseAccounts parses all accounts from configuration and groups them by provider type
func (p *AccountParser) ParseAccounts(config *cli.RootConfig) (map[core.ProviderType][]core.Selectable, error) {
	// Initialize result map
	accounts := make(map[core.ProviderType][]core.Selectable)

	for i, raw := range config.Accounts {
		providerStr, ok := raw["provider"].(string)
		if !ok || providerStr == "" {
			return nil, fmt.Errorf("accounts[%d] missing provider field", i)
		}

		providerType := core.ProviderType(providerStr)
		account, err := p.ParseAccount(providerType, raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse account at index %d: %w", i, err)
		}

		// Group accounts by provider type
		accounts[providerType] = append(accounts[providerType], account)
	}

	return accounts, nil
}

// Helper methods for specific providers - convenience methods that leverage ParseAccount

// ParseSMSAccount parses SMS account configuration
func (p *AccountParser) ParseSMSAccount(raw map[string]interface{}) (*sms.Account, error) {
	account, err := p.ParseAccount(core.ProviderTypeSMS, raw)
	if err != nil {
		return nil, err
	}
	return account.(*sms.Account), nil
}

// ParseEmailAccount parses Email account configuration
func (p *AccountParser) ParseEmailAccount(raw map[string]interface{}) (*email.Account, error) {
	account, err := p.ParseAccount(core.ProviderTypeEmail, raw)
	if err != nil {
		return nil, err
	}
	return account.(*email.Account), nil
}

// ParseWeComBotAccount parses WeComBot account configuration
func (p *AccountParser) ParseWeComBotAccount(raw map[string]interface{}) (*wecombot.Account, error) {
	account, err := p.ParseAccount(core.ProviderTypeWecombot, raw)
	if err != nil {
		return nil, err
	}
	return account.(*wecombot.Account), nil
}

// ParseServerChanAccount parses ServerChan account configuration
func (p *AccountParser) ParseServerChanAccount(raw map[string]interface{}) (*serverchan.Account, error) {
	account, err := p.ParseAccount(core.ProviderTypeServerChan, raw)
	if err != nil {
		return nil, err
	}
	return account.(*serverchan.Account), nil
}
