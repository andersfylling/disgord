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

func (s *shardSync) queueShard(shardID uint, cb func() error) error {
	var delay time.Duration
	now := time.Now()
	start := now

	s.Lock()
	defer s.Unlock()

	if s.next.After(now) {
		delay = s.next.Sub(now)
	} else {
		delay = time.Duration(0)
		s.next = now
	}
	s.next = s.next.Add(s.timeoutMs)

	s.logger.Debug("shard", shardID, "will wait in connect queue for", delay)
	select {
	case <-time.After(delay):
		s.logger.Debug("shard", shardID, "waited", delay, "and is now being connected")
		if err := cb(); err != nil {
			return err
		}
		execDuration := time.Now().Sub(start)
		s.next = s.next.Add(execDuration)

	case <-s.shutdownChan:
		s.logger.Debug("shard", shardID, "got shutdown signal while waiting in connect queue")
		return errors.New("shutting down")
	}

	return nil
}
