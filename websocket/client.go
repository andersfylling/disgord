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

	"github.com/andersfylling/disgord/constant"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket/cmd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
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
		log:               config.Logger,
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
		log:               config.Logger,
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

	Logger constant.Logger

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
	ReadyCounter uint

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

	log constant.Logger
}

func (c *Client) Info(msg string) {
	if c.log != nil {
		c.log.Info("[ws-client, shard:" + strconv.FormatUint(uint64(c.ShardID), 10) + "] " + msg)
	}
}
func (c *Client) Debug(msg string) {
	if c.log != nil {
		c.log.Debug("[ws-client, shard:" + strconv.FormatUint(uint64(c.ShardID), 10) + "] " + msg)
	}
}
func (c *Client) Error(msg string) {
	if c.log != nil {
		c.log.Error("[ws-client, shard:" + strconv.FormatUint(uint64(c.ShardID), 10) + "] " + msg)
	}
}

var _ constant.Logger = (*Client)(nil)

// Connect establishes a socket connection with the Discord API
func (c *Client) Connect() (err error) {
	c.Lock()
	defer c.Unlock()

	// c.conn.Disconnected can always tell us if we are disconnected, but it cannot with
	// certainty say if we are connected
	if !c.disconnected {
		err = errors.New("cannot connect while a connection already exist")
		return
	}

	if c.conf.Endpoint == "" {
		panic("missing websocket endpoint. Must be set before constructing the sockets")
		//c.conf.Endpoint, err = getGatewayRoute(c.conf.HTTPClient, c.conf.Version)
		//if err != nil {
		//	return
		//}
	}

	// establish ws connection
	err = c.conn.Open(c.conf.Endpoint, nil)
	if err != nil {
		if !c.conn.Disconnected() {
			c.conn.Close()
		}
		return
	}

	// we can now interact with Discord
	c.haveConnectedOnce = true
	c.disconnected = false
	go c.receiver()
	go c.emitter()
	return
}

// Disconnect disconnects the socket connection
func (c *Client) Disconnect() (err error) {
	c.Lock()
	defer c.Unlock()
	if c.conn.Disconnected() || !c.haveConnectedOnce {
		c.disconnected = true
		err = errors.New("already disconnected")
		return
	}

	// use the emitter to dispatch the close message
	c.Emit(event.Close, nil)
	c.disconnected = true

	// close connection
	<-time.After(time.Second * 1 * time.Duration(c.timeoutMultiplier))

	// wait for processes
	<-time.After(time.Millisecond * 10)
	return
}

func (c *Client) lockReconnect() bool {
	c.restartMutex.Lock()
	defer c.restartMutex.Unlock()

	now := time.Now().UnixNano()
	locked := (now - c.lastRestart) < (time.Second.Nanoseconds() / 2)

	if !locked && !c.isReconnecting {
		c.lastRestart = now
		c.isReconnecting = true
		return true
	}

	return false
}

func (c *Client) unlockReconnect() {
	c.restartMutex.Lock()
	defer c.restartMutex.Unlock()

	c.isReconnecting = false
}

func (c *Client) reconnect() (err error) {
	// make sure there aren't multiple reconnect processes running
	if !c.lockReconnect() {
		return
	}
	defer c.unlockReconnect()

	c.restart <- 1
	_ = c.Disconnect()

	var try uint
	var delay time.Duration = 3 // seconds
	for {
		c.Debug(fmt.Sprintf("Reconnect attempt #%d\n", try))
		err = c.Connect()
		if err == nil {
			c.Info("successfully reconnected")
			break
		}

		// wait N seconds
		c.Info(fmt.Sprintf("reconnect failed, trying again in N seconds; N =  %d", uint(delay)))
		c.Info(err.Error())
		select {
		case <-time.After(delay * time.Second):
			delay += 4 + time.Duration(try*2)
		case <-c.shutdown:
			return
		}

		if uint(delay) > 5*60 {
			delay = 60
		}
	}

	return
}

