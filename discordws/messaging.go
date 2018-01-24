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
	defer func() {
		c.conn.Close()
	}()

	logrus.Debug("Listening for packets...")

	for {
		select {
		case <-c.disconnected:
			logrus.Debug("closing readPump")
			return
		default:
			messageType, packet, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logrus.Errorf("error: %v", err)
				}
				close(c.disconnected)
				continue
			}

			logrus.Debugf("<-: %+v\n", string(packet))

			// TODO: zlib decompression support
			if messageType != websocket.TextMessage {
				logrus.Fatalf("Cannot handle pacaket type: %d", messageType)
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
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(ws *websocket.Conn) {
	defer func() {
		c.conn.Close()
	}()

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
