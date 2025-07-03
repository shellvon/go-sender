package sms

// @ProviderName: Luosimao / 螺丝帽
// @Website: https://luosimao.com
// @APIDoc: https://luosimao.com/docs/api
//
// 官方文档:
//   - 短信API: https://luosimao.com/docs/api
//
// builder 仅支持 text（普通短信）类型。

// LuosimaoSMSBuilder provides Luosimao-specific SMS message creation.
type LuosimaoSMSBuilder struct {
	*BaseBuilder[*LuosimaoSMSBuilder]
}

// newLuosimaoSMSBuilder creates a new Luosimao SMS builder.
func newLuosimaoSMSBuilder() *LuosimaoSMSBuilder {
	b := &LuosimaoSMSBuilder{}
	b.BaseBuilder = &BaseBuilder[*LuosimaoSMSBuilder]{subProvider: SubProviderLuosimao, self: b}
	return b
}
