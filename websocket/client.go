package websocket

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket/cmd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"github.com/sirupsen/logrus"
)

const (
	// Deprecated
	maxReconnectTries = 5
)

// NewManager creates a new socket client manager for handling behavior and Discord events. Note that this
// function initiates a go routine.
func NewClient(config *Config, shardID uint) (client *Client, err error) {
	ws, err := newConn(config.HTTPClient)
	if err != nil {
		return nil, err
	}

	if config.TrackedEvents == nil {
		config.TrackedEvents = &UniqueStringSlice{}
	}

	var eChan chan<- *Event
	if config.EventChan != nil {
		eChan = config.EventChan
	} else {
		panic("missing event channel")
	}

	client = &Client{
		conf:              config,
		ShardID:           shardID,
		trackedEvents:     config.TrackedEvents,
		shutdown:          make(chan interface{}),
		restart:           make(chan interface{}),
		eventChan:         eChan,
		receiveChan:       make(chan *discordPacket),
		emitChan:          make(chan *clientPacket),
		conn:              ws,
		ratelimit:         newRatelimiter(),
		timeoutMultiplier: 1,
		disconnected:      true,
	}
	client.Start()

	return
}

func NewTestClient(config *Config, shardID uint, conn Conn) (*Client, chan interface{}) {
	s := make(chan interface{})

	c := &Client{
		conf:              config,
		ShardID:           shardID,
		trackedEvents:     config.TrackedEvents,
		shutdown:          s,
		restart:           make(chan interface{}),
		eventChan:         make(chan *Event),
		receiveChan:       make(chan *discordPacket),
		emitChan:          make(chan *clientPacket),
		conn:              conn,
		ratelimit:         newRatelimiter(),
		timeoutMultiplier: 1,
		disconnected:      true,
	}
	c.Start()
	go c.receiver()

	return c, s
}

// Event is dispatched by the socket layer after parsing and extracting Discord data from a incoming packet.
// This is the data structure used by Disgord for triggering handlers and channels with an event.
type Event struct {
	Name string
	Data []byte
}

// Config ws
// TODO: remove shardID, such that this struct can be reused for every shard
type Config struct {
	// BotToken Discord bot token
	BotToken string

	// HTTPClient custom http client to support the use of proxy
	HTTPClient *http.Client

	// ChannelBuffer is used to set the event channel buffer
	ChannelBuffer uint

	// Endpoint for establishing socket connection. Either endpoints, `Gateway` or `Gateway Bot`, is used to retrieve
	// a valid socket endpoint from Discord
	Endpoint string

	// Encoding make sure we support the correct encoding
	Encoding string

	// Version make sure we support the correct Discord version
	Version int

	// TrackedEvents holds a list of predetermined events that should not be ignored.
	// This is especially useful for creating multiple shards, to reuse the same slice
	TrackedEvents *UniqueStringSlice

	// EventChan can be used to inject a channel instead of letting the ws client construct one
	// useful in sharding to avoid complicated patterns to handle N channels.
	EventChan chan<- *Event

	// for identify packets
	Browser             string
	Device              string
	GuildLargeThreshold uint
	ShardCount          uint
}

type Client struct {
	sync.RWMutex
	conf         *Config
	shutdown     chan interface{}
	restart      chan interface{}
	lastRestart  int64 //unix
	restartMutex sync.Mutex

	ShardID uint

	eventChan     chan<- *Event
	trackedEvents *UniqueStringSlice

	heartbeatInterval uint
	heartbeatLatency  time.Duration
	lastHeartbeatAck  time.Time

	sessionID      string
	trace          []string
	sequenceNumber uint

	ratelimit ratelimiter

	pulsating  uint8
	pulseMutex sync.Mutex

	receiveChan       chan *discordPacket
	emitChan          chan *clientPacket
	conn              Conn
	disconnected      bool
	haveConnectedOnce bool

	isReconnecting bool

	// is the go routine started
	isReceiving  bool
	isEmitting   bool
	recEmitMutex sync.Mutex

	// identify timeout on invalid session
	timeoutMultiplier int
}

