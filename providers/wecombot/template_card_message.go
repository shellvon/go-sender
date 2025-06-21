package wecombot

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// TemplateCardType defines the allowed types for template cards.
type TemplateCardType string

const (
	// CardTypeTextNotice represents a text notice template card.
	CardTypeTextNotice TemplateCardType = "text_notice"
	// CardTypeNewsNotice represents a news notice template card.
	CardTypeNewsNotice TemplateCardType = "news_notice"
)

// Source represents the source section of a template card.
type Source struct {
	// URL of the source icon image.
	IconURL string `json:"icon_url"`
	// Description of the source, recommended to be no more than 13 characters.
	Desc string `json:"desc"`
	// Color of the source text. Currently supported values:
	// 0 (default) Gray, 1 Black, 2 Red, 3 Green.
	DescColor int `json:"desc_color"`
}

// MainTitle represents the main title section of a template card.
type MainTitle struct {
	// Primary title, recommended to be no more than 26 characters.
	// Either main_title.title or sub_title_text must be provided for template cards.
	Title string `json:"title"`
	// Auxiliary title information, recommended to be no more than 30 characters.
	Desc string `json:"desc"`
}

// EmphasisContent represents the emphasized content section (applicable to text_notice only).
type EmphasisContent struct {
	// Data content for the key data style, recommended to be no more than 10 characters.
	Title string `json:"title"`
	// Description content for the key data style, recommended to be no more than 15 characters.
	Desc string `json:"desc"`
}

// QuoteArea represents the quote area style, generally not recommended for use with EmphasisContent.
type QuoteArea struct {
	// Click event type for the quote area: 0 or empty means no click event,
	// 1 means jump to URL, 2 means jump to mini-program.
	Type int `json:"type"`
	// URL to jump to when clicked. Required if quote_area.type is 1.
	URL string `json:"url,omitempty"`
	// AppID of the mini-program to jump to when clicked. Required if quote_area.type is 2.
	AppID string `json:"appid,omitempty"`
	// Page path of the mini-program to jump to when clicked. Optional if quote_area.type is 2.
	PagePath string `json:"pagepath,omitempty"`
	// Title of the quote area.
	Title string `json:"title"`
	// Text content of the quote.
	QuoteText string `json:"quote_text"`
}

// CardImage represents the card image section (applicable to news_notice only).
type CardImage struct {
	// URL of the image.
	URL string `json:"url"`
	// Aspect ratio of the image.
	AspectRatio float64 `json:"aspect_ratio"`
}

// ImageTextArea represents the image text area section (applicable to news_notice only).
type ImageTextArea struct {
	// Type of the image text area.
	Type int `json:"type"`
	// URL to jump to when clicked.
	URL string `json:"url"`
	// Title of the image text area.
	Title string `json:"title"`
	// Description of the image text area.
	Desc string `json:"desc"`
	// URL of the image.
	ImageURL string `json:"image_url"`
}

// HorizontalContent represents an item in the horizontal content list.
// This field can be an empty array, but if data is present, corresponding fields
// must be validated. The list length must not exceed 6 items.
type HorizontalContent struct {
	// Secondary title, recommended to be no more than 5 characters.
	Keyname string `json:"keyname"`
	// Secondary text. If horizontal_content_list.type is 2, this field represents
	// the file name (including file type), recommended to be no more than 26 characters.
	Value string `json:"value"`
	// Type of the secondary title information content:
	// 1 for URL, 2 for file attachment, 3 for jump to member details.
	Type int `json:"type,omitempty"`
	// URL to jump to. Required if horizontal_content_list.type is 1.
	URL string `json:"url,omitempty"`
	// Media ID of the attachment. Required if horizontal_content_list.type is 2.
	MediaID string `json:"media_id,omitempty"`
}

// VerticalContent represents an item in the vertical content list (applicable to news_notice only).
type VerticalContent struct {
	// Title of the vertical content item.
	Title string `json:"title"`
	// Description of the vertical content item.
	Desc string `json:"desc"`
}

// Jump represents an item in the jump list.
type Jump struct {
	// Type of the jump.
	Type int `json:"type"`
	// URL to jump to. Required if Type is 1.
	URL string `json:"url,omitempty"`
	// AppID of the mini-program to jump to. Required if Type is 2.
	AppID string `json:"appid,omitempty"`
	// Page path of the mini-program. Optional if Type is 2.
	PagePath string `json:"pagepath,omitempty"`
	// Title of the jump.
	Title string `json:"title"`
}

