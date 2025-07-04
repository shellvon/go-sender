package email_test

import (
	"testing"

	"github.com/shellvon/go-sender/providers/email"
)

func TestNewMessage(t *testing.T) {
	to := []string{"test@example.com"}
	body := "Test email body"

	msg := email.NewMessage(to, body)

	if msg == nil {
		t.Fatal("NewMessage returned nil")
	}

	if len(msg.To) != 1 || msg.To[0] != "test@example.com" {
		t.Errorf("Expected To to be %v, got %v", to, msg.To)
	}

	if msg.Body != body {
		t.Errorf("Expected Body to be %s, got %s", body, msg.Body)
	}

	if msg.IsHTML {
		t.Error("Expected IsHTML to be false by default")
	}
}

func TestMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		msg     *email.Message
		wantErr bool
	}{
		{
			name:    "valid message",
			msg:     email.NewMessage([]string{"test@example.com"}, "Test body"),
			wantErr: false,
		},
		{
			name:    "empty recipients",
			msg:     email.NewMessage([]string{}, "Test body"),
			wantErr: true,
		},
		{
			name:    "nil recipients",
			msg:     email.NewMessage(nil, "Test body"),
			wantErr: true,
		},
		{
			name:    "empty body",
			msg:     email.NewMessage([]string{"test@example.com"}, ""),
			wantErr: true,
		},
		{
			name:    "invalid recipient email",
			msg:     email.NewMessage([]string{"invalid-email"}, "Test body"),
			wantErr: true,
		},
		{
			name:    "invalid CC email",
			msg:     email.Email().To("test@example.com").Cc("invalid-cc").Body("Test body").Build(),
			wantErr: true,
		},
		{
			name:    "invalid BCC email",
			msg:     email.Email().To("test@example.com").Bcc("invalid-bcc").Body("Test body").Build(),
			wantErr: true,
		},
		{
			name:    "valid with name format",
			msg:     email.NewMessage([]string{"Test User <test@example.com>"}, "Test body"),
			wantErr: false,
		},
		{
			name: "valid with CC and BCC",
			msg: email.Email().
				To("test@example.com").
				Cc("cc@example.com").
				Bcc("bcc@example.com").
				Body("Test body").
				Build(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageOptions(t *testing.T) {
	msg := email.Email().
		To("test@example.com").
		Body("Test body").
		From("sender@example.com").
		Subject("Test Subject").
		Cc("cc@example.com").
		Bcc("bcc@example.com").
		ReplyTo("reply1@example.com").
		HTML().
		AddAttach("file1.txt", "file2.pdf").
		Build()

	if msg.From != "sender@example.com" {
		t.Errorf("Expected From to be 'sender@example.com', got %s", msg.From)
	}

	if msg.Subject != "Test Subject" {
		t.Errorf("Expected Subject to be 'Test Subject', got %s", msg.Subject)
	}

	if len(msg.Cc) != 1 || msg.Cc[0] != "cc@example.com" {
		t.Errorf("Expected Cc to be ['cc@example.com'], got %v", msg.Cc)
	}

	if len(msg.Bcc) != 1 || msg.Bcc[0] != "bcc@example.com" {
		t.Errorf("Expected Bcc to be ['bcc@example.com'], got %v", msg.Bcc)
	}

	if msg.ReplyTo != "reply1@example.com" {
		t.Errorf("Expected ReplyTo to be 'reply1@example.com', got %s", msg.ReplyTo)
	}

	if !msg.IsHTML {
		t.Error("Expected IsHTML to be true")
	}

	if len(msg.Attachments) != 2 {
		t.Errorf("Expected 2 attachments, got %d", len(msg.Attachments))
	}

	expectedAttachments := []string{"file1.txt", "file2.pdf"}
	for i, att := range msg.Attachments {
		if att != expectedAttachments[i] {
			t.Errorf("Expected attachment %d to be %s, got %s", i, expectedAttachments[i], att)
		}
	}
}

func TestMessageProviderType(t *testing.T) {
	msg := email.NewMessage([]string{"test@example.com"}, "Test body")

	providerType := msg.ProviderType()
	if providerType != "email" {
		t.Errorf("Expected ProviderType to be 'email', got %s", providerType)
	}
}

func TestMessageMsgID(t *testing.T) {
	msg := email.NewMessage([]string{"test@example.com"}, "Test body")

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

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid with name", "Test User <test@example.com>", false},
		{"empty email", "", true},
		{"invalid format", "invalid-email", true},
		{"missing domain", "test@", true},
		{"missing local", "@example.com", true},
		{"multiple @", "test@@example.com", true},
		{"valid complex", "test.user+tag@example.co.uk", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := email.NewMessage([]string{tt.email}, "Test body")
			err := msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Email validation for %s: error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}
