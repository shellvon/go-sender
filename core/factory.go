package core

import (
	"errors"
	"fmt"
)

// Common errors for factory functions.
var (
	ErrNoItems = errors.New("at least one item is required")
)

// CreateAccount creates a new account with the given configuration.
func CreateAccount[T BasicAccount, OptT ~func(T)](
	providerType ProviderType,
	defaultName string,
	subType string,
	credentials Credentials,
	constructor func(BaseAccount) T,
	nameGenerator func(string, string) string,
	opts ...OptT,
) T {
	name := nameGenerator(defaultName, subType)

	baseAccount := BaseAccount{
		AccountMeta: AccountMeta{
			Name:     name,
			SubType:  subType,
			Provider: string(providerType),
			Weight:   1,
			Disabled: false,
		},
		Credentials: credentials,
	}

	account := constructor(baseAccount)

	for _, opt := range opts {
		opt(account)
	}

	return account
}

// CreateProvider creates a new provider with the given configuration.
func CreateProvider[T Selectable, ConfigT interface{ Validate() error }, OptT ~func(ConfigT), ProviderT Provider](
	items []T,
	providerType ProviderType,
	configConstructor func(ProviderMeta, []T) ConfigT,
	providerConstructor func(ConfigT) (ProviderT, error),
	opts ...OptT,
) (ProviderT, error) {
	var zero ProviderT

	if len(items) == 0 {
		return zero, ErrNoItems
	}

	// Create config with default values
	config := configConstructor(ProviderMeta{
		Strategy: StrategyRoundRobin,
		Disabled: false,
	}, items)

	for _, opt := range opts {
		opt(config)
	}

	if err := config.Validate(); err != nil {
		return zero, fmt.Errorf("invalid %s provider configuration: %w", string(providerType), err)
	}

	return providerConstructor(config)
}
