package serverchan

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// Message for Server酱
type Message struct {
	core.DefaultMessage
	Title   string `json:"title"`   // Message title
	Content string `json:"content"` // Message content (supports Markdown)
}

func NewMessage(title, content string) *Message {
	return &Message{
		Title:   title,
		Content: content,
	}
}

func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeServerChan
}

func (m *Message) Validate() error {
	if m.Title == "" {
		return errors.New("title cannot be empty")
	}
	if m.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
