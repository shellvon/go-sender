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
	*BaseBuilder[*YuntongxunSMSBuilder]
}

func newYuntongxunSMSBuilder() *YuntongxunSMSBuilder {
	b := &YuntongxunSMSBuilder{}
	b.BaseBuilder = &BaseBuilder[*YuntongxunSMSBuilder]{subProvider: SubProviderYuntongxun, self: b}
	return b
}

// MediaNameType sets the mediaNameType field for YunTongXun voice SMS.
// 语音文件名的类型，默认值为0，表示用户语音文件；　值为1表示平台通用文件。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) MediaNameType(mediaNameType string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunMediaNameTypeKey, mediaNameType)
}

// MediaTxt sets the mediaTxt field for YunTongXun voice SMS.
// 文本内容，文本中汉字要求utf8编码，默认值为空。当mediaName为空才有效
// 注：此方法本质上就是Conetent的别名，只是为了与官方文档保持一致，方便用户使用。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) MediaTxt(txt string) *YuntongxunSMSBuilder {
	return b.Content(txt)
}

// PlayTimes sets the playTimes field for YunTongXun voice SMS.
// 循环播放次数，1－3次，默认播放1次。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) PlayTimes(times string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunPlayTimesKey, times)
}

// MediaName sets the mediaName field for YunTongXun voice SMS.
// 语音文件名称，格式 wav，播放多个文件用英文分号隔开。与mediaTxt不能同时为空。当不为空时mediaTxt属性失效。测试用默认语音：ccp_marketingcall.wav
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) MediaName(name string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunMediaNameKey, name)
}

// DisplayNum sets the displayNum field for YunTongXun voice SMS.
// 来电显示的号码，根据平台侧显号规则控制(有显号需求请联系云通讯商务，并且说明显号的方式)，不在平台规则内或空则显示云通讯平台默认号码。默认值空。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) DisplayNum(num string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunDisplayNumKey, num)
}

// UserData sets the userData field for YunTongXun voice SMS.
// 可选 用户数据，透传字段，可填入任意字符串，如：用户id，用户名等。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) UserData(data string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunUserDataKey, data)
}

// TxtSpeed sets the txtSpeed field for YunTongXun voice SMS.
// 文本转语音的语速，默认值为空。文本转语音后的发音速度，取值范围：-50至50，当mediaTxt有效才生效,默认值为0。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) TxtSpeed(speed string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunTxtSpeedKey, speed)
}

// TxtPitch sets the txtPitch field for YunTongXun voice SMS.
// 文本转语音后的音调，取值范围：-500至500，当mediaTxt有效才生效，默认值为0。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) TxtPitch(pitch string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunTxtPitchKey, pitch)
}

// TxtVolume sets the txtVolume field for YunTongXun voice SMS.
// 文本转语音后的音量大小，取值范围：-20至20，当mediaTxt有效才生效，默认值为0。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) TxtVolume(volume string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunTxtVolumeKey, volume)
}

// TxtBgsound sets the txtBgsound field for YunTongXun voice SMS.
// 文本转语音后的背景音编号，目前云通讯平台支持6种背景音，1到6的六种背景音编码，0为不需要背景音。暂时不支持第三方自定义背景音。当mediaTxt有效才生效。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) TxtBgsound(bgsound string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunTxtBgsoundKey, bgsound)
}

// PlayMode sets the playMode field for YunTongXun voice SMS.
//
// 是否同时播放文本和语音文件 , 0、否 1、是，默认0。优先播放文本。
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (b *YuntongxunSMSBuilder) PlayMode(mode string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunPlayModeKey, mode)
}

// MaxCallTime sets the maxCallTime field for YunTongXun voice SMS.
// 该通通话最大通话时长，到时间自动挂机
//   - 语音验证码: https://doc.yuntongxun.com/pe/5a533de43b8496dd00dce07e
func (b *YuntongxunSMSBuilder) MaxCallTime(time string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunMaxCallTimeKey, time)
}

// WelcomePrompt sets the welcomePrompt field for YunTongXun voice SMS.
// wav格式的文件名，欢迎提示音，在播放验证码语音前播放此内容，配合verifyCode使用，默认值为空。
//   - 语音验证码: https://doc.yuntongxun.com/pe/5a533de43b8496dd00dce07e
func (b *YuntongxunSMSBuilder) WelcomePrompt(prompt string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunWelcomePromptKey, prompt)
}

// PlayVerifyCode sets the playVerifyCode field for YunTongXun voice SMS.
// 播放验证码语音，默认值为空。
//   - 语音验证码: https://doc.yuntongxun.com/pe/5a533de43b8496dd00dce07e
func (b *YuntongxunSMSBuilder) PlayVerifyCode(code string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunPlayVerifyCodeKey, code)
}

// Region 设置云通讯短信的区域, 默认为国内, cn，否则使用港澳台以及境外的接入点
//   - 国内: https://api.cloopen.com:8883
//   - 港澳台以及境外: https://hksms.cloopen.com:8883
//
// 文档地址:
//   - 国内短信: https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
//   - 国际短信: https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d
//
// 此值目前影响国际短信的请求地址.
func (b *YuntongxunSMSBuilder) Region(region string) *YuntongxunSMSBuilder {
	return b.meta(yuntongxunRegionKey, region)
}
