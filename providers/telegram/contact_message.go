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

// NewContactMessage creates a new ContactMessage instance.
// Based on SendContactParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendcontact
//   - Only chat_id and phoneNumber/firstName are required.
func NewContactMessage(chatID string, phoneNumber, firstName string) *ContactMessage {
	return &ContactMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeContact,
			ChatID:  chatID,
		},
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
	}
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
