package wecomapp

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// 无需长度常量 - 企业微信会在需要时自动截断

// NewsArticle 代表图文消息中的单篇文章
type NewsArticle struct {
	// Title 文章标题
	Title string `json:"title"`
	// Description 文章描述
	Description string `json:"description"`
	// URL 点击文章时跳转的URL
	URL string `json:"url"`
	// PicURL 文章的图片URL。支持JPG和PNG格式
	PicURL string `json:"picurl"`

	// AppID 小程序appid，必须是与当前应用关联的小程序，appid和pagepath必须同时填写，填写后会忽略url字段
	AppID string `json:"appid"`

	// PagePath 点击消息卡片后的小程序页面，最长128字节，仅限本小程序内的页面。appid和pagepath必须同时填写，填写后会忽略url字段
	PagePath string `json:"pagepath"`
}

// NewsMessageContent 代表企业微信应用API的图文内容
type NewsMessageContent struct {
	Articles []*NewsArticle `json:"articles"`
}

// NewsMessage 代表企业微信应用的图文消息
// 更多详情请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/90372#%E5%9B%BE%E6%96%87%E6%B6%88%E6%81%AF
type NewsMessage struct {
	BaseMessage

	News NewsMessageContent `json:"news"`
}

// NewNewsMessage 创建新的NewsMessage
func NewNewsMessage(articles []*NewsArticle) *NewsMessage {
	return &NewsMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeNews,
		},
		News: NewsMessageContent{
			Articles: articles,
		},
	}
}

// Validate 验证NewsMessage以确保满足企业微信API要求
func (m *NewsMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	if len(m.News.Articles) == 0 {
		return core.NewParamError("articles cannot be empty")
	}

	// 无需限制验证 - 企业微信会处理限制

	for i, article := range m.News.Articles {
		if article.Title == "" {
			return core.NewParamError(fmt.Sprintf("article %d: title is required", i+1))
		}

		// 必须提供URL或小程序（AppID + PagePath）之一
		hasURL := article.URL != ""
		hasMiniprogram := article.AppID != "" && article.PagePath != ""

		if !hasURL && !hasMiniprogram {
			return core.NewParamError(
				fmt.Sprintf("article %d: either URL or mini-program (AppID + PagePath) must be provided", i+1),
			)
		}

		// 如果两者都提供，小程序优先（URL将被忽略）
	}

	return nil
}
