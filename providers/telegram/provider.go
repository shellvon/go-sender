package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

// buildApiURL returns the full Telegram API URL for the given message type
func (p *Provider) buildApiURL(account *core.Account, msg Message) (string, error) {
	var endpoint string
	switch msg.(type) {
	case *TextMessage:
		endpoint = "sendMessage"
	case *PhotoMessage:
		endpoint = "sendPhoto"
	case *DocumentMessage:
		endpoint = "sendDocument"
	case *LocationMessage:
		endpoint = "sendLocation"
	case *ContactMessage:
		endpoint = "sendContact"
	case *PollMessage:
		endpoint = "sendPoll"
	case *AudioMessage:
		endpoint = "sendAudio"
	case *VideoMessage:
		endpoint = "sendVideo"
	case *AnimationMessage:
		endpoint = "sendAnimation"
	case *VoiceMessage:
		endpoint = "sendVoice"
	case *VideoNoteMessage:
		endpoint = "sendVideoNote"
	case *DiceMessage:
		endpoint = "sendDice"
	default:
		return "", fmt.Errorf("unsupported message type: %T", msg)
	}
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", account.Key, endpoint), nil
}

// doSend sends a message using the specified account
func (p *Provider) doSend(ctx context.Context, account *core.Account, msg Message) error {
	apiURL, err := p.buildApiURL(account, msg)
	if err != nil {
		return err
	}

	bodyBytes, _, err := utils.DoRequest(ctx, apiURL, utils.RequestOptions{
		Method: "POST",
		JSON:   msg,
	})
	if err != nil {
		return err
	}

	// Parse Telegram API response
	var response struct {
		OK          bool        `json:"ok"`
		ErrorCode   int         `json:"error_code,omitempty"`
		Description string      `json:"description,omitempty"`
		Result      interface{} `json:"result,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return fmt.Errorf("failed to parse Telegram API response: %w", err)
	}

	// Check if the request was successful
	if !response.OK {
		if response.ErrorCode != 0 {
			return fmt.Errorf("telegram API error %d: %s", response.ErrorCode, response.Description)
		}
		return fmt.Errorf("telegram API request failed: %s", response.Description)
	}

	return nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return string(core.ProviderTypeTelegram)
}
