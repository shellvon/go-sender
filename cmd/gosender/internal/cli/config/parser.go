package config

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/dingtalk"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/webhook"
	"github.com/shellvon/go-sender/providers/wecombot"
)

// AccountParser handles parsing of account configurations without mapstructure
type AccountParser struct{}

// NewAccountParser creates a new account parser
func NewAccountParser() *AccountParser {
	return &AccountParser{}
}

// ParseSMSAccount parses SMS account configuration
func (p *AccountParser) ParseSMSAccount(raw map[string]interface{}) (*sms.Account, error) {
	acc := &sms.Account{}

	// Parse AccountMeta fields
	if provider, ok := raw["provider"].(string); ok {
		acc.Provider = provider
	}
	if subtype, ok := raw["subtype"].(string); ok {
		acc.SubType = subtype
	}
	if name, ok := raw["name"].(string); ok {
		acc.Name = name
	}
	if weight, ok := raw["weight"]; ok {
		if w, ok := weight.(int); ok {
			acc.Weight = w
		} else if w, ok := weight.(float64); ok {
			acc.Weight = int(w)
		}
	}
	if disabled, ok := raw["disabled"].(bool); ok {
		acc.Disabled = disabled
	}

	// Parse Credentials fields
	if apiKey, ok := raw["api_key"].(string); ok {
		acc.APIKey = apiKey
	}
	if apiSecret, ok := raw["api_secret"].(string); ok {
		acc.APISecret = apiSecret
	}
	if appID, ok := raw["app_id"].(string); ok {
		acc.AppID = appID
	}

	// Parse SMS-specific fields
	if region, ok := raw["region"].(string); ok {
		acc.Region = region
	}
	if callback, ok := raw["callback"].(string); ok {
		acc.Callback = callback
	}
	if signName, ok := raw["sign_name"].(string); ok {
		acc.SignName = signName
	}

	return acc, nil
}

// ParseEmailAccount parses Email account configuration
func (p *AccountParser) ParseEmailAccount(raw map[string]interface{}) (*email.Account, error) {
	acc := &email.Account{}

	// Parse AccountMeta fields
	if provider, ok := raw["provider"].(string); ok {
		acc.Provider = provider
	}
	if name, ok := raw["name"].(string); ok {
		acc.Name = name
	}
	if weight, ok := raw["weight"]; ok {
		if w, ok := weight.(int); ok {
			acc.Weight = w
		} else if w, ok := weight.(float64); ok {
			acc.Weight = int(w)
		}
	}
	if disabled, ok := raw["disabled"].(bool); ok {
		acc.Disabled = disabled
	}

	// Parse Credentials fields
	if apiKey, ok := raw["api_key"].(string); ok {
		acc.APIKey = apiKey
	}
	if apiSecret, ok := raw["api_secret"].(string); ok {
		acc.APISecret = apiSecret
	}

	// Parse Email-specific fields
	if host, ok := raw["host"].(string); ok {
		acc.Host = host
	}
	if port, ok := raw["port"]; ok {
		if p, ok := port.(int); ok {
			acc.Port = p
		} else if p, ok := port.(float64); ok {
			acc.Port = int(p)
		}
	}
	if from, ok := raw["from"].(string); ok {
		acc.From = from
	}

	return acc, nil
}

// ParseDingTalkAccount parses DingTalk account configuration
func (p *AccountParser) ParseDingTalkAccount(raw map[string]interface{}) (*dingtalk.Account, error) {
	acc := &dingtalk.Account{}

	// Parse AccountMeta fields
	if provider, ok := raw["provider"].(string); ok {
		acc.Provider = provider
	}
	if name, ok := raw["name"].(string); ok {
		acc.Name = name
	}
	if weight, ok := raw["weight"]; ok {
		if w, ok := weight.(int); ok {
			acc.Weight = w
		} else if w, ok := weight.(float64); ok {
			acc.Weight = int(w)
		}
	}
	if disabled, ok := raw["disabled"].(bool); ok {
		acc.Disabled = disabled
	}

	// Parse Credentials fields
	if apiKey, ok := raw["api_key"].(string); ok {
		acc.APIKey = apiKey
	}
	if apiSecret, ok := raw["api_secret"].(string); ok {
		acc.APISecret = apiSecret
	}

	return acc, nil
}

