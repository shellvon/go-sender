package telegram

import "github.com/shellvon/go-sender/core"

// DocumentMessage represents a document message for Telegram
// Based on SendDocumentParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddocument
type DocumentMessage struct {
	MediaMessage

	// File to send. Pass a file_id as String to send a file that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get a file from the Internet, or upload a new one using multipart/form-data.
	// The file must be at most maxDocumentSizeMB MB in size.
	Document string `json:"document"`

	// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
	// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	Thumbnail string `json:"thumbnail,omitempty"`

	// Disables automatic server-side content type detection for files uploaded using multipart/form-data
	//  - Currently, This option is unsupported and will be ignored.
	DisableContentTypeDetection bool `json:"disable_content_type_detection,omitempty"`
}

// NewDocumentMessage creates a new DocumentMessage instance.
// Based on SendDocumentParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddocument
//   - Only chat_id and document are required.
//   - Currently, only file_id or http URL is supported.
func NewDocumentMessage(chatID string, document string) *DocumentMessage {
	return Document().Chat(chatID).File(document).Build()
}

func (m *DocumentMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Document == "" {
		return core.NewParamError("document cannot be empty")
	}
	return nil
}
