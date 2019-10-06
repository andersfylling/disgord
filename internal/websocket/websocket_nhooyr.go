// +build !disgord_websocket_gorilla

package websocket

import (
	"context"
	"errors"
	"io"
	"net/http"

	httd2 "github.com/andersfylling/disgord/internal/httd"

	"golang.org/x/net/proxy"

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
	g.c, _, err = websocket.Dial(ctx, endpoint, &websocket.DialOptions{
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

func (g *nhooyr) WriteJSON(v interface{}) (err error) {
	// TODO: move unmarshalling out of here?
	var w io.WriteCloser
	w, err = g.c.Writer(context.Background(), websocket.MessageText)
	if err != nil {
		return err
	}
	err = httd2.JSONEncode(w, v)
	return
}

func (g *nhooyr) Close() (err error) {
	err = g.c.Close(websocket.StatusNormalClosure, "Bot is shutting down")
	g.c = nil
	return err
}

func (g *nhooyr) Read(ctx context.Context) (packet []byte, err error) {
	var messageType websocket.MessageType
	messageType, packet, err = g.c.Read(ctx)
	if err != nil {
		var closeErr *websocket.CloseError
		if errors.As(err, &closeErr) {
			err = &CloseErr{
				code: int(closeErr.Code),
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
