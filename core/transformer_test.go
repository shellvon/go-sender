package core_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

func TestHTTPRequestSpecFields(t *testing.T) {
	spec := &core.HTTPRequestSpec{
		Method:      "POST",
		URL:         "http://example.com",
		Headers:     map[string]string{"k": "v"},
		QueryParams: map[string]string{"q": "1"},
		Body:        []byte("data"),
		BodyType:    "json",
		Timeout:     2 * time.Second,
	}
	if spec.Method != http.MethodPost || spec.URL != "http://example.com" || spec.Headers["k"] != "v" ||
		spec.QueryParams["q"] != "1" ||
		string(spec.Body) != "data" ||
		spec.BodyType != "json" ||
		spec.Timeout != 2*time.Second {
		t.Errorf("HTTPRequestSpec fields not set correctly: %+v", spec)
	}
}

func TestResponseHandler(t *testing.T) {
	h := func(status int, body []byte) error {
		if status != 200 || string(body) != "ok" {
			return errors.New("bad response")
		}
		return nil
	}
	if err := h(200, []byte("ok")); err != nil {
		t.Errorf("ResponseHandler should succeed: %v", err)
	}
	if err := h(500, []byte("fail")); err == nil {
		t.Error("ResponseHandler should fail on bad input")
	}
}
