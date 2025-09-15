package wecomapp

// BaseMediaBuilder 为所有媒体消息构建器提供通用功能.
type BaseMediaBuilder struct {
	mediaID   string
	localPath string
	toUser    string
	toParty   string
	toTag     string
	agentID   string
	safe      int
}

// MediaID 设置从媒体上传API获得的媒体文件ID.
func (b *BaseMediaBuilder) MediaID(mediaID string) *BaseMediaBuilder {
	b.mediaID = mediaID
	return b
}

// LocalPath 设置本地文件路径用于自动上传.
func (b *BaseMediaBuilder) LocalPath(localPath string) *BaseMediaBuilder {
	b.localPath = localPath
	return b
}

// ToUser 设置发送给的用户ID，用"|"分隔。使用"@all"发送给所有用户.
func (b *BaseMediaBuilder) ToUser(toUser string) *BaseMediaBuilder {
	b.toUser = toUser
	return b
}

// ToParty sets the department IDs to send to, separated by "|".
func (b *BaseMediaBuilder) ToParty(toParty string) *BaseMediaBuilder {
	b.toParty = toParty
	return b
}

// ToTag sets the tag IDs to send to, separated by "|".
func (b *BaseMediaBuilder) ToTag(toTag string) *BaseMediaBuilder {
	b.toTag = toTag
	return b
}

// AgentID sets the application ID (required).
func (b *BaseMediaBuilder) AgentID(agentID string) *BaseMediaBuilder {
	b.agentID = agentID
	return b
}

// Safe sets whether to enable safe mode (0: no, 1: yes).
func (b *BaseMediaBuilder) Safe(safe int) *BaseMediaBuilder {
	b.safe = safe
	return b
}

// buildCommonFields creates CommonFields from builder state.
func (b *BaseMediaBuilder) buildCommonFields() CommonFields {
	return CommonFields{
		ToUser:  b.toUser,
		ToParty: b.toParty,
		ToTag:   b.toTag,
		AgentID: b.agentID,
		Safe:    b.safe,
	}
}

// ImageBuilder provides a fluent API to construct WeChat Work Application image messages.
type ImageBuilder struct {
	BaseMediaBuilder
}

// Image creates a new ImageBuilder instance.
func Image() *ImageBuilder {
	return &ImageBuilder{}
}

// MediaID sets the media file ID (fluent interface).
func (b *ImageBuilder) MediaID(mediaID string) *ImageBuilder {
	b.BaseMediaBuilder.MediaID(mediaID)
	return b
}

// LocalPath sets the local file path (fluent interface).
func (b *ImageBuilder) LocalPath(localPath string) *ImageBuilder {
	b.BaseMediaBuilder.LocalPath(localPath)
	return b
}

// ToUser sets the user IDs (fluent interface).
func (b *ImageBuilder) ToUser(toUser string) *ImageBuilder {
	b.BaseMediaBuilder.ToUser(toUser)
	return b
}

// ToParty sets the department IDs (fluent interface).
func (b *ImageBuilder) ToParty(toParty string) *ImageBuilder {
	b.BaseMediaBuilder.ToParty(toParty)
	return b
}

// ToTag sets the tag IDs (fluent interface).
func (b *ImageBuilder) ToTag(toTag string) *ImageBuilder {
	b.BaseMediaBuilder.ToTag(toTag)
	return b
}

// AgentID sets the application ID (fluent interface).
func (b *ImageBuilder) AgentID(agentID string) *ImageBuilder {
	b.BaseMediaBuilder.AgentID(agentID)
	return b
}

// Safe sets safe mode (fluent interface).
func (b *ImageBuilder) Safe(safe int) *ImageBuilder {
	b.BaseMediaBuilder.Safe(safe)
	return b
}

// Build assembles a ready-to-send *ImageMessage.
func (b *ImageBuilder) Build() *ImageMessage {
	msg := &ImageMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{
				CommonFields: b.buildCommonFields(),
				MsgType:      TypeImage,
			},
		},
	}

	if b.mediaID != "" {
		msg.Image.MediaID = b.mediaID
	}

	if b.localPath != "" {
		msg.LocalPath = b.localPath
	}

	return msg
}

// VoiceBuilder provides a fluent API to construct WeChat Work Application voice messages.
type VoiceBuilder struct {
	BaseMediaBuilder
}

// Voice creates a new VoiceBuilder instance.
func Voice() *VoiceBuilder {
	return &VoiceBuilder{}
}

// MediaID sets the media file ID (fluent interface).
func (b *VoiceBuilder) MediaID(mediaID string) *VoiceBuilder {
	b.BaseMediaBuilder.MediaID(mediaID)
	return b
}

// LocalPath sets the local file path (fluent interface).
func (b *VoiceBuilder) LocalPath(localPath string) *VoiceBuilder {
	b.BaseMediaBuilder.LocalPath(localPath)
	return b
}

// ToUser sets the user IDs (fluent interface).
func (b *VoiceBuilder) ToUser(toUser string) *VoiceBuilder {
	b.BaseMediaBuilder.ToUser(toUser)
	return b
}

// ToParty sets the department IDs (fluent interface).
func (b *VoiceBuilder) ToParty(toParty string) *VoiceBuilder {
	b.BaseMediaBuilder.ToParty(toParty)
	return b
}

// ToTag sets the tag IDs (fluent interface).
func (b *VoiceBuilder) ToTag(toTag string) *VoiceBuilder {
	b.BaseMediaBuilder.ToTag(toTag)
	return b
}

