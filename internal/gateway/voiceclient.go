package gateway

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/gateway/cmd"
	"github.com/andersfylling/disgord/internal/gateway/opcode"
	"github.com/andersfylling/disgord/internal/logger"
	"github.com/andersfylling/disgord/json"
)

type VoiceConfig struct {
	// Guild ID to connect to
	GuildID Snowflake

	// User ID that is connecting
	UserID Snowflake

	// Session ID
	SessionID string

	// Token to connect with the voice websocket
	Token string

	// proxy allows for use of a custom proxy
	HTTPClient *http.Client

	// Endpoint for establishing voice connection
	Endpoint string

	// MessageQueueLimit number of outgoing messages that can be queued and sent correctly.
	MessageQueueLimit uint

	Logger logger.Logger

	SystemShutdown chan interface{}
}

func (conf *VoiceConfig) validate() {
	if conf.SystemShutdown == nil {
		panic("missing conf.SystemShutdown channel")
	}
}

type VoiceClient struct {
	*client
	conf *VoiceConfig

	haveIdentifiedOnce bool

	active         chan interface{}
	SystemShutdown chan interface{}
}

func NewVoiceClient(conf *VoiceConfig) (client *VoiceClient, err error) {
	conf.validate()

	client = &VoiceClient{
		conf: conf,
	}
	client.client, err = newClient(0, &config{
		Logger:     conf.Logger,
		Endpoint:   conf.Endpoint,
		HTTPClient: conf.HTTPClient,
		DiscordPktPool: &sync.Pool{
			New: func() interface{} {
				return &DiscordPacket{}
			},
		},
		messageQueueLimit: conf.MessageQueueLimit,
		SystemShutdown:    conf.SystemShutdown,
	}, client.internalConnect)
	if err != nil {
		return nil, err
	}
	client.clientType = clientTypeVoice
	client.setupBehaviors()

	return
}

//////////////////////////////////////////////////////
//
// BEHAVIORS
//
//////////////////////////////////////////////////////

func (c *VoiceClient) Active() <-chan interface{} {
	return c.SystemShutdown
}

