package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/json"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/andersfylling/disgord/internal/gateway/cmd"
	"github.com/andersfylling/disgord/internal/gateway/event"
	"github.com/andersfylling/disgord/internal/gateway/opcode"
	"github.com/andersfylling/disgord/internal/logger"
)

// NewManager creates a new socket client manager for handling behavior and Discord events. Note that this
// function initiates a go routine.
func NewEventClient(shardID uint, conf *EvtConfig) (client *EvtClient, err error) {
	conf.validate()

	var eChan chan<- *Event
	if conf.EventChan != nil {
		eChan = conf.EventChan
	} else {
		err = errors.New("missing event channel")
		return nil, err
	}

	client = &EvtClient{
		evtConf:      conf,
		ignoreEvents: conf.IgnoreEvents,
		eventChan:    eChan,
	}
	client.client, err = newClient(shardID, &config{
		Logger:            conf.Logger,
		Endpoint:          conf.Endpoint,
		DiscordPktPool:    conf.DiscordPktPool,
		HTTPClient:        conf.HTTPClient,
		conn:              conf.conn,
		messageQueueLimit: conf.MessageQueueLimit,

		SystemShutdown: conf.SystemShutdown,
	}, client.internalConnect)
	if err != nil {
		return nil, err
	}
	client.setupBehaviors()

	client.identity = &evtIdentity{
		Token: conf.BotToken,
		Properties: struct {
			OS      string `json:"$os"`
			Browser string `json:"$browser"`
			Device  string `json:"$device"`
		}{runtime.GOOS, conf.Browser, conf.Device},
		LargeThreshold:     conf.GuildLargeThreshold,
		Shard:              &[2]uint{client.ShardID, conf.ShardCount},
		GuildSubscriptions: conf.GuildSubscriptions,
		Intents:            conf.Intents,
	}
	if conf.Presence != nil {
		if err = client.SetPresence(conf.Presence); err != nil {
			return nil, err
		}
	}

	return
}

// Event is dispatched by the socket layer after parsing and extracting Discord data from a incoming packet.
// This is the data structure used by Disgord for triggering handlers and channels with an event.
type Event struct {
	Name    string
	Data    []byte
	ShardID uint
}

// EvtConfig ws
type EvtConfig struct {
	// BotToken Discord bot token
	BotToken   string
	HTTPClient *http.Client

	// for testing only
	conn Conn

	// IgnoreEvents holds a list of predetermined events that should be ignored.
	IgnoreEvents []string

	Intents Intent

	// EventChan can be used to inject a channel instead of letting the ws client construct one
	// useful in sharding to avoid complicated patterns to handle N channels.
	EventChan chan<- *Event

	connectQueue connectQueue

	discordErrListener discordErrListener

	Presence *UpdateStatusPayload

	// Endpoint for establishing socket connection. Either endpoints, `Gateway` or `Gateway Bot`, is used to retrieve
	// a valid socket endpoint from Discord
	Endpoint string

	// Encoding make sure we support the correct encoding
	Encoding string

	// Version make sure we support the correct Discord version
	Version int

	// for identify packets
	Browser             string
	Device              string
	GuildLargeThreshold uint
	ShardCount          uint
	GuildSubscriptions  bool

	DiscordPktPool *sync.Pool

	// MessageQueueLimit number of outgoing messages that can be queued and sent correctly.
	MessageQueueLimit uint

	Logger logger.Logger

	SystemShutdown chan interface{}
}

func (conf *EvtConfig) validate() {
	if conf.BotToken == "" {
		panic("missing bot token in gateway event config")
	}
	if conf.SystemShutdown == nil {
		panic("missing conf.SystemShutdown channel in gateway event config")
	}
}

type EvtClient struct {
	evtConf *EvtConfig

	*client
	ReadyCounter uint

	eventChan    chan<- *Event
	ignoreEvents []string

	sessionID      string
	sequenceNumber atomic.Uint32

	rdyPool *sync.Pool

	identity *evtIdentity
	idMu     sync.RWMutex
}

func (c *EvtClient) SetPresence(data interface{}) (err error) {
	// marshalling is done to avoid race
	var presence json.RawMessage
	if presence, err = json.Marshal(data); err != nil {
		return err
	}
	c.idMu.Lock()
	c.identity.Presence = presence
	c.idMu.Unlock()

	return nil
}

func (c *EvtClient) Emit(command string, data CmdPayload) (err error) {
	if command == cmd.UpdateStatus {
		if err = c.SetPresence(data); err != nil {
			return err
		}
	}
	return c.client.queueRequest(command, data)
}

//////////////////////////////////////////////////////
//
// BEHAVIORS
//
//////////////////////////////////////////////////////