// AgentID sets the application ID (fluent interface).
func (b *VoiceBuilder) AgentID(agentID string) *VoiceBuilder {
	b.BaseMediaBuilder.AgentID(agentID)
	return b
}

// Safe sets safe mode (fluent interface).
func (b *VoiceBuilder) Safe(safe int) *VoiceBuilder {
	b.BaseMediaBuilder.Safe(safe)
	return b
}

// Build assembles a ready-to-send *VoiceMessage.
func (b *VoiceBuilder) Build() *VoiceMessage {
	msg := &VoiceMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{
				CommonFields: b.buildCommonFields(),
				MsgType:      TypeVoice,
			},
		},
	}

	if b.mediaID != "" {
		msg.Voice.MediaID = b.mediaID
	}

	if b.localPath != "" {
		msg.LocalPath = b.localPath
	}

	return msg
}

// VideoBuilder provides a fluent API to construct WeChat Work Application video messages.
type VideoBuilder struct {
	BaseMediaBuilder

	title       string
	description string
}

// Video creates a new VideoBuilder instance.
func Video() *VideoBuilder {
	return &VideoBuilder{}
}

// MediaID sets the media file ID (fluent interface).
func (b *VideoBuilder) MediaID(mediaID string) *VideoBuilder {
	b.BaseMediaBuilder.MediaID(mediaID)
	return b
}

// LocalPath sets the local file path (fluent interface).
func (b *VideoBuilder) LocalPath(localPath string) *VideoBuilder {
	b.BaseMediaBuilder.LocalPath(localPath)
	return b
}

// Title sets the video title (optional).
func (b *VideoBuilder) Title(title string) *VideoBuilder {
	b.title = title
	return b
}

// Description sets the video description (optional).
func (b *VideoBuilder) Description(description string) *VideoBuilder {
	b.description = description
	return b
}

// ToUser sets the user IDs (fluent interface).
func (b *VideoBuilder) ToUser(toUser string) *VideoBuilder {
	b.BaseMediaBuilder.ToUser(toUser)
	return b
}

// ToParty sets the department IDs (fluent interface).
func (b *VideoBuilder) ToParty(toParty string) *VideoBuilder {
	b.BaseMediaBuilder.ToParty(toParty)
	return b
}

// ToTag sets the tag IDs (fluent interface).
func (b *VideoBuilder) ToTag(toTag string) *VideoBuilder {
	b.BaseMediaBuilder.ToTag(toTag)
	return b
}

// AgentID sets the application ID (fluent interface).
func (b *VideoBuilder) AgentID(agentID string) *VideoBuilder {
	b.BaseMediaBuilder.AgentID(agentID)
	return b
}

// Safe sets safe mode (fluent interface).
func (b *VideoBuilder) Safe(safe int) *VideoBuilder {
	b.BaseMediaBuilder.Safe(safe)
	return b
}

// Build assembles a ready-to-send *VideoMessage.
func (b *VideoBuilder) Build() *VideoMessage {
	msg := &VideoMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{
				CommonFields: b.buildCommonFields(),
				MsgType:      TypeVideo,
			},
		},
	}

	if b.mediaID != "" {
		msg.Video.MediaID = b.mediaID
	}

	if b.title != "" {
		msg.Video.Title = b.title
	}

	if b.description != "" {
		msg.Video.Description = b.description
	}

	if b.localPath != "" {
		msg.LocalPath = b.localPath
	}

	return msg
}

// FileBuilder provides a fluent API to construct WeChat Work Application file messages.
type FileBuilder struct {
	BaseMediaBuilder
}

// File creates a new FileBuilder instance.
func File() *FileBuilder {
	return &FileBuilder{}
}

// MediaID sets the media file ID (fluent interface).
func (b *FileBuilder) MediaID(mediaID string) *FileBuilder {
	b.BaseMediaBuilder.MediaID(mediaID)
	return b
}

// LocalPath sets the local file path (fluent interface).
func (b *FileBuilder) LocalPath(localPath string) *FileBuilder {
	b.BaseMediaBuilder.LocalPath(localPath)
	return b
}

// ToUser sets the user IDs (fluent interface).
func (b *FileBuilder) ToUser(toUser string) *FileBuilder {
	b.BaseMediaBuilder.ToUser(toUser)
	return b
}

// ToParty sets the department IDs (fluent interface).
func (b *FileBuilder) ToParty(toParty string) *FileBuilder {
	b.BaseMediaBuilder.ToParty(toParty)
	return b
}

// ToTag sets the tag IDs (fluent interface).
func (b *FileBuilder) ToTag(toTag string) *FileBuilder {
	b.BaseMediaBuilder.ToTag(toTag)
	return b
}

// AgentID sets the application ID (fluent interface).
func (b *FileBuilder) AgentID(agentID string) *FileBuilder {
	b.BaseMediaBuilder.AgentID(agentID)
	return b
}

// Safe sets safe mode (fluent interface).
func (b *FileBuilder) Safe(safe int) *FileBuilder {
	b.BaseMediaBuilder.Safe(safe)
	return b
}

// Build assembles a ready-to-send *FileMessage.
func (b *FileBuilder) Build() *FileMessage {
	msg := &FileMessage{
		BaseMediaMessage: BaseMediaMessage{
			BaseMessage: BaseMessage{
				CommonFields: b.buildCommonFields(),
				MsgType:      TypeFile,
			},
		},
	}

	if b.mediaID != "" {
		msg.File.MediaID = b.mediaID
	}

	if b.localPath != "" {
		msg.LocalPath = b.localPath
	}

	return msg
}
