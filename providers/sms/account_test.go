package sms_test

import (
	"encoding/json"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/sms"
)

func TestAccount_SignName(t *testing.T) {
	// 测试 Account 结构体包含 SignName 字段
	account := &sms.Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "sms",
				SubType:  "aliyun",
				Name:     "test-account",
			},
			Credentials: core.Credentials{
				APIKey:    "test-key",
				APISecret: "test-secret",
			},
		},
		Region:   "cn-hangzhou",
		Callback: "https://example.com/callback",
		SignName: "测试签名",
	}

	// 验证 SignName 字段
	if account.SignName != "测试签名" {
		t.Errorf("Expected SignName '测试签名', got %s", account.SignName)
	}

	// 验证其他字段
	if account.Region != "cn-hangzhou" {
		t.Errorf("Expected Region 'cn-hangzhou', got %s", account.Region)
	}
	if account.Callback != "https://example.com/callback" {
		t.Errorf("Expected Callback 'https://example.com/callback', got %s", account.Callback)
	}
}

func TestAccount_Validate_WithSignName(t *testing.T) {
	// 测试带 SignName 的 Account 验证
	account := &sms.Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "sms",
				SubType:  "aliyun",
				Name:     "test-account",
			},
			Credentials: core.Credentials{
				APIKey:    "test-key",
				APISecret: "test-secret",
			},
		},
		SignName: "测试签名",
	}

	err := account.Validate()
	if err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestAccount_JSON_Serialization(t *testing.T) {
	// 测试 JSON 序列化/反序列化
	account := &sms.Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: "sms",
				SubType:  "aliyun",
				Name:     "test-account",
			},
			Credentials: core.Credentials{
				APIKey:    "test-key",
				APISecret: "test-secret",
			},
		},
		Region:   "cn-hangzhou",
		Callback: "https://example.com/callback",
		SignName: "测试签名",
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	// 反序列化回来
	var newAccount sms.Account
	err = json.Unmarshal(jsonData, &newAccount)
	if err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	// 验证往返一致性
	if newAccount.SignName != account.SignName {
		t.Errorf("SignName mismatch: got %s, want %s", newAccount.SignName, account.SignName)
	}
	if newAccount.Region != account.Region {
		t.Errorf("Region mismatch: got %s, want %s", newAccount.Region, account.Region)
	}
	if newAccount.Callback != account.Callback {
		t.Errorf("Callback mismatch: got %s, want %s", newAccount.Callback, account.Callback)
	}
}
