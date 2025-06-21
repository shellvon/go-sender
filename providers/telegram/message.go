package telegram

import (
	"strings"

	"github.com/shellvon/go-sender/core"
)

// MessageType represents the type of Telegram message
type MessageType string

const (
	TypeText     MessageType = "text"
	TypePhoto    MessageType = "photo"
	TypeDocument MessageType = "document"
	TypeLocation MessageType = "location"
	TypeContact  MessageType = "contact"
	TypePoll     MessageType = "poll"
)

// Message interface definition
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// BaseMessage is the base message structure
type BaseMessage struct {
	core.DefaultMessage
	MsgType MessageType `json:"msgtype"`
}

// GetMsgType implements the Message interface
func (m *BaseMessage) GetMsgType() MessageType {
	return m.MsgType
}

// TextMessage represents a text message for Telegram
type TextMessage struct {
	BaseMessage
	ChatID            string `json:"chat_id"`
	Text              string `json:"text"`
	ParseMode         string `json:"parse_mode,omitempty"`
	DisableWebPreview bool   `json:"disable_web_page_preview,omitempty"`
	Silent            bool   `json:"disable_notification,omitempty"`
	ProtectContent    bool   `json:"protect_content,omitempty"`
	ReplyToMsg        int    `json:"reply_to_message_id,omitempty"`
}

// PhotoMessage represents a photo message for Telegram
type PhotoMessage struct {
	BaseMessage
	ChatID         string `json:"chat_id"`
	Photo          string `json:"photo"` // File path, URL, or file_id
	Caption        string `json:"caption,omitempty"`
	ParseMode      string `json:"parse_mode,omitempty"`
	Silent         bool   `json:"disable_notification,omitempty"`
	ProtectContent bool   `json:"protect_content,omitempty"`
	ReplyToMsg     int    `json:"reply_to_message_id,omitempty"`
}

// DocumentMessage represents a document message for Telegram
type DocumentMessage struct {
	BaseMessage
	ChatID         string `json:"chat_id"`
	Document       string `json:"document"` // File path, URL, or file_id
	Caption        string `json:"caption,omitempty"`
	ParseMode      string `json:"parse_mode,omitempty"`
	Silent         bool   `json:"disable_notification,omitempty"`
	ProtectContent bool   `json:"protect_content,omitempty"`
	ReplyToMsg     int    `json:"reply_to_message_id,omitempty"`
}

