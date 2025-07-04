package sms

import "github.com/shellvon/go-sender/core"

// Account represents a single SMS service account (any sub-provider).
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey, APISecret, AppID (vendor-specific) (from core.BaseAccount)
//   - Extra: Region, Callback (optional defaults for SMS vendors)
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces. Region / Callback are optional defaults that can be
// overridden on a per-message basis via Message.Extras.
type Account struct {
	core.BaseAccount

	// Region is the SMS service region.
	Region string `json:"region,omitempty"   yaml:"region,omitempty"`
	// Callback is the callback URL for delivery reports.
	Callback string `json:"callback,omitempty" yaml:"callback,omitempty"`
}
