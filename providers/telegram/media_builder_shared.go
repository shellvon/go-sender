package telegram

// mediaBuilder is an abstract type that encapsulates common logic for media message builders.
// It provides a uniform interface for setting caption, parse mode, entities, and show caption above media.
// Concrete builders embed *mediaBuilder[*XxxBuilder] to inherit baseBuilder methods and caption-related setters.
// This avoids code duplication in each builder and ensures type safety.

type mediaBuilder[T any] struct {
	*baseBuilder[T]

	caption        string
	parseMode      ParseMode
	entities       []MessageEntity
	showCaptionTop bool
}

// newMediaBuilder creates a new mediaBuilder instance with the self type parameter.
// It correctly sets the self field of the internal baseBuilder to ensure that the chain call returns the specific builder type.
func newMediaBuilder[T any](self T) *mediaBuilder[T] {
	return &mediaBuilder[T]{
		baseBuilder: &baseBuilder[T]{self: self},
	}
}

// Caption sets the caption field.
func (b *mediaBuilder[T]) Caption(c string) T {
	b.caption = c
	return b.baseBuilder.self
}

// WithParseMode sets the parse mode for the text message.
func (b *mediaBuilder[T]) WithParseMode(mode ParseMode) T {
	b.parseMode = mode
	return b.self
}

// WithMarkdown sets the parse mode to Markdown.
func (b *mediaBuilder[T]) WithMarkdown() T {
	return b.WithParseMode(ParseModeMarkdown)
}

// WithMarkdownV2 sets the parse mode to MarkdownV2.
func (b *mediaBuilder[T]) WithMarkdownV2() T {
	return b.WithParseMode(ParseModeMarkdownV2)
}

// WithHTML sets the parse mode to HTML.
func (b *mediaBuilder[T]) WithHTML() T {
	return b.WithParseMode(ParseModeHTML)
}

// Entities sets the caption_entities field.
func (b *mediaBuilder[T]) Entities(ents []MessageEntity) T {
	b.entities = ents
	return b.baseBuilder.self
}

// ShowCaptionAboveMedia sets the show_caption_above_media field.
func (b *mediaBuilder[T]) ShowCaptionAboveMedia(show bool) T {
	b.showCaptionTop = show
	return b.baseBuilder.self
}
