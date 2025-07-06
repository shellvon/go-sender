package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

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

// telegramTransformer 实现 providers.HTTPTransformer[*Account].
type telegramTransformer struct{}

// newTelegramTransformer creates a new Telegram transformer instance.
func newTelegramTransformer() core.HTTPTransformer[*Account] {
	return &telegramTransformer{}
}

// CanTransform 判断是否为 Telegram 消息.
func (t *telegramTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeTelegram
}

// Transform 构造 Telegram HTTPRequestSpec
// 参数:
//   - ctx: 上下文
//   - msg: Telegram 消息体
//   - account: 账号配置
//
// 返回:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *telegramTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
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

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/%s", account.APIKey, endpoint)
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal telegram payload: %w", err)
	}
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      apiURL,
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}
	return reqSpec, handleTelegramResponse, nil
}

// handleTelegramResponse 处理 Telegram API 响应.
func handleTelegramResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("telegram API returned non-OK status: %d", statusCode)
	}
	var result struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code,omitempty"`
		Description string `json:"description,omitempty"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse telegram response: %w", err)
	}
	if !result.OK {
		if result.ErrorCode != 0 {
			return fmt.Errorf("telegram API error %d: %s", result.ErrorCode, result.Description)
		}
		return fmt.Errorf("telegram API request failed: %s", result.Description)
	}
	return nil
}
