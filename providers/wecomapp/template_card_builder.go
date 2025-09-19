package wecomapp

// TemplateCardBuilder 为创建企业微信应用template_card消息提供流畅的API
// 它支持text_notice和news_notice卡片类型，并提供全面的配置选项
//
// text_notice示例:
//   msg := wecomapp.NewTemplateCardBuilder(wecomapp.CardTypeTextNotice).
//            MainTitle("Server Alert", "High CPU Usage").
//            SubTitle("Click to view details").
//            JumpURL("https://example.com/alert/123").
//            ToUser("@all").
//            Build()
//
// news_notice示例:
//   msg := wecomapp.NewTemplateCardBuilder(wecomapp.CardTypeNewsNotice).
//            MainTitle("Product Update", "New features available").
//            CardImage("https://example.com/image.jpg", 1.3).
//            AddVerticalContent("Feature 1", "Description of feature 1").
//            JumpURL("https://example.com/updates").
//            ToUser("user1|user2").
//            Build()

type TemplateCardBuilder struct {
	message *TemplateCardMessage
}

// NewTemplateCardBuilder 使用给定的卡片类型创建新的TemplateCardBuilder.
func NewTemplateCardBuilder(cardType TemplateCardType) *TemplateCardBuilder {
	return &TemplateCardBuilder{
		message: NewTemplateCardMessage(cardType),
	}
}

// ToUser sets the user IDs to send to, separated by "|". Use "@all" to send to all users.
func (b *TemplateCardBuilder) ToUser(users string) *TemplateCardBuilder {
	b.message.CommonFields.ToUser = users
	return b
}

// ToParty sets the department IDs to send to, separated by "|".
func (b *TemplateCardBuilder) ToParty(parties string) *TemplateCardBuilder {
	b.message.CommonFields.ToParty = parties
	return b
}

// ToTag sets the tag IDs to send to, separated by "|".
func (b *TemplateCardBuilder) ToTag(tags string) *TemplateCardBuilder {
	b.message.CommonFields.ToTag = tags
	return b
}

// AgentID sets the application ID.
func (b *TemplateCardBuilder) AgentID(agentID string) *TemplateCardBuilder {
	b.message.CommonFields.AgentID = agentID
	return b
}

// Safe sets whether to enable safe mode (0: no, 1: yes).
func (b *TemplateCardBuilder) Safe(safe int) *TemplateCardBuilder {
	b.message.CommonFields.Safe = safe
	return b
}

// EnableIDTrans sets whether to enable ID translation (0: no, 1: yes).
func (b *TemplateCardBuilder) EnableIDTrans(enable int) *TemplateCardBuilder {
	b.message.CommonFields.EnableIDTrans = enable
	return b
}

// EnableDuplicateCheck sets whether to enable duplicate message check (0: no, 1: yes).
func (b *TemplateCardBuilder) EnableDuplicateCheck(enable int) *TemplateCardBuilder {
	b.message.CommonFields.EnableDuplicateCheck = enable
	return b
}

// DuplicateCheckInterval sets the duplicate check interval in seconds.
func (b *TemplateCardBuilder) DuplicateCheckInterval(interval int) *TemplateCardBuilder {
	b.message.CommonFields.DuplicateCheckInterval = interval
	return b
}

// MainTitle sets the main title section (title + desc).
func (b *TemplateCardBuilder) MainTitle(title, desc string) *TemplateCardBuilder {
	b.message.TemplateCard.MainTitle = MainTitle{
		Title: title,
		Desc:  desc,
	}
	return b
}

// SubTitle sets optional sub title text (text_notice only).
func (b *TemplateCardBuilder) SubTitle(text string) *TemplateCardBuilder {
	b.message.TemplateCard.SubTitleText = text
	return b
}

// Source sets the source section with icon URL, description, and color.
func (b *TemplateCardBuilder) Source(iconURL, desc string, descColor int) *TemplateCardBuilder {
	b.message.TemplateCard.Source = &Source{
		IconURL:   iconURL,
		Desc:      desc,
		DescColor: descColor,
	}
	return b
}

// EmphasisContent sets the emphasized content section (text_notice only).
func (b *TemplateCardBuilder) EmphasisContent(title, desc string) *TemplateCardBuilder {
	b.message.TemplateCard.EmphasisContent = &EmphasisContent{
		Title: title,
		Desc:  desc,
	}
	return b
}

// QuoteArea sets the quote area section.
func (b *TemplateCardBuilder) QuoteArea(title, quoteText string, clickType int) *TemplateCardBuilder {
	b.message.TemplateCard.QuoteArea = &QuoteArea{
		Type:      clickType,
		Title:     title,
		QuoteText: quoteText,
	}
	return b
}

// QuoteAreaWithURL sets the quote area section with URL jump.
func (b *TemplateCardBuilder) QuoteAreaWithURL(title, quoteText, url string) *TemplateCardBuilder {
	b.message.TemplateCard.QuoteArea = &QuoteArea{
		Type:      1,
		Title:     title,
		QuoteText: quoteText,
		URL:       url,
	}
	return b
}

// QuoteAreaWithMiniProgram sets the quote area section with mini-program jump.
func (b *TemplateCardBuilder) QuoteAreaWithMiniProgram(title, quoteText, appID, pagePath string) *TemplateCardBuilder {
	const miniProgramType = 2
	b.message.TemplateCard.QuoteArea = &QuoteArea{
		Type:      miniProgramType,
		Title:     title,
		QuoteText: quoteText,
		AppID:     appID,
		PagePath:  pagePath,
	}
	return b
}

