//nolint:dupl // intentional duplication for explicit message type separation
package telegram

import "github.com/shellvon/go-sender/core"

// VoiceMessage represents a voice message for Telegram
// Based on SendVoiceParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvoice
type VoiceMessage struct {
	MediaMessage

	// Audio file to send. Pass a file_id as String to send a file that exists on the Telegram servers (recommended),
	// pass an HTTP URL as a String for Telegram to get a file from the Internet, or upload a new one using multipart/form-data.
	// The audio must be at most maxVoiceSizeMB MB in size.
	Voice string `json:"voice"`

	// Duration of the voice message in seconds
	Duration int `json:"duration,omitempty"`
}

// NewVoiceMessage creates a new VoiceMessage instance.
func NewVoiceMessage(chatID string, voice string, opts ...interface{}) *VoiceMessage {
	return NewVoiceMessageWithBuilder(chatID, voice, opts...)
}

func (m *VoiceMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *VoiceMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *VoiceMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Voice == "" {
		return core.NewParamError("voice cannot be empty")
	}
	return nil
}

type VoiceMessageOption func(*VoiceMessage)

// WithVoiceDuration sets the duration of the voice message in seconds
// This is optional and can be used to provide metadata about the voice message.
func WithVoiceDuration(duration int) VoiceMessageOption {
	return func(m *VoiceMessage) { m.Duration = duration }
}
