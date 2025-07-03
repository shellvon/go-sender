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

	return &TemplateCardMessage{
		BaseMessage:  BaseMessage{MsgType: TypeTemplateCard},
		TemplateCard: tpl,
	}
}
