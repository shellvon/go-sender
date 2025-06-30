package sms

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

// 集成测试 - 需要真实的短信账号和手机号
// 环境变量：
//
//	SMS_PROVIDER_TYPE=smsbao/cl253/aliyun等
//	SMS_APP_ID=xxx
//	SMS_APP_SECRET=xxx
//	SMS_SIGN_NAME=xxx
//	SMS_TO=手机号,多个逗号分隔
//	SMS_TEMPLATE_CODE=模板ID（可选）
//	SMS_TEMPLATE_PARAM_key=val（可选，模板参数）
func TestSMSProviderIntegration(t *testing.T) {
	providerType := os.Getenv("SMS_PROVIDER_TYPE")
	appID := os.Getenv("SMS_APP_ID")
	appSecret := os.Getenv("SMS_APP_SECRET")
	signName := os.Getenv("SMS_SIGN_NAME")
	to := os.Getenv("SMS_TO")
	templateCode := os.Getenv("SMS_TEMPLATE_CODE")

	if providerType == "" || appID == "" || appSecret == "" || signName == "" || to == "" {
		t.Skip("未设置必要的环境变量，跳过集成测试")
	}

	mobiles := splitAndTrim(to, ",")
	provider := SMSProvider{
		Name:      "integration-test",
		Type:      ProviderType(providerType),
		AppID:     appID,
		AppSecret: appSecret,
	}
	config := Config{
		Providers: []SMSProvider{provider},
		Strategy:  core.StrategyRoundRobin,
	}
	p, err := New(config)
	if err != nil {
		t.Fatalf("创建provider失败: %v", err)
	}
	ctx := context.Background()

	t.Run("普通短信", func(t *testing.T) {
		msg := &Message{
			Mobiles:  mobiles,
			Content:  "🧪 集成测试: 普通短信发送" + time.Now().Format("15:04:05"),
			SignName: signName,
		}
		err := p.Send(ctx, msg, nil)
		if err != nil {
			t.Errorf("普通短信发送失败: %v", err)
		} else {
			t.Log("✅ 普通短信发送成功")
		}
	})

	t.Run("模板短信", func(t *testing.T) {
		if templateCode == "" {
			t.Skip("未设置模板ID，跳过模板短信测试")
		}
		params := map[string]string{}
		for _, e := range os.Environ() {
			if len(e) > 18 && e[:18] == "SMS_TEMPLATE_PARAM" {
				kv := splitAndTrim(e, "=")
				if len(kv) == 2 {
					key := kv[0][19:]
					params[key] = kv[1]
				}
			}
		}
		msg := &Message{
			Mobiles:        mobiles,
			TemplateID:     templateCode,
			TemplateParams: params,
			SignName:       signName,
		}
		err := p.Send(ctx, msg, nil)
		if err != nil {
			t.Errorf("模板短信发送失败: %v", err)
		} else {
			t.Log("✅ 模板短信发送成功")
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
