package gateway

import (
	"context"
	"net/http"

	"github.com/andersfylling/disgord/internal/util"
)

type Snowflake = util.Snowflake

type Conn interface {
	Close() error
	Open(ctx context.Context, endpoint string, requestHeader http.Header) error
	WriteJSON(v interface{}) error
	Read(ctx context.Context) (packet []byte, err error)

	Disconnected() bool
}

type CloseErr struct {
	code int
	info string
}

func (e *CloseErr) Error() string {
	return e.info
}

// WebsocketErr is used internally when the websocket package returns an error. It does not represent a Discord error!
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
