package websocket

import (
	"bytes"
	"compress/zlib"
	"io"

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

func decompressBytes(input []byte) (output []byte, err error) {
	b := bytes.NewReader(input)
	var r io.ReadCloser

	r, err = zlib.NewReader(b)
	if err != nil {
		return
	}
	defer r.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(r)
	if err != nil {
		return
	}

	output = buffer.Bytes()
	return
}
