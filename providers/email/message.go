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

var (
	_ core.Message = (*Message)(nil)
)

// NewMessage creates a new email message with required fields only.
func NewMessage(to []string, body string) *Message {
	return &Message{
		To:   to,
		Body: body,
	}
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

// --- Builder Pattern for Email ---

// EmailBuilder is used for building email messages in a chainable style.
type EmailBuilder struct {
	recipients  []string // Recipients
	subject     string   // Subject
	body        string   // Body
	from        string   // Sender
	cc          []string // CC
	bcc         []string // BCC
	replyTo     string   // Reply-To
	html        bool     // Is HTML
	attachments []string // Attachments
}

// Email creates a new EmailBuilder as the entry point for builder style.
func Email() *EmailBuilder {
	return &EmailBuilder{}
}

// To sets the recipients.
func (b *EmailBuilder) To(recipients ...string) *EmailBuilder {
	b.recipients = append(b.recipients, recipients...)
	return b
}

// Subject sets the subject.
func (b *EmailBuilder) Subject(subject string) *EmailBuilder {
	b.subject = subject
	return b
}

// Body sets the body.
func (b *EmailBuilder) Body(body string) *EmailBuilder {
	b.body = body
	return b
}

// From sets the sender.
func (b *EmailBuilder) From(from string) *EmailBuilder {
	b.from = from
	return b
}

// Cc sets the CC recipients.
func (b *EmailBuilder) Cc(cc ...string) *EmailBuilder {
	b.cc = append(b.cc, cc...)
	return b
}

// Bcc sets the BCC recipients.
func (b *EmailBuilder) Bcc(bcc ...string) *EmailBuilder {
	b.bcc = append(b.bcc, bcc...)
	return b
}

// ReplyTo sets the reply-to address.
func (b *EmailBuilder) ReplyTo(replyTo string) *EmailBuilder {
	b.replyTo = replyTo
	return b
}

// HTML marks the message as HTML content.
func (b *EmailBuilder) HTML() *EmailBuilder {
	b.html = true
	return b
}

// Attach replaces the attachment list.
func (b *EmailBuilder) Attach(files ...string) *EmailBuilder {
	b.attachments = files
	return b
}

// AddAttach appends files to the attachment list.
func (b *EmailBuilder) AddAttach(files ...string) *EmailBuilder {
	b.attachments = append(b.attachments, files...)
	return b
}

// Build creates the Message instance from the builder.
func (b *EmailBuilder) Build() *Message {
	msg := NewMessage(b.recipients, b.body)
	if b.subject != "" {
		msg.Subject = b.subject
	}
	if b.from != "" {
		msg.From = b.from
	}
	if len(b.cc) > 0 {
		msg.Cc = b.cc
	}
	if len(b.bcc) > 0 {
		msg.Bcc = b.bcc
	}
	if b.replyTo != "" {
		msg.ReplyTo = b.replyTo
	}
	if b.html {
		msg.IsHTML = true
	}
	if len(b.attachments) > 0 {
		msg.Attachments = b.attachments
	}
	return msg
}
