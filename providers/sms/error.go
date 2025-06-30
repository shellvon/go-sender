package sms

import "fmt"

// Error represents a unified SMS error type.
// code=0 means success
// provider: "cl253"/"aliyun"/"aliyun_globe"/"luosimao"/etc
// message: error description
// code: platform error code
type Error struct {
	Code     string
	Message  string
	Provider string
}

// NewUnsupportedMessageTypeError creates an error for unsupported message types.
func NewUnsupportedMessageTypeError(provider string, messageType, category string) *Error {
	return &Error{
		Code: ErrorCodeUnsupportedMessageType,
		Message: fmt.Sprintf(
			"Provider %s does not support %s messages with category %s",
			provider,
			messageType,
			category,
		),
		Provider: provider,
	}
}

// NewUnsupportedInternationalError creates an error for unsupported international numbers.
func NewUnsupportedInternationalError(provider string, feature string) *Error {
	return &Error{
		Code:     ErrorCodeUnsupportedInternational,
		Message:  fmt.Sprintf("Provider %s does not support international numbers for %s", provider, feature),
		Provider: provider,
	}
}

// Error implements the error interface for Error.
func (e *Error) Error() string {
	return fmt.Sprintf("Error(provider=%s, code=%s, msg=%s)", e.Provider, e.Code, e.Message)
}

// IsRetryable returns whether the error is retryable.
func (e *Error) IsRetryable() bool { return false }

// 通用错误码常量.
const (
	ErrorCodeUnsupportedMessageType   = "UNSUPPORTED_MESSAGE_TYPE"
	ErrorCodeUnsupportedInternational = "UNSUPPORTED_INTERNATIONAL"
)
