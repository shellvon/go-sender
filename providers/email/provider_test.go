package email

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestEmailProvider_New(t *testing.T) {
	config := Config{
		BaseConfig: core.BaseConfig{Strategy: core.StrategyRoundRobin},
		Accounts: []Account{
			{
				Name:     "test",
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password",
				From:     "test@example.com",
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	if provider == nil {
		t.Fatal("Provider should not be nil")
	}

	if len(provider.accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(provider.accounts))
	}
}

func TestEmailProvider_New_NoAccounts(t *testing.T) {
	config := Config{
		BaseConfig: core.BaseConfig{Strategy: core.StrategyRoundRobin},
		Accounts:   []Account{},
	}

	_, err := New(config)
	if err == nil {
		t.Fatal("Expected error when no accounts provided")
	}
}

func TestEmailProvider_New_AllDisabled(t *testing.T) {
	config := Config{
		BaseConfig: core.BaseConfig{Strategy: core.StrategyRoundRobin},
		Accounts: []Account{
			{
				Name:     "test",
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password",
				From:     "test@example.com",
				Disabled: true,
			},
		},
	}

	_, err := New(config)
	if err == nil {
		t.Fatal("Expected error when all accounts are disabled")
	}
}

func TestEmailMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
		wantErr bool
	}{
		{
			name: "valid message",
			message: &Message{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr: false,
		},
		{
			name: "valid message without subject",
			message: &Message{
				To:   []string{"test@example.com"},
				Body: "Test Body",
			},
			wantErr: false,
		},
		{
			name: "missing to",
			message: &Message{
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr: true,
		},
		{
			name: "empty to",
			message: &Message{
				To:      []string{},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr: true,
		},
		{
			name: "missing body",
			message: &Message{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmailProvider_Send_InvalidMessageType(t *testing.T) {
	config := Config{
		BaseConfig: core.BaseConfig{Strategy: core.StrategyRoundRobin},
		Accounts: []Account{
			{
				Name:     "test",
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password",
				From:     "test@example.com",
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Test with wrong message type
	wrongMessage := &mockMessage{}
	err = provider.Send(context.Background(), wrongMessage, nil)
	if err == nil {
		t.Fatal("Expected error when sending wrong message type")
	}
}

func TestNewMessageWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		to       []string
		body     string
		opts     []MessageOption
		expected *Message
	}{
		{
			name: "basic message",
			to:   []string{"test@example.com"},
			body: "Test body",
			opts: []MessageOption{},
			expected: &Message{
				To:   []string{"test@example.com"},
				Body: "Test body",
			},
		},
		{
			name: "message with subject",
			to:   []string{"test@example.com"},
			body: "Test body",
			opts: []MessageOption{
				WithSubject("Test Subject"),
			},
			expected: &Message{
				To:      []string{"test@example.com"},
				Body:    "Test body",
				Subject: "Test Subject",
			},
		},
		{
			name: "message with all options",
			to: []string{
				"test@example.com",
				"test2@example.com",
			},
			body: "Test body",
			opts: []MessageOption{
				WithFrom("sender@example.com"),
				WithSubject("Test Subject"),
				WithCc("cc@example.com"),
				WithBcc("bcc@example.com"),
				WithReplyTo("reply@example.com"),
				WithHTML(),
				WithAttachments("file1.pdf", "file2.txt"),
			},
			expected: &Message{
				To: []string{
					"test@example.com",
					"test2@example.com",
				},
				Body:        "Test body",
				From:        "sender@example.com",
				Subject:     "Test Subject",
				Cc:          []string{"cc@example.com"},
				Bcc:         []string{"bcc@example.com"},
				ReplyTo:     "reply@example.com",
				IsHTML:      true,
				Attachments: []string{"file1.pdf", "file2.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage(tt.to, tt.body, tt.opts...)

			// Check required fields
			if len(msg.To) != len(tt.expected.To) {
				t.Errorf("To field mismatch: got %v, want %v", msg.To, tt.expected.To)
			}
			if msg.Body != tt.expected.Body {
				t.Errorf("Body field mismatch: got %s, want %s", msg.Body, tt.expected.Body)
			}

			// Check optional fields
			if tt.expected.From != "" {
				if msg.From != tt.expected.From {
					t.Errorf("From field mismatch: got %s, want %s", msg.From, tt.expected.From)
				}
			}
			if msg.Subject != tt.expected.Subject {
				t.Errorf("Subject field mismatch: got %s, want %s", msg.Subject, tt.expected.Subject)
			}
			if len(msg.Cc) != len(tt.expected.Cc) {
				t.Errorf("Cc field mismatch: got %v, want %v", msg.Cc, tt.expected.Cc)
			}
			if len(msg.Bcc) != len(tt.expected.Bcc) {
				t.Errorf("Bcc field mismatch: got %v, want %v", msg.Bcc, tt.expected.Bcc)
			}
			if msg.IsHTML != tt.expected.IsHTML {
				t.Errorf("IsHTML field mismatch: got %t, want %t", msg.IsHTML, tt.expected.IsHTML)
			}
			if len(msg.Attachments) != len(tt.expected.Attachments) {
				t.Errorf("Attachments field mismatch: got %v, want %v", msg.Attachments, tt.expected.Attachments)
			}
		})
	}
}
