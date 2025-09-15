package wecomapp

// NewsBuilder 提供流畅的API来构造企业微信应用图文消息
//
// 使用示例:
//
//	msg := wecomapp.News().
//	         AddArticle("Title", "Description", "https://example.com", "https://img.com/pic.png").
//	         AddArticle("Title 2", "Description 2", "https://example2.com", "").
//	         ToUser("user1|user2").
//	         AgentID("1000001").
//	         Build()
//
// 图文消息支持1-8篇文章。这遵循与其他provider相同的构建器风格模式以保持一致性
type NewsBuilder struct {
	articles               []*NewsArticle
	toUser                 string
	toParty                string
	toTag                  string
	safe                   int
	enableIDTrans          int
	enableDuplicateCheck   int
	duplicateCheckInterval int
}

// News 创建新的NewsBuilder实例
func News() *NewsBuilder {
	return &NewsBuilder{}
}

// AddArticle 向图文消息添加一篇文章
// 返回构建器以支持链式调用。Build()中强制最大8篇文章
//
// 参数:
//   - title: 文章标题（必需，最大128字节）
//   - description: 文章描述（可选，最大512字节）
//   - url: 点击时跳转的URL（必需）
//   - picURL: 文章图片URL（可选，支持JPG和PNG）
func (b *NewsBuilder) AddArticle(title, description, url, picURL string) *NewsBuilder {
	b.articles = append(b.articles, &NewsArticle{
		Title:       title,
		Description: description,
		URL:         url,
		PicURL:      picURL,
	})
	return b
}

// AddArticleWithMiniprogram 向图文消息添加一篇带小程序跳转的文章
// 返回构建器以支持链式调用。Build()中强制最大8篇文章
//
// 参数:
//   - title: 文章标题（必需，最大128字节）
//   - description: 文章描述（可选，最大512字节）
//   - appID: 小程序应用ID（必需，必须与当前应用关联）
//   - pagePath: 小程序页面路径（必需，最大128字节）
//   - picURL: 文章图片URL（可选，支持JPG和PNG）
func (b *NewsBuilder) AddArticleWithMiniprogram(title, description, appID, pagePath, picURL string) *NewsBuilder {
	b.articles = append(b.articles, &NewsArticle{
		Title:       title,
		Description: description,
		AppID:       appID,
		PagePath:    pagePath,
		PicURL:      picURL,
	})
	return b
}

// Articles 设置完整的文章列表（覆盖之前的文章）
func (b *NewsBuilder) Articles(articles []*NewsArticle) *NewsBuilder {
	b.articles = articles
	return b
}

// ToUser 设置发送给的用户ID，用"|"分隔。使用"@all"发送给所有用户
func (b *NewsBuilder) ToUser(toUser string) *NewsBuilder {
	b.toUser = toUser
	return b
}

// ToParty 设置发送给的部门ID，用"|"分隔
func (b *NewsBuilder) ToParty(toParty string) *NewsBuilder {
	b.toParty = toParty
	return b
}

// ToTag 设置发送给的标签ID，用"|"分隔
func (b *NewsBuilder) ToTag(toTag string) *NewsBuilder {
	b.toTag = toTag
	return b
}

// Safe sets whether to enable safe mode (0: no, 1: yes).
func (b *NewsBuilder) Safe(safe int) *NewsBuilder {
	b.safe = safe
	return b
}

// EnableIDTrans sets whether to enable ID translation (0: no, 1: yes).
func (b *NewsBuilder) EnableIDTrans(enable int) *NewsBuilder {
	b.enableIDTrans = enable
	return b
}

// EnableDuplicateCheck sets whether to enable duplicate message check (0: no, 1: yes).
func (b *NewsBuilder) EnableDuplicateCheck(enable int) *NewsBuilder {
	b.enableDuplicateCheck = enable
	return b
}

// DuplicateCheckInterval sets the duplicate check interval in seconds (max 4 hours).
func (b *NewsBuilder) DuplicateCheckInterval(interval int) *NewsBuilder {
	b.duplicateCheckInterval = interval
	return b
}

// Build assembles a ready-to-send *NewsMessage.
func (b *NewsBuilder) Build() *NewsMessage {
	return &NewsMessage{
		BaseMessage: BaseMessage{
			CommonFields: CommonFields{
				ToUser:                 b.toUser,
				ToParty:                b.toParty,
				ToTag:                  b.toTag,
				Safe:                   b.safe,
				EnableIDTrans:          b.enableIDTrans,
				EnableDuplicateCheck:   b.enableDuplicateCheck,
				DuplicateCheckInterval: b.duplicateCheckInterval,
			},
			MsgType: TypeNews,
		},
		News: NewsMessageContent{
			Articles: b.articles,
		},
	}
}
