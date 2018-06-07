package discordws

const (
	// MaxReconnectTries restrict the number of reconnect attempt to the discord websocket service
	MaxReconnectTries = 5

	MaxBytesSize = 4096

	DefaultHTTPTimeout = 10

	LowestAPIVersion  = 6
	HighestAPIVersion = 6

	// BaseURL The base URL for all API requests
	BaseURL = "https://discordapp.com/api"

	EncodingJSON = "json"
	EncodingETF  = "etf"

	LibVersionMajor = 0
	LibVersionMinor = 0
	LibVersionPatch = 0
	LibName         = "Disgord"
)

// Encodings legal
var Encodings = []string{
	EncodingJSON,
	EncodingETF,
}
