package telegram

import "github.com/shellvon/go-sender/core"

// AnimationMessage represents an animation message for Telegram
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
type AnimationMessage struct {
	MediaMessage

	// Animation to send. Pass a file_id as String to send an animation that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get an animation from the Internet, or upload a new animation using multipart/form-data.
	// The animation must be at most maxAnimationSizeMB MB in size.
	Animation string `json:"animation"`

	// Duration of sent animation in seconds
	Duration int `json:"duration,omitempty"`

	// Animation width
	Width int `json:"width,omitempty"`

	// Animation height
	Height int `json:"height,omitempty"`

	// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
	// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	Thumbnail string `json:"thumbnail,omitempty"`

	// Pass True if the animation needs to be covered with a spoiler animation
	HasSpoiler bool `json:"has_spoiler,omitempty"`
}

// NewAnimationMessage creates a new AnimationMessage instance.
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
//   - Only chat_id and animation are required.
//   - Currently, only file_id or http URL is supported.
func NewAnimationMessage(chatID string, animation string) *AnimationMessage {
	return Animation().Chat(chatID).File(animation).Build()
}

func (m *AnimationMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Animation == "" {
		return core.NewParamError("animation cannot be empty")
	}
	return nil
}
