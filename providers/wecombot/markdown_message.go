package wecombot

import (
	"encoding/json"

	"github.com/shellvon/go-sender/core"
)

const maxMarkdownContentLength = 4096

// MarkdownVersion represents the version of markdown format
type MarkdownVersion string

const (
	// MarkdownVersionLegacy is the legacy version of markdown format.
	MarkdownVersionLegacy MarkdownVersion = "legacy"
	// MarkdownVersionV2 is the new version of markdown format.
	// It is recommended to use the latest client version to experience the message.
	// only available on client version 4.1.36 or higher (Android 4.1.38 or higher), lower version will be treated as plain text
	MarkdownVersionV2 MarkdownVersion = "v2"
)

// String implements the Stringer interface.
func (v MarkdownVersion) String() string {
	return string(v)
}

// IsValid checks if the markdown version is valid.
func (v MarkdownVersion) IsValid() bool {
	switch v {
	case MarkdownVersionLegacy, MarkdownVersionV2:
		return true
	default:
		return false
	}
}

// MarkdownContent represents the markdown content for a WeCom message.
type MarkdownContent struct {
	// Content of the markdown message. Maximum length is 4096 bytes, and it must be UTF-8 encoded.
	Content string `json:"content"`
	// Version of the markdown message.
	// Currently, v2 or legacy is supported.
	Version MarkdownVersion `json:"version,omitempty"`
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
//
// See https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B for more details.
func NewMarkdownMessage(content string, version MarkdownVersion) *MarkdownMessage {
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
	if m.Markdown.Version != "" && !m.Markdown.Version.IsValid() {
		return core.NewParamError("invalid markdown version: " + string(m.Markdown.Version))
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
