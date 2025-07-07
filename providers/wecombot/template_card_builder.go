package wecombot

// TemplateCardBuilder provides the sole, fluent API for creating WeCom `template_card` messages.
// It focuses on the most common fields: card type, main title, subtitle and a URL jump action.
//
// Example:
//   msg := wecombot.Card(wecombot.CardTypeTextNotice).
//            MainTitle("Server Alert", "High CPU").
//            SubTitle("点击查看详情").
//            JumpURL("https://example.com/alert/123").
//            Build()

type TemplateCardBuilder struct {
	cardType TemplateCardType

	mainTitle MainTitle
	subTitle  string

	// For news_notice type
	cardImage           *CardImage
	imageTextArea       *ImageTextArea
	verticalContentList []VerticalContent

	// currently support only URL jump
	jumpURL string
}

// Card creates a new TemplateCardBuilder with the given card type.
func Card(cardType TemplateCardType) *TemplateCardBuilder {
	return &TemplateCardBuilder{cardType: cardType}
}

// MainTitle sets the main title section (title + desc).
func (b *TemplateCardBuilder) MainTitle(title, desc string) *TemplateCardBuilder {
	b.mainTitle = MainTitle{Title: title, Desc: desc}
	return b
}

// SubTitle sets optional sub title text (text_notice only).
func (b *TemplateCardBuilder) SubTitle(text string) *TemplateCardBuilder {
	b.subTitle = text
	return b
}

// CardImage sets the card image section (news_notice only).
// aspectRatio is the width/height ratio of the image.
func (b *TemplateCardBuilder) CardImage(url string, aspectRatio float64) *TemplateCardBuilder {
	b.cardImage = &CardImage{
		URL:         url,
		AspectRatio: aspectRatio,
	}
	return b
}

// ImageTextArea sets the image text area section (news_notice only).
func (b *TemplateCardBuilder) ImageTextArea(title, desc, imageURL, jumpURL string) *TemplateCardBuilder {
	b.imageTextArea = &ImageTextArea{
		Type:     1, // URL type
		URL:      jumpURL,
		Title:    title,
		Desc:     desc,
		ImageURL: imageURL,
	}
	return b
}

// AddVerticalContent adds a vertical content item (news_notice only).
func (b *TemplateCardBuilder) AddVerticalContent(title, desc string) *TemplateCardBuilder {
	b.verticalContentList = append(b.verticalContentList, VerticalContent{
		Title: title,
		Desc:  desc,
	})
	return b
}

// JumpURL configures a URL jump action.
func (b *TemplateCardBuilder) JumpURL(url string) *TemplateCardBuilder {
	b.jumpURL = url
	return b
}

// Build assembles a *TemplateCardMessage.
func (b *TemplateCardBuilder) Build() *TemplateCardMessage {
	// default card action (URL)
	action := CardAction{Type: 1, URL: b.jumpURL}

	tpl := TemplateCard{
		CardType:     b.cardType,
		MainTitle:    b.mainTitle,
		SubTitleText: b.subTitle,
		CardAction:   action,
	}

	// Add news_notice specific fields
	if b.cardType == CardTypeNewsNotice {
		tpl.CardImage = b.cardImage
		tpl.ImageTextArea = b.imageTextArea
		if len(b.verticalContentList) > 0 {
			tpl.VerticalContentList = b.verticalContentList
		}
	}

	return &TemplateCardMessage{
		BaseMessage:  BaseMessage{MsgType: TypeTemplateCard},
		TemplateCard: tpl,
	}
}
