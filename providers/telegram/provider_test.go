package telegram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shellvon/go-sender/core"
)

// TestTelegramProviderWithTextMessage 测试文本消息发送
func TestTelegramProviderWithTextMessage(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 检查 Content-Type
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		// 解析 JSON 数据
		var requestData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		// 检查必需字段
		chatID, ok := requestData["chat_id"]
		if !ok || chatID == "" {
			http.Error(w, "Missing chat_id", http.StatusBadRequest)
			return
		}
		text, ok := requestData["text"]
		if !ok || text == "" {
			http.Error(w, "Missing text", http.StatusBadRequest)
			return
		}

		// 检查可选字段（用于验证但不使用）
		_ = requestData["parse_mode"]
		_ = requestData["disable_web_page_preview"]
		_ = requestData["disable_notification"]
		_ = requestData["protect_content"]
		_ = requestData["reply_to_message_id"]

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12345,
				"chat": map[string]interface{}{
					"id": chatID,
				},
				"text": text,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建配置
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试用例
	tests := []struct {
		name    string
		message *TextMessage
		wantErr bool
	}{
		{
			name:    "basic text message",
			message: NewTextMessage("@test_channel", "Hello from test"),
			wantErr: false,
		},
		{
			name: "text message with markdown",
			message: NewTextMessage("@test_channel", "Hello *bold* and _italic_",
				WithParseMode("Markdown")),
			wantErr: false,
		},
		{
			name: "text message with all options",
			message: NewTextMessage("@test_channel", "Test message",
				WithParseMode("HTML"),
				WithDisableWebPreview(true),
				WithSilent(true),
				WithProtectContent(true),
				WithReplyTo(123)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := provider.Send(ctx, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTelegramProviderWithPhotoMessage 测试图片消息发送
func TestTelegramProviderWithPhotoMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 检查 Content-Type
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		// 解析 JSON 数据
		var requestData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		chatID, ok := requestData["chat_id"]
		if !ok || chatID == "" {
			http.Error(w, "Missing chat_id", http.StatusBadRequest)
			return
		}
		photo, ok := requestData["photo"]
		if !ok || photo == "" {
			http.Error(w, "Missing photo", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12346,
				"chat": map[string]interface{}{
					"id": chatID,
				},
				"photo": []map[string]interface{}{
					{"file_id": "test_photo_id"},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试图片消息
	message := NewPhotoMessage("@test_channel", "https://example.com/test.jpg",
		WithCaption("Test photo caption"),
		WithParseMode("Markdown"),
		WithSilent(true),
		WithProtectContent(true),
		WithReplyParameters(&ReplyParameters{MessageID: 123}))

	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err != nil {
		t.Errorf("Failed to send photo message: %v", err)
	}
}

// TestTelegramProviderWithDocumentMessage 测试文档消息发送
func TestTelegramProviderWithDocumentMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		chatID := r.FormValue("chat_id")
		document := r.FormValue("document")
		if chatID == "" || document == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12347,
				"chat": map[string]interface{}{
					"id": chatID,
				},
				"document": map[string]interface{}{
					"file_id": "test_document_id",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试文档消息
	message := NewDocumentMessage("@test_channel", "https://example.com/test.pdf",
		WithCaption("Test document"),
		WithParseMode("HTML"),
		WithSilent(true),
		WithProtectContent(true),
		WithReplyParameters(&ReplyParameters{MessageID: 123}))

	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err != nil {
		t.Errorf("Failed to send document message: %v", err)
	}
}

// TestTelegramProviderWithLocationMessage 测试位置消息发送
func TestTelegramProviderWithLocationMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		chatID := r.FormValue("chat_id")
		latitude := r.FormValue("latitude")
		longitude := r.FormValue("longitude")
		if chatID == "" || latitude == "" || longitude == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12348,
				"chat": map[string]interface{}{
					"id": chatID,
				},
				"location": map[string]interface{}{
					"latitude":  latitude,
					"longitude": longitude,
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试位置消息
	message := NewLocationMessage("@test_channel", 40.7128, -74.0060,
		WithSilent(true),
		WithProtectContent(true),
		WithReplyParameters(&ReplyParameters{MessageID: 123}))

	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err != nil {
		t.Errorf("Failed to send location message: %v", err)
	}
}

// TestTelegramProviderWithContactMessage 测试联系人消息发送
func TestTelegramProviderWithContactMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		chatID := r.FormValue("chat_id")
		phoneNumber := r.FormValue("phone_number")
		firstName := r.FormValue("first_name")
		if chatID == "" || phoneNumber == "" || firstName == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12349,
				"chat": map[string]interface{}{
					"id": chatID,
				},
				"contact": map[string]interface{}{
					"phone_number": phoneNumber,
					"first_name":   firstName,
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试联系人消息
	message := NewContactMessage("@test_channel", "+1234567890", "John",
		WithContactLastName("Doe"),
		WithContactVCard("BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+1234567890\nEND:VCARD"),
		WithSilent(true),
		WithProtectContent(true),
		WithReplyParameters(&ReplyParameters{MessageID: 123}))

	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err != nil {
		t.Errorf("Failed to send contact message: %v", err)
	}
}

// TestTelegramProviderWithPollMessage 测试投票消息发送
func TestTelegramProviderWithPollMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		chatID := r.FormValue("chat_id")
		question := r.FormValue("question")
		options := r.FormValue("options")
		if chatID == "" || question == "" || options == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12350,
				"chat": map[string]interface{}{
					"id": chatID,
				},
				"poll": map[string]interface{}{
					"question": question,
					"options":  strings.Split(options, "\n"),
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试投票消息
	options := []InputPollOption{
		{Text: "Red"},
		{Text: "Blue"},
		{Text: "Green"},
		{Text: "Yellow"},
	}
	message := NewPollMessage("@test_channel", "What's your favorite color?",
		options,
		WithPollIsAnonymous(true),
		WithPollType("quiz"),
		WithPollAllowsMultipleAnswers(false),
		WithSilent(true),
		WithProtectContent(true),
		WithReplyParameters(&ReplyParameters{MessageID: 123}))

	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err != nil {
		t.Errorf("Failed to send poll message: %v", err)
	}
}

// TestTelegramProviderErrorHandling 测试错误处理
func TestTelegramProviderErrorHandling(t *testing.T) {
	// 测试服务器返回错误
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: chat not found",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试错误响应
	message := NewTextMessage("@invalid_channel", "This should fail")
	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err == nil {
		t.Error("Expected error but got success")
	}
}

// TestTelegramProviderValidation 测试消息验证
func TestTelegramProviderValidation(t *testing.T) {
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试无效消息类型
	ctx := context.Background()

	// 创建一个不实现Message接口的值
	var invalidMessage interface{} = "not a message"

	// 使用类型断言来测试错误处理
	err = provider.Send(ctx, invalidMessage.(core.Message))
	if err == nil {
		t.Error("Expected error for invalid message type but got success")
	}

	// 测试空文本消息
	emptyTextMessage := NewTextMessage("@test_channel", "")
	err = provider.Send(ctx, emptyTextMessage)
	if err == nil {
		t.Error("Expected error for empty text but got success")
	}

	// 测试空聊天ID
	emptyChatMessage := NewTextMessage("", "Test message")
	err = provider.Send(ctx, emptyChatMessage)
	if err == nil {
		t.Error("Expected error for empty chat_id but got success")
	}
}

// TestTelegramProviderMultipleAccounts 测试多账户配置
func TestTelegramProviderMultipleAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12351,
				"chat": map[string]interface{}{
					"id": "@test_channel",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "bot1",
				Key:      "bot1-token",
				Weight:   100,
				Disabled: false,
			},
			{
				Name:     "bot2",
				Key:      "bot2-token",
				Weight:   50,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试多账户发送
	message := NewTextMessage("@test_channel", "Test with multiple accounts")
	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err != nil {
		t.Errorf("Failed to send message with multiple accounts: %v", err)
	}
}

// TestTelegramProviderDisabledAccount 测试禁用账户
func TestTelegramProviderDisabledAccount(t *testing.T) {
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "disabled-bot",
				Key:      "disabled-bot-token",
				Weight:   100,
				Disabled: true, // 禁用账户
			},
		},
	}

	_, err := New(config)
	if err == nil {
		t.Error("Expected error for disabled account but got success")
	}
}

// TestTelegramProviderEmptyConfig 测试空配置
func TestTelegramProviderEmptyConfig(t *testing.T) {
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{}, // 空账户列表
	}

	_, err := New(config)
	if err == nil {
		t.Error("Expected error for empty accounts but got success")
	}
}

// TestTelegramProviderJSONSerialization 测试 JSON 序列化
func TestTelegramProviderJSONSerialization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求体是否为有效的 JSON
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// 验证必需的字段
		if requestBody["chat_id"] == nil {
			http.Error(w, "Missing chat_id", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"ok": true,
			"result": map[string]interface{}{
				"message_id": 12345,
				"chat": map[string]interface{}{
					"id": requestBody["chat_id"],
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "test-bot",
				Key:      "test-bot-token",
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	// 测试文本消息的 JSON 序列化
	message := NewTextMessage("@test_channel", "Test message with JSON serialization",
		WithSilent(true),
		WithProtectContent(true),
		WithReplyParameters(&ReplyParameters{MessageID: 123}))

	ctx := context.Background()
	err = provider.Send(ctx, message)
	if err != nil {
		t.Errorf("Failed to send message with JSON serialization: %v", err)
	}
}
