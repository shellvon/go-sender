package sms

// @ProviderName: Yuntongxun / 容联云
// @Website: https://www.yuntongxun.com
// @APIDoc: https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
//
// 官方文档:
//   - 国内短信: https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
//   - 国际短信: https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d
//
// builder 仅支持 text（普通短信）类型。

// NewYuntongxunTextMessage creates a new Yuntongxun text SMS message.
func NewYuntongxunTextMessage(mobiles []string, content, sign string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderYuntongxun)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
		WithSignName(sign),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// WithYuntongxunPlayTimes sets the playTimes field for YunTongXun voice SMS.
func WithYuntongxunPlayTimes(playTimes string) MessageOption {
	return WithExtra(yuntongxunPlayTimesKey, playTimes)
}

// WithYuntongxunMediaName sets the mediaName field for YunTongXun voice SMS.
func WithYuntongxunMediaName(mediaName string) MessageOption {
	return WithExtra(yuntongxunMediaNameKey, mediaName)
}

// WithYuntongxunDisplayNum sets the displayNum field for YunTongXun voice SMS.
func WithYuntongxunDisplayNum(displayNum string) MessageOption {
	return WithExtra(yuntongxunDisplayNumKey, displayNum)
}

// WithYuntongxunUserData sets the userData field for YunTongXun voice SMS.
func WithYuntongxunUserData(userData string) MessageOption {
	return WithExtra(yuntongxunUserDataKey, userData)
}

// WithYuntongxunMaxCallTime sets the maxCallTime field for YunTongXun voice SMS.
func WithYuntongxunMaxCallTime(maxCallTime string) MessageOption {
	return WithExtra(yuntongxunMaxCallTimeKey, maxCallTime)
}
