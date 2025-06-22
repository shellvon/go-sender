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

// NewMessage creates a new email message
func NewMessage(subject, body string) *Message {
	return &Message{
		Subject: subject,
		Body:    body,
	}
}

// WithFrom sets the sender email address
func (m *Message) WithFrom(from string) *Message {
	m.From = from
	return m
}

// WithTo sets the recipient email addresses
func (m *Message) WithTo(to ...string) *Message {
	m.To = to
	return m
}

// WithCc sets the CC email addresses
func (m *Message) WithCc(cc ...string) *Message {
	m.Cc = cc
	return m
}

// WithBcc sets the BCC email addresses
func (m *Message) WithBcc(bcc ...string) *Message {
	m.Bcc = bcc
	return m
}

// WithSubject sets the email subject
func (m *Message) WithSubject(subject string) *Message {
	m.Subject = subject
	return m
}

// WithBody sets the email body
func (m *Message) WithBody(body string) *Message {
	m.Body = body
	return m
}

// WithHTML sets the email as HTML content
func (m *Message) WithHTML() *Message {
	m.IsHTML = true
	return m
}

// WithAttachments sets the attachment file paths
func (m *Message) WithAttachments(attachments ...string) *Message {
	m.Attachments = attachments
	return m
}

// MsgID returns the unique id of the message.
func (m *Message) MsgID() string {
	return m.DefaultMessage.MsgID()
}
