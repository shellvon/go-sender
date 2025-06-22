package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// ImageMessage represents an image message for Lark/Feishu
type ImageMessage struct {
	BaseMessage
	Content ImageContent `json:"content"`
}

// ImageContent represents the content of an image message
type ImageContent struct {
	Image Image `json:"image"`
}

// Image represents the image structure
type Image struct {
	ImageKey string `json:"image_key"`
}

// NewImageMessage creates a new image message
func NewImageMessage(imageKey string) *ImageMessage {
	return &ImageMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeImage,
		},
		Content: ImageContent{
			Image: Image{
				ImageKey: imageKey,
			},
		},
	}
}

// GetMsgType returns the message type
func (m *ImageMessage) GetMsgType() MessageType {
	return TypeImage
}

// ProviderType returns the provider type
func (m *ImageMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the image message
func (m *ImageMessage) Validate() error {
	if m.Content.Image.ImageKey == "" {
		return errors.New("image_key cannot be empty")
	}
	return nil
}
