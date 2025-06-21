package wecombot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider supports multiple bots and strategy selection
type Provider struct {
	bots     []*Bot
	selector *utils.Selector[*Bot]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new wecombot provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("wecombot provider is not configured or is disabled")
	}

	// Convert to pointer slice
	bots := make([]*Bot, len(config.Bots))
	for i := range config.Bots {
		bots[i] = &config.Bots[i]
	}

	// Use common initialization logic
	enabledBots, selector, err := utils.InitProvider(&config, bots)
	if err != nil {
		return nil, errors.New("no enabled wecombot bots found")
	}

	return &Provider{
		bots:     enabledBots,
		selector: selector,
	}, nil
}

// Send sends message, automatically selects bot
func (p *Provider) Send(ctx context.Context, message core.Message) error {
	wecomMsg, ok := message.(Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected wecombot.Message interface, got %T", message))
	}
	if err := wecomMsg.Validate(); err != nil {
		return err
	}

	selectedBot := p.selector.Select(ctx)
	if selectedBot == nil {
		return errors.New("no available bot")
	}
	return p.doSendWecom(ctx, selectedBot, wecomMsg)
}

func (p *Provider) doSendWecom(ctx context.Context, bot *Bot, message core.Message) error {
	endpoint := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", bot.Key)
	reqBody, err := json.Marshal(message)
	if err != nil {
		return err
	}
	body, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method: http.MethodPost,
		Body:   reqBody,
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
		return errors.New(reply.ErrMsg)
	}
	return nil
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return string(core.ProviderTypeWecombot)
}

// UploadMedia uploads a file to WeCom and returns its media_id and the bot used.
// The media_id is valid for only three days and can only be used by the bot that uploaded it.
// If you want to specify a bot, use context.WithValue(ctx, core.ItemNameKey, botName) or core.WithItemName(ctx, botName).
func (p *Provider) UploadMedia(ctx context.Context, filepath string, bodyReader io.Reader) (mediaId string, bot *Bot, err error) {
	selectedBot := p.selector.Select(ctx)
	if selectedBot == nil {
		return "", nil, errors.New("no available bot for media upload")
	}
	endpoint := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file", selectedBot.Key)

	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		BodyReader:  bodyReader,
		FilePath:    filepath,
		ContentType: "application/octet-stream",
	})

	if err != nil {
		return "", selectedBot, fmt.Errorf("failed to upload media: %w", err)
	}

	var reply struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		Type      string `json:"type"`
		MediaId   string `json:"media_id"`
		CreatedAt string `json:"created_at"`
	}

	if err = json.Unmarshal(resp, &reply); err != nil {
		return "", selectedBot, fmt.Errorf("failed to unmarshal WeCom media upload response: %w\nraw response: %s", err, resp)
	}

	if reply.ErrCode != 0 {
		return "", selectedBot, fmt.Errorf("WeCom media upload API error [%d]: %s", reply.ErrCode, reply.ErrMsg)
	}
	return reply.MediaId, selectedBot, nil
}
