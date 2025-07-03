package telegram

// PollBuilder constructs Telegram poll messages with compile-time safety.
// Example:
//   msg := telegram.Poll().
//            Chat("123").
//            Question("Your favourite color?").
//            Options(
//                telegram.Option("Red"),
//                telegram.Option("Blue"),
//            ).
//            AllowsMultipleAnswers(true).
//            Build()

type PollBuilder struct {
	*baseBuilder[*PollBuilder]

	question string
	options  []InputPollOption

	// optional fields mirroring PollMessage
	questionParseMode    string
	questionEntities     []MessageEntity
	isAnonymous          bool
	pollType             string
	allowsMultiple       bool
	correctOptionID      int
	explanation          string
	explanationParseMode string
	explanationEntities  []MessageEntity
	openPeriod           int
	closeDate            int64
	isClosed             bool
}

// Poll returns a new PollBuilder.
func Poll() *PollBuilder {
	b := &PollBuilder{correctOptionID: -1}
	b.baseBuilder = &baseBuilder[*PollBuilder]{self: b}
	return b
}

// Question sets the poll question (1-300 chars). Required.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// Poll question, 1-300 characters.
func (b *PollBuilder) Question(q string) *PollBuilder {
	b.question = q
	return b
}

// Options sets the answer options (2-10). Required.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// A JSON-serialized list of 2-10 answer options.
func (b *PollBuilder) Options(opts ...InputPollOption) *PollBuilder {
	b.options = opts
	return b
}

// QuestionParseMode sets parse mode for question.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// Mode for parsing entities in the question. See formatting options for more details. Currently, only custom emoji entities are allowed
// See  https://core.telegram.org/bots/api#formatting-options
func (b *PollBuilder) QuestionParseMode(mode string) *PollBuilder {
	b.questionParseMode = mode
	return b
}

// QuestionEntities sets question entities.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// A JSON-serialized list of special entities that appear in the poll question. It can be specified instead of question_parse_mode.
func (b *PollBuilder) QuestionEntities(ents []MessageEntity) *PollBuilder {
	b.questionEntities = ents
	return b
}

// IsAnonymous toggles anonymity.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// True, if the poll needs to be anonymous, defaults to True.
func (b *PollBuilder) IsAnonymous(a bool) *PollBuilder {
	b.isAnonymous = a
	return b
}

// Type sets poll type ("quiz" or "regular"). default is "regular"
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
//
// all types:
//   - quiz
//   - regular (default)
func (b *PollBuilder) Type(t string) *PollBuilder {
	b.pollType = t
	return b
}

// AllowsMultipleAnswers toggles multi-select.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// True, if the poll allows multiple answers, ignored for polls in quiz mode, defaults to False.
func (b *PollBuilder) AllowsMultipleAnswers(allow bool) *PollBuilder {
	b.allowsMultiple = allow
	return b
}

// CorrectOptionID sets correct option index for quiz polls.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// 0-based identifier of the correct answer option, required for polls in quiz mode.
func (b *PollBuilder) CorrectOptionID(idx int) *PollBuilder {
	b.correctOptionID = idx
	return b
}

// Explanation sets explanation text for quiz polls.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// Text that is shown when a user chooses an incorrect answer or taps on the lamp icon in a quiz-style poll, 0-200 characters with at most 2 line feeds after entities parsing.
func (b *PollBuilder) Explanation(exp string) *PollBuilder {
	b.explanation = exp
	return b
}

// ExplanationParseMode sets explanation parse mode.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// Mode for parsing entities in the explanation. See formatting options for more details.
// See  https://core.telegram.org/bots/api#formatting-options
func (b *PollBuilder) ExplanationParseMode(mode string) *PollBuilder {
	b.explanationParseMode = mode
	return b
}

// ExplanationEntities sets explanation entities.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// A JSON-serialized list of special entities that appear in the poll explanation. It can be specified instead of explanation_parse_mode.
func (b *PollBuilder) ExplanationEntities(ents []MessageEntity) *PollBuilder {
	b.explanationEntities = ents
	return b
}

// OpenPeriod sets open period seconds (5-600).
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// Amount of time in seconds the poll will be active after creation, 5-600. Can't be used together with close_date.
func (b *PollBuilder) OpenPeriod(sec int) *PollBuilder {
	b.openPeriod = sec
	return b
}

// CloseDate sets close date unix timestamp.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// Point in time (Unix timestamp) when the poll will be automatically closed. Must be at least 5 and no more than 600 seconds in the future. Can't be used together with open_period.
func (b *PollBuilder) CloseDate(ts int64) *PollBuilder {
	b.closeDate = ts
	return b
}

// IsClosed marks poll closed immediately.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
// Pass True if the poll needs to be immediately closed. This can be useful for poll preview.
func (b *PollBuilder) IsClosed(closed bool) *PollBuilder {
	b.isClosed = closed
	return b
}

// Build assembles the *PollMessage.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
//   - Only chat_id and question/options are required.
func (b *PollBuilder) Build() *PollMessage {
	msg := &PollMessage{
		BaseMessage:           b.baseBuilder.toBaseMessage(TypePoll),
		Question:              b.question,
		Options:               b.options,
		QuestionParseMode:     b.questionParseMode,
		QuestionEntities:      b.questionEntities,
		IsAnonymous:           b.isAnonymous,
		Type:                  b.pollType,
		AllowsMultipleAnswers: b.allowsMultiple,
		CorrectOptionID:       b.correctOptionID,
		Explanation:           b.explanation,
		ExplanationParseMode:  b.explanationParseMode,
		ExplanationEntities:   b.explanationEntities,
		OpenPeriod:            b.openPeriod,
		CloseDate:             b.closeDate,
		IsClosed:              b.isClosed,
	}
	return msg
}

// Option is a helper to create InputPollOption more concisely.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
func Option(text string) InputPollOption {
	return InputPollOption{Text: text}
}