// Connect establishes a socket connection with the Discord API
func (m *Client) Connect() (err error) {
	m.Lock()
	defer m.Unlock()

	// m.conn.Disconnected can always tell us if we are disconnected, but it cannot with
	// certainty say if we are connected
	if !m.disconnected {
		err = errors.New("cannot connect while a connection already exist")
		return
	}

	if m.conf.Endpoint == "" {
		panic("missing websocket endpoint. Must be set before constructing the sockets")
		//m.conf.Endpoint, err = getGatewayRoute(m.conf.HTTPClient, m.conf.Version)
		//if err != nil {
		//	return
		//}
	}

	// establish ws connection
	err = m.conn.Open(m.conf.Endpoint, nil)
	if err != nil {
		if !m.conn.Disconnected() {
			m.conn.Close()
		}
		return
	}

	// we can now interact with Discord
	m.haveConnectedOnce = true
	m.disconnected = false
	go m.receiver()
	go m.emitter()
	return
}

// Disconnect disconnects the socket connection
func (m *Client) Disconnect() (err error) {
	m.Lock()
	defer m.Unlock()
	if m.conn.Disconnected() || !m.haveConnectedOnce {
		m.disconnected = true
		err = errors.New("already disconnected")
		return
	}

	// use the emitter to dispatch the close message
	m.Emit(event.Close, nil)
	m.disconnected = true

	// close connection
	<-time.After(time.Second * 1 * time.Duration(m.timeoutMultiplier))

	// wait for processes
	<-time.After(time.Millisecond * 10)
	return
}

func (m *Client) lockReconnect() bool {
	m.restartMutex.Lock()
	defer m.restartMutex.Unlock()

	now := time.Now().UnixNano()
	locked := (now - m.lastRestart) < (time.Second.Nanoseconds() / 2)

	if !locked && !m.isReconnecting {
		m.lastRestart = now
		m.isReconnecting = true
		return true
	}

	return false
}

func (m *Client) unlockReconnect() {
	m.restartMutex.Lock()
	defer m.restartMutex.Unlock()

	m.isReconnecting = false
}

func (m *Client) reconnect() (err error) {
	// make sure there aren't multiple reconnect processes running
	if !m.lockReconnect() {
		return
	}
	defer m.unlockReconnect()

	m.restart <- 1
	_ = m.Disconnect()

	var try uint
	var delay time.Duration = 3 // seconds
	for {
		logrus.Debugf("Reconnect attempt #%d\n", try)
		err = m.Connect()
		if err == nil {
			logrus.Info("successfully reconnected")
			break
		}

		// wait N seconds
		logrus.Infof("reconnect failed, trying again in N seconds; N =  %d", uint(delay))
		logrus.Info(err)
		select {
		case <-time.After(delay * time.Second):
			delay += 4 + time.Duration(try*2)
		case <-m.shutdown:
			return
		}

		if uint(delay) > 5*60 {
			delay = 60
		}
	}

	return
}

// Emit emits a command, if supported, and its data to the Discord Socket API
func (m *Client) Emit(command string, data interface{}) (err error) {
	if !m.haveConnectedOnce {
		return errors.New("race condition detected: you must connect to the socket API/Gateway before you can send gateway commands!")
	}

	var op uint
	switch command {
	case event.Shutdown:
		op = opcode.Shutdown
	case event.Close:
		op = opcode.Close
	case event.Heartbeat:
		op = opcode.Heartbeat
	case event.Identify:
		op = opcode.Identify
	case event.Resume:
		op = opcode.Resume
	case cmd.RequestGuildMembers:
		op = opcode.RequestGuildMembers
	case cmd.UpdateVoiceState:
		op = opcode.VoiceStateUpdate
	case cmd.UpdateStatus:
		op = opcode.StatusUpdate
	default:
		err = errors.New("unsupported command: " + command)
		return
	}

	accepted := m.ratelimit.Request(command)
	if !accepted {
		return errors.New("rate limited")
	}

	m.emitChan <- &clientPacket{
		Op:   op,
		Data: data,
	}
	return
}

// Receive returns the channel for receiving Discord packets
func (m *Client) Receive() <-chan *discordPacket {
	return m.receiveChan
}