func (c *EvtClient) setupBehaviors() {
	// operation handlers
	c.addBehavior(&behavior{
		addresses: discordOperations,
		actions: behaviorActions{
			opcode.EventDiscordEvent:   c.onDiscordEvent,
			opcode.EventHeartbeat:      c.onHeartbeatRequest,
			opcode.EventHeartbeatAck:   c.onHeartbeatAck,
			opcode.EventHello:          c.onHello,
			opcode.EventInvalidSession: c.onSessionInvalidated,
			opcode.EventReconnect: func(i interface{}) error {
				c.log.Info(c.getLogPrefix(), "Discord requested a reconnect")
				// There might be duplicate EventReconnect requests from Discord
				// this is therefore a goroutine such that reconnect requests that takes
				// place at the same time as the current one is discarded
				go c.reconnect()
				return nil
			},
		},
	})

	c.addBehavior(&behavior{
		addresses: heartbeating,
		actions: behaviorActions{
			sendHeartbeat: c.sendHeartbeat,
		},
	})
}

//////////////////////////////////////////////////////
//
// BEHAVIOR: Discord Operations & helpers
//
//////////////////////////////////////////////////////

func (c *EvtClient) synchronizeSnr(p *DiscordPacket) (err error) {
	c.Lock()
	defer c.Unlock()

	// validate the sequence numbers
	// ws/tcp only
	if p.SequenceNumber != c.sequenceNumber.Load()+1 {
		go c.reconnect()

		err = fmt.Errorf("websocket sequence numbers missmatch, forcing reconnect. Got %d, wants %d", p.SequenceNumber, c.sequenceNumber.Load()+1)
		return
	}

	// increment the sequence number for each event to make sure everything is synced with discord
	c.sequenceNumber.Inc()
	return nil
}

func (c *EvtClient) virginConnection() bool {
	return c.sessionID == "" && c.sequenceNumber.Load() == 0
}

func (c *EvtClient) onReady(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	// always store the session id
	ready := evtReadyPacket{}
	if err = json.Unmarshal(p.Data, &ready); err != nil {
		return err
	}

	c.Lock()
	c.sessionID = ready.SessionID
	c.ReadyCounter++
	c.Unlock()

	//if ch := c.onceChannels.Acquire(opcode.EventReadyResumed); ch != nil {
	//	ch <- ready
	//}

	return nil
}

func (c *EvtClient) onDiscordEvent(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	if err = c.synchronizeSnr(p); err != nil {
		return
	}

	if p.EventName == event.Ready {
		if err = c.onReady(p); err != nil {
			return err
		}
	}
	//} else if p.EventName == event.Resumed {
	//	if ch := c.onceChannels.Acquire(opcode.EventReadyResumed); ch != nil {
	//		// WARNING! does not return a ready event on resume!
	//		// TODO: clean up
	//		ch <- event.Resumed
	//	}
	//}

	if !c.eventOfInterest(p.EventName) {
		return nil
	}

	// dispatch event through out the Disgord system
	c.eventChan <- &Event{
		Name:    p.EventName,
		Data:    p.Data,
		ShardID: c.ShardID,
	}

	return nil
} // end onDiscordEvent

func (c *EvtClient) onHeartbeatRequest(v interface{}) error {
	return c.sendHeartbeat(v)
}

func (c *EvtClient) onHeartbeatAck(v interface{}) error {
	c.Lock()
	c.lastHeartbeatAck = time.Now()
	c.Unlock()

	return nil
}

func (c *EvtClient) onHello(v interface{}) error {
	p := v.(*DiscordPacket)

	helloPk := &helloPacket{}
	if err := json.Unmarshal(p.Data, helloPk); err != nil {
		return err
	}

	c.Lock()
	c.heartbeatInterval = helloPk.HeartbeatInterval
	c.Unlock()

	c.activateHeartbeats <- true

	// if this is a new connection we can drop the resume packet
	if c.virginConnection() {
		return sendIdentityPacket(false, c)
	}

	c.sendHelloPacket()
	return nil
}

func (c *EvtClient) onSessionInvalidated(v interface{}) error {
	// invalid session. Must respond with a identify packet
	c.log.Info(c.getLogPrefix(), "Discord invalidated session")

	// session is invalidated, reset the sequence number
	c.sequenceNumber.Store(0)

	rand.Seed(time.Now().UnixNano())
	delay := rand.Intn(4) + 1
	delay *= c.timeoutMultiplier
	randomDelay := time.Second * time.Duration(delay)

	// This ignores the identify rate limit of 1/5s, because of the documentation stating:
	//  It's also possible that your client cannot reconnect in time to resume, in which case
	//  the client will receive a Opcode 9 Invalid Session and is expected to wait a random
	//  amount of time—between 1 and 5 seconds—then send a fresh Opcode 2 EventIdentify.
	select {
	case <-time.After(randomDelay):
	case <-c.SystemShutdown:
		return errors.New("system is shutting down")
	}

	return sendIdentityPacket(true, c)
}

