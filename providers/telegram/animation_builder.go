package telegram

// AnimationBuilder constructs Telegram animation (GIF/WEBM) messages.
// Example:
//   msg := telegram.Animation().
//            Chat("123").
//            File("CgACAgQAAxkBA...").
//            Caption("funny gif").
//            Duration(3).
//            Width(320).
//            Height(240).
//            HasSpoiler(true).
//            Build()

type AnimationBuilder struct {
	*mediaBuilder[*AnimationBuilder]

	animation  string
	duration   int
	width      int
	height     int
	thumbnail  string
	hasSpoiler bool
}

// Animation returns a new AnimationBuilder.
func Animation() *AnimationBuilder {
	b := &AnimationBuilder{}
	b.mediaBuilder = newMediaBuilder(b)
	return b
}

// File sets the required animation file id / URL / attach:// reference.
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
// Animation to send. Pass a file_id as String to send an animation that exists on the Telegram servers (recommended),
// pass an HTTP URL as a String for Telegram to get an animation from the Internet, or upload a new animation using multipart/form-data.
// The animation must be at most maxAnimationSizeMB MB in size.
//   - Currently, only file_id or http URL is supported.
func (b *AnimationBuilder) File(id string) *AnimationBuilder {
	b.animation = id
	return b
}

// Duration sets animation duration seconds.
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
// Duration of sent animation in seconds.
func (b *AnimationBuilder) Duration(sec int) *AnimationBuilder {
	b.duration = sec
	return b
}

// Width sets animation width.
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
// Animation width.
func (b *AnimationBuilder) Width(w int) *AnimationBuilder {
	b.width = w
	return b
}

// Height sets animation height.
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
// Animation height.
func (b *AnimationBuilder) Height(h int) *AnimationBuilder {
	b.height = h
	return b
}

// Thumbnail sets custom thumbnail file_id or attach:// reference.
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
//   - Currently, only file_id or http URL is supported.
func (b *AnimationBuilder) Thumbnail(th string) *AnimationBuilder {
	b.thumbnail = th
	return b
}

// HasSpoiler toggles spoiler animation overlay.
// Based on SendAnimationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendanimation
// Pass True if the animation needs to be covered with a spoiler animation.
func (b *AnimationBuilder) HasSpoiler(spoiler bool) *AnimationBuilder {
	b.hasSpoiler = spoiler
	return b
}

// Build assembles the *AnimationMessage.
func (b *AnimationBuilder) Build() *AnimationMessage {
	msg := &AnimationMessage{
		MediaMessage: MediaMessage{
			BaseMessage:           b.mediaBuilder.toBaseMessage(TypeAnimation),
			Caption:               b.caption,
			ParseMode:             b.parseMode,
			CaptionEntities:       b.entities,
			ShowCaptionAboveMedia: b.showCaptionTop,
		},
		Animation:  b.animation,
		Duration:   b.duration,
		Width:      b.width,
		Height:     b.height,
		Thumbnail:  b.thumbnail,
		HasSpoiler: b.hasSpoiler,
	}
	return msg
}
