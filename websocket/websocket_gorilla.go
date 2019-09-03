// +build disgord_websocket_gorilla

package websocket

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
)

func newConn(proxy proxy.Dialer, httpClient *http.Client) (*connection, error) {
	return &connection{
		proxy: proxy,
	}, nil
}

// rwc is a wrapper for the Conn interface (not net.Conn).
// Interface can be found at https://golang.org/pkg/net/#Conn
// See original code at https://github.com/gorilla/websocket/issues/282
type connection struct {
	c            *websocket.Conn
	proxy        proxy.Dialer
	lastActivity lastActivity
}

var _ Conn = (*connection)(nil)

func (g *connection) Open(ctx context.Context, endpoint string, requestHeader http.Header) (err error) {
	defer g.lastActivity.Update()

	// by default we use connection's websocket dialer here, but if the passed http client uses a custom transport
	// we make sure we open the websocket over the same transport/proxy, in case the user uses this
	dialer := websocket.DefaultDialer
	if g.proxy != nil {
		dialer = &websocket.Dialer{
			NetDial: g.proxy.Dial,
		}
	}

	// establish ws connection
	g.c, _, err = dialer.Dial(endpoint, requestHeader)
	if err != nil && !g.Disconnected() {
		_ = g.Close()
	}
	return
}

func (g *connection) WriteJSON(v interface{}) (err error) {
	defer g.lastActivity.Update()

	// TODO: move unmarshalling out of here?
	var w io.WriteCloser
	w, err = g.c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	err = httd.JSONEncode(w, v)
	return
}

func (g *connection) Close() (err error) {
	defer g.lastActivity.Update()

	err = g.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	err2 := g.c.Close()
	g.c = nil

	if err == nil && err2 != nil {
		return err2
	}
	return err
}

func (g *connection) Read(ctx context.Context) (packet []byte, err error) {
	if g.Disconnected() {
		// this gets triggered when losing internet connection -> trying to reconnect for a while -> re-establishing a connection
		// as discord then sends a invalid session package and disgord tries to reconnect again, a panic takes place.
		// this check is a tmp hack to fix that, as the actual issue is not clearly understood/defined yet.
		err = errors.New("no connection is established. Can not read new messages")
		return
	}
	defer g.lastActivity.Update()

	var messageType int
	messageType, packet, err = g.c.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			err = &CloseErr{
				info: err.Error(),
			}
		}

		return nil, err
	}

	if messageType == websocket.BinaryMessage {
		packet, err = decompressBytes(packet)
	}
	return packet, nil
}

func (g *connection) Disconnected() bool {
	return g.c == nil
}

func (g *connection) Inactive() bool {
	return g.lastActivity.OlderThan(MaxReconnectDelay + 2*time.Minute)
}

func (g *connection) InactiveSince() time.Time {
	return g.lastActivity.Time()
}
