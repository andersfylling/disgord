// +build !disgord_websocket_gorilla

package websocket

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord/httd"
	"nhooyr.io/websocket"
)

func newConn(proxy proxy.Dialer, httpClient *http.Client) (Conn, error) {
	return &nhooyr{
		httpClient: httpClient,
	}, nil
}

type nhooyr struct {
	c          *websocket.Conn
	httpClient *http.Client
}

func (g *nhooyr) Open(ctx context.Context, endpoint string, requestHeader http.Header) (err error) {
	// establish ws connection
	g.c, _, err = websocket.Dial(ctx, endpoint, websocket.DialOptions{
		HTTPClient: g.httpClient,
		HTTPHeader: requestHeader,
	})
	g.c.SetReadLimit(32768 * 10000) // discord.. Can we add stream support?
	return
}

func (g *nhooyr) WriteJSON(v interface{}) (err error) {
	// TODO: move unmarshalling out of here?
	var w io.WriteCloser
	w, err = g.c.Writer(context.Background(), websocket.MessageText)
	if err != nil {
		return err
	}
	err = httd.JSONEncode(w, v)
	return
}

func (g *nhooyr) Close() (err error) {
	err = g.c.Close(websocket.StatusNormalClosure, "Bot is shutting down")
	g.c = nil
	return err
}

func (g *nhooyr) Read() (packet []byte, err error) {
	var messageType websocket.MessageType
	messageType, packet, err = g.c.Read(context.Background())
	if err != nil {
		if closeErr, ok := err.(*websocket.CloseError); ok {
			err = &ErrorUnexpectedClose{
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

func (g *nhooyr) Disconnected() bool {
	status := g.disconnected()

	return status
}

func (g *nhooyr) disconnected() bool {
	return g.c == nil
}

var _ Conn = (*nhooyr)(nil)
