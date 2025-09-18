package core

import (
	"strconv"
	"strings"
)

// ExtraFieldsAware is an optional interface for messages that support extra fields.
// This interface is primarily used by providers that need additional configuration
// parameters that are specific to certain sub-providers or advanced features.
//
// For example
//   - SMS provider: for provider-specific parameters (region, endpoint, etc.)
//   - EmailAPI provider: for API-specific configurations
type ExtraFieldsAware interface {
	// GetExtra returns a value from extras by key.
	GetExtra(key string) (interface{}, bool)

	// SetExtra sets a value in extras by key.
	SetExtra(key string, value interface{})

	// GetExtraString returns a string value from extras.
	GetExtraString(key string) (string, bool)

	// GetExtraStringOrDefault returns a string value from extras, or the default value if not found.
	GetExtraStringOrDefault(key, defaultValue string) string

	// GetExtraStringOrDefaultEmpty returns a string value from extras, or empty string if not found.
	GetExtraStringOrDefaultEmpty(key string) string

	// GetExtraInt returns an int value from extras.
	GetExtraInt(key string) (int, bool)

	// GetExtraIntOrDefault returns an int value from extras, or the default value if not found.
	GetExtraIntOrDefault(key string, defaultValue int) int

	// GetExtraBool returns a bool value from extras.
	GetExtraBool(key string) (bool, bool)

	// GetExtraBoolOrDefault returns a bool value from extras, or the default value if not found.
	GetExtraBoolOrDefault(key string, defaultValue bool) bool
}

// WithExtraFields provides a concrete implementation of ExtraFieldsAware.
// Embed this struct in your message types when you need extra fields functionality.
//
// Example:
//
//	type MyMessage struct {
//	    *core.BaseMessage
//	    *core.WithExtraFields  // Add extra fields capability
//	    // other fields...
//	}
type WithExtraFields struct {
	Extras map[string]interface{} `json:"extras,omitempty"`
}

// NewWithExtraFields creates a new WithExtraFields instance.
func NewWithExtraFields() *WithExtraFields {
	return &WithExtraFields{
		Extras: make(map[string]interface{}),
	}
}

// GetExtra returns a value from extras by key.
func (w *WithExtraFields) GetExtra(key string) (interface{}, bool) {
	if w.Extras == nil {
		return nil, false
	}
	value, exists := w.Extras[key]
	return value, exists
}

// SetExtra sets a value in extras by key.
func (w *WithExtraFields) SetExtra(key string, value interface{}) {
	if w.Extras == nil {
		w.Extras = make(map[string]interface{})
	}
	w.Extras[key] = value
}

// GetExtraString returns a string value from extras.
func (w *WithExtraFields) GetExtraString(key string) (string, bool) {
	if value, exists := w.GetExtra(key); exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetExtraStringOrDefault returns a string value from extras, or the default value if not found.
func (w *WithExtraFields) GetExtraStringOrDefault(key, defaultValue string) string {
	if value, exists := w.GetExtraString(key); exists {
		return value
	}
	return defaultValue
}

// GetExtraStringOrDefaultEmpty returns a string value from extras, or empty string if not found.
func (w *WithExtraFields) GetExtraStringOrDefaultEmpty(key string) string {
	return w.GetExtraStringOrDefault(key, "")
}

// GetExtraInt returns an int value from extras.
func (w *WithExtraFields) GetExtraInt(key string) (int, bool) {
	if value, exists := w.GetExtra(key); exists {
		switch v := value.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i, true
			}
		}
	}
	return 0, false
}

// GetExtraIntOrDefault returns an int value from extras, or the default value if not found.
func (w *WithExtraFields) GetExtraIntOrDefault(key string, defaultValue int) int {
	if value, exists := w.GetExtraInt(key); exists {
		return value
	}
	return defaultValue
}

// GetExtraBool returns a bool value from extras.
func (w *WithExtraFields) GetExtraBool(key string) (bool, bool) {
	if value, exists := w.GetExtra(key); exists {
		switch v := value.(type) {
		case bool:
			return v, true
		case string:
			return strings.ToLower(v) == "true", true
		case int:
			return v != 0, true
		}
	}
	return false, false
}

// GetExtraBoolOrDefault returns a bool value from extras, or the default value if not found.
func (w *WithExtraFields) GetExtraBoolOrDefault(key string, defaultValue bool) bool {
	if value, exists := w.GetExtraBool(key); exists {
		return value
	}
	return defaultValue
}

// Compile-time assertion: WithExtraFields implements ExtraFieldsAware interface.
var _ ExtraFieldsAware = (*WithExtraFields)(nil)
