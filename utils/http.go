package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	// For a default timeout
)

// RequestOptions holds the parameters for your HTTP request.
// Similar to Python requests, supports data, json, form fields for easy data handling.
type HTTPRequestOptions struct {
	Method  string            // e.g., http.MethodGet, http.MethodPost
	Headers map[string]string // Custom headers
	Timeout time.Duration     // Request timeout

	// Query string parameters
	// Supports:
	//   - string: single value (e.g., "api_key": "your-key")
	//   - []string: array values (e.g., "tags": []string{"tag1", "tag2"})
	Query map[string]interface{} // Query string parameters (only string or []string allowed)

	// Data handling (similar to Python requests)
	Data      map[string]string // Form data (application/x-www-form-urlencoded)
	JSON      interface{}       // JSON data (application/json)
	Form      map[string]string // Multipart form data
	Raw       []byte            // Raw body data
	RawReader io.Reader         // Custom reader for body

	// File upload support
	Files map[string]string // Field name -> file path for multipart uploads

	// Client allows custom HTTP client for this request. Only affects HTTP-based providers; SMTP/email is not affected.
	Client *http.Client // Optional: custom HTTP client (proxy, timeout, etc.)
}

// DoRequest performs an HTTP request and returns the response body.
// It automatically handles data formatting based on the provided fields.
//
// Data handling priority:
// 1. Raw/RawReader (highest priority)
// 2. JSON
// 3. Form (multipart)
// 4. Data (urlencoded)
// 5. Files
//
// Returns:
//   - []byte: Response body
//   - int: HTTP status code
//   - error: Request error
func DoRequest(ctx context.Context, requestURL string, options HTTPRequestOptions) ([]byte, int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Handle timeout
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	// Determine request body and content type
	reqBody, contentType, err := buildRequestBody(options)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build request body: %w", err)
	}

	// Handle query string parameters
	finalURL := requestURL
	if len(options.Query) > 0 {
		parsedURL, err := url.Parse(requestURL)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse URL: %w", err)
		}

		// Get existing query parameters
		query := parsedURL.Query()

		// Add new query parameters
		for key, value := range options.Query {
			switch v := value.(type) {
			case string:
				query.Set(key, v)
			case []string:
				for _, val := range v {
					query.Add(key, val)
				}
			default:
				return nil, 0, fmt.Errorf("unsupported query parameter type for key '%s': %T (only string and []string are supported)", key, value)
			}
		}

		// Update URL with query parameters
		parsedURL.RawQuery = query.Encode()
		finalURL = parsedURL.String()
	}

	// Set default method if not provided
	if options.Method == "" {
		if reqBody != nil {
			options.Method = http.MethodPost
		} else {
			options.Method = http.MethodGet
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, options.Method, finalURL, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// Set default User-Agent if not provided
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", core.DefaultUserAgent)
	}

	// Set content type if determined
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Use custom client if provided, otherwise default
	client := core.EnsureHTTPClient(options.Client)

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", readErr)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return bodyBytes, resp.StatusCode, fmt.Errorf("HTTP request failed with status code %d. Response body: %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, resp.StatusCode, nil
}

// buildRequestBody builds the request body and determines content type
func buildRequestBody(options HTTPRequestOptions) (io.Reader, string, error) {
	// Priority 1: Raw data
	if options.Raw != nil {
		return bytes.NewReader(options.Raw), "", nil
	}

	if options.RawReader != nil {
		return options.RawReader, "", nil
	}

	// Priority 2: JSON data
	if options.JSON != nil {
		jsonData, err := json.Marshal(options.JSON)
		if err != nil {
			return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return bytes.NewReader(jsonData), "application/json", nil
	}

	// Priority 3: Multipart form data
	if len(options.Form) > 0 || len(options.Files) > 0 {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add form fields
		for key, value := range options.Form {
			if err := writer.WriteField(key, value); err != nil {
				return nil, "", fmt.Errorf("failed to write form field %s: %w", key, err)
			}
		}

		// Add files
		for fieldName, filePath := range options.Files {
			if err := addFileToMultipart(writer, fieldName, filePath); err != nil {
				return nil, "", err
			}
		}

		if err := writer.Close(); err != nil {
			return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
		}

		return body, writer.FormDataContentType(), nil
	}

	// Priority 4: URL-encoded form data
	if len(options.Data) > 0 {
		values := url.Values{}
		for key, value := range options.Data {
			values.Set(key, value)
		}
		encoded := values.Encode()
		return strings.NewReader(encoded), "application/x-www-form-urlencoded", nil
	}

	// No body
	return nil, "", nil
}

// addFileToMultipart adds a file to multipart form
func addFileToMultipart(writer *multipart.Writer, fieldName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file part: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}
