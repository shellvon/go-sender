package serverchan

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
)

// TestServerChanProviderIntegration 测试 ServerChan 提供者的集成功能
// 这个测试需要设置环境变量 SERVERCHAN_KEY 来运行
// 可以通过以下方式设置环境变量：
//
//	export SERVERCHAN_KEY="your_sendkey_here"
//	或者在运行测试时设置：
//	SERVERCHAN_KEY="your_sendkey_here" go test -v -run TestServerChanProviderIntegration
func TestServerChanProviderIntegration(t *testing.T) {
	// 检查环境变量
	sendKey := os.Getenv("SERVERCHAN_KEY")
	if sendKey == "" {
		t.Skip("SERVERCHAN_KEY environment variable not set, skipping integration test")
	}

	// 创建 ServerChan 配置
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "main",
				Key:      sendKey,
				Weight:   100,
				Disabled: false,
			},
		},
	}

	// 创建 ServerChan provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create ServerChan provider: %v", err)
	}

	// 创建发送器
	sender := gosender.NewSender(nil) // 使用默认 logger
	sender.RegisterProvider(core.ProviderTypeServerChan, provider, nil)

	// 创建上下文
	ctx := context.Background()

	fmt.Println("🚀 开始测试 ServerChan 发送功能...")
	fmt.Println(strings.Repeat("=", 50))

	// 测试用例：基础文本消息
	t.Run("Basic Text Message", func(t *testing.T) {
		msg := NewMessage(
			"测试消息",
			"这是一条来自 go-sender 的测试消息\n\n时间: "+time.Now().Format("2006-01-02 15:04:05"),
		)

		err := sender.Send(ctx, msg)
		if err != nil {
			t.Errorf("Failed to send basic message: %v", err)
		} else {
			t.Log("Basic message sent successfully")
		}

		// 等待一段时间避免发送过快
		time.Sleep(2 * time.Second)
	})

	// 测试用例：带短描述的卡片消息
	t.Run("Card Message with Short Description", func(t *testing.T) {
		msg := NewMessage(
			"系统通知",
			"## 系统状态报告\n\n- CPU 使用率: 45%\n- 内存使用率: 60%\n- 磁盘空间: 75%\n\n**状态**: 正常",
			WithShort("系统运行正常"),
		)

		err := sender.Send(ctx, msg)
		if err != nil {
			t.Errorf("Failed to send card message: %v", err)
		} else {
			t.Log("Card message sent successfully")
		}

		time.Sleep(2 * time.Second)
	})

	// 测试用例：隐藏 IP 的消息
	t.Run("Message with Hidden IP", func(t *testing.T) {
		msg := NewMessage(
			"隐私消息",
			"这是一条隐藏发送 IP 的消息\n\n发送时间: "+time.Now().Format("15:04:05"),
			WithNoIP(),
		)

		err := sender.Send(ctx, msg)
		if err != nil {
			t.Errorf("Failed to send message with hidden IP: %v", err)
		} else {
			t.Log("Message with hidden IP sent successfully")
		}

		time.Sleep(2 * time.Second)
	})

	// 测试用例：指定渠道的消息
	t.Run("Multi-Channel Message", func(t *testing.T) {
		msg := NewMessage(
			"多渠道消息",
			"这条消息将发送到指定的渠道\n\n支持多种推送方式",
			WithChannel("wecom|dingtalk"), // 企业微信 + 钉钉
		)

		err := sender.Send(ctx, msg)
		if err != nil {
			t.Errorf("Failed to send multi-channel message: %v", err)
		} else {
			t.Log("Multi-channel message sent successfully")
		}

		time.Sleep(2 * time.Second)
	})

	// 测试用例：长文本消息
	t.Run("Long Text Message", func(t *testing.T) {
		longContent := `# 详细报告

## 项目状态
这是一个使用 go-sender 库发送的详细报告。

### 功能特性
- ✅ 支持多种消息类型
- ✅ 支持多渠道推送
- ✅ 支持 Markdown 格式
- ✅ 支持自定义配置

### 技术栈
- Go 语言
- ServerChan API
- HTTP 客户端

### 时间信息
发送时间: ` + time.Now().Format("2006-01-02 15:04:05") + `

---
*此消息由 go-sender 自动发送*`

		msg := NewMessage(
			"详细报告",
			longContent,
			WithShort("项目状态报告"),
		)

		err := sender.Send(ctx, msg)
		if err != nil {
			t.Errorf("Failed to send long text message: %v", err)
		} else {
			t.Log("Long text message sent successfully")
		}

		time.Sleep(2 * time.Second)
	})

	// 测试用例：验证支持的渠道
	t.Run("Supported Channels", func(t *testing.T) {
		channels := GetSupportedChannels()
		if len(channels) == 0 {
			t.Error("No supported channels found")
		} else {
			t.Logf("Found %d supported channels", len(channels))
			for name, code := range channels {
				t.Logf("  - %s (%s)", name, code)
			}
		}
	})

	fmt.Println("\n🎉 测试完成！")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("💡 提示:")
	fmt.Println("1. 请确保将 YOUR_SENDKEY_HERE 替换为你的真实 SendKey")
	fmt.Println("2. SendKey 可以在 ServerChan 官网获取: https://sct.ftqq.com/")
	fmt.Println("3. 如果使用企业版，SendKey 格式为 sctp{num}t{key}")
	fmt.Println("4. 可以在 ServerChan 官网查看消息发送记录")
}

