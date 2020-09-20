// +build !integration

package disgord

import (
	"context"
	"sync"
	"testing"
)

func Test_isHandler(t *testing.T) {
	handler := make(chan *MessageCreate)
	if !isHandler(handler) {
		t.Error("chan *Message create was identified as not a handler")
	}

	handlers := []interface{}{
		make(chan *MessageCreate),
		make(chan *MessageCreate, 20),
		make(chan *Ready, 20),
		func() {},
		func(s Session) {},
		func(s Session, e *MessageCreate) {},
	}
	for i := range handlers {
		if !isHandler(handlers[i]) {
			t.Error("identified as not a handler")
		}
	}
}

func TestRegister(t *testing.T) {
	d := newDispatcher()
	handler := make(chan *MessageCreate)

	if err := d.register(EvtMessageCreate, handler); err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-handler
		wg.Done()
	}()
	d.dispatch(context.Background(), EvtMessageCreate, &MessageCreate{})
	wg.Wait()
}

func TestCtrl_CloseChannel(t *testing.T) {
	d := newDispatcher()
	handler := make(chan *MessageCreate)
	ctrl := &Ctrl{Channel: handler}

	if err := d.register(EvtMessageCreate, handler, ctrl); err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-handler
		wg.Done()
	}()
	d.dispatch(context.Background(), EvtMessageCreate, &MessageCreate{})
	wg.Wait()

	// close channel
	ctrl.CloseChannel()
	if _, open := <-handler; open {
		t.Fatal("expected channel to be closed")
	}

	// should not hang
	d.dispatch(context.Background(), EvtMessageCreate, &MessageCreate{})
}
