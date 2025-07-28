package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/sms"
)

// GetSMSSubProviders 返回CLI支持的SMS子提供者列表.
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

// createSMSProvider 从账户列表创建 SMS Provider.
func createSMSProvider(accounts []*sms.Account) (core.Provider, error) {
	if len(accounts) == 0 {
		return nil, errors.New("no valid sms accounts found")
	}

	cfg := &sms.Config{Items: accounts}
	return sms.New(cfg)
}

// buildSMSMessage 从 CLI 标志构建 SMS 消息.
func buildSMSMessage(flags *cli.CLIFlags) (core.Message, error) {
	// 检查子提供者是否受支持
	if !isSMSSubProviderSupported(flags.SubProvider) {
		return nil, fmt.Errorf("unsupported sub-provider: %s. Supported providers: %s",
			flags.SubProvider, strings.Join(GetSMSSubProviders(), ", "))
	}

	// 使用Generic Builder API构建消息
	builder := sms.NewMessage().
		SubProvider(sms.SubProviderType(flags.SubProvider)).
		To(flags.To...).
		Content(flags.Content)

	// 添加模板ID和参数(如果有)
	if flags.TemplateID != "" {
		if len(flags.TemplateParams) > 0 {
			builder.Template(flags.TemplateID, flags.TemplateParams)
		} else {
			builder.TemplateID(flags.TemplateID)
		}
	}

	// 设置消息类型
	if flags.MessageType != "" {
		builder.Type(getMessageType(flags.MessageType))
	}

	// 添加常用元数据字段
	if signName := flags.Metadata["sign_name"]; signName != "" {
		builder.SignName(signName)
	}

	if regionCode := getIntValue(flags.Metadata["region_code"]); regionCode > 0 {
		builder.RegionCode(regionCode)
	}

	if callbackURL := flags.Metadata["callback_url"]; callbackURL != "" {
		builder.CallbackURL(callbackURL)
	}

	if extend := flags.Metadata["extend"]; extend != "" {
		builder.Extend(extend)
	}

	if uid := flags.Metadata["uid"]; uid != "" {
		builder.UID(uid)
	}

	// 处理参数顺序
	if paramsOrder := getJSONArray(flags.Metadata["params_order"]); len(paramsOrder) > 0 {
		builder.ParamsOrder(paramsOrder)
	}

	// 添加其他元数据
	for key, value := range flags.Metadata {
		// 跳过已处理的常用字段
		if key == "sign_name" || key == "region_code" || key == "callback_url" ||
			key == "extend" || key == "uid" || key == "params_order" {
			continue
		}
		builder.Meta(key, value)
	}

	return builder.Build(), nil
}

func getMessageType(messageType string) sms.MessageType {
	types := map[string]sms.MessageType{
		"voice": sms.Voice,
		"mms":   sms.MMS,
		"text":  sms.SMSText,
	}
	return types[messageType]
}

// getIntValue 安全地将字符串转换为整数.
func getIntValue(s string) int {
	value, _ := strconv.Atoi(s) // 忽略错误，错误时返回0
	return value
}

// getJSONArray 解析JSON数组字符串.
func getJSONArray(s string) []string {
	var result []string
	_ = json.Unmarshal([]byte(s), &result) // 忽略错误，错误时返回空数组
	return result
}

// isSMSSubProviderSupported 检查子提供者是否受支持.
func isSMSSubProviderSupported(subProvider string) bool {
	for _, p := range GetSMSSubProviders() {
		if p == subProvider {
			return true
		}
	}
	return false
}

// validateSMSFlags 验证 CLI 标志是否符合 SMS 发送要求.
func validateSMSFlags(flags *cli.CLIFlags) error {
	if len(flags.To) == 0 {
		return errors.New("sms requires at least one recipient (--to)")
	}

	// 验证子提供者
	if flags.SubProvider == "" {
		return errors.New("sub-provider is required for SMS, please use --sub-provider (aliyun, tencent, etc.)")
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

// isValidPhoneNumber 简单验证电话号码格式.
func isValidPhoneNumber(phone string) bool {
	// 这里提供一个简单的验证，实际应用中可能需要更复杂的逻辑
	return len(phone) >= 7 && len(phone) <= 15
}

// NewSMSBuilder 创建一个新的 SMS GenericBuilder.
func NewSMSBuilder() *GenericBuilder[*sms.Account, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeSMS,
		createSMSProvider,
		buildSMSMessage,
		validateSMSFlags,
	)
}
