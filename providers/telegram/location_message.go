package telegram

import "github.com/shellvon/go-sender/core"

// LocationMessage represents a location message for Telegram
// Based on SendLocationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendlocation
type LocationMessage struct {
	BaseMessage

	// Latitude of the location
	Latitude float64 `json:"latitude"`

	// Longitude of the location
	Longitude float64 `json:"longitude"`

	// The radius of uncertainty for the location, measured in meters; 0-1500
	HorizontalAccuracy float64 `json:"horizontal_accuracy,omitempty"`

	// Period in seconds for which the location can be updated, should be between 60 and 86400
	LivePeriod int `json:"live_period,omitempty"`

	// For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.
	Heading int `json:"heading,omitempty"`

	// For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters.
	// Must be between 1 and 100000 if specified.
	ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
}

// NewLocationMessage creates a new LocationMessage instance.
// Based on SendLocationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendlocation
//   - Only chat_id and latitude/longitude are required.
func NewLocationMessage(chatID string, latitude, longitude float64) *LocationMessage {
	return &LocationMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeLocation,
			ChatID:  chatID,
		},
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func (m *LocationMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	return nil
}
