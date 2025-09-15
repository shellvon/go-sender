package wecombot

// ImageBuilder 提供了一个流式 API 来构建企业微信图片消息。
// 示例：
//
//	msg := wecombot.Image().
//	         Base64(imgB64).
//	         MD5(imgMD5).
//	         Build()
type ImageBuilder struct {
	base64 string
	md5    string
}

// Image 创建一个新的 ImageBuilder 实例。
// 基于企业微信机器人 API 的 SendImageParams
// 参考：https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
// 返回值：*ImageBuilder - 新创建的 ImageBuilder 实例，用于构建 ImageMessage。
func Image() *ImageBuilder {
	return &ImageBuilder{}
}

// Base64 设置图片内容的 Base64 编码（必需）。
// 参数：data string - 图片的 Base64 编码字符串。
// 返回值：*ImageBuilder - 返回 ImageBuilder 实例以支持链式调用。
func (b *ImageBuilder) Base64(data string) *ImageBuilder {
	b.base64 = data
	return b
}

// MD5 设置原始图片的 MD5 哈希值（必需）。
// 基于企业微信机器人 API 的 SendImageParams。
// 参数：hash string - 图片的 MD5 哈希值。
// 返回值：*ImageBuilder - 返回 ImageBuilder 实例以支持链式调用。
func (b *ImageBuilder) MD5(hash string) *ImageBuilder {
	b.md5 = hash
	return b
}

// Build 构建并返回一个 ImageMessage 实例。
// 返回值：*ImageMessage - 基于 ImageBuilder 配置创建的图片消息实例。
func (b *ImageBuilder) Build() *ImageMessage {
	return &ImageMessage{
		BaseMessage: BaseMessage{MsgType: TypeImage},
		Image: ImageContent{
			Base64: b.base64,
			MD5:    b.md5,
		},
	}
}
