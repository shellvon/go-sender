package telegram

import "github.com/shellvon/go-sender/core"

const (
	minPollOptions = 2
	maxPollOptions = 10
)

// PollMessage represents a poll message for Telegram
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
type PollMessage struct {
	BaseMessage

	// Poll question, 1-300 characters
	Question string `json:"question"`

	// Mode for parsing entities in the question. See formatting options for more details. Currently, only custom emoji entities are allowed
	// See  https://core.telegram.org/bots/api#formatting-options
	QuestionParseMode string `json:"question_parse_mode,omitempty"`

	// A JSON-serialized list of special entities that appear in the poll question. It can be specified instead of question_parse_mode
	QuestionEntities []MessageEntity `json:"question_entities,omitempty"`

	// A JSON-serialized list of 2-10 answer options
	Options []InputPollOption `json:"options"`

	// True, if the poll needs to be anonymous, defaults to True
	IsAnonymous bool `json:"is_anonymous,omitempty"`

	// Poll type, "quiz" or "regular", defaults to "regular"
	Type PollType `json:"type,omitempty"`

	// True, if the poll allows multiple answers, ignored for polls in quiz mode, defaults to False
	AllowsMultipleAnswers bool `json:"allows_multiple_answers,omitempty"`

	// 0-based identifier of the correct answer option, required for polls in quiz mode
	CorrectOptionID int `json:"correct_option_id,omitempty"`

	// Text that is shown when a user chooses an incorrect answer or taps on the lamp icon in a quiz-style poll, 0-200 characters with at most 2 line feeds after entities parsing
	Explanation string `json:"explanation,omitempty"`

	// Mode for parsing entities in the explanation. See formatting options for more details.
	// See  https://core.telegram.org/bots/api#formatting-options
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

// NewPollMessage creates a new PollMessage instance.
// Based on SendPollParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendpoll
//   - Only chat_id and question/options are required.
func NewPollMessage(chatID string, question string, options []InputPollOption) *PollMessage {
	return &PollMessage{
		BaseMessage: BaseMessage{
			MsgType: TypePoll,
			ChatID:  chatID,
		},
		Question: question,
		Options:  options,
	}
}

func (m *PollMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Question == "" {
		return core.NewParamError("question cannot be empty")
	}
	if len(m.Options) < minPollOptions || len(m.Options) > maxPollOptions {
		return core.NewParamError("poll must have between 2 and 10 options")
	}
	return nil
}

// Option helpers removed – use PollBuilder instead.
