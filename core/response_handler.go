package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// ---------------- New minimal config -------------------------------------

// MatchMode defines comparison types.
type MatchMode string

const (
	MatchEq       MatchMode = "eq"
	MatchContains MatchMode = "contains"
	MatchRegex    MatchMode = "regex"
	MatchGt       MatchMode = "gt"
	MatchGte      MatchMode = "gte"
	MatchLt       MatchMode = "lt"
	MatchLte      MatchMode = "lte"
)

// MatchRule defines how to evaluate success.
type MatchRule struct {
	Path  []string  `json:"path,omitempty"` // e.g. ["Response","Data","0","Code"]. Empty means whole body.
	Mode  MatchMode `json:"mode,omitempty"` // default eq
	Value string    `json:"value,omitempty"`
}

// ResponseHandlerConfig – simplified version.
type ResponseHandlerConfig struct {
	AcceptStatus []int    `json:"accept_status,omitempty"`
	CheckBody    bool     `json:"check_body,omitempty"`
	BodyType     BodyType `json:"body_type,omitempty"`

	Path   string    `json:"path"`
	Expect string    `json:"expect"`
	Mode   MatchMode `json:"mode,omitempty"`

	CodePath string            `json:"code_path,omitempty"`
	MsgPath  string            `json:"msg_path,omitempty"`
	CodeMap  map[string]string `json:"code_map,omitempty"`
}

// ---------------- Core handler -------------------------------------------

// NewSendResultHandler builds a SendResultHandler based on cfg.
//

func NewSendResultHandler(cfg *ResponseHandlerConfig) SendResultHandler {
	if cfg == nil {
		cfg = &ResponseHandlerConfig{CheckBody: false}
	}

	// If CodePath is not specified but Path is, use Path as the default value for CodePath
	if cfg.CheckBody && cfg.CodePath == "" && cfg.Path != "" {
		cfg.CodePath = cfg.Path
	}

	return func(result *SendResult) error {
		if result == nil {
			return errors.New("result is nil")
		}

		bodyBytes := result.Body

		// 1. HTTP status
		if !isStatusOK(cfg.AcceptStatus, result.StatusCode) {
			return fmt.Errorf("HTTP status %d not acceptable", result.StatusCode)
		}
		if !cfg.CheckBody {
			return nil
		}

		// Determine body type
		bType := cfg.BodyType
		if bType == BodyTypeNone {
			ct := ""
			if result.Headers != nil {
				ct = result.Headers.Get("Content-Type")
			}
			bType = detectBodyType(ct)
		}

		success, evalErr := evaluateSuccessSimple(cfg, bType, bodyBytes)
		if evalErr != nil {
			return evalErr
		}
		if success {
			return nil
		}

		// Failure – build error message
		code := extractValue(splitDotPath(cfg.CodePath), bType, bodyBytes)
		msg := extractValue(splitDotPath(cfg.MsgPath), bType, bodyBytes)
		if mapped, ok := cfg.CodeMap[code]; ok {
			msg = mapped
		}
		if msg == "" {
			msg = "unknown error"
		}
		return fmt.Errorf("api error: %s (code=%s)", msg, code)
	}
}

func isStatusOK(white []int, code int) bool {
	if len(white) > 0 {
		return slices.Contains(white, code)
	}
	return code >= http.StatusOK && code < http.StatusBadRequest
}

func evaluateSuccessSimple(cfg *ResponseHandlerConfig, bType BodyType, body []byte) (bool, error) {
	val := extractValue(splitDotPath(cfg.Path), bType, body)
	mode := cfg.Mode
	switch mode {
	case "", MatchEq:
		return val == cfg.Expect, nil
	case MatchContains:
		return strings.Contains(val, cfg.Expect), nil
	case MatchRegex:
		matched, err := regexp.MatchString(cfg.Expect, val)
		return matched, err
	case MatchGt, MatchGte, MatchLt, MatchLte:
		f1, err1 := strconv.ParseFloat(val, 64)
		f2, err2 := strconv.ParseFloat(cfg.Expect, 64)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("non-numeric compare for %s", mode)
		}
		//nolint:exhaustive // only 4 cases.
		switch mode {
		case MatchGt:
			return f1 > f2, nil
		case MatchGte:
			return f1 >= f2, nil
		case MatchLt:
			return f1 < f2, nil
		case MatchLte:
			return f1 <= f2, nil
		}
	default:
		return false, fmt.Errorf("unknown match mode %s", mode)
	}
	return false, nil
}