// CardAction represents the card action section.
type CardAction struct {
	// Type of the card action.
	Type int `json:"type"`
	// URL to jump to. Required if Type is 1.
	URL string `json:"url,omitempty"`
	// AppID of the mini-program to jump to. Required if Type is 2.
	AppID string `json:"appid,omitempty"`
	// Page path of the mini-program. Optional if Type is 2.
	PagePath string `json:"pagepath,omitempty"`
}

// TemplateCard represents the template_card section of the message.
type TemplateCard struct {
	CardType              TemplateCardType    `json:"card_type"`
	Source                *Source             `json:"source"`
	MainTitle             MainTitle           `json:"main_title"`
	EmphasisContent       *EmphasisContent    `json:"emphasis_content,omitempty"` // text_notice only
	QuoteArea             *QuoteArea          `json:"quote_area,omitempty"`
	SubTitleText          string              `json:"sub_title_text,omitempty"`        // text_notice only
	CardImage             *CardImage          `json:"card_image,omitempty"`            // news_notice only
	ImageTextArea         *ImageTextArea      `json:"image_text_area,omitempty"`       // news_notice only
	VerticalContentList   []VerticalContent   `json:"vertical_content_list,omitempty"` // news_notice only
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list,omitempty"`
	JumpList              []Jump              `json:"jump_list,omitempty"`
	CardAction            CardAction          `json:"card_action"`
}

// TemplateCardMessage represents a WeCom template card message.
// For more details, refer to the WeCom API documentation:
// https://developer.work.weixin.qq.com/document/path/91770#%E6%A8%A1%E7%89%88%E5%8D%A1%E7%89%87%E7%B1%BB%E5%9E%8B
type TemplateCardMessage struct {
	BaseMessage
	TemplateCard TemplateCard `json:"template_card"`
}

