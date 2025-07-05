package wecombot

// NewsBuilder provides a fluent API to assemble WeCom news messages.
// Usage:
//   msg := wecombot.News().
//            AddArticle("title","desc","https://example.com","https://img.com/pic.png").
//            AddArticle(...).
//            Build()
//
// A news message supports 1~8 articles.

type NewsBuilder struct {
	articles []*Article
}

// News creates a new NewsBuilder.
// Based on SendNewsParams from WeCom Bot API
// https://developer.work.weixin.qq.com/document/path/91770#%E6%96%B0%E9%97%BB%E7%B1%BB%E5%9E%8B
func News() *NewsBuilder {
	return &NewsBuilder{}
}

// AddArticle appends an article to the news message.
// Returns the builder for chaining. Max 8 articles enforced in Build().
func (b *NewsBuilder) AddArticle(title, desc, url, picURL string) *NewsBuilder {
	b.articles = append(b.articles, &Article{
		Title:       title,
		Description: desc,
		URL:         url,
		Picurl:      picURL,
	})
	return b
}

// Articles sets the complete article slice (overwrites previous).
func (b *NewsBuilder) Articles(arts []*Article) *NewsBuilder {
	b.articles = arts
	return b
}

// Build assembles a *NewsMessage.
func (b *NewsBuilder) Build() *NewsMessage {
	return &NewsMessage{
		BaseMessage: BaseMessage{MsgType: TypeNews},
		News: NewsContent{
			Articles: b.articles,
		},
	}
}
