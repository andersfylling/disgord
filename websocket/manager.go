package websocket

import (
	"errors"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"github.com/sirupsen/logrus"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	maxReconnectTries = 5
)

// NewManager creates a new socket client manager for handling behavior and Discord events. Note that this
// function initiates a go routine.
func NewManager(config *ManagerConfig) (manager *Manager, err error) {
	var client Client
	client, err = NewDefaultClient(&config.DefaultClientConfig)
	if err != nil {
		return
	}

	manager = &Manager{
		conf:      config,
		shutdown:  make(chan interface{}),
		restart:   make(chan interface{}),
		Client:    client,
		eventChan: make(chan *Event),
	}
	go manager.operationHandlers()

	return
}

// Event is dispatched by the socket layer after parsing and extracting Discord data from a incoming packet.
// This is the data structure used by Disgord for triggering handlers and channels with an event.
type Event struct {
	Name string
	Data []byte
}

type ManagerConfig struct {
	DefaultClientConfig

	// for identify packets
	Browser             string
	Device              string
	GuildLargeThreshold uint
	ShardID             uint
	ShardCount          uint
}

type Manager struct {
	sync.RWMutex
	conf     *ManagerConfig
	shutdown chan interface{}
	restart  chan interface{}

	Client // interface

	eventChan     chan *Event
	trackedEvents []string
	evtMutex      sync.RWMutex

	heartbeatInterval uint
	heartbeatLatency  time.Duration
	lastHeartbeatAck  time.Time

	sessionID      string
	trace          []string
	sequenceNumber uint

	pulsating  uint8
	pulseMutex sync.Mutex
}

// HeartbeatLatency get the time diff between sending a heartbeat and Discord replying with a heartbeat ack
func (m *Manager) HeartbeatLatency() (duration time.Duration, err error) {
	duration = m.heartbeatLatency
	if duration == 0 {
		err = errors.New("latency not determined yet")
	}

	return
}

// RegisterEvent tells the socket layer which event types are of interest. Any event that are not registered
// will be discarded once the socket info is extracted from the event.
func (m *Manager) RegisterEvent(event string) {
	m.evtMutex.Lock()
	defer m.evtMutex.Unlock()

	for i := range m.trackedEvents {
		if event == m.trackedEvents[i] {
			return
		}
	}

	m.trackedEvents = append(m.trackedEvents, event)
}

// RemoveEvent removes an event type from the registry. This will cause the event type to be discarded
// by the socket layer.
func (m *Manager) RemoveEvent(event string) {
	m.evtMutex.Lock()
	defer m.evtMutex.Unlock()

	for i := range m.trackedEvents {
		if event == m.trackedEvents[i] {
			m.trackedEvents[i] = m.trackedEvents[len(m.trackedEvents)-1]
			m.trackedEvents = m.trackedEvents[:len(m.trackedEvents)-1]
			break
		}
	}
	return
}

func (m *Manager) EventChan() <-chan *Event {
	return m.eventChan
}

func (m *Manager) Shutdown() (err error) {
	m.Disconnect()
	close(m.shutdown)
	return
}

