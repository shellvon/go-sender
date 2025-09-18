package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// InteractiveMessage represents a Lark/Feishu card message (schema 2.0).
// See:
//   - https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot#5a997364
//   - https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot#478cb64f
//   - https://open.feishu.cn/document/feishu-cards/card-json-v2-structure
type InteractiveMessage struct {
	BaseMessage

	Card Card `json:"card"`
}

// Compile-time assertion: InteractiveMessage implements Message interface.
var (
	_ core.Validatable = (*InteractiveMessage)(nil)
)

type Card struct {
	Schema string `json:"schema"` // should be "2.0"
	// 配置卡片的全局行为，包括流式更新模式（JSON 2.0 新增能力）、是否允许被转发、是否为共享卡片等。

	Config *CardConfig `json:"config,omitempty"`
	// 指定卡片整体的点击跳转链接。你可以配置一个默认链接，也可以分别为 PC 端、Android 端、iOS 端配置不同的跳转链接。
	CardLink *CardLink `json:"card_link,omitempty"`
	// 标题组件相关配置。详情参考标题组件。
	//  - https://open.feishu.cn/document/feishu-cards/card-json-v2-components/content-components/title
	Header *CardHeader `json:"header,omitempty"`
	// 卡片正文，包含一个名为 elements 的数组，用于放置各类组件。
	Body *CardBody `json:"body,omitempty"`
}

// PrintFrequency for print_frequency_ms and print_step fields.
type PrintFrequency struct {
	Default int `json:"default,omitempty"`
	Android int `json:"android,omitempty"`
	IOS     int `json:"ios,omitempty"`
	PC      int `json:"pc,omitempty"`
}

// CardStreamingConfig for streaming_config field.
type CardStreamingConfig struct {
	PrintFrequencyMs *PrintFrequency `json:"print_frequency_ms,omitempty"`
	PrintStep        *PrintFrequency `json:"print_step,omitempty"`
	PrintStrategy    string          `json:"print_strategy,omitempty"` // fast/delay
}

// CardConfig for card config field.
// https://open.feishu.cn/document/feishu-cards/card-json-v2-structure
type CardConfig struct {
	StreamingMode            bool                 `json:"streaming_mode,omitempty"`
	StreamingConfig          *CardStreamingConfig `json:"streaming_config,omitempty"`
	Summary                  *CardSummary         `json:"summary,omitempty"`
	Locales                  []string             `json:"locales,omitempty"`
	EnableForward            *bool                `json:"enable_forward,omitempty"`
	UpdateMulti              *bool                `json:"update_multi,omitempty"`
	WidthMode                string               `json:"width_mode,omitempty"`
	UseCustomTranslation     *bool                `json:"use_custom_translation,omitempty"`
	EnableForwardInteraction *bool                `json:"enable_forward_interaction,omitempty"`
	Style                    *CardStyle           `json:"style,omitempty"`
}

type CardSummary struct {
	Content     string            `json:"content,omitempty"`
	I18nContent map[string]string `json:"i18n_content,omitempty"`
}

type CardStyle struct {
	TextSize map[string]CardTextSize `json:"text_size,omitempty"`
	Color    map[string]CardColor    `json:"color,omitempty"`
}

