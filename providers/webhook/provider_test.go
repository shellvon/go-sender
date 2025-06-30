package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

// TestWebhookProviderWithJSONResponse 测试 webhook provider 的 JSON 响应验证功能
func TestWebhookProviderWithJSONResponse(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 检查 Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		// 读取请求体
		var requestBody map[string]interface{}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		fmt.Printf("DEBUG: Received request body: %s\n", string(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			fmt.Printf("DEBUG: JSON unmarshal error: %v\n", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// 检查请求体是否包含必要的字段
		if _, exists := requestBody["message"]; !exists {
			// 返回错误响应，但状态码为200，让响应验证处理业务逻辑错误
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // 状态码 200，但业务逻辑失败
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  1, // 错误状态
				"error":   "Missing required field: message",
				"message": "Request validation failed",
			})
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  0, // 成功状态
			"message": "Message sent successfully",
			"id":      "msg_12345",
		})
	}))
	defer server.Close()

	// 创建 webhook 配置
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Endpoints: []Endpoint{
			{
				Name:   "test-webhook",
				URL:    server.URL,
				Method: "POST",
				ResponseConfig: &ResponseConfig{
					ValidateResponse: true,
					ResponseType:     "json",
					SuccessField:     "status",
					SuccessValue:     "0", // 期望 status 字段为 0 表示成功
					ErrorField:       "error",
					MessageField:     "message",
				},
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 测试用例
	tests := []struct {
		name    string
		message *Message
		wantErr bool
	}{
		{
			name: "successful message",
			message: &Message{
				DefaultMessage: core.DefaultMessage{},
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: func() []byte {
					data := map[string]interface{}{
						"message": "Hello from test",
						"time":    time.Now().Unix(),
					}
					body, _ := json.Marshal(data)
					return body
				}(),
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			message: &Message{
				DefaultMessage: core.DefaultMessage{},
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: func() []byte {
					data := map[string]interface{}{
						"time": time.Now().Unix(),
						// 故意不包含 message 字段
					}
					body, _ := json.Marshal(data)
					return body
				}(),
			},
			wantErr: true, // 期望失败，因为缺少 message 字段
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := provider.Send(ctx, tt.message, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// 打印结果用于调试
			if err != nil {
				t.Logf("Error: %v", err)
			} else {
				t.Logf("Success: Message sent successfully")
			}
		})
	}
}

// TestWebhookProviderWithCustomStatusCodes 测试自定义状态码功能
func TestWebhookProviderWithCustomStatusCodes(t *testing.T) {
	// 创建测试服务器，返回 201 状态码
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 返回 201
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  0,
			"message": "Resource created successfully",
		})
	}))
	defer server.Close()

	// 创建配置，只接受 201 状态码
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Endpoints: []Endpoint{
			{
				Name:   "test-webhook",
				URL:    server.URL,
				Method: "POST",
				ResponseConfig: &ResponseConfig{
					SuccessStatusCodes: []int{201}, // 只接受 201 状态码
					ValidateResponse:   true,
					ResponseType:       "json",
					SuccessField:       "status",
					SuccessValue:       "0",
				},
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建消息
	message := &Message{
		DefaultMessage: core.DefaultMessage{},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: func() []byte {
			data := map[string]interface{}{
				"message": "Test message",
			}
			body, _ := json.Marshal(data)
			return body
		}(),
	}

	// 发送消息
	ctx := context.Background()
	err = provider.Send(ctx, message, nil)

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	} else {
		t.Logf("Success: Message sent with custom status code")
	}
}

