package wecomapp

import (
	"github.com/shellvon/go-sender/core"
)

// MediaMessageContent 代表媒体消息的通用结构.
type MediaMessageContent struct {
	// MediaID 从媒体上传API获得的媒体文件ID
	MediaID string `json:"media_id"`
}

// VideoMediaMessageContent 使用视频特定字段扩展MediaMessageContent.
type VideoMediaMessageContent struct {
	MediaMessageContent

	// Title 视频标题（可选）
	Title string `json:"title,omitempty"`
	// Description 视频描述（可选）
	Description string `json:"description,omitempty"`
}

// BaseMediaMessage 为所有媒体消息提供通用功能.
type BaseMediaMessage struct {
	BaseMessage

	// LocalPath 用于自动媒体上传（不发送到API）
	LocalPath string `json:"-"`
}

// Validate validates the BaseMediaMessage to ensure it meets WeChat Work API requirements.
func (m *BaseMediaMessage) Validate() error {
	if err := m.BaseMessage.Validate(); err != nil {
		return err
	}

	if m.AgentID == "" {
		return core.NewParamError("agentid is required")
	}

	// At least one target must be specified
	if m.ToUser == "" && m.ToParty == "" && m.ToTag == "" {
		return core.NewParamError("at least one of touser, toparty, or totag must be specified")
	}

	return nil
}

// Common uploadTarget interface implementation helpers.
func (m *BaseMediaMessage) getLocalPath() string { return m.LocalPath }

// MediaMessage interface for type-specific operations.
type MediaMessage interface {
	Message
	uploadTarget
	getMediaContent() *MediaMessageContent
}

// ImageMessage 代表企业微信应用的图片消息.
type ImageMessage struct {
	BaseMediaMessage

	Image MediaMessageContent `json:"image"`
}

func NewImageMessage(mediaID string) *ImageMessage {
	return &ImageMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeImage},
		},
		Image: MediaMessageContent{MediaID: mediaID},
	}
}

func NewImageMessageFromPath(localPath string) *ImageMessage {
	return &ImageMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeImage},
			LocalPath:   localPath,
		},
	}
}

func (m *ImageMessage) getMediaContent() *MediaMessageContent { return &m.Image }
func (m *ImageMessage) getMediaID() string                    { return m.Image.MediaID }
func (m *ImageMessage) setMediaID(id string)                  { m.Image.MediaID = id }
func (m *ImageMessage) mediaType() string                     { return string(MediaTypeImage) }

func (m *ImageMessage) Validate() error {
	if err := m.BaseMediaMessage.Validate(); err != nil {
		return err
	}
	if m.Image.MediaID == "" && m.LocalPath == "" {
		return core.NewParamError("either media_id or local file path must be provided")
	}
	return nil
}

// VoiceMessage 代表企业微信应用的语音消息.
type VoiceMessage struct {
	BaseMediaMessage

	Voice MediaMessageContent `json:"voice"`
}

func NewVoiceMessage(mediaID string) *VoiceMessage {
	return &VoiceMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeVoice},
		},
		Voice: MediaMessageContent{MediaID: mediaID},
	}
}

func NewVoiceMessageFromPath(localPath string) *VoiceMessage {
	return &VoiceMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeVoice},
			LocalPath:   localPath,
		},
	}
}

func (m *VoiceMessage) getMediaContent() *MediaMessageContent { return &m.Voice }
func (m *VoiceMessage) getMediaID() string                    { return m.Voice.MediaID }
func (m *VoiceMessage) setMediaID(id string)                  { m.Voice.MediaID = id }
func (m *VoiceMessage) mediaType() string                     { return string(MediaTypeVoice) }

func (m *VoiceMessage) Validate() error {
	if err := m.BaseMediaMessage.Validate(); err != nil {
		return err
	}
	if m.Voice.MediaID == "" && m.LocalPath == "" {
		return core.NewParamError("either media_id or local file path must be provided")
	}
	return nil
}

// VideoMessage 代表企业微信应用的视频消息.
type VideoMessage struct {
	BaseMediaMessage

	Video VideoMediaMessageContent `json:"video"`
}

func NewVideoMessage(mediaID string) *VideoMessage {
	return &VideoMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeVideo},
		},
		Video: VideoMediaMessageContent{
			MediaMessageContent: MediaMessageContent{MediaID: mediaID},
		},
	}
}

func NewVideoMessageFromPath(localPath string) *VideoMessage {
	return &VideoMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeVideo},
			LocalPath:   localPath,
		},
	}
}

func (m *VideoMessage) getMediaContent() *MediaMessageContent { return &m.Video.MediaMessageContent }
func (m *VideoMessage) getMediaID() string                    { return m.Video.MediaID }
func (m *VideoMessage) setMediaID(id string)                  { m.Video.MediaID = id }
func (m *VideoMessage) mediaType() string                     { return string(MediaTypeVideo) }

func (m *VideoMessage) Validate() error {
	if err := m.BaseMediaMessage.Validate(); err != nil {
		return err
	}
	if m.Video.MediaID == "" && m.LocalPath == "" {
		return core.NewParamError("either media_id or local file path must be provided")
	}
	return nil
}

// FileMessage 代表企业微信应用的文件消息.
type FileMessage struct {
	BaseMediaMessage

	File MediaMessageContent `json:"file"`
}

func NewFileMessage(mediaID string) *FileMessage {
	return &FileMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeFile},
		},
		File: MediaMessageContent{MediaID: mediaID},
	}
}

func NewFileMessageFromPath(localPath string) *FileMessage {
	return &FileMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{MsgType: TypeFile},
			LocalPath:   localPath,
		},
	}
}

func (m *FileMessage) getMediaContent() *MediaMessageContent { return &m.File }
func (m *FileMessage) getMediaID() string                    { return m.File.MediaID }
func (m *FileMessage) setMediaID(id string)                  { m.File.MediaID = id }
func (m *FileMessage) mediaType() string                     { return string(MediaTypeFile) }

func (m *FileMessage) Validate() error {
	if err := m.BaseMediaMessage.Validate(); err != nil {
		return err
	}
	if m.File.MediaID == "" && m.LocalPath == "" {
		return core.NewParamError("either media_id or local file path must be provided")
	}
	return nil
}
