package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/telegram"
)

// main is the unified entry for all provider examples.
// Usage:
//
//	go run main.go --provider=email     # Test email sending
//	go run main.go --provider=telegram  # Test telegram bot
//	go run main.go --provider=sms       # Test SMS (aliyun/tencent/cl253)
//	go run main.go --provider=lark      # Test lark (TODO)
//
// Each demo requires you to set the corresponding environment variables.
func main() {
	provider := flag.String("provider", "", "Provider to demo: telegram, lark, email, sms")
	flag.Parse()

	switch *provider {
	case "telegram":
		runTelegramDemo()
	case "lark":
		runLarkDemo()
	case "email":
		runEmailDemo()
	case "sms":
		runSMSDemo()
	default:
		log.Println("Usage: go run main.go --provider=[telegram|lark|email|sms]")
		os.Exit(1)
	}
}

// runTelegramDemo demonstrates how to send a message via Telegram Bot API.
// Required environment variables:
//
//	TELEGRAM_BOT_TOKEN - your bot token
//	TELEGRAM_CHAT_ID   - target chat id (user or group)
//
// This demo will send a simple text message to the specified chat.
func runTelegramDemo() {
	log.Println("[DEMO] Telegram provider")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if botToken == "" || chatID == "" {
		log.Println("Please set TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables.")
		return
	}
	cfg := telegram.Config{
		Accounts: []*telegram.Account{{
			BaseAccount: core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Name: "default",
				},
				Credentials: core.Credentials{
					APIKey: botToken,
				},
			},
		}},
	}
	prov, err := telegram.New(cfg)
	if err != nil {
		log.Println("Init telegram provider failed:", err)
		return
	}
	msg := telegram.NewTextMessage(chatID, "Hello from go-sender example!")
	err = prov.Send(context.Background(), msg, nil)
	if err != nil {
		log.Println("Send failed:", err)
	} else {
		log.Println("Send success!")
	}
}

// runLarkDemo is a placeholder for Lark provider demo.
// TODO: Implement Lark message sending example.
func runLarkDemo() {
	log.Println("[DEMO] Lark provider: please implement your test logic here.")
	// TODO: Add real lark send example
}

// runEmailDemo demonstrates how to send an email using the email provider.
// Required environment variables:
//
//	EMAIL_HOST - SMTP server host
//	EMAIL_PORT - SMTP server port (int)
//	EMAIL_USER - SMTP username (also used as From)
//	EMAIL_PASS - SMTP password
//	EMAIL_TO   - recipient email address
//
// This demo will send a simple email to the specified recipient.
func runEmailDemo() {
	log.Println("[DEMO] Email provider")
	host := os.Getenv("EMAIL_HOST")
	portStr := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")
	to := os.Getenv("EMAIL_TO")
	if host == "" || portStr == "" || user == "" || pass == "" || to == "" {
		log.Println("Please set EMAIL_HOST, EMAIL_PORT, EMAIL_USER, EMAIL_PASS, EMAIL_TO environment variables.")
		return
	}
	// Convert port to int
	var port int
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		log.Println("EMAIL_PORT must be an integer.")
		return
	}
	cfg := email.Config{
		Accounts: []*email.Account{{
			BaseAccount: core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Name: "default",
				},
				Credentials: core.Credentials{
					APIKey:    user,
					APISecret: pass,
				},
			},
			Host: host,
			Port: port,
			From: user,
		}},
	}
	prov, err := email.New(cfg)
	if err != nil {
		log.Println("Init email provider failed:", err)
		return
	}
	msg := email.Email().
		To(to).
		Body("This is a test email from go-sender example.").
		From(user).
		Subject("Go-Sender Example").
		Attach("main.go").
		Build()
	err = prov.Send(context.Background(), msg, nil)
	if err != nil {
		log.Println("Send failed:", err)
	} else {
		log.Println("Send success!")
	}
}

