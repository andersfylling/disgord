package websocket

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord/websocket/cmd"
)

func newRatelimiter() ratelimiter {
	rl := ratelimiter{
		buckets: map[string]rlBucket{},
		global:  newRatelimitBucket(120, 60),
	}
	rl.buckets[cmd.UpdateStatus] = newRatelimitBucket(5, 60)

	return rl
}

type rlEntry struct {
	unix int64
	cmd  string
}

func newRatelimitBucket(requests, seconds int64) rlBucket {
	return newRatelimitBucketNano(requests, seconds*int64(time.Second))
}

func newRatelimitBucketNano(requests, nano int64) rlBucket {
	return rlBucket{
		entries:  make([]rlEntry, requests),
		duration: (time.Duration(nano) * time.Nanosecond).Nanoseconds(),
	}
}

type rlBucket struct {
	entries  []rlEntry
	duration int64
}

func (b *rlBucket) Blocked() bool {
	last := b.entries[len(b.entries)-1]
	return time.Now().UnixNano()-last.unix <= b.duration
}

func (b *rlBucket) Insert(cmd string) {
	// TODO: we could shift the last valid element to the bottom and then not shift on every insert
	// b.entries = append(b.entries[1:], b.entries[:len(b.entries)-2])
	for i := len(b.entries) - 1; i > 0; i-- {
		b.entries[i] = b.entries[i-1]
	}
	b.entries[0] = rlEntry{
		unix: time.Now().UnixNano(),
		cmd:  cmd,
	}
}

type ratelimiter struct {
	sync.RWMutex
	buckets map[string]rlBucket
	global  rlBucket
}

func (rl *ratelimiter) Request(command string) (accepted bool) {
	rl.Lock()
	defer rl.Unlock()

	// global
	if rl.global.Blocked() {
		return false
	}
	rl.global.Insert(command)

	// bucket specific
	if bucket, exists := rl.buckets[command]; exists {
		if bucket.Blocked() {
			return false
		}
		bucket.Insert(command)
	}

	return true
}
