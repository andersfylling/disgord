package websocket

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/andersfylling/disgord/httd"

	"github.com/andersfylling/disgord/logger"

	"github.com/andersfylling/disgord/websocket/cmd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"golang.org/x/net/proxy"
)

const (
	clientTypeEvent = iota
	clientTypeVoice
)

// Link is used to establish basic commands to create and destroy a link.
// See client.Disconnect() and client.Connect() for linking to the Discord servers
type Link interface {
	Connect() error
	Disconnect() error
}

//////////////////////////////////////////////////////
//
// synchronization & rate limiting
// By default, no such restrictions exist
//
//////////////////////////////////////////////////////

type connectPermit interface {
	requestConnectPermit() error
	releaseConnectPermit() error
}

type emptyConnectPermit struct {
}

func (emptyConnectPermit) requestConnectPermit() error {
	return nil
}

func (emptyConnectPermit) releaseConnectPermit() error {
	return nil
}

var _ connectPermit = (*emptyConnectPermit)(nil)

type connectSignature = func() (evt interface{}, err error)

// newClient ...
func newClient(shardID uint, conf *config, connect connectSignature) (c *client, err error) {
	var ws Conn
	if conf.conn == nil {
		ws, err = newConn(conf.Proxy, conf.HTTPClient)
		if err != nil {
			return nil, err
		}
	} else {
		ws = conf.conn
	}

	c = &client{
		conf:              conf,
		ShardID:           shardID,
		receiveChan:       make(chan *DiscordPacket, 50),
		internalEmitChan:  make(chan *clientPacket, 50),
		emitChan:          make(chan *clientPacket, 50),
		conn:              ws,
		ratelimit:         newRatelimiter(),
		timeoutMultiplier: 1,
		disconnected:      true,
		log:               conf.Logger,
		behaviors:         map[string]*behavior{},
		poolDiscordPkt:    conf.DiscordPktPool,
		onceChannels:      newOnceChannels(),
		connect:           connect,

		activateHeartbeats: make(chan interface{}),
		SystemShutdown:     conf.SystemShutdown,
	}
	c.emitChanMutex.Lock()

	return
}

type config struct {
	Proxy      proxy.Dialer
	HTTPClient *http.Client

	// for testing only
	conn Conn

	// Endpoint for establishing socket connection. Either endpoints, `Gateway` or `Gateway Bot`, is used to retrieve
	// a valid socket endpoint from Discord
	Endpoint string

	DiscordPktPool *sync.Pool

	Logger logger.Logger

	SystemShutdown chan interface{}
}

// client can be used as a base for other ws clients; voice, event. Note the use of
// client.ReleasePermit() and client.RequestPermit() in Connect (and then automatically reconnect()).
// these are used for synchronizing connecting and you must therefore correctly release the permit once you have
// such that the next shard or yourself can Connect in the future.
//
// If you do not care about these. Please overwrite both methods.
type client struct {
	sync.RWMutex
	clientType   int
	conf         *config
	lastRestart  int64      //unix
	restartMutex sync.Mutex // TODO: atomic bool

	pulsating          uint8
	pulseMutex         sync.Mutex
	heartbeatLatency   time.Duration
	heartbeatInterval  uint
	lastHeartbeatAck   time.Time
	lastHeartbeatSent  time.Time
	activateHeartbeats chan interface{}

	ShardID uint

	// sending and receiving data
	ratelimit        ratelimiter
	receiveChan      chan *DiscordPacket
	internalEmitChan chan *clientPacket
	emitChan         chan *clientPacket
	emitChanMutex    sync.Mutex
	conn             Conn

	// connect is blocking until a websocket connection has completed it's setup.
	// eg. Normal shards that handles events are considered connected once the
	// identity/resume has been sent. While for voice we wait until a ready event
	// is returned.
	connect connectSignature

	// states
	disconnected      bool
	haveConnectedOnce bool
	isReconnecting    bool
	isReceiving       bool // has the go routine started
	isEmitting        bool // has the go routine started
	recEmitMutex      sync.Mutex
	onceChannels      onceChannels

	isRestarting bool

	// identify timeout on invalid session
	// useful in unit tests when you want to drop any actual timeouts
	timeoutMultiplier int

	// proxy allows for use of a custom proxy
	Proxy proxy.Dialer

	// ChannelBuffer is used to set the event channel buffer
	ChannelBuffer uint

	log         logger.Logger
	logSequence atomic.Uint64

	// behaviours - optional
	behaviors map[string]*behavior

	poolDiscordPkt *sync.Pool

	cancel context.CancelFunc

	SystemShutdown <-chan interface{}

	// receiver gets closed when the connection is lost
	requestedDisconnect bool
}

