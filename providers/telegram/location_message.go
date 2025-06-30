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
func NewLocationMessage(chatID string, latitude, longitude float64, opts ...interface{}) *LocationMessage {
	msg := &LocationMessage{
		BaseMessage: BaseMessage{
			MsgType: TypeLocation,
			ChatID:  chatID,
		},
		Latitude:  latitude,
		Longitude: longitude,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case LocationMessageOption:
			o(msg)
		case MessageOption:
			o(msg)
		}
	}
	return msg
}

func (m *LocationMessage) GetBase() *BaseMessage {
	return &m.BaseMessage
}

func (m *LocationMessage) ProviderType() core.ProviderType {
	return core.ProviderTypeTelegram
}

func (m *LocationMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	return nil
}

type LocationMessageOption func(*LocationMessage)

// WithLocationHorizontalAccuracy sets the radius of uncertainty for the location
// Measured in meters, should be between 0-1500.
func WithLocationHorizontalAccuracy(accuracy float64) LocationMessageOption {
	return func(m *LocationMessage) { m.HorizontalAccuracy = accuracy }
}

// WithLocationLivePeriod sets the period for which the location can be updated
// Should be between 60 and 86400 seconds.
func WithLocationLivePeriod(period int) LocationMessageOption {
	return func(m *LocationMessage) { m.LivePeriod = period }
}

// WithLocationHeading sets the direction in which the user is moving
// Must be between 1 and 360 degrees if specified.
func WithLocationHeading(heading int) LocationMessageOption {
	return func(m *LocationMessage) { m.Heading = heading }
}

// WithLocationProximityAlertRadius sets the maximum distance for proximity alerts
// Must be between 1 and 100000 meters if specified.
func WithLocationProximityAlertRadius(radius int) LocationMessageOption {
	return func(m *LocationMessage) { m.ProximityAlertRadius = radius }
}
