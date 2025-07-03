package telegram

// PhotoBuilder constructs Telegram photo messages using a fluent API.
// Example:
//   msg := telegram.Photo().
//            Chat("123456").
//            File("AgACAgQAAxkBA...").
//            Caption("Sunset").
//            HasSpoiler(true).
//            Build()

type PhotoBuilder struct {
	*mediaBuilder[*PhotoBuilder]

	photo   string
	spoiler bool
}

// Photo returns a new PhotoBuilder.
func Photo() *PhotoBuilder {
	b := &PhotoBuilder{}
	b.mediaBuilder = newMediaBuilder(b)
	return b
}

// File sets the photo file (file_id, http URL). Required.
// Based on SendPhotoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendphoto
// Photo to send. Pass a file_id as String to send a photo that exists on the Telegram servers (recommended),
// pass an HTTP URL as a String for Telegram to get a photo from the Internet, or upload a new photo using multipart/form-data.
// The photo must be at most 10 MB in size. The photo's width and height must not exceed maxPhotoDimensionSum in total.
// Width and height ratio must be at most 20.
//   - Currently, only file_id or http URL is supported.
func (b *PhotoBuilder) File(id string) *PhotoBuilder {
	b.photo = id
	return b
}

// HasSpoiler toggles has_spoiler.
// Based on SendPhotoParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendphoto
// Pass True if the photo needs to be covered with a spoiler animation.
func (b *PhotoBuilder) HasSpoiler(spoiler bool) *PhotoBuilder {
	b.spoiler = spoiler
	return b
}

// Build converts builder state into *PhotoMessage.
func (b *PhotoBuilder) Build() *PhotoMessage {
	msg := &PhotoMessage{
		MediaMessage: MediaMessage{
			BaseMessage:           b.mediaBuilder.toBaseMessage(TypePhoto),
			Caption:               b.caption,
			ParseMode:             b.parseMode,
			CaptionEntities:       b.entities,
			ShowCaptionAboveMedia: b.showCaptionTop,
		},
		Photo:      b.photo,
		HasSpoiler: b.spoiler,
	}
	return msg
}
