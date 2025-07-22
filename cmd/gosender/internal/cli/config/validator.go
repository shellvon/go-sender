package config

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
)

// Validator handles configuration validation
type Validator struct{}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateConfiguration validates the entire configuration
func (v *Validator) ValidateConfiguration(config *cli.RootConfig) []string {
	var errors []string

	if len(config.Accounts) == 0 {
		errors = append(errors, "no accounts configured")
		return errors
	}

	for i, account := range config.Accounts {
		provider, ok := account["provider"].(string)
		if !ok || provider == "" {
			errors = append(errors, fmt.Sprintf("account[%d]: missing or invalid provider field", i))
			continue
		}

		switch provider {
		case string(core.ProviderTypeSMS):
			if errs := v.validateSMSAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeEmail):
			if errs := v.validateEmailAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeDingtalk):
			if errs := v.validateDingTalkAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeWebhook):
			if errs := v.validateWebhookAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeTelegram):
			if errs := v.validateTelegramAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeLark):
			if errs := v.validateLarkAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeWecombot):
			if errs := v.validateWecomBotAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeServerChan):
			if errs := v.validateServerChanAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		case string(core.ProviderTypeEmailAPI):
			if errs := v.validateEmailAPIAccount(account, i); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		default:
			errors = append(errors, fmt.Sprintf("account[%d]: unsupported provider type: %s", i, provider))
		}
	}

	return errors
}

func (v *Validator) validateSMSAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["api_key"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: SMS provider missing api_key", index))
	}
	if _, ok := account["api_secret"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: SMS provider missing api_secret", index))
	}
	if _, ok := account["sign_name"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: SMS provider missing sign_name", index))
	}

	return errors
}

func (v *Validator) validateEmailAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["host"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Email provider missing host", index))
	}
	if _, ok := account["port"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Email provider missing port", index))
	}
	if _, ok := account["api_key"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Email provider missing api_key (username)", index))
	}
	if _, ok := account["api_secret"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Email provider missing api_secret (password)", index))
	}
	if _, ok := account["from"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Email provider missing from address", index))
	}

	return errors
}

func (v *Validator) validateDingTalkAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["api_key"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: DingTalk provider missing api_key (webhook token)", index))
	}

	return errors
}

func (v *Validator) validateWebhookAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["url"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Webhook provider missing url", index))
	}
	if _, ok := account["method"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Webhook provider missing method", index))
	}

	return errors
}

func (v *Validator) validateTelegramAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["api_key"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Telegram provider missing api_key (bot token)", index))
	}

	return errors
}

func (v *Validator) validateLarkAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["api_key"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: Lark provider missing api_key (webhook url)", index))
	}

	return errors
}

func (v *Validator) validateWecomBotAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["api_key"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: WecomBot provider missing api_key (webhook key)", index))
	}

	return errors
}

func (v *Validator) validateServerChanAccount(account map[string]interface{}, index int) []string {
	var errors []string

	if _, ok := account["api_key"]; !ok {
		errors = append(errors, fmt.Sprintf("account[%d]: ServerChan provider missing api_key (SCT token)", index))
	}

	return errors
}

func (v *Validator) validateEmailAPIAccount(account map[string]interface{}, index int) []string {
	var errors []string

	subtype, ok := account["subtype"].(string)
	if !ok || subtype == "" {
		errors = append(errors, fmt.Sprintf("account[%d]: EmailAPI provider missing subtype", index))
		return errors
	}

	switch subtype {
	case "resend":
		if _, ok := account["api_key"]; !ok {
			errors = append(errors, fmt.Sprintf("account[%d]: Resend EmailAPI missing api_key", index))
		}
		if _, ok := account["from"]; !ok {
			errors = append(errors, fmt.Sprintf("account[%d]: Resend EmailAPI missing from address", index))
		}
	case "emailjs":
		if _, ok := account["api_key"]; !ok {
			errors = append(errors, fmt.Sprintf("account[%d]: EmailJS missing api_key (user_id)", index))
		}
		if _, ok := account["api_secret"]; !ok {
			errors = append(errors, fmt.Sprintf("account[%d]: EmailJS missing api_secret (access_token)", index))
		}
		if _, ok := account["service_id"]; !ok {
			errors = append(errors, fmt.Sprintf("account[%d]: EmailJS missing service_id", index))
		}
		if _, ok := account["template_id"]; !ok {
			errors = append(errors, fmt.Sprintf("account[%d]: EmailJS missing template_id", index))
		}
	default:
		errors = append(errors, fmt.Sprintf("account[%d]: unsupported EmailAPI subtype: %s", index, subtype))
	}

	return errors
}
