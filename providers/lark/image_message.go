package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// ImageMessage represents an image message for Lark/Feishu.
type ImageMessage struct {
	BaseMessage

	Content ImageContent `json:"content"`
}

// ImageContent represents the content of an image message.
type ImageContent struct {
	Image ImagePayload `json:"image"`
}

// ImagePayload represents the image structure.
type ImagePayload struct {
	ImageKey string `json:"image_key"`
}

// imageMsgBuilder provides a fluent API to construct Lark image messages (unexported to avoid conflict).
type imageMsgBuilder struct {
	imageKey string
}

// Image creates a new imageMsgBuilder instance (user-facing API).
func Image() *imageMsgBuilder { return &imageMsgBuilder{} }

// ImageKey sets the image key.
func (b *imageMsgBuilder) ImageKey(key string) *imageMsgBuilder {
	b.imageKey = key
	return b
}

// Build assembles a *ImageMessage.
func (b *imageMsgBuilder) Build() *ImageMessage {
	return &ImageMessage{
		BaseMessage: BaseMessage{MsgType: TypeImage},
		Content:     ImageContent{Image: ImagePayload{ImageKey: b.imageKey}},
	}
}

// NewImageMessage creates a new image message.
func NewImageMessage(imageKey string) *ImageMessage {
	return Image().ImageKey(imageKey).Build()
}

// GetMsgType returns the message type.
func (m *ImageMessage) GetMsgType() MessageType {
	return TypeImage
}

// ProviderType returns the provider type.
func (m *ImageMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeLark
}

// Validate validates the image message.
func (m *ImageMessage) Validate() error {
	if m.Content.Image.ImageKey == "" {
		return errors.New("image_key cannot be empty")
	}
	return nil
}
