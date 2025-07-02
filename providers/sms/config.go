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

	Accounts []*core.Account   `json:"providers"` // Multiple SMS providers configuration
	Strategy core.StrategyType `json:"strategy"`  // Selection strategy
}
