package websocket

import (
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake"
)

const (
	// MaxReconnectTries restrict the number of reconnect attempt to the discord websocket service
	MaxReconnectTries = 5

	//MaxBytesSize = 4096

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

type Snowflake = snowflake.Snowflake

func unmarshal(data []byte, v interface{}) error {
	return httd.Unmarshal(data, v)
}
