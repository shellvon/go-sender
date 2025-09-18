package wecomapp

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// TemplateCardType 定义模板卡片允许的类型.
type TemplateCardType string

const (
	// CardTypeTextNotice 代表文本通知模板卡片.
	CardTypeTextNotice TemplateCardType = "text_notice"
	// CardTypeNewsNotice 代表图文通知模板卡片.
	CardTypeNewsNotice TemplateCardType = "news_notice"
	// CardTypeButtonInteraction 代表按钮交互模板卡片.
	CardTypeButtonInteraction TemplateCardType = "button_interaction"
	// CardTypeVoteInteraction 代表投票交互模板卡片.
	CardTypeVoteInteraction TemplateCardType = "vote_interaction"
)

const (
	HorizontalTypeDefault = 0
	HorizontalTypeURL     = 1
	HorizontalTypeMedia   = 2
	HorizontalTypeMember  = 3
)

// Source 代表模板卡片的来源区域.
type Source struct {
	// IconURL 来源图标图片的URL
	IconURL string `json:"icon_url"`
	// Desc 来源描述，建议不超过13个字符
	Desc string `json:"desc"`
	// DescColor 来源文本的颜色。当前支持的值：
	// 0 (默认) 灰色，1 黑色，2 红色，3 绿色
	DescColor int `json:"desc_color"`
}

// MainTitle 代表模板卡片的主标题区域.
type MainTitle struct {
	// Title 主标题，建议不超过26个字符
	// 模板卡片必须提供main_title.title或sub_title_text之一
	Title string `json:"title"`
	// Desc 辅助标题信息，建议不超过30个字符
	Desc string `json:"desc"`
}

// EmphasisContent 代表强调内容区域（仅适用于text_notice）.
type EmphasisContent struct {
	// Title 关键数据样式的数据内容，建议不超过10个字符
	Title string `json:"title"`
	// Desc 关键数据样式的描述内容，建议不超过15个字符
	Desc string `json:"desc"`
}

// QuoteArea 代表引用区域样式，通常不建议与EmphasisContent同时使用.
type QuoteArea struct {
	// Type 引用区域的点击事件类型：0或空表示无点击事件，
	// 1表示跳转到URL，2表示跳转到小程序
	Type int `json:"type"`
	// URL 点击时跳转的URL。当quote_area.type为1时必填
	URL string `json:"url,omitempty"`
	// AppID 点击时跳转的小程序AppID。当quote_area.type为2时必填
	AppID string `json:"appid,omitempty"`
	// PagePath 点击时跳转的小程序页面路径。当quote_area.type为2时可选
	PagePath string `json:"pagepath,omitempty"`
	// Title 引用区域的标题
	Title string `json:"title"`
	// QuoteText 引用的文本内容
	QuoteText string `json:"quote_text"`
}

// CardImage 代表卡片图片区域（仅适用于news_notice）.
type CardImage struct {
	// URL 图片的URL
	URL string `json:"url"`
	// AspectRatio 图片的宽高比
	AspectRatio float64 `json:"aspect_ratio"`
}

// ImageTextArea 代表图文区域部分（仅适用于news_notice）.
type ImageTextArea struct {
	// Type 图文区域的类型
	Type int `json:"type"`
	// URL 点击时跳转的URL
	URL string `json:"url"`
	// Title 图文区域的标题
	Title string `json:"title"`
	// Desc 图文区域的描述
	Desc string `json:"desc"`
	// ImageURL 图片的URL
	ImageURL string `json:"image_url"`
}

// HorizontalContent 代表水平内容列表中的一项
// 该字段可以是空数组，但如果有数据，必须验证相应字段
// 列表长度不得超过6项.
type HorizontalContent struct {
	// Keyname 次标题，建议不超过5个字符
	Keyname string `json:"keyname"`
	// Value 次文本。如果horizontal_content_list.type为2，该字段表示
	// 文件名（包括文件类型），建议不超过26个字符
	Value string `json:"value"`
	// Type 次标题信息内容的类型：
	// 1表示URL，2表示文件附件，3表示跳转到成员详情
	Type int `json:"type,omitempty"`
	// URL 跳转的URL。当horizontal_content_list.type为1时必填
	URL string `json:"url,omitempty"`
	// MediaID 附件的媒体ID。当horizontal_content_list.type为2时必填
	MediaID string `json:"media_id,omitempty"`
}