type behaviorActions map[interface{}]actionFunc
type actionFunc func(interface{}) error
type behavior struct {
	addresses string
	actions   behaviorActions
}

func (c *client) addBehavior(b *behavior) {
	c.behaviors[b.addresses] = b
}

const (
	discordOperations string = "discord-ops"
	heartbeating      string = "heartbeats"
	sendHeartbeat            = 0
)

func (c *client) startBehaviors(ctx context.Context) {
	for k := range c.behaviors {
		switch k {
		case discordOperations:
			go c.operationHandlers(ctx)
		}
	}
}

// operation handler de-multiplexer
func (c *client) operationHandlers(ctx context.Context) {
	c.log.Debug(c.getLogPrefix(), "Ready to receive operation codes...")
	for {
		var p *DiscordPacket
		var open bool
		select {
		case p, open = <-c.Receive():
			if !open {
				c.log.Debug(c.getLogPrefix(), "operationChan is dead..")
				return
			}
		case <-ctx.Done():
			c.log.Debug(c.getLogPrefix(), "closing operations handler")
			return
		}

		if action, defined := c.behaviors[discordOperations].actions[p.Op]; defined {
			if err := action(p); err != nil {
				c.log.Error(c.getLogPrefix(), err)
			}
		} else {
			c.log.Debug(c.getLogPrefix(), "tried calling undefined discord operation", p.Op)
		}

		// see receiver() for creation/Get()
		c.poolDiscordPkt.Put(p)
	}
}

//////////////////////////////////////////////////////
//
// LOGGING
//
//////////////////////////////////////////////////////

func (c *client) getLogPrefix() string {
	t := "ws-"
	if c.clientType == clientTypeVoice {
		t += "v"
	} else if c.clientType == clientTypeEvent {
		t += "e"
	} else {
		t += "?"
	}

	s := "s:" + strconv.FormatUint(c.logSequence.Inc(), 10)
	shardID := "shard:" + strconv.FormatUint(uint64(c.ShardID), 10)

	// [ws-?, s:0, shard:0]
	return "[" + t + "," + s + "," + shardID + "]"
}

//////////////////////////////////////////////////////
//
// LINKING: CONNECTING / DISCONNECTING / RECONNECTING
//
//////////////////////////////////////////////////////
func (c *client) disconnect() (err error) {
	c.Lock()
	defer c.Unlock()
	if c.conn.Disconnected() || !c.haveConnectedOnce || c.cancel == nil {
		c.disconnected = true
		err = errors.New("already disconnected")
		return
	}

	// stop emitter, receiver and behaviors
	c.cancel()
	c.cancel = nil

	// use the emitter to dispatch the close message
	err = c.conn.Close()
	// a typical err here is that the pipe is closed. Err is returned later

	// c.Emit(event.Close, nil)
	// dont use emit, such that we can call shutdown at the same time as Disconnect (See Shutdown())
	c.disconnected = true

	c.log.Info(c.getLogPrefix(), "disconnected")

	// close connection
	<-time.After(time.Second * 1 * time.Duration(c.timeoutMultiplier))

	return
}

// Disconnect disconnects the socket connection
func (c *client) Disconnect() (err error) {
	c.Lock()
	c.requestedDisconnect = true
	c.Unlock()
	return c.disconnect()
}

