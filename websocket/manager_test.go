package websocket

import (
	"errors"
	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"net/http"
	"sync"
	"testing"
	"time"
)

func newTestClient() (*testClient, error) {
	return &testClient{
		receiveChan:    make(chan *discordPacket),
		emitChan:       make(chan *clientPacket),
		connectingChan: make(chan interface{}),
	}, nil
}

type testClient struct {
	receiveChan    chan *discordPacket
	emitChan       chan *clientPacket
	connectingChan chan interface{}
	connected      int
	disconnected   int
}

func (c *testClient) Connect() error {
	c.connected++
	c.connectingChan <- 1
	return nil
}
func (c *testClient) Disconnect() error {
	c.disconnected++
	return nil
}
func (c *testClient) Emit(command string, data interface{}) error {
	var op uint
	switch command {
	case event.Shutdown:
		op = opcode.Shutdown
	case event.Heartbeat:
		op = opcode.Heartbeat
	case event.Identify:
		op = opcode.Identify
	case event.Resume:
		op = opcode.Resume
	case event.RequestGuildMembers:
		op = opcode.RequestGuildMembers
	case event.VoiceStateUpdate:
		op = opcode.VoiceStateUpdate
	case event.StatusUpdate:
		op = opcode.StatusUpdate
	default:
		return errors.New("unsupported command: " + command)
	}

	c.emitChan <- &clientPacket{
		Op:   op,
		Data: data,
	}
	return nil
}
func (c *testClient) Receive() <-chan *discordPacket {
	return c.receiveChan
}

func TestManager_RegisterEvent(t *testing.T) {
	m := Manager{}
	t1 := "test"
	m.RegisterEvent(t1)

	if len(m.trackedEvents) == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RegisterEvent(t1)
	if len(m.trackedEvents) == 2 {
		t.Error("expected length to be 1, got 2")
	}
}

func TestManager_RemoveEvent(t *testing.T) {
	m := Manager{}
	t1 := "test"
	m.RegisterEvent(t1)

	if len(m.trackedEvents) == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RemoveEvent("sdfsdf")
	if len(m.trackedEvents) == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RemoveEvent(t1)
	if len(m.trackedEvents) == 1 {
		t.Error("expected length to be 0, got 1")
	}
}

func TestManager_reconnect(t *testing.T) {
	client, _ := newTestClient()
	m := &Manager{
		conf: &ManagerConfig{
			// identity
			Browser:             "disgord",
			Device:              "disgord",
			GuildLargeThreshold: 250,
			ShardID:             0,
			ShardCount:          0,

			DefaultClientConfig: DefaultClientConfig{
				// lib specific
				Version:       constant.DiscordVersion,
				Encoding:      constant.JSONEncoding,
				ChannelBuffer: 1,
				Endpoint:      "sfkjsdlfsf",

				// user settings
				Token: "sifhsdoifhsdifhsdf",
				HTTPClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
		},
		shutdown:  make(chan interface{}),
		restart:   make(chan interface{}),
		Client:    client,
		eventChan: make(chan *Event),
	}
	seq := uint(1)

	shutdown := make(chan interface{})
	done := make(chan interface{})

	defer func() {
		m.Shutdown()
		close(done)
	}()

	var wgID sync.WaitGroup

	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case <-client.connectingChan:
			case <-shutdown:
				return
			case <-done:
				return
			}
			wg.Done()
		}
	}()
	wg.Add(1)
	m.Connect()
	wg.Wait()
	if client.connected != 1 {
		t.Error("expected 1 connect, got 0")
	}

	go m.operationHandlers()
	go func(t *testing.T) {
		select {
		case <-time.After(6 * time.Second):
		case <-done:
			return
		}
		close(shutdown)
		t.Error("timeout")
	}(t)

	go func() {
		for {
			select {
			case <-shutdown:
				return
			case <-done:
				return
			case <-m.eventChan:
			}
		}
	}()

	// heartbeat
	go func() {
		for {
			var data *clientPacket
			select {
			case data = <-client.emitChan:
			case <-shutdown:
				return
			case <-done:
				return
			}
			if data.Op != opcode.Heartbeat {
				client.emitChan <- data // pass it along
				continue
			}

			client.receiveChan <- &discordPacket{
				Op: opcode.HeartbeatAck,
			}
		}
	}()

	// identify
	go func() {
		for {
			var data *clientPacket
			select {
			case data = <-client.emitChan:
			case <-shutdown:
				return
			case <-done:
				return
			}
			if data.Op != opcode.Identify {
				client.emitChan <- data
				continue
			}
			wgID.Done()

			client.receiveChan <- &discordPacket{
				Op:             opcode.DiscordEvent,
				SequenceNumber: seq,
				EventName:      event.Ready,
				Data:           []byte(`{}`),
			}
			seq++
		}
	}()

	// send hello packet
	wgID.Add(1)
	client.receiveChan <- &discordPacket{
		Op:   opcode.Hello,
		Data: []byte(`{"heartbeat_interval":45000,"_trace":["discord-gateway-prd-1-99"]}`),
	}
	wgID.Wait()

	// connection is established, now force a reconnect
	wg.Add(1) // only one, cause we only want one reconnect when 2 reconnect commands are received
	// TODO: it would be nicer if we could merge duplicate events in the socket layer to avoid timeouts instead.
	// this would also improve the behaviour for event handlers and channels, although removing duplicate events is
	// more advanced and requires more heavy work and memory.
	client.receiveChan <- &discordPacket{
		Op: opcode.Reconnect,
	}
	client.receiveChan <- &discordPacket{
		Op: opcode.Reconnect,
	}

	wg.Wait()
	if client.connected != 2 {
		t.Error("expected 2 connect, got 1")
	}

	// send hello packet
	client.receiveChan <- &discordPacket{
		Op:   opcode.Hello,
		Data: []byte(`{"heartbeat_interval":45000,"_trace":["discord-gateway-prd-1-99"]}`),
	}

	// resume
	for {
		var data *clientPacket
		select {
		case data = <-client.emitChan:
		case <-shutdown:
			return
		}
		if data.Op != opcode.Resume {
			client.emitChan <- data
			continue
		}
		break
	}

	client.receiveChan <- &discordPacket{
		Op:             opcode.DiscordEvent,
		EventName:      event.Resumed,
		Data:           []byte(`{}`),
		SequenceNumber: seq,
	}

	<-time.After(2 * time.Millisecond) // TODO: don't use timeouts
	if m.sequenceNumber != seq {
		t.Errorf("incorrect sequence number. Got %d, wants %d\n", m.sequenceNumber, seq)
		return
	}
	seq++

	// what if there is a session invalidate event
	wgID.Add(1)
	client.receiveChan <- &discordPacket{
		Op:   opcode.InvalidSession,
		Data: []byte(`false`),
	}

	// wait for identify
	wgID.Wait()
}