// extractValue returns string per path.
func extractValue(path []string, bType BodyType, body []byte) string {
	if len(path) == 0 {
		return string(body)
	}
	switch bType {
	case BodyTypeJSON:
		return extractJSONPath(body, path)
	case BodyTypeXML:
		return extractXMLPath(body, path)
	case BodyTypeText:
		return extractTextRegex(body, path)
	case BodyTypeNone:
		return string(body)
	case BodyTypeRaw:
		return string(body)
	case BodyTypeForm:
		return string(body)
	default:
		return ""
	}
}

// Very simple dot/array JSON path extractor.
//
//nolint:gocognit // simple implementation.
func extractJSONPath(body []byte, path []string) string {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return ""
	}

	curr := data

	for _, seg := range path {
		if seg == "" {
			continue
		}

		// helper to drill into map
		drillMap := func(key string) bool {
			if m, ok := curr.(map[string]interface{}); ok {
				curr = m[key]
				return true
			}
			return false
		}

		// iterate through segment parts: key and [index]...
		i := 0
		for i < len(seg) {
			switch seg[i] {
			case '[':
				// array index token here (no preceding key)
				end := strings.IndexByte(seg[i:], ']')
				if end == -1 {
					return ""
				}
				idxStr := seg[i+1 : i+end]
				idx, err := strconv.Atoi(idxStr)
				if err != nil {
					return ""
				}
				if arr, ok := curr.([]interface{}); ok {
					if idx < 0 || idx >= len(arr) {
						return ""
					}
					curr = arr[idx]
				} else {
					return ""
				}
				i += end + 1 // move past ]
			default:
				// parse key until next '[' or end
				j := i
				for j < len(seg) && seg[j] != '[' {
					j++
				}
				key := seg[i:j]
				if !drillMap(key) {
					// when current is array and key is digits (e.g., "0")
					if arr, ok := curr.([]interface{}); ok {
						idx, err := strconv.Atoi(key)
						if err != nil || idx < 0 || idx >= len(arr) {
							return ""
						}
						curr = arr[idx]
					} else {
						return ""
					}
				}
				i = j
			}
		}
	}
	return fmt.Sprintf("%v", curr)
}

func extractXMLPath(body []byte, path []string) string {
	seg := string(body)

	for _, raw := range path {
		if raw == "" {
			continue
		}

		// support item[1] syntax
		tag := raw
		idx := 0
		if strings.HasSuffix(raw, "]") {
			if l := strings.Index(raw, "["); l != -1 {
				tag = raw[:l]
				idxStr := raw[l+1 : len(raw)-1]
				if v, err := strconv.Atoi(idxStr); err == nil {
					idx = v
				}
			}
		}

		openTag := "<" + tag + ">"
		closeTag := "</" + tag + ">"

		// iterate to the idx-th occurrence
		searchPos := 0
		for i := 0; i <= idx; i++ {
			pos := strings.Index(seg[searchPos:], openTag)
			if pos == -1 {
				return ""
			}
			searchPos += pos + len(openTag)
		}

		end := strings.Index(seg[searchPos:], closeTag)
		if end == -1 {
			return ""
		}

		seg = seg[searchPos : searchPos+end]
	}

	return seg
}

func extractTextRegex(body []byte, pattern []string) string {
	if len(pattern) == 0 || (len(pattern) == 1 && pattern[0] == "") {
		return string(body)
	}
	re := regexp.MustCompile(strings.Join(pattern, ""))
	match := re.FindSubmatch(body)
	if len(match) > 1 {
		return string(match[1])
	}
	// if no capture group, but matches, return whole
	if len(match) > 0 {
		return string(match[0])
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
	default:
		return BodyTypeText
	}
}

// splitDotPath parses a dot-notation path into segments, supporting escaped '\\.' .
func splitDotPath(p string) []string {
	if p == "" {
		return nil
	}
	var segs []string
	var curr strings.Builder
	esc := false
	for _, r := range p {
		switch {
		case esc:
			curr.WriteRune(r)
			esc = false
		case r == '\\':
			esc = true
		case r == '.':
			segs = append(segs, curr.String())
			curr.Reset()
		default:
			curr.WriteRune(r)
		}
	}
	segs = append(segs, curr.String())
	return segs
}
