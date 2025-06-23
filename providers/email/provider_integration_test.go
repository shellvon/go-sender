package email

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

// 集成测试 - 需要真实的邮箱账号和收件人
// 环境变量：
//
//	EMAIL_HOST=smtp.xxx.com
//	EMAIL_PORT=端口
//	EMAIL_USERNAME=xxx
//	EMAIL_PASSWORD=xxx
//	EMAIL_FROM=xxx
//	EMAIL_TO=收件人,多个逗号分隔
func TestEmailProviderIntegration(t *testing.T) {
	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	from := os.Getenv("EMAIL_FROM")
	to := os.Getenv("EMAIL_TO")

	if host == "" || port == 0 || username == "" || password == "" || from == "" || to == "" {
		t.Skip("未设置必要的环境变量，跳过集成测试")
	}

	account := Account{
		Name:     "integration-test",
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
	config := Config{
		BaseConfig: core.BaseConfig{Strategy: core.StrategyRoundRobin},
		Accounts:   []Account{account},
	}
	p, err := New(config)
	if err != nil {
		t.Fatalf("创建email provider失败: %v", err)
	}
	ctx := context.Background()
	toList := splitAndTrim(to, ",")

	t.Run("普通文本邮件", func(t *testing.T) {
		msg := &Message{
			To:      toList,
			Subject: "🧪 集成测试: 普通文本邮件 " + time.Now().Format("15:04:05"),
			Body:    "这是一封集成测试邮件。",
		}
		err := p.Send(ctx, msg)
		if err != nil {
			t.Errorf("普通文本邮件发送失败: %v", err)
		} else {
			t.Log("✅ 普通文本邮件发送成功")
		}
	})

	t.Run("HTML邮件", func(t *testing.T) {
		msg := &Message{
			To:      toList,
			Subject: "🧪 集成测试: HTML邮件 " + time.Now().Format("15:04:05"),
			Body:    "<h1>集成测试</h1><p>这是一封<b>HTML</b>格式的邮件。</p>",
			IsHTML:  true,
		}
		err := p.Send(ctx, msg)
		if err != nil {
			t.Errorf("HTML邮件发送失败: %v", err)
		} else {
			t.Log("✅ HTML邮件发送成功")
		}
	})
}

// splitAndTrim splits a string by sep and trims each element
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
