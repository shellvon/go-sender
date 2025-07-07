//revive:disable:var-naming
package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

	// Query string parameters (url.Values allows single or multiple values)
	Query url.Values // Query string parameters

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
func buildFinalURL(requestURL string, query url.Values) (string, error) {
	if len(query) == 0 {
		return requestURL, nil
	}
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}
	q := parsedURL.Query()
	for key, values := range query {
		for _, val := range values {
			q.Add(key, val)
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
	// 1. apply user headers first (they have highest priority)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 2. ensure default User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", core.DefaultUserAgent)
	}

	// 3. set Content-Type only if not already provided
	if contentType != "" && req.Header.Get("Content-Type") == "" {
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

// SendRequest sends an HTTP request and returns the raw *http.Response.
// The caller is responsible for closing resp.Body (or delegating to ReadAndClose).
func SendRequest(ctx context.Context, requestURL string, options HTTPRequestOptions) (*http.Response, error) {
	ctx, cancel := prepareContext(ctx, &options)
	if cancel != nil {
		defer cancel()
	}

	reqBody, contentType, err := buildRequestBody(options)
	if err != nil {
		return nil, fmt.Errorf("failed to build request body: %w", err)
	}

	finalURL, err := buildFinalURL(requestURL, options.Query)
	if err != nil {
		return nil, err
	}

	method := getDefaultMethod(options.Method, reqBody)
	req, err := http.NewRequestWithContext(ctx, method, finalURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	setRequestHeaders(req, options.Headers, contentType)
	client := core.EnsureHTTPClient(options.Client)
	return client.Do(req)
}

// ReadAndClose reads all bytes from resp.Body and then closes it.
// It does NOT perform any status code validation; callers should decide
// how to interpret resp.StatusCode. Returns the body bytes and response headers.
func ReadAndClose(resp *http.Response) ([]byte, http.Header, error) {
	if resp == nil {
		return nil, nil, errors.New("response is nil")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.Header, fmt.Errorf("failed to read response body: %w", err)
	}
	return bodyBytes, resp.Header, nil
}

func DoRequest(ctx context.Context, requestURL string, options HTTPRequestOptions) ([]byte, int, error) {
	resp, err := SendRequest(ctx, requestURL, options)
	if err != nil {
		return nil, 0, err
	}
	body, _, readErr := ReadAndClose(resp)
	return body, resp.StatusCode, readErr
}
