package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// Telegram is a message provider for Telegram.
// It supports sending text, image, audio, video, document, location, contact, poll, dice, venue, and other types of messages.
// It also supports sending messages to different channels and groups.
//
// Reference:
//   - Official Website: https://telegram.org/
//   - API Docs: https://core.telegram.org/bots/api

const (
	endpointSendMessage   = "sendMessage"
	endpointSendPhoto     = "sendPhoto"
	endpointSendAudio     = "sendAudio"
	endpointSendVoice     = "sendVoice"
	endpointSendDocument  = "sendDocument"
	endpointSendVideo     = "sendVideo"
	endpointSendAnimation = "sendAnimation"
	endpointSendVideoNote = "sendVideoNote"
	endpointSendLocation  = "sendLocation"
	endpointSendContact   = "sendContact"
	endpointSendPoll      = "sendPoll"
	endpointSendDice      = "sendDice"
	endpointSendVenue     = "sendVenue"
)

const telegramAPIURLTemplate = "https://api.telegram.org/bot%s/%s"

type telegramTransformer struct {
	*transformer.BaseHTTPTransformer[Message, *Account]
}

func (tt *telegramTransformer) transform(
	_ context.Context,
	msg Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	var endpoint string

	switch msg.(type) {
	case *TextMessage:
		endpoint = endpointSendMessage
	case *PhotoMessage:
		endpoint = endpointSendPhoto
	case *AudioMessage:
		endpoint = endpointSendAudio
	case *VoiceMessage:
		endpoint = endpointSendVoice
	case *DocumentMessage:
		endpoint = endpointSendDocument
	case *VideoMessage:
		endpoint = endpointSendVideo
	case *AnimationMessage:
		endpoint = endpointSendAnimation
	case *VideoNoteMessage:
		endpoint = endpointSendVideoNote
	case *LocationMessage:
		endpoint = endpointSendLocation
	case *ContactMessage:
		endpoint = endpointSendContact
	case *PollMessage:
		endpoint = endpointSendPoll
	case *DiceMessage:
		endpoint = endpointSendDice
	case *VenueMessage:
		endpoint = endpointSendVenue
	default:
		return nil, nil, fmt.Errorf("unsupported message type for telegram transformer: %T", msg)
	}

	apiURL := fmt.Sprintf(telegramAPIURLTemplate, account.APIKey, endpoint)

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal telegram payload: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      apiURL,
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

// newTelegramTransformer creates a new Telegram transformer instance.
func newTelegramTransformer() core.HTTPTransformer[*Account] {
	respCfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeJSON,
		CheckBody: true,
		Path:      "ok",
		Expect:    "true",
		MsgPath:   "description",
		Mode:      core.MatchEq,
	}

	tt := &telegramTransformer{}
	tt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeTelegram,
		"",
		respCfg,
		tt.transform,
	)

	return tt
}
