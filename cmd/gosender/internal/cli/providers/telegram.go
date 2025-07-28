package providers

import (
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/telegram"
)

// createTelegramProvider 从账户列表创建 Telegram Provider.
func createTelegramProvider(accounts []*telegram.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, errors.New("no valid telegram accounts found")
	}

	cfg := &telegram.Config{Items: accounts}
	return telegram.New(cfg)
}

// buildTelegramMessage 从 CLI 标志构建 Telegram 消息.
func buildTelegramMessage(flags *cli.CLIFlags) (core.Message, error) {
	messageType := flags.MessageType
	if messageType == "" {
		messageType = "text" // Default to text
	}

	if len(flags.To) == 0 {
		return nil, errors.New("at least one recipient (chat_id) is required for telegram messages")
	}

	// Use the first recipient as the chat ID
	chatID := flags.To[0]

	switch messageType {
	case "text":
		return telegram.Text().Chat(chatID).Text(flags.Content).Build(), nil
	default:
		return nil, fmt.Errorf(
			"unsupported message type: %s for telegram, currently only 'text' is supported",
			messageType,
		)
	}
}

// validateTelegramFlags 验证 CLI 标志是否符合 Telegram 发送要求.
func validateTelegramFlags(flags *cli.CLIFlags) error {
	if len(flags.To) == 0 {
		return errors.New("telegram requires at least one recipient (chat_id)")
	}

	if flags.Content == "" {
		return errors.New("telegram requires content")
	}

	if flags.MessageType != "" && flags.MessageType != "text" {
		return fmt.Errorf(
			"invalid message type '%s' for telegram, only 'text' is currently supported",
			flags.MessageType,
		)
	}

	return nil
}

// NewTelegramBuilder 创建一个新的 Telegram GenericBuilder.
func NewTelegramBuilder() *GenericBuilder[*telegram.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeTelegram,
		createTelegramProvider,
		buildTelegramMessage,
		validateTelegramFlags,
	)
}
