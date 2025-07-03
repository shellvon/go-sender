package telegram

import "github.com/shellvon/go-sender/core"

// PhotoMessage represents a photo message for Telegram
// Based on SendPhotoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendphoto
type PhotoMessage struct {
	MediaMessage

	// Photo to send. Pass a file_id as String to send a photo that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get a photo from the Internet, or upload a new photo using multipart/form-data.
	// The photo must be at most 10 MB in size. The photo's width and height must not exceed maxPhotoDimensionSum in total.
	// Width and height ratio must be at most 20.
	Photo string `json:"photo"`

	// Pass True if the photo needs to be covered with a spoiler animation
	HasSpoiler bool `json:"has_spoiler,omitempty"`
}

// NewPhotoMessage creates a new PhotoMessage instance.
// Based on SendPhotoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendphoto
//   - Only chat_id and photo are required.
func NewPhotoMessage(chatID string, photo string) *PhotoMessage {
	return Photo().Chat(chatID).File(photo).Build()
}

// Validate checks if the message is valid.
func (m *PhotoMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Photo == "" {
		return core.NewParamError("photo cannot be empty")
	}
	return nil
}
