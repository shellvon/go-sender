package core

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	// "github.com/shellvon/go-sender/utils".
)

// ResponseHandlerConfig describes how to validate an HTTP response.
// The same config can be reused by any HTTP-based provider (webhook, SMS, etc.).
//
// Key fields:
//   - SuccessStatusCodes  – Optional whitelist of HTTP status codes considered successful.
//   - ValidateResponse   – If false, only the status code is checked.
//   - ResponseType       – Expected body type (json/text/xml/raw/none/form).
//
// JSON responses (ResponseType == BodyTypeJSON):
//   - SuccessField / SuccessValue – Pair that marks a successful result when equal.
//   - ErrorCodeField              – Field holding a platform-specific error code.
//   - ErrorField / MessageField   – Fields holding the error / info message.
//   - ErrorCodeMap                – Map error code ➜ human-friendly message.
//
// Text responses (ResponseType == BodyTypeText):
//   - SuccessPattern / ErrorPattern – Regexes applied to plain-text body.
//
// You can extend this struct without breaking existing providers; unrecognised
// fields are simply ignored by the validation logic.
type ResponseHandlerConfig struct {
	// Explicit list of HTTP status codes that should be treated as success.
	SuccessStatusCodes []int `json:"success_status_codes,omitempty"`

	// When false only the HTTP status code is validated.
	ValidateResponse bool `json:"validate_response,omitempty"`

	// Expected body format.
	ResponseType BodyType `json:"response_type,omitempty"`

	// JSON-specific fields.
	SuccessField      string            `json:"success_field,omitempty"`
	SuccessValue      string            `json:"success_value,omitempty"`
	ErrorCodeField    string            `json:"error_code_field,omitempty"`
	ErrorField        string            `json:"error_field,omitempty"`   // Kept for webhook backward-compat.
	MessageField      string            `json:"message_field,omitempty"` // Kept for webhook backward-compat.
	ErrorMessageField string            `json:"error_message_field,omitempty"`
	ErrorCodeMap      map[string]string `json:"error_code_map,omitempty"`

	// Text-specific fields.
	SuccessPattern string `json:"success_pattern,omitempty"`
	ErrorPattern   string `json:"error_pattern,omitempty"`
}

// NewResponseHandler builds a ResponseHandler using the supplied configuration.
// If cfg is nil, it returns a ResponseHandler that only validates status code.
func NewResponseHandler(cfg *ResponseHandlerConfig) ResponseHandler {
	if cfg == nil {
		cfg = &ResponseHandlerConfig{ValidateResponse: false}
	}

	return func(resp *http.Response) error {
		if resp == nil {
			return errors.New("response is nil")
		}

		// Read body & close
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		defer resp.Body.Close()
		// 1) Check HTTP status code.
		if !cfg.isStatusOK(resp.StatusCode) {
			return fmt.Errorf("HTTP status %d not acceptable", resp.StatusCode)
		}

		// 2) Optionally validate response body.
		if !cfg.ValidateResponse {
			return nil
		}

		// Determine the effective body type.
		effectiveType := cfg.ResponseType
		if effectiveType == BodyTypeNone {
			effectiveType = detectBodyType(resp.Header.Get("Content-Type"))
		}

		switch effectiveType {
		case BodyTypeJSON:
			return cfg.handleJSON(body)
		case BodyTypeText:
			return cfg.handleText(body)
		case BodyTypeXML:
			return cfg.handleXML(body)
		case BodyTypeForm, BodyTypeRaw:
			return fmt.Errorf("unsupported response body type: %s", effectiveType.String())
		case BodyTypeNone:
			// do nothing
			return nil
		default:
			// Fallback treat as text
			return cfg.handleText(body)
		}
	}
}

// isStatusOK returns true if the given status code is considered successful.
// If SuccessStatusCodes is empty, it returns true for 2xx/3xx status codes.
// otherwise, it returns true if the status code is in SuccessStatusCodes.
func (cfg *ResponseHandlerConfig) isStatusOK(statusCode int) bool {
	if len(cfg.SuccessStatusCodes) > 0 {
		return slices.Contains(cfg.SuccessStatusCodes, statusCode)
	}
	return statusCode >= http.StatusOK && statusCode < http.StatusBadRequest
}

