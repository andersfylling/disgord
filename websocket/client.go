package websocket

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/andersfylling/disgord/constant"

	"github.com/andersfylling/disgord/httd"
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
		trackedEvents: config.TrackedEvents,
		eventChan:     eChan,
	}
	client.baseClient, err = newBaseClient(config, shardID)
	if err != nil {
		return nil, err
	}
	client.Start()

	return
}

func NewTestClient(config *Config, shardID uint, conn Conn) (*Client, chan interface{}) {
	s := make(chan interface{})

	c := &Client{
		trackedEvents: config.TrackedEvents,
		eventChan:     make(chan *Event),
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

	A A

	Logger constant.Logger

	// for identify packets
	Browser             string
	Device              string
	GuildLargeThreshold uint
	ShardCount          uint
}

type Client struct {
	*baseClient
	ReadyCounter uint

	eventChan        chan<- *Event
	trackedEvents    *UniqueStringSlice
	heartbeatLatency time.Duration

	sessionID string
	trace     []string
}

var _ constant.Logger = (*Client)(nil)

// Connect establishes a socket connection with the Discord API
func (c *Client) Connect() (err error) {
	return c.connect()
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
	close(c.shutdown)
	c.Disconnect()
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

				// This ignores the identify rate limit of 1/5s, because of the documentation stating:
				//  It's also possible that your client cannot reconnect in time to resume, in which case
				//  the client will receive a Opcode 9 Invalid Session and is expected to wait a random
				//  amount of time—between 1 and 5 seconds—then send a fresh Opcode 2 Identify.
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
			c.baseClient.heartbeatInterval = helloPk.HeartbeatInterval
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

	err = c.connectPermit.releaseConnectPermit()
	if err != nil {
		err = errors.New("unable to release connection permission. Err: " + err.Error())
		c.Error(err.Error())
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
	if err != nil {
		return err
	}

	err = c.connectPermit.releaseConnectPermit()
	if err != nil {
		err = errors.New("unable to release connection permission. Err: " + err.Error())
	}

	return
}
