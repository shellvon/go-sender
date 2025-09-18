package wecombot

// TemplateCardBuilder 提供用于创建企业微信 `template_card` 消息的流式API
// 专注于最常用的字段：卡片类型、主标题、副标题和URL跳转动作
//
// 示例：
//   msg := wecombot.Card(wecombot.CardTypeTextNotice).
//            MainTitle("服务器告警", "CPU使用率过高").
//            SubTitle("点击查看详情").
//            JumpURL("https://example.com/alert/123").
//            Build()

type TemplateCardBuilder struct {
	cardType TemplateCardType

	mainTitle MainTitle
	subTitle  string

	// 用于 news_notice 类型
	cardImage           *CardImage
	imageTextArea       *ImageTextArea
	verticalContentList []VerticalContent

	// 目前仅支持URL跳转
	jumpURL string
}

// Card 使用给定的卡片类型创建新的TemplateCardBuilder.
func Card(cardType TemplateCardType) *TemplateCardBuilder {
	return &TemplateCardBuilder{cardType: cardType}
}

// MainTitle 设置主标题部分（标题 + 描述）.
func (b *TemplateCardBuilder) MainTitle(title, desc string) *TemplateCardBuilder {
	b.mainTitle = MainTitle{Title: title, Desc: desc}
	return b
}

// SubTitle 设置可选的副标题文本（仅text_notice类型支持）.
func (b *TemplateCardBuilder) SubTitle(text string) *TemplateCardBuilder {
	b.subTitle = text
	return b
}

// CardImage 设置卡片图片部分（仅news_notice类型支持）
// aspectRatio 是图片的宽高比.
func (b *TemplateCardBuilder) CardImage(url string, aspectRatio float64) *TemplateCardBuilder {
	b.cardImage = &CardImage{
		URL:         url,
		AspectRatio: aspectRatio,
	}
	return b
}

// ImageTextArea 设置左图右文区域（仅news_notice类型支持）.
func (b *TemplateCardBuilder) ImageTextArea(title, desc, imageURL, jumpURL string) *TemplateCardBuilder {
	b.imageTextArea = &ImageTextArea{
		Type:     1, // URL类型
		URL:      jumpURL,
		Title:    title,
		Desc:     desc,
		ImageURL: imageURL,
	}
	return b
}

// AddVerticalContent 添加垂直内容项（仅news_notice类型支持）.
func (b *TemplateCardBuilder) AddVerticalContent(title, desc string) *TemplateCardBuilder {
	b.verticalContentList = append(b.verticalContentList, VerticalContent{
		Title: title,
		Desc:  desc,
	})
	return b
}

// JumpURL 配置URL跳转动作.
func (b *TemplateCardBuilder) JumpURL(url string) *TemplateCardBuilder {
	b.jumpURL = url
	return b
}

// Build 构建TemplateCardMessage.
func (b *TemplateCardBuilder) Build() *TemplateCardMessage {
	// 默认卡片动作（URL跳转）
	action := CardAction{Type: 1, URL: b.jumpURL}

	tpl := TemplateCard{
		CardType:     b.cardType,
		MainTitle:    b.mainTitle,
		SubTitleText: b.subTitle,
		CardAction:   action,
	}

	// 添加news_notice特有字段
	if b.cardType == CardTypeNewsNotice {
		tpl.CardImage = b.cardImage
		tpl.ImageTextArea = b.imageTextArea
		if len(b.verticalContentList) > 0 {
			tpl.VerticalContentList = b.verticalContentList
		}
	}

	return &TemplateCardMessage{
		BaseMessage:  newBaseMessage(TypeTemplateCard),
		TemplateCard: tpl,
	}
}
