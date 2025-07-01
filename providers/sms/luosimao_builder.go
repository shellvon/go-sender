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
	BaseSMSBuilder
}

// newLuosimaoSMSBuilder creates a new Luosimao SMS builder.
func newLuosimaoSMSBuilder() *LuosimaoSMSBuilder {
	return &LuosimaoSMSBuilder{
		BaseSMSBuilder: BaseSMSBuilder{subProvider: SubProviderLuosimao},
	}
}
