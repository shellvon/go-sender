package wecombot

import "github.com/shellvon/go-sender/core"

const maxArticlesPerNewsMessage = 8

// Article 表示新闻消息中的单个文章。
type Article struct {
	// 文章标题。最大 128 字节，超出部分将被截断。
	Title string `json:"title"`
	// 文章描述。最大 512 字节，超出部分将被截断。
	Description string `json:"description"`
	// 点击文章时跳转的 URL。
	URL string `json:"url"`
	// 文章图片的 URL。支持 JPG 和 PNG 格式。
	// 推荐尺寸：大图 1068*455，小图 150*150。
	Picurl string `json:"picurl"`
}

// NewsContent 封装文章列表；重命名以避免与 News() 构建器工厂冲突。
type NewsContent struct {
	Articles []*Article `json:"articles"`
}

// NewsMessage 表示企业微信的新闻消息。
// 更多详情，请参考企业微信 API 文档。
type NewsMessage struct {
	BaseMessage

	News NewsContent `json:"news"`
}

// NewNewsMessage 创建一个新的 NewsMessage 实例。
// 基于企业微信机器人 API 的 SendNewsParams
// 参考：https://developer.work.weixin.qq.com/document/path/91770#%E6%96%B0%E9%97%BB%E7%B1%BB%E5%9E%8B
//   - 仅 articles 是必需的。
//
// 参数：articles []*Article - 新闻消息的文章列表。
// 返回值：*NewsMessage - 新创建的新闻消息实例。
func NewNewsMessage(articles []*Article) *NewsMessage {
	return News().Articles(articles).Build()
}

// Validate 验证 NewsMessage 是否满足企业微信 API 的要求。
// 该方法检查文章列表是否为空、文章数量是否在 1 到 8 之间，以及每篇文章是否包含必需的标题和 URL。
// 返回值：error - 如果验证失败，返回具体的参数错误；否则返回 nil。
func (m *NewsMessage) Validate() error {
	if len(m.News.Articles) == 0 {
		return core.NewParamError("文章列表不能为空")
	}
	if len(m.News.Articles) > maxArticlesPerNewsMessage {
		return core.NewParamError("新闻消息支持 1 到 8 篇文章")
	}
	for _, article := range m.News.Articles {
		if article.URL == "" || article.Title == "" {
			return core.NewParamError("文章缺少必需的标题或 URL")
		}
	}
	return nil
}
