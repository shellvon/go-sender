package providers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/emailapi"
)

// GetEmailAPISubProviders 返回支持的EmailAPI子提供者列表.
func GetEmailAPISubProviders() []string {
	return []string{
		string(emailapi.SubProviderEmailJS),
		string(emailapi.SubProviderResend),
	}
}

func createEmailAPIProvider(accounts []*emailapi.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, errors.New("no valid emailapi accounts found")
	}

	cfg := &emailapi.Config{Items: accounts}
	return emailapi.New(cfg)
}

func buildEmailAPIMessage(flags *cli.CLIFlags) (core.Message, error) {
	if len(flags.To) == 0 {
		return nil, errors.New("emailapi requires at least one recipient email address")
	}

	subProvider := flags.SubProvider
	if subProvider == "" {
		return nil, errors.New(
			"sub-provider is required for EmailAPI, please use --sub-provider (resend, emailjs, etc.)",
		)
	}

	// 目前EmailAPI可能没有像SMS那样的构建器体系，先直接创建消息
	message := &emailapi.Message{
		SubProvider: subProvider,
		To:          flags.To,
		Subject:     flags.Subject,
	}

	messageType := flags.MessageType
	if messageType == "" {
		if flags.HTML {
			messageType = "html"
		} else {
			messageType = "text"
		}
	}

	switch messageType {
	case "text":
		message.Text = flags.Content
	case "html":
		message.HTML = flags.Content
	default:
		return nil, fmt.Errorf("unsupported message type: %s for EmailAPI, supported types: text, html", messageType)
	}

	// 如果提供了模板ID，作为额外参数
	if flags.TemplateID != "" {
		message.TemplateID = flags.TemplateID
		// 转换模板参数
		templateData := make(map[string]interface{})
		for k, v := range flags.TemplateParams {
			templateData[k] = v
		}
		message.TemplateData = templateData
	}

	return message, nil
}

func validateEmailAPIFlags(flags *cli.CLIFlags) error {
	if len(flags.To) == 0 {
		return errors.New("emailapi requires at least one recipient (--to)")
	}

	// 验证子提供者
	if flags.SubProvider == "" {
		return errors.New("sub-provider is required for EmailAPI, please use --sub-provider (resend, emailjs, etc.)")
	}

	// 验证子提供者类型
	supportedSubProviders := GetEmailAPISubProviders()
	supported := false
	for _, sp := range supportedSubProviders {
		if flags.SubProvider == sp {
			supported = true
			break
		}
	}

	if !supported {
		return fmt.Errorf("unknown sub-provider: %s for EmailAPI, supported sub-providers: %s",
			flags.SubProvider,
			strings.Join(supportedSubProviders, ", "))
	}

	// 验证消息类型
	if flags.MessageType != "" &&
		flags.MessageType != "text" &&
		flags.MessageType != "html" {
		return fmt.Errorf("invalid message type '%s' for EmailAPI, supported types: text, html", flags.MessageType)
	}

	// 验证内容
	if flags.Content == "" {
		return errors.New("content is required for EmailAPI messages")
	}

	// 验证主题
	if flags.Subject == "" {
		return errors.New("subject is required for EmailAPI messages")
	}

	return nil
}

// NewEmailAPIBuilder 创建一个新的 EmailAPI GenericBuilder.
func NewEmailAPIBuilder() *GenericBuilder[*emailapi.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeEmailAPI,
		createEmailAPIProvider,
		buildEmailAPIMessage,
		validateEmailAPIFlags,
	)
}
