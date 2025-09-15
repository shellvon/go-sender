package wecombot

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// TemplateCardType 定义模板卡片允许的类型.
type TemplateCardType string

const (
	// CardTypeTextNotice 文本通知模板卡片.
	CardTypeTextNotice TemplateCardType = "text_notice"
	// CardTypeNewsNotice 图文展示模板卡片.
	CardTypeNewsNotice = "news_notice"

	// 水平内容列表相关常量.
	maxHorizontalContentList = 6 // 水平内容列表最大项数
	maxJumpList              = 5 // 跳转列表最大项数

	// 水平内容类型常量.
	HorizontalTypeDefault = 0 // 默认类型
	HorizontalTypeURL     = 1 // URL跳转类型
	HorizontalTypeMedia   = 2 // 文件附件类型
	HorizontalTypeMember  = 3 // 成员详情跳转类型
)

// Source 模板卡片的来源信息.
type Source struct {
	// IconURL 来源图标的URL
	IconURL string `json:"icon_url"`
	// Desc 来源的描述文字
	Desc string `json:"desc"`
	// DescColor 来源文字的颜色，支持的值：0(默认)灰色，1黑色，2红色，3绿色
	DescColor int `json:"desc_color"`
}

// MainTitle 模板卡片的一级标题.
type MainTitle struct {
	// Title 一级标题，模板卡片必须有main_title.title或sub_title_text中的一个
	Title string `json:"title"`
	// Desc 标题辅助信息
	Desc string `json:"desc"`
}

// EmphasisContent 二级普通文本（仅text_notice类型支持）.
type EmphasisContent struct {
	// Title 关键数据样式的数据内容
	Title string `json:"title"`
	// Desc 关键数据样式的描述内容
	Desc string `json:"desc"`
}

// QuoteArea 引用文献样式，一般建议不与二级普通文本共用.
type QuoteArea struct {
	// Type 点击跳转类型：0或不填代表不是链接，1代表跳转URL，2代表跳转小程序
	Type int `json:"type"`
	// URL 跳转链接，当quote_area.type为1时必填
	URL string `json:"url,omitempty"`
	// AppID 小程序的appid，当quote_area.type为2时必填
	AppID string `json:"appid,omitempty"`
	// PagePath 小程序的页面路径，当quote_area.type为2时选填
	PagePath string `json:"pagepath,omitempty"`
	// Title 引用文献样式区域点击事件的文本
	Title string `json:"title"`
	// QuoteText 引用文献样式的文本
	QuoteText string `json:"quote_text"`
}

// CardImage 图片样式（仅news_notice类型支持）.
type CardImage struct {
	// URL 图片的URL
	URL string `json:"url"`
	// AspectRatio 图片的宽高比
	AspectRatio float64 `json:"aspect_ratio"`
}

// ImageTextArea 左图右文样式（仅news_notice类型支持）.
type ImageTextArea struct {
	// Type 左图右文样式区域点击事件类型
	Type int `json:"type"`
	// URL 点击跳转的URL
	URL string `json:"url"`
	// Title 左图右文样式的标题
	Title string `json:"title"`
	// Desc 左图右文样式的描述
	Desc string `json:"desc"`
	// ImageURL 左图右文样式的图片URL
	ImageURL string `json:"image_url"`
}

// HorizontalContent 二级标题+文本列表的单个项目
// 该字段可为空数组，但有数据的话需确保对应字段的合法性，列表长度不超过6个.
type HorizontalContent struct {
	// Keyname 二级标题
	Keyname string `json:"keyname"`
	// Value 二级文本，如果horizontal_content_list.type是2，该字段代表文件名(要包含文件类型)
	Value string `json:"value"`
	// Type 二级标题信息类型：1跳转URL，2文件附件，3跳转成员详情
	Type int `json:"type,omitempty"`
	// URL 链接跳转的URL，horizontal_content_list.type是1时必填
	URL string `json:"url,omitempty"`
	// MediaID 附件的media_id，horizontal_content_list.type是2时必填
	MediaID string `json:"media_id,omitempty"`
}

