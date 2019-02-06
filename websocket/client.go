package websocket

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

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

// newClient ...
func newClient(conf *config, shardID uint) (c *client, err error) {
	var ws Conn
	if conf.conn == nil {
		ws, err = newConn(conf.Proxy)
		if err != nil {
			return nil, err
		}
	} else {
		ws = conf.conn
	}

	c = &client{
		conf:              conf,
		shutdown:          make(chan interface{}),
		restart:           make(chan interface{}),
		ShardID:           shardID,
		receiveChan:       make(chan *DiscordPacket),
		emitChan:          make(chan *clientPacket),
		conn:              ws,
		ratelimit:         newRatelimiter(),
		timeoutMultiplier: 1,
		disconnected:      true,
		log:               conf.Logger,
		behaviors:         map[string]*behavior{},
		poolDiscordPkt:    conf.DiscordPktPool,
	}
	c.connectPermit = &emptyConnectPermit{}
	c.preCon = func() {}
	c.postCon = func() {}

	return
}

type config struct {
	Proxy proxy.Dialer

	// for testing only
	conn Conn

	// Endpoint for establishing socket connection. Either endpoints, `Gateway` or `Gateway Bot`, is used to retrieve
	// a valid socket endpoint from Discord
	Endpoint string

	DiscordPktPool *sync.Pool

	Logger logger.Logger
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
	shutdown     chan interface{}
	restart      chan interface{}
	lastRestart  int64 //unix
	restartMutex sync.Mutex

	ShardID uint

	// sending and receiving data
	ratelimit   ratelimiter
	receiveChan chan *DiscordPacket
	emitChan    chan *clientPacket
	conn        Conn

	// states
	disconnected      bool
	haveConnectedOnce bool
	isReconnecting    bool
	isReceiving       bool // has the go routine started
	isEmitting        bool // has the go routine started
	recEmitMutex      sync.Mutex

	isRestarting bool

	preCon  func()
	postCon func()

	// identify timeout on invalid session
	// useful in unit tests when you want to drop any actual timeouts
	timeoutMultiplier int

	// Proxy allows for use of a custom proxy
	Proxy proxy.Dialer

	// ChannelBuffer is used to set the event channel buffer
	ChannelBuffer uint

	log logger.Logger

	// choreographic programming to handle rate limit, reconnects, and connects
	connectPermit connectPermit

	// behaviours - optional
	behaviors map[string]*behavior

	poolDiscordPkt *sync.Pool
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
)

func (c *client) start() {
	for k := range c.behaviors {
		switch k {
		case discordOperations:
			go c.operationHandlers()
		}
	}
}