// VerticalContent 代表垂直内容列表中的一项（仅适用于news_notice）.
type VerticalContent struct {
	// Title 垂直内容项的标题
	Title string `json:"title"`
	// Desc 垂直内容项的描述
	Desc string `json:"desc"`
}

// Jump 代表跳转列表中的一项.
type Jump struct {
	// Type 跳转的类型
	Type int `json:"type"`
	// URL 跳转的URL。当Type为1时必填
	URL string `json:"url,omitempty"`
	// AppID 跳转的小程序AppID。当Type为2时必填
	AppID string `json:"appid,omitempty"`
	// PagePath 小程序的页面路径。当Type为2时可选
	PagePath string `json:"pagepath,omitempty"`
	// Title 跳转的标题
	Title string `json:"title"`
}

// CardAction 代表卡片动作区域.
type CardAction struct {
	// Type 卡片动作的类型
	Type int `json:"type"`
	// URL 跳转的URL。当Type为1时必填
	URL string `json:"url,omitempty"`
	// AppID 跳转的小程序AppID。当Type为2时必填
	AppID string `json:"appid,omitempty"`
	// PagePath 小程序的页面路径。当Type为2时可选
	PagePath string `json:"pagepath,omitempty"`
}

// TemplateCard 代表消息的template_card部分.
type TemplateCard struct {
	CardType              TemplateCardType    `json:"card_type"`
	Source                *Source             `json:"source"`
	MainTitle             MainTitle           `json:"main_title"`
	EmphasisContent       *EmphasisContent    `json:"emphasis_content,omitempty"` // 仅text_notice
	QuoteArea             *QuoteArea          `json:"quote_area,omitempty"`
	SubTitleText          string              `json:"sub_title_text,omitempty"`        // 仅text_notice
	CardImage             *CardImage          `json:"card_image,omitempty"`            // 仅news_notice
	ImageTextArea         *ImageTextArea      `json:"image_text_area,omitempty"`       // 仅news_notice
	VerticalContentList   []VerticalContent   `json:"vertical_content_list,omitempty"` // 仅news_notice
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list,omitempty"`
	JumpList              []Jump              `json:"jump_list,omitempty"`
	CardAction            CardAction          `json:"card_action"`
}

// TemplateCardMessage 代表企业微信应用的模板卡片消息
// 更多详情请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/90372#%E6%A8%A1%E7%89%88%E5%8D%A1%E7%89%87%E6%B6%88%E6%81%AF
type TemplateCardMessage struct {
	BaseMessage

	TemplateCard TemplateCard `json:"template_card"`
}

// NewTemplateCardMessage creates a new TemplateCardMessage.
func NewTemplateCardMessage(cardType TemplateCardType) *TemplateCardMessage {
	return &TemplateCardMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeTemplateCard,
		},
		TemplateCard: TemplateCard{
			CardType: cardType,
		},
	}
}

// Validate validates the TemplateCardMessage to ensure it meets WeChat Work API requirements.
func (m *TemplateCardMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	// BaseMessage.Validate() already checks for agentid and target fields

	if err := m.validateCardType(); err != nil {
		return err
	}
	if err := m.validateMainTitleAndSubTitle(); err != nil {
		return err
	}
	if err := m.validateCardAction(); err != nil {
		return err
	}
	if err := m.validateHorizontalContentList(); err != nil {
		return err
	}
	if err := m.validateJumpList(); err != nil {
		return err
	}
	if err := m.validateCardTypeSpecific(); err != nil {
		return err
	}
	if err := m.validateQuoteArea(); err != nil {
		return err
	}
	if err := m.validateImageTextArea(); err != nil {
		return err
	}
	if err := m.validateSource(); err != nil {
		return err
	}
	return nil
}

