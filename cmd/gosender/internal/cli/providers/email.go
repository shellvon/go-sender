package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
)

// createEmailProvider 从账户列表创建 Email Provider
func createEmailProvider(accounts []*email.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid email accounts found")
	}

	cfg := &email.Config{Items: accounts}
	return email.New(cfg)
}

// buildEmailMessage 从 CLI 标志构建邮件消息
func buildEmailMessage(flags *cli.CLIFlags) (*email.Message, error) {
	builder := email.Email().
		To(flags.To...).
		Body(flags.Content)

	if flags.Subject != "" {
		builder = builder.Subject(flags.Subject)
	}

	if flags.HTML {
		builder = builder.HTML()
	}

	if len(flags.Files) > 0 {
		builder = builder.Attach(flags.Files...)
	}

	return builder.Build(), nil
}

// validateEmailFlags 验证 CLI 标志是否符合邮件发送要求
func validateEmailFlags(flags *cli.CLIFlags) error {
	if len(flags.To) == 0 {
		return fmt.Errorf("email requires at least one recipient")
	}

	if flags.Subject == "" {
		return fmt.Errorf("email requires a subject")
	}

	if flags.Content == "" {
		return fmt.Errorf("email requires content")
	}

	return nil
}

// NewEmailBuilder 创建一个新的 Email GenericBuilder
func NewEmailBuilder() *GenericBuilder[*email.Account, *email.Message] {
	return NewGenericBuilder(
		core.ProviderTypeEmail,
		createEmailProvider,
		buildEmailMessage,
		validateEmailFlags,
	)
}
