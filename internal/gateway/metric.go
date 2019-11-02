package gateway

import (
	"sync"
	"time"
)

const MetricReconnectPeriod = time.Hour * 48

type ShardMetric struct {
	sync.Mutex
	Reconnects         []time.Time // last 48h or (see const ReconnectPeriod)
	RequestedReconnect []time.Time // ^
}

func (s *ShardMetric) cleanup() {
	now := time.Now()
	s.Reconnects = rotateByTime(s.Reconnects, now.Add(MetricReconnectPeriod))
	s.RequestedReconnect = rotateByTime(s.RequestedReconnect, now.Add(MetricReconnectPeriod))
}

// ReconnectsSince counts the number of reconnects since t, where t can be no more than ReconnectPeriod (48h?)
func (s *ShardMetric) ReconnectsSince(d time.Duration) (counter uint) {
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
