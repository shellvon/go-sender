package providers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/wecomapp"
)

// createWeComAppProvider 从账户列表创建 WeComApp Provider.
func createWeComAppProvider(accounts []*wecomapp.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, errors.New("no valid wecomapp accounts found")
	}

	return wecomapp.NewProvider(accounts)
}

// buildWeComAppMessage 从 CLI 标志构建企业微信应用消息.
func buildWeComAppMessage(flags *cli.CLIFlags) (core.Message, error) {
	messageType := flags.MessageType
	if messageType == "" {
		messageType = "text" // Default to text
	}

	toUsers := strings.Join(flags.To, "|")

	fmt.Sprintln("--->", toUsers)
	switch messageType {
	case "text":
		return wecomapp.Text().Content(flags.Content).ToUser(toUsers).Build(), nil
	case "markdown":
		return wecomapp.Markdown().Content(flags.Content).ToUser(toUsers).Build(), nil
	case "image":
		if len(flags.Files) == 0 {
			return nil, errors.New("image message requires a file")
		}
		return wecomapp.Image().LocalPath(flags.Files[0]).ToUser(toUsers).Build(), nil
	case "voice":
		if len(flags.Files) == 0 {
			return nil, errors.New("voice message requires a file")
		}
		return wecomapp.Voice().LocalPath(flags.Files[0]).ToUser(toUsers).Build(), nil
	case "video":
		if len(flags.Files) == 0 {
			return nil, errors.New("video message requires a file")
		}
		return wecomapp.Video().LocalPath(flags.Files[0]).ToUser(toUsers).Build(), nil
	case "file":
		if len(flags.Files) == 0 {
			return nil, errors.New("file message requires a file")
		}
		return wecomapp.File().LocalPath(flags.Files[0]).ToUser(toUsers).Build(), nil
	case "news":
		return wecomapp.News().Build(), nil
	case "template_card":
		return wecomapp.NewTemplateCardBuilder(wecomapp.CardTypeTextNotice).
			MainTitle("Template Card", flags.Content).ToUser(toUsers).
			Build(), nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}
}

// validateWeComAppFlags 验证 CLI 标志是否符合企业微信应用发送要求.
func validateWeComAppFlags(flags *cli.CLIFlags) error {
	if flags.Content == "" && flags.MessageType != "news" && flags.MessageType != "template_card" {
		return errors.New("wecomapp requires content for most message types")
	}

	if flags.MessageType != "" {
		validTypes := []string{"text", "markdown", "image", "voice", "video", "file", "news", "template_card"}
		valid := false
		for _, validType := range validTypes {
			if flags.MessageType == validType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf(
				"invalid message type '%s' for wecomapp, supported types: %v",
				flags.MessageType,
				validTypes,
			)
		}
	}

	// File-based message types require files
	fileBasedTypes := []string{"image", "voice", "video", "file"}
	for _, fileType := range fileBasedTypes {
		if flags.MessageType == fileType && len(flags.Files) == 0 {
			return fmt.Errorf("message type '%s' requires at least one file", flags.MessageType)
		}
	}

	return nil
}

// NewWeComAppBuilder 创建一个新的 WeComApp GenericBuilder.
func NewWeComAppBuilder() *GenericBuilder[*wecomapp.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeWecomApp,
		createWeComAppProvider,
		buildWeComAppMessage,
		validateWeComAppFlags,
	)
}
