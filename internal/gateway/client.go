package gateway

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/gateway/opcode"
	"github.com/andersfylling/disgord/internal/logger"
	"github.com/andersfylling/disgord/json"

	"go.uber.org/atomic"
)

type ClientType int

const (
	clientTypeEvent ClientType = iota
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

type connectQueue = func(shardID uint, cb func() error) error
type connectSignature = func() (evt interface{}, err error)
type discordErrListener = func(code int, reason string)

// newClient ...
func newClient(shardID uint, conf *config, connect connectSignature) (c *client, err error) {
	var ws Conn
	if conf.conn == nil {
		ws, err = newConn(conf.HTTPClient)
		if err != nil {
			return nil, err
		}
	} else {
		ws = conf.conn
	}

	var queueLimit int
	if conf.messageQueueLimit == 0 {
		queueLimit = 20
	} else {
		queueLimit = int(conf.messageQueueLimit)
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
		log:               conf.Logger,
		behaviors:         map[string]*behavior{},
		poolDiscordPkt:    conf.DiscordPktPool,
		onceChannels:      newOnceChannels(),
		connect:           connect,
		messageQueue:      newClientPktQueue(queueLimit),

		activateHeartbeats: make(chan interface{}),
		SystemShutdown:     conf.SystemShutdown,
	}
	c.isConnected.Store(false)

	return
}

type config struct {
	HTTPClient *http.Client

	// for testing only
	conn Conn

	// Endpoint for establishing socket connection. Either endpoints, `Gateway` or `Gateway Bot`, is used to retrieve
	// a valid socket endpoint from Discord
	Endpoint string

	DiscordPktPool *sync.Pool

	Logger logger.Logger

	discordErrListener discordErrListener

	// messageQueueLimit number of outgoing messages that can be queued and sent correctly.
	messageQueueLimit uint

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
	clientType  ClientType
	conf        *config
	lastRestart atomic.Int64 // unix nano

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
	conn             Conn
	messageQueue     clientPktQueue

	// connect is blocking until a websocket connection has completed it's setup.
	// eg. Normal shards that handles events are considered connected once the
	// identity/resume has been sent. While for voice we wait until a ready event
	// is returned.
	connect connectSignature

	// states
	isConnected       atomic.Bool
	haveConnectedOnce atomic.Bool
	isReconnecting    atomic.Bool
	isReceiving       atomic.Bool // has the go routine started
	isEmitting        atomic.Bool // has the go routine started
	onceChannels      onceChannels

	isRestarting atomic.Bool

	// identify timeout on invalid session
	// useful in unit tests when you want to drop any actual timeouts
	timeoutMultiplier int

	// ChannelBuffer is used to set the event channel buffer
	ChannelBuffer uint

	log         logger.Logger
	logSequence atomic.Uint32 // ARM 32bit causes panic with 64bit

	// behaviours - optional
	behaviors map[string]*behavior

	poolDiscordPkt *sync.Pool

	cancel context.CancelFunc

	SystemShutdown <-chan interface{}

	// receiver gets closed when the connection is lost
	requestedDisconnect atomic.Bool
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
		case p, open = <-c.receive():
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

func (c *client) inactivityDetector() {
	// make sure that websocket is connecting, connect or reconnecting.
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

	nr := c.logSequence.Add(1)
	s := "s:" + strconv.FormatUint(uint64(nr), 10)
	shardID := "shard:" + strconv.FormatUint(uint64(c.ShardID), 10)

	// [ws-?, s:0, shard:0]
	return "[" + t + "," + s + "," + shardID + "]"
}

//////////////////////////////////////////////////////
//
// LINKING: CONNECTING / DISCONNECTING / RECONNECTING
//
//////////////////////////////////////////////////////
func (c *client) IsDisconnected() bool {
	return !c.isConnected.Load()
}

func (c *client) disconnect() (err error) {
	c.Lock()
	defer c.Unlock()
	alreadyDisconnected := c.conn.Disconnected() || !c.haveConnectedOnce.Load() || c.cancel == nil

	// stop emitter, receiver and behaviors
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}

	// use the emitter to dispatch the close message
	err = c.conn.Close()
	// a typical err here is that the pipe is closed. Err is returned later

	// c.Emit(event.Close, nil)
	// dont use emit, such that we can call shutdown at the same time as Disconnect (See Shutdown())
	c.isConnected.Store(false)

	if alreadyDisconnected {
		return errors.New("already disconnected")
	}
	c.log.Info(c.getLogPrefix(), "disconnected")
	return err
}

// Disconnect disconnects the socket connection
func (c *client) Disconnect() (err error) {
	c.requestedDisconnect.Store(true)
	return c.disconnect()
}

func (c *client) reconnect() (err error) {
	if !c.isReconnecting.CAS(false, true) {
		return
	}
	c.lastRestart.Store(time.Now().UnixNano())
	defer c.isReconnecting.Store(false)

	c.log.Debug(c.getLogPrefix(), "is reconnecting")
	if err := c.disconnect(); err != nil {
		c.log.Debug(c.getLogPrefix(), "reconnecting failed: ", err.Error())
		c.RLock()
		if c.requestedDisconnect.Load() {
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
			c.log.Debug(c.getLogPrefix(), "reconnect attempt ", try)
		}
		if _, err = c.connect(); err == nil {
			c.log.Debug(c.getLogPrefix(), "establishing connection succeeded")
			break
		}
		c.log.Error(c.getLogPrefix(), "establishing connection failed: ", err)
		c.log.Info(c.getLogPrefix(), "next connection attempt in ", delay)

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

// Emit is used by Disgord users for dispatching a socket command to the Discord Gateway.
func (c *client) Emit(command string, data CmdPayload) (err error) {
	return c.queueRequest(command, data)
}

func (c *client) queueRequest(command string, data CmdPayload) (err error) {
	if !c.haveConnectedOnce.Load() {
		return errors.New("race condition detected: you must Connect to the socket API/Gateway before you can send gateway commands: " + command)
	}

	op := CmdNameToOpCode(command, c.clientType)
	if op == opcode.None {
		return errors.New("unsupported command: " + command)
	}

	p := &clientPacket{
		Op:      op,
		Data:    data,
		CmdName: command,
	}
	if accepted := c.ratelimit.Request(command); !accepted {
		// we might be rate limited.. but lets see if there is another
		// presence update in the queue; then it can be overwritten
		if err := c.messageQueue.AddByOverwrite(p); err != nil {
			return errors.New("rate limited")
		} else {
			return nil
		}
	}
	return c.messageQueue.Add(p)
}

func (c *client) emit(command string, data interface{}) (err error) {
	if !c.haveConnectedOnce.Load() {
		return errors.New("race condition detected: you must Connect to the socket API/Gateway before you can send gateway commands: " + command)
	}

	c.internalEmitChan <- &clientPacket{
		Op:      CmdNameToOpCode(command, c.clientType),
		Data:    data,
		CmdName: command,
	}
	return nil
}

// emitter holds the actually dispatching logic for sending data to the Discord Gateway.
// client#Emit depends on this.
func (c *client) emitter(ctx context.Context) {
	if !c.isEmitting.CAS(false, true) {
		return
	}
	defer c.isEmitting.Store(false)
	c.log.Debug(c.getLogPrefix(), "starting emitter")

	internal, cancel := context.WithCancel(context.Background())
	once := sync.Once{}

	write := func(msg *clientPacket) error {
		// save to file
		// build tag: disgord_diagnosews
		saveOutgoingPacket(c, msg)

		if err := c.conn.WriteJSON(msg); err != nil {
			once.Do(cancel)
			return err
		}
		return nil
	}

	for {
		var insight string
		var err error

		select {
		case <-ctx.Done():
			c.log.Debug(c.getLogPrefix(), "closing emitter")
			once.Do(cancel)
			return
		case <-internal.Done():
			c.log.Debug(c.getLogPrefix(), "closing emitter after write error")
			go c.reconnect()
			return
		case <-time.After(300 * time.Millisecond): // TODO: don't use fixed timeout
			if !c.messageQueue.IsEmpty() {
				// try to write the message
				// on failure the message is put back into the queue
				err = c.messageQueue.Try(write)
			}
		case msg, open := <-c.internalEmitChan:
			if !open {
				err = errors.New("emitter channel is closed")
			} else if err = write(msg); err != nil {
				insight = fmt.Sprintf("%v", *msg)
			}
		}

		if err != nil {
			c.log.Error(c.getLogPrefix(), err, insight)
		}
	}
}

//////////////////////////////////////////////////////
//
// RECEIVING
//
//////////////////////////////////////////////////////

// Receive returns the channel for receiving Discord packets
func (c *client) receive() <-chan *DiscordPacket {
	return c.receiveChan
}

func (c *client) receiver(ctx context.Context) {
	if !c.isReceiving.CAS(false, true) {
		return
	}
	defer c.isReceiving.Store(false)
	c.log.Debug(c.getLogPrefix(), "starting receiver")

	var noopCounter int

	internal, cancel := context.WithCancel(context.Background())
	once := sync.Once{}

	for {
		// check if application has closed
		// and clean up
		select {
		case <-ctx.Done():
			c.log.Debug(c.getLogPrefix(), "closing receiver")
			once.Do(cancel) // free
			return
		case <-internal.Done():
			c.log.Debug(c.getLogPrefix(), "closing receiver after read error")
			go func() {
				if err := c.reconnect(); err != nil {
					c.log.Error(c.getLogPrefix(), "reconnecting attempt failed: ", err.Error())
				}
			}()
			return
		default:
		}

		var packet []byte
		var err error
		if packet, err = c.conn.Read(ctx); err != nil {
			if !errors.Is(err, context.Canceled) && ctx.Err() != nil {
				c.log.Debug(c.getLogPrefix(), "read error: ", err.Error())
			}
			reconnect := true
			var closeErr *CloseErr
			isCloseErr := errors.As(err, &closeErr)
			if isCloseErr {
				if c.conf.discordErrListener != nil && closeErr.code >= 4000 && closeErr.code < 5000 {
					go c.conf.discordErrListener(closeErr.code, closeErr.info)
				}
				switch closeErr.code {
				case 4014:
					// Disconnected: Either the channel was deleted or you were kicked. Should not reconnect.
					// https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice-voice-close-event-codes
					c.log.Debug(c.getLogPrefix(), "discord sent a 4014 websocket code and the bot will now disconnect")
					_ = c.Disconnect()
					close(c.receiveChan) // notify client
					reconnect = false
				default:
				}
			}

			select {
			case <-ctx.Done():
				// in this case we dont want reconnect to start, only to stop
			default:
				if reconnect {
					once.Do(cancel)
				}
			}
			continue
		}

		// parse to gateway payload object
		// see operationHandler for return/Put()
		evt := c.poolDiscordPkt.Get().(*DiscordPacket)
		evt.reset()
		//err = evt.UnmarshalJSON(packet) // custom unmarshal
		if err = json.Unmarshal(packet, evt); err != nil {
			c.log.Error(c.getLogPrefix(), err, "SKIPPED ERRONEOUS PACKET CONTENT:", string(packet))
			c.poolDiscordPkt.Put(evt)

			// noop
			if noopCounter >= 10 {
				c.log.Error(c.getLogPrefix(), "json unmarshal failed 10 times for this shard and reconnect is now forced")
				once.Do(cancel) // on 10 continuous errors, we just force a reconnect
			}
			noopCounter++
			continue
		} else {
			noopCounter = 0
		}

		// save to file
		// build tag: disgord_diagnosews
		saveIncomingPacker(c, evt, packet)

		// notify listeners
		c.receiveChan <- evt
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
			// in case there is a race between when the ticker was started and the
			// heartbeat interval was updated
			c.RLock()
			interval2 := time.Millisecond * time.Duration(c.heartbeatInterval)
			if interval != interval2 {
				ticker.Stop()
				interval = interval2
				ticker = time.NewTicker(interval)
			}
			c.RUnlock()
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
