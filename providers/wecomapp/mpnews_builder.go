package wecomapp

// MPNewsBuilder 提供流畅的API来构造企业微信应用mpnews消息
//
// 使用示例:
//
//	msg := wecomapp.MPNews().
//	         AddArticle("Article Title", "thumb_media_id", "Article content...").
//	         AddArticleWithDetails("Title 2", "thumb_media_id_2", "Content 2", "Author", "Digest", "https://example.com", 1).
//	         ToUser("user1|user2").
//	         Build()
//
// mpnews消息支持1-8篇文章。这遵循与其他provider相同的构建器风格模式以保持一致性
type MPNewsBuilder struct {
	articles               []*MPNewsArticle
	toUser                 string
	toParty                string
	toTag                  string
	safe                   int
	enableIDTrans          int
	enableDuplicateCheck   int
	duplicateCheckInterval int
}

// MPNews 创建新的MPNewsBuilder实例
func MPNews() *MPNewsBuilder {
	return &MPNewsBuilder{}
}

// AddArticle 向mpnews消息添加一篇基础文章
// 返回构建器以支持链式调用。Build()中强制最大8篇文章
//
// 参数:
//   - title: 文章标题（必需，最大128字节）
//   - thumbMediaID: 缩略图的媒体ID（必需）
//   - content: 文章的HTML内容（必需）
func (b *MPNewsBuilder) AddArticle(title, thumbMediaID, content string) *MPNewsBuilder {
	b.articles = append(b.articles, &MPNewsArticle{
		Title:        title,
		ThumbMediaID: thumbMediaID,
		Content:      content,
	})
	return b
}

// AddArticleWithDetails 向mpnews消息添加一篇详细文章
// 返回构建器以支持链式调用。Build()中强制最大8篇文章
//
// 参数:
//   - title: 文章标题（必需，最大128字节）
//   - thumbMediaID: 缩略图的媒体ID（必需）
//   - content: 文章的HTML内容（必需）
//   - author: 作者名称（可选，最大64字节）
//   - digest: 简要描述（可选，最大512字节）
//   - contentSourceURL: 原文链接（可选，最大200字节）
//   - showCoverPic: 是否在内容中显示封面图（0：否，1：是）
func (b *MPNewsBuilder) AddArticleWithDetails(
	title, thumbMediaID, content, author, digest, contentSourceURL string,
	showCoverPic int,
) *MPNewsBuilder {
	b.articles = append(b.articles, &MPNewsArticle{
		Title:            title,
		ThumbMediaID:     thumbMediaID,
		Content:          content,
		Author:           author,
		Digest:           digest,
		ContentSourceURL: contentSourceURL,
		ShowCoverPic:     showCoverPic,
	})
	return b
}

// Articles sets the complete article slice (overwrites previous articles).
func (b *MPNewsBuilder) Articles(articles []*MPNewsArticle) *MPNewsBuilder {
	b.articles = articles
	return b
}

// ToUser sets the user IDs to send to, separated by "|". Use "@all" to send to all users.
func (b *MPNewsBuilder) ToUser(toUser string) *MPNewsBuilder {
	b.toUser = toUser
	return b
}

// ToParty sets the department IDs to send to, separated by "|".
func (b *MPNewsBuilder) ToParty(toParty string) *MPNewsBuilder {
	b.toParty = toParty
	return b
}

// ToTag sets the tag IDs to send to, separated by "|".
func (b *MPNewsBuilder) ToTag(toTag string) *MPNewsBuilder {
	b.toTag = toTag
	return b
}

// Safe sets whether to enable safe mode (0: no, 1: yes).
func (b *MPNewsBuilder) Safe(safe int) *MPNewsBuilder {
	b.safe = safe
	return b
}

// EnableIDTrans sets whether to enable ID translation (0: no, 1: yes).
func (b *MPNewsBuilder) EnableIDTrans(enable int) *MPNewsBuilder {
	b.enableIDTrans = enable
	return b
}

// EnableDuplicateCheck sets whether to enable duplicate message check (0: no, 1: yes).
func (b *MPNewsBuilder) EnableDuplicateCheck(enable int) *MPNewsBuilder {
	b.enableDuplicateCheck = enable
	return b
}

// DuplicateCheckInterval sets the duplicate check interval in seconds (max 4 hours).
func (b *MPNewsBuilder) DuplicateCheckInterval(interval int) *MPNewsBuilder {
	b.duplicateCheckInterval = interval
	return b
}

// Build assembles a ready-to-send *MPNewsMessage.
func (b *MPNewsBuilder) Build() *MPNewsMessage {
	return &MPNewsMessage{
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
			MsgType: TypeMPNews,
		},
		MPNews: MPNewsMessageContent{
			Articles: b.articles,
		},
	}
}
