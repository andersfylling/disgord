package websocket

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/logger"
	"github.com/andersfylling/disgord/websocket/cmd"
	"github.com/andersfylling/disgord/websocket/opcode"
	"golang.org/x/net/proxy"
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
	Proxy      proxy.Dialer
	HTTPClient *http.Client

	// Endpoint for establishing voice connection
	Endpoint string

	// MessageQueueLimit number of outgoing messages that can be queued and sent correctly.
	MessageQueueLimit uint

	Logger logger.Logger

	SystemShutdown chan interface{}
}

type VoiceClient struct {
	*client
	conf *VoiceConfig

	haveIdentifiedOnce bool

	SystemShutdown chan interface{}
}

func NewVoiceClient(conf *VoiceConfig) (client *VoiceClient, err error) {
	if conf.SystemShutdown == nil {
		panic("missing conf.SystemShutdown channel")
	}

	client = &VoiceClient{
		conf: conf,
	}
	client.client, err = newClient(0, &config{
		Logger:     conf.Logger,
		Endpoint:   conf.Endpoint,
		Proxy:      conf.Proxy,
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

func (c *VoiceClient) setupBehaviors() {
	// operation handlers
	// we manually link event methods instead of using reflection
	c.addBehavior(&behavior{
		addresses: discordOperations,
		actions: behaviorActions{
			opcode.VoiceReady:              c.onReady,
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
	if err = httd.Unmarshal(p.Data, readyPk); err != nil {
		return err
	}

	if ch := c.onceChannels.Acquire(opcode.VoiceReady); ch != nil {
		ch <- readyPk
	} else {
		panic("once channel for Ready was missing")
	}
	return nil
}

func (c *VoiceClient) onHeartbeatRequest(v interface{}) error {
	// https://discordapp.com/developers/docs/topics/gateway#heartbeating
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

	helloPk := &helloPacket{}
	if err = httd.Unmarshal(p.Data, helloPk); err != nil {
		return err
	}
	c.Lock()
	// From: https://discordapp.com/developers/docs/topics/voice-connections#heartbeating
	// There is currently a bug in the Hello payload heartbeat interval.
	// Until it is fixed, please take your heartbeat interval as `heartbeat_interval` * .75.
	// TODO This warning will be removed and a changelog published when the bug is fixed.
	c.heartbeatInterval = uint(float64(helloPk.HeartbeatInterval) * .75)
	c.Unlock()

	c.activateHeartbeats <- true

	c.sendVoiceHelloPacket()
	return nil
}

func (c *VoiceClient) onVoiceSessionDescription(v interface{}) (err error) {
	p := v.(*DiscordPacket)

	sessionPk := &VoiceSessionDescription{}
	if err = httd.Unmarshal(p.Data, sessionPk); err != nil {
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
	return c.emit(cmd.VoiceHeartbeat, nil)
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
	// c.conn.Disconnected can always tell us if we are disconnected, but it cannot with
	// certainty say if we are connected
	if !c.disconnected {
		err = errors.New("cannot Connect while a connection already exist")
		return nil, err
	}

	if c.conf.Endpoint == "" {
		panic("missing websocket endpoint. Must be set before constructing the sockets")
	}

	waitingChan := make(chan interface{}, 2)
	c.onceChannels.Add(opcode.VoiceReady, waitingChan)
	defer func() {
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

	select {
	case evt = <-waitingChan:
		c.log.Info(c.getLogPrefix(), "connected")
	case <-ctx.Done():
		c.disconnected = true
	case <-time.After(5 * time.Second):
		c.disconnected = true
		err = errors.New("did not receive desired event in time. opcode " + strconv.Itoa(int(opcode.VoiceReady)))
	}
	return evt, err
}

func (c *VoiceClient) sendVoiceHelloPacket() {
	// if this is a new connection we can drop the resume packet
	if !c.haveIdentifiedOnce {
		if err := sendVoiceIdentityPacket(c); err != nil {
			c.log.Error(c.getLogPrefix(), err)
		}
		return
	}

	_ = c.emit(cmd.VoiceResume, struct {
		GuildID   Snowflake `json:"server_id"`
		SessionID string    `json:"session_id"`
		Token     string    `json:"token"`
	}{c.conf.GuildID, c.conf.SessionID, c.conf.Token})
}

func sendVoiceIdentityPacket(m *VoiceClient) (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
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
