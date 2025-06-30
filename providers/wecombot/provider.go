package wecombot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the WeCom Bot provider.
type Provider struct {
	*providers.HTTPProvider[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new WeCom Bot provider instance.
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
	enabledAccounts, _, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled wecombot accounts found")
	}

	strategy := utils.GetStrategy(config.Strategy)

	// Create generic provider
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeWecombot),
		enabledAccounts,
		newWecombotTransformer(),
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// UploadMedia uploads a media file and returns the media_id
// This method can be used to upload images, files, etc. before sending them.
func (p *Provider) UploadMedia(
	ctx context.Context,
	filePath string,
	bodyReader io.Reader,
) (string, *core.Account, error) {
	selectedAccount := p.SelectConfig(ctx)
	if selectedAccount == nil {
		return "", nil, errors.New("no available account")
	}

	// Build upload URL
	uploadURL := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file",
		selectedAccount.Key,
	)

	// Create multipart form with custom reader
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file part
	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content
	if _, errCopy := io.Copy(part, bodyReader); errCopy != nil {
		return "", nil, fmt.Errorf("failed to copy file content: %w", errCopy)
	}

	// Close writer
	if errClose := writer.Close(); errClose != nil {
		return "", nil, fmt.Errorf("failed to close writer: %w", errClose)
	}

	respBody, statusCode, err := utils.DoRequest(ctx, uploadURL, utils.HTTPRequestOptions{
		Method:    "POST",
		RawReader: body,
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to upload media: %w", err)
	}

	// Check response
	if statusCode != http.StatusOK {
		return "", nil, fmt.Errorf("upload API returned non-OK status: %d", statusCode)
	}

	// Parse response
	var result struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		MediaID   string `json:"media_id"`
		Type      string `json:"type"`
		CreatedAt int64  `json:"created_at"`
	}

	if errUnmarshal := json.Unmarshal(respBody, &result); errUnmarshal != nil {
		return "", nil, fmt.Errorf("failed to parse upload response: %w", errUnmarshal)
	}

	if result.ErrCode != 0 {
		return "", nil, fmt.Errorf("upload error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	return result.MediaID, selectedAccount, nil
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeWecombot)
}
