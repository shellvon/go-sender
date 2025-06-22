package telegram

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the Telegram provider
type Provider struct {
	accounts []*core.Account
	selector *utils.Selector[*core.Account]
}

var (
	_ core.Provider = (*Provider)(nil)
)

// New creates a new Telegram provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("telegram provider is not configured or is disabled")
	}

	// Convert to pointer slice
	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	// Use common initialization logic
	enabledAccounts, selector, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled telegram accounts found")
	}

	return &Provider{
		accounts: enabledAccounts,
		selector: selector,
	}, nil
}

// Send sends a Telegram message
func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	tgMsg, ok := msg.(Message)
	if !ok {
		return fmt.Errorf("unsupported message type for telegram provider: %T", msg)
	}

	selectedAccount := p.selector.Select(ctx)
	if selectedAccount == nil {
		return errors.New("no available account")
	}

	return p.doSend(ctx, selectedAccount, tgMsg)
}

// doSend sends a message using the specified account
func (p *Provider) doSend(ctx context.Context, account *core.Account, msg Message) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/send%s", account.Key, capitalize(string(msg.GetMsgType())))

	// Convert message to form data
	data := url.Values{}
	switch m := msg.(type) {
	case *TextMessage:
		data.Set("chat_id", m.ChatID)
		data.Set("text", m.Text)
		if m.ParseMode != "" {
			data.Set("parse_mode", m.ParseMode)
		}
		if m.DisableWebPreview {
			data.Set("disable_web_page_preview", "true")
		}
		if m.Silent {
			data.Set("disable_notification", "true")
		}
		if m.ProtectContent {
			data.Set("protect_content", "true")
		}
		if m.ReplyToMsg > 0 {
			data.Set("reply_to_message_id", strconv.Itoa(m.ReplyToMsg))
		}
	case *PhotoMessage:
		data.Set("chat_id", m.ChatID)
		data.Set("photo", m.Photo)
		if m.Caption != "" {
			data.Set("caption", m.Caption)
			if m.ParseMode != "" {
				data.Set("parse_mode", m.ParseMode)
			}
		}
		if m.Silent {
			data.Set("disable_notification", "true")
		}
		if m.ProtectContent {
			data.Set("protect_content", "true")
		}
		if m.ReplyToMsg > 0 {
			data.Set("reply_to_message_id", strconv.Itoa(m.ReplyToMsg))
		}
	case *DocumentMessage:
		data.Set("chat_id", m.ChatID)
		data.Set("document", m.Document)
		if m.Caption != "" {
			data.Set("caption", m.Caption)
			if m.ParseMode != "" {
				data.Set("parse_mode", m.ParseMode)
			}
		}
		if m.Silent {
			data.Set("disable_notification", "true")
		}
		if m.ProtectContent {
			data.Set("protect_content", "true")
		}
		if m.ReplyToMsg > 0 {
			data.Set("reply_to_message_id", strconv.Itoa(m.ReplyToMsg))
		}
	case *LocationMessage:
		data.Set("chat_id", m.ChatID)
		data.Set("latitude", strconv.FormatFloat(m.Latitude, 'f', -1, 64))
		data.Set("longitude", strconv.FormatFloat(m.Longitude, 'f', -1, 64))
		if m.Silent {
			data.Set("disable_notification", "true")
		}
		if m.ProtectContent {
			data.Set("protect_content", "true")
		}
		if m.ReplyToMsg > 0 {
			data.Set("reply_to_message_id", strconv.Itoa(m.ReplyToMsg))
		}
	case *ContactMessage:
		data.Set("chat_id", m.ChatID)
		data.Set("phone_number", m.PhoneNumber)
		data.Set("first_name", m.FirstName)
		if m.LastName != "" {
			data.Set("last_name", m.LastName)
		}
		if m.VCard != "" {
			data.Set("vcard", m.VCard)
		}
		if m.Silent {
			data.Set("disable_notification", "true")
		}
		if m.ProtectContent {
			data.Set("protect_content", "true")
		}
		if m.ReplyToMsg > 0 {
			data.Set("reply_to_message_id", strconv.Itoa(m.ReplyToMsg))
		}
	case *PollMessage:
		data.Set("chat_id", m.ChatID)
		data.Set("question", m.Question)
		data.Set("options", strings.Join(m.Options, "\n"))
		if m.IsAnonymous {
			data.Set("is_anonymous", "true")
		}
		if m.Type != "" {
			data.Set("type", m.Type)
		}
		if m.AllowsMultipleAnswers {
			data.Set("allows_multiple_answers", "true")
		}
		if m.Silent {
			data.Set("disable_notification", "true")
		}
		if m.ProtectContent {
			data.Set("protect_content", "true")
		}
		if m.ReplyToMsg > 0 {
			data.Set("reply_to_message_id", strconv.Itoa(m.ReplyToMsg))
		}
	default:
		return fmt.Errorf("unsupported message type: %T", msg)
	}

	_, _, err := utils.DoRequest(ctx, apiURL, utils.RequestOptions{
		Method:      "POST",
		Body:        []byte(data.Encode()),
		ContentType: "application/x-www-form-urlencoded",
	})
	return err
}

// Name returns the provider name
func (p *Provider) Name() string {
	return string(core.ProviderTypeTelegram)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
