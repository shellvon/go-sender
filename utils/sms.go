//revive:disable:var-naming
package utils

import (
	"strings"
)

// GetSignatureFromContent extracts the signature from the SMS content.
// Returns the signature content, or an empty string if no signature is found.
// Signature format: starts with 【, ends with 】, and 】 is in a reasonable position (within the first 20 characters).
func GetSignatureFromContent(content string) string {
	if content == "" {
		return ""
	}

	// Check if it starts with 【
	if !strings.HasPrefix(content, "【") {
		return ""
	}

	// Find the position of 】
	endIndex := strings.Index(content, "】")
	if endIndex == -1 || endIndex > 20 || endIndex == 1 {
		return ""
	}

	// Extract the signature content (remove 【 and 】)
	return content[1:endIndex]
}

// HasSignature checks if the content already contains a signature.
func HasSignature(content string) bool {
	return GetSignatureFromContent(content) != ""
}

// AddSignature adds the signature to the content.
//   - If the content already contains a signature, it returns the original content. Otherwise, it adds the signature at the beginning of the content.
//   - If the signName is empty, it returns the original content.
//
// For example, if the content is "Hello, world!" and the signName is "Test", it returns "【Test】Hello, world!".
func AddSignature(content, signName string) string {
	if HasSignature(content) || signName == "" {
		return content
	}

	return "【" + signName + "】" + content
}
