package telegram

// AudioBuilder constructs Telegram audio messages.
// Example:
//   msg := telegram.Audio().
//            Chat("123").
//            File("BQBCAgMEBQ").
//            Title("Podcast").
//            Performer("John").
//            Duration(180).
//            Build()

type AudioBuilder struct {
	*mediaBuilder[*AudioBuilder]

	audio string

	duration  int
	performer string
	title     string
	thumb     string
}

// Audio returns a new AudioBuilder.
func Audio() *AudioBuilder {
	b := &AudioBuilder{}
	b.mediaBuilder = newMediaBuilder(b)
	return b
}

// File sets the audio file (file_id / URL). Required.
// Based on SendAudioParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendaudio
// Audio file to send. Pass a file_id as String to send an audio file that exists on the Telegram servers (recommended),
// pass an HTTP URL as a String for Telegram to get an audio file from the Internet, or upload a new one using multipart/form-data.
// The audio must be at most maxAudioSizeMB MB in size.
//   - Currently, only file_id or http URL is supported.
func (b *AudioBuilder) File(id string) *AudioBuilder {
	b.audio = id
	return b
}

// Duration sets audio duration in seconds.
// Based on SendAudioParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendaudio
// Duration of the audio in seconds.
func (b *AudioBuilder) Duration(sec int) *AudioBuilder {
	b.duration = sec
	return b
}

// Performer sets performer field.
// Based on SendAudioParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendaudio
// Performer of the audio.
func (b *AudioBuilder) Performer(p string) *AudioBuilder {
	b.performer = p
	return b
}

// Title sets track title.
// Based on SendAudioParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendaudio
// Track name.
func (b *AudioBuilder) Title(t string) *AudioBuilder {
	b.title = t
	return b
}

// Thumbnail sets custom thumbnail file_id or attach://.
// Based on SendAudioParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendaudio
// Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
// The thumbnail should be in JPEG format and less than 200 kB in size. A thumbnail's width and height should not exceed 320.
// Ignored if the file is not uploaded using multipart/form-data. Thumbnails can't be reused and can be only uploaded as a new file,
// so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
//   - Currently, only file_id or http URL is supported.
func (b *AudioBuilder) Thumbnail(th string) *AudioBuilder {
	b.thumb = th
	return b
}

// Build converts builder into *AudioMessage.
func (b *AudioBuilder) Build() *AudioMessage {
	msg := &AudioMessage{
		MediaMessage: MediaMessage{
			BaseMessage:     b.mediaBuilder.toBaseMessage(TypeAudio),
			Caption:         b.caption,
			ParseMode:       b.parseMode,
			CaptionEntities: b.entities,
		},
		Audio:     b.audio,
		Duration:  b.duration,
		Performer: b.performer,
		Title:     b.title,
		Thumbnail: b.thumb,
	}
	return msg
}
