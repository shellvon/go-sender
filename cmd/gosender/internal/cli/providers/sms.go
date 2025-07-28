package providers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/sms"
)

// GetSMSSubProviders 返回CLI支持的SMS子提供者列表
func GetSMSSubProviders() []string {
	return []string{
		string(sms.SubProviderAliyun),
		string(sms.SubProviderTencent),
		string(sms.SubProviderHuawei),
		string(sms.SubProviderCl253),
		string(sms.SubProviderJuhe),
		string(sms.SubProviderLuosimao),
		string(sms.SubProviderSmsbao),
		string(sms.SubProviderSubmail),
		string(sms.SubProviderUcp),
		string(sms.SubProviderVolc),
		string(sms.SubProviderYuntongxun),
		string(sms.SubProviderYunpian),
	}
}

// createSMSProvider 从账户列表创建 SMS Provider
func createSMSProvider(accounts []*sms.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid sms accounts found")
	}

	cfg := &sms.Config{Items: accounts}
	return sms.New(cfg)
}

// buildSMSMessage 从 CLI 标志构建 SMS 消息
func buildSMSMessage(flags *cli.CLIFlags) (core.Message, error) {
	// 检查子提供者是否受支持
	if !isSMSSubProviderSupported(flags.SubProvider) {
		return nil, fmt.Errorf("unsupported sub-provider: %s. Supported providers: %s",
			flags.SubProvider, strings.Join(GetSMSSubProviders(), ", "))
	}

	// 直接创建消息对象
	msg := &sms.Message{
		SubProvider:    flags.SubProvider,
		Mobiles:        flags.To,
		Content:        flags.Content,
		TemplateID:     flags.TemplateID,
		TemplateParams: flags.TemplateParams,
		// 基于MessageType设置类型
		Type: getMessageType(flags.MessageType),
		// 从Metadata提取常用字段
		SignName:    flags.Metadata["sign_name"],
		RegionCode:  getIntValue(flags.Metadata["region_code"]),
		CallbackURL: flags.Metadata["callback_url"],
		Extend:      flags.Metadata["extend"],
		UID:         flags.Metadata["uid"],
	}

	// 解析参数顺序
	msg.ParamsOrder = getJSONArray(flags.Metadata["params_order"])

	return msg, nil
}

func getMessageType(messageType string) sms.MessageType {
	types := map[string]sms.MessageType{
		"voice": sms.Voice,
		"mms":   sms.MMS,
		"text":  sms.SMSText,
	}
	return types[messageType]
}

// getIntValue 安全地将字符串转换为整数
func getIntValue(s string) int {
	value, _ := strconv.Atoi(s) // 忽略错误，错误时返回0
	return value
}

// getJSONArray 解析JSON数组字符串
func getJSONArray(s string) []string {
	var result []string
	_ = json.Unmarshal([]byte(s), &result) // 忽略错误，错误时返回空数组
	return result
}

// isSMSSubProviderSupported 检查子提供者是否受支持
func isSMSSubProviderSupported(subProvider string) bool {
	for _, p := range GetSMSSubProviders() {
		if p == subProvider {
			return true
		}
	}
	return false
}

// validateSMSFlags 验证 CLI 标志是否符合 SMS 发送要求
func validateSMSFlags(flags *cli.CLIFlags) error {
	if len(flags.To) == 0 {
		return fmt.Errorf("sms requires at least one recipient (--to)")
	}

	// 验证子提供者
	if flags.SubProvider == "" {
		return fmt.Errorf("sub-provider is required for SMS, please use --sub-provider (aliyun, tencent, etc.)")
	}

	// 验证子提供者是否受支持
	if !isSMSSubProviderSupported(flags.SubProvider) {
		return fmt.Errorf("unsupported sub-provider: %s. Supported providers: %s",
			flags.SubProvider, strings.Join(GetSMSSubProviders(), ", "))
	}

	// 验证电话号码格式
	for _, phone := range flags.To {
		if !isValidPhoneNumber(phone) {
			return fmt.Errorf("invalid phone number format: %s", phone)
		}
	}

	// 验证消息类型
	if flags.MessageType != "" &&
		flags.MessageType != "text" &&
		flags.MessageType != "voice" &&
		flags.MessageType != "mms" {
		return fmt.Errorf("invalid message type '%s' for SMS, supported types: text, voice, mms", flags.MessageType)
	}

	return nil
}

// isValidPhoneNumber 简单验证电话号码格式
func isValidPhoneNumber(phone string) bool {
	// 这里提供一个简单的验证，实际应用中可能需要更复杂的逻辑
	return len(phone) >= 7 && len(phone) <= 15
}

// NewSMSBuilder 创建一个新的 SMS GenericBuilder
func NewSMSBuilder() *GenericBuilder[*sms.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeSMS,
		createSMSProvider,
		buildSMSMessage,
		validateSMSFlags,
	)
}
