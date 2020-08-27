package disgord

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/andersfylling/disgord/internal/gateway"
	"github.com/andersfylling/disgord/internal/gateway/cmd"

	"golang.org/x/crypto/nacl/secretbox"
)

type voiceRepository struct {
	sync.Mutex
	c *Client

	pendingStates  map[Snowflake]chan *VoiceStateUpdate
	pendingServers map[Snowflake]chan *VoiceServerUpdate
}

// VoiceConnection is the interface used to interact with active voice connections.
type VoiceConnection interface {
	// StartSpeaking should be sent before sending voice data.
	StartSpeaking() error
	// StopSpeaking should be sent after sending voice data. If there's a break in the sent data, you should not simply
	// stop sending data. Instead you have to send five frames of silence ([]byte{0xF8, 0xFF, 0xFE}) before stopping
	// to avoid unintended Opus interpolation with subsequent transmissions.
	StopSpeaking() error

	// SendOpusFrame sends a single frame of opus data to the UDP server. Frames are sent every 20ms with 960 samples (48kHz).
	//
	// if the bot has been disconnected or the channel removed, an error will be returned. The voice object must then be properly dealt with to avoid further issues.
	SendOpusFrame(data []byte) error
	// SendDCA reads from a Reader expecting a DCA encoded stream/file and sends them as frames.
	SendDCA(r io.Reader) error

	// MoveTo moves from the current voice channel to the given.
	MoveTo(channelID Snowflake) error

	// Close closes the websocket and UDP connection. This VoiceConnection interface will no
	// longer be usable.
	// It is the callers responsibility to ensure there are no concurrent calls to any other
	// methods of this interface after calling Close.
	Close() error
}

type voiceImpl struct {
	sync.Mutex

	ready atomic.Bool

	ws  *gateway.VoiceClient
	udp net.Conn

	ssrc      uint32
	secretKey [32]byte
	send      chan []byte
	close     chan struct{}

	guildID Snowflake
	c       *Client
}

func newVoiceRepository(c *Client) (voice *voiceRepository) {
	voice = &voiceRepository{
		c: c,

		pendingStates:  make(map[Snowflake]chan *VoiceStateUpdate),
		pendingServers: make(map[Snowflake]chan *VoiceServerUpdate),
	}
	c.On(EvtVoiceServerUpdate, voice.onVoiceServerUpdate)
	c.On(EvtVoiceStateUpdate, voice.onVoiceStateUpdate)

	return voice
}

