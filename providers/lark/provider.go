package lark

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the Lark provider
type Provider struct {
	accounts []*core.Account
	selector *utils.Selector[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new Lark provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("lark provider is not configured or is disabled")
	}

	// Convert to pointer slice
	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	// Use common initialization logic
	enabledAccounts, selector, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled lark accounts found")
	}

	return &Provider{
		accounts: enabledAccounts,
		selector: selector,
	}, nil
}

// Send sends a message using the Lark provider
func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	larkMsg, ok := msg.(Message)
	if !ok {
		return fmt.Errorf("unsupported message type for lark provider: %T", msg)
	}

	selectedAccount := p.selector.Select(ctx)
	if selectedAccount == nil {
		return errors.New("no available account")
	}

	// Build webhook URL
	webhookURL := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", selectedAccount.Key)

	// Prepare the request payload
	payload := map[string]interface{}{
		"msg_type": larkMsg.GetMsgType(),
		"content":  larkMsg,
	}

	// Marshal message to JSON
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	// Send request
	body, statusCode, err := utils.DoRequest(ctx, webhookURL, utils.RequestOptions{
		Method: "POST",
		JSON:   jsonBody,
	})
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Check response
	if statusCode != 200 {
		return fmt.Errorf("lark API returned non-OK status: %d", statusCode)
	}

	// Parse response
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("lark error: code=%d, msg=%s", result.Code, result.Message)
	}

	return nil
}

// generateSign generates the signature for Lark webhook
func (p *Provider) generateSign(timestamp int64, secret string) string {
	// Lark signature format: timestamp + "\n" + secret
	stringToSign := strconv.FormatInt(timestamp, 10) + "\n" + secret

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return string(core.ProviderTypeLark)
}
