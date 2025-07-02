package telegram_test

import (
	"testing"

	"github.com/shellvon/go-sender/providers/telegram"
)

func TestNewTextMessage(t *testing.T) {
	msg := telegram.NewTextMessage("123456789", "Hello, World!")

	if msg == nil {
		t.Fatal("NewTextMessage returned nil")
	}

	if msg.ChatID != "123456789" {
		t.Errorf("Expected ChatID to be '123456789', got '%s'", msg.ChatID)
	}

	if msg.Text != "Hello, World!" {
		t.Errorf("Expected Text to be 'Hello, World!', got '%s'", msg.Text)
	}

	if msg.MsgType != telegram.TypeText {
		t.Errorf("Expected MsgType to be TypeText, got %v", msg.MsgType)
	}

	if msg.ProviderType() != "telegram" {
		t.Errorf("Expected ProviderType to be 'telegram', got '%s'", msg.ProviderType())
	}
}

func TestTextMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *telegram.TextMessage
		wantErr bool
	}{
		{
			name:    "valid message",
			msg:     telegram.NewTextMessage("123456789", "Hello, World!"),
			wantErr: false,
		},
		{
			name:    "empty chat_id",
			msg:     telegram.NewTextMessage("", "Hello, World!"),
			wantErr: true,
		},
		{
			name:    "empty text",
			msg:     telegram.NewTextMessage("123456789", ""),
			wantErr: true,
		},
		{
			name:    "whitespace text",
			msg:     telegram.NewTextMessage("123456789", "   "),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TextMessage.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTextMessage_WithDisableWebPreview(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithDisableWebPreview(true),
	)

	if !msg.DisableWebPreview {
		t.Error("Expected DisableWebPreview to be true")
	}
}

func TestTextMessage_WithReplyTo(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithReplyTo(12345),
	)

	if msg.ReplyToMessageID != 12345 {
		t.Errorf("Expected ReplyToMessageID to be 12345, got %d", msg.ReplyToMessageID)
	}
}

func TestTextMessage_WithParseMode(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithParseMode("HTML"),
	)

	if msg.ParseMode != "HTML" {
		t.Errorf("Expected ParseMode to be 'HTML', got '%s'", msg.ParseMode)
	}
}

func TestTextMessage_WithSilent(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithSilent(true),
	)

	if !msg.DisableNotification {
		t.Error("Expected DisableNotification to be true")
	}
}

func TestTextMessage_WithProtectContent(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithProtectContent(true),
	)

	if !msg.ProtectContent {
		t.Error("Expected ProtectContent to be true")
	}
}

func TestTextMessage_WithAllowPaidBroadcast(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithAllowPaidBroadcast(true),
	)

	if !msg.AllowPaidBroadcast {
		t.Error("Expected AllowPaidBroadcast to be true")
	}
}

func TestTextMessage_WithMessageEffectID(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithMessageEffectID("effect123"),
	)

	if msg.MessageEffectID != "effect123" {
		t.Errorf("Expected MessageEffectID to be 'effect123', got '%s'", msg.MessageEffectID)
	}
}

func TestTextMessage_WithBusinessConnectionID(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithBusinessConnectionID("business123"),
	)

	if msg.BusinessConnectionID != "business123" {
		t.Errorf("Expected BusinessConnectionID to be 'business123', got '%s'", msg.BusinessConnectionID)
	}
}

func TestTextMessage_WithMessageThreadID(t *testing.T) {
	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithMessageThreadID(67890),
	)

	if msg.MessageThreadID != 67890 {
		t.Errorf("Expected MessageThreadID to be 67890, got %d", msg.MessageThreadID)
	}
}

func TestTextMessage_GetBase(t *testing.T) {
	msg := telegram.NewTextMessage("123456789", "Hello, World!")
	base := msg.GetBase()

	if base == nil {
		t.Fatal("GetBase returned nil")
	}

	if base.ChatID != "123456789" {
		t.Errorf("Expected base ChatID to be '123456789', got '%s'", base.ChatID)
	}

	if base.MsgType != telegram.TypeText {
		t.Errorf("Expected base MsgType to be TypeText, got %v", base.MsgType)
	}
}

func TestTextMessage_MsgID(t *testing.T) {
	msg := telegram.NewTextMessage("123456789", "Hello, World!")

	msgID := msg.MsgID()
	if msgID == "" {
		t.Error("Expected MsgID to be non-empty")
	}

	// MsgID should be consistent for the same message
	msgID2 := msg.MsgID()
	if msgID != msgID2 {
		t.Error("Expected MsgID to be consistent")
	}
}
