package websocket

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/disgord/constant"

	"github.com/andersfylling/disgord/websocket/cmd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
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

type connectPermit interface {
	requestConnectPermit() error
	releaseConnectPermit() error
}

// newClient ...
func newClient(config *Config, shardID uint) (c *client, err error) {
	ws, err := newConn(config.HTTPClient)
	if err != nil {
		return nil, err
	}

	c = &client{
		conf:              config,
		shutdown:          make(chan interface{}),
		restart:           make(chan interface{}),
		ShardID:           shardID,
		receiveChan:       make(chan *DiscordPacket),
		emitChan:          make(chan *clientPacket),
		conn:              ws,
		ratelimit:         newRatelimiter(),
		timeoutMultiplier: 1,
		disconnected:      true,
		log:               config.Logger,
		a:                 config.A,
		behaviors:         map[string]*behavior{},
		poolDiscordPkt:    config.DiscordPktPool,
	}
	c.connectPermit = c

	return
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
	conf         *Config
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

	// identify timeout on invalid session
	// useful in unit tests when you want to drop any actual timeouts
	timeoutMultiplier int

	log constant.Logger

	// choreographic programming to handle rate limit, reconnects, and connects
	connectPermit connectPermit
	K             *K
	a             A

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
			err := action(p)
			if err != nil {
				c.Error(err.Error())
			}
		}

		// see receiver() for creation/Get()
		c.poolDiscordPkt.Put(p)
	}
}

//////////////////////////////////////////////////////
//
// SHARD synchronization & rate limiting
//
//////////////////////////////////////////////////////

func (c *client) requestConnectPermit() error {
	c.Debug("trying to get Connect permission")
	b := make(B)
	defer close(b)
	c.a <- b
	c.Debug("waiting")
	var ok bool
	select {
	case c.K, ok = <-b:
		if !ok || c.K == nil {
			c.Debug("unable to get Connect permission")
			return errors.New("channel closed or K was nil")
		}
		c.Debug("got Connect permission")
	case <-c.shutdown:
	}

	return nil
}

func (c *client) releaseConnectPermit() error {
	if c.K == nil {
		return errors.New("K has not been granted yet")
	}

	c.K.Release <- c.K
	c.K = nil
	return nil
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

var _ constant.Logger = (*client)(nil)

//////////////////////////////////////////////////////
//
// LINKING: CONNECTING / DISCONNECTING / RECONNECTING
//
//////////////////////////////////////////////////////

// Connect establishes a socket connection with the Discord API
func (c *client) Connect() (err error) {
	c.Lock()
	defer c.Unlock()

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

	err = c.connectPermit.requestConnectPermit()
	if err != nil {
		err = errors.New("unable to get permission to Connect. Err: " + err.Error())
		return
	}

	// establish ws connection
	err = c.conn.Open(c.conf.Endpoint, nil)
	if err != nil {
		if !c.conn.Disconnected() {
			c.conn.Close()
			// TODO: logging
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
func (c *client) Disconnect() (err error) {
	c.Lock()
	defer c.Unlock()
	if c.conn.Disconnected() || !c.haveConnectedOnce {
		c.disconnected = true
		err = errors.New("already disconnected")
		return
	}

	// use the emitter to dispatch the close message
	err = c.conn.Close()
	if err != nil {
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
		c.Debug(fmt.Sprintf("EventReconnect attempt #%d\n", try))
		err = c.Connect()
		if err == nil {
			c.Info("successfully reconnected")
			break
		}

		err = c.connectPermit.releaseConnectPermit()
		if err != nil {
			err = errors.New("unable to release connection permission. Err: " + err.Error())
			c.Info(err.Error())
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
			err := c.Disconnect()
			if err != nil {
				c.Error(err.Error())
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
		// see operationHandler for return/Put()
		evt := c.poolDiscordPkt.Get().(*DiscordPacket)
		evt.reset()
		err = evt.UnmarshalJSON(packet)
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
