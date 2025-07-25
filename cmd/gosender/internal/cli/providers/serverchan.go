package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/serverchan"
)

// createServerChanProvider 从账户列表创建 ServerChan Provider
func createServerChanProvider(accounts []*serverchan.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid serverchan accounts found")
	}

	cfg := &serverchan.Config{Items: accounts}
	return serverchan.New(cfg)
}

// buildServerChanMessage 从 CLI 标志构建 ServerChan 消息
func buildServerChanMessage(flags *cli.CLIFlags) (*serverchan.Message, error) {
	builder := serverchan.Text().
		Title(flags.Subject).
		Content(flags.Content)

	if channel, ok := flags.Metadata["channel"]; ok {
		builder = builder.Channel(channel)
	}

	if openid, ok := flags.Metadata["openid"]; ok {
		builder = builder.OpenID(openid)
	}

	if _, ok := flags.Metadata["noip"]; ok {
		builder = builder.NoIP()
	}

	if short, ok := flags.Metadata["short"]; ok {
		builder = builder.Short(short)
	}

	return builder.Build(), nil
}

// validateServerChanFlags 验证 CLI 标志是否符合 ServerChan 发送要求
func validateServerChanFlags(flags *cli.CLIFlags) error {
	if flags.Subject == "" {
		return fmt.Errorf("serverchan requires a subject (title)")
	}

	if flags.Content == "" {
		return fmt.Errorf("serverchan requires content")
	}

	return nil
}

// NewServerChanBuilder 创建一个新的 ServerChan GenericBuilder
func NewServerChanBuilder() *GenericBuilder[*serverchan.Account, *serverchan.Message] {
	return NewGenericBuilder(
		core.ProviderTypeServerChan,
		createServerChanProvider,
		buildServerChanMessage,
		validateServerChanFlags,
	)
}
