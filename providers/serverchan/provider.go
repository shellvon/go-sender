package serverchan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the ServerChan provider
type Provider struct {
	accounts []*core.Account
	selector *utils.Selector[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new ServerChan provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("serverchan provider is not configured or is disabled")
	}

	// Convert to pointer slice
	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	// Use common initialization logic
	enabledAccounts, selector, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled serverchan accounts found")
	}

	return &Provider{
		accounts: enabledAccounts,
		selector: selector,
	}, nil
}

func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	scMsg, ok := msg.(*Message)
	if !ok {
		return fmt.Errorf("unsupported message type for serverchan provider: %T", msg)
	}

	selectedAccount := p.selector.Select(ctx)
	if selectedAccount == nil {
		return errors.New("no available account")
	}

	apiURL := p.buildAPIURL(selectedAccount.Key)

	// Directly marshal the message object to JSON
	jsonBody, err := json.Marshal(scMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	body, statusCode, err := utils.DoRequest(ctx, apiURL, utils.RequestOptions{
		Method:      "POST",
		Body:        jsonBody,
		ContentType: "application/json",
	})
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Check response status
	if statusCode != 200 {
		return fmt.Errorf("serverchan API returned non-OK status: %d", statusCode)
	}

	// Turbo version returns JSON
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("serverchan error: code=%d, msg=%s", result.Code, result.Message)
	}

	return nil
}

// buildAPIURL builds the appropriate API URL based on the SendKey format
// https://github.com/easychen/serverchan-demo/blob/481eccec5c3a6da71d6e744207a145d31a92b5a5/golang/main.go#L57C2-L69C3
func (p *Provider) buildAPIURL(key string) string {
	// Check if key starts with 'sctp' (enterprise/professional version)
	if strings.HasPrefix(key, "sctp") {
		// Extract number from sctp{num}t format
		re := regexp.MustCompile(`^sctp(\d+)t`)
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 {
			num := matches[1]
			return fmt.Sprintf("https://%s.push.ft07.com/send/%s.send", num, key)
		}
	}
	// Standard Turbo version
	return fmt.Sprintf("https://sctapi.ftqq.com/%s.send", key)
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeServerChan)
}
