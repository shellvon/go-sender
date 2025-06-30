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

// Image represents the image content for a WeCom message.
type Image struct {
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

	Image Image `json:"image"`
}

// NewImageMessage creates a new ImageMessage with required fields and applies optional configurations.
func NewImageMessage(base64, md5 string, opts ...ImageMessageOption) *ImageMessage {
	// Initialize ImageMessage with required fields.
	msg := &ImageMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeImage,
		},
		Image: Image{
			Base64: base64,
			MD5:    md5,
		},
	}

	// Apply optional configurations, which can override the initial base64 or md5 if provided.
	for _, opt := range opts {
		opt(msg)
	}

	return msg
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

// ImageMessageOption defines a function type for configuring ImageMessage.
type ImageMessageOption func(*ImageMessage)

// WithBase64 sets the Base64 field for ImageMessage.
func WithBase64(base64 string) ImageMessageOption {
	return func(m *ImageMessage) {
		m.Image.Base64 = base64
	}
}

// WithMD5 sets the MD5 field for ImageMessage.
func WithMD5(md5 string) ImageMessageOption {
	return func(m *ImageMessage) {
		m.Image.MD5 = md5
	}
}
