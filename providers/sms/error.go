package sms

import "fmt"

// SMSError 统一短信错误类型
// code=0 表示成功
// provider: "cl253"/"aliyun"/"aliyun_globe"/"luosimao"/等
// message: 错误描述
// code: 平台返回码
type SMSError struct {
	Code     string
	Message  string
	Provider string
}

func (e *SMSError) Error() string {
	return fmt.Sprintf("SMSError(provider=%s, code=%s, msg=%s)", e.Provider, e.Code, e.Message)
}
func (e *SMSError) IsRetryable() bool { return false }

// 通用错误码常量
const (
	ErrorCodeUnsupportedMessageType   = "UNSUPPORTED_MESSAGE_TYPE"
	ErrorCodeUnsupportedInternational = "UNSUPPORTED_INTERNATIONAL"
)

// NewUnsupportedMessageTypeError 创建不支持的消息类型错误
func NewUnsupportedMessageTypeError(provider string, messageType, category string) *SMSError {
	return &SMSError{
		Code:     ErrorCodeUnsupportedMessageType,
		Message:  fmt.Sprintf("Provider %s does not support %s messages with category %s", provider, messageType, category),
		Provider: provider,
	}
}

// NewUnsupportedInternationalError 创建不支持国际号码的错误
func NewUnsupportedInternationalError(provider string, feature string) *SMSError {
	return &SMSError{
		Code:     ErrorCodeUnsupportedInternational,
		Message:  fmt.Sprintf("Provider %s does not support international numbers for %s", provider, feature),
		Provider: provider,
	}
}