// VerticalContent 卡片二级垂直内容（仅news_notice类型支持）.
type VerticalContent struct {
	// Title 卡片二级标题
	Title string `json:"title"`
	// Desc 二级普通文本
	Desc string `json:"desc"`
}

// Jump 跳转指引样式的列表项.
type Jump struct {
	// Type 跳转链接类型
	Type int `json:"type"`
	// URL 跳转链接，Type为1时必填
	URL string `json:"url,omitempty"`
	// AppID 小程序appid，Type为2时必填
	AppID string `json:"appid,omitempty"`
	// PagePath 小程序页面路径，Type为2时选填
	PagePath string `json:"pagepath,omitempty"`
	// Title 跳转链接文案
	Title string `json:"title"`
}

// CardAction 整体卡片的点击事件.
type CardAction struct {
	// Type 卡片跳转类型
	Type int `json:"type"`
	// URL 跳转事件的URL，Type为1时必填
	URL string `json:"url,omitempty"`
	// AppID 跳转事件的小程序的appid，Type为2时必填
	AppID string `json:"appid,omitempty"`
	// PagePath 跳转事件的小程序的页面路径，Type为2时选填
	PagePath string `json:"pagepath,omitempty"`
}

// TemplateCard 模板卡片消息的template_card字段.
type TemplateCard struct {
	CardType              TemplateCardType    `json:"card_type"`                         // 模板卡片类型
	Source                *Source             `json:"source"`                            // 卡片来源样式信息
	MainTitle             MainTitle           `json:"main_title"`                        // 模板卡片的一级标题
	EmphasisContent       *EmphasisContent    `json:"emphasis_content,omitempty"`        // 关键数据样式（仅text_notice支持）
	QuoteArea             *QuoteArea          `json:"quote_area,omitempty"`              // 引用文献样式
	SubTitleText          string              `json:"sub_title_text,omitempty"`          // 二级普通文本（仅text_notice支持）
	CardImage             *CardImage          `json:"card_image,omitempty"`              // 图片样式（仅news_notice支持）
	ImageTextArea         *ImageTextArea      `json:"image_text_area,omitempty"`         // 左图右文样式（仅news_notice支持）
	VerticalContentList   []VerticalContent   `json:"vertical_content_list,omitempty"`   // 卡片二级垂直内容（仅news_notice支持）
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list,omitempty"` // 二级标题+文本列表
	JumpList              []Jump              `json:"jump_list,omitempty"`               // 跳转指引样式的列表
	CardAction            CardAction          `json:"card_action"`                       // 整体卡片的点击事件
}

// TemplateCardMessage 企业微信模板卡片消息
// 详细信息请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/91770#%E6%A8%A1%E7%89%88%E5%8D%A1%E7%89%87%E7%B1%BB%E5%9E%8B
type TemplateCardMessage struct {
	BaseMessage

	TemplateCard TemplateCard `json:"template_card"`
}