// Emit emits a command, if supported, and its data to the Discord Socket API
func (c *Client) Emit(command string, data interface{}) (err error) {
	if !c.haveConnectedOnce {
		return errors.New("race condition detected: you must connect to the socket API/Gateway before you can send gateway commands")
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

	accepted := c.ratelimit.Request(command)
	if !accepted {
		return errors.New("rate limited")
	}

	c.emitChan <- &clientPacket{
		Op:   op,
		Data: data,
	}
	return
}

// Receive returns the channel for receiving Discord packets
func (c *Client) Receive() <-chan *discordPacket {
	return c.receiveChan
}

func (c *Client) lockEmitter() bool {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	if !c.isEmitting {
		c.isEmitting = true
		return true
	}

	return false
}

func (c *Client) unlockEmitter() {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	c.isEmitting = false
}

// emitter holds the actually dispatching logic for the Emit method. See DefaultClient#Emit.
func (c *Client) emitter() {
	if !c.lockEmitter() {
		c.Error("tried to start another websocket emitter go routine")
		return
	}
	defer c.unlockEmitter()

	for {
		var msg *clientPacket
		var open bool

		select {
		case <-c.shutdown:
			// c.connection got closed
		case msg, open = <-c.emitChan:
		}
		if !open || (msg.Data == nil && (msg.Op == opcode.Shutdown || msg.Op == opcode.Close)) {
			// TODO: what if we get a connection error, how do we restart?
			err := c.conn.Close()
			if err != nil {
				c.Error(err.Error())
			}
			return
		}

		err := c.conn.WriteJSON(msg)
		if err != nil {
			c.Error(err.Error())
		}
	}
}

func (c *Client) lockReceiver() bool {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	if !c.isReceiving {
		c.isReceiving = true
		return true
	}

	return false
}

func (c *Client) unlockReceiver() {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	c.isReceiving = false
}

func (c *Client) receiver() {
	if !c.lockReceiver() {
		c.Error("tried to start another websocket receiver go routine")
		return
	}
	defer c.unlockReceiver()

	for {
		packet, err := c.conn.Read()
		if err != nil {
			c.Debug("closing receiver")
			return
		}

		// parse to gateway payload object
		evt := &discordPacket{}
		err = evt.UnmarshalJSON(packet)
		if err != nil {
			c.Error(err.Error())
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
				c.Error(err.Error())
			}
		}

		// notify listeners
		c.receiveChan <- evt

		// check if application has closed
		select {
		case <-c.shutdown:
			return
		default:
		}
	}
}

// HeartbeatLatency get the time diff between sending a heartbeat and Discord replying with a heartbeat ack
func (c *Client) HeartbeatLatency() (duration time.Duration, err error) {
	duration = c.heartbeatLatency
	if duration == 0 {
		err = errors.New("latency not determined yet")
	}

	return
}

// RegisterEvent tells the socket layer which event types are of interest. Any event that are not registered
// will be discarded once the socket info is extracted from the event.
func (c *Client) RegisterEvent(event string) {
	c.trackedEvents.Add(event)
}

// RemoveEvent removes an event type from the registry. This will cause the event type to be discarded
// by the socket layer.
func (c *Client) RemoveEvent(event string) {
	c.trackedEvents.Remove(event)
}

func (c *Client) Start() {
	go c.operationHandlers()
}

func (c *Client) Shutdown() (err error) {
	c.Disconnect()
	close(c.shutdown)
	return
}

func (c *Client) eventHandler(p *discordPacket) {
	// discord events
	// events that directly correlates to the socket layer, will be dealt with here. But still dispatched.

	// increment the sequence number for each event to make sure everything is synced with discord
	c.Lock()
	c.sequenceNumber++

	// validate the sequence numbers
	if p.SequenceNumber != c.sequenceNumber {
		c.Info(fmt.Sprintf("websocket sequence numbers missmatch, forcing reconnect. Got %d, wants %d", p.SequenceNumber, c.sequenceNumber))
		c.sequenceNumber--
		c.Unlock()
		go c.reconnect()
		return
	}
	c.Unlock()

	if p.EventName == event.Ready {

		// always store the session id & update the trace content
		ready := readyPacket{}
		err := httd.Unmarshal(p.Data, &ready)
		if err != nil {
			c.Error(err.Error())
		}

		c.Lock()
		c.sessionID = ready.SessionID
		c.trace = ready.Trace
		c.ReadyCounter++
		c.Unlock()
	} else if p.EventName == event.Resume {
		// eh? debugging.
		// TODO
	} else if p.Op == opcode.DiscordEvent && !c.eventOfInterest(p.EventName) {
		return
	}

	// dispatch event
	c.eventChan <- &Event{
		Name: p.EventName,
		Data: p.Data,
	}
} // end eventHandler()

func (c *Client) eventOfInterest(name string) bool {
	return c.trackedEvents.Exists(name)
}