// TestServerChanProviderWithMultipleAccounts 测试多账号配置
// 这个测试演示了如何使用多个 ServerChan 账号进行负载均衡
func TestServerChanProviderWithMultipleAccounts(t *testing.T) {
	// 检查环境变量
	sendKey1 := os.Getenv("SERVERCHAN_KEY_1")
	sendKey2 := os.Getenv("SERVERCHAN_KEY_2")

	if sendKey1 == "" || sendKey2 == "" {
		t.Skip("SERVERCHAN_KEY_1 and SERVERCHAN_KEY_2 environment variables not set, skipping multi-account test")
	}

	// 创建多账号配置
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "account1",
				Key:      sendKey1,
				Weight:   50,
				Disabled: false,
			},
			{
				Name:     "account2",
				Key:      sendKey2,
				Weight:   50,
				Disabled: false,
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create ServerChan provider with multiple accounts: %v", err)
	}

	ctx := context.Background()

	// 测试多账号发送
	t.Run("Multi-Account Message Sending", func(t *testing.T) {
		msg := NewMessage(
			"多账号测试",
			"这条消息将通过负载均衡发送到多个账号\n\n时间: "+time.Now().Format("2006-01-02 15:04:05"),
		)

		err := provider.Send(ctx, msg, nil)
		if err != nil {
			t.Errorf("Failed to send message with multiple accounts: %v", err)
		} else {
			t.Log("Message sent successfully with multiple accounts")
		}
	})
}

// TestServerChanProviderErrorHandling 测试错误处理
// 这个测试验证了当配置无效或网络问题时 provider 的行为
func TestServerChanProviderErrorHandling(t *testing.T) {
	// 测试无效配置
	t.Run("Invalid Configuration", func(t *testing.T) {
		config := Config{
			BaseConfig: core.BaseConfig{
				Strategy: core.StrategyRoundRobin,
			},
			Accounts: []core.Account{}, // 空账号列表
		}

		_, err := New(config)
		if err == nil {
			t.Error("Expected error for empty accounts, but got none")
		} else {
			t.Logf("Expected error received: %v", err)
		}
	})

	// 测试无效的 SendKey
	t.Run("Invalid SendKey", func(t *testing.T) {
		config := Config{
			BaseConfig: core.BaseConfig{
				Strategy: core.StrategyRoundRobin,
			},
			Accounts: []core.Account{
				{
					Name:     "invalid",
					Key:      "invalid_sendkey",
					Weight:   100,
					Disabled: false,
				},
			},
		}

		provider, err := New(config)
		if err != nil {
			t.Fatalf("Failed to create provider with invalid sendkey: %v", err)
		}

		ctx := context.Background()
		msg := NewMessage("测试", "这是一条测试消息")

		err = provider.Send(ctx, msg, nil)
		if err == nil {
			t.Error("Expected error for invalid sendkey, but got none")
		} else {
			t.Logf("Expected error received: %v", err)
		}
	})
}
