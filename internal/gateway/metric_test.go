// +build !integration

package gateway

import (
	"testing"
	"time"
)

func TestShardMetric_ReconnectsSince(t *testing.T) {
	now := time.Now()

	times := []time.Time{
		now.Add(time.Hour),
		now.Add(time.Minute),
		now.Add(time.Second),
		now,
	}

	metric := &IdentifyMetric{
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
