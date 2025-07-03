package wecombot

// ImageBuilder provides a fluent API to construct WeCom image messages.
// Example:
//
//	msg := wecombot.Image().
//	         Base64(imgB64).
//	         MD5(imgMD5).
//	         Build()
type ImageBuilder struct {
	base64 string
	md5    string
}

// Image creates a new ImageBuilder.
// Based on SendImageParams from WeCom Bot API
// https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
func Image() *ImageBuilder {
	return &ImageBuilder{}
}

// Base64 sets the required base64-encoded image content.
func (b *ImageBuilder) Base64(data string) *ImageBuilder {
	b.base64 = data
	return b
}

// MD5 sets the required md5 hash of the (raw) image.
// Based on SendImageParams from WeCom Bot API.
func (b *ImageBuilder) MD5(hash string) *ImageBuilder {
	b.md5 = hash
	return b
}

// Build assembles an *ImageMessage.
func (b *ImageBuilder) Build() *ImageMessage {
	return &ImageMessage{
		BaseMessage: BaseMessage{MsgType: TypeImage},
		Image: ImageContent{
			Base64: b.base64,
			MD5:    b.md5,
		},
	}
}
