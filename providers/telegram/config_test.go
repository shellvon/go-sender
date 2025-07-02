package telegram_test

import (
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/telegram"
)

func TestConfig_IsConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   telegram.Config
		expected bool
	}{
		{
			name: "valid config",
			config: telegram.Config{
				Accounts: []core.Account{
					{Name: "test", Key: "bot123:token"},
				},
			},
			expected: true,
		},
		{
			name:     "empty config",
			config:   telegram.Config{},
			expected: false,
		},
		{
			name: "disabled config",
			config: telegram.Config{
				BaseConfig: core.BaseConfig{Disabled: true},
				Accounts: []core.Account{
					{Name: "test", Key: "bot123:token"},
				},
			},
			expected: false,
		},
		{
			name: "no accounts",
			config: telegram.Config{
				Accounts: []core.Account{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsConfigured()
			if result != tt.expected {
				t.Errorf("Expected IsConfigured() to be %v, got %v", tt.expected, result)
			}
		})
	}
}
