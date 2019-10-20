package httd

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestLtBucket_AcquireLock(t *testing.T) {
	t.Run("already-locked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		bucket := newLeakyBucket(global)
		bucket.atomicLock.AcquireLock()

		if success := bucket.AcquireLock(); success {
			t.Error("should not be able to acquire lock on locked bucket")
		}
	})
	t.Run("unlocked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		bucket := newLeakyBucket(global)

		if success := bucket.AcquireLock(); !success {
			t.Error("should be able to lock unlocked bucket")
		}
	})
	t.Run("global-unlocked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		global.remaining = 2
		if !global.active() {
			t.Fatal("incorrectly configured global bucket")
		}
		bucket := newLeakyBucket(global)

		if success := bucket.AcquireLock(); !success {
			t.Error("should be able to lock unlocked bucket")
		}
	})
	t.Run("global-locked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		global.atomicLock.AcquireLock()
		global.remaining = 2
		if !global.active() {
			t.Fatal("incorrectly configured global bucket")
		}

		bucket := newLeakyBucket(global)

		if success := bucket.AcquireLock(); success {
			t.Error("should be able to lock when global is locked bucket")
		}
	})
}

func TestLtBucket_updateAfterResponse(t *testing.T) {
	t.Run("update-fresh-bucket", func(t *testing.T) {
		global := newLeakyBucket(nil)
		bucket := newLeakyBucket(global)

		limit := 2
		remaining := 4
		reset := time.Now().Add(500 * time.Millisecond)
		hash := "sdlkfhsdlkfhsdkjafhsdf"

		resp := &http.Response{
			Header:     make(http.Header, 3),
			StatusCode: http.StatusTooManyRequests,
		}
		resp.Header.Set(XRateLimitBucket, hash)
		resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
		resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
		resp.Header.Set(XRateLimitReset, strconv.FormatFloat(float64(reset.UnixNano())/float64(time.Second), 'f', 4, 64))
		resp.Header.Set("date", time.Now().Format(time.RFC1123))

		header, err := NormalizeDiscordHeader(resp.StatusCode, resp.Header, nil)
		if err != nil {
			t.Fatal(err)
		}

		bucket.updateAfterRequest(header, resp.StatusCode)
		if bucket.hash != hash {
			t.Errorf("hash did not update. Got %s, wants %s", bucket.hash, hash)
		}
		if bucket.remaining != remaining {
			t.Errorf("remaining did not update. Got %d, wants %d", bucket.remaining, remaining)
		}
		if bucket.discordResetTime.Nanosecond()/int(time.Millisecond) != reset.Nanosecond()/int(time.Millisecond) {
			t.Errorf("reset did not update. Got %s, wants %s", bucket.discordResetTime.String(), reset.String())
		}
	})
}
