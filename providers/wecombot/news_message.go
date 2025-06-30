package wecombot

import "github.com/shellvon/go-sender/core"

const maxArticlesPerNewsMessage = 8

// Article represents a single article within a news message.
type Article struct {
	// Title of the article. Max 128 bytes, will be truncated if longer.
	Title string `json:"title"`
	// Description of the article. Max 512 bytes, will be truncated if longer.
	Description string `json:"description"`
	// The URL to jump to when the article is clicked.
	URL string `json:"url"`
	// URL of the image for the article. Supports JPG and PNG formats.
	// Recommended sizes: large image 1068*455, small image 150*150.
	Picurl string `json:"picurl"`
}

// News contains a list of articles for a news message.
type News struct {
	Articles []*Article `json:"articles"`
}

// NewsMessage represents a news message for WeCom.
// For more details, refer to the WeCom API documentation.
type NewsMessage struct {
	BaseMessage

	News News `json:"news"`
}

// NewNewsMessage creates a new NewsMessage with required articles and applies optional configurations.
func NewNewsMessage(articles []*Article, opts ...NewsMessageOption) *NewsMessage {
	msg := &NewsMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeNews,
		},
		News: News{
			Articles: articles,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// Validate validates the NewsMessage to ensure it meets WeCom API requirements.
func (m *NewsMessage) Validate() error {
	if len(m.News.Articles) == 0 {
		return core.NewParamError("articles cannot be empty")
	}
	if len(m.News.Articles) > maxArticlesPerNewsMessage {
		return core.NewParamError("a news message supports 1 to 8 articles")
	}
	for _, article := range m.News.Articles {
		if article.URL == "" || article.Title == "" {
			return core.NewParamError("missing required title or URL for an article")
		}
	}
	return nil
}

// NewsMessageOption defines a function type for configuring NewsMessage.
type NewsMessageOption func(*NewsMessage)

// WithArticles sets the Articles for NewsMessage.
func WithArticles(articles []*Article) NewsMessageOption {
	return func(m *NewsMessage) {
		m.News.Articles = articles
	}
}
