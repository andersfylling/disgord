package gateway

import (
	"sync"

	"github.com/andersfylling/disgord/internal/gateway/opcode"
)

// inline
func newOnceChannels() onceChannels {
	return onceChannels{
		channels: map[opcode.OpCode]chan interface{}{},
	}
}

type onceChannels struct {
	mu       sync.Mutex
	channels map[opcode.OpCode]chan interface{}
}

func (o *onceChannels) Acquire(op opcode.OpCode) (ch chan interface{}) {
	var ok bool
	o.mu.Lock()
	if ch, ok = o.channels[op]; ok {
		delete(o.channels, op)
	}
	o.mu.Unlock()

	return ch
}

func (o *onceChannels) Add(op opcode.OpCode, ch chan interface{}) {
	o.mu.Lock()
	o.channels[op] = ch
	o.mu.Unlock()
}
