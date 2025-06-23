package telegram

import "github.com/shellvon/go-sender/core"

// AnimationMessage represents an animation message for Telegram
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
type AnimationMessage struct {
	MediaMessage

	// Animation to send. Pass a file_id as String to send an animation that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get an animation from the Internet, or upload a new animation using multipart/form-data.
	// The animation must be at most 50 MB in size.
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

	// Animation caption (may also be used when resending animation by file_id), 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`

	// Mode for parsing entities in the animation caption. See formatting options for more details on supported modes.
	// Options: "HTML", "Markdown", "MarkdownV2"
	ParseMode string `json:"parse_mode,omitempty"`

	// A JSON-serialized list of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	// Pass True if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`

	// Pass True if the animation needs to be covered with a spoiler animation
	HasSpoiler bool `json:"has_spoiler,omitempty"`
}

func (m *AnimationMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *AnimationMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
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

type AnimationMessageOption func(*AnimationMessage)

// WithAnimationDuration sets the duration of the animation in seconds
// This is optional and can be used to provide metadata about the animation
func WithAnimationDuration(duration int) AnimationMessageOption {
	return func(m *AnimationMessage) { m.Duration = duration }
}

// WithAnimationWidth sets the width of the animation
// This is optional and can be used to provide metadata about the animation
func WithAnimationWidth(width int) AnimationMessageOption {
	return func(m *AnimationMessage) { m.Width = width }
}

// WithAnimationHeight sets the height of the animation
// This is optional and can be used to provide metadata about the animation
func WithAnimationHeight(height int) AnimationMessageOption {
	return func(m *AnimationMessage) { m.Height = height }
}

// WithAnimationThumbnail sets the thumbnail for the animation
// Should be in JPEG format and less than 200 kB in size
// Width and height should not exceed 320
func WithAnimationThumbnail(thumbnail string) AnimationMessageOption {
	return func(m *AnimationMessage) { m.Thumbnail = thumbnail }
}

// WithAnimationHasSpoiler sets whether the animation should be covered with a spoiler animation
// Users will need to tap to reveal the animation content
func WithAnimationHasSpoiler(has bool) AnimationMessageOption {
	return func(m *AnimationMessage) { m.HasSpoiler = has }
}

func NewAnimationMessage(chatID string, animation string, opts ...interface{}) *AnimationMessage {
	msg := &AnimationMessage{
		MediaMessage: MediaMessage{
			BaseMessage: BaseMessage{
				MsgType: TypeAnimation,
				ChatID:  chatID,
			},
		},
		Animation: animation,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case AnimationMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}
