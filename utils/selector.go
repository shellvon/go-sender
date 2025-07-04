package utils

// DefaultStringIfEmpty returns def if s is empty, otherwise returns s.
func DefaultStringIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// FirstNonEmpty returns the first non-empty string from the provided values.
// It follows the priority order: first non-empty value wins.
// This is useful for implementing fallback logic with multiple default values.
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// BuildExtras is a helper function to build extras map from builder fields.
// This helps eliminate code duplication in Build() methods across different providers.
// It filters out empty strings and zero values, only including non-empty fields.
func BuildExtras(fields map[string]interface{}) map[string]interface{} {
	extra := make(map[string]interface{})
	for key, value := range fields {
		if value != "" && value != 0 {
			extra[key] = value
		}
	}
	return extra
}
