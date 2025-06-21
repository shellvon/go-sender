package email

import (
	"errors"
	"net/mail"

	"github.com/shellvon/go-sender/core"
)

// Message represents an email message structure
// Retains the From field, supports multiple recipients, CC, BCC, subject, body, HTML flag, and attachments
// ProviderType method is used for Sender routing
// Validate method is used for parameter validation

type Message struct {
	core.DefaultMessage
	From        string   // Sender's email address
	To          []string // List of recipient email addresses
	Cc          []string // List of CC email addresses
	Bcc         []string // List of BCC email addresses
	Subject     string   // Email subject
	IsHTML      bool     // Indicates if the body is HTML content
	Body        string   // Email body content
	Attachments []string // List of attachment file paths
	AccountName string   // Specific account name to use
}

var (
	_ core.Message = (*Message)(nil)
)

// ProviderType returns the provider type for this message.
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeEmail
}

// Validate checks the validity of the Message fields
func (m *Message) Validate() error {
	if len(m.To) == 0 {
		return core.NewParamError("recipient list cannot be empty")
	}
	for _, email := range m.To {
		if err := validateEmail(email); err != nil {
			return core.NewParamError("invalid recipient email: " + err.Error())
		}
	}
	for _, email := range m.Cc {
		if err := validateEmail(email); err != nil {
			return core.NewParamError("invalid CC email: " + err.Error())
		}
	}
	for _, email := range m.Bcc {
		if err := validateEmail(email); err != nil {
			return core.NewParamError("invalid BCC email: " + err.Error())
		}
	}
	if m.Subject == "" {
		return core.NewParamError("email subject cannot be empty")
	}
	if m.Body == "" {
		return core.NewParamError("email body cannot be empty")
	}
	return nil
}

// validateEmail checks if an email address is valid
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email address cannot be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format: " + err.Error())
	}
	return nil
}

// MsgID returns the unique id of the message.
func (m *Message) MsgID() string {
	return m.DefaultMessage.MsgID()
}
