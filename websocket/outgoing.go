package websocket

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(ws *websocket.Conn) {
	logrus.Debug("Ready to send packets...")

	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.wsMutex.Lock()
			err := ws.WriteJSON(&message)
			c.wsMutex.Unlock()
			if err != nil {
				logrus.Error(err)
			}

			b, err := json.MarshalIndent(&message, "", "  ")
			if err != nil {
				logrus.Debugf("->: %+v\n", message)
			} else {
				logrus.Debugf("->: %+v\n", string(b))
			}
		case <-c.disconnected:
			logrus.Debug("closing writePump")
			return
		}
	}
}

type gatewayPayload struct {
	Op             uint        `json:"op"`
	Data           interface{} `json:"d"`
	SequenceNumber uint        `json:"s,omitempty"`
	EventName      string      `json:"t,omitempty"`
}

func (gp *gatewayPayload) DataToByteArr() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(gp.Data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
