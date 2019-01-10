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

	"github.com/andersfylling/disgord/websocket"
)

func TestClient_Once(t *testing.T) {
	c, err := NewClient(&Config{
		Token: "testing",
	})
	if err != nil {
		panic(err)
	}

	dispatcher := c.evtDispatch
	if len(dispatcher.listenOnceOnly) > 0 {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", len(dispatcher.listenOnceOnly))
	}

	wg := sync.WaitGroup{}
	c.Once(EventMessageCreate, func() {
		wg.Done()
	})
	if len(dispatcher.listenOnceOnly) != 1 {
		t.Errorf("expected dispatch to have 1 listener. Got %d", len(dispatcher.listenOnceOnly))
	}
	wg.Add(1) // only run once

	// trigger the handler
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, nil)
	if len(dispatcher.listenOnceOnly) > 0 {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", len(dispatcher.listenOnceOnly))
	}

	// trigger the handler, again
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, nil)
	if len(dispatcher.listenOnceOnly) > 0 {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", len(dispatcher.listenOnceOnly))
	}

	wg.Wait()
	// if wg.Done() is called more than once, we get a panic.

	// TODO: add a timeout
}

func TestClient_On(t *testing.T) {
	c, err := NewClient(&Config{
		Token: "testing",
	})
	if err != nil {
		panic(err)
	}

	dispatcher := c.evtDispatch
	if len(dispatcher.listeners) > 0 {
		t.Errorf("expected dispatch to have 0 listeners. Got %d", len(dispatcher.listeners))
	}

	wg := sync.WaitGroup{}
	c.On(EventMessageCreate, func() {
		wg.Done()
	})
	if len(dispatcher.listeners) != 1 {
		t.Errorf("expected dispatch to have 1 listener. Got %d", len(dispatcher.listeners))
	}
	wg.Add(2)

	// trigger the handler twice
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, nil)
	dispatcher.triggerHandlers(nil, EventMessageCreate, c, nil)
	dispatcher.triggerHandlers(nil, EventReady, c, nil)
	wg.Wait()

	// TODO: add a timeout
}

// TestClient_System looks for crashes when the DisGord system starts up.
// the websocket logic is excluded to avoid crazy rewrites. At least, for now.
func TestClient_System(t *testing.T) {
	c, err := NewClient(&Config{
		Token: "testing",
	})
	if err != nil {
		panic(err)
	}

	input := make(chan *websocket.Event)
	c.ws = nil
	c.socketEvtChan = input
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

	_, err = c.GetGuild(244200618854580224)
	if err != nil {
		t.Error(err)
	}

	// cleanup
	close(c.evtDispatch.shutdown)
	close(c.shutdownChan)
}

func TestInternalStateHandlers(t *testing.T) {
	c, err := NewClient(&Config{
		Token: "testing",
	})
	if err != nil {
		t.Fatal(err)
	}

	id := Snowflake(123)

	if len(c.connectedGuilds) != 0 {
		t.Errorf("expected no guilds to have been added yet. Got %d, wants %d", len(c.connectedGuilds), 0)
	}

	c.handlerGuildCreate(c, &GuildCreate{
		Guild: NewPartialGuild(id),
	})
	if len(c.connectedGuilds) != 1 {
		t.Errorf("expected one guild to have been added. Got %d, wants %d", len(c.connectedGuilds), 1)
	}

	c.handlerGuildCreate(c, &GuildCreate{
		Guild: NewPartialGuild(id),
	})
	if len(c.connectedGuilds) != 1 {
		t.Errorf("Adding the same guild should not create another entry. Got %d, wants %d", len(c.connectedGuilds), 1)
	}

	c.handlerGuildDelete(c, &GuildDelete{
		UnavailableGuild: &GuildUnavailable{
			ID: 9999,
		},
	})
	if len(c.connectedGuilds) != 1 {
		t.Errorf("Removing a unknown guild should not affect the internal state. Got %d, wants %d", len(c.connectedGuilds), 1)
	}

	c.handlerGuildDelete(c, &GuildDelete{
		UnavailableGuild: &GuildUnavailable{
			ID: id,
		},
	})
	if len(c.connectedGuilds) != 0 {
		t.Errorf("Removing a connected guild should affect the internal state. Got %d, wants %d", len(c.connectedGuilds), 0)
	}
}
