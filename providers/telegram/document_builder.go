package telegram

// DocumentBuilder constructs Telegram document messages with a fluent API.
// Example:
//   msg := telegram.Document().
//            Chat("123").
//            File("BQACAgQAAxkBA...").
//            Caption("report.pdf").
//            Thumbnail("thumb.jpg").
//            DisableContentTypeDetection(true).
//            Build()

type DocumentBuilder struct {
	*mediaBuilder[*DocumentBuilder]

	document                    string
	thumbnail                   string
	disableContentTypeDetection bool
}

// Document returns a new DocumentBuilder.
func Document() *DocumentBuilder {
	b := &DocumentBuilder{}
	b.mediaBuilder = newMediaBuilder(b)
	return b
}

// File sets the required document file (file_id / URL).
// Based on SendDocumentParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddocument
// File to send. Pass a file_id as String to send a file that exists on the Telegram servers (recommended),
// pass an HTTP URL as a String for Telegram to get a file from the Internet, or upload a new one using multipart/form-data.
// The file must be at most maxDocumentSizeMB MB in size.
//   - Currently, only file_id or http URL is supported.
func (b *DocumentBuilder) File(id string) *DocumentBuilder {
	b.document = id
	return b
}

// Thumbnail sets a custom thumbnail file_id or attach:// reference.
// Based on SendDocumentParams from Telegram Bot API
// https://core.telegram.org/bots/api#senddocument
// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
func (b *DocumentBuilder) Thumbnail(th string) *DocumentBuilder {
	b.thumbnail = th
	return b
}

// Build assembles the *DocumentMessage.
func (b *DocumentBuilder) Build() *DocumentMessage {
	msg := &DocumentMessage{
		MediaMessage: MediaMessage{
			BaseMessage:     b.mediaBuilder.toBaseMessage(TypeDocument),
			Caption:         b.caption,
			ParseMode:       b.parseMode,
			CaptionEntities: b.entities,
		},
		Document:                    b.document,
		Thumbnail:                   b.thumbnail,
		DisableContentTypeDetection: b.disableContentTypeDetection,
	}
	return msg
}
