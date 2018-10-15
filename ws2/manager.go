package ws2

import (
	"errors"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

const (
	maxReconnectTries = 5
)

// Event is dispatched by the socket layer after parsing and extracting Discord data from a incoming packet.
// This is the data structure used by Disgord for triggering handlers and channels with an event.
type Event struct {
	Name string
	Data []byte
}

type ManagerConfig struct {
	*DefaultClientConfig

	// for identify packets
	Browser             string
	Device              string
	GuildLargeThreshold uint
	ShardID             uint
	ShardCount          uint
}

type Manager struct {
	conf *ManagerConfig

	Client // interface

	eventChan chan *Event

	sessionID      string
	sequenceNumber uint
}

func (m *Manager) EventChan() <-chan *Event {
	return m.eventChan
}

func (m *Manager) reconnect() (err error) {
	_ = m.Disconnect()
	for try := 0; try <= maxReconnectTries; try++ {
		logrus.Debugf("Reconnect attempt #%d\n", try)
		err = m.Connect()
		if err == nil {
			logrus.Info("successfully reconnected")

			// send resume package

			break
			// TODO voice
		}
		if try == maxReconnectTries {
			err = errors.New("Too many reconnect attempts")
			return err
		}

		// wait N seconds
		logrus.Info("reconnect failed, trying again in N seconds; N = " + strconv.Itoa((try+3)*2))
		<-time.After(time.Duration((try+3)*2) * time.Second)
	}

	return
}

func (m *Manager) eventHandler(p *discordPacket) {
	// discord events
	// events that directly correlates to the socket layer, will be dealt with here. But still dispatched.

	// increment the sequence number for each event to make sure everything is synced with discord
	m.sequenceNumber++

	// always store the session id
	if p.EventName == event.Ready {
		m.updateSession(p)
	} else if p.EventName == event.Resume {
		// eh? debugging.
	} else if p.Op == opcode.DiscordEvent {
		// make sure we care about the event type
		if m.ListensForEvent(p.EventName) == -1 {
			return
		}
	}

	// dispatch event
	m.eventChan <- &Event{
		Name: p.EventName,
		Data: p.Data,
	}
} // end eventHandler()

// operation handler demultiplexer
func (m *Manager) operationHandlers() {
	logrus.Debug("Ready to receive operation codes...")
	for {
		select {
		case p, ok := <-m.Receive():
			if !ok {
				logrus.Debug("operationChan is dead..")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			switch opcode.ExtractFrom(p) {
			case opcode.DiscordEvent:
				c.eventHandler(p)
			case opcode.Reconnect:
				go m.reconnect()
			case opcode.InvalidSession:
				// invalid session. Must respond with a identify packet
				go func() {
					randomDelay := time.Second * time.Duration(rand.Intn(4)+1)
					<-time.After(randomDelay)
					err := sendIdentityPacket(m)
					if err != nil {
						logrus.Error(err)
					}
				}()
			case opcode.Heartbeat:
				// https://discordapp.com/developers/docs/topics/gateway#heartbeating
				_, _, snr := m.GetSocketInfo()
				_ = m.Emit(event.Heartbeat, snr)
			case opcode.Hello:
				// hello
				helloPk := &helloPacket{}
				err := httd.Unmarshal(p.Data, helloPk)
				if err != nil {
					logrus.Debug(err)
				}
				m.Lock()
				m.heartbeatInterval = helloPk.HeartbeatInterval
				m.Unlock()

				sendHelloPacket(m, p, helloPk.HeartbeatInterval)
			case opcode.HeartbeatAck:
				// heartbeat received
				m.Lock()
				m.lastHeartbeatAck = time.Now()
				m.Unlock()
			default:
				// unknown
				logrus.Debugf("Unknown operation: %+v\n", p)
			}

		case <-m.disconnected:
			logrus.Debug("exiting operation handler")
			return
		}
	}
}

func sendHelloPacket(m *Manager, gp *gatewayEvent, heartbeatInterval uint) {
	// TODO, this might create several idle goroutines..
	go pulsate(client, client.conn, client.disconnected)

	// send identify or resume packet
	if m.sessionID == "" && m.sequenceNumber == 0 {
		err := sendIdentityPacket(m)
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
		go func(client Pulsater, last time.Time, sent time.Time) {
			// TODO
			heartbeatResponseDeadline := (3 * time.Second) % (time.Duration(interval) * time.Millisecond)
			<-time.After(heartbeatResponseDeadline)
			if !client.HeartbeatWasRecieved(last) {
				logrus.Debug("heartbeat ACK was not received")
				client.HeartbeatAckMissingFix()
			} else {
				// update latency
				if c, ok := client.(*Client); ok {
					c.heartbeatLatency = c.lastHeartbeatAck.Sub(sent)
				}
			}
		}(client, last, time.Now())

		select {
		case <-ticker.C:
			continue
		case <-disconnected:
			logrus.Debug("Stopping pulse")
			return
		}
	}
}

func sendIdentityPacket(m *Manager) (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	identityPayload := struct {
		Token          string      `json:"token"`
		Properties     interface{} `json:"properties"`
		Compress       bool        `json:"compress"`
		LargeThreshold uint        `json:"large_threshold"`
		Shard          *[2]uint    `json:"shard,omitempty"`
		Presence       interface{} `json:"presence,omitempty"`
	}{
		Token: m.conf.Token,
		Properties: struct {
			OS      string `json:"$os"`
			Browser string `json:"$browser"`
			Device  string `json:"$device"`
		}{runtime.GOOS, m.conf.Browser, m.conf.Device},
		LargeThreshold: m.conf.GuildLargeThreshold,
		// Presence: struct {
		// 	Since  *uint       `json:"since"`
		// 	Game   interface{} `json:"game"`
		// 	Status string      `json:"status"`
		// 	AFK    bool        `json:"afk"`
		// }{Status: "online"},
	}

	if m.conf.ShardCount > 1 {
		identityPayload.Shard = &[2]uint{m.conf.ShardID, m.conf.ShardCount}
	}

	err = m.Emit(event.Identify, identityPayload)
	return
}
