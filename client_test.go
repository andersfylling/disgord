// +build !integration

package disgord

import (
	"github.com/andersfylling/disgord/internal/logger"
	"github.com/andersfylling/disgord/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andersfylling/disgord/internal/gateway"
)

//////////////////////////////////////////////////////
//
// Struct extensions / extra methods for testing only
//
//////////////////////////////////////////////////////

func (d *dispatcher) nrOfAliveHandlers() (counter int) {
	d.RLock()
	defer d.RUnlock()

	for k := range d.handlerSpecs {
		for i := range d.handlerSpecs[k] {
			d.handlerSpecs[k][i].Lock()
			if d.handlerSpecs[k][i].ctrl.IsDead() == false {
				counter++
			}
			d.handlerSpecs[k][i].Unlock()
		}
	}

	return
}

func ensure(inputs ...interface{}) {
	for i := range inputs {
		if err, ok := inputs[i].(error); ok && err != nil {
			panic(err)
		}
	}
}

//////////////////////////////////////////////////////
//
// Tests
//
//////////////////////////////////////////////////////

func TestOn(t *testing.T) {
	c := New(Config{
		BotToken:     "sdkjfhdksfhskdjfhdkfjsd",
		DisableCache: true,
	})

	t.Run("normal Session", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("should not have triggered a panic")
			}
		}()

		c.On(EvtChannelCreate, func(s Session, e *ChannelCreate) {})
	})

	t.Run("normal Session with ctrl", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("should not have triggered a panic")
			}
		}()

		c.On(EvtChannelCreate, func(s Session, e *ChannelCreate) {}, &Ctrl{Runs: 1})
	})

	t.Run("normal Session with multiple ctrl's", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("multiple controllers should trigger a panic")
			}
		}()

		c.On(EvtChannelCreate,
			func(s Session, e *ChannelCreate) {},
			&Ctrl{Runs: 1},
			&Ctrl{Until: time.Now().Add(1 * time.Minute)})
	})

	t.Run("Session pointer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Did not panic on incorrect handler signature")
			}
		}()

		c.On(EvtChannelCreate, func(s *Session, e *ChannelCreate) {})
	})
}

//////////////////////////////////////////////////////
//
// Benchmarks
//
//////////////////////////////////////////////////////

func BenchmarkClient_On(b *testing.B) {
	b.ReportAllocs()
	c := New(Config{
		BotToken:     "testing",
		DisableCache: true,
	})
	c.eventChan = make(chan *gateway.Event)
	c.setupConnectEnv()

	msgData := []byte(`{"attachments":[],"author":{"avatar":"69a7a0e9cb963adfdd69a2224b4ac180","discriminator":"7237","id":"228846961774559232","username":"Anders"},"channel_id":"409359688258551850","content":"https://discord.gg/kaWJsV","edited_timestamp":null,"embeds":[],"id":"409654019611688960","mention_everyone":false,"mention_roles":[],"mentions":[],"nonce":"409653919891849216","pinned":false,"timestamp":"2018-02-04T10:18:49.279000+00:00","tts":false,"type":0}`)

	wg := sync.WaitGroup{}
	c.On(EvtMessageCreate, func() {
		wg.Done()
	})

	for i := 0; i < b.N; i++ {
		wg.Add(1)

		cp := make([]byte, len(msgData))
		copy(cp, msgData)
		evt := &gateway.Event{Name: EvtMessageCreate, Data: cp}
		c.eventChan <- evt
		wg.Wait()
	}
}

//////////////////////////////////////////////////////
//
// TEST funcs
//
//////////////////////////////////////////////////////

