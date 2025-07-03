package telegram

// VenueBuilder constructs Telegram venue messages.
// Example:
//   msg := telegram.Venue().
//            Chat("123").
//            Latitude(40.7).
//            Longitude(-74.0).
//            Title("HQ").
//            Address("Wall Street").
//            Build()

type VenueBuilder struct {
	*baseBuilder[*VenueBuilder]

	latitude  float64
	longitude float64
	title     string
	address   string

	foursquareID    string
	foursquareType  string
	googlePlaceID   string
	googlePlaceType string
}

// Venue returns a new VenueBuilder instance.
func Venue() *VenueBuilder {
	b := &VenueBuilder{}
	b.baseBuilder = &baseBuilder[*VenueBuilder]{self: b}
	return b
}

// Latitude sets the latitude of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Latitude of the location.
func (b *VenueBuilder) Latitude(lat float64) *VenueBuilder { b.latitude = lat; return b }

// Longitude sets the longitude of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Longitude of the location.
func (b *VenueBuilder) Longitude(lon float64) *VenueBuilder { b.longitude = lon; return b }

// Title sets the title of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Title of the venue.
func (b *VenueBuilder) Title(t string) *VenueBuilder { b.title = t; return b }

// Address sets the address of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Address of the venue.
func (b *VenueBuilder) Address(a string) *VenueBuilder { b.address = a; return b }

// FoursquareID sets the Foursquare ID of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Foursquare ID of the venue.
func (b *VenueBuilder) FoursquareID(id string) *VenueBuilder { b.foursquareID = id; return b }

// FoursquareType sets the Foursquare type of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Foursquare type of the venue.
func (b *VenueBuilder) FoursquareType(t string) *VenueBuilder { b.foursquareType = t; return b }

// GooglePlaceID sets the Google Place ID of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Google Place ID of the venue.
func (b *VenueBuilder) GooglePlaceID(id string) *VenueBuilder { b.googlePlaceID = id; return b }

// GooglePlaceType sets the Google Place type of the venue.
// Based on SendVenueParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendvenue
// Google Place type of the venue.
func (b *VenueBuilder) GooglePlaceType(t string) *VenueBuilder { b.googlePlaceType = t; return b }

// Build creates the *VenueMessage.
func (b *VenueBuilder) Build() *VenueMessage {
	return &VenueMessage{
		BaseMessage:     b.baseBuilder.toBaseMessage(TypeVenue),
		Latitude:        b.latitude,
		Longitude:       b.longitude,
		Title:           b.title,
		Address:         b.address,
		FoursquareID:    b.foursquareID,
		FoursquareType:  b.foursquareType,
		GooglePlaceID:   b.googlePlaceID,
		GooglePlaceType: b.googlePlaceType,
	}
}
