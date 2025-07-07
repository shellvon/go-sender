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
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/lark"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/telegram"
)

// newLoggingHTTPClient returns an http.Client with logging capabilities.
func newLoggingHTTPClient() *http.Client {
	return &http.Client{
		Transport: &LoggingTransport{},
	}
}

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
		Items: []*telegram.Account{{
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
	prov, err := telegram.New(&cfg)
	if err != nil {
		log.Println("Init telegram provider failed:", err)
		return
	}

	// Test all message types
	messages := []struct {
		name string
		msg  core.Message
	}{
		{
			name: "Text Message",
			msg: telegram.Text().
				Chat(chatID).
				Text("Hello from go-sender! This is a text message test.").
				ParseMode("HTML").
				Build(),
		},
		{
			name: "Photo Message",
			msg: telegram.Photo().
				Chat(chatID).
				File("https://picsum.photos/200/300").
				Caption("This is a photo message test.").
				Build(),
		},
		{
			name: "Audio Message",
			msg: telegram.Audio().
				Chat(chatID).
				File("https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav").
				Caption("This is an audio message test.").
				Title("Baby Elephant Walk").
				Performer("Henry Mancini").
				Build(),
		},
		{
			name: "Voice Message",
			msg: telegram.Voice().
				Chat(chatID).
				File("https://www2.cs.uic.edu/~i101/SoundFiles/CantinaBand3.wav").
				Caption("This is a voice message test.").
				Build(),
		},
		{
			name: "Document Message",
			msg: telegram.Document().
				Chat(chatID).
				File("https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf").
				Caption("This is a document message test.").
				Build(),
		},
		{
			name: "Video Message",
			msg: telegram.Video().
				Chat(chatID).
				File("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4").
				Caption("This is a video message test.").
				Build(),
		},
		{
			name: "Animation Message",
			msg: telegram.Animation().
				Chat(chatID).
				File("https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExcDd6bGF4OWN1bnF3OWFxbzF1aHBxM2t1ZDV1bWx1c2Vxd2RqcnR6eCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/3o7TKnxHjgNvUVClHy/giphy.gif").
				Caption("This is an animation message test.").
				Build(),
		},
		{
			name: "Location Message",
			msg: telegram.Location().
				Chat(chatID).
				Coordinates(40.7128, -74.0060).
				Build(),
		},
		{
			name: "Contact Message",
			msg: telegram.Contact().
				Chat(chatID).
				Phone("+1234567890").
				FirstName("John").
				LastName("Doe").
				Build(),
		},
		{
			name: "Poll Message",
			msg: telegram.Poll().
				Chat(chatID).
				Question("What's your favorite programming language?").
				Options(
					telegram.Option("Go"),
					telegram.Option("Python"),
					telegram.Option("JavaScript"),
					telegram.Option("Java"),
				).
				IsAnonymous(true).
				AllowsMultipleAnswers(true).
				Build(),
		},
		{
			name: "Dice Message",
			msg: telegram.Dice().
				Chat(chatID).
				Emoji("🎲").
				Build(),
		},
		{
			name: "Venue Message",
			msg: telegram.Venue().
				Chat(chatID).
				Latitude(40.7128).
				Longitude(-74.0060).
				Title("New York City").
				Address("New York, NY, USA").
				Build(),
		},
	}

	// Send all messages
	for _, m := range messages {
		log.Printf("Sending %s...", m.name)
		err = prov.Send(context.Background(), m.msg, &core.ProviderSendOptions{
			HTTPClient: newLoggingHTTPClient(),
		})
		if err != nil {
			log.Printf("Failed to send %s: %v", m.name, err)
		} else {
			log.Printf("Successfully sent %s!", m.name)
		}
		// Sleep for 3 seconds between messages to avoid hitting rate limits
		// Telegram limits:
		// - 30 messages per second to different chats
		// - 20 messages per minute to the same chat
		time.Sleep(3 * time.Second)
	}
}

// runLarkDemo demonstrates how to send a message using the Lark provider.
// Required environment variables:
//
//	LARK_WEBHOOK_KEY - Lark webhook key (the part after /hook/)
//	LARK_WEBHOOK_SECRET - (optional) Lark webhook secret
//
// This demo will send a simple text message to the specified Lark group.
func runLarkDemo() {
	log.Println("[DEMO] Lark provider")
	key := os.Getenv("LARK_WEBHOOK_KEY")
	secret := os.Getenv("LARK_WEBHOOK_SECRET")
	if key == "" {
		log.Println("Please set LARK_WEBHOOK_KEY environment variable.")
		return
	}
	cfg := lark.Config{
		Items: []*lark.Account{{
			BaseAccount: core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Name: "default",
				},
				Credentials: core.Credentials{
					APIKey:    key,
					APISecret: secret,
				},
			},
		}},
	}
	prov, err := lark.New(&cfg)
	if err != nil {
		log.Println("Init lark provider failed:", err)
		return
	}
	msg := lark.Text().Content("Hello from go-sender example!").Build()
	err = prov.Send(context.Background(), msg, &core.ProviderSendOptions{
		HTTPClient: newLoggingHTTPClient(),
	})
	if err != nil {
		log.Println("Send failed:", err)
	} else {
		log.Println("Send success!")
	}
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
		Items: []*email.Account{{
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
	prov, err := email.New(&cfg)
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
	err = prov.Send(context.Background(), msg, &core.ProviderSendOptions{
		HTTPClient: newLoggingHTTPClient(),
	})
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
		Items: []*sms.Account{{
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
			SignName: from, // 使用环境变量中的签名作为默认签名
		}},
	}
	prov, err := sms.New(&cfg)
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
		msg = sms.Tencent().
			Type(sms.Voice).
			VoiceSdkAppID("1401009332").
			SmsSdkAppID("1401009332").
			To(phone).
			Content("【GoSender】Your code is 1234.").
			TemplateID(templateID).
			SignName(from).
			Build()
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
		HTTPClient: newLoggingHTTPClient(),
	})
	if err != nil {
		log.Println("Send failed:", err)
	} else {
		log.Println("Send success!")
	}
}