func (m *Manager) reconnect() (err error) {
	close(m.restart)
	_ = m.Disconnect()
	for try := 0; try <= maxReconnectTries; try++ {
		logrus.Debugf("Reconnect attempt #%d\n", try)
		err = m.Connect()
		if err == nil {
			logrus.Info("successfully reconnected")
			break
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

	// validate the sequence numbers
	if p.SequenceNumber != m.sequenceNumber {
		m.sequenceNumber--
		m.reconnect()
	}

	if p.EventName == event.Ready {
		// always store the session id & update the trace content
		ready := readyPacket{}
		err := httd.Unmarshal(p.Data, &ready)
		if err != nil {
			logrus.Error(err)
		}

		m.Lock()
		m.sessionID = ready.SessionID
		m.trace = ready.Trace
		m.Unlock()
	} else if p.EventName == event.Resume {
		// eh? debugging.
		// TODO
	} else if p.Op == opcode.DiscordEvent && !m.eventOfInterest(p.EventName) {
		return
	}

	// dispatch event
	m.eventChan <- &Event{
		Name: p.EventName,
		Data: p.Data,
	}
} // end eventHandler()

func (m *Manager) eventOfInterest(name string) bool {
	m.evtMutex.RLock()
	defer m.evtMutex.RUnlock()

	for i := range m.trackedEvents {
		if name == m.trackedEvents[i] {
			return true
		}
	}

	return false
}

// operation handler demultiplexer
func (m *Manager) operationHandlers() {
	logrus.Debug("Ready to receive operation codes...")
	for {
		var p *discordPacket
		var open bool
		select {
		case p, open = <-m.Receive():
			if !open {
				logrus.Debug("operationChan is dead..")
				return
			}
		// case <-m.restart:
		case <-m.shutdown:
			logrus.Debug("exiting operation handler")
			return
		}

		// new packet that must be handled by it's Discord operation code
		switch p.Op {
		case opcode.DiscordEvent:
			m.eventHandler(p)
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
			_ = m.Emit(event.Heartbeat, m.sequenceNumber)
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

			m.sendHelloPacket()
		case opcode.HeartbeatAck:
			// heartbeat received
			m.Lock()
			m.lastHeartbeatAck = time.Now()
			m.Unlock()
		default:
			// unknown
			logrus.Debugf("Unknown operation: %+v\n", p)
		}
	}
}

func (m *Manager) sendHelloPacket() {
	// TODO, this might create several idle goroutines..
	go m.pulsate()

	err := sendIdentityPacket(m)
	if err != nil {
		logrus.Error(err)
	}

	// if this is a new connection we can drop the resume packet
	if m.sessionID == "" && m.sequenceNumber == 0 {
		return
	}

	m.RLock()
	token := m.conf.Token
	session := m.sessionID
	sequence := m.sequenceNumber
	m.RUnlock()

	m.Emit(event.Resume, struct {
		Token      string `json:"token"`
		SessionID  string `json:"session_id"`
		SequenceNr *uint  `json:"seq"`
	}{token, session, &sequence})
}

// AllowedToStartPulsating you must notify when you are done pulsating!
func (m *Manager) AllowedToStartPulsating(serviceID uint8) bool {
	m.pulseMutex.Lock()
	defer m.pulseMutex.Unlock()

	if m.pulsating == 0 {
		m.pulsating = serviceID
	}

	return m.pulsating == serviceID
}

// StopPulsating stops sending heartbeats to Discord
func (m *Manager) StopPulsating(serviceID uint8) {
	m.pulseMutex.Lock()
	defer m.pulseMutex.Unlock()

	if m.pulsating == serviceID {
		m.pulsating = 0
	}
}

func (m *Manager) pulsate() {
	serviceID := uint8(rand.Intn(254) + 1) // uint8 cap
	if !m.AllowedToStartPulsating(serviceID) {
		return
	}
	defer m.StopPulsating(serviceID)

	m.RLock()
	ticker := time.NewTicker(time.Millisecond * time.Duration(m.heartbeatInterval))
	m.RUnlock()
	defer ticker.Stop()

	var last time.Time
	var snr uint
	for {
		m.RLock()
		last = m.lastHeartbeatAck
		snr = m.sequenceNumber
		m.RUnlock()

		m.Emit(event.Heartbeat, snr)

		// verify the heartbeat ACK
		go func(m *Manager, last time.Time, sent time.Time) {
			<-time.After(3 * time.Second) // deadline for Discord to respond
			m.RLock()
			receivedHeartbeatAck := m.lastHeartbeatAck.After(last)
			m.RUnlock()

			if !receivedHeartbeatAck {
				logrus.Debug("heartbeat ACK was not received")
				m.reconnect()
			} else {
				// update "latency"
				m.heartbeatLatency = m.lastHeartbeatAck.Sub(sent)
			}
		}(m, last, time.Now())

		var shutdown bool
		select {
		case <-ticker.C:
		case <-m.shutdown:
			shutdown = true
		case <-m.restart:
			shutdown = true
		}
		if !shutdown {
			continue
		}

		logrus.Debug("Stopping pulse")
		return
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
