package telegram

import "github.com/shellvon/go-sender/core"

// VenueMessage represents a venue message for Telegram.
// https://core.telegram.org/bots/api#sendvenue
// Required fields: ChatID, Latitude, Longitude, Title, Address.
// All caption-like fields are not supported by sendVenue.
// Note: sendVenue does NOT support caption/parse_mode; hence we embed BaseMessage, not MediaMessage.
type VenueMessage struct {
	BaseMessage

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`

	FoursquareID    string `json:"foursquare_id,omitempty"`
	FoursquareType  string `json:"foursquare_type,omitempty"`
	GooglePlaceID   string `json:"google_place_id,omitempty"`
	GooglePlaceType string `json:"google_place_type,omitempty"`
}

// NewVenueMessage is a convenience wrapper using the builder.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
//   - Only chat_id and latitude/longitude/title/address are required.
func NewVenueMessage(chatID string, lat, lon float64, title, address string) *VenueMessage {
	return Venue().Chat(chatID).
		Latitude(lat).
		Longitude(lon).
		Title(title).
		Address(address).
		Build()
}

func (m *VenueMessage) Validate() error {
	if m.ChatID == "" {
		return core.NewParamError("chat_id cannot be empty")
	}
	if m.Title == "" {
		return core.NewParamError("title cannot be empty")
	}
	if m.Address == "" {
		return core.NewParamError("address cannot be empty")
	}
	return nil
}