// ParseWebhookEndpoint parses Webhook endpoint configuration
func (p *AccountParser) ParseWebhookEndpoint(raw map[string]interface{}) (*webhook.Endpoint, error) {
	ep := &webhook.Endpoint{}

	// Parse AccountMeta fields
	// if provider, ok := raw["provider"].(string); ok {
	// 	// ep.Provider = provider
	// }
	if name, ok := raw["name"].(string); ok {
		ep.Name = name
	}
	if weight, ok := raw["weight"]; ok {
		if w, ok := weight.(int); ok {
			ep.Weight = w
		} else if w, ok := weight.(float64); ok {
			ep.Weight = int(w)
		}
	}
	if disabled, ok := raw["disabled"].(bool); ok {
		ep.Disabled = disabled
	}

	// Parse Webhook-specific fields
	if url, ok := raw["url"].(string); ok {
		ep.URL = url
	}
	if method, ok := raw["method"].(string); ok {
		ep.Method = method
	}
	if headers, ok := raw["headers"].(map[string]interface{}); ok {
		ep.Headers = make(map[string]string)
		for k, v := range headers {
			if str, ok := v.(string); ok {
				ep.Headers[k] = str
			}
		}
	}

	return ep, nil
}

// ParseWeComBotAccount parses WeComBot account configuration
func (p *AccountParser) ParseWeComBotAccount(raw map[string]interface{}) (*wecombot.Account, error) {
	acc := &wecombot.Account{}

	// Parse AccountMeta fields
	if provider, ok := raw["provider"].(string); ok {
		acc.Provider = provider
	}
	if name, ok := raw["name"].(string); ok {
		acc.Name = name
	}
	if weight, ok := raw["weight"]; ok {
		if w, ok := weight.(int); ok {
			acc.Weight = w
		} else if w, ok := weight.(float64); ok {
			acc.Weight = int(w)
		}
	}
	if disabled, ok := raw["disabled"].(bool); ok {
		acc.Disabled = disabled
	}

	// Parse Credentials fields
	if apiKey, ok := raw["api_key"].(string); ok {
		acc.APIKey = apiKey
	}

	return acc, nil
}

// ParseAccounts parses all accounts from configuration
func (p *AccountParser) ParseAccounts(config *cli.RootConfig) (
	[]*sms.Account,
	[]*email.Account,
	[]*dingtalk.Account,
	[]*webhook.Endpoint,
	[]*wecombot.Account,
	error,
) {
	var (
		smsAccounts      []*sms.Account
		emailAccounts    []*email.Account
		dingtalkAccounts []*dingtalk.Account
		webhookEndpoints []*webhook.Endpoint
		wecomAccounts    []*wecombot.Account
	)

	for i, raw := range config.Accounts {
		provider, ok := raw["provider"].(string)
		if !ok || provider == "" {
			return nil, nil, nil, nil, nil, fmt.Errorf("accounts[%d] missing provider field", i)
		}

		switch provider {
		case string(core.ProviderTypeSMS):
			acc, err := p.ParseSMSAccount(raw)
			if err != nil {
				return nil, nil, nil, nil, nil, fmt.Errorf("parse sms account at index %d: %w", i, err)
			}
			smsAccounts = append(smsAccounts, acc)

		case string(core.ProviderTypeEmail):
			acc, err := p.ParseEmailAccount(raw)
			if err != nil {
				return nil, nil, nil, nil, nil, fmt.Errorf("parse email account at index %d: %w", i, err)
			}
			emailAccounts = append(emailAccounts, acc)

		case string(core.ProviderTypeDingtalk):
			acc, err := p.ParseDingTalkAccount(raw)
			if err != nil {
				return nil, nil, nil, nil, nil, fmt.Errorf("parse dingtalk account at index %d: %w", i, err)
			}
			dingtalkAccounts = append(dingtalkAccounts, acc)

		case string(core.ProviderTypeWebhook):
			ep, err := p.ParseWebhookEndpoint(raw)
			if err != nil {
				return nil, nil, nil, nil, nil, fmt.Errorf("parse webhook endpoint at index %d: %w", i, err)
			}
			webhookEndpoints = append(webhookEndpoints, ep)

		case string(core.ProviderTypeWecombot):
			acc, err := p.ParseWeComBotAccount(raw)
			if err != nil {
				return nil, nil, nil, nil, nil, fmt.Errorf("parse wecombot account at index %d: %w", i, err)
			}
			wecomAccounts = append(wecomAccounts, acc)

		default:
			return nil, nil, nil, nil, nil, fmt.Errorf("unsupported provider type: %s", provider)
		}
	}

	return smsAccounts, emailAccounts, dingtalkAccounts, webhookEndpoints, wecomAccounts, nil
}
