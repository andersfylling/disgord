package disgord

import (
	"github.com/andersfylling/disgord/websocket"
	"github.com/andersfylling/disgord/websocket/event"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

type mockerWSReceiveOnly struct {
	reading chan []byte
}

func (g *mockerWSReceiveOnly) Open(endpoint string, requestHeader http.Header) (err error) {
	return
}

func (g *mockerWSReceiveOnly) WriteJSON(v interface{}) (err error) {
	return
}

func (g *mockerWSReceiveOnly) Close() (err error) {
	return
}

func (g *mockerWSReceiveOnly) Read() (packet []byte, err error) {
	packet = <-g.reading
	return
}

func (g *mockerWSReceiveOnly) Disconnected() bool {
	return true
}

var _ websocket.Conn = (*mockerWSReceiveOnly)(nil)

var sink1 int = 1

// BenchmarkDiscordEventToHandler from the time Disgord gets the raw byte event data, to the event handler is triggered
func Benchmark1000DiscordEventToHandler_cacheDisabled(b *testing.B) {
	mocker := mockerWSReceiveOnly{
		reading: make(chan []byte),
	}
	// starts receiver and operation handler
	wsClient, wsShutdownChan := websocket.NewTestClient(nil, &mocker)

	d := Client{
		shutdownChan: make(chan interface{}),
		config: &Config{
			DisableCache: true,
		},
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		ws:            wsClient,
		socketEvtChan: wsClient.EventChan(),
		evtDispatch:   NewDispatch(wsClient, false, 20),
	}
	go d.eventHandler()

	seq := uint(1)
	wg := &sync.WaitGroup{}
	d.On(event.Ready, func(s Session, evt *Ready) {
		sink1++
		wg.Done()
	})

	f := func(mocker *mockerWSReceiveOnly, wg *sync.WaitGroup, seq *uint) {
		loops := 1000
		wg.Add(loops)
		for i := 0; i < loops; i++ {
			//evt := []byte(`{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{}}`)
			evt := []byte(`{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{"v":6,"user_settings":{},"user":{"verified":true,"username":"Disgord tester","mfa_enabled":false,"id":"486832262592069632","email":null,"discriminator":"9338","bot":true,"avatar":null},"session_id":"d3954ff063fa8d387ec395fe65723624","relationships":[],"private_channels":[],"presences":[],"guilds":[{"unavailable":true,"id":"486833041486905345"},{"unavailable":true,"id":"486833611564253184"}],"_trace":["gateway-prd-main-kg6w","discord-sessions-prd-1-27"]}}`)
			mocker.reading <- evt
			*seq++
		}

		wg.Wait()
	}

	for i := 0; i < b.N; i++ {
		f(&mocker, wg, &seq)
	}

	close(d.shutdownChan)
	close(wsShutdownChan)
}
