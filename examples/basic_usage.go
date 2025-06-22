package main

import (
	"context"
	"fmt"
	"log"
	"time"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/serverchan"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/webhook"
	"github.com/shellvon/go-sender/providers/wecombot"
)

func main2() {
	// 创建sender实例
	s := gosender.NewSender(nil)

	// 注册SMS提供商
	smsProvider, err := sms.New(sms.Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyWeighted,
		},
		Providers: []sms.SMSProvider{
			{
				Name:      "tencent",
				Type:      sms.ProviderTypeTencent,
				AppID:     "your_tencent_app_id",
				AppSecret: "your_tencent_app_secret",
				Weight:    100,
				Disabled:  false,
			},
			{
				Name:      "aliyun",
				Type:      sms.ProviderTypeAliyun,
				AppID:     "your_aliyun_access_key",
				AppSecret: "your_aliyun_access_secret",
				Weight:    80,
				Disabled:  false,
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create SMS provider: %v", err)
	}
	s.RegisterProvider(core.ProviderTypeSMS, smsProvider, nil)

	// 注册Email提供商
	emailProvider, err := email.New(email.Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyWeighted,
		},
		Accounts: []email.Account{
			{
				Name:     "smtp1",
				Host:     "smtp.gmail.com",
				Port:     587,
				Username: "your_email@gmail.com",
				Password: "your_app_password",
				Weight:   100,
				Disabled: false,
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create email provider: %v", err)
	}
	s.RegisterProvider(core.ProviderTypeEmail, emailProvider, nil)

	// 注册Webhook提供商
	webhookProvider, err := webhook.New(webhook.Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyWeighted,
		},
		Endpoints: []webhook.Endpoint{
			{
				Name:     "telegram",
				URL:      "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage",
				Method:   "POST",
				Headers:  map[string]string{"Content-Type": "application/json"},
				Weight:   100,
				Disabled: false,
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create webhook provider: %v", err)
	}
	s.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 注册企业微信机器人提供商
	wecomProvider, err := wecombot.New(wecombot.Config{
		BaseConfig: core.BaseConfig{
			Strategy: core.StrategyWeighted,
		},
		Accounts: []core.Account{
			{

				Name:     "bot1",
				Weight:   100,
				Disabled: false,
				Key:      "YOUR_KEY",
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create wecom bot provider: %v", err)
	}
	s.RegisterProvider(core.ProviderTypeWecombot, wecomProvider, nil)

	// 设置全局重试策略
	s.SetRetryPolicy(&core.RetryPolicy{
		MaxAttempts:   3,
		InitialDelay:  time.Second,
		MaxDelay:      time.Second * 30,
		BackoffFactor: 2.0,
	})

	ctx := context.Background()

	// 示例1: 发送SMS消息
	fmt.Println("=== 发送SMS消息 ===")
	smsMsg := sms.NewMessage("13800138000",
		sms.WithContent("您的验证码是123456，5分钟内有效"),
		sms.WithTemplateCode("SMS_123456789"),
		sms.WithTemplateParams(map[string]string{
			"code": "123456",
			"time": "5分钟",
		}))

	err = s.Send(ctx, smsMsg)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
	} else {
		fmt.Println("SMS sent successfully")
	}

	// 示例2: 发送Email消息
	fmt.Println("\n=== 发送Email消息 ===")
	emailMsg := email.NewMessage("测试邮件", "这是一封测试邮件的内容").
		WithTo("recipient@example.com").
		WithFrom("sender@example.com")

	err = s.Send(ctx, emailMsg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	} else {
		fmt.Println("Email sent successfully")
	}

	// 示例3: 发送Webhook消息
	fmt.Println("\n=== 发送Webhook消息 ===")
	webhookMsg := webhook.NewMessage([]byte(`{"chat_id":"@your_channel","text":"Hello from go-sender!"}`),
		webhook.WithHeaders(map[string]string{"Content-Type": "application/json"}))

	err = s.Send(ctx, webhookMsg)
	if err != nil {
		log.Printf("Failed to send webhook: %v", err)
	} else {
		fmt.Println("Webhook sent successfully")
	}

	// 示例4: 发送企业微信机器人消息
	fmt.Println("\n=== 发送企业微信机器人消息 ===")
	wecomMsg := wecombot.NewTextMessage("Hello from go-sender!",
		wecombot.WithMentionedList([]string{"@all"}))

	err = s.Send(ctx, wecomMsg)
	if err != nil {
		log.Printf("Failed to send wecom bot message: %v", err)
	} else {
		fmt.Println("Wecom bot message sent successfully")
	}

	// 示例5: 发送 ServerChan 消息（函数选项模式）
	fmt.Println("\n=== 发送 ServerChan 消息 ===")
	serverchanMsg := serverchan.NewMessage("测试标题", "这是测试内容",
		serverchan.WithShort("简短摘要"),
		serverchan.WithChannel("wecom|dingtalk"),
		serverchan.WithNoIP())

	err = s.Send(ctx, serverchanMsg)
	if err != nil {
		log.Printf("Failed to send serverchan message: %v", err)
	} else {
		fmt.Println("ServerChan message sent successfully")
	}

	// 示例6: 使用自定义重试策略
	fmt.Println("\n=== 使用自定义重试策略 ===")
	customRetryMsg := sms.NewMessage("13800138000",
		sms.WithContent("重要消息，需要多次重试"))

	// 使用自定义发送选项
	err = s.Send(ctx, customRetryMsg, core.WithSendRetryPolicy(&core.RetryPolicy{
		MaxAttempts:   5,
		InitialDelay:  time.Second * 2,
		MaxDelay:      time.Minute,
		BackoffFactor: 2.0,
	}), core.WithSendTimeout(time.Second*30))

	if err != nil {
		log.Printf("Failed to send with custom retry: %v", err)
	} else {
		fmt.Println("Message sent with custom retry policy")
	}
}
