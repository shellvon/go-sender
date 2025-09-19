package wecomapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
	"github.com/shellvon/go-sender/utils"
)

// AccessToken 代表带有过期时间的企业微信访问令牌.
type AccessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// 企业微信应用API端点.
const (
	// SendMessageEndpoint 发送消息的API端点.
	SendMessageEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	// UploadMediaEndpoint 上传媒体文件的API端点.
	UploadMediaEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	// GetTokenEndpoint 获取访问令牌的API端点.
	GetTokenEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/gettoken" //nolint:gosec // G101: This is WeChat Work API endpoint URL, not credentials
)

// transform 为企业微信应用消息构建HTTPRequestSpec.
func (t *wecomTransformer) transform(
	ctx context.Context,
	msg Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// 直接序列化消息以保留所有字段
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal wecomapp message: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      SendMessageEndpoint,
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

// SendRequest 代表通过企业微信应用API发送消息的请求结构
// 它使用json.RawMessage在转换过程中保留原始消息结构.
type SendRequest struct {
	ToUser  string `json:"touser,omitempty"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	MsgType string `json:"msgtype"`
	AgentID string `json:"agentid"`
	Safe    int    `json:"safe,omitempty"`

	// Use json.RawMessage to preserve the original message content
	Content json.RawMessage `json:"-"`
}

// wecomTransformer 具备完整企业微信API能力的transformer.
type wecomTransformer struct {
	*transformer.BaseHTTPTransformer[Message, *Account]

	tokenCache core.Cache[*AccessToken]
}

// newWecomTransformer 创建企业微信应用transformer.
func newWecomTransformer(tokenCache core.Cache[*AccessToken]) *wecomTransformer {
	if tokenCache == nil {
		tokenCache = core.NewMemoryCache[*AccessToken]()
	}

	respCfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeJSON,
		CheckBody: true,
		Path:      "errcode",
		Expect:    "0",
		Mode:      core.MatchEq,
		CodePath:  "errcode",
		MsgPath:   "errmsg",
	}

	t := &wecomTransformer{
		tokenCache: tokenCache,
	}
	t.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeWecomApp,
		"",
		respCfg,
		t.transform,
	)

	return t
}

// UploadMediaWithClient 使用指定的HTTP客户端上传媒体文件.
func (t *wecomTransformer) UploadMediaWithClient(
	ctx context.Context,
	account *Account,
	localPath, mediaType string,
	httpClient *http.Client,
) (string, error) {
	// 获取access token（使用传入的HTTP client）
	accessToken, err := t.GetValidAccessToken(ctx, account, httpClient)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	// 构建查询参数
	query := url.Values{}
	query.Set("access_token", accessToken)
	query.Set("type", mediaType)

	// 使用utils包的文件上传功能
	options := utils.HTTPRequestOptions{
		Method: http.MethodPost,
		Query:  query,
		Files:  map[string]string{"media": localPath},
		Client: httpClient,
	}

	resp, err := utils.SendRequest(ctx, UploadMediaEndpoint, options)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析响应
	var uploadResp MediaUploadResponse
	if parseErr := json.Unmarshal(bodyBytes, &uploadResp); parseErr != nil {
		return "", fmt.Errorf("failed to parse upload response: %w", parseErr)
	}

	// 检查企业微信API错误
	if uploadResp.ErrCode != 0 {
		wecomErr := &WecomAPIError{
			ErrCode: uploadResp.ErrCode,
			ErrMsg:  uploadResp.ErrMsg,
			Raw:     bodyBytes,
		}
		// 如果是认证错误，清除缓存
		if wecomErr.IsAuthenticationError() {
			tokenKey := t.getTokenKey(account)
			_ = t.tokenCache.Delete(tokenKey) // 忽略删除缓存的错误
		}
		return "", wecomErr
	}

	return uploadResp.MediaID, nil
}

// GetValidAccessToken 从缓存获取或刷新access token（公开方法供Provider调用）.
func (t *wecomTransformer) GetValidAccessToken(
	ctx context.Context,
	account *Account,
	httpClient *http.Client,
) (string, error) {
	tokenKey := t.getTokenKey(account)

	// 尝试从缓存获取
	token, exists := t.tokenCache.Get(tokenKey)
	if exists && token != nil {
		// 检查是否还有效（缓存会自动处理过期）
		return token.Token, nil
	}

	// 获取新token
	newToken, err := t.fetchAccessToken(ctx, account, httpClient)
	if err != nil {
		return "", err
	}

	// 直接缓存新token（TTL已经在fetchAccessToken中调整过了）
	ttl := time.Until(newToken.ExpiresAt)
	if ttl > 0 {
		_ = t.tokenCache.Set(tokenKey, newToken, &ttl)
	}

	return newToken.Token, nil
}

// fetchAccessToken 从API获取新的access token.
func (t *wecomTransformer) fetchAccessToken(
	ctx context.Context,
	account *Account,
	httpClient *http.Client,
) (*AccessToken, error) {
	// 构建查询参数
	query := url.Values{}
	query.Set("corpid", account.CorpID())
	query.Set("corpsecret", account.CorpSecret())

	options := utils.HTTPRequestOptions{
		Method: http.MethodGet,
		Query:  query,
		Client: httpClient,
	}

	resp, err := utils.SendRequest(ctx, GetTokenEndpoint, options)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp TokenResponse
	if parseErr := json.Unmarshal(bodyBytes, &tokenResp); parseErr != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", parseErr)
	}

	if tokenResp.ErrCode != 0 {
		return nil, &WecomAPIError{
			ErrCode: tokenResp.ErrCode,
			ErrMsg:  tokenResp.ErrMsg,
			Raw:     bodyBytes,
		}
	}

	// 创建带过期时间的访问令牌，提前5分钟过期避免临界条件
	const tokenExpiryBuffer = 5 * time.Minute
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn)*time.Second - tokenExpiryBuffer)
	return &AccessToken{
		Token:     tokenResp.AccessToken,
		ExpiresAt: expiresAt,
	}, nil
}

// getTokenKey 根据账号凭据生成访问令牌的缓存键.
func (t *wecomTransformer) getTokenKey(account *Account) string {
	return fmt.Sprintf("wecomapp:%s:%s", account.CorpID(), account.AgentID())
}

// wrapHandler 包装响应处理器以支持认证错误处理.
func (t *wecomTransformer) wrapHandler(
	_ context.Context,
	account *Account,
	originalHandler core.SendResultHandler,
) core.SendResultHandler {
	return func(result *core.SendResult) error {
		// 尝试解析为企业微信错误
		var wecomResp struct {
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}

		if err := json.Unmarshal(result.Body, &wecomResp); err == nil {
			if wecomResp.ErrCode != 0 {
				wecomErr := &WecomAPIError{
					ErrCode: wecomResp.ErrCode,
					ErrMsg:  wecomResp.ErrMsg,
					Raw:     result.Body,
				}
				// 如果是认证错误，清除缓存
				if wecomErr.IsAuthenticationError() {
					tokenKey := t.getTokenKey(account)
					_ = t.tokenCache.Delete(tokenKey)
				}
				return wecomErr
			}
		}

		// 使用原有的handler处理成功响应
		if originalHandler != nil {
			return originalHandler(result)
		}
		return nil
	}
}
