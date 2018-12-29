package disgord

import (
	"sync"
	"testing"
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
