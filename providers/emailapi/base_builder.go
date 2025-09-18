package emailapi

import (
	"time"

	"github.com/shellvon/go-sender/core"
)

// BaseBuilder provides chainable methods for common EmailAPI message fields.
type BaseBuilder[T any] struct {
	to           []string
	cc           []string
	bcc          []string
	from         string
	replyTo      []string
	subject      string
	text         string
	html         string
	attachments  []Attachment
	headers      map[string]string
	callbackURL  string
	extras       map[string]interface{}
	templateID   string
	scheduledAt  *time.Time
	templateData map[string]interface{}
	self         T
}

func (b *BaseBuilder[T]) To(to ...string) T {
	b.to = to
	return b.self
}

func (b *BaseBuilder[T]) Cc(cc ...string) T {
	b.cc = cc
	return b.self
}

func (b *BaseBuilder[T]) Bcc(bcc ...string) T {
	b.bcc = bcc
	return b.self
}

func (b *BaseBuilder[T]) From(from string) T {
	b.from = from
	return b.self
}

func (b *BaseBuilder[T]) ReplyTo(replyTo ...string) T {
	b.replyTo = replyTo
	return b.self
}

func (b *BaseBuilder[T]) Subject(subject string) T {
	b.subject = subject
	return b.self
}

func (b *BaseBuilder[T]) Text(text string) T {
	b.text = text
	return b.self
}

func (b *BaseBuilder[T]) HTML(html string) T {
	b.html = html
	return b.self
}

func (b *BaseBuilder[T]) Attachments(attachments []Attachment) T {
	b.attachments = attachments
	return b.self
}

func (b *BaseBuilder[T]) AddAttachment(attachment Attachment) T {
	b.attachments = append(b.attachments, attachment)
	return b.self
}

func (b *BaseBuilder[T]) Headers(headers map[string]string) T {
	b.headers = headers
	return b.self
}

func (b *BaseBuilder[T]) CallbackURL(url string) T {
	b.callbackURL = url
	return b.self
}

func (b *BaseBuilder[T]) Extras(extras map[string]interface{}) T {
	b.extras = extras
	return b.self
}

func (b *BaseBuilder[T]) TemplateID(templateID string) T {
	b.templateID = templateID
	return b.self
}

func (b *BaseBuilder[T]) TemplateParams(params map[string]interface{}) T {
	b.templateData = params
	return b.self
}

func (b *BaseBuilder[T]) ScheduledAt(scheduledAt *time.Time) T {
	b.scheduledAt = scheduledAt
	return b.self
}

// BuildMessage constructs a *Message from the builder fields.
func (b *BaseBuilder[T]) BuildMessage(subProvider string) *Message {
	msg := &Message{
		BaseMessage:     core.NewBaseMessage(core.ProviderTypeEmailAPI),
		WithExtraFields: core.NewWithExtraFields(),
		SubProvider:     subProvider,
		To:              b.to,
		Cc:              b.cc,
		Bcc:             b.bcc,
		From:            b.from,
		ReplyTo:         b.replyTo,
		Subject:         b.subject,
		Text:            b.text,
		HTML:            b.html,
		Attachments:     b.attachments,
		Headers:         b.headers,
		CallbackURL:     b.callbackURL,
		TemplateID:      b.templateID,
		TemplateData:    b.templateData,
		ScheduledAt:     b.scheduledAt,
	}
	// Copy extras from builder to message
	if b.extras != nil {
		msg.Extras = b.extras
	}
	return msg
}