// LocationMessage represents a location message for Telegram
type LocationMessage struct {
	BaseMessage
	ChatID         string  `json:"chat_id"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Silent         bool    `json:"disable_notification,omitempty"`
	ProtectContent bool    `json:"protect_content,omitempty"`
	ReplyToMsg     int     `json:"reply_to_message_id,omitempty"`
}

// ContactMessage represents a contact message for Telegram
type ContactMessage struct {
	BaseMessage
	ChatID         string `json:"chat_id"`
	PhoneNumber    string `json:"phone_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name,omitempty"`
	VCard          string `json:"vcard,omitempty"`
	Silent         bool   `json:"disable_notification,omitempty"`
	ProtectContent bool   `json:"protect_content,omitempty"`
	ReplyToMsg     int    `json:"reply_to_message_id,omitempty"`
}

// PollMessage represents a poll message for Telegram
type PollMessage struct {
	BaseMessage
	ChatID                string   `json:"chat_id"`
	Question              string   `json:"question"`
	Options               []string `json:"options"`
	IsAnonymous           bool     `json:"is_anonymous,omitempty"`
	Type                  string   `json:"type,omitempty"`
	AllowsMultipleAnswers bool     `json:"allows_multiple_answers,omitempty"`
	Silent                bool     `json:"disable_notification,omitempty"`
	ProtectContent        bool     `json:"protect_content,omitempty"`
	ReplyToMsg            int      `json:"reply_to_message_id,omitempty"`
}

// ProviderType implementations
func (m *TextMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *PhotoMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *DocumentMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *LocationMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *ContactMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *PollMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

// Validate implementations
func (m *TextMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if strings.TrimSpace(m.Text) == "" {
		return core.NewParamError("text cannot be empty")
	}
	return nil
}

func (m *PhotoMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if strings.TrimSpace(m.Photo) == "" {
		return core.NewParamError("photo cannot be empty")
	}
	return nil
}

func (m *DocumentMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if strings.TrimSpace(m.Document) == "" {
		return core.NewParamError("document cannot be empty")
	}
	return nil
}

func (m *LocationMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	return nil
}

func (m *ContactMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if strings.TrimSpace(m.PhoneNumber) == "" {
		return core.NewParamError("phone_number cannot be empty")
	}
	if strings.TrimSpace(m.FirstName) == "" {
		return core.NewParamError("first_name cannot be empty")
	}
	return nil
}

func (m *PollMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if strings.TrimSpace(m.Question) == "" {
		return core.NewParamError("question cannot be empty")
	}
	if len(m.Options) < 2 || len(m.Options) > 10 {
		return core.NewParamError("poll must have between 2 and 10 options")
	}
	return nil
}

// Message options
type TextMessageOption func(*TextMessage)
type PhotoMessageOption func(*PhotoMessage)
type DocumentMessageOption func(*DocumentMessage)
type LocationMessageOption func(*LocationMessage)
type ContactMessageOption func(*ContactMessage)
type PollMessageOption func(*PollMessage)

// Text message options
func WithParseMode(mode string) TextMessageOption {
	return func(m *TextMessage) { m.ParseMode = mode }
}

func WithDisableWebPreview(disable bool) TextMessageOption {
	return func(m *TextMessage) { m.DisableWebPreview = disable }
}

func WithSilent(silent bool) TextMessageOption {
	return func(m *TextMessage) { m.Silent = silent }
}

func WithProtectContent(protect bool) TextMessageOption {
	return func(m *TextMessage) { m.ProtectContent = protect }
}

func WithReplyTo(replyTo int) TextMessageOption {
	return func(m *TextMessage) { m.ReplyToMsg = replyTo }
}

// Photo message options
func WithPhotoCaption(caption string) PhotoMessageOption {
	return func(m *PhotoMessage) { m.Caption = caption }
}

func WithPhotoParseMode(mode string) PhotoMessageOption {
	return func(m *PhotoMessage) { m.ParseMode = mode }
}

func WithPhotoSilent(silent bool) PhotoMessageOption {
	return func(m *PhotoMessage) { m.Silent = silent }
}

func WithPhotoProtectContent(protect bool) PhotoMessageOption {
	return func(m *PhotoMessage) { m.ProtectContent = protect }
}

func WithPhotoReplyTo(replyTo int) PhotoMessageOption {
	return func(m *PhotoMessage) { m.ReplyToMsg = replyTo }
}

// Document message options
func WithDocumentCaption(caption string) DocumentMessageOption {
	return func(m *DocumentMessage) { m.Caption = caption }
}

func WithDocumentParseMode(mode string) DocumentMessageOption {
	return func(m *DocumentMessage) { m.ParseMode = mode }
}

func WithDocumentSilent(silent bool) DocumentMessageOption {
	return func(m *DocumentMessage) { m.Silent = silent }
}

func WithDocumentProtectContent(protect bool) DocumentMessageOption {
	return func(m *DocumentMessage) { m.ProtectContent = protect }
}

func WithDocumentReplyTo(replyTo int) DocumentMessageOption {
	return func(m *DocumentMessage) { m.ReplyToMsg = replyTo }
}

// Location message options
func WithLocationSilent(silent bool) LocationMessageOption {
	return func(m *LocationMessage) { m.Silent = silent }
}

func WithLocationProtectContent(protect bool) LocationMessageOption {
	return func(m *LocationMessage) { m.ProtectContent = protect }
}

func WithLocationReplyTo(replyTo int) LocationMessageOption {
	return func(m *LocationMessage) { m.ReplyToMsg = replyTo }
}

// Contact message options
func WithContactLastName(lastName string) ContactMessageOption {
	return func(m *ContactMessage) { m.LastName = lastName }
}

func WithContactVCard(vcard string) ContactMessageOption {
	return func(m *ContactMessage) { m.VCard = vcard }
}

func WithContactSilent(silent bool) ContactMessageOption {
	return func(m *ContactMessage) { m.Silent = silent }
}

func WithContactProtectContent(protect bool) ContactMessageOption {
	return func(m *ContactMessage) { m.ProtectContent = protect }
}

func WithContactReplyTo(replyTo int) ContactMessageOption {
	return func(m *ContactMessage) { m.ReplyToMsg = replyTo }
}

// Poll message options
func WithPollIsAnonymous(anonymous bool) PollMessageOption {
	return func(m *PollMessage) { m.IsAnonymous = anonymous }
}

func WithPollType(pollType string) PollMessageOption {
	return func(m *PollMessage) { m.Type = pollType }
}

func WithPollAllowsMultipleAnswers(allow bool) PollMessageOption {
	return func(m *PollMessage) { m.AllowsMultipleAnswers = allow }
}

func WithPollSilent(silent bool) PollMessageOption {
	return func(m *PollMessage) { m.Silent = silent }
}

func WithPollProtectContent(protect bool) PollMessageOption {
	return func(m *PollMessage) { m.ProtectContent = protect }
}

func WithPollReplyTo(replyTo int) PollMessageOption {
	return func(m *PollMessage) { m.ReplyToMsg = replyTo }
}

// Constructor functions
func NewTextMessage(chatID, text string, opts ...TextMessageOption) *TextMessage {
	msg := &TextMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeText,
		},
		ChatID: chatID,
		Text:   text,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

func NewPhotoMessage(chatID, photo string, opts ...PhotoMessageOption) *PhotoMessage {
	msg := &PhotoMessage{
		BaseMessage: BaseMessage{
			MsgType: TypePhoto,
		},
		ChatID: chatID,
		Photo:  photo,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

func NewDocumentMessage(chatID, document string, opts ...DocumentMessageOption) *DocumentMessage {
	msg := &DocumentMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeDocument,
		},
		ChatID:   chatID,
		Document: document,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

func NewLocationMessage(chatID string, latitude, longitude float64, opts ...LocationMessageOption) *LocationMessage {
	msg := &LocationMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeLocation,
		},
		ChatID:    chatID,
		Latitude:  latitude,
		Longitude: longitude,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

func NewContactMessage(chatID, phoneNumber, firstName string, opts ...ContactMessageOption) *ContactMessage {
	msg := &ContactMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeContact,
		},
		ChatID:      chatID,
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

func NewPollMessage(chatID, question string, options []string, opts ...PollMessageOption) *PollMessage {
	msg := &PollMessage{
		BaseMessage: BaseMessage{
			MsgType: TypePoll,
		},
		ChatID:   chatID,
		Question: question,
		Options:  options,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}