type LoggingTransport struct {
	Transport http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 打印请求信息
	log.Println("=== HTTP Request ===")

	// 打印请求方法和 URL
	log.Printf("Method: %s\n", req.Method)
	log.Printf("URL: %s\n", req.URL.String())

	// 打印请求头
	log.Println("Headers:")
	for key, values := range req.Header {
		for _, value := range values {
			log.Printf("\t%s: %s\n", key, value)
		}
	}

	// 打印请求体（如果存在）
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
		} else {
			log.Println("Body:")
			log.Printf("\t%s\n", string(bodyBytes))
			// 重新设置请求体，因为 io.ReadAll 会消耗掉原始 Body
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
	log.Println("===================")

	// 使用默认的 Transport 或自定义 Transport
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// 执行请求
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 打印响应信息
	log.Println("\n=== HTTP Response ===")

	// 打印状态码和状态
	log.Printf("Status: %s\n", resp.Status)

	// 打印响应头
	log.Println("Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			log.Printf("\t%s: %s\n", key, value)
		}
	}

	// 打印响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
	} else {
		log.Println("Body:")
		log.Printf("\t%s\n", string(bodyBytes))
		// 重新设置响应体，因为 io.ReadAll 会消耗掉原始 Body
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	log.Println("===================")

	return resp, nil
}

// runSMSDemo demonstrates how to send an SMS using the sms provider.
// Required environment variables:
//
//	SMS_PROVIDER_TYPE - provider type (aliyun, tencent, cl253, ...)
//	SMS_KEY           - API key or AppID
//	SMS_SECRET        - API secret (optional, depends on provider)
//	SMS_SIGN          - SMS sign name (for aliyun/tencent, etc., optional for cl253)
//	SMS_PHONE         - recipient phone number
//
// This demo will send a simple text message to the specified phone number using the selected provider.
func runSMSDemo() {
	log.Println("[DEMO] SMS provider")
	providerType := os.Getenv("SMS_PROVIDER_TYPE") // aliyun, tencent, cl253, etc.
	key := os.Getenv("SMS_KEY")
	secret := os.Getenv("SMS_SECRET")
	from := os.Getenv("SMS_SIGN")
	phone := os.Getenv("SMS_PHONE")
	templateID := os.Getenv("SMS_TEMPLATE_ID")
	if providerType == "" || key == "" || phone == "" {
		log.Println(
			"Please set SMS_PROVIDER_TYPE, SMS_KEY, SMS_PHONE (and optionally SMS_SECRET, SMS_SIGN) environment variables.",
		)
		return
	}
	cfg := sms.Config{
		Accounts: []*sms.Account{{
			BaseAccount: core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Name:    "default",
					SubType: providerType,
				},
				Credentials: core.Credentials{
					APIKey:    key,
					APISecret: secret,
				},
			},
		}},
	}
	prov, err := sms.New(cfg)
	if err != nil {
		log.Println("Init sms provider failed:", err)
		return
	}
	var msg *sms.Message
	switch providerType {
	case "aliyun":
		// 阿里云短信示例
		msg = sms.Aliyun().
			To(phone).
			Content("【GoSender】Your code is 1234.").
			TemplateID(templateID).
			Params(map[string]string{"code": "1234"}).
			SignName(from).
			Build()
	case "tencent":
		// 腾讯云短信示例
		msg = sms.Tencent().To(phone).Content("【GoSender】Your code is 1234.").SignName(from).Build()
	case "cl253":
		// 创蓝253短信示例
		msg = sms.Cl253().To(phone).Content("【GoSender】Your code is 1234.").Build()

	case "volc":
		msg = sms.Volc().To(phone).
			//	Content("【GoSender】Your code is 1234.").
			TemplateID(templateID).
			Params(map[string]string{"code": "1234"}).
			SignName(from).
			Build()
	default:
		log.Println("Unsupported provider type for demo. Supported: aliyun, tencent, cl25, volc")
		return
	}
	err = prov.Send(context.Background(), msg, &core.ProviderSendOptions{
		HTTPClient: &http.Client{
			Transport: &LoggingTransport{},
		},
	})
	if err != nil {
		log.Println("Send failed:", err)
	} else {
		log.Println("Send success!")
	}
}