// Validate 根据企业微信API规则验证模板卡片消息.
func (m *TemplateCardMessage) Validate() error {
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

// validateCardType 验证卡片类型.
func (m *TemplateCardMessage) validateCardType() error {
	if m.TemplateCard.CardType != CardTypeTextNotice && m.TemplateCard.CardType != CardTypeNewsNotice {
		return core.NewParamError(
			fmt.Sprintf(
				"无效的card_type: %s，必须是 %s 或 %s",
				m.TemplateCard.CardType,
				CardTypeTextNotice,
				CardTypeNewsNotice,
			),
		)
	}
	return nil
}

// validateMainTitleAndSubTitle 验证主标题和副标题.
func (m *TemplateCardMessage) validateMainTitleAndSubTitle() error {
	if m.TemplateCard.MainTitle.Title == "" && m.TemplateCard.SubTitleText == "" {
		return core.NewParamError("必须提供 main_title.title 或 sub_title_text 中的一个")
	}
	return nil
}

// validateCardAction 验证卡片动作.
func (m *TemplateCardMessage) validateCardAction() error {
	if m.TemplateCard.CardAction.Type == 0 {
		return core.NewParamError("card_action.type 是必填的")
	}
	if m.TemplateCard.CardAction.Type == 1 && m.TemplateCard.CardAction.URL == "" {
		return core.NewParamError("当 card_action.type 为 1 时，card_action.url 是必填的")
	}
	if m.TemplateCard.CardAction.Type == 2 &&
		(m.TemplateCard.CardAction.AppID == "" || m.TemplateCard.CardAction.PagePath == "") {
		return core.NewParamError("当 card_action.type 为 2 时，card_action.appid 和 pagepath 是必填的")
	}
	if m.TemplateCard.CardAction.Type < 0 || m.TemplateCard.CardAction.Type > 2 {
		return core.NewParamError(
			fmt.Sprintf("无效的 card_action.type: %d，必须是 0、1 或 2", m.TemplateCard.CardAction.Type),
		)
	}
	return nil
}

// validateHorizontalContentList 验证水平内容列表.
func (m *TemplateCardMessage) validateHorizontalContentList() error {
	if len(m.TemplateCard.HorizontalContentList) > maxHorizontalContentList {
		return core.NewParamError("horizontal_content_list 不能超过 6 个项目")
	}
	for i, item := range m.TemplateCard.HorizontalContentList {
		if item.Keyname == "" {
			return core.NewParamError(fmt.Sprintf("horizontal_content_list 第 %d 项：keyname 不能为空", i))
		}
		if item.Value == "" {
			return core.NewParamError(fmt.Sprintf("horizontal_content_list 第 %d 项：value 不能为空", i))
		}
		// 移除长度限制检查，让API服务端决策
		switch item.Type {
		case HorizontalTypeURL:
			if item.URL == "" {
				return core.NewParamError(fmt.Sprintf("horizontal_content_list 第 %d 项：当 type 为 1 时，url 是必填的", i))
			}
		case HorizontalTypeMedia:
			if item.MediaID == "" {
				return core.NewParamError(
					fmt.Sprintf(
						"horizontal_content_list 第 %d 项：当 type 为 2 时，media_id 是必填的",
						i,
					),
				)
			}
		case HorizontalTypeDefault, HorizontalTypeMember:
			// 有效类型，无需额外处理
		default:
			return core.NewParamError(
				fmt.Sprintf("horizontal_content_list 第 %d 项：无效的 type %d，必须是 0、1、2 或 3", i, item.Type),
			)
		}
	}
	return nil
}

// validateJumpList 验证跳转列表.
func (m *TemplateCardMessage) validateJumpList() error {
	if len(m.TemplateCard.JumpList) > maxJumpList {
		return core.NewParamError("jump_list 不能超过 5 个项目")
	}
	for i, jump := range m.TemplateCard.JumpList {
		if jump.Type == 0 || jump.Title == "" {
			return core.NewParamError(fmt.Sprintf("jump_list 第 %d 项：type 和 title 是必填的", i))
		}
		if jump.Type == 1 && jump.URL == "" {
			return core.NewParamError(fmt.Sprintf("jump_list 第 %d 项：当 type 为 1 时，url 是必填的", i))
		}
		if jump.Type == 2 && (jump.AppID == "" || jump.PagePath == "") {
			return core.NewParamError(fmt.Sprintf("jump_list 第 %d 项：当 type 为 2 时，appid 和 pagepath 是必填的", i))
		}
		if jump.Type < 0 || jump.Type > 2 {
			return core.NewParamError(
				fmt.Sprintf("jump_list 第 %d 项：无效的 type %d，必须是 0、1 或 2", i, jump.Type),
			)
		}
		// 移除长度限制检查，让API服务端决策
	}
	return nil
}

// validateCardTypeSpecific 验证特定卡片类型的字段.
func (m *TemplateCardMessage) validateCardTypeSpecific() error {
	switch m.TemplateCard.CardType {
	case CardTypeNewsNotice:
		if m.TemplateCard.CardImage == nil || m.TemplateCard.CardImage.URL == "" {
			return core.NewParamError("news_notice 卡片需要 card_image.url")
		}
		if m.TemplateCard.CardImage.AspectRatio == 0 {
			return core.NewParamError("news_notice 卡片的 card_image.aspect_ratio 不能为零")
		}
		for i, vc := range m.TemplateCard.VerticalContentList {
			if vc.Title == "" {
				return core.NewParamError(fmt.Sprintf("vertical_content_list 第 %d 项：title 不能为空", i))
			}
		}
	case CardTypeTextNotice:
		if m.TemplateCard.CardImage != nil {
			return core.NewParamError("text_notice 卡片不允许使用 card_image")
		}
		if m.TemplateCard.ImageTextArea != nil {
			return core.NewParamError("text_notice 卡片不允许使用 image_text_area")
		}
		if len(m.TemplateCard.VerticalContentList) > 0 {
			return core.NewParamError("text_notice 卡片不允许使用 vertical_content_list")
		}
	}
	return nil
}

// validateQuoteArea 验证引用区域.
func (m *TemplateCardMessage) validateQuoteArea() error {
	if m.TemplateCard.QuoteArea != nil {
		if m.TemplateCard.QuoteArea.Type < 0 || m.TemplateCard.QuoteArea.Type > 2 {
			return core.NewParamError(
				fmt.Sprintf("quote_area 无效的 type: %d，必须是 0、1 或 2", m.TemplateCard.QuoteArea.Type),
			)
		}
		if m.TemplateCard.QuoteArea.Title == "" || m.TemplateCard.QuoteArea.QuoteText == "" {
			return core.NewParamError("quote_area 必须有 title 和 quote_text")
		}
		if m.TemplateCard.QuoteArea.Type == 1 && m.TemplateCard.QuoteArea.URL == "" {
			return core.NewParamError("当 quote_area.type 为 1 时，quote_area.url 是必填的")
		}
		if m.TemplateCard.QuoteArea.Type == 2 &&
			(m.TemplateCard.QuoteArea.AppID == "" || m.TemplateCard.QuoteArea.PagePath == "") {
			return core.NewParamError("当 quote_area.type 为 2 时，quote_area.appid 和 pagepath 是必填的")
		}
	}
	return nil
}

// validateImageTextArea 验证图文区域.
func (m *TemplateCardMessage) validateImageTextArea() error {
	if m.TemplateCard.ImageTextArea != nil {
		if m.TemplateCard.CardType != CardTypeNewsNotice {
			return core.NewParamError("image_text_area 只能用于 news_notice 卡片")
		}
		if m.TemplateCard.ImageTextArea.Type == 0 || m.TemplateCard.ImageTextArea.Title == "" ||
			m.TemplateCard.ImageTextArea.ImageURL == "" {
			return core.NewParamError("image_text_area 必须有 type、title 和 image_url")
		}
		if m.TemplateCard.ImageTextArea.Type == 1 && m.TemplateCard.ImageTextArea.URL == "" {
			return core.NewParamError("当 image_text_area.type 为 1 时，image_text_area.url 是必填的")
		}
		if m.TemplateCard.ImageTextArea.Type < 0 ||
			m.TemplateCard.ImageTextArea.Type > 1 {
			return core.NewParamError(
				fmt.Sprintf("image_text_area 无效的 type: %d，必须是 0 或 1", m.TemplateCard.ImageTextArea.Type),
			)
		}
	}
	return nil
}

// validateMainTitleLength 验证主标题长度（已移除长度限制，让API服务端决策）.
func (m *TemplateCardMessage) validateMainTitleLength() error {
	// 移除长度限制检查，让API服务端决策
	return nil
}

// validateSubTitleTextLength 验证副标题长度（已移除长度限制，让API服务端决策）.
func (m *TemplateCardMessage) validateSubTitleTextLength() error {
	// 移除长度限制检查，让API服务端决策
	return nil
}

// validateSource 验证来源信息.
func (m *TemplateCardMessage) validateSource() error {
	if m.TemplateCard.Source != nil {
		if m.TemplateCard.Source.IconURL == "" {
			return core.NewParamError("提供 source 时，source.icon_url 是必填的")
		}
		// 移除长度限制检查，让API服务端决策
		if m.TemplateCard.Source.DescColor < 0 || m.TemplateCard.Source.DescColor > 3 {
			return core.NewParamError(
				fmt.Sprintf("无效的 source.desc_color: %d，必须是 0、1、2 或 3", m.TemplateCard.Source.DescColor),
			)
		}
	}
	return nil
}
