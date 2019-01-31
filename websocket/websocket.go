package websocket

import (
	"net/http"

	"github.com/pkg/errors"
)

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

// K is used to get the connect permission from the shard manager
type K struct {
	Release B
	Key     interface{}
}

func requestConnectPermission(c *Client) error {
	c.Debug("trying to get connect permission")
	b := make(B)
	defer close(b)
	c.a <- b
	c.Info("waiting")
	var ok bool
	select {
	case c.K, ok = <-b:
		if !ok || c.K == nil {
			c.Debug("unable to get connect permission")
			return errors.New("channel closed or K was nil")
		}
		c.Debug("got connect permission")
	case <-c.shutdown:
	}

	return nil
}

func releaseConnectPermission(c *Client) error {
	if c.K == nil {
		return errors.New("K has not been granted yet")
	}

	c.K.Release <- c.K
	c.K = nil
	return nil
}

// diagnosing
const DiagnosePath = "diagnose-report"
const DiagnosePath_packets = "diagnose-report/packets"
