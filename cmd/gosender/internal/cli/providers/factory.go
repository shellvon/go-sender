package providers

import (
	"fmt"
	"log"

	"github.com/shellvon/go-sender/core"
)

func DefaultProviderRegistry() *ProviderRegistry {
	registry := NewProviderRegistry()

	builders := []ProviderFactory{
		NewEmailBuilder(),      // 邮件服务
		NewWeComBotBuilder(),   // 企业微信机器人
		NewServerChanBuilder(), // ServerChan
	}

	for _, builder := range builders {
		if err := registry.RegisterFactory(builder); err != nil {
			log.Printf("Warning: failed to register builder for provider %s: %v",
				builder.GetProviderType(), err)
		}
	}

	return registry
}

// GetProviderType 将字符串转换为 ProviderType 并进行验证
func GetProviderType(provider string) (core.ProviderType, error) {
	if provider == "" {
		return core.ProviderTypeSMS, nil // 默认
	}

	providerType := core.ProviderType(provider)

	// 验证是否为已知的 provider 类型
	switch providerType {
	case core.ProviderTypeSMS,
		core.ProviderTypeEmail,
		core.ProviderTypeDingtalk,
		core.ProviderTypeWebhook,
		core.ProviderTypeTelegram,
		core.ProviderTypeLark,
		core.ProviderTypeWecombot,
		core.ProviderTypeServerChan,
		core.ProviderTypeEmailAPI:
		return providerType, nil
	default:
		return "", fmt.Errorf("unsupported provider type: %s", provider)
	}
}
