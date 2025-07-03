//nolint:dupl // builder pattern files share similar boilerplate; acceptable duplication
package telegram

// VideoNoteBuilder constructs Telegram video note messages.
// Example:
//   msg := telegram.VideoNote().
//            Chat("123").
//            File("CAACAgQAAxkBA...").
//            Duration(10).
//            Length(240).
//            HasThumbnail("thumb.jpg").
//            Build()

type VideoNoteBuilder struct {
	*baseBuilder[*VideoNoteBuilder]

	videoNote string
	duration  int
	length    int
	thumbnail string
}

// VideoNote returns a new VideoNoteBuilder.
func VideoNote() *VideoNoteBuilder {
	b := &VideoNoteBuilder{}
	b.baseBuilder = &baseBuilder[*VideoNoteBuilder]{self: b}
	return b
}

// File sets the required video note file reference.
// Based on SendVideoNoteParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideonote
// Video note to send. Pass a file_id as String to send a video note that exists on the Telegram servers (recommended),
// pass an HTTP URL as a String for Telegram to get a video note from the Internet, or upload a new video note using multipart/form-data.
// The video note must be at most maxVideoNoteSizeMB MB in size.
//   - Currently, only file_id or http URL is supported.
func (b *VideoNoteBuilder) File(id string) *VideoNoteBuilder {
	b.videoNote = id
	return b
}

// Duration sets the video note duration in seconds.
// Based on SendVideoNoteParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideonote
// Duration of sent video note in seconds.
func (b *VideoNoteBuilder) Duration(sec int) *VideoNoteBuilder {
	b.duration = sec
	return b
}

// Length sets the diameter (width/height) of the video note.
// Based on SendVideoNoteParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideonote
// Video width and height, i.e. diameter of the video message.
func (b *VideoNoteBuilder) Length(l int) *VideoNoteBuilder {
	b.length = l
	return b
}

// Thumbnail sets custom thumbnail reference.
// Based on SendVideoNoteParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideonote
// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
//   - Currently, only file_id or http URL is supported.
func (b *VideoNoteBuilder) Thumbnail(th string) *VideoNoteBuilder {
	b.thumbnail = th
	return b
}

// Build assembles the *VideoNoteMessage.
func (b *VideoNoteBuilder) Build() *VideoNoteMessage {
	msg := &VideoNoteMessage{
		BaseMessage: b.baseBuilder.toBaseMessage(TypeVideoNote),
		VideoNote:   b.videoNote,
		Duration:    b.duration,
		Length:      b.length,
		Thumbnail:   b.thumbnail,
	}
	return msg
}
