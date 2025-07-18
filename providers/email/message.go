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
	return Email().
		To(to...).
		Body(body).
		Build()
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

// Builder is used for building email messages in a chainable style.
type Builder struct {
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

// Email creates a new Builder as the entry point for builder style.
func Email() *Builder {
	return &Builder{}
}

// To sets the recipients.
func (b *Builder) To(recipients ...string) *Builder {
	b.recipients = append(b.recipients, recipients...)
	return b
}

// Subject sets the subject.
func (b *Builder) Subject(subject string) *Builder {
	b.subject = subject
	return b
}

// Body sets the body.
func (b *Builder) Body(body string) *Builder {
	b.body = body
	return b
}

// From sets the sender.
func (b *Builder) From(from string) *Builder {
	b.from = from
	return b
}

// Cc sets the CC recipients.
func (b *Builder) Cc(cc ...string) *Builder {
	b.cc = append(b.cc, cc...)
	return b
}

// Bcc sets the BCC recipients.
func (b *Builder) Bcc(bcc ...string) *Builder {
	b.bcc = append(b.bcc, bcc...)
	return b
}

// ReplyTo sets the reply-to address.
func (b *Builder) ReplyTo(replyTo string) *Builder {
	b.replyTo = replyTo
	return b
}

// HTML marks the message as HTML content.
func (b *Builder) HTML() *Builder {
	b.html = true
	return b
}

// Attach replaces the attachment list.
func (b *Builder) Attach(files ...string) *Builder {
	b.attachments = files
	return b
}

// AddAttach appends files to the attachment list.
func (b *Builder) AddAttach(files ...string) *Builder {
	b.attachments = append(b.attachments, files...)
	return b
}

// Build creates the Message instance from the builder.
func (b *Builder) Build() *Message {
	msg := &Message{
		To:          b.recipients,
		Body:        b.body,
		Subject:     b.subject,
		From:        b.from,
		Cc:          b.cc,
		Bcc:         b.bcc,
		ReplyTo:     b.replyTo,
		IsHTML:      b.html,
		Attachments: b.attachments,
	}
	return msg
}
