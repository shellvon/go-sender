package sms

import (
	"github.com/shellvon/go-sender/core"
)

// ProviderType represents different SMS service providers
type ProviderType string

const (
	ProviderTypeAliyun     ProviderType = "aliyun"     // 阿里云短信（支持国内和国际）
	ProviderTypeCl253      ProviderType = "cl253"      // 蓝创253
	ProviderTypeSmsbao     ProviderType = "smsbao"     // 短信宝
	ProviderTypeJuhe       ProviderType = "juhe"       // 聚合服务
	ProviderTypeLuosimao   ProviderType = "luosimao"   // 螺丝帽
	ProviderTypeHuawei     ProviderType = "huawei"     // 华为云短信
	ProviderTypeUcp        ProviderType = "ucp"        // 云之讯
	ProviderTypeYunpian    ProviderType = "yunpian"    // 云片短信（支持国内和国际）
	ProviderTypeSubmail    ProviderType = "submail"    // 赛邮短信（支持国内和国际）
	ProviderTypeVolc       ProviderType = "volc"       // 火山引擎短信
	ProviderTypeYuntongxun ProviderType = "yuntongxun" // 云讯通（容联云通讯）
	ProviderTypeTencent    ProviderType = "tencent"    // 腾讯云短信（支持国内和国际）
)

// Config holds configuration for the SMS provider
type Config struct {
	core.BaseConfig
	Providers []SMSProvider     `json:"providers"` // Multiple SMS providers configuration
	Strategy  core.StrategyType `json:"strategy"`  // Selection strategy
}

// SMSProvider represents a single SMS service provider configuration
type SMSProvider struct {
	Name     string       `json:"name"`     // Provider name for selection
	Type     ProviderType `json:"type"`     // SMS service provider type
	Weight   int          `json:"weight"`   // Weight for weighted strategy
	Disabled bool         `json:"disabled"` // Whether this provider is disabled

	// Common fields
	AppID        string `json:"app_id"`        // App ID/Account
	AppSecret    string `json:"app_secret"`    // App Secret/Password/Token
	Channel      string `json:"channel"`       // 通道号（如csms100000001）
	Callback     string `json:"callback"`      // 回调地址（华为等部分厂商需要）
	Endpoint     string `json:"endpoint"`      // 自定义API Endpoint（国内短信）
	IntlEndpoint string `json:"intl_endpoint"` // 自定义API Endpoint（国际短信）
}

func (c Config) IsConfigured() bool {
	return len(c.Providers) > 0
}

func (p *SMSProvider) GetName() string {
	return p.Name
}

func (p *SMSProvider) GetWeight() int {
	if p.Weight <= 0 {
		return 1
	}
	return p.Weight
}

func (p *SMSProvider) IsEnabled() bool {
	return !p.Disabled
}

// GetEndpoint 根据 isIntl 获取对应的 endpoint
// isIntl: 是否为国际短信
// defaultEndpoint: 默认的 endpoint
func (p *SMSProvider) GetEndpoint(isIntl bool, defaultEndpoint string) string {
	if isIntl {
		if p.IntlEndpoint != "" {
			return p.IntlEndpoint
		}
	}
	if p.Endpoint != "" {
		return p.Endpoint
	}
	return defaultEndpoint
}
