package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/wneessen/go-mail"

	"github.com/shellvon/go-sender/core"
)

// Config is a type alias for core.BaseConfig[*Account].
type Config = core.BaseConfig[*Account]

// Provider supports multiple accounts and strategy selection.
type Provider struct {
	config *Config
}

var _ core.Provider = (*Provider)(nil)

// New creates a new email provider instance.
func New(config *Config) (*Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &Provider{
		config: config,
	}, nil
}

// ProviderOption represents a function that modifies Email Provider configuration.
type ProviderOption func(*Config)

// NewProvider creates a new Email provider with the given accounts and options.
//
// At least one account is required.
//
// Example:
//
//	provider, err := email.NewProvider([]*email.Account{account1, account2},
//	    email.Strategy(core.StrategyWeighted))
func NewProvider(accounts []*Account, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		accounts,
		core.ProviderTypeEmail,
		func(meta core.ProviderMeta, items []*Account) *Config {
			return &Config{
				ProviderMeta: meta,
				Items:        items,
			}
		},
		New,
		opts...,
	)
}

// Send sends an email message and returns the result.
func (p *Provider) Send(
	ctx context.Context,
	message core.Message,
	_ *core.ProviderSendOptions,
) (*core.SendResult, error) {
	emailMsg, ok := message.(*Message)
	if !ok {
		return nil, core.NewParamError(fmt.Sprintf("invalid message type: expected *email.Message, got %T", message))
	}
	if err := emailMsg.Validate(); err != nil {
		return nil, err
	}
	account, err := p.config.Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	err = p.doSendEmail(ctx, account, emailMsg)

	return &core.SendResult{StatusCode: 0}, err
}

func buildMailOptions(account *Account) []mail.Option {
	opts := []mail.Option{
		mail.WithPort(account.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(account.Username()),
		mail.WithPassword(account.Password()),
		// mail.WithDebugLog(),
	}

	switch account.Port {
	case mail.DefaultPortSSL:
		opts = append(opts, mail.WithSSL())
	case mail.DefaultPortTLS, mail.DefaultPort:
		opts = append(opts, mail.WithTLSPolicy(mail.TLSOpportunistic))
	default:
		opts = append(opts, mail.WithTLSPolicy(mail.TLSOpportunistic))
	}
	return opts
}

// doSendEmail performs the actual email sending.
func (p *Provider) doSendEmail(ctx context.Context, account *Account, emailMsg *Message) error {
	// Set default From if not provided
	if emailMsg.From == "" {
		emailMsg.From = account.From
	}

	// Create new mail message
	m := mail.NewMsg()

	// Set From address
	if err := m.From(emailMsg.From); err != nil {
		return fmt.Errorf("failed to set from: %w", err)
	}

	// Set To addresses
	if err := m.To(emailMsg.To...); err != nil {
		return fmt.Errorf("failed to set to: %w", err)
	}

	// Set CC addresses
	if len(emailMsg.Cc) > 0 {
		if err := m.Cc(emailMsg.Cc...); err != nil {
			return fmt.Errorf("failed to set cc: %w", err)
		}
	}

	// Set BCC addresses
	if len(emailMsg.Bcc) > 0 {
		if err := m.Bcc(emailMsg.Bcc...); err != nil {
			return fmt.Errorf("failed to set bcc: %w", err)
		}
	}

	// Set Reply-To address
	if emailMsg.ReplyTo != "" {
		if err := m.ReplyTo(emailMsg.ReplyTo); err != nil {
			return fmt.Errorf("failed to set reply-to: %w", err)
		}
	}

	// Set subject only if not empty
	if emailMsg.Subject != "" {
		m.Subject(emailMsg.Subject)
	}

	// Set body
	if emailMsg.IsHTML {
		m.SetBodyString(mail.TypeTextHTML, emailMsg.Body)
	} else {
		m.SetBodyString(mail.TypeTextPlain, emailMsg.Body)
	}

	// Add attachments
	for _, att := range emailMsg.Attachments {
		m.AttachFile(att)
	}

	// Create client and send
	client, err := mail.NewClient(account.Host, buildMailOptions(account)...)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	defer client.Close()

	err = client.DialAndSendWithContext(ctx, m)
	var sendErr *mail.SendError
	// https://github.com/wneessen/go-mail/issues/463
	if err != nil && errors.As(err, &sendErr) && sendErr.Reason == mail.ErrSMTPReset {
		err = nil
	}
	return err
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeEmail)
}

func (a *Account) GetName() string {
	return a.Name
}

func (a *Account) GetWeight() int {
	if a.Weight <= 0 {
		return 1
	}
	return a.Weight
}

func (a *Account) IsEnabled() bool {
	return !a.Disabled
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: email.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*email.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)
