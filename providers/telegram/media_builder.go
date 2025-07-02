package telegram

// MediaMessageBuilder is a generic builder for creating media messages
// This reduces duplicate code across different media message types.
type MediaMessageBuilder struct {
	msgType MessageType
	chatID  string
	media   string
	opts    []interface{}
}

// NewMediaMessageBuilder creates a new media message builder.
func NewMediaMessageBuilder(msgType MessageType, chatID, media string) *MediaMessageBuilder {
	return &MediaMessageBuilder{
		msgType: msgType,
		chatID:  chatID,
		media:   media,
	}
}

// WithOptions adds options to the builder.
func (b *MediaMessageBuilder) WithOptions(opts ...interface{}) *MediaMessageBuilder {
	b.opts = opts
	return b
}

// BuildMediaMessage creates a media message with the given type.
func (b *MediaMessageBuilder) BuildMediaMessage() MessageWithBase {
	baseMsg := BaseMessage{
		MsgType: b.msgType,
		ChatID:  b.chatID,
	}

	mediaMsg := MediaMessage{
		BaseMessage: baseMsg,
	}

	// Create the specific message type based on msgType
	var msg MessageWithBase
	switch b.msgType {
	case TypePhoto:
		photoMsg := &PhotoMessage{
			MediaMessage: mediaMsg,
			Photo:        b.media,
		}
		msg = photoMsg
	case TypeVoice:
		voiceMsg := &VoiceMessage{
			MediaMessage: mediaMsg,
			Voice:        b.media,
		}
		msg = voiceMsg
	case TypeAudio:
		audioMsg := &AudioMessage{
			MediaMessage: mediaMsg,
			Audio:        b.media,
		}
		msg = audioMsg
	case TypeDocument:
		docMsg := &DocumentMessage{
			MediaMessage: mediaMsg,
			Document:     b.media,
		}
		msg = docMsg
	case TypeVideo:
		videoMsg := &VideoMessage{
			MediaMessage: mediaMsg,
			Video:        b.media,
		}
		msg = videoMsg
	case TypeAnimation:
		animMsg := &AnimationMessage{
			MediaMessage: mediaMsg,
			Animation:    b.media,
		}
		msg = animMsg
	case TypeVideoNote:
		videoNoteMsg := &VideoNoteMessage{
			BaseMessage: baseMsg,
			VideoNote:   b.media,
		}
		msg = videoNoteMsg
	case TypeText:
		// Text messages don't use media parameter - use NewTextMessage instead
		// This case is included for completeness but should not be used
		msg = &MediaMessage{
			BaseMessage: baseMsg,
		}
	case TypeLocation:
		// Location messages require latitude/longitude coordinates, not media files
		// Use NewLocationMessage instead
		msg = &MediaMessage{
			BaseMessage: baseMsg,
		}
	case TypeContact:
		// Contact messages require phone number and name, not media files
		// Use NewContactMessage instead
		msg = &MediaMessage{
			BaseMessage: baseMsg,
		}
	case TypePoll:
		// Poll messages require question and options, not media files
		// Use NewPollMessage instead
		msg = &MediaMessage{
			BaseMessage: baseMsg,
		}
	case TypeDice:
		// Dice messages only require emoji, not media files
		// Use NewDiceMessage instead
		msg = &MediaMessage{
			BaseMessage: baseMsg,
		}
	default:
		// Fallback to a generic media message
		msg = &MediaMessage{
			BaseMessage: baseMsg,
		}
	}

	// Apply options
	applyMediaMessageOptions(msg, b.opts)
	return msg
}

// NewPhotoMessageWithBuilder creates a new PhotoMessage with builder options.
func NewPhotoMessageWithBuilder(chatID, photo string, opts ...interface{}) *PhotoMessage {
	builder := NewMediaMessageBuilder(TypePhoto, chatID, photo).WithOptions(opts...)
	msg := builder.BuildMediaMessage()
	if photoMsg, ok := msg.(*PhotoMessage); ok {
		return photoMsg
	}
	// Fallback - this should never happen with correct msgType
	return &PhotoMessage{}
}

func NewVoiceMessageWithBuilder(chatID, voice string, opts ...interface{}) *VoiceMessage {
	builder := NewMediaMessageBuilder(TypeVoice, chatID, voice).WithOptions(opts...)
	msg := builder.BuildMediaMessage()
	if voiceMsg, ok := msg.(*VoiceMessage); ok {
		return voiceMsg
	}
	// Fallback - this should never happen with correct msgType
	return &VoiceMessage{}
}

func NewAudioMessageWithBuilder(chatID, audio string, opts ...interface{}) *AudioMessage {
	builder := NewMediaMessageBuilder(TypeAudio, chatID, audio).WithOptions(opts...)
	msg := builder.BuildMediaMessage()
	if audioMsg, ok := msg.(*AudioMessage); ok {
		return audioMsg
	}
	// Fallback - this should never happen with correct msgType
	return &AudioMessage{}
}

func NewDocumentMessageWithBuilder(chatID, document string, opts ...interface{}) *DocumentMessage {
	builder := NewMediaMessageBuilder(TypeDocument, chatID, document).WithOptions(opts...)
	msg := builder.BuildMediaMessage()
	if docMsg, ok := msg.(*DocumentMessage); ok {
		return docMsg
	}
	// Fallback - this should never happen with correct msgType
	return &DocumentMessage{}
}

func NewVideoMessageWithBuilder(chatID, video string, opts ...interface{}) *VideoMessage {
	builder := NewMediaMessageBuilder(TypeVideo, chatID, video).WithOptions(opts...)
	msg := builder.BuildMediaMessage()
	if videoMsg, ok := msg.(*VideoMessage); ok {
		return videoMsg
	}
	// Fallback - this should never happen with correct msgType
	return &VideoMessage{}
}

func NewAnimationMessageWithBuilder(chatID, animation string, opts ...interface{}) *AnimationMessage {
	builder := NewMediaMessageBuilder(TypeAnimation, chatID, animation).WithOptions(opts...)
	msg := builder.BuildMediaMessage()
	if animMsg, ok := msg.(*AnimationMessage); ok {
		return animMsg
	}
	// Fallback - this should never happen with correct msgType
	return &AnimationMessage{}
}

func NewVideoNoteMessageWithBuilder(chatID, videoNote string, opts ...interface{}) *VideoNoteMessage {
	builder := NewMediaMessageBuilder(TypeVideoNote, chatID, videoNote).WithOptions(opts...)
	msg := builder.BuildMediaMessage()
	if videoNoteMsg, ok := msg.(*VideoNoteMessage); ok {
		return videoNoteMsg
	}
	// Fallback - this should never happen with correct msgType
	return &VideoNoteMessage{}
}
