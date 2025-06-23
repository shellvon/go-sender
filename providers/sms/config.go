package sms

import (
	"github.com/shellvon/go-sender/core"
)

// ProviderType represents different SMS service providers
type ProviderType string

const (
	ProviderTypeAliyun      ProviderType = "aliyun"       // 阿里云短信
	ProviderTypeAliyunGlobe ProviderType = "aliyun_globe" // 阿里云国际短信
	ProviderTypeCl253       ProviderType = "cl253"        // 蓝创253
	ProviderTypeSmsbao      ProviderType = "smsbao"       // 短信宝
	ProviderTypeJuhe        ProviderType = "juhe"         // 聚合服务
	ProviderTypeLuosimao    ProviderType = "luosimao"     // 螺丝帽
	ProviderTypeHuawei      ProviderType = "huawei"       // 华为云短信
	ProviderTypeUcp         ProviderType = "ucp"          // 云之讯
)

// Config holds configuration for the SMS provider
type Config struct {
	core.BaseConfig
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
	Channel   string `json:"channel"`    // 通道号（如csms100000001）
	Callback  string `json:"callback"`   // 回调地址（华为等部分厂商需要）
	Endpoint  string `json:"endpoint"`   // 自定义API Endpoint（如华为国际/国内不同节点）
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
