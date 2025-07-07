package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/dingtalk"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/serverchan"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/telegram"
	"github.com/shellvon/go-sender/providers/wecombot"
)

// createEmptyImage creates a 100x100 transparent PNG image and returns its base64 and MD5.
func createEmptyImage() (string, string, error) {
	// Create a new 100x100 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", "", fmt.Errorf("failed to encode PNG: %w", err)
	}

	// Calculate MD5 of the raw image data
	hash := md5.New()
	if _, err := hash.Write(buf.Bytes()); err != nil {
		return "", "", fmt.Errorf("failed to calculate MD5: %w", err)
	}
	imgMD5 := fmt.Sprintf("%x", hash.Sum(nil))

	// Convert to base64
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return imgBase64, imgMD5, nil
}

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
//	go run main.go --provider=serverchan # Test ServerChan
//
// Each demo requires you to set the corresponding environment variables.
func main() {
	provider := flag.String("provider", "", "Provider to demo: telegram, lark, email, sms, dingtalk, serverchan")
	flag.Parse()

	switch *provider {
	case "telegram":
		runTelegramDemo()
	case "lark":
		runLarkDemo()
	case "wecombot":
		runWecombotDemo()
	case "email":
		runEmailDemo()
	case "sms":
		runSMSDemo()
	case "dingtalk":
		runDingTalkDemo()
	case "serverchan":
		runServerChanDemo()
	default:
		log.Println("Usage: go run main.go --provider=[telegram|lark|email|sms|dingtalk|serverchan]")
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
//	LARK_BOT_TOKEN - Lark Bot token
//	LARK_CHAT_ID   - Lark chat ID
//
// This demo will send a simple text message to the specified chat.
func runLarkDemo() {
	log.Println("[DEMO] Lark provider")
}

// runLarkDemo demonstrates how to send a message using the Lark provider.
// Required environment variables:
//
//	WECOM_BOT_KEY - WeCom Bot webhook key (the part after /key=)
//
// This demo will send various types of messages to the specified WeCom group.
func runWecombotDemo() {
	log.Println("[DEMO] WeCom Bot provider")
	key := os.Getenv("WECOM_BOT_KEY")
	if key == "" {
		log.Println("Please set WECOM_BOT_KEY environment variable.")
		return
	}
	cfg := wecombot.Config{
		Items: []*wecombot.Account{{
			BaseAccount: core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Name: "default",
				},
				Credentials: core.Credentials{
					APIKey: key,
				},
			},
		}},
	}
	prov, err := wecombot.New(&cfg)
	if err != nil {
		log.Println("Init wecom bot provider failed:", err)
		return
	}

	// Test all message types
	messages := []struct {
		name string
		msg  core.Message
	}{
		{
			name: "Text Message",
			msg: wecombot.Text().
				Content("Hello from go-sender! This is a text message test.").
				MentionUsers([]string{"@all"}).
				Build(),
		},
		{
			name: "Image Message",
			msg: func() core.Message {
				base64Img, imgMD5, err := createEmptyImage()
				if err != nil {
					log.Printf("Failed to create empty image: %v", err)
					return nil
				}
				return wecombot.Image().
					Base64(base64Img).
					MD5(imgMD5).
					Build()
			}(),
		},
		{
			name: "Markdown Message (Legacy)",
			msg: wecombot.Markdown().
				Content("# Hello from go-sender!\n\n**This** is a *markdown* message test.\n\n> Legacy version").
				Version(wecombot.MarkdownVersionLegacy).
				Build(),
		},
		{
			name: "Markdown Message (V2)",
			msg: wecombot.Markdown().
				Content("# 一、标题\n" +
					"## 二级标题\n" +
					"### 三级标题\n\n" +
					"# 二、字体\n" +
					"*斜体*\n" +
					"**加粗**\n\n" +
					"# 三、列表 \n" +
					"- 无序列表 1 \n" +
					"- 无序列表 2\n" +
					"  - 无序列表 2.1\n" +
					"  - 无序列表 2.2\n" +
					"1. 有序列表 1\n" +
					"2. 有序列表 2\n\n" +
					"# 四、引用\n" +
					"> 一级引用\n" +
					">> 二级引用\n" +
					">>> 三级引用\n\n" +
					"# 五、链接\n" +
					"[这是一个链接](https://work.weixin.qq.com/api/doc)\n" +
					"![](https://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png)\n\n" +
					"# 六、分割线\n" +
					"---\n\n" +
					"# 七、代码\n" +
					"`这是行内代码`\n" +
					"```\n" +
					"这是独立代码块\n" +
					"```\n\n" +
					"# 八、表格\n" +
					"| 姓名 | 文化衫尺寸 | 收货地址 |\n" +
					"| :--- | :---: | ---: |\n" +
					"| 张三 | S | 广州 |\n" +
					"| 李四 | L | 深圳 |",
				).
				Version(wecombot.MarkdownVersionV2).
				Build(),
		},
		{
			name: "File Message",
			msg: wecombot.File().
				LocalPath("main.go").
				Build(),
		},
		{
			name: "News Message",
			msg: wecombot.News().
				AddArticle(
					"Go-Sender Example",
					"Testing WeCom Bot News Message",
					"https://github.com/shellvon/go-sender",
					"https://golang.org/lib/godoc/images/go-logo-blue.svg",
				).
				AddArticle(
					"Second Article",
					"Another article in the news message",
					"https://github.com/shellvon/go-sender/blob/main/docs/getting-started.md",
					"https://golang.org/lib/godoc/images/cloud.png",
				).
				Build(),
		},
		{
			name: "Template Card (Text Notice)",
			msg: wecombot.Card(wecombot.CardTypeTextNotice).
				MainTitle("Go-Sender Demo", "Testing Template Card").
				SubTitle("This is a text notice template card.\nClick to view more details.").
				JumpURL("https://github.com/shellvon/go-sender").
				Build(),
		},
		{
			name: "Template Card (News Notice)",
			msg: wecombot.Card(wecombot.CardTypeNewsNotice).
				MainTitle("Go-Sender News", "Important Updates").
				CardImage("https://golang.org/lib/godoc/images/go-logo-blue.svg", 1.8).
				ImageTextArea(
					"Latest Release",
					"Check out our new features and improvements",
					"https://golang.org/lib/godoc/images/cloud.png",
					"https://github.com/shellvon/go-sender/releases",
				).
				AddVerticalContent("New Features", "Support for all message types").
				AddVerticalContent("Improvements", "Enhanced error handling").
				JumpURL("https://github.com/shellvon/go-sender/releases").
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
		// Sleep for 2 seconds between messages to avoid hitting rate limits
		time.Sleep(2 * time.Second)
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

// runDingTalkDemo demonstrates how to send messages using the DingTalk provider.
// Required environment variables:
//
//	DINGTALK_BOT_TOKEN - DingTalk Bot token
//	DINGTALK_BOT_SECRET - DingTalk Bot secret
//
// This demo will send various types of messages to the specified DingTalk group.
func runDingTalkDemo() {
	log.Println("[DEMO] DingTalk provider")
	key := os.Getenv("DINGTALK_BOT_TOKEN")
	secret := os.Getenv("DINGTALK_BOT_SECRET")
	if key == "" || secret == "" {
		log.Println("Please set DINGTALK_BOT_TOKEN and DINGTALK_BOT_SECRET environment variables.")
		return
	}
	cfg := dingtalk.Config{
		Items: []*dingtalk.Account{{
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
	prov, _ := dingtalk.New(&cfg)
	demos := []struct {
		name string
		msg  core.Message
	}{
		{
			name: "Text Message",
			msg: dingtalk.Text().
				Content("这是一条文本消息\n支持换行和@指定用户").
				AtMobiles([]string{"***REMOVED***"}).
				AtAll().
				Build(),
		},
		{
			name: "Markdown Message",
			msg: dingtalk.Markdown().
				Title("这是一条 Markdown 消息").
				Text("# 标题\n" +
					"## 二级标题\n" +
					"### 三级标题\n\n" +
					"#### 文本样式\n" +
					"- **加粗文本**\n" +
					"- *斜体文本*\n" +
					"- ~~删除线文本~~\n\n" +
					"#### 图片\n" +
					"![](https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png)\n\n" +
					"#### 链接\n" +
					"[点击跳转](https://www.dingtalk.com)\n\n" +
					"#### 代码段\n" +
					"```go\n" +
					"package main\n\n" +
					"func main() {\n" +
					"    println(\"Hello DingTalk\")\n" +
					"}\n" +
					"```").
				Build(),
		},
		{
			name: "Link Message",
			msg: dingtalk.Link().
				Title("这是一条链接消息").
				Text("链接消息支持标题、描述、链接和图片").
				MessageURL("https://www.dingtalk.com").
				PicURL("https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png").
				Build(),
		},
		{
			name: "ActionCard Message (整体跳转)",
			msg: dingtalk.ActionCard().
				Title("这是一条整体跳转 ActionCard").
				Text("# 惊喜\n"+
					"## 惊喜来了\n\n"+
					"![](https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png)\n\n"+
					"新版本发布了！").
				SingleButton("查看详情", "https://www.dingtalk.com").
				BtnOrientation("0").
				Build(),
		},
		{
			name: "ActionCard Message (独立跳转)",
			msg: dingtalk.ActionCard().
				Title("这是一条独立跳转 ActionCard").
				Text("# 惊喜\n"+
					"## 惊喜来了\n\n"+
					"![](https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png)\n\n"+
					"新版本发布了！").
				AddButton("文档", "https://www.dingtalk.com/doc").
				AddButton("示例", "https://www.dingtalk.com/example").
				BtnOrientation("1").
				Build(),
		},
		{
			name: "FeedCard Message",
			msg: dingtalk.FeedCard().
				AddLink(
					"新版本发布",
					"https://www.dingtalk.com/release",
					"https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png",
				).
				AddLink(
					"使用文档",
					"https://www.dingtalk.com/doc",
					"https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png",
				).
				Build(),
		},
	}

	for _, demo := range demos {
		log.Printf("Testing %s...\n", demo.name)
		if err := prov.Send(context.Background(), demo.msg, &core.ProviderSendOptions{
			HTTPClient: newLoggingHTTPClient(),
		}); err != nil {
			log.Printf("%s failed: %v\n", demo.name, err)
		} else {
			log.Printf("%s success!\n", demo.name)
		}
		time.Sleep(3 * time.Second) // 限流保护
	}
}

// runServerChanDemo demonstrates how to send messages using the ServerChan provider.
// Required environment variable:
//
//	SERVERCHAN_KEY - ServerChan key
//
// This demo will send various types of messages to the specified ServerChan group.
func runServerChanDemo() {
	log.Println("[DEMO] ServerChan provider")
	key := os.Getenv("SERVERCHAN_KEY")
	if key == "" {
		log.Println("Please set SERVERCHAN_KEY environment variable.")
		return
	}

	cfg := serverchan.Config{
		ProviderMeta: core.ProviderMeta{
			Strategy: core.StrategyRoundRobin,
		},
		Items: []*serverchan.Account{{
			BaseAccount: core.BaseAccount{
				AccountMeta: core.AccountMeta{
					Name: "default",
				},
				Credentials: core.Credentials{
					APIKey: key,
				},
			},
		}},
	}

	prov, _ := serverchan.New(&cfg)
	demos := []struct {
		name string
		msg  core.Message
	}{
		{
			name: "Simple Text Message",
			msg: serverchan.Text().
				Title("服务器告警").
				Content("CPU 使用率超过 90%，请及时处理！").
				Build(),
		},
		{
			name: "Rich Markdown Message",
			msg: serverchan.Text().
				Title("每日系统报告").
				Content("# 系统状态\n\n" +
					"## 资源使用情况\n" +
					"- CPU: 45%\n" +
					"- 内存: 60%\n" +
					"- 磁盘: 75%\n\n" +
					"## 服务状态\n" +
					"| 服务名称 | 状态 | 运行时间 |\n" +
					"|---------|------|----------|\n" +
					"| nginx   | ✅    | 7天      |\n" +
					"| mysql   | ✅    | 15天     |\n" +
					"| redis   | ✅    | 3天      |\n\n" +
					"## 最近告警\n" +
					"> 过去24小时内无重要告警\n\n" +
					"详细信息请查看 [监控面板](https://monitor.example.com)").
				Build(),
		},
		{
			name: "Message with Channel (WeCom)",
			msg: serverchan.Text().
				Title("紧急通知").
				Content("数据库连接异常，请立即检查！").
				Channel("wecom"). // 发送到企业微信
				Build(),
		},
		{
			name: "Message with Channel (DingTalk)",
			msg: serverchan.Text().
				Title("服务更新通知").
				Content("系统将在今晚 23:00 进行例行维护，预计耗时 30 分钟。").
				Channel("dingtalk"). // 发送到钉钉
				Build(),
		},
		{
			name: "Message with Multiple Channels",
			msg: serverchan.Text().
				Title("重要：安全漏洞修复").
				Content("发现并修复了一个严重的安全漏洞，请所有开发人员关注。\n" +
					"影响范围：所有生产服务器\n" +
					"修复状态：已修复\n" +
								"所需操作：请更新生产环境").
				Channel("wecom|dingtalk"). // 同时发送到企业微信和钉钉
				Build(),
		},
		{
			name: "Message with Short Summary",
			msg: serverchan.Text().
				Title("性能监控报告").
				Content("## 接口性能报告\n" +
					"1. 登录接口: 95% < 100ms\n" +
					"2. 搜索接口: 95% < 200ms\n" +
					"3. 订单接口: 95% < 300ms\n\n" +
							"所有接口响应时间都在正常范围内。").
				Short("接口性能正常"). // 设置简短摘要
				Build(),
		},
	}

	for _, demo := range demos {
		log.Printf("Testing %s...\n", demo.name)
		if err := prov.Send(context.Background(), demo.msg, &core.ProviderSendOptions{
			HTTPClient: newLoggingHTTPClient(),
		}); err != nil {
			log.Printf("%s failed: %v\n", demo.name, err)
		} else {
			log.Printf("%s success!\n", demo.name)
		}
		time.Sleep(3 * time.Second) // 限流保护
	}
}
