package discordws

import (
	"encoding/json"

	"github.com/andersfylling/disgord/event"
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

type payloadData []byte

func (pd *payloadData) UnmarshalJSON(data []byte) error {
	*pd = payloadData(data)
	return nil
}

func (pd *payloadData) ByteArr() []byte {
	return []byte(*pd)
}

type gatewayEvent struct {
	Op             uint          `json:"op"`
	Data           payloadData   `json:"d"`
	SequenceNumber uint          `json:"s"`
	EventName      event.KeyType `json:"t"`
}

type getGatewayResponse struct {
	URL    string `json:"url"`
	Shards uint   `json:"shards,omitempty"`
}

type helloPacket struct {
	HeartbeatInterval uint     `json:"heartbeat_interval"`
	Trace             []string `json:"_trace"`
}

type readyPacket struct {
	SessionID string   `json:"session_id"`
	Trace     []string `json:"_trace"`
}

type EventInterface interface {
	Name() event.KeyType
	Data() []byte
	Unmarshal(interface{}) error
}

type Event struct {
	content *gatewayEvent
}

func (evt *Event) Name() event.KeyType {
	return evt.content.EventName
}

func (evt *Event) Data() []byte {
	return evt.content.Data.ByteArr()
}

func (evt *Event) Unmarshal(v interface{}) error {
	return json.Unmarshal(evt.Data(), v)
}
