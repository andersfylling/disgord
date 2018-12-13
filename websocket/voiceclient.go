package websocket

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket/cmd"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"github.com/andersfylling/snowflake/v3"
	"github.com/sirupsen/logrus"
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
}

type VoiceClient struct {
	sync.RWMutex
	conf         *VoiceConfig
	shutdown     chan interface{}
	restart      chan interface{}
	lastRestart  int64 //unix
	restartMutex sync.Mutex

	heartbeatInterval uint
	heartbeatLatency  time.Duration
	lastHeartbeatAck  time.Time

	pulsating  uint8
	pulseMutex sync.Mutex

	receiveChan        chan *voicePacket
	emitChan           chan *clientPacket
	onceChannels       map[uint]chan interface{}
	conn               Conn
	disconnected       bool
	haveConnectedOnce  bool
	haveIdentifiedOnce bool

	isRestarting bool

	// identify timeout on invalid session
	timeoutMultiplier int
}

func NewVoiceClient(config *VoiceConfig) (client *VoiceClient, err error) {
	ws, err := newConn(config.Proxy)
	if err != nil {
		return nil, err
	}

	client = &VoiceClient{
		conf:              config,
		shutdown:          make(chan interface{}),
		restart:           make(chan interface{}),
		receiveChan:       make(chan *voicePacket),
		emitChan:          make(chan *clientPacket),
		onceChannels:      make(map[uint]chan interface{}),
		conn:              ws,
		timeoutMultiplier: 1,
		disconnected:      true,
	}
	go client.operationHandler()

	return
}

// Connect establishes a socket connection with the Discord API
func (m *VoiceClient) Connect() (rdy *VoiceReady, err error) {
	m.Lock()

	// m.conn.Disconnected can always tell us if we are disconnected, but it cannot with
	// certainty say if we are connected
	if !m.disconnected {
		err = errors.New("cannot connect while a connection already exist")
		m.Unlock()
		return
	}

	ch := make(chan interface{}, 1)
	m.onceChannels[opcode.VoiceReady] = ch

	// establish ws connection
	err = m.conn.Open(m.conf.Endpoint, nil)
	if err != nil {
		if !m.conn.Disconnected() {
			_ = m.conn.Close()
		}
		m.Unlock()
		return
	}

	// we can now interact with Discord
	m.haveConnectedOnce = true
	m.disconnected = false
	go m.receiver()
	go m.emitter()
	m.Unlock()

	timeout := time.After(5 * time.Second)
	select {
	case d := <-ch:
		rdy = d.(*VoiceReady)
		return
	case <-timeout:
		err = errors.New("did not receive voice ready in time")
		return
	}
}

// Disconnect disconnects the socket connection
func (m *VoiceClient) Disconnect() (err error) {
	m.Lock()
	defer m.Unlock()
	if m.conn.Disconnected() || !m.haveConnectedOnce {
		m.disconnected = true
		err = errors.New("already disconnected")
		return
	}

	// use the emitter to dispatch the close message
	_ = m.Emit(event.Close, nil)
	m.disconnected = true

	// close connection
	<-time.After(time.Second * 1 * time.Duration(m.timeoutMultiplier))

	// wait for processes
	<-time.After(time.Millisecond * 10)
	return
}

func (m *VoiceClient) reconnect() (err error) {
	// make sure there aren't multiple reconnect processes running
	if !m.lockRestart() {
		return
	}
	defer func() {
		m.isRestarting = false
	}()

	m.restart <- 1
	_ = m.Disconnect()

	var try uint
	var delay time.Duration = 3 // seconds
	for {
		logrus.Debugf("Reconnect attempt #%d\n", try)
		_, err = m.Connect()
		if err == nil {
			logrus.Info("successfully reconnected")
			break
		}

		// wait N seconds
		logrus.Infof("reconnect failed, trying again in N seconds; N =  %d", uint(delay))
		logrus.Info(err)
		select {
		case <-time.After(delay * time.Second):
			delay += 4 + time.Duration(try*2)
		case <-m.shutdown:
			return
		}

		if uint(delay) > 5*60 {
			delay = 60
		}
	}

	return
}

