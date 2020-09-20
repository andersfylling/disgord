// +build !integration

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
