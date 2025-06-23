package telegram

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

// TestTelegramProviderIntegration 集成测试 - 需要真实的 bot token
// 使用方法：
// 1. 设置环境变量 TELEGRAM_BOT_TOKEN=your_bot_token
// 2. 设置环境变量 TELEGRAM_CHAT_ID=your_chat_id (可以是用户ID、群组ID或频道ID)
// 3. 运行测试: go test -v -run TestTelegramProviderIntegration
func TestTelegramProviderIntegration(t *testing.T) {
	// 检查环境变量
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" {
		t.Skip("TELEGRAM_BOT_TOKEN environment variable not set, skipping integration test")
	}

	if chatID == "" {
		t.Skip("TELEGRAM_CHAT_ID environment variable not set, skipping integration test")
	}

	// 添加调试信息
	t.Logf("🔍 调试信息:")
	t.Logf("   Bot Token: %s...", botToken[:10])
	t.Logf("   Chat ID: %s", chatID)
	t.Logf("   Chat ID 类型: %T", chatID)

	// 创建配置
	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "integration-bot",
				Key:      botToken,
				Weight:   100,
				Disabled: false,
			},
		},
	}

	// 创建 provider
	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	ctx := context.Background()

	// 测试0: 验证 Bot 信息
	t.Run("Bot Info Verification", func(t *testing.T) {
		// 这里可以添加一个简单的 Bot 信息验证
		// 由于我们没有直接的 Bot API 调用，我们通过发送最简单的消息来验证
		t.Log("🔍 验证 Bot 配置...")
		t.Logf("   Bot Token 长度: %d", len(botToken))
		t.Logf("   Chat ID: %s", chatID)

		// 检查 Chat ID 格式
		if len(chatID) > 0 {
			if chatID[0] == '@' {
				t.Log("   Chat ID 格式: 频道/群组用户名")
			} else if chatID[0] == '-' {
				t.Log("   Chat ID 格式: 群组 ID")
			} else {
				t.Log("   Chat ID 格式: 用户 ID")
			}
		}
	})

	// 测试1: 基本文本消息
	t.Run("Basic Text Message", func(t *testing.T) {
		message := NewTextMessage(chatID, "🧪 集成测试: 基本文本消息")
		t.Logf("📤 发送消息到 Chat ID: %s", chatID)
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("❌ Failed to send basic text message: %v", err)
		} else {
			t.Log("✅ 基本文本消息发送成功")
		}
		time.Sleep(1 * time.Second) // 避免发送过快
	})

	// 测试2: 带格式的文本消息
	t.Run("Formatted Text Message", func(t *testing.T) {
		message := NewTextMessage(chatID,
			"🧪 *集成测试*: 格式化文本消息\n"+
				"• 支持 *粗体*\n"+
				"• 支持 _斜体_\n"+
				"• 支持 `代码`\n"+
				"• 支持 [链接](https://telegram.org)",
			WithParseMode("Markdown"))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send formatted text message: %v", err)
		} else {
			t.Log("✅ 格式化文本消息发送成功")
		}
		time.Sleep(1 * time.Second)
	})

	// 测试3: HTML 格式文本消息
	t.Run("HTML Text Message", func(t *testing.T) {
		message := NewTextMessage(chatID,
			"🧪 <b>集成测试</b>: HTML 格式文本消息"+
				"• 支持 <b>粗体</b>"+
				"• 支持 <i>斜体</i>"+
				"• 支持 <code>代码</code>"+
				"• 支持 <a href=\"https://telegram.org\">链接</a>",
			WithParseMode("HTML"))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send HTML text message: %v", err)
		} else {
			t.Log("✅ HTML 文本消息发送成功")
		}
		time.Sleep(1 * time.Second)
	})

	// 测试4: 带选项的文本消息
	t.Run("Text Message with Options", func(t *testing.T) {
		message := NewTextMessage(chatID, "🧪 集成测试: 带选项的文本消息",
			WithDisableWebPreview(true),
			WithSilent(true),
			WithProtectContent(true))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send text message with options: %v", err)
		} else {
			t.Log("✅ 带选项的文本消息发送成功")
		}
		time.Sleep(1 * time.Second)
	})

	// 测试5: 图片消息（使用网络图片）
	t.Run("Photo Message", func(t *testing.T) {
		message := NewPhotoMessage(chatID,
			"https://picsum.photos/400/300", // 使用随机图片服务
			WithCaption("🧪 集成测试: 图片消息\n这是一张测试图片"),
			WithParseMode("Markdown"))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send photo message: %v", err)
		} else {
			t.Log("✅ 图片消息发送成功")
		}
		time.Sleep(2 * time.Second) // 图片消息需要更多时间
	})

	// 测试6: 文档消息（使用网络文档）
	t.Run("Document Message", func(t *testing.T) {
		message := NewDocumentMessage(chatID,
			"https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf", // 测试PDF
			WithCaption("🧪 集成测试: 文档消息\n这是一个测试PDF文档"),
			WithParseMode("Markdown"))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send document message: %v", err)
		} else {
			t.Log("✅ 文档消息发送成功")
		}
		time.Sleep(2 * time.Second)
	})

	// 测试7: 位置消息
	t.Run("Location Message", func(t *testing.T) {
		message := NewLocationMessage(chatID, 40.7128, -74.0060, // 纽约坐标
			WithSilent(true))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send location message: %v", err)
		} else {
			t.Log("✅ 位置消息发送成功")
		}
		time.Sleep(1 * time.Second)
	})

	// 测试8: 联系人消息
	t.Run("Contact Message", func(t *testing.T) {
		message := NewContactMessage(chatID, "+1234567890", "Test User",
			WithContactLastName("Integration"),
			WithContactVCard("BEGIN:VCARD\nVERSION:3.0\nFN:Test User Integration\nTEL:+1234567890\nEND:VCARD"))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send contact message: %v", err)
		} else {
			t.Log("✅ 联系人消息发送成功")
		}
		time.Sleep(1 * time.Second)
	})

	// 测试9: 投票消息
	t.Run("Poll Message", func(t *testing.T) {
		options := []InputPollOption{
			{Text: "Go"},
			{Text: "Python"},
			{Text: "JavaScript"},
			{Text: "Java"},
			{Text: "Rust"},
		}
		message := NewPollMessage(chatID, "🧪 集成测试: 你最喜欢哪种编程语言？",
			options,
			WithPollIsAnonymous(false),
			WithPollType("regular"),
			WithPollAllowsMultipleAnswers(false))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send poll message: %v", err)
		} else {
			t.Log("✅ 投票消息发送成功")
		}
		time.Sleep(1 * time.Second)
	})

	// 测试10: Dice 消息
	t.Run("Dice Message", func(t *testing.T) {
		message := NewDiceMessage(chatID, WithDiceEmoji("🎯"))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send dice message: %v", err)
		} else {
			t.Log("✅ Dice 消息发送成功")
		}
		time.Sleep(1 * time.Second)
	})

	// 测试11: 语音消息
	t.Run("Voice Message", func(t *testing.T) {
		// 使用一个更可靠的音频文件 URL，或者使用 file_id 进行测试
		// 注意：Telegram 对音频文件有特定要求，必须是可访问的音频文件
		message := NewVoiceMessage(chatID, "http://commondatastorage.googleapis.com/codeskulptor-assets/week7-button.m4a",
			WithCaption("🧪 集成测试: 语音消息"),
			WithParseMode("Markdown"),
			WithVoiceDuration(3))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Logf("Voice message failed (URL might be inaccessible): %v", err)
		} else {
			t.Log("✅ 语音消息发送成功")
		}
		time.Sleep(2 * time.Second)
	})

	// 测试12: 视频消息
	t.Run("Video Message", func(t *testing.T) {
		message := NewVideoMessage(chatID, "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
			WithCaption("🧪 集成测试: 视频消息"),
			WithParseMode("Markdown"),
			WithVideoDuration(10),
			WithVideoWidth(1280),
			WithVideoHeight(720))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Logf("Video message failed (URL might be inaccessible): %v", err)
		} else {
			t.Log("✅ 视频消息发送成功")
		}
		time.Sleep(3 * time.Second)
	})

	// 测试13: 动画消息
	t.Run("Animation Message", func(t *testing.T) {
		message := NewAnimationMessage(chatID, "https://media.giphy.com/media/3o7abKhOpu0NwenH3O/giphy.gif",
			WithCaption("🧪 集成测试: 动画消息"),
			WithParseMode("Markdown"),
			WithAnimationDuration(5),
			WithAnimationWidth(480),
			WithAnimationHeight(270))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Logf("Animation message failed (URL might be inaccessible): %v", err)
		} else {
			t.Log("✅ 动画消息发送成功")
		}
		time.Sleep(2 * time.Second)
	})

	// 测试14: 视频笔记消息
	t.Run("Video Note Message", func(t *testing.T) {
		message := NewVideoNoteMessage(chatID, "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
			WithVideoNoteDuration(5),
			WithVideoNoteLength(360))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Logf("Video note message failed (URL might be inaccessible): %v", err)
		} else {
			t.Log("✅ 视频笔记消息发送成功")
		}
		time.Sleep(2 * time.Second)
	})

	// 测试15: 音频消息
	t.Run("Audio Message", func(t *testing.T) {
		message := NewAudioMessage(chatID, "https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav",
			WithCaption("🧪 集成测试: 音频消息"),
			WithParseMode("Markdown"),
			WithAudioDuration(3),
			WithAudioPerformer("Test Performer"),
			WithAudioTitle("Test Audio"))
		err := provider.Send(ctx, message)
		if err != nil {
			t.Logf("Audio message failed (URL might be inaccessible): %v", err)
		} else {
			t.Log("✅ 音频消息发送成功")
		}
		time.Sleep(2 * time.Second)
	})
}

