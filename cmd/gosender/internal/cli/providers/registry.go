package providers

import (
	"errors"
	"fmt"
	"reflect"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/shellvon/go-sender/core"
)

// 反射使用的方法名常量.
const (
	BuildProviderMethodName = "BuildProvider"
)

// ProviderTypeGetter 是一个简单的接口，只要求实现 GetProviderType 方法.
type ProviderTypeGetter interface {
	GetProviderType() core.ProviderType
}

// MessageBuilder 定义了从 CLI 标志创建消息的接口.
type MessageBuilder interface {
	ProviderTypeGetter
	// BuildMessage 从 CLI 标志创建消息
	BuildMessage(flags *cli.CLIFlags) (core.Message, error)
	// ValidateFlags 验证 CLI 标志是否符合此 Provider 的要求
	ValidateFlags(flags *cli.CLIFlags) error
}

// ProviderFactory 是一个更高级的接口，同时提供 Provider 和 Message 构建能力.
type ProviderFactory interface {
	ProviderTypeGetter
	GetProviderBuilder() ProviderTypeGetter
	GetMessageBuilder() MessageBuilder
}

// ProviderRegistry 管理 Provider 和消息构建器.
type ProviderRegistry struct {
	// 存储不同类型的 Builder
	providerBuilders map[core.ProviderType]ProviderTypeGetter
	messageBuilders  map[core.ProviderType]MessageBuilder
}

// NewProviderRegistry 创建一个新的 Provider 注册表.
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providerBuilders: make(map[core.ProviderType]ProviderTypeGetter),
		messageBuilders:  make(map[core.ProviderType]MessageBuilder),
	}
}

// RegisterProvider 注册任何类型的 Provider Builder.
func (r *ProviderRegistry) RegisterProvider(builder ProviderTypeGetter) error {
	if builder == nil {
		return errors.New("cannot register nil builder")
	}

	providerType := builder.GetProviderType()
	r.providerBuilders[providerType] = builder
	return nil
}

// RegisterMessageBuilder 注册一个消息构建器.
func (r *ProviderRegistry) RegisterMessageBuilder(builder MessageBuilder) {
	r.messageBuilders[builder.GetProviderType()] = builder
}

// Register 同时注册 Provider 和对应的 MessageBuilder.
func (r *ProviderRegistry) Register(providerBuilder ProviderTypeGetter, messageBuilder MessageBuilder) error {
	if providerBuilder == nil || messageBuilder == nil {
		return errors.New("cannot register nil builder")
	}

	// 确保两者类型一致
	if providerBuilder.GetProviderType() != messageBuilder.GetProviderType() {
		return fmt.Errorf("provider type mismatch: %s vs %s",
			providerBuilder.GetProviderType(), messageBuilder.GetProviderType())
	}

	// 注册两个组件
	r.providerBuilders[providerBuilder.GetProviderType()] = providerBuilder
	r.messageBuilders[messageBuilder.GetProviderType()] = messageBuilder

	return nil
}

// RegisterFactory 注册实现了 ProviderFactory 接口的工厂.
func (r *ProviderRegistry) RegisterFactory(factory ProviderFactory) error {
	if factory == nil {
		return errors.New("cannot register nil factory")
	}

	return r.Register(factory.GetProviderBuilder(), factory.GetMessageBuilder())
}

// BuildProviders 从配置构建所有 Provider 并将其注册到 Sender.
func (r *ProviderRegistry) BuildProviders(sender *gosender.Sender, rootConfig *cli.RootConfig) error {
	parser := config.NewAccountParser()
	parsedAccounts, err := parser.ParseAccounts(rootConfig)
	if err != nil {
		return fmt.Errorf("failed to parse accounts: %w", err)
	}

	// 为每种 Provider 类型构建并注册 Provider
	for providerType, accounts := range parsedAccounts {
		builder, exists := r.providerBuilders[providerType]
		if !exists {
			return fmt.Errorf("no provider builder registered for type: %s", providerType)
		}

		// 使用反射来确定具体 Provider 类型并构建 Provider
		builderValue := reflect.ValueOf(builder)

		// 查找 BuildProvider 方法
		buildMethod := builderValue.MethodByName(BuildProviderMethodName)
		if !buildMethod.IsValid() {
			return fmt.Errorf("builder for %s does not have %s method", providerType, BuildProviderMethodName)
		}

		// 获取方法类型并检查参数
		methodType := buildMethod.Type()
		if methodType.NumIn() != 1 {
			return fmt.Errorf("%s method for %s has wrong number of parameters", BuildProviderMethodName, providerType)
		}

		// 创建正确类型的参数
		accountsType := methodType.In(0)
		accountsSlice := reflect.MakeSlice(accountsType, 0, len(accounts))

		// 将账户转换为正确的类型
		for i, acc := range accounts {
			accValue := reflect.ValueOf(acc)
			if !accValue.Type().AssignableTo(accountsType.Elem()) {
				return fmt.Errorf("account at index %d is not compatible with provider %s", i, providerType)
			}
			accountsSlice = reflect.Append(accountsSlice, accValue)
		}

		// 调用 BuildProvider 方法
		results := buildMethod.Call([]reflect.Value{accountsSlice})
		if len(results) != 2 {
			return fmt.Errorf(
				"%s method for %s has wrong number of return values",
				BuildProviderMethodName,
				providerType,
			)
		}

		// 检查错误
		errValue := results[1]
		if !errValue.IsNil() {
			return fmt.Errorf("failed to build provider %s: %w", providerType, errValue.Interface().(error))
		}

		// 获取 Provider
		providerValue := results[0].Interface().(core.Provider)
		sender.RegisterProvider(providerType, providerValue, nil)
	}

	return nil
}

// BuildMessage 为指定的 Provider 类型构建消息.
func (r *ProviderRegistry) BuildMessage(providerType core.ProviderType, flags *cli.CLIFlags) (core.Message, error) {
	builder, exists := r.messageBuilders[providerType]
	if !exists {
		return nil, fmt.Errorf("no message builder registered for provider type: %s", providerType)
	}

	if err := builder.ValidateFlags(flags); err != nil {
		return nil, fmt.Errorf("flag validation failed for %s: %w", providerType, err)
	}

	return builder.BuildMessage(flags)
}

// GetSupportedProviders 返回所有支持的 Provider 类型.
func (r *ProviderRegistry) GetSupportedProviders() []core.ProviderType {
	var providers []core.ProviderType
	for providerType := range r.providerBuilders {
		providers = append(providers, providerType)
	}
	return providers
}

// ValidateFlags 验证指定 Provider 的标志.
func (r *ProviderRegistry) ValidateFlags(providerType core.ProviderType, flags *cli.CLIFlags) error {
	builder, exists := r.messageBuilders[providerType]
	if !exists {
		return fmt.Errorf("no message builder registered for provider type: %s", providerType)
	}

	return builder.ValidateFlags(flags)
}
