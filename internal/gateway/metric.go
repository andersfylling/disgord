package gateway

import (
	"sync"
	"time"
)

const MetricReconnectPeriod = time.Hour * 48

// TODO-1: make it specific to individual shards for more insight
// TODO-2: Limit storage period
type IdentifyMetric struct {
	sync.Mutex
	Reconnects []time.Time // last 48h or (see const ReconnectPeriod)
}

func (s *IdentifyMetric) cleanup() {
	now := time.Now()
	s.Reconnects = rotateByTime(s.Reconnects, now.Add(MetricReconnectPeriod))
}

// ReconnectsSince counts the number of reconnects since t, where t can be no more than ReconnectPeriod (48h?)
func (s *IdentifyMetric) ReconnectsSince(d time.Duration) (counter uint) {
	limit := time.Now().Add(d)
	s.Lock()
	defer s.Unlock()

	for i := len(s.Reconnects) - 1; i >= 0; i-- {
		rt := s.Reconnects[i]
		if rt.Before(limit) {
			counter++
		} else {
			break
		}
	}

	s.cleanup()
	return counter
}
