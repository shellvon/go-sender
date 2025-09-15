package wecombot

import (
	"encoding/json"

	"github.com/shellvon/go-sender/core"
)

const maxMarkdownContentLength = 4096

// MarkdownVersion 表示 Markdown 格式的版本。
type MarkdownVersion string

const (
	// MarkdownVersionLegacy 是 Markdown 格式的旧版本。
	MarkdownVersionLegacy MarkdownVersion = "legacy"
	// MarkdownVersionV2 是 Markdown 格式的新版本。
	// 建议使用最新的客户端版本以体验该消息。
	// 仅在客户端版本 4.1.36 或更高版本（Android 4.1.38 或更高版本）上可用，低版本将被视为纯文本。
	MarkdownVersionV2 MarkdownVersion = "v2"
)

// String 实现 Stringer 接口，将 MarkdownVersion 转换为字符串。
// 返回值：string - MarkdownVersion 的字符串表示。
func (v MarkdownVersion) String() string {
	return string(v)
}

// IsValid 检查 Markdown 版本是否有效。
// 返回值：bool - 如果版本是 MarkdownVersionLegacy 或 MarkdownVersionV2，则返回 true；否则返回 false。
func (v MarkdownVersion) IsValid() bool {
	switch v {
	case MarkdownVersionLegacy, MarkdownVersionV2:
		return true
	default:
		return false
	}
}

// MarkdownContent 表示企业微信消息的 Markdown 内容。
type MarkdownContent struct {
	// Markdown 消息的内容。最大长度为 4096 字节，且必须是 UTF-8 编码。
	Content string `json:"content"`
	// Markdown 消息的版本。
	// 当前支持 v2 或 legacy。
	Version MarkdownVersion `json:"version,omitempty"`
}

// MarkdownMessage 表示企业微信的 Markdown 消息。
// 更多详情，请参考企业微信 API 文档：
// https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
type MarkdownMessage struct {
	BaseMessage

	Markdown MarkdownContent `json:"markdown"`
}

// NewMarkdownMessage 创建一个新的 MarkdownMessage 实例。
// 基于企业微信机器人 API 的 SendMarkdownParams
// 参考：https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
//   - 仅 content 和 version 是必需的。
//   - 如果未提供或版本为空，则版本为 "legacy"。
//   - 如果提供的版本为 "v2"，则版本为 "v2"。
//
// 更多详情请参见 https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
//
// 参数：
//   - content string - Markdown 消息内容；
//   - version MarkdownVersion - Markdown 版本。
//
// 返回值：*MarkdownMessage - 新创建的 Markdown 消息实例。
func NewMarkdownMessage(content string, version MarkdownVersion) *MarkdownMessage {
	return Markdown().Content(content).Version(version).Build()
}

// Validate 验证 MarkdownMessage 是否满足企业微信 API 的要求。
// 该方法检查内容是否为空、内容长度是否超过 4096 字符，以及版本是否有效。
// 返回值：error - 如果验证失败，返回具体的参数错误；否则返回 nil。
func (m *MarkdownMessage) Validate() error {
	if m.Markdown.Content == "" {
		return core.NewParamError("Markdown 内容不能为空")
	}
	if len([]rune(m.Markdown.Content)) > maxMarkdownContentLength {
		return core.NewParamError("Markdown 内容超过 4096 个字符")
	}
	if m.Markdown.Version != "" && !m.Markdown.Version.IsValid() {
		return core.NewParamError("无效的 Markdown 版本：" + string(m.Markdown.Version))
	}
	return nil
}

// MarshalJSON 实现自定义 JSON 序列化，以适应 v2 API 的要求。
// v2 API 期望字段名为 "markdown_v2" 而不是 "markdown"，但负载结构保持一致。
// 返回值：[]byte - 序列化后的 JSON 数据；error - 如果序列化失败，返回错误。
func (m *MarkdownMessage) MarshalJSON() ([]byte, error) {
	// 构建最终 JSON 的映射。
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
