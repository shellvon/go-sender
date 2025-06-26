package email

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test configuration structure
type testConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
	Cc       string
	Bcc      string
}

var globalTestConfig testConfig

// hasTestConfig checks if test configuration is available
func hasTestConfig() bool {
	return globalTestConfig.Host != "" && globalTestConfig.Port != 0 &&
		globalTestConfig.Username != "" && globalTestConfig.Password != "" &&
		globalTestConfig.From != "" && globalTestConfig.To != ""
}

// getTestConfig returns the test configuration
func getTestConfig() Config {
	account := Account{
		Name:     "test-account",
		Host:     globalTestConfig.Host,
		Port:     globalTestConfig.Port,
		Username: globalTestConfig.Username,
		Password: globalTestConfig.Password,
		From:     globalTestConfig.From,
	}

	return Config{
		BaseConfig: core.BaseConfig{Strategy: core.StrategyRoundRobin},
		Accounts:   []Account{account},
	}
}

// init initializes test configuration from environment variables
func init() {
	globalTestConfig = testConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     func() int { p, _ := strconv.Atoi(os.Getenv("EMAIL_PORT")); return p }(),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
		To:       os.Getenv("EMAIL_TO"),
		Cc:       os.Getenv("EMAIL_CC"),
		Bcc:      os.Getenv("EMAIL_BCC"),
	}
}

func TestEmailProvider_Integration(t *testing.T) {
	// Skip if no test configuration
	if !hasTestConfig() {
		t.Skip("Skipping integration test - no test configuration available")
	}

	config := getTestConfig()
	provider, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, provider)

	t.Run("SendTextEmail", func(t *testing.T) {
		msg := NewMessage(
			[]string{globalTestConfig.To},
			"This is a test email from go-sender integration test",
			WithSubject("Integration Test - Text Email"),
			WithFrom(globalTestConfig.From),
		)

		err := provider.Send(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("SendHTMLEmail", func(t *testing.T) {
		htmlBody := `
		<html>
			<body>
				<h1>Test HTML Email</h1>
				<p>This is a test HTML email from go-sender integration test.</p>
				<p><strong>Bold text</strong> and <em>italic text</em> should work.</p>
			</body>
		</html>`

		msg := NewMessage(
			[]string{globalTestConfig.To},
			htmlBody,
			WithSubject("Integration Test - HTML Email"),
			WithFrom(globalTestConfig.From),
			WithHTML(),
		)

		err := provider.Send(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("SendEmailWithCC", func(t *testing.T) {
		msg := NewMessage(
			[]string{globalTestConfig.To},
			"This is a test email with CC from go-sender integration test",
			WithSubject("Integration Test - Email with CC"),
			WithFrom(globalTestConfig.From),
			WithCc(globalTestConfig.Cc),
		)

		err := provider.Send(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("SendEmailWithBCC", func(t *testing.T) {
		msg := NewMessage(
			[]string{globalTestConfig.To},
			"This is a test email with BCC from go-sender integration test",
			WithSubject("Integration Test - Email with BCC"),
			WithFrom(globalTestConfig.From),
			WithBcc(globalTestConfig.Bcc),
		)

		err := provider.Send(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("SendEmailWithMultipleRecipients", func(t *testing.T) {
		recipients := []string{globalTestConfig.To}
		if globalTestConfig.Cc != "" {
			recipients = append(recipients, globalTestConfig.Cc)
		}

		msg := NewMessage(
			recipients,
			"This is a test email with multiple recipients from go-sender integration test",
			WithSubject("Integration Test - Multiple Recipients"),
			WithFrom(globalTestConfig.From),
		)

		err := provider.Send(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("SendEmailWithDisplayNames", func(t *testing.T) {
		// Test with display names in email addresses
		fromWithName := "Test Sender <" + globalTestConfig.From + ">"
		toWithName := "Test Recipient <" + globalTestConfig.To + ">"

		msg := NewMessage(
			[]string{toWithName},
			"This is a test email with display names from go-sender integration test",
			WithSubject("Integration Test - Display Names"),
			WithFrom(fromWithName),
		)

		err := provider.Send(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("EmptyBody", func(t *testing.T) {
		msg := NewMessage(
			[]string{globalTestConfig.To},
			"",
			WithSubject("Test"),
		)

		err := provider.Send(context.Background(), msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email body cannot be empty")
	})
}

func TestEmailProvider_Validation(t *testing.T) {
	config := getTestConfig()
	provider, err := New(config)
	require.NoError(t, err)

	t.Run("InvalidMessageType", func(t *testing.T) {
		wrongMessage := &mockMessage{}
		err := provider.Send(context.Background(), wrongMessage)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid message type")
	})

	t.Run("EmptyRecipients", func(t *testing.T) {
		msg := NewMessage(
			[]string{},
			"Test body",
			WithSubject("Test"),
		)

		err := provider.Send(context.Background(), msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "recipient list cannot be empty")
	})

	t.Run("InvalidEmailAddress", func(t *testing.T) {
		msg := NewMessage(
			[]string{"invalid-email"},
			"Test body",
			WithSubject("Test"),
		)

		err := provider.Send(context.Background(), msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recipient email")
	})

	t.Run("EmptyBody", func(t *testing.T) {
		msg := NewMessage(
			[]string{"test@gmail.com"},
			"",
			WithSubject("Test"),
		)

		err := msg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email body cannot be empty")
	})
}

func TestEmailMessage_Validation(t *testing.T) {
	t.Run("ValidMessage", func(t *testing.T) {
		msg := NewMessage(
			[]string{"test@example.com"},
			"Test body",
			WithSubject("Test Subject"),
		)

		err := msg.Validate()
		assert.NoError(t, err)
	})

	t.Run("EmptyRecipients", func(t *testing.T) {
		msg := NewMessage(
			[]string{},
			"Test body",
		)

		err := msg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "recipient list cannot be empty")
	})

	t.Run("InvalidRecipientEmail", func(t *testing.T) {
		msg := NewMessage(
			[]string{"invalid-email"},
			"Test body",
		)

		err := msg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recipient email")
	})

	t.Run("InvalidCCEmail", func(t *testing.T) {
		msg := NewMessage(
			[]string{"test@example.com"},
			"Test body",
			WithCc("invalid-cc"),
		)

		err := msg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid CC email")
	})

	t.Run("InvalidBCCEmail", func(t *testing.T) {
		msg := NewMessage(
			[]string{"test@example.com"},
			"Test body",
			WithBcc("invalid-bcc"),
		)

		err := msg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid BCC email")
	})

	t.Run("EmptyBody", func(t *testing.T) {
		msg := NewMessage(
			[]string{"test@example.com"},
			"",
			WithSubject("Test"),
		)

		err := msg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email body cannot be empty")
	})

	t.Run("ValidEmailWithDisplayName", func(t *testing.T) {
		msg := NewMessage(
			[]string{"Test User <test@example.com>"},
			"Test body",
		)

		err := msg.Validate()
		assert.NoError(t, err)
	})
}

// mockMessage implements core.Message but is not *email.Message
type mockMessage struct {
	core.DefaultMessage
}

func (m *mockMessage) Validate() error                 { return nil }
func (m *mockMessage) ProviderType() core.ProviderType { return "mock" }