// validateCardType validates the card type.
func (m *TemplateCardMessage) validateCardType() error {
	validTypes := []TemplateCardType{
		CardTypeTextNotice,
		CardTypeNewsNotice,
		CardTypeButtonInteraction,
		CardTypeVoteInteraction,
	}

	for _, validType := range validTypes {
		if m.TemplateCard.CardType == validType {
			return nil
		}
	}

	return core.NewParamError(
		fmt.Sprintf(
			"invalid card_type: %s; must be one of: %s, %s, %s, %s",
			m.TemplateCard.CardType,
			CardTypeTextNotice,
			CardTypeNewsNotice,
			CardTypeButtonInteraction,
			CardTypeVoteInteraction,
		),
	)
}

func (m *TemplateCardMessage) validateMainTitleAndSubTitle() error {
	if m.TemplateCard.MainTitle.Title == "" && m.TemplateCard.SubTitleText == "" {
		return core.NewParamError("main_title.title or sub_title_text must be provided")
	}
	return nil
}

func (m *TemplateCardMessage) validateCardAction() error {
	if m.TemplateCard.CardAction.Type == 0 {
		return core.NewParamError("card_action.type is required")
	}
	if m.TemplateCard.CardAction.Type == 1 && m.TemplateCard.CardAction.URL == "" {
		return core.NewParamError("card_action.url is required for type 1")
	}
	if m.TemplateCard.CardAction.Type == 2 &&
		(m.TemplateCard.CardAction.AppID == "" || m.TemplateCard.CardAction.PagePath == "") {
		return core.NewParamError("card_action.appid and pagepath are required for type 2")
	}
	if m.TemplateCard.CardAction.Type < 0 || m.TemplateCard.CardAction.Type > 2 {
		return core.NewParamError(
			fmt.Sprintf("invalid card_action.type: %d; must be 0, 1, or 2", m.TemplateCard.CardAction.Type),
		)
	}
	return nil
}

func (m *TemplateCardMessage) validateHorizontalContentList() error {
	for i, item := range m.TemplateCard.HorizontalContentList {
		if item.Keyname == "" {
			return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: keyname cannot be empty", i))
		}
		if item.Value == "" {
			return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: value cannot be empty", i))
		}
		switch item.Type {
		case HorizontalTypeURL:
			if item.URL == "" {
				return core.NewParamError(fmt.Sprintf("horizontal_content_list item %d: url is required for type 1", i))
			}
		case HorizontalTypeMedia:
			if item.MediaID == "" {
				return core.NewParamError(
					fmt.Sprintf(
						"horizontal_content_list item %d: media_id is required for type 2",
						i,
					),
				)
			}
		case HorizontalTypeDefault, HorizontalTypeMember:
			// valid, do nothing
		default:
			return core.NewParamError(
				fmt.Sprintf("horizontal_content_list item %d: invalid type %d; must be 0, 1, 2, or 3", i, item.Type),
			)
		}
	}
	return nil
}

func (m *TemplateCardMessage) validateJumpList() error {
	for i, jump := range m.TemplateCard.JumpList {
		if jump.Type == 0 || jump.Title == "" {
			return core.NewParamError(fmt.Sprintf("jump_list item %d: type and title are required", i))
		}
		if jump.Type == 1 && jump.URL == "" {
			return core.NewParamError(fmt.Sprintf("jump_list item %d: url is required for type 1", i))
		}
		if jump.Type == 2 && (jump.AppID == "" || jump.PagePath == "") {
			return core.NewParamError(fmt.Sprintf("jump_list item %d: appid and pagepath are required for type 2", i))
		}
		if jump.Type < 0 || jump.Type > 2 {
			return core.NewParamError(
				fmt.Sprintf("jump_list item %d: invalid type: %d; must be 0, 1, or 2", i, jump.Type),
			)
		}
	}
	return nil
}

