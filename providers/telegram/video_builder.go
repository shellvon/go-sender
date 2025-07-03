package telegram

// VideoBuilder constructs Telegram video messages using a fluent API.
// Example:
//   msg := telegram.Video().
//            Chat("123").
//            File("DQACAgQAAxkBA...").
//            Caption("demo video").
//            Duration(15).
//            Width(640).
//            Height(360).
//            SupportsStreaming(true).
//            Build()

type VideoBuilder struct {
	*mediaBuilder[*VideoBuilder]

	video string

	duration          int
	width             int
	height            int
	thumbnail         string
	cover             string
	startTS           int
	hasSpoiler        bool
	supportsStreaming bool
}

// Video returns a new VideoBuilder.
func Video() *VideoBuilder {
	b := &VideoBuilder{}
	b.mediaBuilder = newMediaBuilder(b)
	return b
}

// File sets the required video identifier.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Video to send. Pass a file_id as String to send a video that exists on the Telegram servers (recommended),
// pass an HTTP URL as a String for Telegram to get a video from the Internet, or upload a new video using multipart/form-data.
// The video must be at most maxVideoSizeMB MB in size.
//   - Currently, only file_id or http URL is supported.
func (b *VideoBuilder) File(id string) *VideoBuilder {
	b.video = id
	return b
}

// Duration sets video duration seconds.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
func (b *VideoBuilder) Duration(sec int) *VideoBuilder {
	b.duration = sec
	return b
}

// Width sets video width.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Video width.
func (b *VideoBuilder) Width(w int) *VideoBuilder {
	b.width = w
	return b
}

// Height sets video height.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Video height.
func (b *VideoBuilder) Height(h int) *VideoBuilder {
	b.height = h
	return b
}

// Thumbnail sets custom thumbnail file_id / attach://.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
//   - Currently, only file_id or http URL is supported.
func (b *VideoBuilder) Thumbnail(th string) *VideoBuilder {
	b.thumbnail = th
	return b
}

// Cover sets separate cover image for the video.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Cover for the video in the message. Pass a file_id to send a file that exists on the Telegram servers (recommended),
// pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
func (b *VideoBuilder) Cover(c string) *VideoBuilder {
	b.cover = c
	return b
}

// StartTimestamp sets start timestamp in seconds.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Start timestamp for the video in the message.
func (b *VideoBuilder) StartTimestamp(ts int) *VideoBuilder {
	b.startTS = ts
	return b
}

// HasSpoiler toggles spoiler overlay.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Pass True if the video needs to be covered with a spoiler animation.
func (b *VideoBuilder) HasSpoiler(sp bool) *VideoBuilder {
	b.hasSpoiler = sp
	return b
}

// SupportsStreaming toggles streaming suitability flag.
// Based on SendVideoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvideo
// Pass True if the uploaded video is suitable for streaming.
func (b *VideoBuilder) SupportsStreaming(s bool) *VideoBuilder {
	b.supportsStreaming = s
	return b
}

// Build assembles the *VideoMessage.
func (b *VideoBuilder) Build() *VideoMessage {
	msg := &VideoMessage{
		MediaMessage: MediaMessage{
			BaseMessage:           b.mediaBuilder.toBaseMessage(TypeVideo),
			Caption:               b.caption,
			ParseMode:             b.parseMode,
			CaptionEntities:       b.entities,
			ShowCaptionAboveMedia: b.showCaptionTop,
		},
		Video:             b.video,
		Duration:          b.duration,
		Width:             b.width,
		Height:            b.height,
		Thumbnail:         b.thumbnail,
		Cover:             b.cover,
		StartTimestamp:    b.startTS,
		HasSpoiler:        b.hasSpoiler,
		SupportsStreaming: b.supportsStreaming,
	}
	return msg
}
