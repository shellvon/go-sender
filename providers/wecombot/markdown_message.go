package wecombot

import (
	"encoding/json"

	"github.com/shellvon/go-sender/core"
)

const maxMarkdownContentLength = 4096

// Markdown versions.
const (
	MarkdownVersionLegacy = "legacy"
	MarkdownVersionV2     = "v2"
)

// MarkdownContent represents the markdown content for a WeCom message.
type MarkdownContent struct {
	// Content of the markdown message. Maximum length is 4096 bytes, and it must be UTF-8 encoded.
	Content string `json:"content"`
	// Version of the markdown message.
	// Currently, v2 or legacy is supported.
	Version string `json:"version"`
}

// MarkdownMessage represents a markdown message for WeCom.
// For more details, refer to the WeCom API documentation:
// https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
type MarkdownMessage struct {
	BaseMessage

	Markdown MarkdownContent `json:"markdown"`
}

// NewMarkdownMessage creates a new MarkdownMessage.
// Based on SendMarkdownParams from WeCom Bot API
// https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
//   - Only content and version are required.
//   - version is "legacy" if not provided or empty.
//   - version is "v2" if provided version is "v2".
func NewMarkdownMessage(content string, version string) *MarkdownMessage {
	return Markdown().Content(content).Version(version).Build()
}

// Validate validates the MarkdownMessage to ensure it meets WeCom API requirements.
func (m *MarkdownMessage) Validate() error {
	if m.Markdown.Content == "" {
		return core.NewParamError("markdown content cannot be empty")
	}
	if len([]rune(m.Markdown.Content)) > maxMarkdownContentLength {
		return core.NewParamError("markdown content exceeds 4096 characters")
	}
	return nil
}

// MarshalJSON implements custom JSON marshalling to accommodate the v2 API which
// expects the field name "markdown_v2" instead of "markdown" while keeping the
// payload structure identical.
func (m *MarkdownMessage) MarshalJSON() ([]byte, error) {
	// Build map for final JSON.
	data := map[string]interface{}{
		"msgtype": m.MsgType,
	}

	key := "markdown"
	if m.Markdown.Version == MarkdownVersionV2 {
		key = "markdown_v2"
	}
	data[key] = m.Markdown

	return json.Marshal(data)
}
