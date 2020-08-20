// +build !integration

package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/json"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/andersfylling/disgord/internal/constant"
	"github.com/andersfylling/disgord/internal/gateway/cmd"
	"github.com/andersfylling/disgord/internal/gateway/opcode"
	"github.com/andersfylling/disgord/internal/logger"
)

type testWS struct {
	closing     chan interface{}
	opening     chan interface{}
	writing     chan interface{}
	reading     chan []byte
	isConnected atomic.Bool
}

func (g *testWS) Open(ctx context.Context, endpoint string, requestHeader http.Header) (err error) {
	g.opening <- 1
	g.isConnected.Store(true)
	return
}

func (g *testWS) WriteJSON(v interface{}) (err error) {
	g.writing <- v
	return
}

func (g *testWS) Close() (err error) {
	g.closing <- 1
	g.isConnected.Store(false)
	return
}

func (g *testWS) Read(ctx context.Context) (packet []byte, err error) {
loop:
	for {
		select {
		case packet = <-g.reading:
		case <-ctx.Done():
			break loop
		case <-time.After(1 * time.Millisecond):
			if g.isConnected.Load() {
				continue
			}
		}
		break
	}

	if packet == nil {
		err = errors.New("empty")
	}
	return
}

func (g *testWS) Disconnected() bool {
	return !g.isConnected.Load()
}

var _ Conn = (*testWS)(nil)

func TestEvtIdentify(t *testing.T) {
	i := &evtIdentity{}
	var fields map[string]interface{}

	raw, err := json.Marshal(i)
	if err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Error(err)
	}

	if constant.DiscordVersion == 6 {
		if _, ok := fields["intents"]; ok {
			t.Error("discord gateway 6 states that intents are optional")
			// Don't send out intents if none were specified
		}
	} else if constant.DiscordVersion == 7 {
		if _, ok := fields["intents"]; !ok {
			t.Error("discord gateway 7 states that intents are mandatory")
			// https://discord.com/developers/docs/topics/gateway#gateway-intents
		}
	}
}

// TODO: rewrite. EventClient now waits for a Ready event in the Connect method
func TestEvtClient_communication(t *testing.T) {
	deadline := 1 * time.Second
	conn := &testWS{
		closing: make(chan interface{}),
		opening: make(chan interface{}),
		writing: make(chan interface{}),
		reading: make(chan []byte),
	}
	conn.isConnected.Store(false)

	eChan := make(chan *Event)

	shutdown := make(chan interface{})
	done := make(chan interface{})

	m, err := NewEventClient(0, &EvtConfig{
		// identity
		Browser:             "disgord",
		Device:              "disgord",
		GuildLargeThreshold: 250,

		// lib specific
		Endpoint: "sfkjsdlfsf",
		Version:  constant.DiscordVersion,
		Encoding: constant.JSONEncoding,
		Logger:   &logger.FmtPrinter{},

		// user settings
		BotToken: "sifhsdoifhsdifhsdf",
		DiscordPktPool: &sync.Pool{
			New: func() interface{} {
				return &DiscordPacket{}
			},
		},

		connectQueue: func(shardID uint, cb func() error) error {
			<-time.After(time.Duration(10) * time.Millisecond)
			return cb()
		},

		// injected for testing
		EventChan: eChan,
		conn:      conn,

		SystemShutdown: shutdown,
	})
	if err != nil {
		t.Fatal(err)
	}
	m.timeoutMultiplier = 0
	seq := uint32(1)

	// ###############################
	// RECONNECT
	resume := 0
	identify := 1
	heartbeat := 2
	connecting := 3
	disconnecting := 4
	status := 4
	wg := []sync.WaitGroup{
		sync.WaitGroup{},
		sync.WaitGroup{},
		sync.WaitGroup{},
		sync.WaitGroup{},
		sync.WaitGroup{},
		sync.WaitGroup{},
	}
	defer func() {
		wg[disconnecting].Add(1)
		close(done)
	}()

	// mocked Disgord logic (shard manager and event handler)
	go func() {
		for {
			select {
			case <-eChan:
			case <-shutdown:
				return
			}
		}
	}()

	// mocked websocket server.. ish
	go func(seq *uint32) {
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
				var d string = `{"t":null,"s":null,"op":11,"d":null}`
				conn.reading <- []byte(d)
				//fmt.Printf("discord: ->%+v\n", d)
				wg[heartbeat].Done()
			case opcode.EventIdentify:
				var d string = `{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{}}`
				conn.reading <- []byte(d)
				//fmt.Printf("discord: ->%+v\n", d)
				*seq++
				wg[identify].Done()
			case opcode.EventResume:
				var d string = `{"t":"RESUMED","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{}}`
				conn.reading <- []byte(d)
				//fmt.Printf("discord: ->%+v\n", d)
				*seq++
				wg[resume].Done()
			case opcode.EventStatusUpdate:
				wg[status].Done()
			default:
				// send the event back around
				fmt.Println("wtf")
			}
		}
	}(&seq)
	go func(t *testing.T) {
		select {
		case <-time.After(deadline):
		case <-done:
			return
		}
		close(shutdown)
		t.Error("timeout")
	}(t)

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		// send hello packet
		wg[heartbeat].Add(1)
		wg[identify].Add(1)
		conn.reading <- []byte(`{"t":null,"s":null,"op":10,"d":{"heartbeat_interval":45000,"_trace":["discord-gateway-prd-1-99"]}}`)
		wg[heartbeat].Wait()
		wg[identify].Wait()
		wg2.Done()
	}()

	wg[connecting].Add(1)
	_ = m.Connect()
	wg2.Wait()
	wg[connecting].Wait()

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
	sequence := m.sequenceNumber.Load()
	if sequence != seq-1 {
		t.Errorf("incorrect sequence number. Got %d, wants %d\n", sequence, seq)
		return
	}
	seq++

	// what if there is a session invalidate event
	wg[identify].Add(1)
	seq = 1
	conn.reading <- []byte(`{"t":null,"s":null,"op":9,"d":false}`)

	// wait for identify
	wg[identify].Wait()

	// #########################################
	// emitting user messages
	wg[status].Add(1)
	_ = m.emit(cmd.UpdateStatus, 1)
	wg[status].Wait()

	<-time.After(10 * time.Millisecond)
}