// operation handler demultiplexer
func (c *Client) operationHandlers() {
	c.Debug("Ready to receive operation codes...")
	for {
		var p *discordPacket
		var open bool
		select {
		case p, open = <-c.Receive():
			if !open {
				c.Debug("operationChan is dead..")
				return
			}
		// case <-c.restart:
		case <-c.shutdown:
			c.Debug("exiting operation handler")
			return
		}

		// new packet that must be handled by it's Discord operation code
		switch p.Op {
		case opcode.DiscordEvent:
			c.eventHandler(p)
		case opcode.Reconnect:
			c.Info("Discord requested a reconnect")
			go c.reconnect()
		case opcode.InvalidSession:
			// invalid session. Must respond with a identify packet
			c.Info("Discord invalidated session")
			go func() {
				// session is invalidated, reset the sequence number
				c.sequenceNumber = 0

				rand.Seed(time.Now().UnixNano())
				delay := rand.Intn(4) + 1
				delay *= c.timeoutMultiplier
				randomDelay := time.Second * time.Duration(delay)
				<-time.After(randomDelay)
				err := sendIdentityPacket(c)
				if err != nil {
					c.Error(err.Error())
				}
			}()
		case opcode.Heartbeat:
			// https://discordapp.com/developers/docs/topics/gateway#heartbeating
			_ = c.Emit(event.Heartbeat, c.sequenceNumber)
		case opcode.Hello:
			// hello
			helloPk := &helloPacket{}
			err := httd.Unmarshal(p.Data, helloPk)
			if err != nil {
				c.Debug(err.Error())
			}
			c.Lock()
			c.heartbeatInterval = helloPk.HeartbeatInterval
			c.Unlock()

			c.sendHelloPacket()
		case opcode.HeartbeatAck:
			// heartbeat received
			c.Lock()
			c.lastHeartbeatAck = time.Now()
			c.Unlock()
		default:
			// unknown
			c.Debug(fmt.Sprintf("Unknown operation: %+v\n", p))
		}
	}
}

func (c *Client) sendHelloPacket() {
	// TODO, this might create several idle goroutines..
	go c.pulsate()

	// if this is a new connection we can drop the resume packet
	if c.sessionID == "" && c.sequenceNumber == 0 {
		err := sendIdentityPacket(c)
		if err != nil {
			c.Error(err.Error())
		}
		return
	}

	c.RLock()
	token := c.conf.BotToken
	session := c.sessionID
	sequence := c.sequenceNumber
	c.RUnlock()

	err := c.Emit(event.Resume, struct {
		Token      string `json:"token"`
		SessionID  string `json:"session_id"`
		SequenceNr uint   `json:"seq"`
	}{token, session, sequence})
	if err != nil {
		c.Error(err.Error())
	}
}

// AllowedToStartPulsating you must notify when you are done pulsating!
func (c *Client) AllowedToStartPulsating(serviceID uint8) bool {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == 0 {
		c.pulsating = serviceID
	}

	return c.pulsating == serviceID
}

// StopPulsating stops sending heartbeats to Discord
func (c *Client) StopPulsating(serviceID uint8) {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == serviceID {
		c.pulsating = 0
	}
}

func (c *Client) pulsate() {
	serviceID := uint8(rand.Intn(254) + 1) // uint8 cap
	if !c.AllowedToStartPulsating(serviceID) {
		return
	}
	defer c.StopPulsating(serviceID)

	c.RLock()
	ticker := time.NewTicker(time.Millisecond * time.Duration(c.heartbeatInterval))
	c.RUnlock()
	defer ticker.Stop()

	var last time.Time
	var snr uint
	for {
		c.RLock()
		last = c.lastHeartbeatAck
		snr = c.sequenceNumber
		c.RUnlock()

		c.Emit(event.Heartbeat, snr)

		stopChan := make(chan interface{})

		// verify the heartbeat ACK
		go func(m *Client, last time.Time, sent time.Time, cancel chan interface{}) {
			select {
			case <-cancel:
				return
			case <-time.After(3 * time.Second): // deadline for Discord to respond
			}

			c.RLock()
			receivedHeartbeatAck := c.lastHeartbeatAck.After(last)
			c.RUnlock()

			if !receivedHeartbeatAck {
				c.Info("heartbeat ACK was not received, forcing reconnect")
				err := c.reconnect()
				if err != nil {
					c.Error(err.Error())
				}
			} else {
				// update "latency"
				c.heartbeatLatency = c.lastHeartbeatAck.Sub(sent)
			}
		}(c, last, time.Now(), stopChan)

		select {
		case <-ticker.C:
			continue
		case <-c.shutdown:
		case <-c.restart:
		}

		c.Debug("Stopping pulse")
		close(stopChan)
		return
	}
}

func sendIdentityPacket(c *Client) (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	identityPayload := struct {
		Token          string      `json:"token"`
		Properties     interface{} `json:"properties"`
		Compress       bool        `json:"compress"`
		LargeThreshold uint        `json:"large_threshold"`
		Shard          *[2]uint    `json:"shard,omitempty"`
		Presence       interface{} `json:"presence,omitempty"`
	}{
		Token: c.conf.BotToken,
		Properties: struct {
			OS      string `json:"$os"`
			Browser string `json:"$browser"`
			Device  string `json:"$device"`
		}{runtime.GOOS, c.conf.Browser, c.conf.Device},
		LargeThreshold: c.conf.GuildLargeThreshold,
		Shard:          &[2]uint{c.ShardID, c.conf.ShardCount},
		// Presence: struct {
		// 	Since  *uint       `json:"since"`
		// 	Game   interface{} `json:"game"`
		// 	Status string      `json:"status"`
		// 	AFK    bool        `json:"afk"`
		// }{Status: "online"},
	}

	err = c.Emit(event.Identify, &identityPayload)
	return
}
