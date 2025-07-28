package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/dingtalk"
)

// createDingTalkProvider 从账户列表创建 DingTalk Provider
func createDingTalkProvider(accounts []*dingtalk.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid dingtalk accounts found")
	}

	cfg := &dingtalk.Config{Items: accounts}
	return dingtalk.New(cfg)
}

// buildDingTalkMessage 从 CLI 标志构建 DingTalk 消息
func buildDingTalkMessage(flags *cli.CLIFlags) (core.Message, error) {
	messageType := flags.MessageType
	if messageType == "" {
		messageType = "text" // Default to text
	}

	switch messageType {
	case "text":
		// DingTalk text message
		builder := dingtalk.Text().Content(flags.Content)
		// You could potentially parse AtMobiles or AtUserIDs from flags.To or flags.Metadata
		return builder.Build(), nil
	case "markdown":
		// DingTalk markdown message
		title := flags.Subject
		if title == "" {
			title = "Message from gosender"
		}
		return dingtalk.Markdown().Title(title).Text(flags.Content).Build(), nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s for dingtalk, currently only 'text' and 'markdown' are supported", messageType)
	}
}

// validateDingTalkFlags 验证 CLI 标志是否符合 DingTalk 发送要求
func validateDingTalkFlags(flags *cli.CLIFlags) error {
	if flags.Content == "" {
		return fmt.Errorf("dingtalk requires content")
	}

	if flags.MessageType != "" && flags.MessageType != "text" && flags.MessageType != "markdown" {
		return fmt.Errorf("invalid message type '%s' for dingtalk, only 'text' and 'markdown' are currently supported", flags.MessageType)
	}

	return nil
}

// NewDingTalkBuilder 创建一个新的 DingTalk GenericBuilder
func NewDingTalkBuilder() *GenericBuilder[*dingtalk.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeDingtalk,
		createDingTalkProvider,
		buildDingTalkMessage,
		validateDingTalkFlags,
	)
}
 