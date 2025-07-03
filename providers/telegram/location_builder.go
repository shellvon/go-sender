package telegram

// LocationBuilder constructs Telegram location messages.
// Example:
//   msg := telegram.Location().
//            Chat("123").
//            Coordinates(51.5074, -0.1278).
//            HorizontalAccuracy(15).
//            Build()

type LocationBuilder struct {
	*baseBuilder[*LocationBuilder]

	latitude  float64
	longitude float64

	horizontalAccuracy   float64
	livePeriod           int
	heading              int
	proximityAlertRadius int
}

// Location returns a new LocationBuilder.
func Location() *LocationBuilder {
	b := &LocationBuilder{}
	b.baseBuilder = &baseBuilder[*LocationBuilder]{self: b}
	return b
}

// Coordinates sets latitude and longitude. Required.
// Based on SendLocationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendlocation
// Latitude and longitude of the location.
func (b *LocationBuilder) Coordinates(lat, lon float64) *LocationBuilder {
	b.latitude = lat
	b.longitude = lon
	return b
}

// HorizontalAccuracy sets accuracy in meters (0-1500).
// Based on SendLocationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendlocation
// The radius of uncertainty for the location, measured in meters; 0-1500.
func (b *LocationBuilder) HorizontalAccuracy(acc float64) *LocationBuilder {
	b.horizontalAccuracy = acc
	return b
}

// LivePeriod sets live location period seconds (60-86400).
// Based on SendLocationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendlocation
// Period in seconds for which the location can be updated, should be between 60 and 86400.
func (b *LocationBuilder) LivePeriod(sec int) *LocationBuilder {
	b.livePeriod = sec
	return b
}

// Heading sets movement direction degrees (1-360).
// Based on SendLocationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendlocation
// For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.
func (b *LocationBuilder) Heading(deg int) *LocationBuilder {
	b.heading = deg
	return b
}

// ProximityAlertRadius sets proximity alert radius meters (1-100000).
// Based on SendLocationParams from Telegram Bot API
// https://core.telegram.org/bots/api#sendlocation
// For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters.
// Must be between 1 and 100000 if specified.
func (b *LocationBuilder) ProximityAlertRadius(radius int) *LocationBuilder {
	b.proximityAlertRadius = radius
	return b
}

// Build assembles the *LocationMessage.
func (b *LocationBuilder) Build() *LocationMessage {
	msg := &LocationMessage{
		BaseMessage:          b.baseBuilder.toBaseMessage(TypeLocation),
		Latitude:             b.latitude,
		Longitude:            b.longitude,
		HorizontalAccuracy:   b.horizontalAccuracy,
		LivePeriod:           b.livePeriod,
		Heading:              b.heading,
		ProximityAlertRadius: b.proximityAlertRadius,
	}
	return msg
}
