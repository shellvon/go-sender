package telegram

import "github.com/shellvon/go-sender/core"

// ContactMessage represents a contact message for Telegram
// Based on SendContactParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendcontact
type ContactMessage struct {
	BaseMessage

	// Contact's phone number
	PhoneNumber string `json:"phone_number"`

	// Contact's first name
	FirstName string `json:"first_name"`

	// Contact's last name
	LastName string `json:"last_name,omitempty"`

	// Additional data about the contact in the form of a vCard, 0-2048 bytes
	VCard string `json:"vcard,omitempty"`
}

func (m *ContactMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *ContactMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *ContactMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.PhoneNumber == "" {
		return core.NewParamError("phone_number cannot be empty")
	}
	if m.FirstName == "" {
		return core.NewParamError("first_name cannot be empty")
	}
	return nil
}

type ContactMessageOption func(*ContactMessage)

// WithContactLastName sets the last name of the contact
// This is optional and can be omitted
func WithContactLastName(lastName string) ContactMessageOption {
	return func(m *ContactMessage) { m.LastName = lastName }
}

// WithContactVCard sets additional data about the contact in vCard format
// Should be 0-2048 bytes in size
func WithContactVCard(vCard string) ContactMessageOption {
	return func(m *ContactMessage) { m.VCard = vCard }
}

func NewContactMessage(chatID string, phoneNumber, firstName string, opts ...interface{}) *ContactMessage {
	msg := &ContactMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeContact,
			ChatID:  chatID,
		},
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case ContactMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}