// Emit emits a command, if supported, and its data to the Discord Socket API
func (m *VoiceClient) Emit(command string, data interface{}) (err error) {
	if !m.haveConnectedOnce {
		return errors.New("race condition detected: you must connect to the socket API/Gateway before you can send gateway commands!")
	}

	var op uint
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

	default:
		err = errors.New("unsupported command: " + command)
		return
	}

	m.emitChan <- &clientPacket{
		Op:   op,
		Data: data,
	}
	return
}

// Receive returns the channel for receiving Discord packets
func (m *VoiceClient) Receive() <-chan *voicePacket {
	return m.receiveChan
}

// emitter holds the actually dispatching logic for the Emit method. See DefaultClient#Emit.
func (m *VoiceClient) emitter() {
	for {
		var msg *clientPacket
		var open bool

		select {
		case <-m.shutdown:
			// m.connection got closed
		case msg, open = <-m.emitChan:
		}
		if !open || (msg.Data == nil && (msg.Op == opcode.Shutdown || msg.Op == opcode.Close)) {
			// TODO: what if we get a connection error, how do we restart?
			_ = m.conn.Close()
			return
		}

		logrus.WithField("op", msg.Op).Info("Sending vop")
		err := m.conn.WriteJSON(msg)
		if err != nil {
			// TODO-logging
			fmt.Printf("could not send data to discord: %+v\n", msg)
		}
	}
}

func (m *VoiceClient) receiver() {
	for {
		packet, err := m.conn.Read()
		if err != nil {
			logrus.WithError(err).Debug("closing voice readPump")
			return
		}

		// parse to gateway payload object
		evt := &voicePacket{}
		err = httd.Unmarshal(packet, &evt)
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.WithField("op", evt.Op).Info("Receiving vop")

		// notify listeners
		m.receiveChan <- evt

		// check if application has closed
		select {
		case <-m.shutdown:
			return
		default:
		}
	}
}

func (m *VoiceClient) operationHandler() {
	logrus.Debug("Ready to receive voice operation codes...")
	for {
		var p *voicePacket
		var open bool
		select {
		case p, open = <-m.Receive():
			if !open {
				logrus.Debug("voice operationChan is dead..")
				return
			}
		// case <-m.restart:
		case <-m.shutdown:
			logrus.Debug("exiting voice operation handler")
			return
		}

		// new packet that must be handled by it's Discord operation code
		switch p.Op {
		case opcode.VoiceHeartbeat:
			// https://discordapp.com/developers/docs/topics/gateway#heartbeating
			_ = m.Emit(cmd.VoiceHeartbeat, nil)
		case opcode.VoiceHeartbeatAck:
			// heartbeat received
			m.Lock()
			m.lastHeartbeatAck = time.Now()
			m.Unlock()
		case opcode.VoiceHello:
			// hello
			helloPk := &helloPacket{}
			err := httd.Unmarshal(p.Data, helloPk)
			if err != nil {
				logrus.Debug(err)
			}
			m.Lock()
			// From: https://discordapp.com/developers/docs/topics/voice-connections#heartbeating
			// There is currently a bug in the Hello payload heartbeat interval.
			// Until it is fixed, please take your heartbeat interval as `heartbeat_interval` * .75.
			// TODO This warning will be removed and a changelog published when the bug is fixed.
			m.heartbeatInterval = uint(float64(helloPk.HeartbeatInterval) * .75)
			m.Unlock()

			m.sendVoiceHelloPacket()
		case opcode.VoiceReady:
			readyPk := &VoiceReady{}
			err := httd.Unmarshal(p.Data, readyPk)
			if err != nil {
				logrus.Debug(err)
			}
			m.Lock()
			if ch, ok := m.onceChannels[opcode.VoiceReady]; ok {
				delete(m.onceChannels, opcode.VoiceReady)
				ch <- readyPk
			}
			m.Unlock()
		case opcode.VoiceSessionDescription:
			sessionPk := &VoiceSessionDescription{}
			err := httd.Unmarshal(p.Data, sessionPk)
			if err != nil {
				logrus.Debug(err)
			}
			m.Lock()
			if ch, ok := m.onceChannels[opcode.VoiceSessionDescription]; ok {
				delete(m.onceChannels, opcode.VoiceSessionDescription)
				ch <- sessionPk
			}
			m.Unlock()
		default:
			// unknown
			logrus.Debugf("Unknown voice operation: %+v\n", p)
		}
	}
}

