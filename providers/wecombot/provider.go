// Package wecombot provides wecombot notification support for go-sender.
//
// This package implements the go-sender Provider interface for wecombot
// messaging service, enabling seamless integration with the go-sender
// notification system.
//
// Basic usage:
//
//	provider, err := wecombot.NewProvider([]*wecombot.Account{account})
//	msg := wecombot.Text().Content("Hello World").Build()
//	provider.Send(context.Background(), msg, nil)
package wecombot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

type Config = core.BaseConfig[*Account]

// uploadTarget abstracts Voice/File messages that need auto media upload.
type uploadTarget interface {
	getLocalPath() string
	getMediaID() string
	setMediaID(string)
	mediaType() string // "voice" or "file"
}

// Provider implements the WeCom Bot provider.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new WeCom Bot provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeWecombot),
		newWecombotTransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// ProviderOption represents a function that modifies WeCom Bot Provider configuration.
type ProviderOption func(*Config)

// NewProvider creates a new WeCom Bot provider with the given accounts and options.
//
// At least one account is required.
//
// Example:
//
//	provider, err := wecombot.NewProvider([]*wecombot.Account{account1, account2},
//	    wecombot.Strategy(core.StrategyWeighted))
func NewProvider(accounts []*Account, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		accounts,
		core.ProviderTypeWecombot,
		func(meta core.ProviderMeta, items []*Account) *Config {
			return &Config{
				ProviderMeta: meta,
				Items:        items,
			}
		},
		New,
		opts...,
	)
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeWecombot)
}

// Send overrides embedded HTTPProvider.Send to support automatic media upload for
// VoiceMessage and FileMessage when only a local file path is provided.
//
// Upload constraints:
//   - All media types must be larger than 5 bytes.
//   - Regular files (`file`): size must not exceed 20 MB.
//   - Voice files (`voice`): size must not exceed 2 MB, playback length ≤ 60 seconds, AMR format only.
//
// See https://developer.work.weixin.qq.com/document/path/91770 in the bottom of the page.
//
//nolint:nestif // acceptable nesting for media upload.
func (p *Provider) Send(
	ctx context.Context,
	msg core.Message,
	opts *core.ProviderSendOptions,
) (*core.SendResult, error) {
	if up, ok := msg.(uploadTarget); ok {
		if up.getMediaID() == "" && up.getLocalPath() != "" {
			file, err := os.Open(up.getLocalPath())
			if err != nil {
				return nil, err
			}
			defer file.Close()
			mediaID, acc, err := p.uploadMediaType(ctx, up.getLocalPath(), up.mediaType(), file, opts.HTTPClient)
			if err != nil {
				return nil, err
			}
			up.setMediaID(mediaID)
			// Preserve existing route info (e.g., strategy) while overriding account.
			ri := core.GetRoute(ctx)
			if ri == nil {
				ri = &core.RouteInfo{}
			}
			ri.AccountName = acc.GetName()
			ctx = core.WithRoute(ctx, ri)
		}
	}
	return p.HTTPProvider.Send(ctx, msg, opts)
}

// uploadMediaType is an internal helper allowing custom media type parameter.
func (p *Provider) uploadMediaType(
	ctx context.Context,
	filePath, mediaType string,
	reader io.Reader,
	httpClient *http.Client,
) (string, *Account, error) {
	selectedAccount, err := p.Select(ctx, nil)
	if err != nil {
		return "", nil, err
	}

	uploadURL := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=%s",
		selectedAccount.APIKey,
		mediaType,
	)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return "", nil, err
	}
	if _, errCopy := io.Copy(part, reader); errCopy != nil {
		return "", nil, errCopy
	}
	if errClose := writer.Close(); errClose != nil {
		return "", nil, errClose
	}

	resp, status, err := utils.DoRequest(ctx, uploadURL, utils.HTTPRequestOptions{
		Method:    http.MethodPost,
		RawReader: body,
		Headers:   map[string]string{"Content-Type": writer.FormDataContentType()},
		Client:    httpClient,
	})
	if err != nil {
		return "", nil, err
	}
	if !utils.IsAcceptableStatus(status) {
		return "", nil, fmt.Errorf("upload status %d", status)
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		MediaID string `json:"media_id"`
	}
	var parseErr error
	if parseErr = json.Unmarshal(resp, &result); parseErr != nil {
		return "", nil, parseErr
	}
	if result.ErrCode != 0 {
		return "", nil, fmt.Errorf("upload error %s", result.ErrMsg)
	}
	return result.MediaID, selectedAccount, nil
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: wecombot.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*wecombot.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)
