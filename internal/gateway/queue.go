package gateway

import (
	"errors"
	"sync"

	"github.com/andersfylling/disgord/internal/gateway/opcode"
)

func newClientPktQueue(limit int) clientPktQueue {
	if limit == 0 {
		limit = -1 // no limit
	}
	return clientPktQueue{
		limit: limit,
	}
}

// clientPktQueue is an ordered queue. Entries are not removed unless they are successfully written to the websocket.
type clientPktQueue struct {
	sync.RWMutex
	messages []*clientPacket
	limit    int
}

func (c *clientPktQueue) IsEmpty() bool {
	c.RLock()
	defer c.RUnlock()

	return len(c.messages) == 0
}

func (c *clientPktQueue) AddByOverwrite(msg *clientPacket) error {
	c.Lock()
	defer c.Unlock()

	for i := range c.messages {
		if c.messages[i].Op == msg.Op {
			c.messages[i] = msg
			return nil
		}
	}
	return errors.New("no entry with existing operation code")
}

func (c *clientPktQueue) Add(msg *clientPacket) error {
	if msg.Op == opcode.EventStatusUpdate {
		if err := c.AddByOverwrite(msg); err == nil {
			return nil
		}
	}

	c.Lock()
	defer c.Unlock()
	if len(c.messages) == c.limit {
		return errors.New("can not send anymore messages, queue is full")
	}

	c.messages = append(c.messages, msg)
	return nil
}

func (c *clientPktQueue) Try(cb func(msg *clientPacket) error) error {
	c.Lock()
	defer c.Unlock()
	if len(c.messages) == 0 {
		return nil // nothing to try, this avoid a potential race as well
	}

	next := c.messages[0]
	if err := cb(next); err != nil {
		return err
	}

	// shift to avoid re-allocations
	for i := 0; i < len(c.messages)-1; i++ {
		c.messages[i] = c.messages[i+1]
	}
	c.messages = c.messages[:len(c.messages)-1]
	return nil
}

func (c *clientPktQueue) Steal() (m []*clientPacket) {
	c.Lock()
	defer c.Unlock()

	m = make([]*clientPacket, len(c.messages))
	copy(m, c.messages)
	for i := range c.messages {
		c.messages[i] = nil // redundant?
	}
	c.messages = c.messages[:0]
	return m
}
