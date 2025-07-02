package core_test

import (
	"errors"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestParamError(t *testing.T) {
	err := core.NewParamError("bad param")
	if err.Error() == "" || err.IsRetryable() {
		t.Error("ParamError should not be retryable and must have message")
	}
}

func TestNetworkError(t *testing.T) {
	err := core.NetworkError{Err: errors.New("net fail")}
	if !err.IsRetryable() || err.Error() == "" {
		t.Error("NetworkError should be retryable and have message")
	}
}

func TestTimeoutError(t *testing.T) {
	err := core.TimeoutError{Err: errors.New("timeout fail")}
	if !err.IsRetryable() || err.Error() == "" {
		t.Error("TimeoutError should be retryable and have message")
	}
}

func TestValidationError(t *testing.T) {
	err := core.ValidationError{Err: errors.New("validation fail")}
	if err.IsRetryable() || err.Error() == "" {
		t.Error("ValidationError should not be retryable and must have message")
	}
}

func TestAuthenticationError(t *testing.T) {
	err := core.AuthenticationError{Err: errors.New("auth fail")}
	if err.IsRetryable() || err.Error() == "" {
		t.Error("AuthenticationError should not be retryable and must have message")
	}
}

func TestSenderError_IsUnwrap(t *testing.T) {
	base := errors.New("base")
	err := core.NewSenderError(core.ErrCodeProviderSendFailed, "fail", base)
	if !errors.Is(err, base) {
		t.Error("SenderError should unwrap to base error")
	}
	if !errors.Is(err.Unwrap(), base) {
		t.Error("SenderError Unwrap failed")
	}
	if !core.IsSenderError(err) {
		t.Error("IsSenderError should detect SenderError")
	}
	if core.GetSenderErrorCode(err) != core.ErrCodeProviderSendFailed {
		t.Error("GetSenderErrorCode should extract code")
	}
}

func TestSenderError_Is(t *testing.T) {
	err1 := core.NewSenderError(core.ErrCodeProviderSendFailed, "fail", nil)
	err2 := core.NewSenderError(core.ErrCodeProviderSendFailed, "fail", nil)
	if !errors.Is(err1, err2) {
		t.Error("SenderError Is should match same code")
	}
	err3 := core.NewSenderError(core.ErrCodeProviderUnavailable, "fail", nil)
	if errors.Is(err1, err3) {
		t.Error("SenderError Is should not match different code")
	}
}

func TestDefaultErrorClassifier_NetworkTimeoutSystemCall(t *testing.T) {
	c := core.NewDefaultErrorClassifier()
	// isNetworkError
	err := errors.New("connection refused")
	if !c.IsRetryableError(core.NetworkError{Err: err}) {
		t.Error("should be retryable network error")
	}
	// isTimeoutError
	err = errors.New("timeout")
	if !c.IsRetryableError(core.TimeoutError{Err: err}) {
		t.Error("should be retryable timeout error")
	}
	// isSystemCallError (simulate syscall error by string)
	err = errors.New("broken pipe")
	if !c.IsRetryableError(core.NetworkError{Err: err}) {
		t.Error("should be retryable system call error")
	}
	// isHTTP5xxError
	err = errors.New("500 internal server error")
	if !c.IsRetryableError(core.NetworkError{Err: err}) {
		t.Error("should be retryable http 5xx error")
	}
	// isTypeAssertionError
	err = errors.New("type assertion failed")
	if !c.IsRetryableError(core.NetworkError{Err: err}) {
		t.Error("should be retryable type assertion error")
	}
	// isJSONError
	err = errors.New("invalid character 'a' looking for beginning of value")
	if !c.IsRetryableError(core.NetworkError{Err: err}) {
		t.Error("should be retryable json error")
	}
	// fallback: not retryable
	err = errors.New("some other error")
	if c.IsRetryableError(core.ValidationError{Err: err}) {
		t.Error("should not be retryable for unrelated error")
	}
}

func TestIsRetryableError(t *testing.T) {
	classifier := core.NewDefaultErrorClassifier()
	err := core.NewSenderError(core.ErrCodeInternal, "net", nil)
	if classifier.IsRetryableError(err) {
		t.Error("SenderError with internal code should not be retryable")
	}
	if classifier.IsRetryableError(errors.New("other")) {
		t.Error("should not be retryable")
	}
}

func TestIsNetworkError(t *testing.T) {
	classifier := core.NewDefaultErrorClassifier()
	err := &net.OpError{Op: "dial", Net: "tcp", Err: os.ErrDeadlineExceeded}
	if !classifier.IsRetryableError(core.NetworkError{Err: err}) {
		t.Error("NetworkError should be retryable")
	}
}

func TestIsTimeoutError(t *testing.T) {
	classifier := core.NewDefaultErrorClassifier()
	err := os.ErrDeadlineExceeded
	if !classifier.IsRetryableError(core.TimeoutError{Err: err}) {
		t.Error("TimeoutError should be retryable")
	}
}

func TestIsSystemCallError(t *testing.T) {
	classifier := core.NewDefaultErrorClassifier()
	err := syscall.ECONNRESET
	if !classifier.IsRetryableError(err) {
		t.Error("syscall should be retryable")
	}
}

func TestIsHTTP5xxError(t *testing.T) {
	classifier := core.NewDefaultErrorClassifier()
	err := errors.New("server error 502")
	if !classifier.IsRetryableError(err) {
		t.Error("5xx should be retryable")
	}
}

func TestIsTypeAssertionError(t *testing.T) {
	classifier := core.NewDefaultErrorClassifier()
	err := errors.New("interface conversion: type assertion failed")
	if classifier.IsRetryableError(err) {
		t.Error("type assertion should not be retryable")
	}
}

func TestIsJSONError(t *testing.T) {
	classifier := core.NewDefaultErrorClassifier()
	err := errors.New("json: cannot unmarshal string")
	if classifier.IsRetryableError(err) {
		t.Error("json error should not be retryable")
	}
}

func TestNewSenderErrorfAndError(t *testing.T) {
	err := core.NewSenderErrorf(core.ErrCodeUnknown, "fail: %d", 42)
	if err == nil || err.Error() == "" {
		t.Error("NewSenderErrorf should return error with message")
	}
}

func TestGetSenderErrorCode(t *testing.T) {
	err := core.NewSenderError(core.ErrCodeInternal, "fail", nil)
	if code := core.GetSenderErrorCode(err); code != core.ErrCodeInternal {
		t.Errorf("expected code %v, got %v", core.ErrCodeInternal, code)
	}
	if code := core.GetSenderErrorCode(errors.New("other")); code != core.ErrCodeUnknown {
		t.Errorf("expected unknown code, got %v", code)
	}
}
