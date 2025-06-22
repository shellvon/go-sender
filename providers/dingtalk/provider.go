package dingtalk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider supports multiple bots and strategy selection
type Provider struct {
	bots     []*Bot
	selector *utils.Selector[*Bot]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new dingtalk provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("dingtalk provider is not configured or is disabled")
	}

	// Convert to pointer slice
	bots := make([]*Bot, len(config.Bots))
	for i := range config.Bots {
		bots[i] = &config.Bots[i]
	}

	// Use common initialization logic
	enabledBots, selector, err := utils.InitProvider(&config, bots)
	if err != nil {
		return nil, errors.New("no enabled dingtalk bots found")
	}

	return &Provider{
		bots:     enabledBots,
		selector: selector,
	}, nil
}

// Send sends message, automatically selects bot
func (p *Provider) Send(ctx context.Context, message core.Message) error {
	dingtalkMsg, ok := message.(Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected dingtalk.Message interface, got %T", message))
	}
	if err := dingtalkMsg.Validate(); err != nil {
		return err
	}

	selectedBot := p.selector.Select(ctx)
	if selectedBot == nil {
		return errors.New("no available bot")
	}
	return p.doSendDingtalk(ctx, selectedBot, dingtalkMsg)
}

func (p *Provider) doSendDingtalk(ctx context.Context, bot *Bot, message core.Message) error {
	// Build webhook URL with signature if secret is provided
	webhookURL := bot.Webhook
	if bot.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		sign := p.generateSignature(timestamp, bot.Secret)
		webhookURL = fmt.Sprintf("%s&timestamp=%s&sign=%s", bot.Webhook, timestamp, sign)
	}

	reqBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	body, _, err := utils.DoRequest(ctx, webhookURL, utils.RequestOptions{
		Method: http.MethodPost,
		Body:   reqBody,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	var reply = new(struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	})

	if err = json.Unmarshal(body, reply); err != nil {
		return fmt.Errorf("error response: %w\nraw response: %s", err, body)
	}

	if reply.ErrCode != 0 {
		return fmt.Errorf("dingtalk API error [%d]: %s", reply.ErrCode, reply.ErrMsg)
	}
	return nil
}

// generateSignature generates signature for DingTalk webhook
func (p *Provider) generateSignature(timestamp, secret string) string {
	stringToSign := timestamp + "\n" + secret
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return string(core.ProviderTypeDingtalk)
}
