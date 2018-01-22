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
		messageType, packet, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Infof(">>> %+v\n", string(packet))
				logrus.Errorf("error: %v", err)
			}
			//close(c.disconnected)
			//break
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

		select {
		case <-c.disconnected:
			logrus.Info("closing readPump")
			return
		default:
			continue
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
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

			b, err := json.Marshal(&message)
			if err != nil {
				logrus.Error(err)
				continue
			}

			c.wsMutex.Lock()
			err = c.conn.WriteJSON(b)
			c.wsMutex.Unlock()
			if err != nil {
				logrus.Error(err)
			}
			fmt.Printf("->: %+v\n", message)
		case <-c.disconnected:
			logrus.Info("closing writePump")
			return
		}
	}
}
