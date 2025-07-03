package sms

// @ProviderName: Smsbao / 短信宝
// @Website: https://www.smsbao.com
// @APIDoc: https://www.smsbao.com/openapi
//
// 官方文档:
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//   - 国际短信: https://www.smsbao.com/openapi/299.html
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
//
// builder 支持 text、voice、intl 短信。

type SmsbaoSMSBuilder struct {
	*BaseBuilder[*SmsbaoSMSBuilder]

	productID string
}

func newSmsbaoSMSBuilder() *SmsbaoSMSBuilder {
	b := &SmsbaoSMSBuilder{}
	b.BaseBuilder = &BaseBuilder[*SmsbaoSMSBuilder]{subProvider: SubProviderSmsbao, self: b}
	return b
}

// ProductID sets the productId for Smsbao SMS.
// 当客户使用专用通道产品时，需要指定产品ID，产品ID可在短信宝后台或联系客服获得,不填则默认使用通用短信产品
//   - https://www.smsbao.com/openapi/213.html
func (b *SmsbaoSMSBuilder) ProductID(pid string) *SmsbaoSMSBuilder {
	b.productID = pid
	return b
}

func (b *SmsbaoSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	extra := map[string]interface{}{}
	if b.productID != "" {
		extra[smsbaoProductIDKey] = b.productID
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}
