package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
)

type GenericBuilder[T core.Selectable, M core.Message] struct {
	providerType     core.ProviderType
	createProviderFn func(accounts []T) (core.Provider, error)
	buildMessageFn   func(flags *cli.CLIFlags) (M, error)
	validateFlagsFn  func(flags *cli.CLIFlags) error
}

func NewGenericBuilder[T core.Selectable, M core.Message](
	providerType core.ProviderType,
	createProviderFn func(accounts []T) (core.Provider, error),
	buildMessageFn func(flags *cli.CLIFlags) (M, error),
	validateFlagsFn func(flags *cli.CLIFlags) error,
) *GenericBuilder[T, M] {
	return &GenericBuilder[T, M]{
		providerType:     providerType,
		createProviderFn: createProviderFn,
		buildMessageFn:   buildMessageFn,
		validateFlagsFn:  validateFlagsFn,
	}
}

func (b *GenericBuilder[T, M]) GetProviderType() core.ProviderType {
	return b.providerType
}

func (b *GenericBuilder[T, M]) BuildProvider(accounts []T) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid accounts found for provider %s", b.providerType)
	}
	return b.createProviderFn(accounts)
}

func (b *GenericBuilder[T, M]) BuildMessage(flags *cli.CLIFlags) (core.Message, error) {
	return b.buildMessageFn(flags)
}

func (b *GenericBuilder[T, M]) ValidateFlags(flags *cli.CLIFlags) error {
	return b.validateFlagsFn(flags)
}

func (b *GenericBuilder[T, M]) GetProviderBuilder() ProviderTypeGetter {
	return b
}

func (b *GenericBuilder[T, M]) GetMessageBuilder() MessageBuilder {
	return b
}

func (b *GenericBuilder[T, M]) GetAccountType() core.Selectable {
	var zero T
	return zero
}

func (b *GenericBuilder[T, M]) GetMessageType() core.Message {
	var zero M
	return zero
}
