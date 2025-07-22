package mock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
)

// MockTransport is a custom RoundTripper that captures HTTP requests
type MockTransport struct {
	capturedRequests []*cli.HTTPRequestCapture
}

// NewMockTransport creates a new mock transport
func NewMockTransport() *MockTransport {
	return &MockTransport{
		capturedRequests: make([]*cli.HTTPRequestCapture, 0),
	}
}

// NewHTTPClient creates a new http.Client with mock transport
func NewHTTPClient() *http.Client {
	transport := NewMockTransport()
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// GetMockTransport extracts the mock transport from an http.Client
func GetMockTransport(client *http.Client) *MockTransport {
	if transport, ok := client.Transport.(*MockTransport); ok {
		return transport
	}
	return nil
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Capture the request details
	capture := &cli.HTTPRequestCapture{
		Method:    req.Method,
		URL:       req.URL.String(),
		Headers:   make(map[string]string),
		Timestamp: start,
	}

	// Capture headers
	for key, values := range req.Header {
		capture.Headers[key] = strings.Join(values, ", ")
	}

	// Capture body if present
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil && len(bodyBytes) > 0 {
			// Try to parse as JSON for better display
			var bodyContent interface{}
			if strings.Contains(capture.Headers["Content-Type"], "application/json") {
				// For JSON content, store as parsed object
				if err := parseJSONBody(bodyBytes, &bodyContent); err == nil {
					capture.Body = bodyContent
				} else {
					capture.Body = string(bodyBytes)
				}
			} else {
				capture.Body = string(bodyBytes)
			}
		}
		// Reset the body for potential re-reading
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// Simulate processing time
	time.Sleep(10 * time.Millisecond)
	capture.Duration = time.Since(start)

	// Store the captured request
	m.capturedRequests = append(m.capturedRequests, capture)

	// Create a mock successful response - always return 200 OK for dry-run
	mockResponseBody := m.generateMockResponse(req)

	response := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          io.NopCloser(strings.NewReader(mockResponseBody)),
		Request:       req,
		ContentLength: int64(len(mockResponseBody)),
	}

	response.Header.Set("Content-Type", "application/json")
	response.Header.Set("X-Mock-Response", "true")
	response.Header.Set("X-Request-Captured", "true")
	response.Header.Set("Content-Length", fmt.Sprintf("%d", len(mockResponseBody)))

	return response, nil
}

// GetCapturedRequests returns all captured requests
func (m *MockTransport) GetCapturedRequests() []*cli.HTTPRequestCapture {
	return m.capturedRequests
}

// GetLastCapturedRequest returns the most recent captured request
func (m *MockTransport) GetLastCapturedRequest() *cli.HTTPRequestCapture {
	if len(m.capturedRequests) == 0 {
		return nil
	}
	return m.capturedRequests[len(m.capturedRequests)-1]
}

// Clear clears all captured requests
func (m *MockTransport) Clear() {
	m.capturedRequests = make([]*cli.HTTPRequestCapture, 0)
}

// generateMockResponse generates a mock response based on the request
func (m *MockTransport) generateMockResponse(req *http.Request) string {
	// Generate different mock responses based on the URL pattern
	url := req.URL.String()

	switch {
	case strings.Contains(url, "dingtalk.com") || strings.Contains(url, "oapi.dingtalk.com"):
		return `{"errcode":0,"errmsg":"ok"}`
	case strings.Contains(url, "api.telegram.org"):
		return `{"ok":true,"result":{"message_id":123,"date":1234567890,"chat":{"id":123,"type":"private"},"text":"Test message"}}`
	case strings.Contains(url, "open.feishu.cn") || strings.Contains(url, "open.larksuite.com"):
		return `{"code":0,"msg":"success","data":{}}`
	case strings.Contains(url, "qyapi.weixin.qq.com"):
		return `{"errcode":0,"errmsg":"ok"}`
	case strings.Contains(url, "sctapi.ftqq.com"):
		return `{"code":0,"message":"","data":{"pushid":"123456","readkey":"abc123","error":"SUCCESS","errno":0}}`
	case strings.Contains(url, "sms") || strings.Contains(url, "aliyun") || strings.Contains(url, "tencent") ||
		strings.Contains(url, "dysmsapi.aliyuncs.com") || strings.Contains(url, "sms.tencentcloudapi.com"):
		return `{"Message":"OK","RequestId":"12345678-1234-1234-1234-123456789012","BizId":"987654321","Code":"OK"}`
	case strings.Contains(url, "httpbin.org"):
		// For webhook testing
		return `{"args":{},"data":"","files":{},"form":{},"headers":{"Content-Type":"application/json"},"json":null,"origin":"127.0.0.1","url":"` + url + `"}`
	default:
		// Generic successful response for any other URL
		return fmt.Sprintf(`{"status":"success","message":"Mock response for %s %s","timestamp":"%s","request_id":"mock-%d","url":"%s"}`,
			req.Method, req.URL.Path, time.Now().Format(time.RFC3339), time.Now().Unix(), url)
	}
}

// parseJSONBody attempts to parse JSON body content
func parseJSONBody(bodyBytes []byte, target *interface{}) error {
	// This is a simplified JSON parser - in practice you might want more robust parsing
	bodyStr := strings.TrimSpace(string(bodyBytes))
	if strings.HasPrefix(bodyStr, "{") && strings.HasSuffix(bodyStr, "}") {
		// Simple JSON object detection
		*target = map[string]interface{}{
			"_raw_json": bodyStr,
			"_parsed":   true,
		}
		return nil
	}
	return fmt.Errorf("not valid JSON")
}
