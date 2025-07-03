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
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
//   - Only chat_id and video are required.
//   - Currently, only file_id or http URL is supported.
func NewVideoMessage(chatID string, video string) *VideoMessage {
	return Video().Chat(chatID).File(video).Build()
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
