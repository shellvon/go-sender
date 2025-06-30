package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
	"github.com/wneessen/go-mail"
)

// Provider supports multiple accounts and strategy selection
type Provider struct {
	accounts []*Account
	strategy core.SelectionStrategy
}

var _ core.Provider = (*Provider)(nil)

// New creates a new email provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("email provider is not configured or is disabled")
	}

	// Convert to pointer slice
	accounts := make([]*Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	// Use common initialization logic
	enabledAccounts, strategy, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled email accounts found")
	}

	return &Provider{
		accounts: enabledAccounts,
		strategy: strategy,
	}, nil
}

// Send sends an email message
func (p *Provider) Send(ctx context.Context, message core.Message, opts *core.ProviderSendOptions) error {
	emailMsg, ok := message.(*Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected *email.Message, got %T", message))
	}

	if err := emailMsg.Validate(); err != nil {
		return err
	}

	// 转换为 Selectable 接口
	selectables := make([]core.Selectable, len(p.accounts))
	for i, account := range p.accounts {
		selectables[i] = account
	}

	selected := utils.Select(ctx, selectables, p.strategy)
	if selected == nil {
		return errors.New("no available account")
	}

	// 找到对应的账号
	var account *Account
	for _, acc := range p.accounts {
		if acc.GetName() == selected.GetName() {
			account = acc
			break
		}
	}

	return p.doSendEmail(ctx, account, emailMsg)
}

// doSendEmail performs the actual email sending
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
	client, err := mail.NewClient(account.Host,
		mail.WithPort(account.Port),
		mail.WithUsername(account.Username),
		mail.WithPassword(account.Password),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer client.Close()

	client.SetSSL(true)
	if err := client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
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

// GetType returns the subprovider type of this account
func (a *Account) GetType() string {
	return a.Type
}
