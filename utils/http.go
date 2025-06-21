package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
	// For a default timeout
)

// RequestOptions holds the parameters for your HTTP request.
type RequestOptions struct {
	Method      string // e.g., http.MethodGet, http.MethodPost
	Headers     map[string]string
	Body        []byte    // Raw byte slice for simple bodies (e.g., JSON, text)
	BodyReader  io.Reader // For custom readers, like file streams or multipart writers
	ContentType string    // Explicit Content-Type for the request body

	// New fields for file uploads (optional, for convenience)
	FilePath      string // Path to a file to upload directly
	FileName      string // Optional: Override the filename sent in multipart
	FileFieldName string // Optional: Field name for the file in multipart form (default "file")
	Timeout       time.Duration
}

// DoRequest performs an HTTP request and returns the response body.
// It handles common HTTP client concerns like context, headers, and body reading.
//
// Parameters:
//
//	options: A RequestOptions struct containing all request details.
//
// Returns:
//
//	[]byte: The raw response body. If the request was not successful (non-200 status code)
//	        or an error occurred, this might still contain the server's error message.
//	int: The HTTP status code of the response.
//	error: An error if the HTTP request failed, if the body couldn't be read,
//	       or if the server returned a non-200 status code. The error will include
//	       details from the response body if available.
func DoRequest(ctx context.Context, url string, options RequestOptions) ([]byte, int, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	var reqBody io.Reader
	var finalContentType string

	if options.BodyReader != nil {
		reqBody = options.BodyReader
		finalContentType = options.ContentType
	} else if options.FilePath != "" {
		fileInfo, err := os.Stat(options.FilePath)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to stat file %s: %w", options.FilePath, err)
		}
		if fileInfo.IsDir() {
			return nil, 0, fmt.Errorf("path %s is a directory, not a file", options.FilePath)
		}
		file, err := os.Open(options.FilePath)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to open file %s: %w", options.FilePath, err)
		}
		defer file.Close()

		if options.ContentType == "application/octet-stream" {
			reqBody = file
			finalContentType = "application/octet-stream"
		} else {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			fileFieldName := options.FileFieldName
			if fileFieldName == "" {
				fileFieldName = "file"
			}

			// Determine the filename to use in the multipart form part
			fileNameToUse := options.FileName
			if fileNameToUse == "" {
				fileNameToUse = filepath.Base(options.FilePath)
			}

			// Create the form file part once with the determined filename
			part, err := writer.CreateFormFile(fileFieldName, fileNameToUse)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to create form file part: %w", err)
			}

			_, err = io.Copy(part, file)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to copy file content to form part: %w", err)
			}

			err = writer.Close()
			if err != nil {
				return nil, 0, fmt.Errorf("failed to close multipart writer: %w", err)
			}

			reqBody = body
			finalContentType = writer.FormDataContentType()
		}
	} else if options.Body != nil {
		reqBody = bytes.NewReader(options.Body)
		finalContentType = options.ContentType
		if finalContentType == "" {
			finalContentType = "application/json"
		}
	}

	if reqBody != nil && options.Method == "" {
		options.Method = http.MethodPost
	}

	req, err := http.NewRequestWithContext(ctx, options.Method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	if finalContentType != "" {
		req.Header.Set("Content-Type", finalContentType)
	}

	// Perform the request using the default client
	// For production, you might want a custom http.Client with specific settings
	// like timeouts, TLS configuration, etc.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Always defer closing the response body to ensure resources are released.
	defer resp.Body.Close()

	// Always read the response body completely, regardless of the status code.
	// This helps with connection reuse and provides error details.
	// From offical docs:
	//  > The default HTTP client's Transport may not reuse HTTP/1.x "keep-alive" TCP connections if the Body is not read to completion and closed.
	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", readErr)
	}

	// Check if the status code indicates an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return bodyBytes, resp.StatusCode, fmt.Errorf("HTTP request failed with status code %d. Response body: %s", resp.StatusCode, string(bodyBytes))
	}

	// If everything is successful, return the body bytes and status code
	return bodyBytes, resp.StatusCode, nil
}
