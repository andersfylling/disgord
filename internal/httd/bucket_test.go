// +build !integration

package httd

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
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
		global.resetTime = time.Now().Add(1 * time.Hour)
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
		global.resetTime = time.Now().Add(1 * time.Hour)
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
		resp.Header.Set(XRateLimitReset, strconv.FormatFloat(float64(reset.UnixNano())/float64(time.Second), 'f', 5, 64))
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
		diff := bucket.discordResetTime.Nanosecond()/int(time.Millisecond) - reset.Nanosecond()/int(time.Millisecond)
		if diff > 2 || diff < -2 {
			t.Errorf("reset did not update. Got %s, wants %s", bucket.discordResetTime.String(), reset.String())
		}
	})
}

func TestLtBucket_RespectRateLimit(t *testing.T) {
	// X-Ratelimit-Bucket:[f56681194ebea036dd1297f1184bf7bd] X-Ratelimit-Limit:[2] X-Ratelimit-Remaining:[0] X-Ratelimit-Reset:[1571597130835] X-Ratelimit-Reset-After:[2787.131]

	mngr := NewManager(nil)
	id := "dlfjhdskfhjdskfjsd"
	mngr.Bucket(id, func(bucket RESTBucket) {
		_, _, _ = bucket.Transaction(context.Background(), func() (response *http.Response, bytes []byte, err error) {
			limit := 2
			remaining := 0
			reset := time.Now().Add(2 * time.Hour)
			hash := "f56681194ebea036dd1297f1184bf7bd"

			resp := &http.Response{
				Header:     make(http.Header),
				StatusCode: http.StatusBadRequest,
			}
			resp.Header.Set(XRateLimitBucket, hash)
			resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
			resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
			resp.Header.Set(XRateLimitReset, strconv.FormatFloat(float64(reset.UnixNano())/float64(time.Second), 'f', 4, 64))
			resp.Header.Set("date", time.Now().Format(time.RFC1123))

			resp.Header, _ = NormalizeDiscordHeader(resp.StatusCode, resp.Header, nil)

			return resp, nil, nil
		})
	})

	mngr.Bucket(id, func(bucket RESTBucket) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
		defer cancel()
		_, _, err := bucket.Transaction(ctx, func() (response *http.Response, bytes []byte, err error) {
			return nil, nil, nil
		})

		if err == nil || !strings.Contains(err.Error(), "time out") {
			t.Error("should have been rate limited")
		}
	})

	// make the info outdated
	mngr.Bucket(id, func(bucket RESTBucket) {
		b := bucket.(*ltBucket)
		b.resetTime = time.Unix(0, time.Now().UnixNano()-int64(5*time.Hour))
		b.discordResetTime = b.resetTime.Add(3 * time.Second)

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
		defer cancel()
		_, _, err := bucket.Transaction(ctx, func() (response *http.Response, bytes []byte, err error) {
			return nil, nil, errors.New("resp error yay")
		})

		if !strings.Contains(err.Error(), "resp error yay") {
			t.Error("should have been able to send the request")
		}
	})

}
