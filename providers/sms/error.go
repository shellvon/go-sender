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
