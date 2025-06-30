package sms

import (
	"github.com/shellvon/go-sender/core"
)

// ProviderType represents different SMS service providers.
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

// Config holds configuration for the SMS provider.
type Config struct {
	core.BaseConfig

	Providers []ProviderConfig  `json:"providers"` // Multiple SMS providers configuration
	Strategy  core.StrategyType `json:"strategy"`  // Selection strategy
}

// ProviderConfig holds configuration for an SMS provider.
type ProviderConfig struct {
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Key          string                 `json:"key"`
	Secret       string                 `json:"secret"`
	SignName     string                 `json:"sign_name"`
	Endpoint     string                 `json:"endpoint"`
	IntlEndpoint string                 `json:"intl_endpoint,omitempty"`
	Weight       int                    `json:"weight,omitempty"`
	Disabled     bool                   `json:"disabled,omitempty"`
	AppID        string                 `json:"app_id,omitempty"`
	AppSecret    string                 `json:"app_secret,omitempty"`
	Channel      string                 `json:"channel,omitempty"`
	Callback     string                 `json:"callback,omitempty"`
	Extras       map[string]interface{} `json:"extras,omitempty"`
}

func (c Config) IsConfigured() bool {
	return len(c.Providers) > 0
}

func (p *ProviderConfig) GetName() string {
	return p.Name
}

func (p *ProviderConfig) GetWeight() int {
	if p.Weight <= 0 {
		return 1
	}
	return p.Weight
}

func (p *ProviderConfig) IsEnabled() bool {
	return !p.Disabled
}

// GetEndpoint 根据 isIntl 获取对应的 endpoint
// isIntl: 是否为国际短信
// defaultEndpoint: 默认的 endpoint.
func (p *ProviderConfig) GetEndpoint(isIntl bool, defaultEndpoint string) string {
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
