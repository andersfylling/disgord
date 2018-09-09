package websocket

import (
	"encoding/json"

	"github.com/json-iterator/go"
)

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

// unmarshalJSONIterator https://github.com/json-iterator/go
func unmarshalJSONIterator(data []byte, v interface{}) (err error) {
	err = jsoniter.Unmarshal(data, v)
	return
}

// unmarshalSTD standard GoLang implementation
func unmarshalSTD(data []byte, v interface{}) (err error) {
	err = json.Unmarshal(data, v)
	return
}

func unmarshal(data []byte, v interface{}) error {
	return unmarshalJSONIterator(data, v)
}