// operation handler de-multiplexer
func (c *client) operationHandlers() {
	c.Debug("Ready to receive operation codes...")
	for {
		var p *DiscordPacket
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

		if action, defined := c.behaviors[discordOperations].actions[p.Op]; defined {
			if err := action(p); err != nil {
				c.Error(err)
			}
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
	if c.clientType == clientTypeVoice {
		return "[ws, voice] "
	}

	// [ws, event, shard:0]
	return "" +
		"[ws, " +
		"event, " +
		"shard:" +
		strconv.FormatUint(uint64(c.ShardID), 10) +
		"] "
}

func (c *client) Info(v ...interface{}) {
	if c.log != nil {
		c.log.Info(c.getLogPrefix(), v)
	}
}
func (c *client) Debug(v ...interface{}) {
	if c.log != nil {
		c.log.Debug(c.getLogPrefix(), v)
	}
}
func (c *client) Error(v ...interface{}) {
	if c.log != nil {
		c.log.Error(c.getLogPrefix(), v)
	}
}

var _ logger.Logger = (*client)(nil)

//////////////////////////////////////////////////////
//
// LINKING: CONNECTING / DISCONNECTING / RECONNECTING
//
//////////////////////////////////////////////////////

func (c *client) preConnect(cb func()) {
	c.preCon = cb
}
func (c *client) postConnect(cb func()) {
	c.postCon = cb
}

func (c *client) connect() (err error) {
	// c.conn.Disconnected can always tell us if we are disconnected, but it cannot with
	// certainty say if we are connected
	if !c.disconnected {
		err = errors.New("cannot Connect while a connection already exist")
		return
	}

	if c.conf.Endpoint == "" {
		panic("missing websocket endpoint. Must be set before constructing the sockets")
		//c.conf.Endpoint, err = getGatewayRoute(c.conf.HTTPClient, c.conf.Version)
		//if err != nil {
		//	return
		//}
	}

	if err = c.connectPermit.requestConnectPermit(); err != nil {
		err = errors.New("unable to get permission to Connect. Err: " + err.Error())
		return
	}

	// establish ws connection
	if err = c.conn.Open(c.conf.Endpoint, nil); err != nil {
		if !c.conn.Disconnected() {
			if err2 := c.conn.Close(); err2 != nil {
				c.Error(err2)
			}
		}

		if err3 := c.connectPermit.releaseConnectPermit(); err3 != nil {
			c.Info("unable to release connection permission. Err: ", err3)
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

// Connect establishes a socket connection with the Discord API
func (c *client) Connect() (err error) {
	c.preCon()

	c.Lock()
	if err = c.connect(); err != nil {
		c.Unlock()
		return err
	}
	c.Unlock()

	c.postCon()
	return nil
}

// Disconnect disconnects the socket connection
func (c *client) Disconnect() (err error) {
	c.Lock()
	defer c.Unlock()
	if c.conn.Disconnected() || !c.haveConnectedOnce {
		c.disconnected = true
		err = errors.New("already disconnected")
		return
	}

	// use the emitter to dispatch the close message
	if err = c.conn.Close(); err != nil {
		c.Error(err.Error())
	}

	// c.Emit(event.Close, nil)
	// dont use emit, such that we can call shutdown at the same time as Disconnect (See Shutdown())
	c.disconnected = true

	// close connection
	<-time.After(time.Second * 1 * time.Duration(c.timeoutMultiplier))
	return
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
		return
	}
	defer c.unlockReconnect()

	c.Debug("is reconnecting")

	c.restart <- 1
	_ = c.Disconnect()

	var try uint
	var delay time.Duration = 3 // seconds
	for {
		c.Debug("EvtClient reconnect attempt", try)
		err = c.Connect()
		if err == nil {
			c.Info("successfully reconnected")
			break
		}

		c.Info(fmt.Sprintf("reconnect failed, trying again in N seconds; N =  %d", uint(delay)))
		c.Info(err.Error())

		// wait N seconds
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

var _ Link = (*client)(nil)

//////////////////////////////////////////////////////
//
// EMITTING / DISPATCHING
//
//////////////////////////////////////////////////////

// Emit is used by DisGord users for dispatching a socket command to the Discord Gateway.
func (c *client) Emit(command string, data interface{}) (err error) {
	if !c.haveConnectedOnce {
		return errors.New("race condition detected: you must Connect to the socket API/Gateway before you can send gateway commands")
	}

	op := ^uint(0)
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
	if op == ^uint(0) {
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
func (c *client) emitter() {
	if !c.lockEmitter() {
		c.Debug("tried to start another websocket emitter go routine")
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
			if err := c.Disconnect(); err != nil {
				c.Error(err)
			}
			c.Debug("closing emitter")
			return
		}
		var err error

		// save to file
		// build tag: disgord_diagnosews
		saveOutgoingPacket(c, msg)

		err = c.conn.WriteJSON(msg)
		if err != nil {
			c.Error(err.Error())
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

func (c *client) receiver() {
	if !c.lockReceiver() {
		c.Debug("tried to start another websocket receiver go routine")
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
		// see operationHandler for return/Put()
		evt := c.poolDiscordPkt.Get().(*DiscordPacket)
		evt.reset()
		//err = evt.UnmarshalJSON(packet) // custom unmarshal
		err = httd.Unmarshal(packet, evt) // json.RawMessage
		if err != nil {
			c.Error(err.Error())
			continue
		}

		// save to file
		// build tag: disgord_diagnosews
		saveIncomingPacker(c, evt, packet)

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

func (c *client) Shutdown() (err error) {
	close(c.shutdown)
	c.Disconnect()
	return
}

//////////////////////////////////////////////////////
//
// HEARTBEAT / PULSATING
//
//////////////////////////////////////////////////////

// TODO - is it worth it?