type CardTextSize struct {
	Default string `json:"default,omitempty"`
	PC      string `json:"pc,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
}

type CardColor struct {
	LightMode string `json:"light_mode,omitempty"`
	DarkMode  string `json:"dark_mode,omitempty"`
}

type CardLink struct {
	URL        string `json:"url,omitempty"`
	AndroidURL string `json:"android_url,omitempty"`
	IOSURL     string `json:"ios_url,omitempty"`
	PCURL      string `json:"pc_url,omitempty"`
}

// CardHeader represents the header component of a card.
// https://open.feishu.cn/document/feishu-cards/card-json-v2-components/content-components/title
type CardHeader struct {
	Title       *CardText     `json:"title"`
	Subtitle    *CardText     `json:"subtitle,omitempty"`
	TextTagList []CardTextTag `json:"text_tag_list,omitempty"`
	// 国际化文本标签列表。
	// 支持的语言列表：
	//  - zh_cn：简体中文
	//  - en_us：英文
	//  - ja_jp：日文
	//  - zh_hk：繁体中文（中国香港）
	//  - zh_tw：繁体中文（中国台湾）
	//  - id_id: 印尼语
	//  - vi_vn: 越南语
	//  - th_th: 泰语
	//  - pt_br: 葡萄牙语
	//  - es_es: 西班牙语
	//  - ko_kr: 韩语
	//  - de_de: 德语
	//  - fr_fr: 法语
	//  - it_it: 意大利语
	//  - ru_ru: 俄语
	//  - ms_my: 马来语
	I18nTextTagList map[string][]CardTextTag `json:"i18n_text_tag_list,omitempty"`
	Template        string                   `json:"template,omitempty"` // 标题主题样式颜色。支持 "blue"|"wathet"|"turquoise"|"green"|"yellow"|"orange"|"red"|"carmine"|"violet"|"purple"|"indigo"|"grey"|"default"。默认值 default
	Icon            *CardHeaderIcon          `json:"icon,omitempty"`
	Padding         string                   `json:"padding,omitempty"`
}

type CardText struct {
	Tag     string `json:"tag"` // "plain_text" or "lark_md"
	Content string `json:"content"`
}

type CardTextTag struct {
	Tag       string   `json:"tag"`
	ElementID string   `json:"element_id,omitempty"`
	Text      CardText `json:"text"`
	Color     string   `json:"color,omitempty"`
}

type CardHeaderIcon struct {
	Tag    string `json:"tag"`
	Token  string `json:"token,omitempty"`
	Color  string `json:"color,omitempty"`
	ImgKey string `json:"img_key,omitempty"`
}

type CardBody struct {
	Direction         string        `json:"direction,omitempty"`
	Padding           string        `json:"padding,omitempty"`
	HorizontalSpacing string        `json:"horizontal_spacing,omitempty"`
	HorizontalAlign   string        `json:"horizontal_align,omitempty"`
	VerticalSpacing   string        `json:"vertical_spacing,omitempty"`
	VerticalAlign     string        `json:"vertical_align,omitempty"`
	Elements          []CardElement `json:"elements"`
}

type CardElement map[string]interface{} // Card 2.0 elements are arbitrary JSON objects with a tag

// InteractiveBuilder for schema 2.0
//
// Example:
//
//	msg := lark.Interactive().
//	         Config(...).
//	         CardLink(...).
//	         Header(...).
//	         Body(...).
//	         Build()
type InteractiveBuilder struct {
	config   *CardConfig
	cardLink *CardLink
	header   *CardHeader
	body     *CardBody
}

func Interactive() *InteractiveBuilder { return &InteractiveBuilder{} }

func (b *InteractiveBuilder) Config(cfg *CardConfig) *InteractiveBuilder {
	b.config = cfg
	return b
}
func (b *InteractiveBuilder) CardLink(link *CardLink) *InteractiveBuilder {
	b.cardLink = link
	return b
}

// Header sets the whole header struct.
func (b *InteractiveBuilder) Header(header *CardHeader) *InteractiveBuilder {
	b.header = header
	return b
}

// HeaderTitle sets the header title (creates header if nil).
func (b *InteractiveBuilder) HeaderTitle(tag, content string) *InteractiveBuilder {
	if b.header == nil {
		b.header = &CardHeader{}
	}
	b.header.Title = &CardText{Tag: tag, Content: content}
	return b
}

// HeaderSubtitle sets the header subtitle (creates header if nil).
func (b *InteractiveBuilder) HeaderSubtitle(tag, content string) *InteractiveBuilder {
	if b.header == nil {
		b.header = &CardHeader{}
	}
	b.header.Subtitle = &CardText{Tag: tag, Content: content}
	return b
}

// HeaderTextTagList sets the header text tag list.
func (b *InteractiveBuilder) HeaderTextTagList(tags []CardTextTag) *InteractiveBuilder {
	if b.header == nil {
		b.header = &CardHeader{}
	}
	b.header.TextTagList = tags
	return b
}

// HeaderI18nTextTagList sets the header i18n text tag list.
func (b *InteractiveBuilder) HeaderI18nTextTagList(i18n map[string][]CardTextTag) *InteractiveBuilder {
	if b.header == nil {
		b.header = &CardHeader{}
	}
	b.header.I18nTextTagList = i18n
	return b
}

// HeaderTemplate sets the header template color.
func (b *InteractiveBuilder) HeaderTemplate(template string) *InteractiveBuilder {
	if b.header == nil {
		b.header = &CardHeader{}
	}
	b.header.Template = template
	return b
}

// HeaderIcon sets the header icon.
func (b *InteractiveBuilder) HeaderIcon(icon *CardHeaderIcon) *InteractiveBuilder {
	if b.header == nil {
		b.header = &CardHeader{}
	}
	b.header.Icon = icon
	return b
}

// HeaderPadding sets the header padding.
func (b *InteractiveBuilder) HeaderPadding(padding string) *InteractiveBuilder {
	if b.header == nil {
		b.header = &CardHeader{}
	}
	b.header.Padding = padding
	return b
}

// Body sets the whole body struct.
func (b *InteractiveBuilder) Body(body *CardBody) *InteractiveBuilder {
	b.body = body
	return b
}

// BodyElements sets the body elements array (creates body if nil).
func (b *InteractiveBuilder) BodyElements(elements []CardElement) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.Elements = elements
	return b
}

// AddElement appends a single element to the body (creates body if nil).
func (b *InteractiveBuilder) AddElement(element CardElement) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.Elements = append(b.body.Elements, element)
	return b
}

// BodyDirection sets the body direction.
func (b *InteractiveBuilder) BodyDirection(direction string) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.Direction = direction
	return b
}

// BodyPadding sets the body padding.
func (b *InteractiveBuilder) BodyPadding(padding string) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.Padding = padding
	return b
}

// BodyHorizontalSpacing sets the body horizontal spacing.
func (b *InteractiveBuilder) BodyHorizontalSpacing(spacing string) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.HorizontalSpacing = spacing
	return b
}

// BodyHorizontalAlign sets the body horizontal align.
func (b *InteractiveBuilder) BodyHorizontalAlign(align string) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.HorizontalAlign = align
	return b
}

// BodyVerticalSpacing sets the body vertical spacing.
func (b *InteractiveBuilder) BodyVerticalSpacing(spacing string) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.VerticalSpacing = spacing
	return b
}

// BodyVerticalAlign sets the body vertical align.
func (b *InteractiveBuilder) BodyVerticalAlign(align string) *InteractiveBuilder {
	if b.body == nil {
		b.body = &CardBody{}
	}
	b.body.VerticalAlign = align
	return b
}
func (b *InteractiveBuilder) Build() *InteractiveMessage {
	return &InteractiveMessage{
		BaseMessage: newBaseMessage(TypeInteractive),
		Card: Card{
			Schema:   "2.0",
			Config:   b.config,
			CardLink: b.cardLink,
			Header:   b.header,
			Body:     b.body,
		},
	}
}

func NewInteractiveMessage() *InteractiveMessage {
	return Interactive().Build()
}

func (m *InteractiveMessage) GetMsgType() MessageType {
	return TypeInteractive
}

func (m *InteractiveMessage) Validate() error {
	if m.Card.Schema != "2.0" {
		return errors.New("card schema must be 2.0")
	}
	if m.Card.Header == nil || m.Card.Header.Title == nil || m.Card.Header.Title.Content == "" {
		return errors.New("card header.title.content is required")
	}
	if m.Card.Body == nil || len(m.Card.Body.Elements) == 0 {
		return errors.New("card body must have at least one element")
	}
	return nil
}
