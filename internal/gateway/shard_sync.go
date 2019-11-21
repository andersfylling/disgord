package gateway

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/logger"
)

const DefaultIdentifyRateLimit = 1000

func newShardSync(conf *ShardConfig, l logger.Logger, lPrefix string, shutdownChan chan interface{}) *shardSync {
	return &shardSync{
		identifiesPer24H: conf.IdentifiesPer24H,
		timeout:          conf.ShardRateLimit,
		queue:            make(chan *shardSyncQueueItem, 100), // it's just pointers anyways
		logger:           l,
		lpre:             lPrefix,
		shutdownChan:     shutdownChan,
		metric:           &IdentifyMetric{},
	}
}

type shardSyncQueueItem struct {
	ShardID uint
	run     func() error
	errChan chan error
}

type shardSync struct {
	sync.Mutex

	identifiesPer24H uint
	timeout          time.Duration
	queue            chan *shardSyncQueueItem
	logger           logger.Logger
	lpre             string
	shutdownChan     chan interface{}
	metric           *IdentifyMetric
}

func (s *shardSync) queueShard(shardID uint, cb func() error) (err error) {
	errChan := make(chan error)
	defer func() {
		close(errChan)
	}()

	start := time.Now()

	s.logger.Debug(s.lpre, "shard", shardID, "is waiting to identify")
	s.queue <- &shardSyncQueueItem{
		ShardID: shardID,
		run:     cb,
		errChan: errChan,
	} // TODO: what if this becomes blocking?

	select {
	case <-s.shutdownChan:
		return nil
	case err = <-errChan:
	}
	s.logger.Debug(s.lpre, "shard", shardID, "waited and finished execution after", time.Since(start))
	return err
}

func (s *shardSync) process() {
	for {
		var item *shardSyncQueueItem
		var open bool
		var penalty time.Duration

		select {
		case <-s.shutdownChan:
			s.logger.Debug(s.lpre, "shard identify-rate-limiter got shutdown signal")
			return
		case item, open = <-s.queue:
			if !open {
				s.logger.Error(s.lpre, "queue unexpectly closed - shards can no longer identify")
				return
			}
		}
		if item == nil {
			continue
		}

		err := item.run()
		item.errChan <- err // panics if shutdown is triggered as errChan is then closed
		if err != nil {
			continue
		}

		s.metric.Lock()
		s.metric.Reconnects = append(s.metric.Reconnects, time.Now())
		s.metric.Unlock()

		// 1000 identify / 24 hours rate limit check
		if s.metric.ReconnectsSince(24*time.Hour) > (s.identifiesPer24H - 1) {
			s.metric.Lock()
			oldest := s.metric.Reconnects[len(s.metric.Reconnects)-int(s.identifiesPer24H)]
			s.metric.Unlock()

			penalty = (24 * time.Hour) - time.Since(oldest)
			s.logger.Info(s.lpre, "shard identifying hit 1k rate limit and connections are halted for", penalty)
		}

		select {
		case <-s.shutdownChan:
			s.logger.Debug(s.lpre, "shard identify-rate-limiter got shutdown signal")
			return
		case <-time.After(s.timeout + penalty):
		}
	}
}
