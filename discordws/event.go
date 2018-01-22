package discordws

//. "github.com/andersfylling/disgord/event"

//
// func (c *Client) onEvent(messageType int, packet []byte) (evt *GatewayPayload, err error) {
// 	logrus.Printf("Event: %+v\n\n", string(packet))
//
// 	// Your client should store the session_id from the Ready, and the sequence number of the last event it received.
//
// 	evt = &GatewayPayload{}
// 	err = json.Unmarshal(packet, evt)
// 	if err != nil {
// 		return
// 	}
//
// 	// ready event
// 	if evt.EventName == Ready {
// 		evt.Data = c
// 		err = json.Unmarshal(packet, evt)
// 		if err != nil {
// 			return
// 		}
//
// 		logrus.Printf("session_id: %s, trace: %s\n\n", c.SessionID, c.Trace)
// 	}
//
// 	// ping
// 	if evt.Op == 1 {
// 		c.wsMutex.Lock()
// 		err = c.conn.WriteJSON(struct {
// 			OP uint  `json:"op"`
// 			d  *uint `json:"d"`
// 		}{1, &c.sequenceNumber})
// 		c.Unlock()
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		return evt, nil
// 	}
//
// 	if evt.Op == 10 {
// 		return
// 	}
//
// 	return
// }
//
// // listenForDiscordEvents runs forever and reads every incoming websocket package from discord.
// // 												This must be executed as a goroutine to avoid blocking.
// func (dws *Client) listenForDiscordEvents(ws *websocket.Conn, connected <-chan struct{}) (err error) {
// 	for {
// 		messageType, packet, err := ws.ReadMessage()
// 		if err != nil {
// 			return err
// 		}
//
// 		select {
// 		case <-connected:
// 			return nil
// 		default:
// 			_, err = dws.onEvent(messageType, packet)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// }
