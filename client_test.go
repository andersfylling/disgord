package disgord

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andersfylling/disgord/websocket"
)

//////////////////////////////////////////////////////
//
// Struct extensions / extra methods for testing only
//
//////////////////////////////////////////////////////

func (d *evtDemultiplexer) nrOfAliveHandlers() (counter int) {
	for k := range d.handlers {
		for i := range d.handlers[k] {
			if d.handlers[k][i].ctrl.IsDead() == false {
				counter++
			}
		}
	}

	return
}

//////////////////////////////////////////////////////
//
// Benchmarks
//
//////////////////////////////////////////////////////

func BenchmarkClient_On(b *testing.B) {
	c := New(&Config{
		BotToken:     "testing",
		DisableCache: true,
	})
	c.shardManager.evtChan = make(chan *websocket.Event)
	c.setupConnectEnv()

	msgData := []byte(`{"attachments":[],"author":{"avatar":"69a7a0e9cb963adfdd69a2224b4ac180","discriminator":"7237","id":"228846961774559232","username":"Anders"},"channel_id":"409359688258551850","content":"https://discord.gg/kaWJsV","edited_timestamp":null,"embeds":[],"id":"409654019611688960","mention_everyone":false,"mention_roles":[],"mentions":[],"nonce":"409653919891849216","pinned":false,"timestamp":"2018-02-04T10:18:49.279000+00:00","tts":false,"type":0}`)
	evt := &websocket.Event{Name: EventMessageCreate, Data: msgData}

	wg := sync.WaitGroup{}
	c.On(EventMessageCreate, func() {
		wg.Done()
	})

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		c.shardManager.evtChan <- evt
		wg.Wait()
	}

}

//////////////////////////////////////////////////////
//
// TEST funcs
//
//////////////////////////////////////////////////////

func TestClient_Once(t *testing.T) {
	c, err := NewClient(&Config{
		BotToken:     "testing",
		DisableCache: true,
	})
	if err != nil {
		panic(err)
	}

	dispatcher := c.evtDemultiplexer
	base := dispatcher.nrOfAliveHandlers()
	if dispatcher.nrOfAliveHandlers() > base {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", dispatcher.nrOfAliveHandlers())
	}

	wg := sync.WaitGroup{}
	c.On(EventMessageCreate, func() {
		wg.Done()
	}, &Ctrl{Remaining: 1})
	if dispatcher.nrOfAliveHandlers() != 1+base {
		t.Errorf("expected dispatch to have 1 listener. Got %d", dispatcher.nrOfAliveHandlers())
	}
	wg.Add(1) // only run once

	// trigger the handler
	box := &MessageCreate{}
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, box)
	if dispatcher.nrOfAliveHandlers() > 0+base {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", dispatcher.nrOfAliveHandlers())
	}

	// trigger the handler, again
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, box)
	if dispatcher.nrOfAliveHandlers() > 0+base {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", dispatcher.nrOfAliveHandlers())
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
	c := New(&Config{
		BotToken:     "testing",
		DisableCache: true,
	})

	dispatcher := c.evtDemultiplexer
	base := dispatcher.nrOfAliveHandlers()
	if dispatcher.nrOfAliveHandlers() > 0+base {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", dispatcher.nrOfAliveHandlers())
	}

	wg := sync.WaitGroup{}
	c.On(EventMessageCreate, func() {
		wg.Done()
	})
	if dispatcher.nrOfAliveHandlers() != 1+base {
		t.Errorf("expected dispatch to have 1 listener. Got %d", dispatcher.nrOfAliveHandlers())
	}
	wg.Add(2)

	// trigger the handler twice
	evt := &MessageCreate{}
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, evt)
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, evt)
	dispatcher.triggerHandlers(nil, EventReady, c, evt)
	wg.Wait()
}

func TestClient_On_Middleware(t *testing.T) {
	c := New(&Config{
		BotToken:     "testing",
		DisableCache: true,
	})
	dispatcher := c.evtDemultiplexer

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
	c.On(EventMessageCreate, func() {
		wg.Done()
	})
	c.On(EventMessageCreate, mdlwHasBotPrefix, func() {
		wg.Done()
	})
	c.On(EventMessageCreate, mdlwHasDifferentPrefix, func() {
		wg.Done()
	})
	wg.Add(2)

	// trigger the handler twice
	evt := &MessageCreate{
		Message: &Message{Content: prefix + " testing"},
	}
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, evt)
	dispatcher.triggerHandlers(nil, EventReady, c, evt)
	wg.Wait()
}

// TestClient_System looks for crashes when the DisGord system starts up.
// the websocket logic is excluded to avoid crazy rewrites. At least, for now.
func TestClient_System(t *testing.T) {
	c, err := NewClient(&Config{
		BotToken: "testing",
	})
	if err != nil {
		panic(err)
	}

	c.shardManager.evtChan = make(chan *websocket.Event, 1)
	c.shardManager.shards = append(c.shardManager.shards, &WSShard{})
	input := c.shardManager.evtChan
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
		err = json.Unmarshal(data, p)
		if err != nil {
			t.Fatal(err)
		}

		// ignore non-event-type packets
		if p.E == "" {
			continue
		}

		input <- &websocket.Event{
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
	close(c.evtDemultiplexer.shutdown)
	close(c.shutdownChan)
}

func TestInternalStateHandlers(t *testing.T) {
	c, err := NewClient(&Config{
		BotToken: "testing",
	})
	if err != nil {
		t.Fatal(err)
	}
	c.shardManager.shards = append(c.shardManager.shards, &WSShard{}) // don't remove

	id := Snowflake(123)

	if len(c.GetConnectedGuilds()) != 0 {
		t.Errorf("expected no guilds to have been added yet. Got %d, wants %d", len(c.GetConnectedGuilds()), 0)
	}

	c.handlerAddToConnectedGuilds(c, &GuildCreate{
		Guild: NewPartialGuild(id),
	})
	if len(c.GetConnectedGuilds()) != 1 {
		t.Errorf("expected one guild to have been added. Got %d, wants %d", len(c.GetConnectedGuilds()), 1)
	}

	c.handlerAddToConnectedGuilds(c, &GuildCreate{
		Guild: NewPartialGuild(id),
	})
	if len(c.GetConnectedGuilds()) != 1 {
		t.Errorf("Adding the same guild should not create another entry. Got %d, wants %d", len(c.GetConnectedGuilds()), 1)
	}

	c.handlerRemoveFromConnectedGuilds(c, &GuildDelete{
		UnavailableGuild: &GuildUnavailable{
			ID: 9999,
		},
	})
	if len(c.GetConnectedGuilds()) != 1 {
		t.Errorf("Removing a unknown guild should not affect the internal state. Got %d, wants %d", len(c.GetConnectedGuilds()), 1)
	}

	c.handlerRemoveFromConnectedGuilds(c, &GuildDelete{
		UnavailableGuild: &GuildUnavailable{
			ID: id,
		},
	})
	if len(c.GetConnectedGuilds()) != 0 {
		t.Errorf("Removing a connected guild should affect the internal state. Got %d, wants %d", len(c.GetConnectedGuilds()), 0)
	}
}
