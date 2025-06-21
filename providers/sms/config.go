package sms

import (
	"github.com/shellvon/go-sender/core"
)

// ProviderType represents different SMS service providers
type ProviderType string

const (
	ProviderTypeTencent  ProviderType = "tencent"  // 腾讯云短信
	ProviderTypeAliyun   ProviderType = "aliyun"   // 阿里云短信
	ProviderTypeYunpian  ProviderType = "yunpian"  // 云片网
	ProviderTypeUcp      ProviderType = "ucp"      // 云之讯
	ProviderTypeCl253    ProviderType = "cl253"    // 蓝创253
	ProviderTypeSmsbao   ProviderType = "smsbao"   // 短信宝
	ProviderTypeJuhe     ProviderType = "juhe"     // 聚合服务
	ProviderTypeLuosimao ProviderType = "luosimao" // 螺丝帽
	ProviderTypeMiaodi   ProviderType = "miaodi"   // 秒滴云
	ProviderTypeNetease  ProviderType = "netease"  // 网易云短信
	ProviderTypeNormal   ProviderType = "normal"   // 其他普通短信
)

// Config holds configuration for the SMS provider
type Config struct {
	Disabled          bool              `json:"disabled"`           // Whether the SMS provider is enabled
	Providers         []SMSProvider     `json:"providers"`          // Multiple SMS providers configuration
	Strategy          core.StrategyType `json:"strategy"`           // Selection strategy
	DisableMiddleware bool              `json:"disable_middleware"` // Whether to disable middleware
}

// SMSProvider represents a single SMS service provider configuration
type SMSProvider struct {
	Name     string       `json:"name"`     // Provider name for selection
	Type     ProviderType `json:"type"`     // SMS service provider type
	Weight   int          `json:"weight"`   // Weight for weighted strategy
	Disabled bool         `json:"disabled"` // Whether this provider is disabled

	// Common fields
	AppID     string `json:"app_id"`     // App ID/Account
	AppSecret string `json:"app_secret"` // App Secret/Password/Token
	SignName  string `json:"sign_name"`  // SMS signature name

	// Provider-specific fields
	TemplateCode   string            `json:"template_code"`   // Template code for template SMS
	TemplateParams map[string]string `json:"template_params"` // Template parameters
	ExtraConfig    map[string]string `json:"extra_config"`    // Extra configuration for specific providers
}

// GetStrategy returns the selection strategy, defaulting to round_robin
func (c *Config) GetStrategy() core.StrategyType {
	if c.Strategy == "" {
		return core.StrategyRoundRobin
	}
	return c.Strategy
}

// IsConfigured checks if the SMS configuration is valid
func (c Config) IsConfigured() bool {
	return !c.Disabled && len(c.Providers) > 0
}

// IsEnabled checks if the provider is enabled
func (p *SMSProvider) IsEnabled() bool {
	return !p.Disabled
}

// GetName returns the provider name for strategy selection
func (p *SMSProvider) GetName() string {
	return p.Name
}

// GetWeight returns the provider weight for strategy selection
func (p *SMSProvider) GetWeight() int {
	if p.Weight <= 0 {
		return 1
	}
	return p.Weight
}
