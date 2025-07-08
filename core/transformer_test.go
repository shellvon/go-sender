package core_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

func TestHTTPRequestSpecFields(t *testing.T) {
	spec := &core.HTTPRequestSpec{
		Method:      "POST",
		URL:         "http://example.com",
		Headers:     map[string]string{"k": "v"},
		QueryParams: url.Values{"q": {"1"}},
		Body:        []byte("data"),
		BodyType:    core.BodyTypeJSON,
		Timeout:     2 * time.Second,
	}
	if spec.Method != http.MethodPost || spec.URL != "http://example.com" || spec.Headers["k"] != "v" ||
		spec.QueryParams.Get("q") != "1" ||
		string(spec.Body) != "data" ||
		spec.BodyType != core.BodyTypeJSON ||
		spec.Timeout != 2*time.Second {
		t.Errorf("HTTPRequestSpec fields not set correctly: %+v", spec)
	}
}
