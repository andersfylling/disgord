package gateway

import (
	"sync"
	"time"
)

const MetricReconnectPeriod = time.Hour * 48

// rotateByTime every timestamp after "limit" is deleted. Assumes the oldest entries are first.
func rotateByTime(times []time.Time, limit time.Time) []time.Time {
	var delim int
	for i := range times {
		if times[i].After(limit) {
			delim++
		} else {
			break
		}
	}

	// shift
	for i := 0; i+delim < len(times); i++ {
		times[i] = times[delim+i]
	}

	return times[:len(times)-delim]
}

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
	now := time.Now()
	s.Lock()
	defer s.Unlock()

	for _, rt := range s.Reconnects {
		if rt.Before(now.Add(d)) {
			counter++
		}
	}

	s.cleanup()
	return counter
}
