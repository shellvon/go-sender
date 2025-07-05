package email_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/sms"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name    string
		config  *email.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &email.Config{
				ProviderMeta: core.ProviderMeta{},
				Items: []*email.Account{
					{
						BaseAccount: core.BaseAccount{
							AccountMeta: core.AccountMeta{
								Provider: "email",
								Name:     "test",
							},
							Credentials: core.Credentials{
								APIKey:    "user",
								APISecret: "pass",
							},
						},
						Host: "smtp.example.com",
						Port: 587,
						From: "test@example.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "disabled config",
			config: &email.Config{
				ProviderMeta: core.ProviderMeta{Disabled: true},
				Items: []*email.Account{
					{
						BaseAccount: core.BaseAccount{
							AccountMeta: core.AccountMeta{
								Provider: "email",
								Name:     "test",
							},
							Credentials: core.Credentials{
								APIKey:    "user",
								APISecret: "pass",
							},
						},
						Host: "smtp.example.com",
						Port: 587,
						From: "test@example.com",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no accounts",
			config: &email.Config{
				ProviderMeta: core.ProviderMeta{},
				Items:        []*email.Account{},
			},
			wantErr: true,
		},
		{
			name: "all accounts disabled",
			config: &email.Config{
				ProviderMeta: core.ProviderMeta{},
				Items: []*email.Account{
					{
						BaseAccount: core.BaseAccount{
							AccountMeta: core.AccountMeta{
								Provider: "email",
								Name:     "test1",
								Disabled: true,
							},
							Credentials: core.Credentials{
								APIKey:    "user",
								APISecret: "pass",
							},
						},
						Host: "smtp.example.com",
						Port: 587,
						From: "test@example.com",
					},
					{
						BaseAccount: core.BaseAccount{
							AccountMeta: core.AccountMeta{
								Provider: "email",
								Name:     "test2",
								Disabled: true,
							},
							Credentials: core.Credentials{
								APIKey:    "user",
								APISecret: "pass",
							},
						},
						Host: "smtp.example.com",
						Port: 587,
						From: "test@example.com",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := email.New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && provider == nil {
				t.Error("New() returned nil provider when no error expected")
			}
		})
	}
}

func TestProviderName(t *testing.T) {
	config := email.Config{
		ProviderMeta: core.ProviderMeta{},
		Items: []*email.Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{
						Provider: "email",
						Name:     "test",
					},
					Credentials: core.Credentials{
						APIKey:    "user",
						APISecret: "pass",
					},
				},
				Host: "smtp.example.com",
				Port: 587,
				From: "test@example.com",
			},
		},
	}

	provider, err := email.New(&config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	name := provider.Name()
	if name != "email" {
		t.Errorf("Expected provider name to be 'email', got %s", name)
	}
}

func TestProviderSendInvalidMessageType(t *testing.T) {
	config := email.Config{
		ProviderMeta: core.ProviderMeta{},
		Items: []*email.Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{
						Provider: "email",
						Name:     "test",
					},
					Credentials: core.Credentials{
						APIKey:    "user",
						APISecret: "pass",
					},
				},
				Host: "smtp.example.com",
				Port: 587,
				From: "test@example.com",
			},
		},
	}

	provider, err := email.New(&config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Create a different message type - use SMS message instead of email message
	// This will cause a type assertion error in the email provider
	invalidMsg := &sms.Message{}

	err = provider.Send(context.Background(), invalidMsg, nil)
	if err == nil {
		t.Error("Expected error for invalid message type, got nil")
	}
}

func TestProviderSendInvalidMessage(t *testing.T) {
	config := email.Config{
		ProviderMeta: core.ProviderMeta{},
		Items: []*email.Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{
						Provider: "email",
						Name:     "test",
					},
					Credentials: core.Credentials{
						APIKey:    "user",
						APISecret: "pass",
					},
				},
				Host: "smtp.example.com",
				Port: 587,
				From: "test@example.com",
			},
		},
	}

	provider, err := email.New(&config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Create invalid message (empty recipients)
	invalidMsg := email.NewMessage([]string{}, "Test body")

	err = provider.Send(context.Background(), invalidMsg, nil)
	if err == nil {
		t.Error("Expected error for invalid message, got nil")
	}
}

func TestAccountMethods(t *testing.T) {
	account := &email.Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "email",
				Name:     "test",
				SubType:  "smtp",
				Weight:   10,
				Disabled: false,
			},
			Credentials: core.Credentials{
				APIKey:    "user",
				APISecret: "pass",
			},
		},
		Host: "smtp.example.com",
		Port: 587,
		From: "test@example.com",
	}

	// Test GetName
	if name := account.GetName(); name != "test" {
		t.Errorf("Expected GetName() to return 'test', got %s", name)
	}

	// Test GetWeight
	if weight := account.GetWeight(); weight != 10 {
		t.Errorf("Expected GetWeight() to return 10, got %d", weight)
	}

	// Test IsEnabled
	if enabled := account.IsEnabled(); !enabled {
		t.Error("Expected IsEnabled() to return true")
	}

	// Test GetType
	if accountType := account.GetType(); accountType != "smtp" {
		t.Errorf("Expected GetType() to return 'smtp', got %s", accountType)
	}

	// Test default weight
	account.AccountMeta.Weight = 0
	if weight := account.GetWeight(); weight != 1 {
		t.Errorf("Expected GetWeight() to return 1 for zero weight, got %d", weight)
	}

	// Test disabled account
	account.AccountMeta.Disabled = true
	if enabled := account.IsEnabled(); enabled {
		t.Error("Expected IsEnabled() to return false for disabled account")
	}
}

func TestConfigMethods(t *testing.T) {
	config := email.Config{
		ProviderMeta: core.ProviderMeta{},
		Items: []*email.Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{
						Provider: "email",
						Name:     "test",
					},
					Credentials: core.Credentials{
						APIKey:    "user",
						APISecret: "pass",
					},
				},
				Host: "smtp.example.com",
				Port: 587,
				From: "test@example.com",
			},
		},
	}

	// Test GetStrategy default
	if strategy := config.GetStrategy(); strategy != core.StrategyRoundRobin {
		t.Errorf("Expected GetStrategy() to return StrategyRoundRobin, got %s", strategy)
	}

	// Test GetStrategy custom
	config.Strategy = core.StrategyRandom
	if strategy := config.GetStrategy(); strategy != core.StrategyRandom {
		t.Errorf("Expected GetStrategy() to return StrategyRandom, got %s", strategy)
	}
}
