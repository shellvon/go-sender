package utils_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shellvon/go-sender/utils"
)

func TestDoRequest_Integration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		case "/bad":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
		case "/fail":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
		case "/timeout":
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("late"))
		}
	}))
	defer ts.Close()

	t.Run("2xx", func(t *testing.T) {
		body, code, err := utils.DoRequest(
			context.Background(),
			ts.URL+"/ok",
			utils.HTTPRequestOptions{Timeout: time.Second},
		)
		if err != nil || code != 200 || string(body) != "ok" {
			t.Errorf("2xx failed: code=%d, err=%v, body=%s", code, err, body)
		}
	})
	t.Run("4xx", func(t *testing.T) {
		_, code, err := utils.DoRequest(
			context.Background(),
			ts.URL+"/bad",
			utils.HTTPRequestOptions{Timeout: time.Second},
		)
		if err != nil || code != 400 {
			t.Errorf("4xx should nil, got code=%d, err=%v", code, err)
		}
	})
	t.Run("5xx", func(t *testing.T) {
		_, code, err := utils.DoRequest(
			context.Background(),
			ts.URL+"/fail",
			utils.HTTPRequestOptions{Timeout: time.Second},
		)
		if err != nil || code != 500 {
			t.Errorf("5xx should nil, got code=%d, err=%v", code, err)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		start := time.Now()
		_, _, err := utils.DoRequest(
			context.Background(),
			ts.URL+"/timeout",
			utils.HTTPRequestOptions{Timeout: 50 * time.Millisecond},
		)
		if err == nil || time.Since(start) > 150*time.Millisecond {
			t.Errorf("timeout should error quickly, err=%v", err)
		}
	})
	// 网络错误
	t.Run("network error", func(t *testing.T) {
		_, _, err := utils.DoRequest(
			context.Background(),
			"http://127.0.0.1:0/",
			utils.HTTPRequestOptions{Timeout: 100 * time.Millisecond},
		)
		if err == nil {
			t.Error("should fail on network error")
		}
	})
}

func testMethodAndParamsCase(t *testing.T, url, path, method string, opt utils.HTTPRequestOptions, wantBody string) {
	body, code, err := utils.DoRequest(context.Background(), url+path, opt)
	if err != nil || code != 200 || (wantBody != "" && string(body) != wantBody) {
		t.Errorf("%s %s failed: code=%d, err=%v, body=%s", method, path, code, err, body)
	}
}

func TestDoRequest_MethodsAndParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	}))
	defer ts.Close()

	t.Run("GET with query", func(t *testing.T) {
		params := url.Values{"foo": {"bar"}, "baz": {"qux"}}
		opt := utils.HTTPRequestOptions{Method: "GET", Query: params}
		testMethodAndParamsCase(t, ts.URL, "/get", "GET", opt, "")
	})

	t.Run("POST json", func(t *testing.T) {
		obj := map[string]interface{}{"a": 1, "b": "x"}
		opt := utils.HTTPRequestOptions{Method: "POST", JSON: obj}
		testMethodAndParamsCase(t, ts.URL, "/json", "POST", opt, "")
	})

	t.Run("POST form-urlencoded", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "POST", Data: map[string]string{"foo": "bar"}}
		testMethodAndParamsCase(t, ts.URL, "/form", "POST", opt, "foo=bar")
	})

	t.Run("POST raw", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "POST", Raw: []byte("rawdata")}
		testMethodAndParamsCase(t, ts.URL, "/raw", "POST", opt, "rawdata")
	})

	t.Run("PUT method", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "PUT", Raw: []byte("putdata")}
		testMethodAndParamsCase(t, ts.URL, "/put", "PUT", opt, "putdata")
	})

	t.Run("DELETE method", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "DELETE"}
		testMethodAndParamsCase(t, ts.URL, "/del", "DELETE", opt, "")
	})

	t.Run("HEAD method", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "HEAD"}
		testMethodAndParamsCase(t, ts.URL, "/head", "HEAD", opt, "")
	})

	t.Run("OPTIONS method", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "OPTIONS"}
		testMethodAndParamsCase(t, ts.URL, "/opt", "OPTIONS", opt, "")
	})

	t.Run("query+body", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "POST", Query: url.Values{"q": {"1"}}, Raw: []byte("body")}
		testMethodAndParamsCase(t, ts.URL, "/qbody", "POST", opt, "body")
	})
}

func TestDoRequest_MultipartFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := io.ReadAll(file)
		w.Write(data)
	}))
	defer ts.Close()

	// 创建临时文件
	tmpfile, err := os.CreateTemp(t.TempDir(), "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("filecontent")
	tmpfile.Close()

	opt := utils.HTTPRequestOptions{
		Method: "POST",
		Form:   map[string]string{"foo": "bar"},
		Files:  map[string]string{"file": tmpfile.Name()},
	}
	body, code, err := utils.DoRequest(context.Background(), ts.URL+"/upload", opt)
	if err != nil || code != 200 || string(body) != "filecontent" {
		t.Fatalf("file upload failed: %v, %s", err, body)
	}
}

func TestDoRequest_CustomClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	client := &http.Client{Timeout: 1 * time.Nanosecond}
	opt := utils.HTTPRequestOptions{Client: client}
	_, _, err := utils.DoRequest(context.Background(), ts.URL+"/timeout", opt)
	if err == nil {
		t.Error("should timeout with custom client")
	}
}

func TestDoRequest_ContextCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, _, err := utils.DoRequest(ctx, ts.URL+"/ctx", utils.HTTPRequestOptions{})
	if err == nil {
		t.Error("should error on context cancel")
	}
}

func TestDoRequest_UserAgent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agent := r.Header.Get("User-Agent")
		w.Write([]byte(agent))
	}))
	defer ts.Close()

	t.Run("default user-agent", func(t *testing.T) {
		body, _, err := utils.DoRequest(context.Background(), ts.URL+"/ua", utils.HTTPRequestOptions{})
		if err != nil || !strings.Contains(string(body), "go-sender") {
			t.Errorf("default user-agent not set, got: %s, err: %v", body, err)
		}
	})

	t.Run("custom user-agent", func(t *testing.T) {
		headers := map[string]string{"User-Agent": "my-custom-agent"}
		body, _, err := utils.DoRequest(context.Background(), ts.URL+"/ua", utils.HTTPRequestOptions{Headers: headers})
		if err != nil || string(body) != "my-custom-agent" {
			t.Errorf("custom user-agent not set, got: %s, err: %v", body, err)
		}
	})
}

func TestDoRequest_CustomHeaderHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/header":
			w.Write([]byte(r.Header.Get("X-Custom")))
		default:
			body, _ := io.ReadAll(r.Body)
			w.Write(body)
		}
	}))
	defer ts.Close()

	t.Run("custom header", func(t *testing.T) {
		opt := utils.HTTPRequestOptions{Method: "GET", Headers: map[string]string{"X-Custom": "abc"}}
		body, code, err := utils.DoRequest(context.Background(), ts.URL+"/header", opt)
		if err != nil || code != 200 {
			t.Fatalf("custom header failed: %v", err)
		}
		if string(body) != "abc" {
			t.Errorf("expected custom header value 'abc', got %q", body)
		}
	})
}