func (m *Client) lockEmitter() bool {
	m.recEmitMutex.Lock()
	defer m.recEmitMutex.Unlock()

	if !m.isEmitting {
		m.isEmitting = true
		return true
	}

	return false
}

func (m *Client) unlockEmitter() {
	m.recEmitMutex.Lock()
	defer m.recEmitMutex.Unlock()

	m.isEmitting = false
}

// emitter holds the actually dispatching logic for the Emit method. See DefaultClient#Emit.
func (m *Client) emitter() {
	if !m.lockEmitter() {
		logrus.Error("tried to start another websocket emitter go routine")
		return
	}
	defer m.unlockEmitter()

	for {
		var msg *clientPacket
		var open bool

		select {
		case <-m.shutdown:
			// m.connection got closed
		case msg, open = <-m.emitChan:
		}
		if !open || (msg.Data == nil && (msg.Op == opcode.Shutdown || msg.Op == opcode.Close)) {
			// TODO: what if we get a connection error, how do we restart?
			m.conn.Close()
			return
		}

		err := m.conn.WriteJSON(msg)
		if err != nil {
			// TODO-logging
			fmt.Printf("could not send data to discord: %+v\n", msg)
		}
	}
}

func (m *Client) lockReceiver() bool {
	m.recEmitMutex.Lock()
	defer m.recEmitMutex.Unlock()

	if !m.isReceiving {
		m.isReceiving = true
		return true
	}

	return false
}

func (m *Client) unlockReceiver() {
	m.recEmitMutex.Lock()
	defer m.recEmitMutex.Unlock()

	m.isReceiving = false
}

func (m *Client) receiver() {
	if !m.lockReceiver() {
		logrus.Error("tried to start another websocket receiver go routine")
		return
	}
	defer m.unlockReceiver()

	for {
		packet, err := m.conn.Read()
		if err != nil {
			logrus.Debug("closing readPump")
			return
		}

		// parse to gateway payload object
		evt := &discordPacket{}
		err = evt.UnmarshalJSON(packet)
		if err != nil {
			logrus.Error(err)
			continue
		}

		// save to file
		// build tag: disgord_diagnosews
		if SaveIncomingPackets {
			evtStr := "_" + evt.EventName
			if evtStr == "_" {
				evtStr = ""
			}
			filename := strconv.FormatUint(uint64(evt.SequenceNumber), 10) +
				"_" + strconv.FormatUint(uint64(evt.Op), 10) + evtStr + ".json"
			err = ioutil.WriteFile(DiagnosePath_packets+"/"+filename, packet, 0644)
			if err != nil {
				logrus.WithField("ws-client", "writing packet to file").Error(err)
			}
		}

		// notify listeners
		m.receiveChan <- evt

		// check if application has closed
		select {
		case <-m.shutdown:
			return
		default:
		}
	}
}

// HeartbeatLatency get the time diff between sending a heartbeat and Discord replying with a heartbeat ack
func (m *Client) HeartbeatLatency() (duration time.Duration, err error) {
	duration = m.heartbeatLatency
	if duration == 0 {
		err = errors.New("latency not determined yet")
	}

	return
}

// RegisterEvent tells the socket layer which event types are of interest. Any event that are not registered
// will be discarded once the socket info is extracted from the event.
func (m *Client) RegisterEvent(event string) {
	m.trackedEvents.Add(event)
}

// RemoveEvent removes an event type from the registry. This will cause the event type to be discarded
// by the socket layer.
func (m *Client) RemoveEvent(event string) {
	m.trackedEvents.Remove(event)
}

func (m *Client) Start() {
	go m.operationHandlers()
}

func (m *Client) Shutdown() (err error) {
	m.Disconnect()
	close(m.shutdown)
	return
}

