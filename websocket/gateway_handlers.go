package websocket

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (c *Client) opHandlerEvt(gp *gatewayEvent) {
	// discord events
	// events that directly correlates to the socket layer, will be dealt with here. But still dispatched.

	// increment the sequence number for each event to make sure everything is synced with discord
	c.incrementSequenceNumber()

	// always store the session id
	if gp.EventName == event.Ready {
		c.updateSession(gp)
	} else if gp.EventName == event.Resume {
		// eh? debugging.
	} else if gp.Op == opcode.DiscordEvent {
		// make sure we care about the event type
		if c.ListensForEvent(gp.EventName) == -1 {
			return
		}
	}

	// dispatch events
	eventPkt := &Event{gp}
	c.discordWSEventChan <- eventPkt
} // end opHandlerEvt()

// operation handler demultiplexer
func (c *Client) operationHandlers() {
	logrus.Debug("Ready to receive operation codes...")
	for {
		select {
		case gp, ok := <-c.operationChan:
			if !ok {
				logrus.Debug("operationChan is dead..")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			switch opcode.ExtractFrom(gp) {
			case opcode.DiscordEvent:
				c.opHandlerEvt(gp)
			case opcode.Reconnect:
				// reconnect
				c.Disconnect()
				go c.reconnect()
			case opcode.InvalidSession:
				// invalid session. Must respond with a identify packet
				go func() {
					randomDelay := time.Second * time.Duration(rand.Intn(4)+1)
					<-time.After(randomDelay)
					err := sendIdentityPacket(c)
					if err != nil {
						logrus.Error(err)
					}
				}()
			case opcode.Heartbeat:
				// https://discordapp.com/developers/docs/topics/gateway#heartbeating
				_, _, snr := c.GetSocketInfo()
				c.Emit(event.Heartbeat, snr)
			case opcode.Hello:
				// hello
				helloPk := &helloPacket{}
				err := unmarshal(gp.Data.ByteArr(), helloPk)
				if err != nil {
					logrus.Debug(err)
				}
				c.Lock()
				c.heartbeatInterval = helloPk.HeartbeatInterval
				c.Unlock()

				sendHelloPacket(c, gp, helloPk.HeartbeatInterval)
			case opcode.HeartbeatAck:
				// heartbeat received
				c.Lock()
				c.lastHeartbeatAck = time.Now()
				c.Unlock()
			default:
				// unknown
				logrus.Debugf("Unknown operation: %+v\n", gp)
			}

		case <-c.disconnected:
			logrus.Debug("exiting operation handler")
			return
		}
	}
}

func sendHelloPacket(client *Client, gp *gatewayEvent, heartbeatInterval uint) {
	// TODO, this might create several idle goroutines..
	go pulsate(client, client.conn, client.disconnected)

	// send identify or resume packet
	if client.SessionID == "" && client.sequenceNumber == 0 {
		err := sendIdentityPacket(client)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	client.RLock()
	token := client.conf.Token
	session := client.SessionID
	sequence := client.sequenceNumber
	client.RUnlock()

	client.Emit(event.Resume, struct {
		Token      string `json:"token"`
		SessionID  string `json:"session_id"`
		SequenceNr *uint  `json:"seq"`
	}{token, session, &sequence})
}

func pulsate(client Pulsater, ws *websocket.Conn, disconnected chan struct{}) {
	serviceID := uint8(rand.Intn(255)) // uint8 cap
	if !client.AllowedToStartPulsating(serviceID) {
		return
	}
	defer client.StopPulsating(serviceID)

	ticker := time.NewTicker(time.Millisecond * time.Duration(client.HeartbeatInterval()))
	defer ticker.Stop()

	var last time.Time
	var interval uint
	var snr uint
	for {
		last, interval, snr = client.GetSocketInfo()
		client.SendHeartbeat(snr)

		// verify the heartbeat ACK
		go func(client Pulsater, last time.Time) {
			// TODO
			heartbeatResponseDeadline := (3 * time.Second) % (time.Duration(interval) * time.Millisecond)
			<-time.After(heartbeatResponseDeadline)
			if !client.HeartbeatWasRecieved(last) {
				logrus.Debug("heartbeat ACK was not received")
				client.HeartbeatAckMissingFix()
			}
		}(client, last)

		select {
		case <-ticker.C:
			continue
		case <-disconnected:
			logrus.Debug("Stopping pulse")
			return
		}
	}
}

func sendIdentityPacket(client *Client) (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	identityPayload := struct {
		Token          string      `json:"token"`
		Properties     interface{} `json:"properties"`
		Compress       bool        `json:"compress"`
		LargeThreshold uint        `json:"large_threshold"`
		Shard          *[2]uint    `json:"shard,omitempty"`
		Presence       interface{} `json:"presence,omitempty"`
	}{
		Token: client.conf.Token,
		Properties: struct {
			OS      string `json:"$os"`
			Browser string `json:"$browser"`
			Device  string `json:"$device"`
		}{runtime.GOOS, client.conf.Browser, client.conf.Device},
		LargeThreshold: client.conf.GuildLargeThreshold,
		// Presence: struct {
		// 	Since  *uint       `json:"since"`
		// 	Game   interface{} `json:"game"`
		// 	Status string      `json:"status"`
		// 	AFK    bool        `json:"afk"`
		// }{Status: "online"},
	}

	if client.ShardCount > 1 {
		identityPayload.Shard = &[2]uint{uint(client.ShardID), client.ShardCount}
	}

	err = client.Emit(event.Identify, identityPayload)
	return
}
