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

	logrus.Info("Listening for packets...")

	for {
		messageType, packet, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("error: %v", err)
			}
			close(c.disconnected)
			break
		}

		logrus.Infof("Recieved: %+v\n", string(packet))

		// TODO: zlib decompression support
		if messageType == websocket.BinaryMessage {
			logrus.Fatal("Cannot handle binary packets yet")
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

			err := c.conn.WriteJSON(message)
			if err != nil {
				logrus.Error(err)
			}
		case <-c.disconnected:
			logrus.Info("closing writePump")
			return
		}
	}
}
