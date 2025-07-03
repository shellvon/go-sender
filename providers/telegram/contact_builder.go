//nolint:dupl // builder pattern files share similar boilerplate; acceptable duplication
package telegram

// ContactBuilder constructs Telegram contact messages.
// Example:
//   msg := telegram.Contact().
//            Chat("123").
//            Phone("+8613800138000").
//            FirstName("Alice").
//            LastName("Smith").
//            Build()

type ContactBuilder struct {
	*baseBuilder[*ContactBuilder]

	phoneNumber string
	firstName   string
	lastName    string
	vcard       string
}

// Contact returns a new ContactBuilder.
func Contact() *ContactBuilder {
	b := &ContactBuilder{}
	b.baseBuilder = &baseBuilder[*ContactBuilder]{self: b}
	return b
}

// Phone sets the contact phone number.
// Based on SendContactParams from Telegram Bot API.
func (b *ContactBuilder) Phone(num string) *ContactBuilder {
	b.phoneNumber = num
	return b
}

// FirstName sets the contact first name.
// Based on SendContactParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendcontact
// Contact's first name.
func (b *ContactBuilder) FirstName(name string) *ContactBuilder {
	b.firstName = name
	return b
}

// LastName sets the contact last name.
// Based on SendContactParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendcontact
// Contact's last name.
func (b *ContactBuilder) LastName(name string) *ContactBuilder {
	b.lastName = name
	return b
}

// VCard attaches vCard additional data.
// Based on SendContactParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendcontact
// Additional data about the contact in the form of a vCard, 0-2048 bytes.
func (b *ContactBuilder) VCard(data string) *ContactBuilder {
	b.vcard = data
	return b
}

// Build assembles the *ContactMessage.
func (b *ContactBuilder) Build() *ContactMessage {
	msg := &ContactMessage{
		BaseMessage: b.baseBuilder.toBaseMessage(TypeContact),
		PhoneNumber: b.phoneNumber,
		FirstName:   b.firstName,
		LastName:    b.lastName,
		VCard:       b.vcard,
	}
	return msg
}
