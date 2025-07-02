package email_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
)

func TestEmailIntegration(t *testing.T) {
	// 创建 email provider 配置
	emailConfig := email.Config{
		BaseConfig: core.BaseConfig{},
		Accounts: []email.Account{
			{
				Name:     "test",
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password",
				From:     "test@example.com",
				Weight:   1,
			},
		},
	}

	// 创建 email provider
	emailProvider, err := email.New(emailConfig)
	if err != nil {
		t.Fatalf("Failed to create email provider: %v", err)
	}

	// 创建 Sender 实例
	s := gosender.NewSender()
	defer s.Close()

	// 注册 email provider
	s.RegisterProvider(core.ProviderTypeEmail, emailProvider, nil)

	// 创建测试消息
	msg := email.NewMessage(
		[]string{"recipient@example.com"},
		"Test email body",
		email.WithSubject("Test Subject"),
		email.WithFrom("sender@example.com"),
		email.WithCc("cc@example.com"),
		email.WithBcc("bcc@example.com"),
		email.WithReplyTo("reply1@example.com"),
		email.WithHTML(),
		email.WithAttachments("test.txt"),
	)

	// 验证消息
	if validateErr := msg.Validate(); validateErr != nil {
		t.Fatalf("Message validation failed: %v", validateErr)
	}

	// 测试通过 Sender 发送（这会失败，因为我们没有真实的 SMTP 服务器）
	// 但我们可以验证消息格式和 provider 选择逻辑
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.Send(ctx, msg)
	// 我们期望这里会失败，因为没有真实的 SMTP 服务器
	// 但错误应该是网络相关的，而不是消息格式相关的
	if err != nil {
		// 检查错误类型，确保不是消息格式错误
		var paramErr *core.ParamError
		if errors.As(err, &paramErr) {
			t.Errorf("Expected network/connection error, got parameter error: %v", err)
		}
		// 其他类型的错误（如网络错误）是预期的
		t.Logf("Expected error (no SMTP server): %v", err)
	}

	if msg.ReplyTo != "reply1@example.com" {
		t.Errorf("Expected ReplyTo to be 'reply1@example.com', got %s", msg.ReplyTo)
	}
}

func TestEmailProviderSelection(t *testing.T) {
	// 测试多个账号的选择逻辑
	emailConfig := email.Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []email.Account{
			{
				Name:     "account1",
				Host:     "smtp1.example.com",
				Port:     587,
				Username: "user1@example.com",
				Password: "pass1",
				From:     "user1@example.com",
				Weight:   1,
			},
			{
				Name:     "account2",
				Host:     "smtp2.example.com",
				Port:     587,
				Username: "user2@example.com",
				Password: "pass2",
				From:     "user2@example.com",
				Weight:   2,
			},
		},
	}

	// 创建 email provider
	emailProvider, err := email.New(emailConfig)
	if err != nil {
		t.Fatalf("Failed to create email provider: %v", err)
	}

	// 验证 provider 名称
	if emailProvider.Name() != "email" {
		t.Errorf("Expected provider name 'email', got %s", emailProvider.Name())
	}

	// 创建测试消息
	msg := email.NewMessage(
		[]string{"test@example.com"},
		"Test body",
	)

	// 验证消息类型
	if msg.ProviderType() != core.ProviderTypeEmail {
		t.Errorf("Expected provider type %s, got %s", core.ProviderTypeEmail, msg.ProviderType())
	}
}

func TestEmailMessageOptions(t *testing.T) {
	// 测试所有消息选项的组合
	msg := email.NewMessage(
		[]string{"recipient1@example.com", "recipient2@example.com"},
		"<h1>HTML Test Body</h1>",
		email.WithFrom("Sender Name <sender@example.com>"),
		email.WithSubject("Test Email Subject"),
		email.WithCc("cc1@example.com", "cc2@example.com"),
		email.WithBcc("bcc1@example.com", "bcc2@example.com"),
		email.WithReplyTo("Reply Name <reply@example.com>"),
		email.WithHTML(),
		email.WithAttachments("file1.txt", "file2.pdf", "file3.jpg"),
	)

	// 验证所有字段
	if len(msg.To) != 2 {
		t.Errorf("Expected 2 recipients, got %d", len(msg.To))
	}

	if msg.From != "Sender Name <sender@example.com>" {
		t.Errorf("Expected From 'Sender Name <sender@example.com>', got %s", msg.From)
	}

	if msg.Subject != "Test Email Subject" {
		t.Errorf("Expected Subject 'Test Email Subject', got %s", msg.Subject)
	}

	if len(msg.Cc) != 2 {
		t.Errorf("Expected 2 CC recipients, got %d", len(msg.Cc))
	}

	if len(msg.Bcc) != 2 {
		t.Errorf("Expected 2 BCC recipients, got %d", len(msg.Bcc))
	}

	if msg.ReplyTo != "Reply Name <reply@example.com>" {
		t.Errorf("Expected ReplyTo to be 'Reply Name <reply@example.com>', got %s", msg.ReplyTo)
	}

	if !msg.IsHTML {
		t.Error("Expected IsHTML to be true")
	}

	if len(msg.Attachments) != 3 {
		t.Errorf("Expected 3 attachments, got %d", len(msg.Attachments))
	}

	// 验证消息 ID
	msgID := msg.MsgID()
	if msgID == "" {
		t.Error("Expected non-empty MsgID")
	}

	// 验证消息
	if validateErr := msg.Validate(); validateErr != nil {
		t.Fatalf("Message validation failed: %v", validateErr)
	}
}

func TestEmailConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  email.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: email.Config{
				BaseConfig: core.BaseConfig{},
				Accounts: []email.Account{
					{
						Name:     "test",
						Host:     "smtp.example.com",
						Port:     587,
						Username: "user@example.com",
						Password: "password",
						From:     "user@example.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "disabled config",
			config: email.Config{
				BaseConfig: core.BaseConfig{Disabled: true},
				Accounts: []email.Account{
					{
						Name:     "test",
						Host:     "smtp.example.com",
						Port:     587,
						Username: "user@example.com",
						Password: "password",
						From:     "user@example.com",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no accounts",
			config: email.Config{
				BaseConfig: core.BaseConfig{},
				Accounts:   []email.Account{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := email.New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("email.New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
