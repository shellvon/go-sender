package wecomapp

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// MPNewsArticle 代表mpnews消息中的单篇文章
// mpnews类型的图文消息，跟普通的图文消息一致，唯一的差异是图文内容存储在企业微信。
// 多次发送mpnews，会被认为是不同的图文，阅读、点赞的统计会被分开计算。
type MPNewsArticle struct {
	// Title 文章标题
	Title string `json:"title"`
	// ThumbMediaID 缩略图的媒体ID
	ThumbMediaID string `json:"thumb_media_id"`
	// Author 文章作者
	Author string `json:"author,omitempty"`
	// ContentSourceURL 原文链接
	ContentSourceURL string `json:"content_source_url,omitempty"`
	// Content 文章的HTML内容
	Content string `json:"content"`
	// Digest 文章摘要
	Digest string `json:"digest,omitempty"`
	// ShowCoverPic 是否在内容中显示封面图（0：否，1：是）
	ShowCoverPic int `json:"show_cover_pic,omitempty"`
}

// MPNewsMessageContent 代表企业微信应用API的mpnews内容
type MPNewsMessageContent struct {
	Articles []*MPNewsArticle `json:"articles"`
}

// MPNewsMessage 代表企业微信应用的mpnews消息
// mpnews类型的图文消息，跟普通的图文消息一致，唯一的差异是图文内容存储在企业微信。
// 多次发送mpnews，会被认为是不同的图文，阅读、点赞的统计会被分开计算。
//
// 更多详情请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/90372#%E5%9B%BE%E6%96%87%E6%B6%88%E6%81%AF%EF%BC%88mpnews%EF%BC%89
type MPNewsMessage struct {
	BaseMessage

	MPNews MPNewsMessageContent `json:"mpnews"`
}

// NewMPNewsMessage 创建新的MPNewsMessage
func NewMPNewsMessage(articles []*MPNewsArticle) *MPNewsMessage {
	return &MPNewsMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeMPNews,
		},
		MPNews: MPNewsMessageContent{
			Articles: articles,
		},
	}
}

// Validate 验证MPNewsMessage以确保满足企业微信API要求
func (m *MPNewsMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	if len(m.MPNews.Articles) == 0 {
		return core.NewParamError("mpnews articles cannot be empty")
	}

	for i, article := range m.MPNews.Articles {
		if article.Title == "" {
			return core.NewParamError(fmt.Sprintf("mpnews article %d: title is required", i+1))
		}
		if article.ThumbMediaID == "" {
			return core.NewParamError(fmt.Sprintf("mpnews article %d: thumb_media_id is required", i+1))
		}
		if article.Content == "" {
			return core.NewParamError(fmt.Sprintf("mpnews article %d: content is required", i+1))
		}
	}

	return nil
}
