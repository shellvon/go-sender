package telegram

import "github.com/shellvon/go-sender/core"

// AudioMessage represents an audio message for Telegram
// Based on SendAudioParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendaudio
type AudioMessage struct {
	MediaMessage

	// Audio file to send. Pass a file_id as String to send an audio file that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get an audio file from the Internet, or upload a new one using multipart/form-data.
	// The audio must be at most 50 MB in size.
	Audio string `json:"audio"`

	// Duration of the audio in seconds
	Duration int `json:"duration,omitempty"`

	// Performer of the audio
	Performer string `json:"performer,omitempty"`

	// Track name
	Title string `json:"title,omitempty"`

	// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
	// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	Thumbnail string `json:"thumbnail,omitempty"`
}

func (m *AudioMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *AudioMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *AudioMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Audio == "" {
		return core.NewParamError("audio cannot be empty")
	}
	return nil
}

type AudioMessageOption func(*AudioMessage)

// WithAudioDuration sets the duration of the audio in seconds
// This is optional and can be used to provide metadata about the audio
func WithAudioDuration(duration int) AudioMessageOption {
	return func(m *AudioMessage) { m.Duration = duration }
}

// WithAudioPerformer sets the performer of the audio
// This is optional and can be used to provide metadata about the audio
func WithAudioPerformer(performer string) AudioMessageOption {
	return func(m *AudioMessage) { m.Performer = performer }
}

// WithAudioTitle sets the title of the audio
// This is optional and can be used to provide metadata about the audio
func WithAudioTitle(title string) AudioMessageOption {
	return func(m *AudioMessage) { m.Title = title }
}

// WithAudioThumbnail sets the thumbnail for the audio
// Should be in JPEG format and less than 200 kB in size
// Width and height should not exceed 320
func WithAudioThumbnail(thumbnail string) AudioMessageOption {
	return func(m *AudioMessage) { m.Thumbnail = thumbnail }
}

func NewAudioMessage(chatID string, audio string, opts ...interface{}) *AudioMessage {
	msg := &AudioMessage{
		MediaMessage: MediaMessage{
			BaseMessage: BaseMessage{
				MsgType: TypeAudio,
				ChatID:  chatID,
			},
		},
		Audio: audio,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case AudioMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}
