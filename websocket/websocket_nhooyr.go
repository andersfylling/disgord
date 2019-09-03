// +build !disgord_websocket_gorilla

package websocket

import (
	"context"
	"io"
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord/httd"
	"nhooyr.io/websocket"
)

func newConn(proxy proxy.Dialer, httpClient *http.Client) (*connection, error) {
	return &connection{
		httpClient: httpClient,
	}, nil
}

type connection struct {
	c            *websocket.Conn
	httpClient   *http.Client
	lastActivity lastActivity
}

var _ Conn = (*connection)(nil)

func (g *connection) Open(ctx context.Context, endpoint string, requestHeader http.Header) (err error) {
	defer g.lastActivity.Update()

	// establish ws connection
	g.c, _, err = websocket.Dial(ctx, endpoint, websocket.DialOptions{
		HTTPClient: g.httpClient,
		HTTPHeader: requestHeader,
	})
	if err != nil {
		if g.c != nil {
			_ = g.Close()
		}
		return err
	}

	g.c.SetReadLimit(32768 * 10000) // discord.. Can we add stream support?
	return
}

func (g *connection) WriteJSON(v interface{}) (err error) {
	defer g.lastActivity.Update()

	// TODO: move unmarshalling out of here?
	var w io.WriteCloser
	w, err = g.c.Writer(context.Background(), websocket.MessageText)
	if err != nil {
		return err
	}
	err = httd.JSONEncode(w, v)
	return
}

func (g *connection) Close() (err error) {
	defer g.lastActivity.Update()

	err = g.c.Close(websocket.StatusNormalClosure, "Bot is shutting down")
	g.c = nil
	return err
}

func (g *connection) Read(ctx context.Context) (packet []byte, err error) {
	defer g.lastActivity.Update()

	var messageType websocket.MessageType
	messageType, packet, err = g.c.Read(ctx)
	if err != nil {
		var closeErr *websocket.CloseError
		if xerrors.As(err, &closeErr) {
			err = &CloseErr{
				info: closeErr.Error(),
			}
		}
		return nil, err
	}

	if messageType == websocket.MessageBinary {
		packet, err = decompressBytes(packet)
	}
	return packet, nil
}

func (g *connection) Disconnected() bool {
	status := g.disconnected()

	return status
}

func (g *connection) Inactive() bool {
	return g.lastActivity.OlderThan(MaxReconnectDelay + 2*time.Minute)
}

func (g *connection) disconnected() bool {
	return g.c == nil
}

func (g *connection) InactiveSince() time.Time {
	return g.lastActivity.Time()
}
