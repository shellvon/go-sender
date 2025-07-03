package telegram

// VoiceBuilder constructs Telegram voice messages using a fluent, type-safe API.
// Example usage:
//   msg := telegram.Voice().
//            Chat("123").
//            File("BQACAgQAAxkBA...").
//            Duration(60).
//            Caption("voice note").
//            Build()

// VoiceBuilder provides chainable setters for fields relevant to *VoiceMessage.
// It embeds a generic *baseBuilder so that common Chat/Silent/Protect methods are
// inherited and still return the concrete *VoiceBuilder type.
//
// Only the subset of fields exposed by VoiceMessage are present here, ensuring
// compile-time safety and IDE auto-completion.
type VoiceBuilder struct {
	*mediaBuilder[*VoiceBuilder]

	voice    string
	duration int
}

// Voice returns a new VoiceBuilder instance.
func Voice() *VoiceBuilder {
	b := &VoiceBuilder{}
	b.mediaBuilder = newMediaBuilder(b)
	return b
}

// File sets the required voice file identifier/URL/attach:// reference.
// Based on SendVoiceParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvoice
// Voice to send. Pass a file_id as String to send a file that exists on the Telegram servers (recommended),
// pass an HTTP URL as a String for Telegram to get a file from the Internet, or upload a new one using multipart/form-data.
// The audio must be at most maxVoiceSizeMB MB in size.
//   - Currently, only file_id or http URL is supported.
func (b *VoiceBuilder) File(id string) *VoiceBuilder {
	b.voice = id
	return b
}

// Duration sets the duration of the voice message in seconds.
// Based on SendVoiceParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvoice
// Duration of the voice message in seconds.
func (b *VoiceBuilder) Duration(sec int) *VoiceBuilder {
	b.duration = sec
	return b
}

// Build assembles the *VoiceMessage ready for sending via the provider.
func (b *VoiceBuilder) Build() *VoiceMessage {
	msg := &VoiceMessage{
		MediaMessage: MediaMessage{
			BaseMessage:     b.mediaBuilder.toBaseMessage(TypeVoice),
			Caption:         b.caption,
			ParseMode:       b.parseMode,
			CaptionEntities: b.entities,
		},
		Voice:    b.voice,
		Duration: b.duration,
	}
	return msg
}
