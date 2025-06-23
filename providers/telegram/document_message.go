package telegram

import "github.com/shellvon/go-sender/core"

// DocumentMessage represents a document message for Telegram
// Based on SendDocumentParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddocument
type DocumentMessage struct {
	MediaMessage

	// File to send. Pass a file_id as String to send a file that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get a file from the Internet, or upload a new one using multipart/form-data.
	// The file must be at most 50 MB in size.
	Document string `json:"document"`

	// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
	// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	Thumbnail string `json:"thumbnail,omitempty"`

	// Disables automatic server-side content type detection for files uploaded using multipart/form-data
	DisableContentTypeDetection bool `json:"disable_content_type_detection,omitempty"`
}

func (m *DocumentMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *DocumentMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
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

type DocumentMessageOption func(*DocumentMessage)

// WithDocumentThumbnail sets the thumbnail for the document
// Should be in JPEG format and less than 200 kB in size
// Width and height should not exceed 320
func WithDocumentThumbnail(thumbnail string) DocumentMessageOption {
	return func(m *DocumentMessage) { m.Thumbnail = thumbnail }
}

// WithDocumentDisableContentTypeDetection disables automatic server-side content type detection
// This is useful when you want to control the MIME type of the uploaded file
func WithDocumentDisableContentTypeDetection(disable bool) DocumentMessageOption {
	return func(m *DocumentMessage) { m.DisableContentTypeDetection = disable }
}

func NewDocumentMessage(chatID string, document string, opts ...interface{}) *DocumentMessage {
	msg := &DocumentMessage{
		MediaMessage: MediaMessage{
			BaseMessage: BaseMessage{
				MsgType: TypeDocument,
				ChatID:  chatID,
			},
		},
		Document: document,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case DocumentMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}
