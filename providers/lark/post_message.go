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

// Compile-time assertion: PostMessage implements Message interface.
var (
	_ core.Validatable = (*PostMessage)(nil)
)

// PostContent represents the content of a post message.
type PostContent struct {
	Post PostPayload `json:"post"`
}

// PostPayload represents the post structure.
type PostPayload struct {
	ZhCN *PostLang `json:"zh_cn,omitempty"`
	EnUS *PostLang `json:"en_us,omitempty"`
}

// PostLang represents post content in a specific language.
type PostLang struct {
	// Title is the title of the post.
	Title   string          `json:"title,omitempty"`
	Content [][]PostElement `json:"content"`
}

type PostElement struct {
	Tag  string `json:"tag"`
	Text string `json:"text,omitempty"`
	// When tag is text, unescape is used to unescape the text. default is false.
	UnEscape bool `json:"un_escape,omitempty"`
	// When tag is a, href is used to set the href of the link.
	Href string `json:"href,omitempty"`
	// When tag is at, user_id is used to set the user id of the mention.
	UserID string `json:"user_id,omitempty"`
	// When tag is at, user_name is used to set the user name of the mention.
	UserName string `json:"user_name,omitempty"`
	// When tag is img, image_key is used to set the image key of the image.
	// See https://open.feishu.cn/document/server-docs/im-v1/image/create
	ImageKey string `json:"image_key,omitempty"`
}

// PostBuilder provides a fluent API to construct Lark post (rich text) messages (unexported).
type PostBuilder struct {
	zhCN *PostLang
	enUS *PostLang
}

// Post creates a new postBuilder instance (user-facing API).
func Post() *PostBuilder { return &PostBuilder{} }

// ZhCN sets the Chinese content.
func (b *PostBuilder) ZhCN(title string, content [][]PostElement) *PostBuilder {
	b.zhCN = &PostLang{Title: title, Content: content}
	return b
}

// EnUS sets the English content.
func (b *PostBuilder) EnUS(title string, content [][]PostElement) *PostBuilder {
	b.enUS = &PostLang{Title: title, Content: content}
	return b
}

// Build assembles a *PostMessage.
func (b *PostBuilder) Build() *PostMessage {
	return &PostMessage{
		BaseMessage: newBaseMessage(TypePost),
		Content: PostContent{
			Post: PostPayload{
				ZhCN: b.zhCN,
				EnUS: b.enUS,
			},
		},
	}
}

// NewPostMessage creates a new post message.
func NewPostMessage() *PostMessage {
	return Post().Build()
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
