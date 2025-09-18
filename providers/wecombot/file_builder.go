package wecombot

type FileBuilder struct {
	mediaID   string
	localPath string
}

// File 返回一个新的 FileBuilder 实例。
// 返回值：*FileBuilder - 新创建的 FileBuilder 实例，用于构建 FileMessage。
func File() *FileBuilder { return &FileBuilder{} }

// MediaID 设置文件消息的 media_id。
// 参数：id string - 要设置的 media_id。
// 返回值：*FileBuilder - 返回 FileBuilder 实例以支持链式调用。
func (b *FileBuilder) MediaID(id string) *FileBuilder { b.mediaID = id; return b }

// LocalPath 设置文件消息的本地文件路径。
// 参数：path string - 要设置的本地文件路径。
// 返回值：*FileBuilder - 返回 FileBuilder 实例以支持链式调用。
func (b *FileBuilder) LocalPath(path string) *FileBuilder { b.localPath = path; return b }

// Build 构建并返回一个 FileMessage 实例。
// 返回值：*FileMessage - 基于 FileBuilder 配置创建的文件消息实例。
func (b *FileBuilder) Build() *FileMessage {
	return &FileMessage{
		BaseMessage: newBaseMessage(TypeFile),
		File:        FilePayload{MediaID: b.mediaID},
		localPath:   b.localPath,
	}
}