// CardImage sets the card image section (news_notice only).
// aspectRatio 图片的宽高比.
func (b *TemplateCardBuilder) CardImage(url string, aspectRatio float64) *TemplateCardBuilder {
	b.message.TemplateCard.CardImage = &CardImage{
		URL:         url,
		AspectRatio: aspectRatio,
	}
	return b
}

// ImageTextArea sets the image text area section (news_notice only).
func (b *TemplateCardBuilder) ImageTextArea(title, desc, imageURL string, clickType int) *TemplateCardBuilder {
	b.message.TemplateCard.ImageTextArea = &ImageTextArea{
		Type:     clickType,
		Title:    title,
		Desc:     desc,
		ImageURL: imageURL,
	}
	return b
}

// ImageTextAreaWithURL sets the image text area section with URL jump (news_notice only).
func (b *TemplateCardBuilder) ImageTextAreaWithURL(title, desc, imageURL, jumpURL string) *TemplateCardBuilder {
	b.message.TemplateCard.ImageTextArea = &ImageTextArea{
		Type:     1,
		URL:      jumpURL,
		Title:    title,
		Desc:     desc,
		ImageURL: imageURL,
	}
	return b
}

// AddVerticalContent adds a vertical content item (news_notice only).
func (b *TemplateCardBuilder) AddVerticalContent(title, desc string) *TemplateCardBuilder {
	b.message.TemplateCard.VerticalContentList = append(b.message.TemplateCard.VerticalContentList, VerticalContent{
		Title: title,
		Desc:  desc,
	})
	return b
}

// AddHorizontalContent adds a horizontal content item with default type.
func (b *TemplateCardBuilder) AddHorizontalContent(keyname, value string) *TemplateCardBuilder {
	b.message.TemplateCard.HorizontalContentList = append(
		b.message.TemplateCard.HorizontalContentList,
		HorizontalContent{
			Keyname: keyname,
			Value:   value,
			Type:    HorizontalTypeDefault,
		},
	)
	return b
}

// AddHorizontalContentWithURL adds a horizontal content item with URL jump.
func (b *TemplateCardBuilder) AddHorizontalContentWithURL(keyname, value, url string) *TemplateCardBuilder {
	b.message.TemplateCard.HorizontalContentList = append(
		b.message.TemplateCard.HorizontalContentList,
		HorizontalContent{
			Keyname: keyname,
			Value:   value,
			Type:    HorizontalTypeURL,
			URL:     url,
		},
	)
	return b
}

// AddHorizontalContentWithMedia adds a horizontal content item with media attachment.
func (b *TemplateCardBuilder) AddHorizontalContentWithMedia(keyname, value, mediaID string) *TemplateCardBuilder {
	b.message.TemplateCard.HorizontalContentList = append(
		b.message.TemplateCard.HorizontalContentList,
		HorizontalContent{
			Keyname: keyname,
			Value:   value,
			Type:    HorizontalTypeMedia,
			MediaID: mediaID,
		},
	)
	return b
}

// AddHorizontalContentWithMember adds a horizontal content item with member details jump.
func (b *TemplateCardBuilder) AddHorizontalContentWithMember(keyname, value string) *TemplateCardBuilder {
	b.message.TemplateCard.HorizontalContentList = append(
		b.message.TemplateCard.HorizontalContentList,
		HorizontalContent{
			Keyname: keyname,
			Value:   value,
			Type:    HorizontalTypeMember,
		},
	)
	return b
}

// AddJump adds a jump item with the specified type and title.
func (b *TemplateCardBuilder) AddJump(title string, jumpType int) *TemplateCardBuilder {
	b.message.TemplateCard.JumpList = append(b.message.TemplateCard.JumpList, Jump{
		Type:  jumpType,
		Title: title,
	})
	return b
}

// AddJumpWithURL adds a jump item with URL.
func (b *TemplateCardBuilder) AddJumpWithURL(title, url string) *TemplateCardBuilder {
	b.message.TemplateCard.JumpList = append(b.message.TemplateCard.JumpList, Jump{
		Type:  1,
		Title: title,
		URL:   url,
	})
	return b
}

// AddJumpWithMiniProgram adds a jump item with mini-program.
func (b *TemplateCardBuilder) AddJumpWithMiniProgram(title, appID, pagePath string) *TemplateCardBuilder {
	const miniProgramType = 2
	b.message.TemplateCard.JumpList = append(b.message.TemplateCard.JumpList, Jump{
		Type:     miniProgramType,
		Title:    title,
		AppID:    appID,
		PagePath: pagePath,
	})
	return b
}

// JumpURL configures a URL jump action for the card.
func (b *TemplateCardBuilder) JumpURL(url string) *TemplateCardBuilder {
	b.message.TemplateCard.CardAction = CardAction{
		Type: 1,
		URL:  url,
	}
	return b
}

// JumpMiniProgram configures a mini-program jump action for the card.
func (b *TemplateCardBuilder) JumpMiniProgram(appID, pagePath string) *TemplateCardBuilder {
	const miniProgramType = 2
	b.message.TemplateCard.CardAction = CardAction{
		Type:     miniProgramType,
		AppID:    appID,
		PagePath: pagePath,
	}
	return b
}

// CardAction sets the card action with custom configuration.
func (b *TemplateCardBuilder) CardAction(actionType int, url, appID, pagePath string) *TemplateCardBuilder {
	b.message.TemplateCard.CardAction = CardAction{
		Type:     actionType,
		URL:      url,
		AppID:    appID,
		PagePath: pagePath,
	}
	return b
}

// Build assembles and returns the TemplateCardMessage.
func (b *TemplateCardBuilder) Build() *TemplateCardMessage {
	return b.message
}
