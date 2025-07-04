//nolint:dupl // Huawei builder 和 Volc builder 结构类似但业务独立，重复为误报。
package sms

import "github.com/shellvon/go-sender/utils"

// @ProviderName: Huawei Cloud / 华为云
// @Website: https://www.huaweicloud.com
// @APIDoc: https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 官方文档:
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 仅支持国内/国际短信，不支持彩信和语音。
//
// builder 仅支持 text（模板短信）类型。

// HuaweiSMSBuilder provides Huawei-specific SMS message creation.
type HuaweiSMSBuilder struct {
	*BaseBuilder[*HuaweiSMSBuilder]

	from   string
	region string
}

// newHuaweiSMSBuilder creates a new Huawei SMS builder.
func newHuaweiSMSBuilder() *HuaweiSMSBuilder {
	b := &HuaweiSMSBuilder{}
	b.BaseBuilder = &BaseBuilder[*HuaweiSMSBuilder]{subProvider: SubProviderHuawei, self: b}
	return b
}

// From sets the sender number for Huawei SMS.
// 华为云短信服务中，from 字段用于指定发送方号码。
//   - 文档地址: https://support.huaweicloud.com/api-msgsms/sms_05_0002.html
func (b *HuaweiSMSBuilder) From(from string) *HuaweiSMSBuilder {
	b.from = from
	return b
}

// Region 设置华为云短信的区域, 默认为cn-north-1
//
// 文档地址:  https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0000.html
//
// 该值影响当次发放的服务接入点
// API请求地址不区分区域、省份，华为云-国际站-短信服务使用同一地址。API请求地址由"APP接入地址"和"访问URI"组成，数据来源如下：
//   - APP接入地址：华为云-国际站-短信服务使用同一地址。
//   - 访问URI：华为云-国际站-短信服务使用同一地址。
//
// 对于API接入地址，您需要管理控制台，从全球短信"应用管理"或中国大陆短信"应用管理"页面获取。
// 则最终请求地址是 https://rtcsms.${region}.myhuaweicloud.com/sms/batchSendSms/v1
func (b *HuaweiSMSBuilder) Region(region string) *HuaweiSMSBuilder {
	b.region = region
	return b
}

// Build constructs the Huawei SMS message with all the configured options.
func (b *HuaweiSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	fields := map[string]interface{}{
		huaweiFromKey:   b.from,
		huaweiRegionKey: b.region,
	}
	if extra := utils.BuildExtras(fields); len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}
