package sms

// @ProviderName: Volc / 火山引擎
// @Website: https://www.volcengine.com
// @APIDoc: https://www.volcengine.com/docs/6348/70146
//
// 官方文档:
//   - 短信API: https://www.volcengine.com/docs/6348/70146
//
// builder 仅支持 text（普通短信）类型。

type VolcSMSBuilder struct {
	*BaseBuilder

	tag string
}

func newVolcSMSBuilder() *VolcSMSBuilder {
	return &VolcSMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderVolc},
	}
}

// Tag 设置火山引擎短信标签。
// 用于标识短信业务类型，便于后续统计分析和管理。
//   - 短信API: https://www.volcengine.com/docs/6348/70146
func (b *VolcSMSBuilder) Tag(tag string) *VolcSMSBuilder {
	b.tag = tag
	return b
}

func (b *VolcSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	// 火山引擎专属参数写入Extras
	if b.tag != "" {
		if msg.Extras == nil {
			msg.Extras = make(map[string]interface{})
		}
		msg.Extras[volcTagKey] = b.tag
	}
	return msg
}
