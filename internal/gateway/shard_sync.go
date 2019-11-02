package gateway

import (
	"errors"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/logger"
)

type shardSync struct {
	timeoutMs time.Duration
	sync.Mutex
	next         time.Time
	logger       logger.Logger
	shutdownChan chan interface{}
	metric       *ShardMetric
}

func (s *shardSync) queueShard(shardID uint, cb func() error) error {
	var success bool
	var delay time.Duration
	now := time.Now()

	s.metric.Lock()
	s.metric.RequestedReconnect = append(s.metric.RequestedReconnect, now)
	s.metric.Unlock()

	defer func() {
		if !success {
			return
		}
		s.metric.Lock()
		s.metric.Reconnects = append(s.metric.Reconnects, time.Now())
		s.metric.Unlock()
	}()

	s.Lock()
	defer s.Unlock()

	if s.next.After(now) {
		delay = s.next.Sub(now)
	} else {
		delay = time.Duration(0)
		s.next = now
	}
	s.next = s.next.Add(s.timeoutMs)

	// 1000 identify / 24 hours rate limit check
	if s.metric.ReconnectsSince(24*time.Hour) > 999 {
		s.metric.Lock()
		oldest := s.metric.Reconnects[len(s.metric.Reconnects)-1000]
		s.metric.Unlock()

		delay += (24 * time.Hour) - time.Since(oldest)
		s.next = s.next.Add(delay) // might add excess milliseconds, but it really doesn't matter at this stage
	}

	s.logger.Debug("shard", shardID, "will wait in connect queue for", delay)
	select {
	case <-time.After(delay):
		s.logger.Debug("shard", shardID, "waited", delay, "and is now being connected")
		start := time.Now()
		if err := cb(); err != nil {
			return err
		}
		execDuration := time.Since(start)
		s.next = s.next.Add(execDuration)
		success = true // store reconnect timestamp in metrics

	case <-s.shutdownChan:
		s.logger.Debug("shard", shardID, "got shutdown signal while waiting in connect queue")
		return errors.New("shutting down")
	}

	return nil
}
