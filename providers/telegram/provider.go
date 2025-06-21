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

// Provider represents a Telegram bot provider
type Provider struct {
	bots     []*Bot
	selector *utils.Selector[*Bot]
}

var (
	_ core.Provider = (*Provider)(nil)
)

// New creates a new Telegram provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("invalid telegram configuration: no bots configured or provider disabled")
	}

	// Convert to pointer slice
	bots := make([]*Bot, len(config.Bots))
	for i := range config.Bots {
		bots[i] = &config.Bots[i]
	}

	// Use common initialization logic
	enabledBots, selector, err := utils.InitProvider(&config, bots)
	if err != nil {
		return nil, errors.New("no enabled bots")
	}

	return &Provider{
		bots:     enabledBots,
		selector: selector,
	}, nil
}

// Send sends a Telegram message
func (p *Provider) Send(ctx context.Context, message core.Message) error {
	tmsg, ok := message.(Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected telegram.Message interface, got %T", message))
	}
	if err := tmsg.Validate(); err != nil {
		return err
	}
	bot := p.selector.Select(ctx)
	if bot == nil {
		return errors.New("no available bot")
	}
	return p.doSend(ctx, bot, tmsg)
}

// doSend performs the actual Telegram API request
func (p *Provider) doSend(ctx context.Context, bot *Bot, msg Message) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/send%s", bot.Token, capitalize(string(msg.GetMsgType())))

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
	return strings.ToUpper(string(s[0])) + s[1:]
}
