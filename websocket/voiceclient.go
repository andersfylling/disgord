package websocket

import (
	"errors"
	"sync"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/logger"
	"github.com/andersfylling/disgord/websocket/cmd"
	"github.com/andersfylling/disgord/websocket/opcode"
	"github.com/andersfylling/snowflake/v3"
	"golang.org/x/net/proxy"
)

type VoiceConfig struct {
	// Guild ID to connect to
	GuildID snowflake.Snowflake

	// User ID that is connecting
	UserID snowflake.Snowflake

	// Session ID
	SessionID string

	// Token to connect with the voice websocket
	Token string

	// Proxy allows for use of a custom proxy
	Proxy proxy.Dialer

	// Endpoint for establishing voice connection
	Endpoint string

	Logger logger.Logger

	SystemShutdown chan interface{}
}

type VoiceClient struct {
	*client
	conf *VoiceConfig

	haveConnectedOnce  bool
	haveIdentifiedOnce bool

	onceChannels map[uint]chan interface{}
	ready        *VoiceReady

	SystemShutdown chan interface{}
}

func NewVoiceClient(conf *VoiceConfig) (client *VoiceClient, err error) {
	if conf.SystemShutdown == nil {
		panic("missing conf.SystemShutdown channel")
	}

	client = &VoiceClient{
		conf:         conf,
		onceChannels: make(map[uint]chan interface{}),
	}
	client.client, err = newClient(&config{
		Logger:   conf.Logger,
		Endpoint: conf.Endpoint,
		DiscordPktPool: &sync.Pool{
			New: func() interface{} {
				return &DiscordPacket{}
			},
		},

		SystemShutdown: conf.SystemShutdown,
	}, 0)
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

	ch := make(chan interface{}, 1)
	c.preConnect(func() {
		c.ready = nil
		c.onceChannels[opcode.VoiceReady] = ch
	})
	c.postConnect(func() {
		timeout := time.After(5 * time.Second)
		select {
		case d := <-ch:
			c.ready = d.(*VoiceReady)
		case <-timeout:
			c.Error("did not receive voice ready in time")
		}
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

	c.Lock()
	if ch, ok := c.onceChannels[opcode.VoiceReady]; ok {
		delete(c.onceChannels, opcode.VoiceReady)
		ch <- readyPk
	}
	c.Unlock()
	return nil
}

func (c *VoiceClient) onHeartbeatRequest(v interface{}) error {
	// https://discordapp.com/developers/docs/topics/gateway#heartbeating
	return c.Emit(cmd.VoiceHeartbeat, nil)
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

	c.Lock()
	if ch, ok := c.onceChannels[opcode.VoiceSessionDescription]; ok {
		delete(c.onceChannels, opcode.VoiceSessionDescription)
		ch <- sessionPk
	}
	c.Unlock()
	return nil
}

//////////////////////////////////////////////////////
//
// BEHAVIOR: heartbeat
//
//////////////////////////////////////////////////////

func (c *VoiceClient) sendHeartbeat(i interface{}) error {
	return c.Emit(cmd.VoiceHeartbeat, nil)
}

//////////////////////////////////////////////////////
//
// GENERAL: unique to voice client
//
//////////////////////////////////////////////////////

// Connect establishes a socket connection with the Discord API
func (c *VoiceClient) Connect() (rdy *VoiceReady, err error) {
	if err = c.client.Connect(); err != nil {
		return nil, err
	}

	// TODO: plausible race condition
	c.Lock()
	defer c.Unlock()
	return c.ready, nil
}

func (c *VoiceClient) sendVoiceHelloPacket() {
	// if this is a new connection we can drop the resume packet
	if !c.haveIdentifiedOnce {
		if err := sendVoiceIdentityPacket(c); err != nil {
			c.Error(err)
		}
		return
	}

	_ = c.Emit(cmd.VoiceResume, struct {
		GuildID   snowflake.Snowflake `json:"server_id"`
		SessionID string              `json:"session_id"`
		Token     string              `json:"token"`
	}{c.conf.GuildID, c.conf.SessionID, c.conf.Token})
}

func sendVoiceIdentityPacket(m *VoiceClient) (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	err = m.Emit(cmd.VoiceIdentify, &voiceIdentify{
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
	c.onceChannels[opcode.VoiceSessionDescription] = ch

	err = c.Emit(cmd.VoiceSelectProtocol, &voiceSelectProtocol{
		Protocol: "udp",
		Data:     data,
	})
	if err != nil {
		return
	}

	timeout := time.After(5 * time.Second)
	select {
	case d := <-ch:
		ret = d.(*VoiceSessionDescription)
		return
	case <-timeout:
		err = errors.New("did not receive voice session description in time")
		return
	}
}
