package core_test

import (
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestBaseAccount_GetMeta(t *testing.T) {
	account := &core.BaseAccount{
		AccountMeta: core.AccountMeta{
			Provider: "email",
			SubType:  "mailgun",
			Name:     "test-account",
			Weight:   10,
			Disabled: false,
		},
	}

	meta := account.GetMeta()
	if meta.Provider != "email" {
		t.Errorf("Expected provider 'email', got %s", meta.Provider)
	}
	if meta.SubType != "mailgun" {
		t.Errorf("Expected subtype 'mailgun', got %s", meta.SubType)
	}
	if meta.Name != "test-account" {
		t.Errorf("Expected name 'test-account', got %s", meta.Name)
	}
	if meta.Weight != 10 {
		t.Errorf("Expected weight 10, got %d", meta.Weight)
	}
	if meta.Disabled {
		t.Error("Expected disabled false, got true")
	}
}

func TestBaseAccount_GetCredentials(t *testing.T) {
	account := &core.BaseAccount{
		AccountMeta: core.AccountMeta{
			Name: "test-account",
		},
		Credentials: core.Credentials{
			AppID:     "test-app",
			APIKey:    "test-key",
			APISecret: "test-secret",
		},
	}

	creds := account.GetCredentials()
	if creds.AppID != "test-app" {
		t.Errorf("Expected AppID 'test-app', got %s", creds.AppID)
	}
	if creds.APIKey != "test-key" {
		t.Errorf("Expected APIKey 'test-key', got %s", creds.APIKey)
	}
	if creds.APISecret != "test-secret" {
		t.Errorf("Expected APISecret 'test-secret', got %s", creds.APISecret)
	}
}

func TestBaseAccount_Validate(t *testing.T) {
	tests := []struct {
		name    string
		account *core.BaseAccount
		wantErr bool
	}{
		{
			name: "valid account",
			account: &core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Provider: "email",
					Name:     "test-account",
				},
				Credentials: core.Credentials{
					APIKey: "test-key",
				},
			},
			wantErr: false,
		},
		{
			name: "empty provider",
			account: &core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Provider: "",
					Name:     "test-account",
				},
				Credentials: core.Credentials{
					APIKey: "test-key",
				},
			},
			wantErr: true,
		},
		{
			name: "empty name",
			account: &core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Provider: "email",
					Name:     "",
				},
				Credentials: core.Credentials{
					APIKey: "test-key",
				},
			},
			wantErr: true,
		},
		{
			name: "empty credentials",
			account: &core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Provider: "email",
					Name:     "test-account",
				},
				Credentials: core.Credentials{},
			},
			wantErr: true,
		},
		{
			name: "valid with AppID only",
			account: &core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Provider: "email",
					Name:     "test-account",
				},
				Credentials: core.Credentials{
					AppID: "test-app",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccountMeta_GetName(t *testing.T) {
	meta := &core.AccountMeta{
		Name: "test-account",
	}

	name := meta.GetName()
	if name != "test-account" {
		t.Errorf("Expected name 'test-account', got %s", name)
	}
}

func TestAccountMeta_GetWeight(t *testing.T) {
	tests := []struct {
		name     string
		meta     *core.AccountMeta
		expected int
	}{
		{
			name: "positive weight",
			meta: &core.AccountMeta{
				Weight: 15,
			},
			expected: 15,
		},
		{
			name: "zero weight",
			meta: &core.AccountMeta{
				Weight: 0,
			},
			expected: 1,
		},
		{
			name: "negative weight",
			meta: &core.AccountMeta{
				Weight: -5,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weight := tt.meta.GetWeight()
			if weight != tt.expected {
				t.Errorf("GetWeight() = %d, want %d", weight, tt.expected)
			}
		})
	}
}

func TestAccountMeta_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		meta     *core.AccountMeta
		expected bool
	}{
		{
			name: "enabled account",
			meta: &core.AccountMeta{
				Disabled: false,
			},
			expected: true,
		},
		{
			name: "disabled account",
			meta: &core.AccountMeta{
				Disabled: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled := tt.meta.IsEnabled()
			if enabled != tt.expected {
				t.Errorf("IsEnabled() = %v, want %v", enabled, tt.expected)
			}
		})
	}
}

func TestAccountMeta_GetType(t *testing.T) {
	meta := &core.AccountMeta{
		SubType: "mailgun",
	}

	accountType := meta.GetType()
	if accountType != "mailgun" {
		t.Errorf("Expected type 'mailgun', got %s", accountType)
	}
}

func TestProviderMeta_GetStrategy(t *testing.T) {
	tests := []struct {
		name     string
		meta     *core.ProviderMeta
		expected core.StrategyType
	}{
		{
			name: "with strategy",
			meta: &core.ProviderMeta{
				Strategy: core.StrategyRandom,
			},
			expected: core.StrategyRandom,
		},
		{
			name: "empty strategy",
			meta: &core.ProviderMeta{
				Strategy: "",
			},
			expected: core.StrategyRoundRobin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := tt.meta.GetStrategy()
			if strategy != tt.expected {
				t.Errorf("GetStrategy() = %v, want %v", strategy, tt.expected)
			}
		})
	}
}
