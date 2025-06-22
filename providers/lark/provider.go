package lark

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the Lark/Feishu Robot provider
type Provider struct {
	config *Config
	client *http.Client
}

// NewProvider creates a new Lark provider
func NewProvider(config *Config) *Provider {
	return &Provider{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Send sends a message using the Lark provider
func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	larkMsg, ok := msg.(Message)
	if !ok {
		return fmt.Errorf("unsupported message type for Lark provider: %T", msg)
	}

	// Select a bot based on strategy
	bot, err := p.selectBot(ctx)
	if err != nil {
		return fmt.Errorf("failed to select bot: %w", err)
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"msg_type": larkMsg.GetMsgType(),
		"content":  larkMsg,
	}

	// Add timestamp and sign if secret is provided
	if bot.Secret != "" {
		timestamp := time.Now().Unix()
		sign := p.generateSign(timestamp, bot.Secret)

		payload["timestamp"] = timestamp
		payload["sign"] = sign
	}

	// Marshal the payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", bot.Webhook, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lark API returned non-OK status: %d", resp.StatusCode)
	}

	// Parse response to check for errors
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Code != 0 {
		return fmt.Errorf("lark API error: code=%d, msg=%s", response.Code, response.Msg)
	}

	return nil
}

// selectBot selects a bot based on the configured strategy
func (p *Provider) selectBot(ctx context.Context) (*Bot, error) {
	enabledBots := make([]*Bot, 0)
	for i := range p.config.Bots {
		if p.config.Bots[i].IsEnabled() {
			enabledBots = append(enabledBots, &p.config.Bots[i])
		}
	}

	if len(enabledBots) == 0 {
		return nil, fmt.Errorf("no enabled bots found")
	}

	// Use the selector to choose a bot
	strategy := utils.GetStrategy(p.config.GetStrategy())
	selector := utils.NewSelector(enabledBots, strategy)
	selected := selector.Select(ctx)

	if selected == nil {
		return nil, fmt.Errorf("no bot selected")
	}

	return selected, nil
}

// generateSign generates the signature for Lark webhook
func (p *Provider) generateSign(timestamp int64, secret string) string {
	// Lark signature format: timestamp + "\n" + secret
	stringToSign := strconv.FormatInt(timestamp, 10) + "\n" + secret

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GetConfig returns the provider configuration
func (p *Provider) GetConfig() interface{} {
	return p.config
}

// IsConfigured checks if the provider is properly configured
func (p *Provider) IsConfigured() bool {
	return p.config.IsConfigured()
}
