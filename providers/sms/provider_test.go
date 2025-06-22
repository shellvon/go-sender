package sms

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Providers: []SMSProvider{
					{
						Name:      "test",
						Type:      ProviderTypeTencent,
						AppID:     "test_app_id",
						AppSecret: "test_app_secret",
						Weight:    100,
						Disabled:  false,
					},
				},
				Strategy: core.StrategyWeighted,
			},
			wantErr: false,
		},
		{
			name: "empty providers",
			config: Config{
				Providers: []SMSProvider{},
				Strategy:  core.StrategyWeighted,
			},
			wantErr: true,
		},
		{
			name: "disabled config",
			config: Config{
				BaseConfig: core.BaseConfig{
					Disabled: true,
					Strategy: core.StrategyWeighted,
				},
				Providers: []SMSProvider{
					{
						Name:      "test",
						Type:      ProviderTypeTencent,
						AppID:     "test_app_id",
						AppSecret: "test_app_secret",
						Weight:    100,
						Disabled:  false,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
			}
		})
	}
}

func TestProvider_Send(t *testing.T) {
	config := Config{
		Providers: []SMSProvider{
			{
				Name:      "tencent",
				Type:      ProviderTypeTencent,
				AppID:     "test_app_id",
				AppSecret: "test_app_secret",
				Weight:    100,
				Disabled:  false,
			},
			{
				Name:      "aliyun",
				Type:      ProviderTypeAliyun,
				AppID:     "test_aliyun_key",
				AppSecret: "test_aliyun_secret",
				Weight:    80,
				Disabled:  false,
			},
		},
		Strategy: core.StrategyWeighted,
	}

	provider, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, provider)

	tests := []struct {
		name    string
		message *Message
		wantErr bool
	}{
		{
			name: "valid message",
			message: &Message{
				Mobile:  "***REMOVED***",
				Content: "test message",
			},
			wantErr: false, // 会失败因为配置是假的，但不会因为消息格式错误
		},
		{
			name: "empty mobile",
			message: &Message{
				Mobile:  "",
				Content: "test message",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			message: &Message{
				Mobile:  "***REMOVED***",
				Content: "",
			},
			wantErr: true,
		},
		{
			name: "with template",
			message: &Message{
				Mobile:       "***REMOVED***",
				TemplateCode: "SMS_123456",
				TemplateParams: map[string]string{
					"code": "123456",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := provider.Send(ctx, tt.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// 即使发送失败，也不应该是消息格式错误
				if err != nil {
					assert.NotContains(t, err.Error(), "mobile")
					assert.NotContains(t, err.Error(), "content")
				}
			}
		})
	}
}

func TestProvider_SelectProvider(t *testing.T) {
	config := Config{
		Providers: []SMSProvider{
			{
				Name:      "tencent",
				Type:      ProviderTypeTencent,
				AppID:     "test_app_id",
				AppSecret: "test_app_secret",
				Weight:    100,
				Disabled:  false,
			},
			{
				Name:      "aliyun",
				Type:      ProviderTypeAliyun,
				AppID:     "test_aliyun_key",
				AppSecret: "test_aliyun_secret",
				Weight:    80,
				Disabled:  false,
			},
			{
				Name:      "disabled",
				Type:      ProviderTypeYunpian,
				AppID:     "test_yunpian_key",
				AppSecret: "test_yunpian_secret",
				Weight:    50,
				Disabled:  true,
			},
		},
		Strategy: core.StrategyWeighted,
	}

	provider, err := New(config)
	require.NoError(t, err)

	tests := []struct {
		name          string
		providerName  string
		expectedNames []string
		expectNil     bool
	}{
		{
			name:          "select by name",
			providerName:  "tencent",
			expectedNames: []string{"tencent"},
			expectNil:     false,
		},
		{
			name:          "select by name aliyun",
			providerName:  "aliyun",
			expectedNames: []string{"aliyun"},
			expectNil:     false,
		},
		{
			name:          "select by strategy weighted",
			providerName:  "",
			expectedNames: []string{"tencent", "aliyun"},
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.providerName != "" {
				ctx = core.WithCtxItemName(ctx, tt.providerName)
			}
			selected := provider.selectProvider(ctx)
			if tt.expectNil {
				assert.Nil(t, selected)
			} else {
				assert.NotNil(t, selected)
				assert.Contains(t, tt.expectedNames, selected.Name)
			}
		})
	}
}

func TestProvider_HashFunctions(t *testing.T) {
	provider := &Provider{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "md5 hash",
			input:    "test",
			expected: "098f6bcd4621d373cade4e832627b4f6", // MD5 of "test"
		},
		{
			name:     "sha1 hash",
			input:    "test",
			expected: "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", // SHA1 of "test"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.name == "md5 hash" {
				result = provider.md5Hash(tt.input)
			} else {
				result = provider.sha1Hash(tt.input)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
		wantErr bool
	}{
		{
			name: "valid message",
			message: &Message{
				Mobile:  "***REMOVED***",
				Content: "test message",
			},
			wantErr: false,
		},
		{
			name: "valid template message",
			message: &Message{
				Mobile:       "***REMOVED***",
				TemplateCode: "SMS_123456",
				TemplateParams: map[string]string{
					"code": "123456",
				},
			},
			wantErr: false,
		},
		{
			name: "empty mobile",
			message: &Message{
				Mobile:  "",
				Content: "test message",
			},
			wantErr: true,
		},
		{
			name: "empty content and template",
			message: &Message{
				Mobile: "***REMOVED***",
			},
			wantErr: true,
		},
		{
			name: "template code without params",
			message: &Message{
				Mobile:       "***REMOVED***",
				TemplateCode: "SMS_123456",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMessage_ProviderType(t *testing.T) {
	message := &Message{
		Mobile:  "***REMOVED***",
		Content: "test message",
	}

	assert.Equal(t, core.ProviderTypeSMS, message.ProviderType())
}

func TestMessageOptions(t *testing.T) {
	message := &Message{}

	// Test WithMobile
	WithMobile("***REMOVED***")(message)
	assert.Equal(t, "***REMOVED***", message.Mobile)

	// Test WithContent
	WithContent("test content")(message)
	assert.Equal(t, "test content", message.Content)

	// Test WithTemplateCode
	WithTemplateCode("SMS_123456")(message)
	assert.Equal(t, "SMS_123456", message.TemplateCode)

	// Test WithTemplateParams
	params := map[string]string{"code": "123456"}
	WithTemplateParams(params)(message)
	assert.Equal(t, params, message.TemplateParams)

}

func TestNewMessage(t *testing.T) {
	message := NewMessage("***REMOVED***", WithContent("test content"),
		WithTemplateCode("SMS_123456"),
		WithTemplateParams(map[string]string{"code": "123456"}),
	)

	assert.Equal(t, "***REMOVED***", message.Mobile)
	assert.Equal(t, "test content", message.Content)
	assert.Equal(t, "SMS_123456", message.TemplateCode)
	assert.Equal(t, map[string]string{"code": "123456"}, message.TemplateParams)
}