// TestTelegramProviderIntegrationErrorCases 集成测试错误情况
func TestTelegramProviderIntegrationErrorCases(t *testing.T) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		t.Skip("TELEGRAM_BOT_TOKEN environment variable not set, skipping integration test")
	}

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name:     "integration-bot",
				Key:      botToken,
				Weight:   100,
				Disabled: false,
			},
		},
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	ctx := context.Background()

	// 测试1: 无效的聊天ID
	t.Run("Invalid Chat ID", func(t *testing.T) {
		message := NewTextMessage("invalid_chat_id", "This should fail")
		err := provider.Send(ctx, message)
		if err == nil {
			t.Error("Expected error for invalid chat ID but got success")
		} else {
			t.Logf("✅ 正确捕获了无效聊天ID错误: %v", err)
		}
	})

	// 测试2: 空文本
	t.Run("Empty Text", func(t *testing.T) {
		message := NewTextMessage("@test_channel", "")
		err := provider.Send(ctx, message)
		if err == nil {
			t.Error("Expected error for empty text but got success")
		} else {
			t.Logf("✅ 正确捕获了空文本错误: %v", err)
		}
	})

	// 测试3: 无效的图片URL
	t.Run("Invalid Photo URL", func(t *testing.T) {
		message := NewPhotoMessage("@test_channel", "https://invalid-url-that-does-not-exist.com/image.jpg")
		err := provider.Send(ctx, message)
		if err == nil {
			t.Error("Expected error for invalid photo URL but got success")
		} else {
			t.Logf("✅ 正确捕获了无效图片URL错误: %v", err)
		}
	})
}

