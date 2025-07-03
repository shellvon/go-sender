package telegram

// mediaBuilder 抽象了所有带媒体消息的 Caption 相关公共逻辑。
// 具体的 XxxBuilder 通过嵌入 *mediaBuilder[*XxxBuilder] 获得统一的
// Chat / Silent / Protect 等 baseBuilder 方法以及 Caption 等 setter，
// 避免在每个 builder 中重复编写同样的字段和链式方法。

type mediaBuilder[T any] struct {
	*baseBuilder[T]

	caption        string
	parseMode      string
	entities       []MessageEntity
	showCaptionTop bool
}

// newMediaBuilder 创建一个带 self 泛型的 mediaBuilder 实例。
// 该函数会正确设置内部 baseBuilder 的 self 字段，确保链式调用返回具体 builder 类型。
func newMediaBuilder[T any](self T) *mediaBuilder[T] {
	return &mediaBuilder[T]{
		baseBuilder: &baseBuilder[T]{self: self},
	}
}

// Caption 设置 Caption 文案。
func (b *mediaBuilder[T]) Caption(c string) T {
	b.caption = c
	return b.baseBuilder.self
}

// ParseMode 设置解析模式 (HTML / Markdown / MarkdownV2)。
func (b *mediaBuilder[T]) ParseMode(mode string) T {
	b.parseMode = mode
	return b.baseBuilder.self
}

// Entities 设置 caption_entities。
func (b *mediaBuilder[T]) Entities(ents []MessageEntity) T {
	b.entities = ents
	return b.baseBuilder.self
}

// ShowCaptionAboveMedia 决定 caption 是否显示在媒体上方。
func (b *mediaBuilder[T]) ShowCaptionAboveMedia(show bool) T {
	b.showCaptionTop = show
	return b.baseBuilder.self
}
