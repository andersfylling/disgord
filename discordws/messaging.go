package discordws

import (
	"encoding/json"
	"fmt"

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

	logrus.Info("Listening for packets...")

	for {
		select {
		case <-c.disconnected:
			logrus.Info("closing readPump")
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

			logrus.Infof("<-: %+v\n", string(packet))

			// TODO: zlib decompression support
			if messageType != websocket.TextMessage {
				logrus.Fatalf("Cannot handle pacaket type: %d", messageType)
			}

			// parse to gateway payload object
			evt := GatewayPayload{}
			err = json.Unmarshal(packet, &evt)
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

	logrus.Info("Ready to send packets...")

	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				// The hub closed the channel.
				logrus.Error("oh no...")
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
				fmt.Printf("->: %+v\n", message)
			} else {
				fmt.Printf("->: %+v\n", string(b))
			}
		case <-c.disconnected:
			logrus.Info("closing writePump")
			return
		}
	}
}