func (c *client) lockReconnect() bool {
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

func (c *client) unlockReconnect() {
	c.restartMutex.Lock()
	defer c.restartMutex.Unlock()

	c.isReconnecting = false
}

func (c *client) reconnect() (err error) {
	// make sure there aren't multiple reconnect processes running
	if !c.lockReconnect() {
		c.log.Debug(c.getLogPrefix(), "tried to start reconnect when already reconnecting")
		return
	}
	defer c.unlockReconnect()

	c.log.Debug(c.getLogPrefix(), "is reconnecting")
	if err := c.disconnect(); err != nil {
		c.RLock()
		if c.requestedDisconnect {
			c.RUnlock()
			c.log.Debug(c.getLogPrefix(), err)
			return errors.New("already disconnected, cannot reconnect")
		}
		c.RUnlock()
	}

	return c.reconnectLoop()
}

func (c *client) reconnectLoop() (err error) {
	var try uint
	var delay = 3 * time.Second
	for {
		if try == 0 {
			c.log.Debug(c.getLogPrefix(), "trying to connect")
		} else {
			c.log.Debug(c.getLogPrefix(), "reconnect attempt", try)
		}
		if _, err = c.connect(); err == nil {
			c.log.Debug(c.getLogPrefix(), "establishing connection succeeded")
			break
		}

		c.log.Info(c.getLogPrefix(), "establishing connection failed, trying again in ", delay)
		c.log.Info(c.getLogPrefix(), err)

		// wait N seconds
		select {
		case <-time.After(delay):
			delay += (4 + time.Duration(try*2)) * time.Second
		case <-c.SystemShutdown:
			c.log.Debug(c.getLogPrefix(), "stopping reconnect attempt", try)
			return
		}

		if delay > 5*60*time.Second {
			delay = 60 * time.Second
		}
	}

	return
}

//////////////////////////////////////////////////////
//
// EMITTING / DISPATCHING
//
//////////////////////////////////////////////////////

// Emit is used by DisGord users for dispatching a socket command to the Discord Gateway.
func (c *client) Emit(internal bool, command string, data interface{}) (err error) {
	if !c.haveConnectedOnce {
		return errors.New("race condition detected: you must Connect to the socket API/Gateway before you can send gateway commands: " + command)
	}

	noMatch := ^uint(0)
	op := noMatch
	// TODO: refactor command and event name to avoid conversion (?)
	if c.clientType == clientTypeVoice {
		switch command {
		case cmd.VoiceSpeaking:
			op = opcode.VoiceSpeaking
		case cmd.VoiceIdentify:
			op = opcode.VoiceIdentify
		case cmd.VoiceSelectProtocol:
			op = opcode.VoiceSelectProtocol
		case cmd.VoiceHeartbeat:
			op = opcode.VoiceHeartbeat
		case cmd.VoiceResume:
			op = opcode.VoiceResume
		}
	} else if c.clientType == clientTypeEvent {
		switch command {
		case event.Shutdown:
			op = opcode.Shutdown
		case event.Close:
			op = opcode.Close
		case event.Heartbeat:
			op = opcode.EventHeartbeat
		case event.Identify:
			op = opcode.EventIdentify
		case event.Resume:
			op = opcode.EventResume
		case cmd.RequestGuildMembers:
			op = opcode.EventRequestGuildMembers
		case cmd.UpdateVoiceState:
			op = opcode.EventVoiceStateUpdate
		case cmd.UpdateStatus:
			op = opcode.EventStatusUpdate
		}
	}

	if op == noMatch {
		return errors.New("unsupported command: " + command)
	}

	if accepted := c.ratelimit.Request(command); !accepted {
		return errors.New("rate limited")
	}

	// TODO: que messages when disconnected( or suspended)
	p := &clientPacket{
		Op:   op,
		Data: data,
	}
	if internal {
		c.internalEmitChan <- p
	} else {
		c.emitChan <- p
	}
	return
}

func (c *client) lockEmitter() bool {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	if !c.isEmitting {
		c.isEmitting = true
		return true
	}

	return false
}

func (c *client) unlockEmitter() {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	c.isEmitting = false
}

// emitter holds the actually dispatching logic for sending data to the Discord Gateway.
// client#Emit depends on this.
func (c *client) emitter(ctx context.Context) {
	if !c.lockEmitter() {
		c.log.Debug(c.getLogPrefix(), "tried to startBehaviors another websocket emitter go routine")
		return
	}
	defer c.unlockEmitter()
	c.log.Debug(c.getLogPrefix(), "starting emitter")

	internal, cancel := context.WithCancel(context.Background())
	for {
		var msg *clientPacket
		var open bool

		select {
		case <-ctx.Done():
			c.log.Debug(c.getLogPrefix(), "closing emitter")
			return
		case <-internal.Done():
			c.log.Debug(c.getLogPrefix(), "closing emitter after write error")
			return
		case msg, open = <-c.internalEmitChan:
		case msg, open = <-c.emitChan:
		}
		if !open || (msg.Data == nil && (msg.Op == opcode.Shutdown || msg.Op == opcode.Close)) {
			if err := c.Disconnect(); err != nil {
				c.log.Error(c.getLogPrefix(), err)
			}
			c.log.Debug(c.getLogPrefix(), "closing emitter after read")
			return
		}
		var err error

		// save to file
		// build tag: disgord_diagnosews
		saveOutgoingPacket(c, msg)

		if err = c.conn.WriteJSON(msg); err != nil {
			c.log.Error(c.getLogPrefix(), err, fmt.Sprintf("%+v", *msg))
			cancel()
		}
	}
}

//////////////////////////////////////////////////////
//
// RECEIVING
//
//////////////////////////////////////////////////////

// Receive returns the channel for receiving Discord packets
func (c *client) Receive() <-chan *DiscordPacket {
	return c.receiveChan
}

func (c *client) lockReceiver() bool {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	if !c.isReceiving {
		c.isReceiving = true
		return true
	}

	return false
}

func (c *client) unlockReceiver() {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	c.isReceiving = false
}

func (c *client) receiver(ctx context.Context) {
	if !c.lockReceiver() {
		c.log.Debug(c.getLogPrefix(), "tried to start another receiver")
		return
	}
	defer c.unlockReceiver()
	c.log.Debug(c.getLogPrefix(), "starting receiver")

	for {
		var packet []byte
		var err error
		if packet, err = c.conn.Read(); err != nil {
			c.log.Debug(c.getLogPrefix(), "closing receiver", err)
			// TODO: should be able to tag c.conn as disconnected at this stage
			return
		}

		// parse to gateway payload object
		// see operationHandler for return/Put()
		evt := c.poolDiscordPkt.Get().(*DiscordPacket)
		evt.reset()
		//err = evt.UnmarshalJSON(packet) // custom unmarshal
		if err = httd.Unmarshal(packet, evt); err != nil {
			c.log.Error(c.getLogPrefix(), err)
			continue
		}

		// save to file
		// build tag: disgord_diagnosews
		saveIncomingPacker(c, evt, packet)

		// notify listeners
		c.receiveChan <- evt

		// check if application has closed
		select {
		case <-ctx.Done():
			c.log.Debug(c.getLogPrefix(), "closing receiver")
			return
		default:
		}
	}
}

//////////////////////////////////////////////////////
//
// HEARTBEAT / PULSATING
//
//////////////////////////////////////////////////////

// AllowedToStartPulsating you must notify when you are done pulsating!
func (c *client) AllowedToStartPulsating(serviceID uint8) bool {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == 0 {
		c.pulsating = serviceID
	}

	return c.pulsating == serviceID
}

// StopPulsating stops sending heartbeats to Discord
func (c *client) StopPulsating(serviceID uint8) {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == serviceID {
		c.pulsating = 0
	}
}

func (c *client) prepareHeartbeating(ctx context.Context) {
	serviceID := uint8(rand.Intn(254) + 1) // uint8 cap
	if !c.AllowedToStartPulsating(serviceID) {
		c.log.Debug(c.getLogPrefix(), "tried to start an additional pulse")
		return
	}
	defer c.StopPulsating(serviceID)

	select {
	case <-ctx.Done():
		c.log.Debug(c.getLogPrefix(), "heartbeat preparations cancelled")
		return
	case <-c.activateHeartbeats:
	}

	c.pulsate(ctx)
}

func (c *client) pulsate(ctx context.Context) {
	c.RLock()
	c.lastHeartbeatSent = time.Now()
	c.lastHeartbeatAck = time.Now()
	interval := time.Millisecond * time.Duration(c.heartbeatInterval)
	c.RUnlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var lastAck time.Time
	var lastSent time.Time
	for {
		c.RLock()
		lastAck = c.lastHeartbeatAck
		lastSent = c.lastHeartbeatSent
		c.RUnlock()

		// make sure that Discord replied to the last heartbeat signal (heartbeat ack)
		if lastSent.After(lastAck) {
			c.log.Info(c.getLogPrefix(), "heartbeat ACK was not received, forcing reconnect")
			go c.reconnect()
			break
		} else {
			c.log.Debug(c.getLogPrefix(), "heartbeat ACK ok")
		}

		// update heartbeat latency record & send new heartbeat signal
		c.Lock()
		c.heartbeatLatency = lastAck.Sub(lastSent)
		c.lastHeartbeatSent = time.Now()
		c.Unlock()
		if err := c.behaviors[heartbeating].actions[sendHeartbeat](nil); err != nil {
			c.log.Error(c.getLogPrefix(), err)
		} else {
			c.log.Debug(c.getLogPrefix(), "sent heartbeat")
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
		}
		break
	}
	c.log.Debug(c.getLogPrefix(), "stopping pulse")
}

// HeartbeatLatency get the time diff between sending a heartbeat and Discord replying with a heartbeat ack
func (c *client) HeartbeatLatency() (duration time.Duration, err error) {
	c.RLock()
	defer c.RUnlock()

	duration = c.heartbeatLatency
	if duration == 0 {
		err = errors.New("latency not determined yet")
	}

	return
}
