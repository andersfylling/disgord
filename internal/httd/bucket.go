package httd

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/util"
)

type bucketID struct {
	global     bool
	mustReturn bool
	resetTime  time.Time
}

func newBucket(global *bucket) (b *bucket) {
	b = &bucket{
		global:    global,
		remaining: -1,
		resetTime: time.Now().Add(1 * time.Hour),
	}

	return b
}

type bucketTransaction func() (resp *http.Response, body []byte, err error)

// bucket holds the rate limit info for a given hash or endpoint
type bucket struct {
	mu         sync.RWMutex
	atomicLock util.AtomicLock
	hash       string // discord designated hash

	queue util.TicketQueue

	remaining uint // remaining requests
	resetTime time.Time

	usingGlobal bool
	global      *bucket
}

func (b *bucket) AcquireLock() (locked bool) {
	if locked = b.atomicLock.AcquireLock(); !locked {
		return false
	}

	if b != b.global {
		// peek global bucket
		b.global.mu.RLock()
		globalLock := b.global.active()
		b.global.mu.RUnlock()

		if globalLock {
			// only the one with an acquired lock can write
			// so check if the globalLock has changed since the read
			locked = b.global.atomicLock.AcquireLock()
			b.global.mu.RLock()
			if !b.global.active() {
				b.global.atomicLock.Unlock()
				locked = true
			} else {
				b.usingGlobal = true
			}
			b.global.mu.RUnlock()
		}
	}

	return locked
}

func (b *bucket) Transaction(ctx context.Context, do bucketTransaction) (resp *http.Response, body []byte, err error) {
	ticket := b.queue.NewTicket()
	for {
		select {
		case <-ctx.Done():
			b.queue.Delete(ticket)
			return nil, nil, errors.New("time out")
		case <-time.After(10 * time.Millisecond):
		}

		if !b.queue.Next(ticket, b.AcquireLock) {
			continue
		}
		break
	}
	defer b.atomicLock.Unlock()
	if b.usingGlobal {
		defer b.global.atomicLock.Unlock()
	}

	// set active bucket
	var bucket *bucket
	if b.usingGlobal {
		bucket = b.global
	} else {
		bucket = b
	}

	// check if rate limited and try to wait it out
	var wait time.Duration
	now := time.Now()
	if bucket.resetTime.After(now) && bucket.remaining == 0 {
		wait = bucket.resetTime.Sub(now)
	}
	select {
	case <-ctx.Done():
		return nil, nil, errors.New("time out")
	case <-time.After(wait):
	}

	// send request
	resp, body, err = do()
	if err != nil {
		return nil, nil, err
	}

	// update bucket info

	// reduce remaining if rate limited
	// remaining == -1 when no rate limit info currently exists
	if bucket.remaining > 0 {
		bucket.remaining--
	}

	// to synchronize the timestamp between the bot and the discord server
	// we assume the current time is equal the header date
	discordTime, err := HeaderToTime(resp.Header)
	if err != nil {
		discordTime = time.Now()
	}

	localTime := time.Now()
	diff := localTime.Sub(discordTime)

	var isGlobal bool
	bucketHash := resp.Header.Get(XRateLimitBucket)
	if _, ok := resp.Header[XRateLimitBucket]; ok && bucketHash == "" {
		isGlobal = true
	}
	isGlobal = isGlobal || resp.Header.Get(XRateLimitGlobal) == "true"

	// if this is not a 429 error we can determine if the local bucket is a global one or not
	if resp.StatusCode != http.StatusTooManyRequests && b.hash == "" {
		if isGlobal {
			b.hash = GlobalHash
		} else if bucketHash != "" {
			b.hash = bucketHash
		}
	}

	var reset time.Time
	var remaining uint
	if resetStr := resp.Header.Get(XRateLimitReset); resetStr != "" {
		epoch, _ := strconv.ParseInt(resetStr, 10, 64)
		reset = time.Unix(0, epoch+diff.Nanoseconds())
	}

	if remainingStr := resp.Header.Get(XRateLimitRemaining); remainingStr != "" {
		remainingInt, _ := strconv.ParseInt(remainingStr, 10, 64)
		remaining = uint(remainingInt)
	}

	// update bucket reference to whatever the header regards
	if isGlobal {
		bucket = b.global // TODO-?: AcquireLock?
	} else {
		if !reset.IsZero() {
			b.resetTime = reset
		}

		if b.resetTime.Before(reset) {
			b.resetTime = reset
		} else {

		}
	}
}

func (b *bucket) active() bool {
	return b.remaining >= 0
}
