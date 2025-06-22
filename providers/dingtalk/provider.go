package dingtalk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the DingTalk provider
type Provider struct {
	accounts []*core.Account
	selector *utils.Selector[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new DingTalk provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("dingtalk provider is not configured or is disabled")
	}

	// Convert to pointer slice
	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	// Use common initialization logic
	enabledAccounts, selector, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled dingtalk accounts found")
	}

	return &Provider{
		accounts: enabledAccounts,
		selector: selector,
	}, nil
}

func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	dingMsg, ok := msg.(Message)
	if !ok {
		return fmt.Errorf("unsupported message type for dingtalk provider: %T", msg)
	}

	selectedAccount := p.selector.Select(ctx)
	if selectedAccount == nil {
		return errors.New("no available account")
	}

	return p.doSendDingtalk(ctx, selectedAccount, dingMsg)
}

// doSendDingtalk sends a message using the specified account
func (p *Provider) doSendDingtalk(ctx context.Context, account *core.Account, message Message) error {
	// Build webhook URL
	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", account.Key)

	// Marshal message to JSON
	jsonBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	// Send request
	body, statusCode, err := utils.DoRequest(ctx, webhookURL, utils.RequestOptions{
		Method:      "POST",
		Body:        jsonBody,
		ContentType: "application/json",
	})
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Check response
	if statusCode != 200 {
		return fmt.Errorf("dingtalk API returned non-OK status: %d", statusCode)
	}

	// Parse response
	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("dingtalk error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
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
