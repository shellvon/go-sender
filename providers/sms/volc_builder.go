package sms

import "github.com/shellvon/go-sender/utils"

// @ProviderName: Volc / 火山引擎
// @Website: https://www.volcengine.com
// @APIDoc: https://www.volcengine.com/docs/6361/67380
//
// 官方文档:
//   - 短信API: https://www.volcengine.com/docs/6361/67380
//
// builder 仅支持 text（普通短信）类型。

type VolcSMSBuilder struct {
	*BaseBuilder[*VolcSMSBuilder]

	tag string
}

func newVolcSMSBuilder() *VolcSMSBuilder {
	b := &VolcSMSBuilder{}
	b.BaseBuilder = &BaseBuilder[*VolcSMSBuilder]{subProvider: SubProviderVolc, self: b}
	return b
}

// Tag 设置火山引擎短信标签。
// 用于标识短信业务类型，便于后续统计分析和管理。
//   - 短信API: https://www.volcengine.com/docs/6361/67380
func (b *VolcSMSBuilder) Tag(tag string) *VolcSMSBuilder {
	b.tag = tag
	return b
}

func (b *VolcSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	fields := map[string]interface{}{
		volcTagKey: b.tag,
	}
	if extra := utils.BuildExtras(fields); len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}
