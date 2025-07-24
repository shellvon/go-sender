package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/wecombot"
)

// WeComBotProviderBuilder builds wecombot providers
type WeComBotProviderBuilder struct{}

func (b *WeComBotProviderBuilder) GetProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}

func (b *WeComBotProviderBuilder) BuildProvider(accounts []any) (core.Provider, error) {
	parser := config.NewAccountParser()
	var wecomAccounts []*wecombot.Account

	for i, rawAccount := range accounts {
		accountMap, ok := rawAccount.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid account format at index %d", i)
		}

		account, err := parser.ParseWeComBotAccount(accountMap)
		if err != nil {
			return nil, fmt.Errorf("failed to parse wecombot account at index %d: %w", i, err)
		}

		wecomAccounts = append(wecomAccounts, account)
	}

	if len(wecomAccounts) == 0 {
		return nil, fmt.Errorf("no valid wecombot accounts found")
	}

	cfg := &wecombot.Config{Items: wecomAccounts}
	return wecombot.New(cfg)
}

// WeComBotMessageBuilder builds wecombot messages
type WeComBotMessageBuilder struct{}

func (b *WeComBotMessageBuilder) GetProviderType() core.ProviderType {
	return core.ProviderTypeWecombot
}

func (b *WeComBotMessageBuilder) ValidateFlags(flags *cli.CLIFlags) error {
	// WeComBot doesn't need --to parameter as it sends to groups
	if flags.Content == "" {
		return fmt.Errorf("wecombot requires content")
	}

	// Validate message type if specified
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

	// File-based message types require files
	if flags.MessageType == "image" || flags.MessageType == "file" {
		if len(flags.Files) == 0 {
			return fmt.Errorf("message type '%s' requires at least one file", flags.MessageType)
		}
	}

	return nil
}

func (b *WeComBotMessageBuilder) BuildMessage(flags *cli.CLIFlags) (core.Message, error) {
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
		// For image messages, we need to handle file upload
		// This might need enhancement in the wecombot package
		return wecombot.Image().Build(), nil

	case "news":
		// News messages typically need structured data
		// For now, we'll create a simple news item with the content
		return wecombot.News().Build(), nil

	case "file":
		if len(flags.Files) == 0 {
			return nil, fmt.Errorf("file message requires a file")
		}
		// File messages need file handling
		return wecombot.File().Build(), nil

	case "template_card":
		// Template card messages need structured data
		return wecombot.Card(wecombot.CardTypeTextNotice).
			MainTitle("Template Card", "This is a template card message").
			SubTitle(flags.Content).
			Build(), nil

	default:
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}
}
