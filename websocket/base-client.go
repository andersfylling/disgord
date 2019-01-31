package websocket

import (
	"errors"
	"fmt"
	"math/rand"
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

type connectPermit interface {
	requestConnectPermit() error
	releaseConnectPermit() error
}

// newBaseClient ...
func newBaseClient(config *Config, shardID uint) (client *baseClient, err error) {
	ws, err := newConn(config.HTTPClient)
	if err != nil {
		return nil, err
	}

	client = &baseClient{
		conf:              config,
		shutdown:          make(chan interface{}),
		restart:           make(chan interface{}),
		ShardID:           shardID,
		receiveChan:       make(chan *discordPacket),
		emitChan:          make(chan *clientPacket),
		conn:              ws,
		ratelimit:         newRatelimiter(),
		timeoutMultiplier: 1,
		disconnected:      true,
		log:               config.Logger,
		a:                 config.A,
	}
	client.connectPermit = client

	return
}

// baseClient can be used as a base for other ws clients; voice, event. Note the use of
// baseClient.ReleasePermit() and baseClient.RequestPermit() in connect (and then automatically reconnect()).
// these are used for synchronizing connecting and you must therefore correctly release the permit once you have
// such that the next shard or yourself can connect in the future.
//
// If you do not care about these. Please overwrite both methods.
type baseClient struct {
	sync.RWMutex
	conf         *Config
	shutdown     chan interface{}
	restart      chan interface{}
	lastRestart  int64 //unix
	restartMutex sync.Mutex
	ReadyCounter uint

	clientType int

	ShardID uint

	heartbeatInterval uint
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

	connectPermit connectPermit

	// choreographic programming to handle rate limit, reconnects, and connects
	K *K
	a A
}

func (c *baseClient) requestConnectPermit() error {
	c.Debug("trying to get connect permission")
	b := make(B)
	defer close(b)
	c.a <- b
	c.Info("waiting")
	var ok bool
	select {
	case c.K, ok = <-b:
		if !ok || c.K == nil {
			c.Debug("unable to get connect permission")
			return errors.New("channel closed or K was nil")
		}
		c.Debug("got connect permission")
	case <-c.shutdown:
	}

	return nil
}

func (c *baseClient) releaseConnectPermit() error {
	if c.K == nil {
		return errors.New("K has not been granted yet")
	}

	c.K.Release <- c.K
	c.K = nil
	return nil
}

func (c *baseClient) getLogPrefix() string {
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

func (c *baseClient) Info(msg string) {
	if c.log != nil {
		c.log.Info(c.getLogPrefix() + msg)
	}
}
func (c *baseClient) Debug(msg string) {
	if c.log != nil {
		c.log.Debug(c.getLogPrefix() + msg)
	}
}
func (c *baseClient) Error(msg string) {
	if c.log != nil {
		c.log.Error(c.getLogPrefix() + msg)
	}
}

var _ constant.Logger = (*baseClient)(nil)

// Connect establishes a socket connection with the Discord API
func (c *baseClient) connect() (err error) {
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

	err = c.connectPermit.requestConnectPermit()
	if err != nil {
		err = errors.New("unable to get permission to connect. Err: " + err.Error())
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
func (c *baseClient) Disconnect() (err error) {
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

func (c *baseClient) lockReconnect() bool {
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

func (c *baseClient) unlockReconnect() {
	c.restartMutex.Lock()
	defer c.restartMutex.Unlock()

	c.isReconnecting = false
}

func (c *baseClient) reconnect() (err error) {
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
		c.Debug(fmt.Sprintf("Reconnect attempt #%d\n", try))
		err = c.connect()
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

// Emit emits a command, if supported, and its data to the Discord Socket API
func (c *baseClient) Emit(command string, data interface{}) (err error) {
	if !c.haveConnectedOnce {
		return errors.New("race condition detected: you must connect to the socket API/Gateway before you can send gateway commands")
	}

	op := ^uint(0)
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

// Receive returns the channel for receiving Discord packets
func (c *baseClient) Receive() <-chan *discordPacket {
	return c.receiveChan
}

func (c *baseClient) lockEmitter() bool {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	if !c.isEmitting {
		c.isEmitting = true
		return true
	}

	return false
}

func (c *baseClient) unlockEmitter() {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	c.isEmitting = false
}

// emitter holds the actually dispatching logic for the Emit method. See DefaultClient#Emit.
func (c *baseClient) emitter() {
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

func (c *baseClient) lockReceiver() bool {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	if !c.isReceiving {
		c.isReceiving = true
		return true
	}

	return false
}

func (c *baseClient) unlockReceiver() {
	c.recEmitMutex.Lock()
	defer c.recEmitMutex.Unlock()

	c.isReceiving = false
}

func (c *baseClient) receiver() {
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

func (c *baseClient) Shutdown() (err error) {
	close(c.shutdown)
	c.Disconnect()
	return
}

// AllowedToStartPulsating you must notify when you are done pulsating!
func (c *baseClient) AllowedToStartPulsating(serviceID uint8) bool {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == 0 {
		c.pulsating = serviceID
	}

	return c.pulsating == serviceID
}

// StopPulsating stops sending heartbeats to Discord
func (c *baseClient) StopPulsating(serviceID uint8) {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == serviceID {
		c.pulsating = 0
	}
}

func (c *baseClient) pulsate() {
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

		command := event.Heartbeat
		if c.clientType == clientTypeVoice {
			command = cmd.VoiceHeartbeat
		}
		_ = c.Emit(command, snr)

		stopChan := make(chan interface{})

		// verify the heartbeat ACK
		go func(m *baseClient, last time.Time, sent time.Time, cancel chan interface{}) {
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
