package wecombot

import (
	"os"

	"github.com/shellvon/go-sender/core"
)

// FilePayload 包含 media_id。
type FilePayload struct {
	MediaID string `json:"media_id"`
}

// FileMessage 表示企业微信机器人的文件消息。
// 参考：https://developer.work.weixin.qq.com/document/path/91770#文件类型
// 仅需要 media_id。
type FileMessage struct {
	BaseMessage

	File FilePayload `json:"file"`

	localPath string `json:"-"`
}

func (m *FileMessage) Validate() error {
	if m.File.MediaID == "" && m.localPath == "" {
		return core.NewParamError("文件消息需要 media_id 或本地文件路径")
	}
	if m.File.MediaID == "" && m.localPath != "" {
		if _, err := os.Stat(m.localPath); err != nil {
			return core.NewParamError("本地文件未找到")
		}
	}
	return nil
}

// uploadTarget 实现文件相关方法。

// getLocalPath 获取文件消息的本地文件路径。
// 返回值：string - 文件的本地路径。
func (m *FileMessage) getLocalPath() string { return m.localPath }

// getMediaID 获取文件消息的 media_id。
// 返回值：string - 文件的 media_id。
func (m *FileMessage) getMediaID() string { return m.File.MediaID }

// setMediaID 设置文件消息的 media_id。
// 参数：id string - 要设置的 media_id。
func (m *FileMessage) setMediaID(id string) { m.File.MediaID = id }

// mediaType 返回文件消息的媒体类型。
// 返回值：string - 媒体类型，固定为 "file"。
func (m *FileMessage) mediaType() string { return "file" }