func (m *VoiceClient) sendVoiceHelloPacket() {
	go m.pulsate()

	// if this is a new connection we can drop the resume packet
	if !m.haveIdentifiedOnce {
		err := sendVoiceIdentityPacket(m)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	_ = m.Emit(cmd.VoiceResume, struct {
		GuildID   snowflake.Snowflake `json:"server_id"`
		SessionID string              `json:"session_id"`
		Token     string              `json:"token"`
	}{m.conf.GuildID, m.conf.SessionID, m.conf.Token})
}

func sendVoiceIdentityPacket(m *VoiceClient) (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	identityPayload := struct {
		GuildID   snowflake.Snowflake `json:"server_id"` // Yay for eventual consistency
		UserID    snowflake.Snowflake `json:"user_id"`
		SessionID string              `json:"session_id"`
		Token     string              `json:"token"`
	}{
		GuildID:   m.conf.GuildID,
		UserID:    m.conf.UserID,
		SessionID: m.conf.SessionID,
		Token:     m.conf.Token,
	}

	err = m.Emit(cmd.VoiceIdentify, &identityPayload)
	m.haveIdentifiedOnce = true
	return
}

func (m *VoiceClient) SendUDPInfo(data *VoiceSelectProtocolParams) (ret *VoiceSessionDescription, err error) {
	ch := make(chan interface{}, 1)
	m.onceChannels[opcode.VoiceSessionDescription] = ch

	err = m.Emit(cmd.VoiceSelectProtocol, &voiceSelectProtocol{
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

func (m *VoiceClient) lockRestart() bool {
	m.restartMutex.Lock()
	defer m.restartMutex.Unlock()

	now := time.Now().UnixNano()
	locked := (now - m.lastRestart) > (time.Second.Nanoseconds() / 2)

	if locked && !m.isRestarting {
		m.lastRestart = now
		m.isRestarting = true
	}

	return locked
}

// AllowedToStartPulsating you must notify when you are done pulsating!
func (m *VoiceClient) AllowedToStartPulsating(serviceID uint8) bool {
	m.pulseMutex.Lock()
	defer m.pulseMutex.Unlock()

	if m.pulsating == 0 {
		m.pulsating = serviceID
	}

	return m.pulsating == serviceID
}

// StopPulsating stops sending heartbeats to Discord
func (m *VoiceClient) StopPulsating(serviceID uint8) {
	m.pulseMutex.Lock()
	defer m.pulseMutex.Unlock()

	if m.pulsating == serviceID {
		m.pulsating = 0
	}
}

func (m *VoiceClient) pulsate() {
	serviceID := uint8(rand.Intn(254) + 1) // uint8 cap
	if !m.AllowedToStartPulsating(serviceID) {
		return
	}
	defer m.StopPulsating(serviceID)

	m.RLock()
	ticker := time.NewTicker(time.Millisecond * time.Duration(m.heartbeatInterval))
	m.RUnlock()
	defer ticker.Stop()

	var last time.Time
	for {
		m.RLock()
		last = m.lastHeartbeatAck
		m.RUnlock()

		_ = m.Emit(cmd.VoiceHeartbeat, nil)

		stopChan := make(chan interface{})

		// verify the heartbeat ACK
		go func(m *VoiceClient, last time.Time, sent time.Time, cancel chan interface{}) {
			select {
			case <-cancel:
				return
			case <-time.After(3 * time.Second): // deadline for Discord to respond
			}

			m.RLock()
			receivedHeartbeatAck := m.lastHeartbeatAck.After(last)
			m.RUnlock()

			if !receivedHeartbeatAck {
				logrus.Info("heartbeat ACK was not received, forcing reconnect")
				_ = m.reconnect()
			} else {
				// update "latency"
				m.heartbeatLatency = m.lastHeartbeatAck.Sub(sent)
			}
		}(m, last, time.Now(), stopChan)

		select {
		case <-ticker.C:
			continue
		case <-m.shutdown:
		case <-m.restart:
		}

		logrus.Debug("Stopping pulse")
		close(stopChan)
		return
	}
}
