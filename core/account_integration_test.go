package core_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/sms"
)

func TestAccount_JSON_Email_MetaFields(t *testing.T) {
	jsonConfig := `{
		"provider": "email",
		"subtype": "smtp",
		"name": "primary-smtp",
		"weight": 10,
		"disabled": false,
		"app_id": "",
		"api_key": "user@example.com",
		"api_secret": "password123",
		"host": "smtp.gmail.com",
		"port": 587,
		"from": "noreply@example.com"
	}`
	var emailAccount email.Account
	if err := json.Unmarshal([]byte(jsonConfig), &emailAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if emailAccount.Provider != "email" {
		t.Errorf("Expected provider 'email', got %s", emailAccount.Provider)
	}
	if emailAccount.SubType != "smtp" {
		t.Errorf("Expected subtype 'smtp', got %s", emailAccount.SubType)
	}
	if emailAccount.Name != "primary-smtp" {
		t.Errorf("Expected name 'primary-smtp', got %s", emailAccount.Name)
	}
	if emailAccount.Weight != 10 {
		t.Errorf("Expected weight 10, got %d", emailAccount.Weight)
	}
	if emailAccount.Disabled {
		t.Error("Expected disabled false, got true")
	}
}

func TestAccount_JSON_Email_CredentialsFields(t *testing.T) {
	jsonConfig := `{
		"provider": "email",
		"subtype": "smtp",
		"name": "primary-smtp",
		"weight": 10,
		"disabled": false,
		"app_id": "",
		"api_key": "user@example.com",
		"api_secret": "password123",
		"host": "smtp.gmail.com",
		"port": 587,
		"from": "noreply@example.com"
	}`
	var emailAccount email.Account
	if err := json.Unmarshal([]byte(jsonConfig), &emailAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if emailAccount.APIKey != "user@example.com" {
		t.Errorf("Expected APIKey 'user@example.com', got %s", emailAccount.APIKey)
	}
	if emailAccount.APISecret != "password123" {
		t.Errorf("Expected APISecret 'password123', got %s", emailAccount.APISecret)
	}
}

func TestAccount_JSON_Email_SpecificFields(t *testing.T) {
	jsonConfig := `{
		"provider": "email",
		"subtype": "smtp",
		"name": "primary-smtp",
		"weight": 10,
		"disabled": false,
		"app_id": "",
		"api_key": "user@example.com",
		"api_secret": "password123",
		"host": "smtp.gmail.com",
		"port": 587,
		"from": "noreply@example.com"
	}`
	var emailAccount email.Account
	if err := json.Unmarshal([]byte(jsonConfig), &emailAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if emailAccount.Host != "smtp.gmail.com" {
		t.Errorf("Expected host 'smtp.gmail.com', got %s", emailAccount.Host)
	}
	if emailAccount.Port != 587 {
		t.Errorf("Expected port 587, got %d", emailAccount.Port)
	}
	if emailAccount.From != "noreply@example.com" {
		t.Errorf("Expected from 'noreply@example.com', got %s", emailAccount.From)
	}
}

func TestAccount_JSON_Email_InterfaceImplementation(_ *testing.T) {
	var emailAccount email.Account
	var _ core.Selectable = &emailAccount
	var _ core.BasicAccount = &emailAccount
	var _ core.Validatable = &emailAccount
}

func TestAccount_JSON_Email_MethodCalls(t *testing.T) {
	jsonConfig := `{
		"provider": "email",
		"subtype": "smtp",
		"name": "primary-smtp",
		"weight": 10,
		"disabled": false,
		"app_id": "",
		"api_key": "user@example.com",
		"api_secret": "password123",
		"host": "smtp.gmail.com",
		"port": 587,
		"from": "noreply@example.com"
	}`
	var emailAccount email.Account
	if err := json.Unmarshal([]byte(jsonConfig), &emailAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	meta := emailAccount.GetMeta()
	if meta.Provider != "email" {
		t.Errorf("GetMeta().Provider expected 'email', got %s", meta.Provider)
	}
	creds := emailAccount.GetCredentials()
	if creds.APIKey != "user@example.com" {
		t.Errorf("GetCredentials().APIKey expected 'user@example.com', got %s", creds.APIKey)
	}
	if emailAccount.Username() != "user@example.com" {
		t.Errorf("Username() expected 'user@example.com', got %s", emailAccount.Username())
	}
	if emailAccount.Password() != "password123" {
		t.Errorf("Password() expected 'password123', got %s", emailAccount.Password())
	}
}

func TestAccount_JSON_SMS_MetaFields(t *testing.T) {
	jsonConfig := `{
		"provider": "sms",
		"subtype": "aliyun",
		"name": "primary-sms",
		"weight": 5,
		"disabled": false,
		"app_id": "LTAI5tRqF123456",
		"api_key": "your-access-key",
		"api_secret": "your-secret-key",
		"region": "cn-hangzhou",
		"callback": "https://example.com/sms/callback"
	}`
	var smsAccount sms.Account
	if err := json.Unmarshal([]byte(jsonConfig), &smsAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if smsAccount.Provider != "sms" {
		t.Errorf("Expected provider 'sms', got %s", smsAccount.Provider)
	}
	if smsAccount.SubType != "aliyun" {
		t.Errorf("Expected subtype 'aliyun', got %s", smsAccount.SubType)
	}
	if smsAccount.Name != "primary-sms" {
		t.Errorf("Expected name 'primary-sms', got %s", smsAccount.Name)
	}
	if smsAccount.Weight != 5 {
		t.Errorf("Expected weight 5, got %d", smsAccount.Weight)
	}
	if smsAccount.Disabled {
		t.Error("Expected disabled false, got true")
	}
}

func TestAccount_JSON_SMS_CredentialsFields(t *testing.T) {
	jsonConfig := `{
		"provider": "sms",
		"subtype": "aliyun",
		"name": "primary-sms",
		"weight": 5,
		"disabled": false,
		"app_id": "LTAI5tRqF123456",
		"api_key": "your-access-key",
		"api_secret": "your-secret-key",
		"region": "cn-hangzhou",
		"callback": "https://example.com/sms/callback"
	}`
	var smsAccount sms.Account
	if err := json.Unmarshal([]byte(jsonConfig), &smsAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if smsAccount.AppID != "LTAI5tRqF123456" {
		t.Errorf("Expected AppID 'LTAI5tRqF123456', got %s", smsAccount.AppID)
	}
	if smsAccount.APIKey != "your-access-key" {
		t.Errorf("Expected APIKey 'your-access-key', got %s", smsAccount.APIKey)
	}
	if smsAccount.APISecret != "your-secret-key" {
		t.Errorf("Expected APISecret 'your-secret-key', got %s", smsAccount.APISecret)
	}
}

func TestAccount_JSON_SMS_SpecificFields(t *testing.T) {
	jsonConfig := `{
		"provider": "sms",
		"subtype": "aliyun",
		"name": "primary-sms",
		"weight": 5,
		"disabled": false,
		"app_id": "LTAI5tRqF123456",
		"api_key": "your-access-key",
		"api_secret": "your-secret-key",
		"region": "cn-hangzhou",
		"callback": "https://example.com/sms/callback"
	}`
	var smsAccount sms.Account
	if err := json.Unmarshal([]byte(jsonConfig), &smsAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if smsAccount.Region != "cn-hangzhou" {
		t.Errorf("Expected region 'cn-hangzhou', got %s", smsAccount.Region)
	}
	if smsAccount.Callback != "https://example.com/sms/callback" {
		t.Errorf("Expected callback 'https://example.com/sms/callback', got %s", smsAccount.Callback)
	}
}

func TestAccount_JSON_SMS_InterfaceImplementation(_ *testing.T) {
	var smsAccount sms.Account
	var _ core.Selectable = &smsAccount
	var _ core.BasicAccount = &smsAccount
	var _ core.Validatable = &smsAccount
}

func TestAccount_JSON_SMS_ValidationLogic(t *testing.T) {
	jsonConfig := `{
		"provider": "sms",
		"subtype": "aliyun",
		"name": "primary-sms",
		"weight": 5,
		"disabled": false,
		"app_id": "LTAI5tRqF123456",
		"api_key": "your-access-key",
		"api_secret": "your-secret-key",
		"region": "cn-hangzhou",
		"callback": "https://example.com/sms/callback"
	}`
	var smsAccount sms.Account
	if err := json.Unmarshal([]byte(jsonConfig), &smsAccount); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if err := smsAccount.Validate(); err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
	smsAccount.SubType = ""
	if err := smsAccount.Validate(); err == nil {
		t.Error("Expected validation error for missing subType")
	}
}

func TestAccount_Serialization(t *testing.T) {
	// 创建Email Account实例
	emailAccount := email.Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "email",
				SubType:  "smtp",
				Name:     "test-smtp",
				Weight:   10,
				Disabled: false,
			},
			Credentials: core.Credentials{
				APIKey:    "test@example.com",
				APISecret: "test-password",
			},
		},
		Host: "smtp.example.com",
		Port: 587,
		From: "test@example.com",
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(emailAccount)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	// 反序列化回来
	var newEmailAccount email.Account
	err = json.Unmarshal(jsonData, &newEmailAccount)
	if err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	// 验证往返一致性
	if newEmailAccount.Provider != emailAccount.Provider {
		t.Errorf("Provider mismatch: got %s, want %s", newEmailAccount.Provider, emailAccount.Provider)
	}
	if newEmailAccount.Name != emailAccount.Name {
		t.Errorf("Name mismatch: got %s, want %s", newEmailAccount.Name, emailAccount.Name)
	}
	if newEmailAccount.Host != emailAccount.Host {
		t.Errorf("Host mismatch: got %s, want %s", newEmailAccount.Host, emailAccount.Host)
	}
}

func TestAccount_ConfigurationExample(t *testing.T) {
	// 演示完整的配置文件结构
	configExample := `{
		"providers": {
			"email": {
				"disabled": false,
				"strategy": "round-robin",
				"accounts": [
					{
						"provider": "email",
						"subtype": "smtp",
						"name": "gmail-primary",
						"weight": 10,
						"disabled": false,
						"api_key": "user@gmail.com",
						"api_secret": "app-password",
						"host": "smtp.gmail.com",
						"port": 587,
						"from": "noreply@gmail.com"
					},
					{
						"provider": "email",
						"subtype": "smtp",
						"name": "company-smtp",
						"weight": 5,
						"disabled": false,
						"api_key": "system@company.com",
						"api_secret": "company-password",
						"host": "smtp.company.com",
						"port": 465,
						"from": "system@company.com"
					}
				]
			},
			"sms": {
				"disabled": false,
				"strategy": "weighted",
				"accounts": [
					{
						"provider": "sms",
						"subtype": "aliyun",
						"name": "aliyun-primary",
						"weight": 8,
						"disabled": false,
						"app_id": "LTAI5tRqF123456",
						"api_key": "aliyun-access-key",
						"api_secret": "aliyun-secret-key",
						"region": "cn-hangzhou",
						"callback": "https://api.company.com/sms/callback"
					},
					{
						"provider": "sms",
						"subtype": "tencent",
						"name": "tencent-backup",
						"weight": 4,
						"disabled": false,
						"app_id": "1400000000",
						"api_key": "tencent-secret-id",
						"api_secret": "tencent-secret-key",
						"region": "ap-guangzhou",
						"callback": "https://api.company.com/sms/callback"
					}
				]
			}
		}
	}`

	// 解析配置
	var config map[string]interface{}
	err := json.Unmarshal([]byte(configExample), &config)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	// 验证配置结构
	providers, ok := config["providers"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected providers object")
	}

	// 验证Email配置
	emailConfig, ok := providers["email"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected email config")
	}

	emailAccounts, ok := emailConfig["accounts"].([]interface{})
	if !ok {
		t.Fatal("Expected email accounts array")
	}

	if len(emailAccounts) != 2 {
		t.Errorf("Expected 2 email accounts, got %d", len(emailAccounts))
	}

	// 验证SMS配置
	smsConfig, ok := providers["sms"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected sms config")
	}

	smsAccounts, ok := smsConfig["accounts"].([]interface{})
	if !ok {
		t.Fatal("Expected sms accounts array")
	}

	if len(smsAccounts) != 2 {
		t.Errorf("Expected 2 sms accounts, got %d", len(smsAccounts))
	}
}

//nolint:gocyclo,gocognit,cyclop  // test function naturally complex with many subtests and branches
func TestAccountConfigDynamicUpdate(t *testing.T) {
	t.Run("Basic CRUD Operations", func(t *testing.T) {
		config := &core.BaseConfig[*core.BaseAccount]{}
		// 创建初始账号
		account1 := &core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "test-provider",
				Name:     "test1",
				Weight:   100,
			},
			Credentials: core.Credentials{
				APIKey:    "key1",
				APISecret: "secret1",
			},
		}
		account2 := &core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "test-provider",
				Name:     "test2",
				Weight:   50,
			},
			Credentials: core.Credentials{
				APIKey:    "key2",
				APISecret: "secret2",
			},
		}

		// 添加账号
		if err := config.Add(account1); err != nil {
			t.Errorf("Failed to add account1: %v", err)
		}
		if err := config.Add(account2); err != nil {
			t.Errorf("Failed to add account2: %v", err)
		}

		// 验证账号数量
		if len(config.GetItems()) != 2 {
			t.Errorf("Expected 2 accounts, got %d", len(config.GetItems()))
		}

		// 更新账号
		account1.Weight = 80
		if err := config.Update(account1); err != nil {
			t.Errorf("Failed to update account: %v", err)
		}

		// 验证更新结果
		items := config.GetItems()
		var found bool
		for _, item := range items {
			if item.Name == "test1" && item.Weight == 80 {
				found = true
				break
			}
		}
		if !found {
			t.Error("Account update not reflected")
		}

		// 删除账号
		config.Delete("test2")

		// 验证删除结果
		if len(config.GetItems()) != 1 {
			t.Error("Account deletion not reflected")
		}
	})

	t.Run("Concurrent Operations", func(_ *testing.T) {
		config := &core.BaseConfig[*core.BaseAccount]{}
		var wg sync.WaitGroup
		accountCount := 10
		wg.Add(accountCount * 2) // 一半goroutine添加，一半删除

		// 并发添加账号
		for i := range accountCount {
			go func(idx int) {
				defer wg.Done()
				account := &core.BaseAccount{
					AccountMeta: core.AccountMeta{
						Provider: "test-provider",
						Name:     fmt.Sprintf("concurrent-%d", idx),
						Weight:   100,
					},
					Credentials: core.Credentials{
						APIKey:    fmt.Sprintf("key-%d", idx),
						APISecret: fmt.Sprintf("secret-%d", idx),
					},
				}
				_ = config.Add(account)
			}(i)
		}

		// 并发删除账号
		for i := range accountCount {
			go func(idx int) {
				defer wg.Done()
				time.Sleep(time.Millisecond * 10) // 稍微延迟删除操作
				config.Delete(fmt.Sprintf("concurrent-%d", idx))
			}(i)
		}

		wg.Wait()
	})

	t.Run("Weight Adjustment", func(t *testing.T) {
		config := &core.BaseConfig[*core.BaseAccount]{}
		// 添加多个账号
		accounts := []*core.BaseAccount{
			{
				AccountMeta: core.AccountMeta{
					Provider: "test-provider",
					Name:     "high",
					Weight:   100,
				},
				Credentials: core.Credentials{APIKey: "k"},
			},
			{
				AccountMeta: core.AccountMeta{
					Provider: "test-provider",
					Name:     "medium",
					Weight:   50,
				},
				Credentials: core.Credentials{APIKey: "k"},
			},
			{
				AccountMeta: core.AccountMeta{
					Provider: "test-provider",
					Name:     "low",
					Weight:   25,
				},
				Credentials: core.Credentials{APIKey: "k"},
			},
		}
		for _, acc := range accounts {
			if err := config.Add(acc); err != nil {
				t.Errorf("Failed to add account %s: %v", acc.Name, err)
			}
		}

		// 动态调整权重
		accounts[1].Weight = 75
		if err := config.Update(accounts[1]); err != nil {
			t.Errorf("Failed to update account weight: %v", err)
		}

		// 验证选择策略是否反映新权重
		ctx := core.WithRoute(context.Background(), &core.RouteInfo{StrategyType: core.StrategyWeighted})

		counts := make(map[string]int)
		iterations := 1000

		for range iterations {
			selected, err := config.Select(ctx, nil)
			if err != nil {
				t.Errorf("Failed to select account: %v", err)
				continue
			}
			counts[selected.Name]++
		}

		// 验证高权重账号被选择的次数更多
		if counts["high"] <= counts["low"] {
			t.Error("Weight-based selection not working as expected")
		}
	})

	t.Run("Account Status Toggle", func(t *testing.T) {
		config := &core.BaseConfig[*core.BaseAccount]{}
		account := &core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "test-provider",
				Name:     "status-test",
				Weight:   100,
			},
			Credentials: core.Credentials{
				APIKey:    "key",
				APISecret: "secret",
			},
		}

		// 添加账号
		if err := config.Add(account); err != nil {
			t.Errorf("Failed to add account: %v", err)
		}

		// 禁用账号
		account.Disabled = true
		if err := config.Update(account); err != nil {
			t.Errorf("Failed to disable account: %v", err)
		}

		// 验证禁用状态
		items := config.GetItems()
		var found bool
		for _, item := range items {
			if item.Name == "status-test" && item.Disabled {
				found = true
				break
			}
		}
		if !found {
			t.Error("Account status update not reflected")
		}

		// 验证禁用账号不会被选择策略选中
		ctx := core.WithRoute(context.Background(), &core.RouteInfo{StrategyType: core.StrategyWeighted})
		_, err := config.Select(ctx, nil)
		if err == nil {
			t.Error("Expected error when selecting disabled account")
		}
	})

	t.Run("Atomic Configuration Update", func(t *testing.T) {
		config := &core.BaseConfig[*core.BaseAccount]{}
		// 创建一组初始账号
		initialAccounts := []*core.BaseAccount{
			{
				AccountMeta: core.AccountMeta{
					Provider: "test-provider",
					Name:     "atomic1",
					Weight:   100,
				},
				Credentials: core.Credentials{APIKey: "k"},
			},
			{
				AccountMeta: core.AccountMeta{
					Provider: "test-provider",
					Name:     "atomic2",
					Weight:   100,
				},
				Credentials: core.Credentials{APIKey: "k"},
			},
		}
		for _, acc := range initialAccounts {
			if err := config.Add(acc); err != nil {
				t.Errorf("Failed to add account %s: %v", acc.Name, err)
			}
		}

		// 并发更新和读取
		var wg sync.WaitGroup
		updateCount := 100
		wg.Add(updateCount + 1) // +1 for the reader goroutine

		// 启动一个持续读取的 goroutine
		go func() {
			defer wg.Done()
			for range updateCount {
				items := config.GetItems()
				// 验证账号列表的一致性
				for _, acc := range items {
					if acc.Weight < 0 || acc.Weight > 200 {
						t.Errorf("Invalid weight detected during concurrent update: %d", acc.Weight)
					}
				}
				time.Sleep(time.Millisecond)
			}
		}()

		// 并发更新权重
		for i := range updateCount {
			go func(idx int) {
				defer wg.Done()
				weight := 50 + idx%150 // 保持在合理范围内
				acc := &core.BaseAccount{
					AccountMeta: core.AccountMeta{
						Provider: "test-provider",
						Name:     "atomic1",
						Weight:   weight,
					},
					Credentials: core.Credentials{APIKey: "k"},
				}
				_ = config.Update(acc)
			}(i)
		}

		wg.Wait()
	})
}
