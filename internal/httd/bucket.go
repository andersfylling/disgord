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

func newLeakyBucket(global *ltBucket) (b *ltBucket) {
	b = &ltBucket{
		remaining: -1,
		resetTime: time.Now(),
		global:    global,
	}

	return b
}

type bucketTransaction = func() (resp *http.Response, body []byte, err error)

// ltBucket combines leaky and token buckets to allow time aware of the REST requests while they're in queue.
type ltBucket struct {
	mu         sync.RWMutex
	atomicLock util.AtomicLock
	hash       string // discord designated hash

	queue util.TicketQueue // Ticket => Token

	remaining        int       // remaining requests
	resetTime        time.Time // affected by time diff
	discordResetTime time.Time // unaffected by time diff

	updatedAt time.Time // use date from discord header

	// this bucket is global if this.global is nil or this == this.global
	global      *ltBucket
	usingGlobal bool
}

var _ RESTBucket = (*ltBucket)(nil)

func (b *ltBucket) AcquireLock() (locked bool) {
	if locked = b.atomicLock.AcquireLock(); !locked {
		return false
	}

	if _, err := b.SelectiveGlobalLock(); err != nil {
		b.atomicLock.Unlock()
		return false
	}

	return true
}

func (b *ltBucket) SelectiveGlobalLock() (locked bool, err error) {
	if b != b.global {
		// peek global ltBucket
		b.global.mu.RLock()
		globalLock := b.global.active()
		b.global.mu.RUnlock()
		// TODO: can this cause http 429?
		if globalLock {
			// so check if the globalLock has changed since the read
			if locked = b.global.atomicLock.AcquireLock(); !locked {
				return false, errors.New("unable to acquire needed global lock")
			}

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

	return locked, nil
}

func (b *ltBucket) Transaction(ctx context.Context, do bucketTransaction) (resp *http.Response, body []byte, err error) {
	// wait until you are next in line and you can acquire a lock
	// this is to support timeout/cancellation for stacked requests
	// TODO: on success, every request with same endpoint or a valid subset can be fulfilled locally
	// reqA = /guilds/1/members?limit=100
	// reqB = /guilds/1/members?limit=10
	// reqB is a subset of A, and therefore reqA can create a response for reqB locally (must be deep copy - djp)
	token := b.queue.NewTicket()
	for {
		select {
		case <-ctx.Done():
			b.queue.Delete(token)
			return nil, nil, errors.New("time out")
		case <-time.After(10 * time.Millisecond):
			// TODO-perf: this wastes a lot of CPU usage
		}

		if !b.queue.Next(token, b.AcquireLock) {
			continue
		}
		break
	}
	defer b.atomicLock.Unlock()
	if b.usingGlobal {
		defer b.global.atomicLock.Unlock()
		defer func() {
			b.usingGlobal = false
		}()
	}

	// set active ltBucket
	var bucket *ltBucket
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
	if deadline, ok := ctx.Deadline(); ok && deadline.Before(time.Now().Add(wait)) {
		return nil, nil, errors.New("time out, bucket resets in " + wait.String())
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

	// update ltBucket info
	// reduce remaining if needed
	if !b.updateAfterRequest(resp.Header, resp.StatusCode) && bucket.remaining > 0 {
		bucket.remaining--
	}

	return resp, body, nil
}

// updateAfterRequests updates the bucket with the latest rate limit info from http responses.
//
// Note! you must call NormalizeDiscordHeader before using this.
func (b *ltBucket) updateAfterRequest(header http.Header, statusCode int) (adjustedRemaining bool) {
	if normalized := header.Get(DisgordNormalizedHeader); normalized == "" {
		panic("headers were not normalized to use milliseconds")
	}

	// to synchronize the timestamp between the bot and the discord server
	// we assume the current time is equal the header date
	discordTime, err := HeaderToTime(header)
	if err != nil {
		discordTime = time.Now()
	}

	localTime := time.Now()
	diff := localTime.Sub(discordTime)

	var isGlobal bool
	bucketHash := header.Get(XRateLimitBucket)
	if _, ok := header[XRateLimitBucket]; ok && bucketHash == "" {
		isGlobal = true
	}
	isGlobal = isGlobal || header.Get(XRateLimitGlobal) == "true"

	// if this is not a 429 error we can determine if the local ltBucket is a global one or not
	if statusCode != http.StatusTooManyRequests && b.hash == "" {
		if isGlobal {
			b.hash = GlobalHash
		} else if bucketHash != "" {
			b.hash = bucketHash
		}
	}

	var reset time.Time
	var discordReset time.Time
	var remaining int = -1
	if resetStr := header.Get(XRateLimitReset); resetStr != "" {
		epoch, _ := strconv.ParseInt(resetStr, 10, 64)
		epoch *= int64(time.Millisecond) // ms => nano
		reset = time.Unix(0, epoch+diff.Nanoseconds())
		discordReset = time.Unix(0, epoch)
	}

	if remainingStr := header.Get(XRateLimitRemaining); remainingStr != "" {
		remainingInt64, _ := strconv.ParseInt(remainingStr, 10, 64)
		if remainingInt64 >= 0 {
			remaining = int(remainingInt64)
		}
	}

	// update ltBucket reference to whatever the header regards
	var bucket *ltBucket
	if isGlobal {
		if b.global == b {
			bucket = b
		} else {
			bucket = b.global
		}
		bucket.mu.Lock()
		defer bucket.mu.Unlock()
	} else {
		bucket = b // no need to lock normal buckets
		if !(b.global == nil || b == b.global) && bucketHash != "" {
			b.hash = bucketHash
		}
	}

	if discordReset.Before(time.Unix(0, int64(time.Hour))) {
		return false
	}

	// TODO: this can be simpler
	// use discord reset time, as the local reset can be different in ms or s per request.
	if discordReset.After(bucket.discordResetTime) {
		bucket.resetTime = reset
		bucket.discordResetTime = discordReset
		bucket.remaining = remaining
		bucket.updatedAt = discordTime
		adjustedRemaining = true
	} else if bucket.discordResetTime == discordReset {
		if bucket.remaining == -1 || bucket.remaining > remaining {
			bucket.remaining = remaining
			bucket.updatedAt = discordTime
			bucket.discordResetTime = discordReset
			adjustedRemaining = true
		}
	}

	return adjustedRemaining
}

func (b *ltBucket) active() bool {
	return b.remaining >= 0 && !time.Now().After(b.resetTime)
}
