package wecomapp

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// 无需长度常量 - 企业微信会在需要时自动截断

// MiniprogramNoticeContentItem 代表小程序通知中的内容项
type MiniprogramNoticeContentItem struct {
	// Key 内容项的键
	Key string `json:"key"`
	// Value 内容项的值
	Value string `json:"value"`
}

// MiniprogramNoticeEmphasisFirstItem 代表第一个强调项
type MiniprogramNoticeEmphasisFirstItem struct {
	// Key 强调项的键
	Key string `json:"key"`
	// Value 强调项的值
	Value string `json:"value"`
}

// MiniprogramNoticeEmphasisSecondItem 代表第二个强调项
type MiniprogramNoticeEmphasisSecondItem struct {
	// Key 强调项的键
	Key string `json:"key"`
	// Value 强调项的值
	Value string `json:"value"`
}

// MiniprogramNoticeMessageContent 代表企业微信应用API的小程序通知内容
type MiniprogramNoticeMessageContent struct {
	// AppID 小程序应用ID
	AppID string `json:"appid"`
	// Page 小程序页面路径
	Page string `json:"page,omitempty"`
	// Title 小程序通知标题
	Title string `json:"title"`
	// Description 小程序通知描述
	Description string `json:"description,omitempty"`
	// EmphasisFirstItem 第一个强调项
	EmphasisFirstItem *MiniprogramNoticeEmphasisFirstItem `json:"emphasis_first_item,omitempty"`
	// EmphasisSecondItem 第二个强调项
	EmphasisSecondItem *MiniprogramNoticeEmphasisSecondItem `json:"emphasis_second_item,omitempty"`
	// ContentItem 内容项数组
	ContentItem []*MiniprogramNoticeContentItem `json:"content_item,omitempty"`
}

// MiniprogramNoticeMessage 代表企业微信应用的小程序通知消息
// 更多详情请参考企业微信API文档：
// https://developer.work.weixin.qq.com/document/path/90372#%E5%B0%8F%E7%A8%8B%E5%BA%8F%E9%80%9A%E7%9F%A5%E6%B6%88%E6%81%AF
type MiniprogramNoticeMessage struct {
	BaseMessage

	MiniprogramNotice MiniprogramNoticeMessageContent `json:"miniprogram_notice"`
}

// NewMiniprogramNoticeMessage 创建新的MiniprogramNoticeMessage
func NewMiniprogramNoticeMessage(appID, title string) *MiniprogramNoticeMessage {
	return &MiniprogramNoticeMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeMiniprogramNotice,
		},
		MiniprogramNotice: MiniprogramNoticeMessageContent{
			AppID: appID,
			Title: title,
		},
	}
}

// Validate 验证MiniprogramNoticeMessage以确保满足企业微信API要求
func (m *MiniprogramNoticeMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	if m.MiniprogramNotice.AppID == "" {
		return core.NewParamError("miniprogram_notice appid cannot be empty")
	}

	if m.MiniprogramNotice.Title == "" {
		return core.NewParamError("miniprogram_notice title cannot be empty")
	}

	// Validate required content items
	for i, item := range m.MiniprogramNotice.ContentItem {
		if item.Key == "" {
			return core.NewParamError(fmt.Sprintf("miniprogram_notice content_item %d key cannot be empty", i+1))
		}
	}

	return nil
}
