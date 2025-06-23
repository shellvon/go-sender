package telegram

import "github.com/shellvon/go-sender/core"

// PollMessage represents a poll message for Telegram
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
type PollMessage struct {
	BaseMessage

	// Poll question, 1-300 characters
	Question string `json:"question"`

	// Mode for parsing entities in the question. See formatting options for more details. Currently, only custom emoji entities are allowed
	QuestionParseMode string `json:"question_parse_mode,omitempty"`

	// A JSON-serialized list of special entities that appear in the poll question. It can be specified instead of question_parse_mode
	QuestionEntities []MessageEntity `json:"question_entities,omitempty"`

	// A JSON-serialized list of 2-10 answer options
	Options []InputPollOption `json:"options"`

	// True, if the poll needs to be anonymous, defaults to True
	IsAnonymous bool `json:"is_anonymous,omitempty"`

	// Poll type, "quiz" or "regular", defaults to "regular"
	Type string `json:"type,omitempty"`

	// True, if the poll allows multiple answers, ignored for polls in quiz mode, defaults to False
	AllowsMultipleAnswers bool `json:"allows_multiple_answers,omitempty"`

	// 0-based identifier of the correct answer option, required for polls in quiz mode
	CorrectOptionID int `json:"correct_option_id,omitempty"`

	// Text that is shown when a user chooses an incorrect answer or taps on the lamp icon in a quiz-style poll, 0-200 characters with at most 2 line feeds after entities parsing
	Explanation string `json:"explanation,omitempty"`

	// Mode for parsing entities in the explanation. See formatting options for more details.
	ExplanationParseMode string `json:"explanation_parse_mode,omitempty"`

	// A JSON-serialized list of special entities that appear in the poll explanation. It can be specified instead of explanation_parse_mode
	ExplanationEntities []MessageEntity `json:"explanation_entities,omitempty"`

	// Amount of time in seconds the poll will be active after creation, 5-600. Can't be used together with close_date.
	OpenPeriod int `json:"open_period,omitempty"`

	// Point in time (Unix timestamp) when the poll will be automatically closed. Must be at least 5 and no more than 600 seconds in the future. Can't be used together with open_period.
	CloseDate int64 `json:"close_date,omitempty"`

	// Pass True if the poll needs to be immediately closed. This can be useful for poll preview.
	IsClosed bool `json:"is_closed,omitempty"`
}

func (m *PollMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *PollMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *PollMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Question == "" {
		return core.NewParamError("question cannot be empty")
	}
	if len(m.Options) < 2 || len(m.Options) > 10 {
		return core.NewParamError("poll must have between 2 and 10 options")
	}
	return nil
}

type PollMessageOption func(*PollMessage)

// WithPollQuestionParseMode sets the parse mode for the poll question
// Currently, only custom emoji entities are allowed
func WithPollQuestionParseMode(mode string) PollMessageOption {
	return func(m *PollMessage) { m.QuestionParseMode = mode }
}

// WithPollQuestionEntities sets the entities for the poll question
// A JSON-serialized list of special entities that appear in the poll question
func WithPollQuestionEntities(entities []MessageEntity) PollMessageOption {
	return func(m *PollMessage) { m.QuestionEntities = entities }
}

// WithPollIsAnonymous sets whether the poll should be anonymous
// Defaults to True if not specified
func WithPollIsAnonymous(anonymous bool) PollMessageOption {
	return func(m *PollMessage) { m.IsAnonymous = anonymous }
}

// WithPollType sets the type of the poll
// Options: "quiz" or "regular", defaults to "regular"
func WithPollType(pollType string) PollMessageOption {
	return func(m *PollMessage) { m.Type = pollType }
}

// WithPollAllowsMultipleAnswers sets whether the poll allows multiple answers
// Ignored for polls in quiz mode, defaults to False
func WithPollAllowsMultipleAnswers(allow bool) PollMessageOption {
	return func(m *PollMessage) { m.AllowsMultipleAnswers = allow }
}

// WithPollCorrectOptionID sets the 0-based identifier of the correct answer option
// Required for polls in quiz mode
func WithPollCorrectOptionID(correctOptionID int) PollMessageOption {
	return func(m *PollMessage) { m.CorrectOptionID = correctOptionID }
}

// WithPollExplanation sets the explanation text for quiz polls
// Text shown when a user chooses an incorrect answer or taps on the lamp icon
// 0-200 characters with at most 2 line feeds after entities parsing
func WithPollExplanation(explanation string) PollMessageOption {
	return func(m *PollMessage) { m.Explanation = explanation }
}

// WithPollExplanationParseMode sets the parse mode for the explanation text
// Supported modes: "HTML", "Markdown", "MarkdownV2"
func WithPollExplanationParseMode(mode string) PollMessageOption {
	return func(m *PollMessage) { m.ExplanationParseMode = mode }
}

// WithPollExplanationEntities sets the entities for the poll explanation
// A JSON-serialized list of special entities that appear in the poll explanation
func WithPollExplanationEntities(entities []MessageEntity) PollMessageOption {
	return func(m *PollMessage) { m.ExplanationEntities = entities }
}

// WithPollOpenPeriod sets the amount of time in seconds the poll will be active
// Range: 5-600 seconds. Can't be used together with close_date
func WithPollOpenPeriod(period int) PollMessageOption {
	return func(m *PollMessage) { m.OpenPeriod = period }
}

// WithPollCloseDate sets the point in time when the poll will be automatically closed
// Unix timestamp, must be at least 5 and no more than 600 seconds in the future
// Can't be used together with open_period
func WithPollCloseDate(closeDate int64) PollMessageOption {
	return func(m *PollMessage) { m.CloseDate = closeDate }
}

// WithPollIsClosed sets whether the poll should be immediately closed
// Useful for poll preview
func WithPollIsClosed(isClosed bool) PollMessageOption {
	return func(m *PollMessage) { m.IsClosed = isClosed }
}

func NewPollMessage(chatID string, question string, options []InputPollOption, opts ...interface{}) *PollMessage {
	msg := &PollMessage{
		BaseMessage: BaseMessage{
			MsgType: TypePoll,
			ChatID:  chatID,
		},
		Question: question,
		Options:  options,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case PollMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}
