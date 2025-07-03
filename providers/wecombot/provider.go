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
	"os"
	"path/filepath"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// uploadTarget abstracts Voice/File messages that need auto media upload.
type uploadTarget interface {
	getLocalPath() string
	getMediaID() string
	setMediaID(string)
	mediaType() string // "voice" or "file"
}

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
func (p *Provider) Send(ctx context.Context, msg core.Message, opts *core.ProviderSendOptions) error {
	if up, ok := msg.(uploadTarget); ok {
		if up.getMediaID() == "" && up.getLocalPath() != "" {
			file, err := os.Open(up.getLocalPath())
			if err != nil {
				return err
			}
			defer file.Close()
			mediaID, acc, err := p.uploadMediaType(ctx, up.getLocalPath(), up.mediaType(), file, opts.HTTPClient)
			if err != nil {
				return err
			}
			up.setMediaID(mediaID)
			ctx = core.WithCtxItemName(ctx, acc.GetName())
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
) (string, *core.Account, error) {
	selectedAccount := p.SelectConfig(ctx)
	if selectedAccount == nil {
		return "", nil, errors.New("no available account")
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
