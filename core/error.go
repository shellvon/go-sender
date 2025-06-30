package core

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"syscall"
)

// RetryableError interface for retryable errors.
type RetryableError interface {
	error
	IsRetryable() bool
}

// ParamError represents a parameter error type.
type ParamError struct {
	message string
}

// NewParamError creates a new ParamError.
func NewParamError(message string) *ParamError {
	return &ParamError{message: message}
}

// Error implements the Error method of the error interface.
func (e *ParamError) Error() string {
	return fmt.Sprintf("parameter error: %s", e.message)
}

// IsRetryable implements the RetryableError interface, parameter errors are not retryable.
func (e *ParamError) IsRetryable() bool {
	return false
}

// NetworkError represents a network error.
type NetworkError struct {
	Err error
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

// IsRetryable returns whether the network error is retryable.
func (e NetworkError) IsRetryable() bool {
	return true
}

// TimeoutError represents a timeout error.
type TimeoutError struct {
	Err error
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("timeout error: %v", e.Err)
}

// IsRetryable returns whether the timeout error is retryable.
func (e TimeoutError) IsRetryable() bool {
	return true
}

// ValidationError represents a validation error (non-retryable).
type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", e.Err)
}

// IsRetryable returns whether the validation error is retryable.
func (e ValidationError) IsRetryable() bool {
	return false
}

// AuthenticationError represents an authentication error (non-retryable).
type AuthenticationError struct {
	Err error
}

func (e AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %v", e.Err)
}

// IsRetryable returns whether the authentication error is retryable.
func (e AuthenticationError) IsRetryable() bool {
	return false
}

// =============================================================================
// Smart Error Classifier
// =============================================================================

// ErrorClassifier defines the interface for determining if an error is retryable.
type ErrorClassifier interface {
	IsRetryableError(err error) bool
}

// defaultErrorClassifier is the default implementation of ErrorClassifier.
type defaultErrorClassifier struct {
	// Can add fields for custom rules if needed
}

// NewDefaultErrorClassifier creates a new default ErrorClassifier.
func NewDefaultErrorClassifier() ErrorClassifier {
	return &defaultErrorClassifier{}
}

// IsRetryableError determines if an error is retryable.
func (c *defaultErrorClassifier) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// 1. Errors implementing RetryableError interface
	var retryableErr RetryableError
	if errors.As(err, &retryableErr) {
		return retryableErr.IsRetryable()
	}

	// 2. Network-related errors
	if c.isNetworkError(err) {
		return true
	}

	// 3. Timeout errors
	if c.isTimeoutError(err) {
		return true
	}

	// 4. System call errors
	if c.isSystemCallError(err) {
		return true
	}

	// 5. HTTP 5xx errors (if error message contains them)
	if c.isHTTP5xxError(err) {
		return true
	}

	// 6. Type assertion errors and obvious program errors are not retryable
	if c.isTypeAssertionError(err) {
		return false
	}

	// 7. Context cancellation/timeout errors are not retryable
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// 8. JSON serialization/deserialization errors are not retryable
	if c.isJSONError(err) {
		return false
	}

	return false
}

// isNetworkError checks if the error is a network error.
func (c *defaultErrorClassifier) isNetworkError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return c.isNetworkError(urlErr.Err)
	}

	// Check for network-related keywords in error message
	errMsg := strings.ToLower(err.Error())
	networkKeywords := []string{
		"connection refused",
		"connection reset",
		"connection timeout",
		"network unreachable",
		"host unreachable",
		"no route to host",
		"broken pipe",
		"dns",
	}

	for _, keyword := range networkKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// isTimeoutError checks if the error is a timeout error.
func (c *defaultErrorClassifier) isTimeoutError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded")
}

// isSystemCallError checks if the error is a system call error.
func (c *defaultErrorClassifier) isSystemCallError(err error) bool {
	var syscallErr syscall.Errno
	if errors.As(err, &syscallErr) {
		// Retry only specific system call errors
		//nolint:exhaustive // intentionally not all cases handled, default covers the rest
		switch syscallErr {
		case syscall.ECONNREFUSED, syscall.ECONNRESET, syscall.EHOSTUNREACH, syscall.ENETUNREACH:
			return true
		default:
			return false
		}
	}
	return false
}

