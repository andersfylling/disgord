package websocket

import (
	"errors"
	"sync"
	"time"

	"github.com/andersfylling/disgord/logger"
)

type shardSync struct {
	timeoutMs time.Duration
	sync.Mutex
	next         time.Time
	logger       logger.Logger
	shutdownChan chan interface{}
}

func (s *shardSync) acquireTimestamp() time.Time {
	var delay time.Duration
	now := time.Now()

	s.Lock()
	if s.next.After(now) {
		delay = s.next.Sub(now)
	} else {
		delay = time.Duration(0)
		s.next = now
	}
	s.next = s.next.Add(s.timeoutMs)
	s.Unlock()

	return now.Add(delay)
}

func (s *shardSync) queueShard(shardID uint, cb func() error) error {
	now := time.Now()
	execTimestamp := s.acquireTimestamp()
	delay := execTimestamp.Sub(now)

	s.logger.Debug("shard", shardID, "will wait in connect queue for", delay)
	select {
	case <-time.After(delay):
		s.logger.Debug("shard", shardID, "waited", delay, "and is now being connected")
		return cb()
	case <-s.shutdownChan:
		s.logger.Debug("shard", shardID, "got shutdown signal while waiting in connect queue")
		return errors.New("shutting down")
	}
}