// Validate validates the TemplateCardMessage according to WeCom API rules.
func (m *TemplateCardMessage) Validate() error {
	// Check card type
	if m.TemplateCard.CardType != CardTypeTextNotice && m.TemplateCard.CardType != CardTypeNewsNotice {
		return core.NewParamError(fmt.Sprintf("invalid card_type: %s; must be %s or %s", m.TemplateCard.CardType, CardTypeTextNotice, CardTypeNewsNotice))
	}

	// Validate MainTitle and SubTitleText coexistence/mutual exclusivity
	if m.TemplateCard.MainTitle.Title == "" && m.TemplateCard.SubTitleText == "" {
		return core.NewParamError("main_title.title or sub_title_text must be provided")
	}

	// Validate card_action
	if m.TemplateCard.CardAction.Type == 0 {
		return core.NewParamError("card_action.type is required")
	}
	if m.TemplateCard.CardAction.Type == 1 && m.TemplateCard.CardAction.URL == "" {
		return core.NewParamError("card_action.url is required for type 1")
	}
	if m.TemplateCard.CardAction.Type == 2 && (m.TemplateCard.CardAction.AppID == "" || m.TemplateCard.CardAction.PagePath == "") {
		return core.NewParamError("card_action.appid and pagepath are required for type 2")
	}
	// Check for unsupported card_action types
	if m.TemplateCard.CardAction.Type < 0 || m.TemplateCard.CardAction.Type > 2 {
		return core.NewParamError(fmt.Sprintf("invalid card_action.type: %d; must be 0, 1, or 2", m.TemplateCard.CardAction.Type))
	}

	// Validate list lengths
	// The comment for HorizontalContent says "list length cannot exceed 6", but the code validates against 5.
	if len(m.TemplateCard.HorizontalContentList) > 6 { // Changed to 6 based on updated comment.
		return core.NewParamError("horizontal_content_list cannot exceed 6 items")
	}
	if len(m.TemplateCard.JumpList) > 5 {
		return core.NewParamError("jump_list cannot exceed 5 items")
	}
	if m.TemplateCard.CardType == CardTypeNewsNotice && len(m.TemplateCard.VerticalContentList) > 5 {
		return core.NewParamError("vertical_content_list cannot exceed 5 items")
	}

	// Validate horizontal_content_list items
	for i, item := range m.TemplateCard.HorizontalContentList {
		if item.Keyname == "" { // Keyname is mandatory
			return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: keyname cannot be empty", i))
		}
		if item.Value == "" { // Value is mandatory
			return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: value cannot be empty", i))
		}

		// Enforce length limits as per comments (using byte length as common for API limits)
		if len([]byte(item.Keyname)) > 5 {
			return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: keyname exceeds 5 bytes", i))
		}
		if len([]byte(item.Value)) > 26 {
			return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: value exceeds 26 bytes", i))
		}

		if item.Type == 1 {
			if item.URL == "" {
				return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: url is required for type 1", i))
			}
		} else if item.Type == 2 {
			if item.MediaID == "" {
				return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: media_id is required for type 2", i))
			}
		} else if item.Type < 0 || item.Type > 3 { // Check for invalid types for HorizontalContent
			return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: invalid type %d; must be 0, 1, 2, or 3", i, item.Type))
		}
	}

	// Validate jump_list items
	for i, jump := range m.TemplateCard.JumpList {
		if jump.Type == 0 || jump.Title == "" { // Type 0 is usually 'no action' but API requires a type to be specified. Title is always required.
			return core.NewParamError(fmt.Sprintf("jump_list item %d: type and title are required", i))
		}
		if jump.Type == 1 && jump.URL == "" {
			return core.NewParamError(fmt.Sprintf("jump_list item %d: url is required for type 1", i))
		}
		if jump.Type == 2 && (jump.AppID == "" || jump.PagePath == "") {
			return core.NewParamError(fmt.Sprintf("jump_list item %d: appid and pagepath are required for type 2", i))
		}
		// Check for unsupported jump.Type
		if jump.Type < 0 || jump.Type > 2 {
			return core.NewParamError(fmt.Sprintf("jump_list item %d: invalid type: %d; must be 0, 1, or 2", i, jump.Type))
		}
		if len([]byte(jump.Title)) > 128 { // Assuming a reasonable limit if not specified, often aligns with main title.
			return core.NewParamError(fmt.Sprintf("jump_list item %d: title exceeds 128 bytes", i))
		}
	}

	// Validate news_notice specific fields
	if m.TemplateCard.CardType == CardTypeNewsNotice {
		if m.TemplateCard.CardImage == nil || m.TemplateCard.CardImage.URL == "" {
			return core.NewParamError("card_image.url is required for news_notice cards")
		}
		if m.TemplateCard.CardImage.AspectRatio == 0 { // Aspect ratio should be a positive value
			return core.NewParamError("card_image.aspect_ratio cannot be zero for news_notice cards")
		}
		// Validate VerticalContentList for news_notice
		for i, vc := range m.TemplateCard.VerticalContentList {
			if vc.Title == "" {
				return core.NewParamError(fmt.Sprintf("vertical_content_list item %d: title cannot be empty", i))
			}
			// Add length validation for title and desc if API specifies. (Currently missing in comments)
		}

	} else if m.TemplateCard.CardType == CardTypeTextNotice {
		// Ensure news_notice-specific fields are NOT set for text_notice
		if m.TemplateCard.CardImage != nil {
			return core.NewParamError("card_image is not allowed for text_notice cards")
		}
		if m.TemplateCard.ImageTextArea != nil {
			return core.NewParamError("image_text_area is not allowed for text_notice cards")
		}
		if len(m.TemplateCard.VerticalContentList) > 0 {
			return core.NewParamError("vertical_content_list is not allowed for text_notice cards")
		}
	}

	// Validate quote_area if present
	if m.TemplateCard.QuoteArea != nil {
		if m.TemplateCard.QuoteArea.Type < 0 || m.TemplateCard.QuoteArea.Type > 2 { // Check for invalid types
			return core.NewParamError(fmt.Sprintf("quote_area: invalid type: %d; must be 0, 1, or 2", m.TemplateCard.QuoteArea.Type))
		}
		if m.TemplateCard.QuoteArea.Title == "" || m.TemplateCard.QuoteArea.QuoteText == "" {
			return core.NewParamError("quote_area must have title and quote_text")
		}
		if m.TemplateCard.QuoteArea.Type == 1 && m.TemplateCard.QuoteArea.URL == "" {
			return core.NewParamError("quote_area.url is required for type 1")
		}
		if m.TemplateCard.QuoteArea.Type == 2 && (m.TemplateCard.QuoteArea.AppID == "" || m.TemplateCard.QuoteArea.PagePath == "") {
			return core.NewParamError("quote_area.appid and pagepath are required for type 2")
		}
	}

	// Validate image_text_area if present (news_notice only)
	if m.TemplateCard.ImageTextArea != nil {
		if m.TemplateCard.CardType != CardTypeNewsNotice {
			return core.NewParamError("image_text_area is only allowed for news_notice cards")
		}
		if m.TemplateCard.ImageTextArea.Type == 0 || m.TemplateCard.ImageTextArea.Title == "" || m.TemplateCard.ImageTextArea.ImageURL == "" {
			return core.NewParamError("image_text_area must have type, title, and image_url")
		}
		if m.TemplateCard.ImageTextArea.Type == 1 && m.TemplateCard.ImageTextArea.URL == "" {
			return core.NewParamError("image_text_area.url is required for type 1")
		}
		if m.TemplateCard.ImageTextArea.Type < 0 || m.TemplateCard.ImageTextArea.Type > 1 { // Assuming only 0 and 1 as valid types
			return core.NewParamError(fmt.Sprintf("image_text_area: invalid type: %d; must be 0 or 1", m.TemplateCard.ImageTextArea.Type))
		}
	}

	// Validate MainTitle length and its auxiliary description
	if len([]byte(m.TemplateCard.MainTitle.Title)) > 26 {
		return core.NewParamError("main_title.title exceeds 26 bytes")
	}
	if len([]byte(m.TemplateCard.MainTitle.Desc)) > 30 {
		return core.NewParamError("main_title.desc exceeds 30 bytes")
	}

	// Validate SubTitleText length if present
	if m.TemplateCard.SubTitleText != "" && len([]byte(m.TemplateCard.SubTitleText)) > 4096 { // A common large text limit, specific to WeCom
		return core.NewParamError("sub_title_text exceeds 4096 bytes") // Assuming byte limit based on other fields
	}

	// Validate Source field if present
	if m.TemplateCard.Source != nil {
		if m.TemplateCard.Source.IconURL == "" { // IconURL typically required if Source is used
			return core.NewParamError("source.icon_url is required if source is provided")
		}
		if len([]byte(m.TemplateCard.Source.Desc)) > 13 {
			return core.NewParamError("source.desc exceeds 13 bytes")
		}
		if m.TemplateCard.Source.DescColor < 0 || m.TemplateCard.Source.DescColor > 3 {
			return core.NewParamError(fmt.Sprintf("invalid source.desc_color: %d; must be 0, 1, 2, or 3", m.TemplateCard.Source.DescColor))
		}
	}

	// Validate EmphasisContent if present (text_notice only)
	if m.TemplateCard.EmphasisContent != nil {
		if m.TemplateCard.CardType != CardTypeTextNotice {
			return core.NewParamError("emphasis_content is only allowed for text_notice cards")
		}
		if len([]byte(m.TemplateCard.EmphasisContent.Title)) > 10 {
			return core.NewParamError("emphasis_content.title exceeds 10 bytes")
		}
		if len([]byte(m.TemplateCard.EmphasisContent.Desc)) > 15 {
			return core.NewParamError("emphasis_content.desc exceeds 15 bytes")
		}
	}

	return nil
}

