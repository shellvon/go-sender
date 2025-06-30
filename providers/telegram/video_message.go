package telegram

import "github.com/shellvon/go-sender/core"

// VideoMessage represents a video message for Telegram
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
type VideoMessage struct {
	MediaMessage

	// Video to send. Pass a file_id as String to send a video that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get a video from the Internet, or upload a new video using multipart/form-data.
	// The video must be at most maxVideoSizeMB MB in size.
	Video string `json:"video"`

	// Duration of sent video in seconds
	Duration int `json:"duration,omitempty"`

	// Video width
	Width int `json:"width,omitempty"`

	// Video height
	Height int `json:"height,omitempty"`

	// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
	// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	Thumbnail string `json:"thumbnail,omitempty"`

	// Cover for the video in the message. Pass a file_id to send a file that exists on the Telegram servers (recommended),
	// pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Cover string `json:"cover,omitempty"`

	// Start timestamp for the video in the message
	StartTimestamp int `json:"start_timestamp,omitempty"`

	// Pass True if the video needs to be covered with a spoiler animation
	HasSpoiler bool `json:"has_spoiler,omitempty"`

	// Pass True if the uploaded video is suitable for streaming
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
}

// NewVideoMessage creates a new VideoMessage instance.
func NewVideoMessage(chatID string, video string, opts ...interface{}) *VideoMessage {
	return NewVideoMessageWithBuilder(chatID, video, opts...)
}

func (m *VideoMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *VideoMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *VideoMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Video == "" {
		return core.NewParamError("video cannot be empty")
	}
	return nil
}

type VideoMessageOption func(*VideoMessage)

// WithVideoDuration sets the duration of the video in seconds
// This is optional and can be used to provide metadata about the video.
func WithVideoDuration(duration int) VideoMessageOption {
	return func(m *VideoMessage) { m.Duration = duration }
}

// WithVideoWidth sets the width of the video
// This is optional and can be used to provide metadata about the video.
func WithVideoWidth(width int) VideoMessageOption {
	return func(m *VideoMessage) { m.Width = width }
}

// WithVideoHeight sets the height of the video
// This is optional and can be used to provide metadata about the video.
func WithVideoHeight(height int) VideoMessageOption {
	return func(m *VideoMessage) { m.Height = height }
}

// WithVideoThumbnail sets the thumbnail for the video
// Should be in JPEG format and less than 200 kB in size
// Width and height should not exceed 320.
func WithVideoThumbnail(thumbnail string) VideoMessageOption {
	return func(m *VideoMessage) { m.Thumbnail = thumbnail }
}

// WithVideoCover sets the video cover/thumbnail
// Can only be uploaded as a separate file
// Only jpeg and png formats are supported.
func WithVideoCover(cover string) VideoMessageOption {
	return func(m *VideoMessage) { m.Cover = cover }
}

// WithVideoStartTimestamp sets the timestamp of the video start
// This can be used to specify where the video should start playing.
func WithVideoStartTimestamp(startTimestamp int) VideoMessageOption {
	return func(m *VideoMessage) { m.StartTimestamp = startTimestamp }
}

// WithVideoHasSpoiler sets whether the video should be covered with a spoiler animation
// Users will need to tap to reveal the video content.
func WithVideoHasSpoiler(has bool) VideoMessageOption {
	return func(m *VideoMessage) { m.HasSpoiler = has }
}

// WithVideoSupportsStreaming sets whether the uploaded video is suitable for streaming
// This helps Telegram optimize the video for streaming playback.
func WithVideoSupportsStreaming(supports bool) VideoMessageOption {
	return func(m *VideoMessage) { m.SupportsStreaming = supports }
}
