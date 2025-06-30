package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// InteractiveMessage represents an interactive (card) message for Lark/Feishu.
type InteractiveMessage struct {
	BaseMessage

	Content InteractiveContent `json:"content"`
}

// InteractiveContent represents the content of an interactive message.
type InteractiveContent struct {
	Card Card `json:"card"`
}

// Card represents the card structure.
type Card struct {
	Header   *CardHeader   `json:"header,omitempty"`
	Elements []CardElement `json:"elements"`
	I18n     *CardI18n     `json:"i18n,omitempty"`
}

// CardHeader represents the card header.
type CardHeader struct {
	Title    *CardText `json:"title,omitempty"`
	Subtitle *CardText `json:"subtitle,omitempty"`
	Icon     *CardIcon `json:"icon,omitempty"`
}

// CardText represents text in a card.
type CardText struct {
	Tag      string `json:"tag"`
	Content  string `json:"content"`
	Lines    int    `json:"lines,omitempty"`
	Size     string `json:"size,omitempty"`
	Type     string `json:"type,omitempty"`
	Color    string `json:"color,omitempty"`
	MaxWidth int    `json:"max_width,omitempty"`
}

// CardIcon represents an icon in a card.
type CardIcon struct {
	Tag       string `json:"tag"`
	ImageKey  string `json:"image_key,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	Emoji     string `json:"emoji,omitempty"`
	Size      string `json:"size,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Align     string `json:"align,omitempty"`
	AltText   string `json:"alt_text,omitempty"`
	Preview   bool   `json:"preview,omitempty"`
	FitMode   string `json:"fit_mode,omitempty"`
	ImageType string `json:"image_type,omitempty"`
}

// CardElement represents an element in a card.
type CardElement struct {
	Tag      string                 `json:"tag"`
	Text     *CardText              `json:"text,omitempty"`
	Image    *CardIcon              `json:"image,omitempty"`
	Action   *CardAction            `json:"action,omitempty"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
	Elements []CardElement          `json:"elements,omitempty"`
}

// CardAction represents an action in a card.
type CardAction struct {
	Tag      string                 `json:"tag"`
	Text     *CardText              `json:"text,omitempty"`
	URL      string                 `json:"url,omitempty"`
	MultiURL *CardMultiURL          `json:"multi_url,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Value    map[string]interface{} `json:"value,omitempty"`
	Confirm  *CardConfirm           `json:"confirm,omitempty"`
}

// CardMultiURL represents multiple URLs for different platforms.
type CardMultiURL struct {
	PC     string `json:"pc,omitempty"`
	Mobile string `json:"mobile,omitempty"`
}

// CardConfirm represents a confirmation dialog.
type CardConfirm struct {
	Title   *CardText `json:"title,omitempty"`
	Text    *CardText `json:"text,omitempty"`
	Confirm *CardText `json:"confirm,omitempty"`
	Deny    *CardText `json:"deny,omitempty"`
}

// CardI18n represents internationalization for a card.
type CardI18n struct {
	ZhCN *Card `json:"zh_cn,omitempty"`
	EnUS *Card `json:"en_us,omitempty"`
}

// NewInteractiveMessage creates a new interactive message.
func NewInteractiveMessage() *InteractiveMessage {
	return &InteractiveMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeInteractive,
		},
		Content: InteractiveContent{
			Card: Card{
				Elements: []CardElement{},
			},
		},
	}
}

// SetHeader sets the card header.
func (m *InteractiveMessage) SetHeader(header *CardHeader) *InteractiveMessage {
	m.Content.Card.Header = header
	return m
}

// AddElement adds an element to the card.
func (m *InteractiveMessage) AddElement(element CardElement) *InteractiveMessage {
	m.Content.Card.Elements = append(m.Content.Card.Elements, element)
	return m
}

// SetI18n sets the internationalization for the card.
func (m *InteractiveMessage) SetI18n(i18n *CardI18n) *InteractiveMessage {
	m.Content.Card.I18n = i18n
	return m
}

// GetMsgType returns the message type.
func (m *InteractiveMessage) GetMsgType() MessageType {
	return TypeInteractive
}

// ProviderType returns the provider type.
func (m *InteractiveMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the interactive message.
func (m *InteractiveMessage) Validate() error {
	if len(m.Content.Card.Elements) == 0 {
		return errors.New("card must have at least one element")
	}
	return nil
}
