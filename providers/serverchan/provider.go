package serverchan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the Server酱 provider
type Provider struct {
	config *Config
	client *http.Client
}

func NewProvider(config *Config) *Provider {
	return &Provider{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	scMsg, ok := msg.(*Message)
	if !ok {
		return fmt.Errorf("unsupported message type for Server酱 provider: %T", msg)
	}

	account, err := p.selectAccount(ctx)
	if err != nil {
		return fmt.Errorf("failed to select account: %w", err)
	}

	apiURL := p.buildAPIURL(account.Key)
	form := url.Values{}
	form.Set("title", scMsg.Title)
	form.Set("desp", scMsg.Content)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

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

func (p *Provider) selectAccount(ctx context.Context) (*Account, error) {
	enabled := make([]*Account, 0)
	for i := range p.config.Accounts {
		if p.config.Accounts[i].IsEnabled() {
			enabled = append(enabled, &p.config.Accounts[i])
		}
	}
	if len(enabled) == 0 {
		return nil, fmt.Errorf("no enabled account found")
	}
	strategy := utils.GetStrategy(p.config.GetStrategy())
	selector := utils.NewSelector(enabled, strategy)
	selected := selector.Select(ctx)
	if selected == nil {
		return nil, fmt.Errorf("no account selected")
	}
	return selected, nil
}

func (p *Provider) GetConfig() interface{} {
	return p.config
}

func (p *Provider) IsConfigured() bool {
	return p.config.IsConfigured()
}

func (p *Provider) Name() string {
	return "serverchan"
}