// TestWebhookProviderWithTextResponse 测试文本响应验证
func TestWebhookProviderWithTextResponse(t *testing.T) {
	// 创建测试服务器，返回文本响应
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// 创建配置，验证文本响应
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Endpoints: []Endpoint{
			{
				Name:   "test-webhook",
				URL:    server.URL,
				Method: "POST",
				ResponseConfig: &ResponseConfig{
					ValidateResponse: true,
					ResponseType:     "text",
					SuccessPattern:   "^OK$",
					ErrorPattern:     "^ERROR:",
				},
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建消息
	message := &Message{
		DefaultMessage: core.DefaultMessage{},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: func() []byte {
			data := map[string]interface{}{
				"message": "Test message",
			}
			body, _ := json.Marshal(data)
			return body
		}(),
	}

	// 发送消息
	ctx := context.Background()
	err = provider.Send(ctx, message, nil)

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	} else {
		t.Logf("Success: Text response validation passed")
	}
}

// TestWebhookProviderWithErrorResponse 测试错误响应处理
func TestWebhookProviderWithErrorResponse(t *testing.T) {
	// 创建测试服务器，返回错误响应
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 状态码 200，但业务逻辑失败
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  1, // 错误状态
			"error":   "Invalid API key",
			"message": "Authentication failed",
		})
	}))
	defer server.Close()

	// 创建配置
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Endpoints: []Endpoint{
			{
				Name:   "test-webhook",
				URL:    server.URL,
				Method: "POST",
				ResponseConfig: &ResponseConfig{
					ValidateResponse: true,
					ResponseType:     "json",
					SuccessField:     "status",
					SuccessValue:     "0",
					ErrorField:       "error",
					MessageField:     "message",
				},
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建消息
	message := &Message{
		DefaultMessage: core.DefaultMessage{},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: func() []byte {
			data := map[string]interface{}{
				"message": "Test message",
			}
			body, _ := json.Marshal(data)
			return body
		}(),
	}

	// 发送消息，期望失败
	ctx := context.Background()
	err = provider.Send(ctx, message, nil)

	if err == nil {
		t.Errorf("Expected error but got none")
	} else {
		t.Logf("Expected error received: %v", err)
	}
}

// TestWebhookProviderWithoutValidation 测试不进行响应验证的情况
func TestWebhookProviderWithoutValidation(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": 1, // 即使返回错误状态，也不应该验证
			"error":  "Some error",
		})
	}))
	defer server.Close()

	// 创建配置，不进行响应验证
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Endpoints: []Endpoint{
			{
				Name:   "test-webhook",
				URL:    server.URL,
				Method: "POST",
				// 不设置 ResponseConfig，或者设置 ValidateResponse: false
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建消息
	message := &Message{
		DefaultMessage: core.DefaultMessage{},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: func() []byte {
			data := map[string]interface{}{
				"message": "Test message",
			}
			body, _ := json.Marshal(data)
			return body
		}(),
	}

	// 发送消息，应该成功（只检查状态码）
	ctx := context.Background()
	err = provider.Send(ctx, message, nil)

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	} else {
		t.Logf("Success: No validation mode works correctly")
	}
}

// TestWebhookProviderWithDynamicURL 测试动态URL构建功能
func TestWebhookProviderWithDynamicURL(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法和路径
		if r.Method != "PUT" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 检查路径是否包含用户ID
		if !strings.Contains(r.URL.Path, "/users/12345/") {
			http.Error(w, "Invalid path", http.StatusNotFound)
			return
		}

		// 检查查询参数
		userID := r.URL.Query().Get("user_id")
		notificationType := r.URL.Query().Get("type")
		if userID != "12345" || notificationType != "alert" {
			http.Error(w, "Invalid query parameters", http.StatusBadRequest)
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  0,
			"message": "Notification sent successfully",
		})
	}))
	defer server.Close()

	// 创建配置，使用URL模板
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Endpoints: []Endpoint{
			{
				Name:   "dynamic-webhook",
				URL:    server.URL + "/api/users/{user_id}/notifications",
				Method: "POST", // 会被message中的method覆盖
				ResponseConfig: &ResponseConfig{
					ValidateResponse: true,
					ResponseType:     "json",
					SuccessField:     "status",
					SuccessValue:     "0",
				},
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建消息，包含动态参数
	message := &Message{
		DefaultMessage: core.DefaultMessage{},
		Method:         "PUT", // 覆盖endpoint的method
		PathParams: map[string]string{
			"user_id": "12345",
		},
		QueryParams: map[string]string{
			"user_id": "12345",
			"type":    "alert",
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: func() []byte {
			data := map[string]interface{}{
				"message":  "Test notification",
				"priority": "high",
			}
			body, _ := json.Marshal(data)
			return body
		}(),
	}

	// 发送消息
	ctx := context.Background()
	err = provider.Send(ctx, message, nil)

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	} else {
		t.Logf("Success: Dynamic URL building works correctly")
	}
}

// TestWebhookProviderWithQueryParamsOnly 测试仅查询参数的情况
func TestWebhookProviderWithQueryParamsOnly(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 检查查询参数
		action := r.URL.Query().Get("action")
		timestamp := r.URL.Query().Get("timestamp")
		if action != "ping" || timestamp == "" {
			http.Error(w, "Invalid query parameters", http.StatusBadRequest)
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  0,
			"message": "Ping successful",
		})
	}))
	defer server.Close()

	// 创建配置
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Endpoints: []Endpoint{
			{
				Name:   "query-webhook",
				URL:    server.URL + "/api/health",
				Method: "POST", // 会被message中的method覆盖
				ResponseConfig: &ResponseConfig{
					ValidateResponse: true,
					ResponseType:     "json",
					SuccessField:     "status",
					SuccessValue:     "0",
				},
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建消息，只包含查询参数
	message := &Message{
		DefaultMessage: core.DefaultMessage{},
		Method:         "GET",
		QueryParams: map[string]string{
			"action":    "ping",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		},
		Body: []byte{}, // GET请求通常不需要body
	}

	// 发送消息
	ctx := context.Background()
	err = provider.Send(ctx, message, nil)

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	} else {
		t.Logf("Success: Query parameters only works correctly")
	}
}
