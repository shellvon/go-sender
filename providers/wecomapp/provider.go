package wecomapp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// Config 代表企业微信应用provider的配置结构
type Config = core.BaseConfig[*Account]

// Provider 实现企业微信应用provider
type Provider struct {
	*providers.HTTPProvider[*Account]
	transformer *WecomAppTransformer
}

var _ core.Provider = (*Provider)(nil)

// uploadTarget 抽象需要自动上传媒体文件的消息接口
type uploadTarget interface {
	getLocalPath() string
	getMediaID() string
	setMediaID(string)
	mediaType() string // "image", "voice", "video", or "file"
}

// ProviderOption 代表修改企业微信应用Provider配置的函数
type ProviderOption func(*Config)

// New 创建一个使用WecomAppTransformer的企业微信应用provider实例
func New(config *Config, tokenCache TokenCache) (*Provider, error) {
	// 创建transformer，传入用户设置的tokenCache
	wecomTransformer := NewWecomAppTransformer(tokenCache)

	// 创建HTTPProvider时使用transformer
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeWecomApp),
		wecomTransformer,
		config,
	)
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		HTTPProvider: httpProvider,
		transformer:  wecomTransformer,
	}

	return provider, nil
}

// Name 返回provider名称
func (p *Provider) Name() string {
	return string(core.ProviderTypeWecomApp)
}

// Send 重写了内嵌HTTPProvider.Send方法，用于处理访问令牌管理
// 和带有本地文件路径的媒体消息的自动上传功能。
func (p *Provider) Send(
	ctx context.Context,
	msg core.Message,
	opts *core.ProviderSendOptions,
) (*core.SendResult, error) {
	if opts == nil {
		opts = &core.ProviderSendOptions{}
	}

	selectedAccount, err := p.HTTPProvider.Select(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 设置AgentID
	if setter, ok := msg.(agentIDSetter); ok {
		setter.setAgentID(selectedAccount.AgentID())
	}

	// 设置路由信息
	ri := core.GetRoute(ctx)
	if ri == nil {
		ri = &core.RouteInfo{}
	}
	ri.AccountName = selectedAccount.GetName()
	ctx = core.WithRoute(ctx, ri)

	// 在发送之前，手动处理需要上传的媒体文件
	if err := p.handleMediaUpload(ctx, msg, selectedAccount, opts); err != nil {
		return nil, err
	}

	// 获取access token（使用opts中的HTTP client）
	var httpClient *http.Client
	if opts != nil {
		httpClient = opts.HTTPClient
	}

	accessToken, err := p.transformer.GetValidAccessToken(ctx, selectedAccount, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 调用基础transformer构建请求
	reqSpec, handler, err := p.transformer.Transform(ctx, msg, selectedAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to transform message: %w", err)
	}

	// 手动添加access_token到查询参数
	if reqSpec.QueryParams == nil {
		reqSpec.QueryParams = make(url.Values)
	}
	reqSpec.QueryParams.Set("access_token", accessToken)

	// 包装handler以处理企业微信错误
	handler = p.transformer.wrapHandler(ctx, selectedAccount, handler)

	// 直接使用父类的公开方法发送HTTP请求
	return p.ExecuteHTTPRequest(ctx, reqSpec, handler, opts)
}

// handleMediaUpload 手动处理媒体文件上传
func (p *Provider) handleMediaUpload(
	ctx context.Context,
	msg core.Message,
	account *Account,
	opts *core.ProviderSendOptions,
) error {
	// 检查消息是否需要上传媒体文件
	uploadMsg, needsUpload := msg.(uploadTarget)
	if !needsUpload {
		return nil // 不需要上传，直接返回
	}

	// 如果已经有mediaID，不需要重复上传
	if uploadMsg.getMediaID() != "" {
		return nil
	}

	// 如果没有本地文件路径，也不需要上传
	localPath := uploadMsg.getLocalPath()
	if localPath == "" {
		return nil
	}

	// 使用transformer上传媒体文件
	mediaID, err := p.transformer.UploadMediaWithClient(
		ctx,
		account,
		localPath,
		uploadMsg.mediaType(),
		opts.HTTPClient,
	)
	if err != nil {
		return fmt.Errorf("failed to upload media: %w", err)
	}

	// 设置上传后的mediaID
	uploadMsg.setMediaID(mediaID)
	return nil
}

// NewProvider 使用给定的账号和选项创建新的企业微信应用provider
//
// 至少需要一个账号。
//
// 示例:
//
//	provider, err := wecomapp.NewProvider([]*wecomapp.Account{account1, account2},
//	    wecomapp.Strategy(core.StrategyWeighted))
func NewProvider(accounts []*Account, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		accounts,
		core.ProviderTypeWecomApp,
		func(meta core.ProviderMeta, items []*Account) *Config {
			return &Config{
				ProviderMeta: meta,
				Items:        items,
			}
		},
		func(config *Config) (*Provider, error) {
			return New(config, nil)
		},
		opts...,
	)
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: wecomapp.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*wecomapp.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)
