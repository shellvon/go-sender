//nolint:dupl // FileMessage and VoiceMessage share similar boilerplate by design.
package wecombot

import (
	"os"

	"github.com/shellvon/go-sender/core"
)

// FilePayload holds media_id.
type FilePayload struct {
	MediaID string `json:"media_id"`
}

// FileMessage represents a file message for WeCom bot.
// https://developer.work.weixin.qq.com/document/path/91770#文件类型
// Only media_id required.
type FileMessage struct {
	BaseMessage

	File FilePayload `json:"file"`

	localPath string `json:"-"`
}

func (m *FileMessage) Validate() error {
	if m.File.MediaID == "" && m.localPath == "" {
		return core.NewParamError("file message requires media_id or local file path")
	}
	if m.File.MediaID == "" && m.localPath != "" {
		if _, err := os.Stat(m.localPath); err != nil {
			return core.NewParamError("file local file not found")
		}
	}
	return nil
}

// Builder

type FileBuilder struct {
	mediaID   string
	localPath string
}

// File returns new FileBuilder.
func File() *FileBuilder { return &FileBuilder{} }

func (b *FileBuilder) MediaID(id string) *FileBuilder { b.mediaID = id; return b }

func (b *FileBuilder) LocalPath(path string) *FileBuilder { b.localPath = path; return b }

func (b *FileBuilder) Build() *FileMessage {
	return &FileMessage{
		BaseMessage: BaseMessage{MsgType: TypeFile},
		File:        FilePayload{MediaID: b.mediaID},
		localPath:   b.localPath,
	}
}

// uploadTarget implementation for file.
func (m *FileMessage) getLocalPath() string { return m.localPath }
func (m *FileMessage) getMediaID() string   { return m.File.MediaID }
func (m *FileMessage) setMediaID(id string) { m.File.MediaID = id }
func (m *FileMessage) mediaType() string    { return "file" }
