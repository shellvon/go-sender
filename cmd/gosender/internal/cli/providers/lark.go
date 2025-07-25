package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/lark"
)

// createLarkProvider 从账户列表创建 Lark Provider
func createLarkProvider(accounts []*lark.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid lark accounts found")
	}

	cfg := &lark.Config{Items: accounts}
	return lark.New(cfg)
}

// buildLarkMessage 从 CLI 标志构建 Lark 消息
func buildLarkMessage(flags *cli.CLIFlags) (core.Message, error) {
	messageType := flags.MessageType
	if messageType == "" {
		messageType = "text" // Default to text
	}

	switch messageType {
	case "text":
		// Lark text message
		return lark.Text().Content(flags.Content).Build(), nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s for lark, currently only 'text' is supported", messageType)
	}
}

// validateLarkFlags 验证 CLI 标志是否符合 Lark 发送要求
func validateLarkFlags(flags *cli.CLIFlags) error {
	if flags.Content == "" {
		return fmt.Errorf("lark requires content")
	}

	if flags.MessageType != "" && flags.MessageType != "text" {
		return fmt.Errorf("invalid message type '%s' for lark, only 'text' is currently supported", flags.MessageType)
	}

	return nil
}

// NewLarkBuilder 创建一个新的 Lark GenericBuilder
func NewLarkBuilder() *GenericBuilder[*lark.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeLark,
		createLarkProvider,
		buildLarkMessage,
		validateLarkFlags,
	)
}