func (r *voiceRepository) VoiceConnectOptions(guildID, channelID Snowflake, selfDeaf, selfMute bool) (ret VoiceConnection, err error) {
	if guildID.IsZero() {
		err = errors.New("guildID must be set to connect to a voice channel")
		return
	}
	if channelID.IsZero() {
		err = errors.New("channelID must be set to connect to a voice channel")
		return
	}

	// Set up some listeners for this connection attempt
	stateCh := make(chan *VoiceStateUpdate, 1)
	serverCh := make(chan *VoiceServerUpdate, 1)
	r.Lock()
	r.pendingStates[guildID] = stateCh
	r.pendingServers[guildID] = serverCh
	r.Unlock()
	defer func(r *voiceRepository, guildID Snowflake) {
		r.Lock()
		defer r.Unlock()

		delete(r.pendingStates, guildID)
		delete(r.pendingServers, guildID)
	}(r, guildID)

	// Tell Discord we want to connect to a channel
	_, err = r.c.Emit(UpdateVoiceState, &UpdateVoiceStatePayload{
		GuildID:   guildID,
		ChannelID: channelID,
		SelfDeaf:  true, //selfDeaf,
		SelfMute:  selfMute,
	})
	if err != nil {
		return
	}

	var (
		state  *VoiceStateUpdate
		server *VoiceServerUpdate
	)

	// Wait for the VoiceStateUpdate and VoiceServerUpdate, or else time out
	timeout := time.After(10 * time.Second)
waiter:
	for {
		select {
		case state = <-stateCh:
			if server != nil {
				break waiter
			}
		case server = <-serverCh:
			if state != nil {
				break waiter
			}
		case <-timeout:
			err = errors.New("timeout on receiving voice channel information from discord")
			return
		}
	}

	voice := voiceImpl{
		guildID: guildID,
		c:       r.c,
		send:    make(chan []byte),
		close:   make(chan struct{}),
	}
	// Defer a cleanup just in case
	defer func(v *voiceImpl) {
		if !v.ready.Load() {
			if v.ws != nil {
				_ = v.ws.Disconnect()
			}
			if v.udp != nil {
				_ = v.udp.Close()
			}
			close(v.close)
		}
	}(&voice)

	// Connect to the websocket
	voice.ws, err = gateway.NewVoiceClient(&gateway.VoiceConfig{
		GuildID:        server.GuildID,
		UserID:         r.c.myID,
		SessionID:      state.SessionID,
		Token:          server.Token,
		HTTPClient:     r.c.config.HTTPClient,
		Endpoint:       "wss://" + strings.TrimSuffix(server.Endpoint, ":80") + "/?v=4",
		Logger:         r.c.log,
		SystemShutdown: r.c.shutdownChan,
	})
	if err != nil {
		return
	}

	var ready *gateway.VoiceReady
	if ready, err = voice.ws.Connect(); err != nil {
		return
	}
	voice.ssrc = ready.SSRC

	// Connect to UDP
	dialer := net.Dial
	if r.c.config.Proxy != nil {
		dialer = r.c.config.Proxy.Dial
	}
	voice.udp, err = dialer("udp", ready.IP+":"+strconv.Itoa(ready.Port))
	if err != nil {
		return
	}

	// SendOpusFrame our SSRC with no further data for the IP discovery process.
	ssrcBuffer := make([]byte, 70)
	binary.BigEndian.PutUint32(ssrcBuffer, ready.SSRC)
	_, err = voice.udp.Write(ssrcBuffer)
	if err != nil {
		return
	}

	ipBuffer := make([]byte, 70)
	var n int
	n, err = voice.udp.Read(ipBuffer)
	if err != nil {
		return
	}
	if n < 70 {
		err = errors.New("udp packet received from discord is not the required 70 bytes")
		return
	}

	ipb := string(ipBuffer[4:68])
	nullPos := strings.Index(ipb, "\x00")
	if nullPos < 0 {
		err = errors.New("udp ip discovery did not contain a null terminator")
		return
	}
	ip := ipb[:nullPos]
	port := binary.LittleEndian.Uint16(ipBuffer[68:70])

	// Tell the websocket which encryption mode we want to use. We'll go with XSalsa20 and Poly1305 since that's what
	// libSodium/NaCl and golang.org/x/crypto/nacl/secretbox use. If both Discord and Go both start supporting more
	// modes "out of the box" we might want to consider implementing a "preferred mode selection" algorithm here.
	var session *gateway.VoiceSessionDescription
	session, err = voice.ws.SendUDPInfo(&gateway.VoiceSelectProtocolParams{
		Mode:    "xsalsa20_poly1305",
		Address: ip,
		Port:    port,
	})
	if err != nil {
		return
	}
	if session.Mode != "xsalsa20_poly1305" {
		err = errors.New("discord selected mismatching encryption algorithm")
		return
	}

	voice.secretKey = session.SecretKey
	voice.ready.Store(true)

	go voice.opusSendLoop()
	go voice.watcherDiscordCloseEvt()

	ret = &voice
	return
}

func (r *voiceRepository) onVoiceStateUpdate(_ Session, event *VoiceStateUpdate) {
	r.Lock()
	if event.UserID != r.c.myID {
		r.Unlock()
		return
	}

	if ch, exists := r.pendingStates[event.VoiceState.GuildID]; exists {
		delete(r.pendingStates, event.VoiceState.GuildID)
		r.Unlock()

		ch <- event
	} else {
		r.Unlock()
	}
}

func (r *voiceRepository) onVoiceServerUpdate(_ Session, event *VoiceServerUpdate) {
	r.Lock()

	if ch, exists := r.pendingServers[event.GuildID]; exists {
		delete(r.pendingStates, event.GuildID)
		r.Unlock()

		ch <- event
	} else {
		r.Unlock()
	}
}

func (v *voiceImpl) StartSpeaking() error {
	return v.speakingImpl(true)
}

func (v *voiceImpl) StopSpeaking() error {
	return v.speakingImpl(false)
}

func (v *voiceImpl) speakingImpl(b bool) error {
	v.Lock()
	defer v.Unlock()

	if !v.ready.Load() {
		return errors.New("attempting to interact with a closed voice connection")
	}

	return v.ws.Emit(cmd.VoiceSpeaking, &voiceSpeakingData{
		Speaking: b,
		SSRC:     v.ssrc,
	})
}

