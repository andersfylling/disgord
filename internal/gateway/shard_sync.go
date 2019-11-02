package gateway

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/logger"
	"github.com/andersfylling/disgord/internal/util"
)

type shardSync struct {
	timeoutMs time.Duration
	sync.Mutex
	queue        util.Queue
	logger       logger.Logger
	lpre         string
	shutdownChan chan interface{}
	metric       *IdentifyMetric
}

func (s *shardSync) queueShard(shardID uint, cb func() error) error {
	waitChan := make(chan (chan bool))
	s.queue.Push(waitChan)

	s.logger.Debug(s.lpre, "shard", shardID, "is waiting to identify")
	waited := time.Now()
	resultChan := <-waitChan

	s.logger.Debug(s.lpre, "shard", shardID, "waited", time.Since(waited), "and is now executing")
	if err := cb(); err != nil {
		resultChan <- false
		return err
	}
	resultChan <- true

	return nil
}

func (s *shardSync) process() {
	var timeout time.Duration
	resultChan := make(chan bool)

	for {
		select {
		case <-s.shutdownChan:
			s.logger.Debug("shard identify-rate-limiter got shutdown signal")
			return
		case <-time.After(timeout):
		}
		timeout = s.timeoutMs * time.Millisecond

		// 1000 identify / 24 hours rate limit check
		if s.metric.ReconnectsSince(24*time.Hour) > 999 {
			s.metric.Lock()
			oldest := s.metric.Reconnects[len(s.metric.Reconnects)-1000]
			s.metric.Unlock()

			timeout += (24 * time.Hour) - time.Since(oldest)
			continue // go back to the top to wait
		}

		x, err := s.queue.Pop()
		if err != nil {
			continue
		}
		if wChan, ok := x.(chan (chan bool)); ok {
			wChan <- resultChan
		} else {
			continue
		}

		// wait for result
		if reconnected := <-resultChan; !reconnected {
			timeout = 0 // no need to wait if the identify was _not_ sent
			continue
		}

		s.metric.Lock()
		s.metric.Reconnects = append(s.metric.Reconnects, time.Now())
		s.metric.Unlock()
	}
}