// handleJSON validates a JSON response body according to the config.
func (cfg *ResponseHandlerConfig) handleJSON(body []byte) error {
	if len(body) == 0 {
		return nil
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	if cfg.SuccessField != "" {
		if v, ok := resp[cfg.SuccessField]; ok && fmt.Sprintf("%v", v) == cfg.SuccessValue {
			return nil
		}
	}
	errCode := extractStringField(resp, cfg.ErrorCodeField)
	if errCode == "" {
		errCode = extractStringField(resp, "code")
	}
	errMsg := extractFirstNonEmpty(
		extractStringField(resp, cfg.ErrorField),
		extractStringField(resp, cfg.MessageField),
		extractStringField(resp, cfg.ErrorMessageField),
		"unknown error",
	)
	if cfg.ErrorCodeMap != nil {
		if mapped, ok := cfg.ErrorCodeMap[errCode]; ok {
			errMsg = mapped
		}
	}
	return fmt.Errorf("api returned error: %s (code=%s)", errMsg, errCode)
}

// handleText validates a plain-text response body according to the config.
func (cfg *ResponseHandlerConfig) handleText(body []byte) error {
	text := string(body)
	if cfg.ErrorPattern != "" {
		if matched, _ := regexp.MatchString(cfg.ErrorPattern, text); matched {
			return fmt.Errorf("api returned error response: %s", text)
		}
	}
	if cfg.SuccessPattern != "" {
		if matched, _ := regexp.MatchString(cfg.SuccessPattern, text); !matched {
			return fmt.Errorf("response does not match success pattern: %s", text)
		}
	}
	return nil
}

// extractXMLField returns the character data of the first element whose name matches field.
func extractXMLField(body []byte, field string) string {
	if field == "" {
		return ""
	}
	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		//nolint:gocritic // just for xml decoding
		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == field {
				var val string
				if decodeErr := decoder.DecodeElement(&val, &se); decodeErr == nil {
					return val
				}
			}
		}
	}
	return ""
}

// handleXML validates XML body – currently basic well-formedness check.
func (cfg *ResponseHandlerConfig) handleXML(body []byte) error {
	if len(body) == 0 {
		return nil
	}
	if err := xml.Unmarshal(body, new(interface{})); err != nil {
		return fmt.Errorf("failed to parse XML response: %w", err)
	}

	// Success detection similar to JSON.
	if cfg.SuccessField != "" {
		if val := extractXMLField(body, cfg.SuccessField); val == cfg.SuccessValue {
			return nil
		}
	}

	errCode := extractXMLField(body, cfg.ErrorCodeField)
	errMsg := extractFirstNonEmpty(
		extractXMLField(body, cfg.ErrorField),
		extractXMLField(body, cfg.MessageField),
		extractXMLField(body, cfg.ErrorMessageField),
	)
	if errCode == "" {
		errCode = "UNKNOWN"
	}
	if cfg.ErrorCodeMap != nil {
		if mapped, ok := cfg.ErrorCodeMap[errCode]; ok {
			errMsg = mapped
		}
	}
	if errMsg == "" {
		errMsg = "unknown error"
	}
	return fmt.Errorf("api returned error: %s (code=%s)", errMsg, errCode)
}

// extractStringField returns resp[field] converted to string (empty if missing).
func extractStringField(resp map[string]interface{}, field string) string {
	if field == "" {
		return ""
	}
	if v, ok := resp[field]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// extractFirstNonEmpty returns the first non-empty string.
func extractFirstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// detectBodyType decides body type from Content-Type header.
func detectBodyType(ct string) BodyType {
	ct = strings.ToLower(ct)
	switch {
	case strings.Contains(ct, "application/json") || strings.Contains(ct, "+json"):
		return BodyTypeJSON
	case strings.Contains(ct, "xml"):
		return BodyTypeXML
	case strings.Contains(ct, "text/") || strings.Contains(ct, "plain"):
		return BodyTypeText
	case strings.Contains(ct, "x-www-form-urlencoded"):
		return BodyTypeForm
	case strings.Contains(ct, "octet-stream"):
		return BodyTypeRaw
	default:
		return BodyTypeText
	}
}
