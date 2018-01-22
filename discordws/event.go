package discordws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (dws *Client) onEvent(messageType int, packet []byte) (event *GatewayPayload, err error) {
	logrus.Printf("Event: %+v\n\n", string(packet))

	// Your client should store the session_id from the Ready, and the sequence number of the last event it received.

	event = &GatewayPayload{}
	err = json.Unmarshal(packet, event)
	if err != nil {
		return
	}

	// ready event
	if event.EventName == "READY" {
		event.Data = dws
		err = json.Unmarshal(packet, event)
		if err != nil {
			return
		}

		logrus.Printf("session_id: %s, trace: %s\n\n", dws.SessionID, dws.Trace)
	}

	// ping
	if event.Op == 1 {
		dws.wsMutex.Lock()
		err = dws.WriteJSON(struct {
			OP uint  `json:"op"`
			d  *uint `json:"d"`
		}{1, &dws.sequenceNumber})
		dws.Unlock()
		if err != nil {
			return nil, err
		}

		return event, nil
	}

	if event.Op == 10 {
		return
	}

	return
}

// listenForDiscordEvents runs forever and reads every incoming websocket package from discord.
// 												This must be executed as a goroutine to avoid blocking.
func (dws *Client) listenForDiscordEvents(ws *websocket.Conn, connected <-chan struct{}) (err error) {
	for {
		messageType, packet, err := ws.ReadMessage()
		if err != nil {
			return err
		}

		select {
		case <-connected:
			return nil
		default:
			_, err = dws.onEvent(messageType, packet)
			if err != nil {
				return err
			}
		}
	}
}
