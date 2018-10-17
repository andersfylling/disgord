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
