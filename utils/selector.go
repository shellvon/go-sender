//revive:disable:var-naming
package utils

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
