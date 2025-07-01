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

type YuntongxunSMSBuilder struct {
	*BaseBuilder

	playTimes   string
	mediaName   string
	displayNum  string
	userData    string
	maxCallTime string
}

func newYuntongxunSMSBuilder() *YuntongxunSMSBuilder {
	return &YuntongxunSMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderYuntongxun},
	}
}

// PlayTimes sets the playTimes field for YunTongXun voice SMS.
// 语音短信播放次数。
//   - https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
func (b *YuntongxunSMSBuilder) PlayTimes(times string) *YuntongxunSMSBuilder {
	b.playTimes = times
	return b
}

// MediaName sets the mediaName field for YunTongXun voice SMS.
// 语音文件名称。
//   - https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
func (b *YuntongxunSMSBuilder) MediaName(name string) *YuntongxunSMSBuilder {
	b.mediaName = name
	return b
}

// DisplayNum sets the displayNum field for YunTongXun voice SMS.
// 显示号码。
//   - https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
func (b *YuntongxunSMSBuilder) DisplayNum(num string) *YuntongxunSMSBuilder {
	b.displayNum = num
	return b
}

// UserData sets the userData field for YunTongXun voice SMS.
// 用户数据。
//   - https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
func (b *YuntongxunSMSBuilder) UserData(data string) *YuntongxunSMSBuilder {
	b.userData = data
	return b
}

// MaxCallTime sets the maxCallTime field for YunTongXun voice SMS.
// 最大通话时长。
//   - https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
func (b *YuntongxunSMSBuilder) MaxCallTime(time string) *YuntongxunSMSBuilder {
	b.maxCallTime = time
	return b
}

func (b *YuntongxunSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	extra := map[string]interface{}{}
	if b.playTimes != "" {
		extra[yuntongxunPlayTimesKey] = b.playTimes
	}
	if b.mediaName != "" {
		extra[yuntongxunMediaNameKey] = b.mediaName
	}
	if b.displayNum != "" {
		extra[yuntongxunDisplayNumKey] = b.displayNum
	}
	if b.userData != "" {
		extra[yuntongxunUserDataKey] = b.userData
	}
	if b.maxCallTime != "" {
		extra[yuntongxunMaxCallTimeKey] = b.maxCallTime
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}
