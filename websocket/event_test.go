package websocket

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/websocket/opcode"
)

type testWS struct {
	closing      chan interface{}
	opening      chan interface{}
	writing      chan interface{}
	reading      chan []byte
	disconnected bool
	sync.Mutex
}

func (g *testWS) Open(endpoint string, requestHeader http.Header) (err error) {
	g.opening <- 1
	g.Lock()
	g.disconnected = false
	g.Unlock()
	return
}

func (g *testWS) WriteJSON(v interface{}) (err error) {
	g.writing <- v
	return
}

func (g *testWS) Close() (err error) {
	g.closing <- 1
	g.Lock()
	g.disconnected = true
	g.Unlock()
	return
}

func (g *testWS) Read() (packet []byte, err error) {
	packet = <-g.reading
	if packet == nil {
		err = errors.New("empty")
	}
	return
}

func (g *testWS) Disconnected() bool {
	return g.disconnected
}

var _ Conn = (*testWS)(nil)

func TestManager_RegisterEvent(t *testing.T) {
	m := EvtClient{
		trackedEvents: &UniqueStringSlice{},
	}
	t1 := "test"
	m.RegisterEvent(t1)

	if m.trackedEvents.Len() == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RegisterEvent(t1)
	if m.trackedEvents.Len() == 2 {
		t.Error("expected length to be 1, got 2")
	}
}

func TestManager_RemoveEvent(t *testing.T) {
	m := EvtClient{
		trackedEvents: &UniqueStringSlice{},
	}
	t1 := "test"
	m.RegisterEvent(t1)

	if m.trackedEvents.Len() == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RemoveEvent("sdfsdf")
	if m.trackedEvents.Len() == 0 {
		t.Error("expected length to be 1, got 0")
	}

	m.RemoveEvent(t1)
	if m.trackedEvents.Len() == 1 {
		t.Error("expected length to be 0, got 1")
	}
}

func TestManager_reconnect(t *testing.T) {
	conn := &testWS{
		closing:      make(chan interface{}),
		opening:      make(chan interface{}),
		writing:      make(chan interface{}),
		reading:      make(chan []byte),
		disconnected: true,
	}

	eChan := make(chan *Event)
	aChan := make(A)
	m := &EvtClient{
		trackedEvents: &UniqueStringSlice{},
		eventChan:     eChan,
		client: &client{
			conf: &Config{
				// identity
				Browser:             "disgord",
				Device:              "disgord",
				GuildLargeThreshold: 250,

				// lib specific
				Version:       constant.DiscordVersion,
				Encoding:      constant.JSONEncoding,
				ChannelBuffer: 3,
				Endpoint:      "sfkjsdlfsf",
				A:             aChan,

				// user settings
				BotToken: "sifhsdoifhsdifhsdf",
				HTTPClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			behaviors:    map[string]*behavior{},
			shutdown:     make(chan interface{}),
			restart:      make(chan interface{}),
			receiveChan:  make(chan *DiscordPacket),
			emitChan:     make(chan *clientPacket),
			conn:         conn,
			disconnected: true,
			ratelimit:    newRatelimiter(),
			a:            aChan,
			poolDiscordPkt: &sync.Pool{
				New: func() interface{} {
					return &DiscordPacket{}
				},
			},
		},
	}
	m.connectPermit = m.client
	seq := uint(1)

	shutdown := make(chan interface{})
	done := make(chan interface{})

	resume := 0
	identify := 1
	heartbeat := 2
	connecting := 3
	disconnecting := 4
	wg := []sync.WaitGroup{
		sync.WaitGroup{},
		sync.WaitGroup{},
		sync.WaitGroup{},
		sync.WaitGroup{},
		sync.WaitGroup{},
	}
	defer func() {
		wg[disconnecting].Add(1)
		m.Shutdown()
		close(done)
	}()

	// mocked DisGord logic (shard manager and event handler)
	go func() {
		for {
			select {
			case <-eChan:
				continue
			case b, ok := <-aChan:
				if !ok {
					continue
				}
				<-time.After(10 * time.Millisecond)
				releaser := make(B)
				b <- &K{
					Release: releaser,
					Key:     412, // random
				}
				<-releaser
			}
		}
	}()

	// mocked websocket server.. ish
	go func(seq *uint) {
		for {
			var data *clientPacket
			select {
			case v := <-conn.writing:
				data = v.(*clientPacket)
			case <-conn.opening:
				wg[connecting].Done()
				continue
			case <-conn.closing:
				wg[disconnecting].Done()
				continue
			case <-shutdown:
				return
			case <-done:
				return
			}
			switch data.Op {
			case opcode.EventHeartbeat:
				conn.reading <- []byte(`{"t":null,"s":null,"op":11,"d":null}`)
				wg[heartbeat].Done()
			case opcode.EventIdentify:
				conn.reading <- []byte(`{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{}}`)
				*seq++
				wg[identify].Done()
			case opcode.EventResume:
				conn.reading <- []byte(`{"t":"RESUMED","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{}}`)
				*seq++
				wg[resume].Done()
			default:
				// send the event back around
				fmt.Println("wtf")
			}
		}
	}(&seq)
	go func(t *testing.T) {
		select {
		case <-time.After(1 * time.Second):
		case <-done:
			return
		}
		close(shutdown)
		t.Error("timeout")
	}(t)

	wg[connecting].Add(1)
	m.Connect()
	wg[connecting].Wait()

	m.setupBehaviors()
	m.start()

	// send hello packet
	wg[heartbeat].Add(1)
	wg[identify].Add(1)
	conn.reading <- []byte(`{"t":null,"s":null,"op":10,"d":{"heartbeat_interval":45000,"_trace":["discord-gateway-prd-1-99"]}}`)
	wg[heartbeat].Wait()
	wg[identify].Wait()

	// connection is established, now force a reconnect
	wg[connecting].Add(1)
	wg[disconnecting].Add(1)
	conn.reading <- []byte(`{"t":null,"s":null,"op":7,"d":null}`)
	wg[disconnecting].Wait()
	wg[connecting].Wait()

	// send hello packet
	wg[resume].Add(1)
	wg[heartbeat].Add(1)
	conn.reading <- []byte(`{"t":null,"s":null,"op":10,"d":{"heartbeat_interval":45000,"_trace":["discord-gateway-prd-1-99"]}}`)
	wg[resume].Wait()
	wg[heartbeat].Wait()

	// during testing, most timeouts are 0, so we experience moments where not all
	// channels have finished syncing. TODO: remove timeout requirement.
	<-time.After(time.Millisecond * 10)
	m.RLock()
	sequence := m.sequenceNumber
	m.RUnlock()
	if sequence != seq-1 {
		t.Errorf("incorrect sequence number. Got %d, wants %d\n", sequence, seq)
		return
	}
	seq++

	// what if there is a session invalidate event
	wg[identify].Add(1)
	conn.reading <- []byte(`{"t":null,"s":null,"op":9,"d":false}`)

	// wait for identify
	wg[identify].Wait()
}
