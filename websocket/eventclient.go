package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/andersfylling/disgord/websocket/cmd"

	"github.com/andersfylling/disgord/logger"

	"golang.org/x/net/proxy"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
)

// NewManager creates a new socket client manager for handling behavior and Discord events. Note that this
// function initiates a go routine.
func NewEventClient(shardID uint, conf *EvtConfig) (client *EvtClient, err error) {
	if conf.TrackedEvents == nil {
		conf.TrackedEvents = &UniqueStringSlice{}
	}

	if conf.SystemShutdown == nil {
		panic("missing conf.SystemShutdown channel")
	}

	var eChan chan<- *Event
	if conf.EventChan != nil {
		eChan = conf.EventChan
	} else {
		err = errors.New("missing event channel")
		return nil, err
	}

	client = &EvtClient{
		evtConf:       conf,
		trackedEvents: conf.TrackedEvents,
		eventChan:     eChan,
	}
	client.client, err = newClient(shardID, &config{
		Logger:            conf.Logger,
		Endpoint:          conf.Endpoint,
		DiscordPktPool:    conf.DiscordPktPool,
		Proxy:             conf.Proxy,
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
		LargeThreshold: conf.GuildLargeThreshold,
		Shard:          &[2]uint{client.ShardID, conf.ShardCount},
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
	Proxy      proxy.Dialer
	HTTPClient *http.Client

	// for testing only
	conn Conn

	// TrackedEvents holds a list of predetermined events that should not be ignored.
	// This is especially useful for creating multiple shards, to reuse the same slice
	TrackedEvents *UniqueStringSlice

	// EventChan can be used to inject a channel instead of letting the ws client construct one
	// useful in sharding to avoid complicated patterns to handle N channels.
	EventChan chan<- *Event

	connectQueue func(shardID uint, cb func() error) error

	Presence interface{}

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

	DiscordPktPool *sync.Pool

	// MessageQueueLimit number of outgoing messages that can be queued and sent correctly.
	MessageQueueLimit uint

	Logger logger.Logger

	SystemShutdown chan interface{}
}

type EvtClient struct {
	evtConf *EvtConfig

	*client
	ReadyCounter uint

	eventChan     chan<- *Event
	trackedEvents *UniqueStringSlice

	sessionID      string
	sequenceNumber uint

	rdyPool *sync.Pool

	identity *evtIdentity
	idMu     sync.RWMutex
}

func (c *EvtClient) SetPresence(data interface{}) (err error) {
	// marshalling is done to avoid race
	var presence json.RawMessage
	if presence, err = httd.Marshal(data); err != nil {
		return err
	}
	c.idMu.Lock()
	c.identity.Presence = presence
	c.idMu.Unlock()

	return nil
}

func (c *EvtClient) Emit(command string, data interface{}) (err error) {
	if command == cmd.UpdateStatus {
		if err = c.SetPresence(data); err != nil {
			return err
		}
	}
	return c.client.emit(false, command, data)
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
	if p.SequenceNumber != c.sequenceNumber+1 {
		go c.reconnect()

		err = fmt.Errorf("websocket sequence numbers missmatch, forcing reconnect. Got %d, wants %d", p.SequenceNumber, c.sequenceNumber+1)
		return
	}

	// increment the sequence number for each event to make sure everything is synced with discord
	c.sequenceNumber++
	return nil
}

func (c *EvtClient) virginConnection() bool {
	return c.sessionID == "" && c.sequenceNumber == 0
}

func (c *EvtClient) onReady(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	// always store the session id
	ready := evtReadyPacket{}
	if err = httd.Unmarshal(p.Data, &ready); err != nil {
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

	// dispatch event through out the DisGord system
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
	if err := httd.Unmarshal(p.Data, helloPk); err != nil {
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
	c.Lock()
	c.sequenceNumber = 0
	c.Unlock()

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

	return c.emit(true, event.Heartbeat, snr)
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
	// c.conn.Disconnected can always tell us if we are disconnected, but it cannot with
	// certainty say if we are connected
	if !c.disconnected {
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
	c.haveConnectedOnce = true
	c.disconnected = false
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

// RegisterEvent tells the socket layer which event types are of interest. Any event that are not registered
// will be discarded once the socket info is extracted from the event.
func (c *EvtClient) RegisterEvent(event string) {
	c.trackedEvents.Add(event)
}

// RemoveEvent removes an event type from the registry. This will cause the event type to be discarded
// by the socket layer.
func (c *EvtClient) RemoveEvent(event string) {
	c.trackedEvents.Remove(event)
}

func (c *EvtClient) eventOfInterest(name string) bool {
	return c.trackedEvents.Exists(name)
}

func (c *EvtClient) sendHelloPacket() {
	c.RLock()
	token := c.evtConf.BotToken
	session := c.sessionID
	sequence := c.sequenceNumber
	c.RUnlock()

	err := c.emit(true, event.Resume, struct {
		Token      string `json:"token"`
		SessionID  string `json:"session_id"`
		SequenceNr uint   `json:"seq"`
	}{token, session, sequence})
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
	var id = &evtIdentity{}
	*id = *c.identity
	// copy it to avoid data race
	c.idMu.RUnlock()
	err = c.emit(true, event.Identify, id)

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
