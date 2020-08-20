package gateway

// TODO: merge websocket_nhooyr.go with the client.go. Figure out mocking for the client as well.

import (
	"context"
	"errors"
	"github.com/andersfylling/disgord/json"
	"io"
	"net/http"

	"go.uber.org/atomic"

	"nhooyr.io/websocket"
)

func newConn(httpClient *http.Client) (Conn, error) {
	return &nhooyr{
		httpClient: httpClient,
	}, nil
}

type nhooyr struct {
	c           *websocket.Conn
	httpClient  *http.Client
	isConnected atomic.Bool
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
	g.isConnected.Store(true)

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

	// TODO: implement custom json handler - switch to new gateway project
	err1 := json.NewEncoder(w).Encode(v)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

func (g *nhooyr) Close() (err error) {
	err = g.c.Close(websocket.StatusNormalClosure, "Bot is shutting down")
	if !g.isConnected.Load() {
		err = nil // discard error if we're already closed, should be a noop anyways
	}
	g.isConnected.Store(false)
	return err
}

func (g *nhooyr) Read(ctx context.Context) (packet []byte, err error) {
	var messageType websocket.MessageType
	messageType, packet, err = g.c.Read(ctx)
	if err != nil {
		// Cancelling Read by ctx results in closed WS, see issue
		// https://github.com/nhooyr/websocket/issues/242
		if ctx.Err() != nil && errors.Is(err, context.Canceled) {
			g.isConnected.Store(false)
			return nil, context.Canceled
		}
		var closeErr websocket.CloseError
		if errors.As(err, &closeErr) {
			g.isConnected.Store(false)
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
	return !g.isConnected.Load()
}

var _ Conn = (*nhooyr)(nil)
