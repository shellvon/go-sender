package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// PostMessage represents a post (rich text) message for Lark/Feishu.
type PostMessage struct {
	BaseMessage

	Content PostContent `json:"content"`
}

// PostContent represents the content of a post message.
type PostContent struct {
	Post Post `json:"post"`
}

// Post represents the post structure.
type Post struct {
	ZhCN *PostLang `json:"zh_cn,omitempty"`
	EnUS *PostLang `json:"en_us,omitempty"`
}

// PostLang represents post content in a specific language.
type PostLang struct {
	Title   string          `json:"title,omitempty"`
	Content [][]PostElement `json:"content"`
}

// PostElement represents a post element (text, link, image, etc.)
type PostElement struct {
	Tag      string                 `json:"tag"`
	Text     string                 `json:"text,omitempty"`
	Href     string                 `json:"href,omitempty"`
	UserID   string                 `json:"user_id,omitempty"`
	UserName string                 `json:"user_name,omitempty"`
	ImageKey string                 `json:"image_key,omitempty"`
	Width    int                    `json:"width,omitempty"`
	Height   int                    `json:"height,omitempty"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
}

// NewPostMessage creates a new post message.
func NewPostMessage() *PostMessage {
	return &PostMessage{
		BaseMessage: BaseMessage{
			MsgType: TypePost,
		},
		Content: PostContent{
			Post: Post{},
		},
	}
}

// SetChineseContent sets the Chinese content.
func (m *PostMessage) SetChineseContent(title string, content [][]PostElement) *PostMessage {
	m.Content.Post.ZhCN = &PostLang{
		Title:   title,
		Content: content,
	}
	return m
}

// SetEnglishContent sets the English content.
func (m *PostMessage) SetEnglishContent(title string, content [][]PostElement) *PostMessage {
	m.Content.Post.EnUS = &PostLang{
		Title:   title,
		Content: content,
	}
	return m
}

// GetMsgType returns the message type.
func (m *PostMessage) GetMsgType() MessageType {
	return TypePost
}

// ProviderType returns the provider type.
func (m *PostMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the post message.
func (m *PostMessage) Validate() error {
	if m.Content.Post.ZhCN == nil && m.Content.Post.EnUS == nil {
		return errors.New("at least one language content must be provided")
	}

	if m.Content.Post.ZhCN != nil && len(m.Content.Post.ZhCN.Content) == 0 {
		return errors.New("chinese content cannot be empty")
	}

	if m.Content.Post.EnUS != nil && len(m.Content.Post.EnUS.Content) == 0 {
		return errors.New("english content cannot be empty")
	}

	return nil
}
