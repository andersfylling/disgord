package websocket

import (
	"net/http"

	"github.com/andersfylling/disgord/depalias"
)

type Snowflake = depalias.Snowflake

type Conn interface {
	Close() error
	Open(endpoint string, requestHeader http.Header) error
	WriteJSON(v interface{}) error
	Read() (packet []byte, err error)

	Disconnected() bool
}

type ErrorUnexpectedClose struct {
	info string
}

func (e *ErrorUnexpectedClose) Error() string {
	return e.info
}

// WebsocketErr is used internally when the websocket package returns an error. It does not represent a Discord error(!)
type WebsocketErr struct {
	ID      uint
	message string
}

func (e *WebsocketErr) Error() string {
	return e.message
}

const (
	encodingJSON = "json"
)

// Choreographic programming.. TODO: rename channels and structs

type A chan B
type B chan *K

// K is used to get the Connect permission from the shard manager
type K struct {
	Release B
	Key     interface{}
}