func TestClient_Once(t *testing.T) {
	c := New(Config{
		BotToken:     "testing",
		DisableCache: true,
		Logger:       &logger.FmtPrinter{},
	})
	defer close(c.dispatcher.shutdown)

	dispatcher := c.dispatcher
	input := make(chan *gateway.Event)
	go c.demultiplexer(dispatcher, input)

	trigger := func() {
		input <- &gateway.Event{Name: EvtMessageCreate, Data: []byte(`{"content":"testing"}`)}
	}

	base := dispatcher.nrOfAliveHandlers()
	wg := sync.WaitGroup{}
	c.On(EvtMessageCreate, func() {
		wg.Done()
	}, &Ctrl{Runs: 1})
	got := dispatcher.nrOfAliveHandlers() - base
	if got != 1 {
		t.Errorf("expected dispatch to have 1 listener. Got %d", got)
	}
	wg.Add(1) // only run once

	// make sure the handler is called
	trigger()
	<-time.After(100 * time.Millisecond)
	got = dispatcher.nrOfAliveHandlers() - base
	if got > 0 {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", got)
	}

	// make sure it is not called once more
	trigger()
	<-time.After(100 * time.Millisecond)
	got = dispatcher.nrOfAliveHandlers() - base
	if got > 0 {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", got)
	}

	done := make(chan interface{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-time.After(20 * time.Millisecond):
		t.Fail()
	case <-done:
	}
}

func TestClient_On(t *testing.T) {
	c := New(Config{
		BotToken:     "testing",
		DisableCache: true,
		Cache:        &CacheNop{},
	})
	defer close(c.dispatcher.shutdown)

	dispatcher := c.dispatcher
	input := make(chan *gateway.Event)
	go c.demultiplexer(dispatcher, input)

	base := dispatcher.nrOfAliveHandlers()
	if dispatcher.nrOfAliveHandlers() > 0+base {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", dispatcher.nrOfAliveHandlers())
	}

	wg := sync.WaitGroup{}
	c.On(EvtMessageCreate, func() {
		wg.Done()
	})
	if dispatcher.nrOfAliveHandlers() != 1+base {
		t.Errorf("expected dispatch to have 1 listener. Got %d", dispatcher.nrOfAliveHandlers())
	}
	wg.Add(2)

	// trigger the handler twice
	input <- &gateway.Event{Name: EvtMessageCreate, Data: []byte(`{}`)}
	input <- &gateway.Event{Name: EvtMessageCreate, Data: []byte(`{}`)}
	input <- &gateway.Event{Name: EvtReady, Data: []byte(`{}`)}
	wg.Wait()
}

func TestClient_On_Middleware(t *testing.T) {
	c := New(Config{
		BotToken:     "testing",
		DisableCache: true,
		Cache:        &CacheNop{},
	})
	defer close(c.dispatcher.shutdown)
	dispatcher := c.dispatcher
	input := make(chan *gateway.Event)
	go c.demultiplexer(dispatcher, input)

	const prefix = "this cool prefix"
	var mdlwHasBotPrefix Middleware = func(evt interface{}) interface{} {
		msg := (evt.(*MessageCreate)).Message
		if strings.HasPrefix(msg.Content, prefix) {
			return evt
		}

		return nil
	}
	var mdlwHasDifferentPrefix Middleware = func(evt interface{}) interface{} {
		msg := (evt.(*MessageCreate)).Message
		if strings.HasPrefix(msg.Content, "random unknown prefix") {
			return evt
		}

		return nil
	}

	wg := sync.WaitGroup{}
	c.On(EvtMessageCreate, func() {
		wg.Done()
	})
	c.On(EvtMessageCreate, mdlwHasBotPrefix, func() {
		wg.Done()
	})
	c.On(EvtMessageCreate, mdlwHasDifferentPrefix, func() {
		wg.Done()
	})
	wg.Add(2)

	input <- &gateway.Event{Name: EvtMessageCreate, Data: []byte(`{"content":"` + prefix + ` testing"}`)}
	input <- &gateway.Event{Name: EvtReady, Data: []byte(`{"content":"testing"}`)}
	wg.Wait()
}

// TestClient_System looks for crashes when the Disgord system starts up.
// the websocket logic is excluded to avoid crazy rewrites. At least, for now.
func TestClient_System(t *testing.T) {
	c, err := NewClient(Config{
		BotToken: "testing",
	})
	if err != nil {
		panic(err)
	}

	input := make(chan *gateway.Event, 1)
	c.eventChan = input
	c.setupConnectEnv()

	var files []string

	root := "testdata/phases/startup-smooth-1"
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := range files {
		if files[i] == root {
			files[i] = files[len(files)-1]
			files = files[:len(files)-1]
			break
		}
	}
	sort.Slice(files, func(i, j int) bool {
		starti := strings.Split(files[i][len(root+"/"):], "_")
		startj := strings.Split(files[j][len(root+"/"):], "_")

		if _, err := strconv.Atoi(starti[0]); err != nil {
			t.Fatal(err)
		}
		if _, err := strconv.Atoi(startj[0]); err != nil {
			t.Fatal(err)
		}

		a, _ := strconv.Atoi(starti[0])
		b, _ := strconv.Atoi(startj[0])

		return a < b
	})
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		p := &struct {
			E string          `json:"t"`
			D json.RawMessage `json:"d"`
		}{}
		if err = json.Unmarshal(data, p); err != nil {
			t.Fatal(err)
		}
		executeInternalUpdater(p)

		// ignore non-event-type packets
		if p.E == "" {
			continue
		}

		input <- &gateway.Event{
			Name: p.E,
			Data: p.D,
		}
	}

	// TODO: race - / don't have another way to "sync" the go routines
	//if _, err = c.cache.GetGuild(244200618854580224); err != nil {
	//	t.Error(err)
	//}

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//c.On(event.GuildMembersChunk, func(s Session, evt *GuildMembersChunk) {
	//	var msg string
	//	for i := range evt.Members {
	//		if evt.Members[i].User == nil {
	//			msg = fmt.Sprintf("expected user in member to not be nil. Got %+v", evt.Members[i])
	//			break
	//		}
	//	}
	//
	//	if msg != "" {
	//		t.Error(msg)
	//	}
	//	wg.Done()
	//})
	//wg.Wait()

	// cleanup
	close(c.dispatcher.shutdown)
	close(c.shutdownChan)
}

