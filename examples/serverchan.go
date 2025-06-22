package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/serverchan"
)

func main() {
	// 创建 ServerChan 配置
	// 请将 YOUR_SENDKEY_HERE 替换为你的真实 SendKey
	config := serverchan.Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyRoundRobin,
		},
		Accounts: []core.Account{
			{
				Name: "main",
				Key:  os.Getenv("SERVERCHAN_KEY"),
			},
			// 如果有多个账号，可以添加更多
			// {
			// 	Name: "backup",
			// 	Key:  "YOUR_BACKUP_SENDKEY_HERE",
			// },
		},
	}

	// 创建 ServerChan provider
	provider, err := serverchan.New(config)
	if err != nil {
		log.Fatalf("Failed to create ServerChan provider: %v", err)
	}

	// 创建发送器
	sender := gosender.NewSender(nil) // 使用默认 logger
	sender.RegisterProvider(core.ProviderTypeServerChan, provider, nil)

	// 创建上下文
	ctx := context.Background()

	fmt.Println("🚀 开始测试 ServerChan 发送功能...")
	fmt.Println(strings.Repeat("=", 50))

	// 测试 1: 基础文本消息
	fmt.Println("📝 测试 1: 基础文本消息")
	basicMsg := serverchan.NewMessage(
		"测试消息",
		"这是一条来自 go-sender 的测试消息\n\n时间: "+time.Now().Format("2006-01-02 15:04:05"),
	)

	err = sender.Send(ctx, basicMsg)
	if err != nil {
		log.Printf("❌ 基础消息发送失败: %v", err)
	} else {
		fmt.Println("✅ 基础消息发送成功")
	}

	time.Sleep(2 * time.Second)

	// 测试 2: 带短描述的卡片消息
	fmt.Println("\n📋 测试 2: 卡片消息")
	cardMsg := serverchan.NewMessage(
		"系统通知",
		"## 系统状态报告\n\n- CPU 使用率: 45%\n- 内存使用率: 60%\n- 磁盘空间: 75%\n\n**状态**: 正常",
		serverchan.WithShort("系统运行正常"),
	)

	err = sender.Send(ctx, cardMsg)
	if err != nil {
		log.Printf("❌ 卡片消息发送失败: %v", err)
	} else {
		fmt.Println("✅ 卡片消息发送成功")
	}

	time.Sleep(2 * time.Second)

	// 测试 3: 隐藏 IP 的消息
	fmt.Println("\n🔒 测试 3: 隐藏 IP 的消息")
	noIPMsg := serverchan.NewMessage(
		"隐私消息",
		"这是一条隐藏发送 IP 的消息\n\n发送时间: "+time.Now().Format("15:04:05"),
		serverchan.WithNoIP(),
	)

	err = sender.Send(ctx, noIPMsg)
	if err != nil {
		log.Printf("❌ 隐藏IP消息发送失败: %v", err)
	} else {
		fmt.Println("✅ 隐藏IP消息发送成功")
	}

	time.Sleep(2 * time.Second)

	// 测试 4: 指定渠道的消息
	fmt.Println("\n📱 测试 4: 指定渠道的消息")
	channelMsg := serverchan.NewMessage(
		"多渠道消息",
		"这条消息将发送到指定的渠道\n\n支持多种推送方式",
		serverchan.WithChannel("wecom|dingtalk"), // 企业微信 + 钉钉
	)

	err = sender.Send(ctx, channelMsg)
	if err != nil {
		log.Printf("❌ 多渠道消息发送失败: %v", err)
	} else {
		fmt.Println("✅ 多渠道消息发送成功")
	}

	time.Sleep(2 * time.Second)

	// 测试 5: 长文本消息
	fmt.Println("\n📄 测试 5: 长文本消息")
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

	longMsg := serverchan.NewMessage(
		"详细报告",
		longContent,
		serverchan.WithShort("项目状态报告"),
	)

	err = sender.Send(ctx, longMsg)
	if err != nil {
		log.Printf("❌ 长文本消息发送失败: %v", err)
	} else {
		fmt.Println("✅ 长文本消息发送成功")
	}

	time.Sleep(2 * time.Second)

	// 显示支持的渠道
	fmt.Println("\n📋 支持的推送渠道:")
	channels := serverchan.GetSupportedChannels()
	for name, code := range channels {
		fmt.Printf("  - %s (%s)\n", name, code)
	}

	fmt.Println("\n🎉 测试完成！")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("💡 提示:")
	fmt.Println("1. 请确保将 YOUR_SENDKEY_HERE 替换为你的真实 SendKey")
	fmt.Println("2. SendKey 可以在 ServerChan 官网获取: https://sct.ftqq.com/")
	fmt.Println("3. 如果使用企业版，SendKey 格式为 sctp{num}t{key}")
	fmt.Println("4. 可以在 ServerChan 官网查看消息发送记录")
}