func (m *TemplateCardMessage) validateCardTypeSpecific() error {
	switch m.TemplateCard.CardType {
	case CardTypeNewsNotice:
		if m.TemplateCard.CardImage == nil || m.TemplateCard.CardImage.URL == "" {
			return core.NewParamError("card_image.url is required for news_notice cards")
		}
		if m.TemplateCard.CardImage.AspectRatio == 0 {
			return core.NewParamError("card_image.aspect_ratio cannot be zero for news_notice cards")
		}
		for i, vc := range m.TemplateCard.VerticalContentList {
			if vc.Title == "" {
				return core.NewParamError(fmt.Sprintf("vertical_content_list item %d: title cannot be empty", i))
			}
		}
	case CardTypeTextNotice:
		if m.TemplateCard.CardImage != nil {
			return core.NewParamError("card_image is not allowed for text_notice cards")
		}
		if m.TemplateCard.ImageTextArea != nil {
			return core.NewParamError("image_text_area is not allowed for text_notice cards")
		}
		if len(m.TemplateCard.VerticalContentList) > 0 {
			return core.NewParamError("vertical_content_list is not allowed for text_notice cards")
		}
	case CardTypeButtonInteraction:
		// Button interaction cards have their own validation logic
		// TODO: Add specific validation for button interaction cards if needed
	case CardTypeVoteInteraction:
		// Vote interaction cards have their own validation logic
		// TODO: Add specific validation for vote interaction cards if needed
	}
	return nil
}

func (m *TemplateCardMessage) validateQuoteArea() error {
	if m.TemplateCard.QuoteArea != nil {
		if m.TemplateCard.QuoteArea.Type < 0 || m.TemplateCard.QuoteArea.Type > 2 {
			return core.NewParamError(
				fmt.Sprintf("quote_area: invalid type: %d; must be 0, 1, or 2", m.TemplateCard.QuoteArea.Type),
			)
		}
		if m.TemplateCard.QuoteArea.Title == "" || m.TemplateCard.QuoteArea.QuoteText == "" {
			return core.NewParamError("quote_area must have title and quote_text")
		}
		if m.TemplateCard.QuoteArea.Type == 1 && m.TemplateCard.QuoteArea.URL == "" {
			return core.NewParamError("quote_area.url is required for type 1")
		}
		if m.TemplateCard.QuoteArea.Type == 2 &&
			(m.TemplateCard.QuoteArea.AppID == "" || m.TemplateCard.QuoteArea.PagePath == "") {
			return core.NewParamError("quote_area.appid and pagepath are required for type 2")
		}
	}
	return nil
}

func (m *TemplateCardMessage) validateImageTextArea() error {
	if m.TemplateCard.ImageTextArea != nil {
		if m.TemplateCard.CardType != CardTypeNewsNotice {
			return core.NewParamError("image_text_area is only allowed for news_notice cards")
		}
		if m.TemplateCard.ImageTextArea.Type == 0 || m.TemplateCard.ImageTextArea.Title == "" ||
			m.TemplateCard.ImageTextArea.ImageURL == "" {
			return core.NewParamError("image_text_area must have type, title, and image_url")
		}
		if m.TemplateCard.ImageTextArea.Type == 1 && m.TemplateCard.ImageTextArea.URL == "" {
			return core.NewParamError("image_text_area.url is required for type 1")
		}
		if m.TemplateCard.ImageTextArea.Type < 0 ||
			m.TemplateCard.ImageTextArea.Type > 1 {
			return core.NewParamError(
				fmt.Sprintf("image_text_area: invalid type: %d; must be 0 or 1", m.TemplateCard.ImageTextArea.Type),
			)
		}
	}
	return nil
}

func (m *TemplateCardMessage) validateSource() error {
	if m.TemplateCard.Source != nil {
		if m.TemplateCard.Source.IconURL == "" {
			return core.NewParamError("source.icon_url is required if source is provided")
		}
		if m.TemplateCard.Source.DescColor < 0 || m.TemplateCard.Source.DescColor > 3 {
			return core.NewParamError(
				fmt.Sprintf("invalid source.desc_color: %d; must be 0, 1, 2, or 3", m.TemplateCard.Source.DescColor),
			)
		}
	}
	return nil
}