func TestInternalStateHandlers(t *testing.T) {
	c, err := NewClient(Config{
		BotToken: "testing",
	})
	if err != nil {
		t.Fatal(err)
	}

	id := Snowflake(123)

	if len(c.GetConnectedGuilds()) != 0 {
		t.Errorf("expected no Guilds to have been added yet. Got %d, wants %d", len(c.GetConnectedGuilds()), 0)
	}

	c.handlers.saveGuildID(c, &GuildCreate{
		Guild: NewPartialGuild(id),
	})
	if len(c.GetConnectedGuilds()) != 1 {
		t.Errorf("expected one guild to have been added. Got %d, wants %d", len(c.GetConnectedGuilds()), 1)
	}

	c.handlers.saveGuildID(c, &GuildCreate{
		Guild: NewPartialGuild(id),
	})
	if len(c.GetConnectedGuilds()) != 1 {
		t.Errorf("Adding the same guild should not create another entry. Got %d, wants %d", len(c.GetConnectedGuilds()), 1)
	}

	c.handlers.deleteGuildID(c, &GuildDelete{
		UnavailableGuild: &GuildUnavailable{
			ID: 9999,
		},
	})
	if len(c.GetConnectedGuilds()) != 1 {
		t.Errorf("Removing a unknown guild should not affect the internal state. Got %d, wants %d", len(c.GetConnectedGuilds()), 1)
	}

	c.handlers.deleteGuildID(c, &GuildDelete{
		UnavailableGuild: &GuildUnavailable{
			ID: id,
		},
	})
	if len(c.GetConnectedGuilds()) != 0 {
		t.Errorf("Removing a connected guild should affect the internal state. Got %d, wants %d", len(c.GetConnectedGuilds()), 0)
	}
}
