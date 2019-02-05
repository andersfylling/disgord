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

	"github.com/andersfylling/disgord/websocket/cmd"

	"github.com/andersfylling/disgord/websocket"
	"golang.org/x/crypto/nacl/secretbox"
)

type voiceRepository struct {
	c *Client
	sync.Mutex

	pendingStates  map[Snowflake]chan *VoiceStateUpdate
	pendingServers map[Snowflake]chan *VoiceServerUpdate
}

type VoiceConnection interface {
	StartSpeaking() error
	StopSpeaking() error

	SendOpusFrame(data []byte)
	SendDCA(r io.Reader) error

	Close() error
}

type voiceImpl struct {
	sync.Mutex

	ready bool

	ws  *websocket.VoiceClient
	udp net.Conn

	ssrc      uint32
	secretKey [32]byte
	send      chan []byte
	close     chan struct{}
}

func newVoiceRepository(c *Client) (voice *voiceRepository) {
	return &voiceRepository{
		c: c,

		pendingStates:  make(map[Snowflake]chan *VoiceStateUpdate),
		pendingServers: make(map[Snowflake]chan *VoiceServerUpdate),
	}
}

func (r *voiceRepository) VoiceConnect(guildID, channelID Snowflake) (ret VoiceConnection, err error) {
	if guildID.Empty() {
		err = errors.New("guildID must be set to connect to a voice channel")
		return
	}
	if channelID.Empty() {
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
	err = r.c.Emit(CommandUpdateVoiceState, UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: &channelID,
		SelfDeaf:  false,
		SelfMute:  false,
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
		send:  make(chan []byte),
		close: make(chan struct{}),
	}
	// Defer a cleanup just in case
	defer func(v *voiceImpl) {
		if !v.ready {
			if v.ws != nil {
				_ = v.ws.Disconnect()
			}
			if v.udp != nil {
				_ = v.udp.Close()
			}
		}
	}(&voice)

	// Connect to the websocket
	voice.ws, err = websocket.NewVoiceClient(&websocket.VoiceConfig{
		GuildID:   server.GuildID,
		UserID:    r.c.myID,
		SessionID: state.SessionID,
		Token:     server.Token,
		Proxy:     r.c.config.Proxy,
		Endpoint:  "wss://" + strings.TrimSuffix(server.Endpoint, ":80") + "/?v=3",
	})
	if err != nil {
		return
	}

	var ready *websocket.VoiceReady
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
	var session *websocket.VoiceSessionDescription
	session, err = voice.ws.SendUDPInfo(&websocket.VoiceSelectProtocolParams{
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
	voice.ready = true

	go voice.opusSendLoop()

	ret = &voice
	return
}

func (r *voiceRepository) onVoiceStateUpdate(event *VoiceStateUpdate) {
	if event.UserID != r.c.myID {
		return
	}

	r.Lock()

	if ch, exists := r.pendingStates[event.VoiceState.GuildID]; exists {
		delete(r.pendingStates, event.VoiceState.GuildID)
		r.Unlock()

		ch <- event
	} else {
		r.Unlock()
	}
}

func (r *voiceRepository) onVoiceServerUpdate(event *VoiceServerUpdate) {
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

	if !v.ready {
		panic("Attempting to interact with a closed voice connection")
	}

	return v.ws.Emit(cmd.VoiceSpeaking, &voiceSpeakingData{
		Speaking: b,
		SSRC:     v.ssrc,
	})
}

func (v *voiceImpl) SendOpusFrame(data []byte) {
	if !v.ready {
		panic("Attempting to send to a closed voice connection")
	}
	v.send <- data
}

func (v *voiceImpl) SendDCA(r io.Reader) error {
	if !v.ready {
		panic("Attempting to send to a closed voice connection")
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

func (v *voiceImpl) Close() (err error) {
	v.Lock()
	defer v.Unlock()

	if !v.ready {
		panic("Attempting to close a closed Voice Connection")
	}

	close(v.close)
	close(v.send)
	_ = v.udp.Close()
	_ = v.ws.Disconnect()

	return
}

type voiceSpeakingData struct {
	Speaking bool   `json:"speaking"`
	Delay    int    `json:"delay"`
	SSRC     uint32 `json:"ssrc"`
}

func (v *voiceImpl) opusSendLoop() {
	// https://discordapp.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice
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
