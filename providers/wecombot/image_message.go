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

// ImageContent represents the image content for a WeCom message.
type ImageContent struct {
	// Base64 encoded content of the image.
	Base64 string `json:"base64"`
	// MD5 hash of the image content before Base64 encoding.
	MD5 string `json:"md5"`
}

// ImageMessage represents an image message for WeCom.
// For more details, refer to the WeCom API documentation:
// https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
// Note: The original image size (before Base64 encoding) must not exceed 2MB. JPG and PNG formats are supported.
type ImageMessage struct {
	BaseMessage

	Image ImageContent `json:"image"`
}

// NewImageMessage creates a new ImageMessage.
// Based on SendImageParams from WeCom Bot API
// https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
//   - Only base64 and md5 are required.
func NewImageMessage(base64, md5 string) *ImageMessage {
	return Image().Base64(base64).MD5(md5).Build()
}

// Validate validates the ImageMessage to ensure it meets WeCom API requirements.
func (m *ImageMessage) Validate() error {
	if m.Image.Base64 == "" {
		return core.NewParamError("image base64 cannot be empty")
	}
	if m.Image.MD5 == "" {
		return core.NewParamError("image md5 cannot be empty")
	}

	// Estimate the raw size of the image from its Base64 encoded string.
	estimatedRawSize := (len(m.Image.Base64) * base64ToRawRatio) / base64Divisor
	if estimatedRawSize > maxImageSizeBytes { // Image size limit is 2MB.
		return core.NewParamError(
			fmt.Sprintf(
				"image size exceeds %dMB: estimated size %d bytes",
				maxImageSizeBytes/bytesPerMB,
				estimatedRawSize,
			),
		)
	}
	return nil
}
