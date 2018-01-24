package discordws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	logrus.Debug("Listening for packets...")

	for {
		messageType, packet, err := c.conn.ReadMessage()
		if err != nil {
			var die bool
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// logrus.Errorf("error(%d): %v", messageType, err)
				die = true
			} else if c.disconnected == nil {
				// connection was closed
				die = true
			}

			if die {
				logrus.Debug("closing readPump")
				return
			}
		}

		logrus.Debugf("<-: %+v\n", string(packet))

		// TODO: zlib decompression support
		if messageType == websocket.BinaryMessage {
			logrus.Fatalf("Cannot handle packet type: %d", messageType)
		}

		// parse to gateway payload object
		evt := &gatewayEvent{}
		err = json.Unmarshal(packet, evt)
		if err != nil {
			logrus.Error(err)
		}

		// notify operation listeners
		c.operationChan <- evt
	}
}

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
