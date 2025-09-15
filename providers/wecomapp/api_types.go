package wecomapp

import "fmt"

// WecomAPIError 企业微信API的结构化错误
type WecomAPIError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Raw     []byte `json:"-"` // 原始响应，便于调试
}

func (e *WecomAPIError) Error() string {
	return fmt.Sprintf("WeChat Work API Error %d: %s", e.ErrCode, e.ErrMsg)
}

// IsAuthenticationError 检查是否为认证相关错误
func (e *WecomAPIError) IsAuthenticationError() bool {
	authCodes := []int{
		40001, // 不合法的secret参数
		40014, // 不合法的access_token
		41001, // 缺少access_token参数
		42001, // access_token超时
		48001, // api接口未授权
	}
	for _, code := range authCodes {
		if e.ErrCode == code {
			return true
		}
	}
	return false
}

// IsRetryable 检查错误是否可重试
func (e *WecomAPIError) IsRetryable() bool {
	return e.IsAuthenticationError() || e.ErrCode == 45009 // 接口调用超过限制
}

// IsSuccess 检查响应是否成功
func (e *WecomAPIError) IsSuccess() bool {
	return e.ErrCode == 0
}

// TokenRequest 代表获取访问令牌的请求结构
type TokenRequest struct {
	CorpID     string `json:"corpid"`
	CorpSecret string `json:"corpsecret"`
}

// TokenResponse 代表令牌API的响应结构
type TokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// SendResponse 代表发送消息API的响应结构
type SendResponse struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	InvalidUser    string `json:"invaliduser,omitempty"`
	InvalidParty   string `json:"invalidparty,omitempty"`
	InvalidTag     string `json:"invalidtag,omitempty"`
	UnlicensedUser string `json:"unlicenseduser,omitempty"`
	MsgID          string `json:"msgid,omitempty"`
	ResponseCode   string `json:"response_code,omitempty"`
}

// MediaUploadResponse 代表媒体上传API的响应结构
type MediaUploadResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}
