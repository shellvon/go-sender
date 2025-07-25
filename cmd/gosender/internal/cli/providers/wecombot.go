package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/wecombot"
)

// createWeComBotProvider 从账户列表创建 WeComBot Provider
func createWeComBotProvider(accounts []*wecombot.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid wecombot accounts found")
	}

	cfg := &wecombot.Config{Items: accounts}
	return wecombot.New(cfg)
}

// buildWeComBotMessage 从 CLI 标志构建企业微信机器人消息
func buildWeComBotMessage(flags *cli.CLIFlags) (core.Message, error) {
	messageType := flags.MessageType
	if messageType == "" {
		messageType = "text" // Default to text
	}

	switch messageType {
	case "text":
		return wecombot.Text().Content(flags.Content).Build(), nil
	case "markdown":
		return wecombot.Markdown().Content(flags.Content).Build(), nil
	case "image":
		if len(flags.Files) == 0 {
			return nil, fmt.Errorf("image message requires a file")
		}
		return wecombot.Image().Build(), nil
	case "news":
		return wecombot.News().Build(), nil
	case "file":
		if len(flags.Files) == 0 {
			return nil, fmt.Errorf("file message requires a file")
		}
		return wecombot.File().Build(), nil
	case "template_card":
		return wecombot.Card(wecombot.CardTypeTextNotice).
			MainTitle("Template Card", "This is a template card message").
			SubTitle(flags.Content).
			Build(), nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}
}

// validateWeComBotFlags 验证 CLI 标志是否符合企业微信机器人发送要求
func validateWeComBotFlags(flags *cli.CLIFlags) error {
	if flags.Content == "" {
		return fmt.Errorf("wecombot requires content")
	}

	if flags.MessageType != "" {
		validTypes := []string{"text", "markdown", "image", "news", "file", "template_card"}
		valid := false
		for _, validType := range validTypes {
			if flags.MessageType == validType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid message type '%s' for wecombot, supported types: %v", flags.MessageType, validTypes)
		}
	}

	if flags.MessageType == "image" || flags.MessageType == "file" {
		if len(flags.Files) == 0 {
			return fmt.Errorf("message type '%s' requires at least one file", flags.MessageType)
		}
	}

	return nil
}

// NewWeComBotBuilder 创建一个新的 WeComBot GenericBuilder
func NewWeComBotBuilder() *GenericBuilder[*wecombot.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeWecombot,
		createWeComBotProvider,
		buildWeComBotMessage,
		validateWeComBotFlags,
	)
}