//////////////////////////////////////////////////////
//
// BEHAVIOR: heartbeat
//
//////////////////////////////////////////////////////

func (c *EvtClient) sendHeartbeat(i interface{}) error {
	c.RLock()
	snr := c.sequenceNumber
	c.RUnlock()

	return c.emit(event.Heartbeat, snr)
}

//////////////////////////////////////////////////////
//
// GENERAL: unique to event
//
//////////////////////////////////////////////////////
func (c *EvtClient) Connect() (err error) {
	_, err = c.internalConnect()
	return
}

func (c *EvtClient) internalConnect() (evt interface{}, err error) {
	if c.isConnected.Load() {
		err = errors.New("cannot Connect while a connection already exist")
		return nil, err
	}

	if c.conf.Endpoint == "" {
		err = errors.New("missing websocket endpoint. Must be set before constructing the sockets")
		return nil, err
	}

	var sessionCtx context.Context
	sessionCtx, c.cancel = context.WithCancel(context.Background())

	err = c.evtConf.connectQueue(c.ShardID, func() error {
		sentIdentifyResume := make(chan interface{})
		c.onceChannels.Add(opcode.EventIdentify, sentIdentifyResume)
		c.onceChannels.Add(opcode.EventResume, sentIdentifyResume)
		defer func() {
			// cleanup once channels
			c.onceChannels.Acquire(opcode.EventIdentify)
			c.onceChannels.Acquire(opcode.EventResume)
		}()

		if err := c.openConnection(sessionCtx); err != nil {
			return err
		}

		c.log.Debug(c.getLogPrefix(), "waiting to send identify/resume")
		select {
		case <-sessionCtx.Done():
			c.log.Info(c.getLogPrefix(), "session context was closed")
		case <-sentIdentifyResume:
			c.log.Debug(c.getLogPrefix(), "sent identify/resume")
		case <-time.After(3 * time.Minute):
			c.log.Error(c.getLogPrefix(), "discord timeout during connect (3 minutes). No idea what went wrong..")
			go c.reconnect()
			return errors.New("websocket connected but was not able to send identify packet within 3 minutes")
		}
		return nil
	})
	return nil, err
}

func (c *EvtClient) openConnection(ctx context.Context) error {
	// establish ws connection
	if err := c.conn.Open(ctx, c.conf.Endpoint, nil); err != nil {
		return err
	}

	// we can now interact with Discord
	c.haveConnectedOnce.Store(true)
	c.isConnected.Store(true)
	go c.receiver(ctx)
	go c.emitter(ctx)
	go c.startBehaviors(ctx)
	go c.prepareHeartbeating(ctx)
	go func() {
		select {
		case <-ctx.Done():
		case <-c.SystemShutdown:
			_ = c.Disconnect()
		}
	}()

	return nil
}

func (c *EvtClient) eventOfInterest(name string) bool {
	for i := range c.ignoreEvents {
		if c.ignoreEvents[i] == name {
			return false
		}
	}
	return true
}

func (c *EvtClient) sendHelloPacket() {
	c.RLock()
	token := c.evtConf.BotToken
	session := c.sessionID
	c.RUnlock()
	sequence := c.sequenceNumber.Load()

	err := c.emit(event.Resume, &evtResume{token, session, sequence})
	if err != nil {
		c.log.Error(c.getLogPrefix(), err)
	}

	c.log.Debug(c.getLogPrefix(), "sendHelloPacket is acquiring once channel")
	channel := c.onceChannels.Acquire(opcode.EventResume)
	c.log.Debug(c.getLogPrefix(), "writing to once channel", channel)
	channel <- true
	c.log.Debug(c.getLogPrefix(), "finished writing to once channel", channel)
}

func sendIdentityPacket(invalidSession bool, c *EvtClient) (err error) {
	c.idMu.RLock()
	var id = &evtIdentity{} // TODO: read only?
	*id = *c.identity
	// copy it to avoid data race
	c.idMu.RUnlock()
	err = c.emit(event.Identify, id)

	if !invalidSession {
		c.log.Debug(c.getLogPrefix(), "sendIdentityPacket is acquiring once channel")
		channel := c.onceChannels.Acquire(opcode.EventIdentify)
		c.log.Debug(c.getLogPrefix(), "writing to once channel", channel)
		channel <- true
		c.log.Debug(c.getLogPrefix(), "finished writing to once channel", channel)
	}
	return
}

var _ Link = (*EvtClient)(nil)
