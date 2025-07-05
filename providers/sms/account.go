package sms

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// Account represents a single SMS service account (any sub-provider).
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey, APISecret, AppID (vendor-specific) (from core.BaseAccount)
//   - Extra: Region, Callback, SignName (optional defaults for SMS vendors)
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces. Region / Callback / SignName are optional defaults that can be
// overridden on a per-message basis via Message.
type Account struct {
	core.BaseAccount

	// Region is the SMS service region.
	Region string `json:"region,omitempty"    yaml:"region,omitempty"`
	// Callback is the callback URL for delivery reports.
	Callback string `json:"callback,omitempty"  yaml:"callback,omitempty"`
	// SignName is the default SMS signature for this account.
	// This is commonly required by Chinese SMS providers like Aliyun, Tencent, etc.
	// Can be overridden per message via Message.
	SignName string `json:"sign_name,omitempty" yaml:"sign_name,omitempty"`
}

// Validate checks if the account is valid.
// It ensures that the subType is set for SMS providers.
func (a *Account) Validate() error {
	if a.SubType == "" {
		return errors.New("subType is required for SMS provider")
	}

	if a.SubType == string(SubProviderVolc) && a.AppID == "" {
		return errors.New("smsAccount(appID) is required for Volc SMS provider")
	}

	return a.BaseAccount.Validate()
}
