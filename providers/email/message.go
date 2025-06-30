package email

import (
	"errors"
	"net/mail"

	"github.com/shellvon/go-sender/core"
)

// Message represents an email message structure
// Retains the From field, supports multiple recipients, CC, BCC, subject, body, HTML flag, and attachments
// ProviderType method is used for Sender routing
// Validate method is used for parameter validation.
type Message struct {
	core.DefaultMessage

	From        string   // Sender's email address (supports "Name <address>" format)
	To          []string // List of recipient email addresses (supports "Name <address>" format)
	Cc          []string // List of CC email addresses (supports "Name <address>" format)
	Bcc         []string // List of BCC email addresses (supports "Name <address>" format)
	ReplyTo     string   // Reply-to email address (supports "Name <address>" format)
	Subject     string   // Email subject
	IsHTML      bool     // Indicates if the body is HTML content
	Body        string   // Email body content
	Attachments []string // List of attachment file paths
}

// MessageOption is a function that configures a Message.
type MessageOption func(*Message)

var (
	_ core.Message = (*Message)(nil)
)

// NewMessage creates a new email message with required to and body, plus optional configurations.
func NewMessage(to []string, body string, opts ...MessageOption) *Message {
	msg := &Message{
		To:   to,
		Body: body,
	}
	// Apply optional configurations
	for _, opt := range opts {
		if opt != nil {
			opt(msg)
		}
	}
	return msg
}

// ProviderType returns the provider type for this message.
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeEmail
}

// Validate checks the validity of the Message fields.
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
	if m.Body == "" {
		return core.NewParamError("email body cannot be empty")
	}
	// Subject is optional, no validation needed
	return nil
}

// validateEmail checks if an email address is valid.
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

// MessageOption functions

// WithFrom sets the sender email address.
func WithFrom(from string) MessageOption {
	return func(m *Message) {
		m.From = from
	}
}

// WithSubject sets the email subject.
func WithSubject(subject string) MessageOption {
	return func(m *Message) {
		m.Subject = subject
	}
}

// WithCc sets the CC email addresses.
func WithCc(cc ...string) MessageOption {
	return func(m *Message) {
		m.Cc = cc
	}
}

// WithBcc sets the BCC email addresses.
func WithBcc(bcc ...string) MessageOption {
	return func(m *Message) {
		m.Bcc = bcc
	}
}

// WithReplyTo sets the Reply-To email address.
func WithReplyTo(replyTo string) MessageOption {
	return func(m *Message) {
		m.ReplyTo = replyTo
	}
}

// WithHTML marks the email as HTML content.
func WithHTML() MessageOption {
	return func(m *Message) {
		m.IsHTML = true
	}
}

// WithAttachments sets the attachment file paths.
func WithAttachments(attachments ...string) MessageOption {
	return func(m *Message) {
		m.Attachments = attachments
	}
}

// MsgID returns the unique id of the message.
func (m *Message) MsgID() string {
	return m.DefaultMessage.MsgID()
}
