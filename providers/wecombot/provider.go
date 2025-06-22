package wecombot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the WeCom Bot provider
type Provider struct {
	accounts []*core.Account
	selector *utils.Selector[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new WeCom Bot provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("wecombot provider is not configured or is disabled")
	}

	// Convert to pointer slice
	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	// Use common initialization logic
	enabledAccounts, selector, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled wecombot accounts found")
	}

	return &Provider{
		accounts: enabledAccounts,
		selector: selector,
	}, nil
}

func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	wecomMsg, ok := msg.(Message)
	if !ok {
		return fmt.Errorf("unsupported message type for wecombot provider: %T", msg)
	}

	selectedAccount := p.selector.Select(ctx)
	if selectedAccount == nil {
		return errors.New("no available account")
	}

	return p.doSendWecom(ctx, selectedAccount, wecomMsg)
}

// doSendWecom sends a message using the specified account
func (p *Provider) doSendWecom(ctx context.Context, account *core.Account, message Message) error {
	// Build webhook URL
	webhookURL := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", account.Key)

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
		return fmt.Errorf("wecom API returned non-OK status: %d", statusCode)
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
		return fmt.Errorf("wecom error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

// UploadMedia uploads a media file and returns the media_id
// This method can be used to upload images, files, etc. before sending them
func (p *Provider) UploadMedia(ctx context.Context, filePath string, bodyReader io.Reader) (mediaId string, account *core.Account, err error) {
	selectedAccount := p.selector.Select(ctx)
	if selectedAccount == nil {
		return "", nil, errors.New("no available account")
	}

	// Build upload URL
	uploadURL := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file", selectedAccount.Key)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file part
	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content
	if _, err := io.Copy(part, bodyReader); err != nil {
		return "", nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return "", nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Send request
	respBody, statusCode, err := utils.DoRequest(ctx, uploadURL, utils.RequestOptions{
		Method:      "POST",
		Body:        body.Bytes(),
		ContentType: writer.FormDataContentType(),
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to upload media: %w", err)
	}

	// Check response
	if statusCode != 200 {
		return "", nil, fmt.Errorf("upload API returned non-OK status: %d", statusCode)
	}

	// Parse response
	var result struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		MediaId   string `json:"media_id"`
		Type      string `json:"type"`
		CreatedAt int64  `json:"created_at"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", nil, fmt.Errorf("failed to parse upload response: %w", err)
	}

	if result.ErrCode != 0 {
		return "", nil, fmt.Errorf("upload error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	return result.MediaId, selectedAccount, nil
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeWecombot)
}
