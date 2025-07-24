package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
)

// EmailProviderBuilder builds email providers
type EmailProviderBuilder struct{}

func (b *EmailProviderBuilder) GetProviderType() core.ProviderType {
	return core.ProviderTypeEmail
}

func (b *EmailProviderBuilder) BuildProvider(accounts []any) (core.Provider, error) {
	parser := config.NewAccountParser()
	var emailAccounts []*email.Account

	for i, rawAccount := range accounts {
		accountMap, ok := rawAccount.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid account format at index %d", i)
		}

		account, err := parser.ParseEmailAccount(accountMap)
		if err != nil {
			return nil, fmt.Errorf("failed to parse email account at index %d: %w", i, err)
		}

		emailAccounts = append(emailAccounts, account)
	}

	if len(emailAccounts) == 0 {
		return nil, fmt.Errorf("no valid email accounts found")
	}

	cfg := &email.Config{Items: emailAccounts}
	return email.New(cfg)
}

// EmailMessageBuilder builds email messages
type EmailMessageBuilder struct{}

func (b *EmailMessageBuilder) GetProviderType() core.ProviderType {
	return core.ProviderTypeEmail
}

func (b *EmailMessageBuilder) ValidateFlags(flags *cli.CLIFlags) error {
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

func (b *EmailMessageBuilder) BuildMessage(flags *cli.CLIFlags) (core.Message, error) {
	builder := email.Email().
		To(flags.To...).
		Body(flags.Content)

	if flags.Subject != "" {
		builder = builder.Subject(flags.Subject)
	}

	if flags.HTML {
		builder = builder.HTML()
	}

	// Add file attachments if provided
	if len(flags.Files) > 0 {
		builder = builder.Attach(flags.Files...)
	}

	// TODO: We could extend CLI flags to support CC, BCC, ReplyTo, From
	// For now, these would be set via the builder if needed:
	// if flags.Cc != nil {
	//     builder = builder.Cc(flags.Cc...)
	// }
	// if flags.Bcc != nil {
	//     builder = builder.Bcc(flags.Bcc...)
	// }
	// if flags.ReplyTo != "" {
	//     builder = builder.ReplyTo(flags.ReplyTo)
	// }
	// if flags.From != "" {
	//     builder = builder.From(flags.From)
	// }

	return builder.Build(), nil
}
