package websocket

// TODO: if we add any other websocket packages, add build constraints to this file.

import (
	"errors"
	"io"
	"net/http"

	"github.com/andersfylling/disgord/httd"
	"github.com/gorilla/websocket"
)

func newConn(HTTPClient *http.Client) (Conn, error) {
	return &gorilla{
		HTTPClient: HTTPClient,
	}, nil
}

// rwc is a wrapper for the Conn interface (not net.Conn).
// Interface can be found at https://golang.org/pkg/net/#Conn
// See original code at https://github.com/gorilla/websocket/issues/282
type gorilla struct {
	c          *websocket.Conn
	HTTPClient *http.Client
}

func (g *gorilla) Open(endpoint string, requestHeader http.Header) (err error) {
	// by default we use gorilla's websocket dialer here, but if the passed http client uses a custom transport
	// we make sure we open the websocket over the same transport/proxy, in case the user uses this
	dialer := websocket.DefaultDialer
	if t, ok := g.HTTPClient.Transport.(*http.Transport); ok {
		dialer = &websocket.Dialer{
			HandshakeTimeout: dialer.HandshakeTimeout,
			Proxy:            t.Proxy,
			NetDialContext:   t.DialContext,
			NetDial:          t.Dial, // even though Dial is deprecated in http.Transport, it isn't in websocket
		}
	}

	// establish ws connection
	g.c, _, err = dialer.Dial(endpoint, requestHeader)
	return
}

func (g *gorilla) WriteJSON(v interface{}) (err error) {
	// TODO: move unmarshalling out of here?
	var w io.WriteCloser
	w, err = g.c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	err = httd.JSONEncode(w, v)
	return
}

func (g *gorilla) Close() (err error) {
	err = g.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	g.c = nil
	return
}

func (g *gorilla) Read() (packet []byte, err error) {
	if g.Disconnected() {
		// this gets triggered when loosing interconnection -> trying to reconnect for a while -> re-establishing a connection
		// as discord then sends a invalid session package and disgord tries to reconnect again, a panic takes place.
		// this check is a tmp hack to fix that, as the actual issue is not clearly understood/defined yet.
		err = errors.New("no connection is established. Can not read new messages")
		return
	}
	var messageType int
	messageType, packet, err = g.c.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			err = &ErrorUnexpectedClose{
				info: err.Error(),
			}
		}

		return
	}

	if messageType == websocket.BinaryMessage {
		packet, err = decompressBytes(packet)
	}
	return
}

func (g *gorilla) Disconnected() bool {
	return g.c == nil
}

var _ Conn = (*gorilla)(nil)
