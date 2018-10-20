package websocket

import "net/http"

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
