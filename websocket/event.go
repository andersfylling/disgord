package websocket

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/andersfylling/disgord/websocket/cmd"

	"github.com/andersfylling/disgord/constant"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
)

// NewManager creates a new socket client manager for handling behavior and Discord events. Note that this
// function initiates a go routine.
func NewEventClient(config *Config, shardID uint) (client *EvtClient, err error) {
	if config.TrackedEvents == nil {
		config.TrackedEvents = &UniqueStringSlice{}
	}

	var eChan chan<- *Event
	if config.EventChan != nil {
		eChan = config.EventChan
	} else {
		panic("missing event channel")
	}

	client = &EvtClient{
		trackedEvents: config.TrackedEvents,
		eventChan:     eChan,
	}
	client.client, err = newClient(config, shardID)
	if err != nil {
		return nil, err
	}
	client.setupBehaviors()
	client.start()

	return
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

	A A

	DiscordPktPool *sync.Pool

	Logger constant.Logger

	// for identify packets
	Browser             string
	Device              string
	GuildLargeThreshold uint
	ShardCount          uint
}

type EvtClient struct {
	*client
	ReadyCounter uint

	eventChan        chan<- *Event
	trackedEvents    *UniqueStringSlice
	heartbeatLatency time.Duration

	heartbeatInterval uint
	lastHeartbeatAck  time.Time

	sessionID      string
	trace          []string
	sequenceNumber uint

	pulsating  uint8
	pulseMutex sync.Mutex

	rdyPool *sync.Pool

	identity *evtIdentity
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
				c.Info("Discord requested a reconnect")
				// There might be duplicate EventReconnect requests from Discord
				// this is therefore a goroutine such that reconnect requests that takes
				// place at the same time as the current one is discarded
				go c.reconnect()
				return nil
			},
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

		err = fmt.Errorf("websocket sequence numbers missmatch, forcing reconnect. Got %d, wants %d", p.SequenceNumber, c.sequenceNumber)
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

	// always store the session id & update the trace content
	ready := evtReadyPacket{}
	if err = httd.Unmarshal(p.Data, &ready); err != nil {
		return err
	}

	c.Lock()
	c.sessionID = ready.SessionID
	c.trace = ready.Trace
	c.ReadyCounter++
	c.Unlock()

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

	if !c.eventOfInterest(p.EventName) {
		return nil
	}

	// dispatch event through out the DisGord system
	c.eventChan <- &Event{
		Name: p.EventName,
		Data: p.Data,
	}

	return nil
} // end onDiscordEvent

func (c *EvtClient) onHeartbeatRequest(v interface{}) error {
	// https://discordapp.com/developers/docs/topics/gateway#heartbeating
	c.RLock()
	snr := c.sequenceNumber
	c.RUnlock()

	return c.Emit(event.Heartbeat, snr)
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

	// TODO, this might create several idle goroutines..
	go c.pulsate()

	// if this is a new connection we can drop the resume packet
	if c.virginConnection() {
		return sendIdentityPacket(c)
	}

	c.sendHelloPacket()
	return nil
}

func (c *EvtClient) onSessionInvalidated(v interface{}) error {
	// invalid session. Must respond with a identify packet
	c.Info("Discord invalidated session")

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
	<-time.After(randomDelay)
	return sendIdentityPacket(c)
}

//////////////////////////////////////////////////////
//
// BEHAVIOR: EventHeartbeat
//
//////////////////////////////////////////////////////

// HeartbeatLatency get the time diff between sending a heartbeat and Discord replying with a heartbeat ack
func (c *EvtClient) HeartbeatLatency() (duration time.Duration, err error) {
	duration = c.heartbeatLatency
	if duration == 0 {
		err = errors.New("latency not determined yet")
	}

	return
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

	err = c.connectPermit.releaseConnectPermit()
	if err != nil {
		err = errors.New("unable to release connection permission. Err: " + err.Error())
		c.Error(err.Error())
	}
}

func sendIdentityPacket(c *EvtClient) (err error) {
	if c.identity == nil {
		// https://discordapp.com/developers/docs/topics/gateway#identify
		c.identity = &evtIdentity{
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
	}

	err = c.Emit(event.Identify, c.identity)
	err2 := c.connectPermit.releaseConnectPermit()
	if err == nil && err2 != nil {
		err = errors.New("unable to release connection permission. Err: " + err2.Error())
	}

	return
}

// AllowedToStartPulsating you must notify when you are done pulsating!
func (c *EvtClient) AllowedToStartPulsating(serviceID uint8) bool {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == 0 {
		c.pulsating = serviceID
	}

	return c.pulsating == serviceID
}

// StopPulsating stops sending heartbeats to Discord
func (c *EvtClient) StopPulsating(serviceID uint8) {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()

	if c.pulsating == serviceID {
		c.pulsating = 0
	}
}

func (c *EvtClient) pulsate() {
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
		go func(m *EvtClient, last time.Time, sent time.Time, cancel chan interface{}) {
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

var _ Link = (*EvtClient)(nil)