// isHTTP5xxError checks if the error is an HTTP 5xx error.
func (c *defaultErrorClassifier) isHTTP5xxError(err error) bool {
	errMsg := err.Error()
	// Check for 5xx status codes
	for code := 500; code < 600; code++ {
		if strings.Contains(errMsg, strconv.Itoa(code)) {
			return true
		}
	}
	return false
}

// isTypeAssertionError checks if the error is a type assertion error.
func (c *defaultErrorClassifier) isTypeAssertionError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	typeAssertionKeywords := []string{
		"type assertion failed",
		"interface conversion",
		"cannot convert",
		"invalid type assertion",
	}

	for _, keyword := range typeAssertionKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}
	return false
}

// isJSONError checks if the error is a JSON error.
func (c *defaultErrorClassifier) isJSONError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	jsonKeywords := []string{
		"invalid character",
		"unexpected end of json",
		"json:",
		"cannot unmarshal",
		"cannot marshal",
	}

	for _, keyword := range jsonKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}
	return false
}

// ErrorCode represents a specific error code for sender errors.
type ErrorCode int

const (
	// ErrCodeUnknown represents unknown error (0).
	ErrCodeUnknown ErrorCode = iota

	// ErrCodeInvalidConfig represents configuration errors (1000-1999).
	ErrCodeInvalidConfig ErrorCode = 1000 + iota
	ErrCodeProviderNotConfigured
	ErrCodeMissingRequiredField
	ErrCodeInvalidProviderType

	// ErrCodeProviderUnavailable represents provider errors (2000-2999).
	ErrCodeProviderUnavailable ErrorCode = 2000 + iota
	ErrCodeProviderSendFailed
	ErrCodeProviderTimeout
	ErrCodeProviderRateLimited

	// ErrCodeQueueFull represents queue errors (3000-3999).
	ErrCodeQueueFull ErrorCode = 3000 + iota
	ErrCodeQueueTimeout
	ErrCodeQueueSerializationFailed
	ErrCodeQueueDeserializationFailed

	// ErrCodeMaxRetriesExceeded represents retry errors (4000-4999).
	ErrCodeMaxRetriesExceeded ErrorCode = 4000 + iota
	ErrCodeRetryPolicyInvalid
	ErrCodeRetryFilterError

	// ErrCodeCircuitBreakerOpen represents circuit breaker errors (5000-5999).
	ErrCodeCircuitBreakerOpen ErrorCode = 5000 + iota
	ErrCodeCircuitBreakerTimeout

	// ErrCodeRateLimitExceeded represents rate limiter errors (6000-6999).
	ErrCodeRateLimitExceeded ErrorCode = 6000 + iota
	ErrCodeRateLimiterInvalid

	// ErrCodeMetricsCollectionFailed represents metrics errors (7000-7999).
	ErrCodeMetricsCollectionFailed ErrorCode = 7000 + iota

	// ErrCodeInternal represents general errors (9000-9999).
	ErrCodeInternal ErrorCode = 9000 + iota
	ErrCodeContextCancelled
	ErrCodeTimeout
	ErrCodeValidationFailed
)

// SenderError represents a structured error with code, message, and cause.
type SenderError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Cause   error     `json:"cause,omitempty"`
}

// NewSenderError creates a new SenderError with the given code, message, and cause.
func NewSenderError(code ErrorCode, message string, cause error) *SenderError {
	return &SenderError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewSenderErrorf creates a new SenderError with the given code and formatted message.
func NewSenderErrorf(code ErrorCode, format string, args ...interface{}) *SenderError {
	return &SenderError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e *SenderError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *SenderError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target error.
func (e *SenderError) Is(target error) bool {
	if targetErr, ok := target.(*SenderError); ok {
		return e.Code == targetErr.Code
	}
	return false
}

// IsSenderError checks if an error is a SenderError.
func IsSenderError(err error) bool {
	senderError := &SenderError{}
	ok := errors.As(err, &senderError)
	return ok
}

// GetSenderErrorCode returns the error code if the error is a SenderError.
func GetSenderErrorCode(err error) ErrorCode {
	senderErr := &SenderError{}
	if errors.As(err, &senderErr) {
		return senderErr.Code
	}
	return ErrCodeUnknown
}
