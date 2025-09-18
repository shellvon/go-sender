package core

import (
	"github.com/google/uuid"
)

// Message is the core interface for all messages.
//
//   - ProviderType returns the provider type.
//   - MsgID returns the message unique identifier.
type Message interface {
	ProviderType() ProviderType
	MsgID() string
}

// SubProviderAware is an optional interface for messages that support sub-providers.
//   - for example, SMS provider supports sub-providers like aliyun, tencent, etc.
//   - EmailAPI provider supports sub-providers like emailjs, resend, etc.
type SubProviderAware interface {
	GetSubProvider() string
}

// BaseMessage is the base message structure.
type BaseMessage struct {
	msgID        string       `json:"-"` // 消息ID，不序列化
	providerType ProviderType `json:"-"` // 提供者类型，不序列化
}

// Compile-time assertion: BaseMessage implements Message interface.
var _ Message = (*BaseMessage)(nil)

// NewBaseMessage Sets the provider type and creates a new BaseMessage.
func NewBaseMessage(providerType ProviderType) *BaseMessage {
	return &BaseMessage{
		providerType: providerType,
	}
}

// ProviderType Implements the Message interface.
// Returns the provider type.
func (m *BaseMessage) ProviderType() ProviderType {
	return m.providerType
}

// MsgID Implements the Message interface.
// Returns the message ID.
func (m *BaseMessage) MsgID() string {
	if m.msgID == "" {
		m.msgID = uuid.NewString()
	}
	return m.msgID
}

// SetMsgID Sets the message ID.
// Does not generate a new UUID if the message ID is already set.
func (m *BaseMessage) SetMsgID(id string) {
	m.msgID = id
}
