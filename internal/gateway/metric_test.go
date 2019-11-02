package gateway

import (
	"testing"
	"time"
)

func TestRotateByTimes(t *testing.T) {
	now := time.Now()

	times := []time.Time{
		now.Add(time.Hour),
		now.Add(time.Minute),
		now.Add(time.Second),
		now,
	}

	times2 := rotateByTime(times, now.Add(5*time.Minute))
	if len(times) == len(times2) {
		t.Errorf("original times slice should be more than the new times slice. Got %d, wants %d", len(times2), len(times))
	}
	if len(times)-1 != len(times2) {
		t.Errorf("new times slice should be one less than the original slice. Got %d, wants %d", len(times2), len(times))
	}
}

func TestShardMetric_ReconnectsSince(t *testing.T) {
	now := time.Now()

	times := []time.Time{
		now.Add(time.Hour),
		now.Add(time.Minute),
		now.Add(time.Second),
		now,
	}

	metric := &ShardMetric{
		Reconnects: times,
	}

	since := metric.ReconnectsSince(2 * time.Millisecond)
	if int(since) != 1 {
		t.Errorf("expected to count one entry. Got %d, wants %d", since, 1)
	}

	since = metric.ReconnectsSince(2 * time.Hour)
	if int(since) != len(times) {
		t.Errorf("expected to count every entry. Got %d, wants %d", since, len(times))
	}
}
