package websocket

import (
	"sync"
)

// inline
func newOnceChannels() onceChannels {
	return onceChannels{
		channels: map[uint]chan interface{}{},
	}
}

type onceChannels struct {
	mu       sync.Mutex
	channels map[uint]chan interface{}
}

func (o *onceChannels) Acquire(op uint) (ch chan interface{}) {
	var ok bool
	o.mu.Lock()
	if ch, ok = o.channels[op]; ok {
		delete(o.channels, op)
	}
	o.mu.Unlock()

	return ch
}

func (o *onceChannels) Add(op uint, ch chan interface{}) {
	o.mu.Lock()
	o.channels[op] = ch
	o.mu.Unlock()
}
