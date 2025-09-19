package wecombot

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

const (
	base64ToRawRatio  = 3
	base64Divisor     = 4
	maxImageSizeBytes = 2 * 1024 * 1024 // 2MB
	bytesPerMB        = 1024 * 1024
)

// ImageContent 表示企业微信消息的图片内容。
type ImageContent struct {
	// Base64 编码的图片内容。
	Base64 string `json:"base64"`
	// Base64 编码前图片内容的 MD5 哈希值。
	MD5 string `json:"md5"`
}

// ImageMessage 表示企业微信的图片消息。
// 更多详情，请参考企业微信 API 文档：
// https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
// 注意：原始图片大小（Base64 编码前）不得超过 2MB。支持 JPG 和 PNG 格式。
type ImageMessage struct {
	BaseMessage

	Image ImageContent `json:"image"`
}

// NewImageMessage 创建一个新的 ImageMessage 实例。
// 基于企业微信机器人 API 的 SendImageParams
// 参考：https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
//
// 参数：
//   - base64 string - 图片的 Base64 编码内容（必需）
//   - md5 string - 图片的 MD5 哈希值（必需）
//
// 返回值：*ImageMessage - 新创建的图片消息实例.
func NewImageMessage(base64, md5 string) *ImageMessage {
	return Image().Base64(base64).MD5(md5).Build()
}

// Validate 验证 ImageMessage 是否满足企业微信 API 要求。
// 该方法检查 Base64 和 MD5 是否为空，并估算图片原始大小是否超过 2MB 限制。
// 返回值：error - 如果验证失败，返回具体的错误信息；否则返回 nil。
func (m *ImageMessage) Validate() error {
	if m.Image.Base64 == "" {
		return core.NewParamError("图片 Base64 不能为空")
	}
	if m.Image.MD5 == "" {
		return core.NewParamError("图片 MD5 不能为空")
	}

	// 根据 Base64 编码字符串估算图片的原始大小。
	estimatedRawSize := (len(m.Image.Base64) * base64ToRawRatio) / base64Divisor
	if estimatedRawSize > maxImageSizeBytes { // 图片大小限制为 2MB。
		return core.NewParamError(
			fmt.Sprintf(
				"图片大小超过 %dMB：估算大小 %d 字节",
				maxImageSizeBytes/bytesPerMB,
				estimatedRawSize,
			),
		)
	}
	return nil
}