// TemplateCardMessageOption defines a function type for configuring TemplateCardMessage.
type TemplateCardMessageOption func(*TemplateCardMessage)

// WithSource sets the Source for the TemplateCardMessage.
func WithSource(source *Source) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.Source = source
	}
}

// WithEmphasisContent sets the EmphasisContent for the TemplateCardMessage (text_notice only).
func WithEmphasisContent(emphasis *EmphasisContent) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.EmphasisContent = emphasis
	}
}

// WithQuoteArea sets the QuoteArea for the TemplateCardMessage.
func WithQuoteArea(quote *QuoteArea) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.QuoteArea = quote
	}
}

// WithSubTitleText sets the SubTitleText for the TemplateCardMessage (text_notice only).
func WithSubTitleText(subTitle string) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.SubTitleText = subTitle
	}
}

// WithCardImage sets the CardImage for the TemplateCardMessage (news_notice only).
func WithCardImage(image *CardImage) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.CardImage = image
	}
}

// WithImageTextArea sets the ImageTextArea for the TemplateCardMessage (news_notice only).
func WithImageTextArea(imageText *ImageTextArea) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.ImageTextArea = imageText
	}
}

// WithVerticalContentList sets the VerticalContentList for the TemplateCardMessage (news_notice only).
func WithVerticalContentList(list []VerticalContent) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.VerticalContentList = list
	}
}

// WithHorizontalContentList sets the HorizontalContentList for the TemplateCardMessage.
func WithHorizontalContentList(list []HorizontalContent) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.HorizontalContentList = list
	}
}

// WithJumpList sets the JumpList for the TemplateCardMessage.
func WithJumpList(list []Jump) TemplateCardMessageOption {
	return func(m *TemplateCardMessage) {
		m.TemplateCard.JumpList = list
	}
}

// NewTemplateCardMessage creates a new TemplateCardMessage with required fields and applies optional configurations.
func NewTemplateCardMessage(cardType TemplateCardType, mainTitle MainTitle, cardAction CardAction, opts ...TemplateCardMessageOption) *TemplateCardMessage {
	msg := &TemplateCardMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeTemplateCard,
		},
		TemplateCard: TemplateCard{
			CardType:   cardType,
			MainTitle:  mainTitle,
			CardAction: cardAction,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}
