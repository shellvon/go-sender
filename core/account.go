// Package core provides the core functionality for the go-sender library.
package core

import "errors"

// ProviderMeta contains the minimal cross-cutting configuration that every
// Provider needs, regardless of the concrete business domain.
//
//  1. Disabled —— whether the whole provider is temporarily turned off. When
//     true, Sender will skip registration or selection of this provider.
//  2. Strategy —— load-balancing strategy used to pick one account/endpoint
//     during a Send operation. If not set, defaults to round-robin.
//
// NOTE: This struct purposefully stays extremely small. Anything that is not a
// global, per-provider concern (e.g. region, callback, template set…) should
// live in the provider-specific Config or Account types instead.
type ProviderMeta struct {
	Disabled bool         `json:"disabled" yaml:"disabled"` // When true the provider is globally disabled
	Strategy StrategyType `json:"strategy" yaml:"strategy"` // Load-balancing strategy
}

// GetStrategy returns the configured strategy, falling back to round-robin if
// the field is empty.
func (m *ProviderMeta) GetStrategy() StrategyType {
	if m.Strategy == "" {
		return StrategyRoundRobin
	}
	return m.Strategy
}

// AccountMeta holds cross-cutting metadata for an account. It is the only information
// the core load-balancing / middleware layer relies on.
type AccountMeta struct {
	// Provider is the provider name. common values are sms / email / webhook …
	Provider string `json:"provider"           yaml:"provider"`
	// SubType is the subtype of the provider. common values are aliyun / tencent / mailgun …
	SubType string `json:"subtype,omitempty"  yaml:"subtype,omitempty"` // aliyun / tencent / mailgun …
	// Name is the unique identifier inside same provider.
	Name string `json:"name"               yaml:"name"`
	// Weight is the load-balancing weight (>=1)
	Weight int `json:"weight"             yaml:"weight"`
	// Disabled is whether the account is disabled.
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// GetName returns the account name.
func (m *AccountMeta) GetName() string { return m.Name }

// GetWeight returns weight with sane default.
func (m *AccountMeta) GetWeight() int {
	if m.Weight <= 0 {
		return 1
	}
	return m.Weight
}

// IsEnabled reports whether the account is enabled.
func (m *AccountMeta) IsEnabled() bool { return !m.Disabled }

// GetType returns the subtype string for filtering.
func (m *AccountMeta) GetType() string { return m.SubType }

// Credentials groups the common 3-tuple used by almost every SaaS.
// AppID can represent smsAccount / domain / AppId, etc.
// APIKey can represent username / accessKey / appKey, etc.
// APISecret can represent password / secretKey / appSecret, etc.
type Credentials struct {
	// AppID can represent smsAccount / domain / AppId, etc.
	AppID string `json:"app_id,omitempty"     yaml:"app_id,omitempty"`
	// APIKey can represent username / accessKey / appKey, etc.
	APIKey string `json:"api_key"              yaml:"api_key"`
	// APISecret can represent password / secretKey / appSecret, etc.
	APISecret string `json:"api_secret,omitempty" yaml:"api_secret,omitempty"`
}

// BasicAccount is the minimal interface every provider-specific account must satisfy.
type BasicAccount interface {
	GetMeta() *AccountMeta
	GetCredentials() *Credentials
}

// BaseAccount provides a ready-to-embed implementation that already satisfies
// BasicAccount and core.Selectable. Provider packages can embed it and then add
// their own extension fields.
type BaseAccount struct {
	AccountMeta
	Credentials
}

// BaseAccount implements core.Selectable and core.BasicAccount.
var _ Selectable = (*BaseAccount)(nil)
var _ BasicAccount = (*BaseAccount)(nil)
var _ Validatable = (*BaseAccount)(nil)

// GetMeta returns the embedded metadata.
func (a *BaseAccount) GetMeta() *AccountMeta { return &a.AccountMeta }

// GetCredentials returns the embedded credentials.
func (a *BaseAccount) GetCredentials() *Credentials { return &a.Credentials }

func (a *BaseAccount) Validate() error {
	if a.AccountMeta.Provider == "" {
		return errors.New("provider is required")
	}
	if a.AccountMeta.Name == "" {
		return errors.New("name is required")
	}

	if a.Credentials.AppID == "" && a.Credentials.APIKey == "" && a.Credentials.APISecret == "" {
		return errors.New("credentials are required")
	}

	return nil
}
