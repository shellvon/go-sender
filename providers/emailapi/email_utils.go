package emailapi

import (
	"net/mail"
	"strings"
)

// EmailAddress represents a parsed email address with name and email parts.
// This is the unified structure used across all email providers to avoid duplication.
type EmailAddress struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

// parseEmailAddress parses an email address string using Go's standard library.
// It supports formats like:
// - "john@example.com"
// - "John Doe <john@example.com>"
// - "\"John Doe\" <john@example.com>".
func parseEmailAddress(addr string) EmailAddress {
	if addr == "" {
		return EmailAddress{}
	}

	// Try to parse using net/mail package
	parsed, err := mail.ParseAddress(addr)
	if err != nil {
		// If parsing fails, assume it's just an email address without name
		return EmailAddress{
			Email: strings.TrimSpace(addr),
			Name:  "",
		}
	}

	return EmailAddress{
		Name:  parsed.Name,
		Email: parsed.Address,
	}
}

// parseEmailAddresses parses multiple email addresses.
func parseEmailAddresses(addrs []string) []EmailAddress {
	if len(addrs) == 0 {
		return nil
	}

	result := make([]EmailAddress, len(addrs))
	for i, addr := range addrs {
		result[i] = parseEmailAddress(addr)
	}
	return result
}
