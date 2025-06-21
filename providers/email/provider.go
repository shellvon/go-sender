package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
	"gopkg.in/gomail.v2"
)

// Provider supports multiple accounts and strategy selection
type Provider struct {
	accounts []*Account
	selector *utils.Selector[*Account]
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
	enabledAccounts, selector, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled email accounts found")
	}

	return &Provider{
		accounts: enabledAccounts,
		selector: selector,
	}, nil
}

// Send sends an email message
func (p *Provider) Send(ctx context.Context, message core.Message) error {
	emailMsg, ok := message.(*Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected *email.Message, got %T", message))
	}

	if err := emailMsg.Validate(); err != nil {
		return err
	}

	account := p.selector.Select(ctx)
	if account == nil {
		return errors.New("no available account")
	}
	return p.doSendEmail(ctx, account, emailMsg)
}


// doSendEmail performs the actual email sending
func (p *Provider) doSendEmail(ctx context.Context, account *Account, emailMsg *Message) error {
	if emailMsg.From == "" {
		emailMsg.From = account.From
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailMsg.From)
	m.SetHeader("To", emailMsg.To...)
	if len(emailMsg.Cc) > 0 {
		m.SetHeader("Cc", emailMsg.Cc...)
	}
	if len(emailMsg.Bcc) > 0 {
		m.SetHeader("Bcc", emailMsg.Bcc...)
	}
	m.SetHeader("Subject", emailMsg.Subject)
	if emailMsg.IsHTML {
		m.SetBody("text/html", emailMsg.Body)
	} else {
		m.SetBody("text/plain", emailMsg.Body)
	}
	for _, att := range emailMsg.Attachments {
		m.Attach(att)
	}
	d := gomail.NewDialer(account.Host, account.Port, account.Username, account.Password)
	return d.DialAndSend(m)
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
