package sms

// @ProviderName: Juhe / 聚合数据
// @Website: https://www.juhe.cn
// @APIDoc: https://www.juhe.cn/docs/api/id/54
//
// 官方文档:
//   - 短信API文档: https://www.juhe.cn/docs/api/id/54
//   - 国内短信API: https://www.juhe.cn/docs/api/id/54
//   - 国际短信API: https://www.juhe.cn/docs/api/id/357
//   - 视频短信API: https://www.juhe.cn/docs/api/id/363
//
// builder 仅支持 text（普通短信）类型。

// JuheSMSBuilder provides Juhe-specific SMS message creation.
type JuheSMSBuilder struct {
	BaseSMSBuilder
}

// newJuheSMSBuilder creates a new Juhe SMS builder.
func newJuheSMSBuilder() *JuheSMSBuilder {
	return &JuheSMSBuilder{
		BaseSMSBuilder: BaseSMSBuilder{subProvider: SubProviderJuhe},
	}
}
