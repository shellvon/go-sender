//nolint:dupl // VoiceMessage and FileMessage share similar boilerplate by design.
package wecombot

import (
	"os"

	"github.com/shellvon/go-sender/core"
)

// VoicePayload holds media_id field as required by WeCom.
type VoicePayload struct {
	MediaID string `json:"media_id"`
}

// VoiceMessage represents a voice message for WeCom bot.
// API: https://developer.work.weixin.qq.com/document/path/91770#语音类型
// Only media_id is needed.
// localPath is internal helper for auto-upload.
type VoiceMessage struct {
	BaseMessage

	Voice VoicePayload `json:"voice"`

	localPath string `json:"-"`
}

// Validate ensures media_id is present (or localPath to be uploaded).
func (m *VoiceMessage) Validate() error {
	if m.Voice.MediaID == "" && m.localPath == "" {
		return core.NewParamError("voice message requires media_id or local file path")
	}
	// optional: check file exists
	if m.Voice.MediaID == "" && m.localPath != "" {
		if _, err := os.Stat(m.localPath); err != nil {
			return core.NewParamError("voice local file not found")
		}
	}
	return nil
}

//--------------------------------------------------------------------
// Builder
//--------------------------------------------------------------------

type VoiceBuilder struct {
	mediaID   string
	localPath string
}

// Voice returns a new VoiceBuilder.
func Voice() *VoiceBuilder {
	return &VoiceBuilder{}
}

// MediaID sets an existing media_id.
func (b *VoiceBuilder) MediaID(id string) *VoiceBuilder {
	b.mediaID = id
	return b
}

// LocalPath sets a local file path to be uploaded automatically.
func (b *VoiceBuilder) LocalPath(path string) *VoiceBuilder {
	b.localPath = path
	return b
}

// Build constructs *VoiceMessage.
func (b *VoiceBuilder) Build() *VoiceMessage {
	return &VoiceMessage{
		BaseMessage: BaseMessage{MsgType: TypeVoice},
		Voice:       VoicePayload{MediaID: b.mediaID},
		localPath:   b.localPath,
	}
}

// uploadTarget implementation.
func (m *VoiceMessage) getLocalPath() string { return m.localPath }
func (m *VoiceMessage) getMediaID() string   { return m.Voice.MediaID }
func (m *VoiceMessage) setMediaID(id string) { m.Voice.MediaID = id }
func (m *VoiceMessage) mediaType() string    { return "voice" }
