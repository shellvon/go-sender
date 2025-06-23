package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDoRequest_GET(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check query parameters
		expectedParams := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		for key, expectedValue := range expectedParams {
			if got := r.URL.Query().Get(key); got != expectedValue {
				t.Errorf("Expected query param %s=%s, got %s", key, expectedValue, got)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	// Test GET request with query parameters
	options := RequestOptions{
		Method: http.MethodGet,
		Query: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	body, status, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}

	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	expectedBody := `{"status": "ok"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(body))
	}
}

func TestDoRequest_QueryArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check array query parameters
		tags := r.URL.Query()["tags"]
		expectedTags := []string{"tag1", "tag2", "tag3"}

		if len(tags) != len(expectedTags) {
			t.Errorf("Expected %d tags, got %d", len(expectedTags), len(tags))
		}

		for i, tag := range expectedTags {
			if i >= len(tags) || tags[i] != tag {
				t.Errorf("Expected tag[%d]=%s, got %s", i, tag, tags[i])
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"tags": "received"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method: http.MethodGet,
		Query: map[string]interface{}{
			"tags": []string{"tag1", "tag2", "tag3"},
		},
	}

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}

func TestDoRequest_POST_JSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Parse JSON body
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			t.Errorf("Failed to decode JSON: %v", err)
		}

		expectedData := map[string]interface{}{
			"name": "John",
			"age":  float64(30), // JSON numbers are float64
		}

		if data["name"] != expectedData["name"] {
			t.Errorf("Expected name %v, got %v", expectedData["name"], data["name"])
		}
		if data["age"] != expectedData["age"] {
			t.Errorf("Expected age %v, got %v", expectedData["age"], data["age"])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method: http.MethodPost,
		JSON: map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}

func TestDoRequest_POST_Form(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			t.Errorf("Expected Content-Type application/x-www-form-urlencoded, got %s", contentType)
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse form: %v", err)
		}

		expectedData := map[string]string{
			"field1": "value1",
			"field2": "value2",
		}

		for key, expectedValue := range expectedData {
			if got := r.FormValue(key); got != expectedValue {
				t.Errorf("Expected form field %s=%s, got %s", key, expectedValue, got)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method: http.MethodPost,
		Data: map[string]string{
			"field1": "value1",
			"field2": "value2",
		},
	}

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}

func TestDoRequest_POST_Multipart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			t.Errorf("Expected Content-Type multipart/form-data, got %s", contentType)
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}

		// Check form fields
		expectedFields := map[string]string{
			"field1": "value1",
			"field2": "value2",
		}

		for key, expectedValue := range expectedFields {
			if got := r.FormValue(key); got != expectedValue {
				t.Errorf("Expected form field %s=%s, got %s", key, expectedValue, got)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method: http.MethodPost,
		Form: map[string]string{
			"field1": "value1",
			"field2": "value2",
		},
	}

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}

func TestDoRequest_FileUpload(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	fileContent := "Hello, World!"
	if err := os.WriteFile(tempFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			t.Errorf("Expected Content-Type multipart/form-data, got %s", contentType)
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}

		// Check file
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Errorf("Failed to get file: %v", err)
		}
		defer file.Close()

		if header.Filename != "test.txt" {
			t.Errorf("Expected filename test.txt, got %s", header.Filename)
		}

		// Read file content
		content, err := io.ReadAll(file)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}

		if string(content) != fileContent {
			t.Errorf("Expected file content %s, got %s", fileContent, string(content))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method: http.MethodPost,
		Files: map[string]string{
			"file": tempFile,
		},
	}

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}

func TestDoRequest_Timeout(t *testing.T) {
	// Create a slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Sleep for 2 seconds
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method:  http.MethodGet,
		Timeout: 1 * time.Second, // 1 second timeout
	}

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "deadline exceeded") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestDoRequest_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method: http.MethodGet,
	}

	body, status, err := DoRequest(context.Background(), server.URL, options)
	if err == nil {
		t.Error("Expected error for 404 status, got nil")
	}

	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}

	expectedBody := `{"error": "not found"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(body))
	}
}

func TestDoRequest_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check custom headers
		expectedHeaders := map[string]string{
			"X-Custom-Header": "custom-value",
			"Authorization":   "Bearer token123",
		}

		for key, expectedValue := range expectedHeaders {
			if got := r.Header.Get(key); got != expectedValue {
				t.Errorf("Expected header %s=%s, got %s", key, expectedValue, got)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		Method: http.MethodGet,
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
			"Authorization":   "Bearer token123",
		},
	}

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}

func TestDoRequest_DefaultMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test default GET method (no body)
		if r.Method != http.MethodGet {
			t.Errorf("Expected default GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{} // No method specified

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}

func TestDoRequest_DefaultPostMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test default POST method (with body)
		if r.Method != http.MethodPost {
			t.Errorf("Expected default POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	options := RequestOptions{
		JSON: map[string]interface{}{
			"test": "data",
		},
	} // No method specified, but has body

	_, _, err := DoRequest(context.Background(), server.URL, options)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
}
