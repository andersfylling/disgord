package discordws

const (
	// MaxReconnectTries restrict the number of reconnect attempt to the discord websocket service
	MaxReconnectTries = 5

	DefaultHTTPTimeout = 10

	LowestAPIVersion  = 6
	HighestAPIVersion = 6

	// BaseURL The base URL for all API requests
	BaseURL string = "https://discordapp.com/api"

	EncodingJSON string = "json"
	EncodingETF  string = "etf"

	LibVersionMajor        = 0
	LibVersionMinor        = 0
	LibVersionPatch        = 0
	LibName         string = "Disgord"
)

// Encodings legal
var Encodings = []string{
	EncodingJSON,
	EncodingETF,
}