func (c *VoiceClient) setupBehaviors() {
	// operation handlers
	// we manually link event methods instead of using reflection
	c.addBehavior(&behavior{
		addresses: discordOperations,
		actions: behaviorActions{
			opcode.VoiceReady:              c.onReady,
			opcode.VoiceResumed:            c.onResumed,
			opcode.VoiceHeartbeat:          c.onHeartbeatRequest,
			opcode.VoiceHeartbeatAck:       c.onHeartbeatAck,
			opcode.VoiceHello:              c.onHello,
			opcode.VoiceSessionDescription: c.onVoiceSessionDescription,
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
// BEHAVIOR: Discord Operations
//
//////////////////////////////////////////////////////

func (c *VoiceClient) onReady(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	readyPk := &VoiceReady{}
	if err = json.Unmarshal(p.Data, readyPk); err != nil {
		return err
	}

	if ch := c.onceChannels.Acquire(opcode.VoiceReady); ch != nil {
		ch <- readyPk
	} else {
		panic("once channel for Ready was missing")
	}
	return nil
}

func (c *VoiceClient) onResumed(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	resumedPk := &voicePacket{}
	if err = json.Unmarshal(p.Data, resumedPk); err != nil {
		return err
	}

	// TODO: use resumed instead..
	if ch := c.onceChannels.Acquire(opcode.VoiceReady); ch != nil {
		ch <- resumedPk
	} else {
		panic("once channel for Resumed was missing")
	}
	return nil
}

func (c *VoiceClient) onHeartbeatRequest(v interface{}) error {
	// https://discord.com/developers/docs/topics/gateway#heartbeating
	return c.emit(cmd.VoiceHeartbeat, nil)
}

func (c *VoiceClient) onHeartbeatAck(v interface{}) error {
	// heartbeat received
	c.Lock()
	c.lastHeartbeatAck = time.Now()
	c.Unlock()

	return nil
}

func (c *VoiceClient) onHello(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	type packet struct {
		// sometimes discord sends a float..............
		// How do you fuck up an integer, Discord?
		HeartbeatInterval float32 `json:"heartbeat_interval"`
	}
	helloPk := &packet{}
	if err = json.Unmarshal(p.Data, helloPk); err != nil {
		return err
	}
	interval := uint(helloPk.HeartbeatInterval)
	c.Lock()
	if interval == c.heartbeatInterval {
		c.Unlock()
		return nil
	} else if c.heartbeatInterval > 0 {
		c.Unlock()
		return errors.New("a new hello packet was sent, with a different interval - please make a github issue at https://github.com/andersfylling/disgord")
	}
	c.heartbeatInterval = interval
	c.Unlock()

	c.sendVoiceHelloPacket()
	c.activateHeartbeats <- true
	return nil
}

func (c *VoiceClient) onVoiceSessionDescription(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	sessionPk := &VoiceSessionDescription{}
	if err = json.Unmarshal(p.Data, sessionPk); err != nil {
		return err
	}

	if ch := c.onceChannels.Acquire(opcode.VoiceSessionDescription); ch != nil {
		ch <- sessionPk
	}
	return nil
}

//////////////////////////////////////////////////////
//
// BEHAVIOR: heartbeat
//
//////////////////////////////////////////////////////

func (c *VoiceClient) sendHeartbeat(i interface{}) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return c.emit(cmd.VoiceHeartbeat, r.Uint32())
}

//////////////////////////////////////////////////////
//
// GENERAL: unique to voice client
//
//////////////////////////////////////////////////////

// Connect establishes a socket connection with the Discord API
func (c *VoiceClient) Connect() (rdy *VoiceReady, err error) {
	var rdyI interface{}
	if rdyI, err = c.internalConnect(); rdyI != nil && err == nil {
		return rdyI.(*VoiceReady), nil
	}

	return nil, err
}

func (c *VoiceClient) internalConnect() (evt interface{}, err error) {
	// c.conn.Disconnected can always tell us if we are isConnected, but it cannot with
	// certainty say if we are connected
	if c.isConnected.Load() {
		err = errors.New("cannot Connect while a connection already exist")
		return nil, err
	}

	if c.conf.Endpoint == "" {
		panic("missing websocket endpoint. Must be set before constructing the sockets")
	}

	waitingChan := make(chan interface{}, 2)
	c.onceChannels.Add(opcode.VoiceReady, waitingChan)
	// TODO: explicitly add resumed as well
	defer func() {
		// cleanup
		c.onceChannels.Acquire(opcode.VoiceReady)
		close(waitingChan)
	}()

	// establish ws connection
	if err := c.conn.Open(context.Background(), c.conf.Endpoint, nil); err != nil {
		return nil, err
	}

	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())

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

	//errIdentify := c.emit(cmd.VoiceIdentify, &voiceIdentify{
	//	GuildID: c.conf.GuildID,
	//	UserID: c.conf.UserID,
	//	SessionID: c.conf.SessionID,
	//	Token: c.conf.Token,
	//})
	//if errIdentify != nil {
	//	c.log.Error(c.getLogPrefix(), "unable to send identify", errIdentify)
	//}

	select {
	case evt = <-waitingChan:
		c.log.Info(c.getLogPrefix(), "connected")
	case <-ctx.Done():
		c.isConnected.Store(false)
		err = errors.New("context cancelled")
	case <-time.After(5 * time.Second):
		c.isConnected.Store(false)
		err = errors.New("did not receive desired event in time. opcode " + strconv.Itoa(int(opcode.VoiceReady)))
	}
	return evt, err
}

// Emit for voice client needs to bypass the normal Emit restrictions.
// TODO: put more of the code flow of disgord/voiceclient.go into the websocket pkg.
func (c *VoiceClient) Emit(name string, data interface{}) error {
	return c.emit(name, data)
}

func (c *VoiceClient) sendVoiceHelloPacket() {
	// if this is a new connection we can drop the resume packet
	if !c.haveIdentifiedOnce {
		if err := sendVoiceIdentityPacket(c); err != nil {
			c.log.Error(c.getLogPrefix(), err)
		}
		return
	}

	resumeData := struct {
		GuildID   Snowflake `json:"server_id"`
		SessionID string    `json:"session_id"`
		Token     string    `json:"token"`
	}{c.conf.GuildID, c.conf.SessionID, c.conf.Token}
	_ = c.emit(cmd.VoiceResume, &resumeData)
}

func sendVoiceIdentityPacket(m *VoiceClient) (err error) {
	// https://discord.com/developers/docs/topics/gateway#identify
	err = m.emit(cmd.VoiceIdentify, &voiceIdentify{
		GuildID:   m.conf.GuildID,
		UserID:    m.conf.UserID,
		SessionID: m.conf.SessionID,
		Token:     m.conf.Token,
	})

	m.haveIdentifiedOnce = true
	return
}

func (c *VoiceClient) SendUDPInfo(data *VoiceSelectProtocolParams) (ret *VoiceSessionDescription, err error) {
	ch := make(chan interface{}, 1)
	c.onceChannels.Add(opcode.VoiceSessionDescription, ch)

	err = c.emit(cmd.VoiceSelectProtocol, &voiceSelectProtocol{
		Protocol: "udp",
		Data:     data,
	})
	if err != nil {
		return nil, err
	}

	select {
	case d := <-ch:
		ret = d.(*VoiceSessionDescription)
		return
	case <-time.After(5 * time.Second):
		err = errors.New("did not receive voice session description in time")
		return
	}
}
