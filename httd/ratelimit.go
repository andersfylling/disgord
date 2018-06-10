package httd

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	XRateLimitLimit     = "X-RateLimit-Limit"
	XRateLimitRemaining = "X-RateLimit-Remaining"
	XRateLimitReset     = "X-RateLimit-Reset"
	XRateLimitGlobal    = "X-RateLimit-Global"
)

// const
var majorEndpointPrefixes = []string{
	"/channels/",
	"/guilds/",
	"/webhooks/",
}

type RateLimiter interface {
	Bucket(key string) *Bucket
	RateLimitTimeout(key string) int64
	RateLimited(key string) bool
	HandleResponse(key string, res *http.Response)
}

func NewRateLimit() *RateLimit {
	return &RateLimit{
		buckets: make(map[string]*Bucket),
	}
}

// RateLimit
// TODO: a bucket is created for every request. Might want to delete them after a while. seriously.
// `/users/1` has the same ratelimiter as `/users/2`
// but any major endpoint prefix does not: `/channels/1` != `/channels/2`
type RateLimit struct {
	buckets map[string]*Bucket
	mu      sync.RWMutex
}

func (r *RateLimit) Bucket(key string) *Bucket {
	var bucket *Bucket
	var exists bool

	// check for major endpoints
	// TODO: this feels frail
	var endpoint string
	for _, major := range majorEndpointPrefixes {
		if !strings.HasPrefix(key, major) {
			continue
		}
		pathAfterMajor := strings.TrimPrefix(key, major)
		endpoint = major
		for _, r := range pathAfterMajor {
			if r == '/' {
				break
			}
			endpoint += string(r)
		}
	}
	if endpoint == "" {
		endpoint = key
	}

	r.mu.Lock()
	if bucket, exists = r.buckets[key]; !exists {
		r.buckets[key] = &Bucket{
			endpoint: key,
			reset:    time.Now().UnixNano() / 1000,
		}
		bucket = r.buckets[key]
	}
	r.mu.Unlock()

	return bucket
}

func (r *RateLimit) RateLimitTimeout(key string) int64 {
	bucket := r.Bucket(key)
	return bucket.timeout()
}

func (r *RateLimit) RateLimited(key string) bool {
	bucket := r.Bucket(key)
	return bucket.limited()
}

func (r *RateLimit) HandleResponse(key string, res *http.Response) {
	bucket := r.Bucket(key)
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if !bucket.global && res.Header.Get(XRateLimitGlobal) == "true" {
		bucket.global = true
	}

	if bucket.limit == 0 && res.Header.Get(XRateLimitLimit) != "" {
		limit, err := strconv.ParseUint(res.Header.Get(XRateLimitLimit), 10, 64)
		if err == nil {
			bucket.limit = limit
		} else {
			// TODO: log
		}
	}

	remainingStr := res.Header.Get(XRateLimitRemaining)
	if remainingStr != "" {
		remaining, err := strconv.ParseUint(remainingStr, 10, 64)
		if err == nil {
			bucket.remaining = remaining
		} else {
			// TODO: log
		}
	}

	resetStr := res.Header.Get(XRateLimitReset)
	if resetStr != "" {
		reset, err := strconv.ParseInt(remainingStr, 10, 64)
		if err == nil {
			bucket.reset = reset * 1000 //convert seconds to milliseconds
			// TODO: if global, reset should be remaining milliseconds. convert it => (time.now + globalReset)
		} else {
			// TODO: log
		}
	}
}

// ---------------------

type Bucket struct {
	endpoint  string // endpoint where rate limit is applied. endpoint = key
	limit     uint64 // total allowed requests before rate limit
	remaining uint64 // remaining requests
	reset     int64  // unix milliseconds, even tho discord prefers seconds. global uses milliseconds however.
	global    bool   // global rate limiter

	mu sync.RWMutex
}

func (b *Bucket) limited() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.reset > (time.Now().UnixNano() / 1000)
}

func (b *Bucket) timeout() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	now := time.Now().UnixNano() / 1000
	if b.reset > now {
		return b.reset - now
	}
	return 0
}