// TestTelegramProviderIntegrationMultipleAccounts 测试多账户集成
func TestTelegramProviderIntegrationMultipleAccounts(t *testing.T) {
	botToken1 := os.Getenv("TELEGRAM_BOT_TOKEN")
	botToken2 := os.Getenv("TELEGRAM_BOT_TOKEN_2") // 可选的第二个bot token
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken1 == "" {
		t.Skip("TELEGRAM_BOT_TOKEN environment variable not set, skipping integration test")
	}

	if chatID == "" {
		t.Skip("TELEGRAM_CHAT_ID environment variable not set, skipping integration test")
	}

	// 创建多账户配置
	accounts := []core.Account{
		{
			Name:     "bot1",
			Key:      botToken1,
			Weight:   100,
			Disabled: false,
		},
	}

	// 如果有第二个bot token，添加到配置中
	if botToken2 != "" {
		accounts = append(accounts, core.Account{
			Name:     "bot2",
			Key:      botToken2,
			Weight:   50,
			Disabled: false,
		})
	}

	config := Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: accounts,
	}

	provider, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create telegram provider: %v", err)
	}

	ctx := context.Background()

	// 测试多账户发送
	t.Run("Multiple Accounts", func(t *testing.T) {
		message := NewTextMessage(chatID, "🧪 集成测试: 多账户配置测试")
		err := provider.Send(ctx, message)
		if err != nil {
			t.Errorf("Failed to send message with multiple accounts: %v", err)
		} else {
			t.Log("✅ 多账户配置测试成功")
		}
	})
}
