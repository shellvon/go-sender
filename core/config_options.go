package core

// Account configuration options

// WithName 设置账号名称 - 使用泛型约束，类型安全.
func WithName[T BasicAccount](name string) func(T) {
	return func(account T) {
		account.GetMeta().Name = name
	}
}

// WithWeight 设置账号权重 - 使用泛型约束，类型安全.
func WithWeight[T BasicAccount](weight int) func(T) {
	return func(account T) {
		account.GetMeta().Weight = weight
	}
}

// WithDisabled 禁用账号 - 使用泛型约束，类型安全.
func WithDisabled[T BasicAccount]() func(T) {
	return func(account T) {
		account.GetMeta().Disabled = true
	}
}

// WithAppID 设置App ID - 使用泛型约束，类型安全.
func WithAppID[T BasicAccount](appID string) func(T) {
	return func(account T) {
		account.GetCredentials().AppID = appID
	}
}

// Provider configuration options

// ProviderMetaAccessor 定义访问ProviderMeta的接口.
type ProviderMetaAccessor interface {
	GetProviderMeta() *ProviderMeta
}

// WithStrategy 设置负载均衡策略 - 使用接口约束，类型安全.
func WithStrategy[T ProviderMetaAccessor](strategy StrategyType) func(T) {
	return func(config T) {
		config.GetProviderMeta().Strategy = strategy
	}
}

// WithProviderDisabled 禁用整个Provider - 使用接口约束，类型安全.
func WithProviderDisabled[T ProviderMetaAccessor]() func(T) {
	return func(config T) {
		config.GetProviderMeta().Disabled = true
	}
}
