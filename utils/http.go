//revive:disable:var-naming
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
	// For a default timeout.
)

// HTTPRequestOptions holds the parameters for your HTTP request.
// Similar to Python requests, supports data, json, form fields for easy data handling.
type HTTPRequestOptions struct {
	Method  string            // e.g., http.MethodGet, http.MethodPost
	Headers map[string]string // Custom headers
	Timeout time.Duration     // Request timeout

	// Query string parameters
	//
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

	// ThrowHTTPErrors: if true, non-2xx/3xx status codes will return error; otherwise, always return body+status.
	ThrowHTTPErrors bool
}

// DoRequest performs an HTTP request and returns the response body.
// It automatically handles data formatting based on the provided fields.
//
// Data handling priority:
//  1. Raw/RawReader (highest priority)
//  2. JSON
//  3. Form (multipart)
//  4. Data (urlencoded)
//  5. Files
//
// Returns:
//   - []byte: Response body
//   - int: HTTP status code
//   - error: Request error
func DoRequest(ctx context.Context, requestURL string, options HTTPRequestOptions) ([]byte, int, error) {
	ctx, cancel := prepareContext(ctx, &options)
	if cancel != nil {
		defer cancel()
	}

	reqBody, contentType, err := buildRequestBody(options)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to build request body: %w", err)
	}

	finalURL, err := buildFinalURL(requestURL, options.Query)
	if err != nil {
		return nil, 0, err
	}

	method := getDefaultMethod(options.Method, reqBody)
	req, err := http.NewRequestWithContext(ctx, method, finalURL, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	setRequestHeaders(req, options.Headers, contentType)
	client := core.EnsureHTTPClient(options.Client)

	resp, err := sendRequest(client, req)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := readResponseBody(resp)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if options.ThrowHTTPErrors && !IsAcceptableStatus(resp.StatusCode) {
		return bodyBytes, resp.StatusCode, fmt.Errorf(
			"HTTP request failed with status code %d. Response body: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	return bodyBytes, resp.StatusCode, nil
}

func prepareContext(ctx context.Context, options *HTTPRequestOptions) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if options.Timeout == 0 {
		options.Timeout = core.DefaultTimeout
	}
	if options.Timeout > 0 {
		return context.WithTimeout(ctx, options.Timeout)
	}
	return ctx, nil
}

// buildFinalURL builds the final request URL by appending query parameters to the base URL.
// Supports string and []string values in the query map. Returns the final URL or an error.
func buildFinalURL(requestURL string, query map[string]interface{}) (string, error) {
	if len(query) == 0 {
		return requestURL, nil
	}
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}
	q := parsedURL.Query()
	for key, value := range query {
		switch v := value.(type) {
		case string:
			q.Set(key, v)
		case []string:
			for _, val := range v {
				q.Add(key, val)
			}
		default:
			return "", fmt.Errorf("unsupported query parameter type for key '%s': %T (only string and []string are supported)", key, value)
		}
	}
	parsedURL.RawQuery = q.Encode()
	return parsedURL.String(), nil
}

// getDefaultMethod returns the HTTP method to use based on the provided method and request body.
// If method is not empty, it is returned as is. If there is a request body, POST is returned. Otherwise, GET is returned.
func getDefaultMethod(method string, reqBody io.Reader) string {
	if method != "" {
		return method
	}
	if reqBody != nil {
		return http.MethodPost
	}
	return http.MethodGet
}

// setRequestHeaders sets the headers for the HTTP request, including the Content-Type if provided.
// If User-Agent is not set, it sets a default User-Agent. Content-Type is set if not empty.
func setRequestHeaders(req *http.Request, headers map[string]string, contentType string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", core.DefaultUserAgent)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
}

// buildRawBody creates an io.Reader from a raw byte slice for the request body.
// Returns the reader, an empty content type, and any error encountered.
func buildRawBody(raw []byte) (io.Reader, string, error) {
	return bytes.NewReader(raw), "", nil
}

// buildRawReaderBody uses the provided io.Reader as the request body.
// Returns the reader, an empty content type, and any error encountered.
func buildRawReaderBody(reader io.Reader) (io.Reader, string, error) {
	return reader, "", nil
}

// buildJSONBody marshals the given object to JSON and returns an io.Reader for the request body.
// Returns the reader, the content type "application/json", and any error encountered.
func buildJSONBody(jsonObj interface{}) (io.Reader, string, error) {
	jsonData, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return bytes.NewReader(jsonData), "application/json", nil
}

// sendRequest sends the HTTP request using the provided client and returns the response.
// Returns the HTTP response and any error encountered.
func sendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

// readResponseBody reads and returns the response body as a byte slice.
// Returns the body and any error encountered.
func readResponseBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}

// IsAcceptableStatus returns true if the status code is 2xx or 3xx.
func IsAcceptableStatus(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusBadRequest
}

// buildRequestBody builds the request body and determines content type.
func buildRequestBody(options HTTPRequestOptions) (io.Reader, string, error) {
	if options.Raw != nil {
		return buildRawBody(options.Raw)
	}
	if options.RawReader != nil {
		return buildRawReaderBody(options.RawReader)
	}
	if options.JSON != nil {
		return buildJSONBody(options.JSON)
	}
	if len(options.Form) > 0 || len(options.Files) > 0 {
		return buildMultipartBody(options.Form, options.Files)
	}
	if len(options.Data) > 0 {
		return buildURLEncodedBody(options.Data)
	}
	return nil, "", nil
}

// buildMultipartBody constructs a multipart/form-data request body from form fields and files.
// Returns the body as an io.Reader, the content type, and any error encountered.
func buildMultipartBody(form map[string]string, files map[string]string) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range form {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", fmt.Errorf("failed to write form field %s: %w", key, err)
		}
	}

	// Add files
	for fieldName, filePath := range files {
		if err := addFileToMultipart(writer, fieldName, filePath); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}

// buildURLEncodedBody constructs an application/x-www-form-urlencoded request body from form data.
// Returns the body as an io.Reader, the content type, and any error encountered.
func buildURLEncodedBody(data map[string]string) (io.Reader, string, error) {
	values := url.Values{}
	for key, value := range data {
		values.Set(key, value)
	}
	encoded := values.Encode()
	return strings.NewReader(encoded), "application/x-www-form-urlencoded", nil
}

// addFileToMultipart adds a file to a multipart.Writer with the given field name and file path.
// Returns any error encountered during file reading or writing.
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

	if _, errCopy := io.Copy(part, file); errCopy != nil {
		return fmt.Errorf("failed to copy file content: %w", errCopy)
	}

	return nil
}
