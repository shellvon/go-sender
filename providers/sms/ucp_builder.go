package sms

// @ProviderName: UCP / 云之讯
// @Website: https://www.ucpaas.com
// @APIDoc: http://docs.ucpaas.com
//
// 官方文档:
//   - 短信API: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:about_sms
//
// builder 仅支持 text（普通短信）类型。

type UcpSMSBuilder struct {
	BaseSMSBuilder
}

func newUcpSMSBuilder() *UcpSMSBuilder {
	return &UcpSMSBuilder{
		BaseSMSBuilder: BaseSMSBuilder{subProvider: SubProviderUcp},
	}
}