func (m *Client) eventHandler(p *discordPacket) {
	// discord events
	// events that directly correlates to the socket layer, will be dealt with here. But still dispatched.

	// increment the sequence number for each event to make sure everything is synced with discord
	m.Lock()
	m.sequenceNumber++

	// validate the sequence numbers
	if p.SequenceNumber != m.sequenceNumber {
		logrus.Infof("websocket sequence numbers missmatch, forcing reconnect. Got %d, wants %d", p.SequenceNumber, m.sequenceNumber)
		m.sequenceNumber--
		m.Unlock()
		go m.reconnect()
		return
	}
	m.Unlock()

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

func (m *Client) eventOfInterest(name string) bool {
	return m.trackedEvents.Exists(name)
}

// operation handler demultiplexer
func (m *Client) operationHandlers() {
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
			logrus.Info("Discord requested a reconnect")
			go m.reconnect()
		case opcode.InvalidSession:
			// invalid session. Must respond with a identify packet
			logrus.Info("Discord invalidated session")
			go func() {
				// session is invalidated, reset the sequence number
				m.sequenceNumber = 0

				rand.Seed(time.Now().UnixNano())
				delay := rand.Intn(4) + 1
				delay *= m.timeoutMultiplier
				randomDelay := time.Second * time.Duration(delay)
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

func (m *Client) sendHelloPacket() {
	// TODO, this might create several idle goroutines..
	go m.pulsate()

	// if this is a new connection we can drop the resume packet
	if m.sessionID == "" && m.sequenceNumber == 0 {
		err := sendIdentityPacket(m)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	m.RLock()
	token := m.conf.BotToken
	session := m.sessionID
	sequence := m.sequenceNumber
	m.RUnlock()

	m.Emit(event.Resume, struct {
		Token      string `json:"token"`
		SessionID  string `json:"session_id"`
		SequenceNr uint   `json:"seq"`
	}{token, session, sequence})
}

// AllowedToStartPulsating you must notify when you are done pulsating!
func (m *Client) AllowedToStartPulsating(serviceID uint8) bool {
	m.pulseMutex.Lock()
	defer m.pulseMutex.Unlock()

	if m.pulsating == 0 {
		m.pulsating = serviceID
	}

	return m.pulsating == serviceID
}

// StopPulsating stops sending heartbeats to Discord
func (m *Client) StopPulsating(serviceID uint8) {
	m.pulseMutex.Lock()
	defer m.pulseMutex.Unlock()

	if m.pulsating == serviceID {
		m.pulsating = 0
	}
}

func (m *Client) pulsate() {
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

		stopChan := make(chan interface{})

		// verify the heartbeat ACK
		go func(m *Client, last time.Time, sent time.Time, cancel chan interface{}) {
			select {
			case <-cancel:
				return
			case <-time.After(3 * time.Second): // deadline for Discord to respond
			}

			m.RLock()
			receivedHeartbeatAck := m.lastHeartbeatAck.After(last)
			m.RUnlock()

			if !receivedHeartbeatAck {
				logrus.Info("heartbeat ACK was not received, forcing reconnect")
				m.reconnect()
			} else {
				// update "latency"
				m.heartbeatLatency = m.lastHeartbeatAck.Sub(sent)
			}
		}(m, last, time.Now(), stopChan)

		select {
		case <-ticker.C:
			continue
		case <-m.shutdown:
		case <-m.restart:
		}

		logrus.Debug("Stopping pulse")
		close(stopChan)
		return
	}
}

func sendIdentityPacket(m *Client) (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	identityPayload := struct {
		Token          string      `json:"token"`
		Properties     interface{} `json:"properties"`
		Compress       bool        `json:"compress"`
		LargeThreshold uint        `json:"large_threshold"`
		Shard          *[2]uint    `json:"shard,omitempty"`
		Presence       interface{} `json:"presence,omitempty"`
	}{
		Token: m.conf.BotToken,
		Properties: struct {
			OS      string `json:"$os"`
			Browser string `json:"$browser"`
			Device  string `json:"$device"`
		}{runtime.GOOS, m.conf.Browser, m.conf.Device},
		LargeThreshold: m.conf.GuildLargeThreshold,
		Shard:          &[2]uint{m.ShardID, m.conf.ShardCount},
		// Presence: struct {
		// 	Since  *uint       `json:"since"`
		// 	Game   interface{} `json:"game"`
		// 	Status string      `json:"status"`
		// 	AFK    bool        `json:"afk"`
		// }{Status: "online"},
	}

	err = m.Emit(event.Identify, &identityPayload)
	return
}