func (v *voiceImpl) SendOpusFrame(data []byte) error {
	if !v.ready.Load() {
		return errors.New("attempting to send to a closed voice connection")
	}
	v.send <- data
	return nil
}

func (v *voiceImpl) SendDCA(r io.Reader) error {
	if !v.ready.Load() {
		return errors.New("attempting to send to a closed voice connection")
	}

	var sampleSize uint16
	for {
		if err := binary.Read(r, binary.LittleEndian, &sampleSize); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		buf := make([]byte, sampleSize)
		if _, err := io.ReadFull(r, buf); err != nil {
			panic(err)
		}

		v.send <- buf
	}
}

func (v *voiceImpl) MoveTo(channelID Snowflake) error {
	if channelID.IsZero() {
		return errors.New("channelID must be set to move to a voice channel")
	}

	v.Lock()
	defer v.Unlock()

	if !v.ready.Load() {
		return errors.New("attempting to move in a closed Voice Connection")
	}

	_, _ = v.c.Emit(UpdateVoiceState, &UpdateVoiceStatePayload{
		GuildID:   v.guildID,
		ChannelID: channelID,
		SelfDeaf:  true, //false,
		SelfMute:  false,
	})

	return nil
}

func (v *voiceImpl) watcherDiscordCloseEvt() {
	for {
		var open bool
		select {
		case <-v.close:
			return
		case _, open = <-v.ws.Active():
		}
		if !open {
			break
		}
	}

	v.Lock()
	defer v.Unlock()

	if !v.ready.Load() {
		return
	}
	v.ready.Store(false)

	close(v.close)
	// clear send channel
	select {
	case <-v.send:
	default:
	}

	_ = v.udp.Close()
	_ = v.ws.Disconnect()
	close(v.send)

	//for range v.ws.Receive() {} // drain

	v.c.Logger().Info("Discord closed voice connection")
}

func (v *voiceImpl) Close() (err error) {
	v.Lock()
	defer v.Unlock()

	if !v.ready.Load() {
		return errors.New("attempting to close a closed Voice Connection")
	}

	defer func() {
		close(v.close)
		// clear send channel
		select {
		case <-v.send:
		default:
		}
		close(v.send)
	}()

	// if discord have already closed the connection
	// there is no need to send out a bunch of events
	if v.ws.IsDisconnected() {
		return v.udp.Close()
	}

	// Tell Discord we want to disconnect from channel/guild
	_, _ = v.c.Emit(UpdateVoiceState, &UpdateVoiceStatePayload{
		GuildID:   v.guildID,
		ChannelID: 0, // disconnect "code/value" (disgord implementation specific)
		SelfDeaf:  true,
		SelfMute:  true,
	})

	err1 := v.udp.Close()
	err2 := v.ws.Disconnect()

	if err1 != nil || err2 != nil {
		var errMsg string
		if err1 != nil {
			errMsg += err1.Error()
		}
		if err2 != nil {
			errMsg += err2.Error()
		}

		return errors.New(errMsg)
	}

	return nil
}

type voiceSpeakingData struct {
	Speaking bool   `json:"speaking"`
	Delay    int    `json:"delay"`
	SSRC     uint32 `json:"ssrc"`
}

func (v *voiceImpl) opusSendLoop() {
	// https://discord.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice
	header := make([]byte, 12)
	header[0] = 0x80
	header[1] = 0x78
	binary.BigEndian.PutUint32(header[8:12], v.ssrc)

	var (
		sequence  uint16
		timestamp uint32
		nonce     [24]byte

		msg  []byte
		open bool
	)

	frequency := time.NewTicker(time.Millisecond * 20) // 50 sends per sec, 960 samples each at 48kHz
	defer frequency.Stop()

	for {
		select {
		case msg, open = <-v.send:
			if !open {
				return
			}
		case <-v.close:
			return
		}

		binary.BigEndian.PutUint16(header[2:4], sequence)
		sequence++

		binary.BigEndian.PutUint32(header[4:8], timestamp)
		timestamp += 960 // samples

		copy(nonce[:], header)

		toSend := secretbox.Seal(header, msg, &nonce, &v.secretKey)
		select {
		case <-frequency.C:
		case <-v.close:
			return
		}

		_, _ = v.udp.Write(toSend)
		// err on udp write? hahahahahah... hahah.. good joke.
	}
}
