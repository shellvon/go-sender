package telegram

import "github.com/shellvon/go-sender/core"

// VideoNoteMessage represents a video note message for Telegram
// Based on SendVideoNoteParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideonote
type VideoNoteMessage struct {
	BaseMessage

	// Video note to send. Pass a file_id as String to send a video note that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get a video note from the Internet, or upload a new video note using multipart/form-data.
	// The video note must be at most 50 MB in size.
	VideoNote string `json:"video_note"`

	// Duration of sent video note in seconds
	Duration int `json:"duration,omitempty"`

	// Video width and height, i.e. diameter of the video message
	Length int `json:"length,omitempty"`

	// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
	// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	Thumbnail string `json:"thumbnail,omitempty"`
}

func (m *VideoNoteMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *VideoNoteMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *VideoNoteMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.VideoNote == "" {
		return core.NewParamError("video_note cannot be empty")
	}
	return nil
}

type VideoNoteMessageOption func(*VideoNoteMessage)

// WithVideoNoteDuration sets the duration of the video note in seconds
// This is optional and can be used to provide metadata about the video note
func WithVideoNoteDuration(duration int) VideoNoteMessageOption {
	return func(m *VideoNoteMessage) { m.Duration = duration }
}

// WithVideoNoteLength sets the length (diameter) of the video note
// Video notes are always square, so this represents both width and height
func WithVideoNoteLength(length int) VideoNoteMessageOption {
	return func(m *VideoNoteMessage) { m.Length = length }
}

// WithVideoNoteThumbnail sets the thumbnail for the video note
// Should be in JPEG format and less than 200 kB in size
// Width and height should not exceed 320
func WithVideoNoteThumbnail(thumbnail string) VideoNoteMessageOption {
	return func(m *VideoNoteMessage) { m.Thumbnail = thumbnail }
}

func NewVideoNoteMessage(chatID string, videoNote string, opts ...interface{}) *VideoNoteMessage {
	msg := &VideoNoteMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeVideoNote,
			ChatID:  chatID,
		},
		VideoNote: videoNote,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case VideoNoteMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}
