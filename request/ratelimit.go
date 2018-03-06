package request

import (
	"net/http"
	"strconv"
	"strings"
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

type RateLimiter interface{}

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

	if bucket, exists = r.buckets[key]; !exists {
		bucket = &Bucket{
			endpoint: key,
		}
	}

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
	reset     int64  // milliseconds, even tho discord prefers seconds. global uses milliseconds however.
	global    bool   // global rate limiter
}

func (b *Bucket) limited() bool {
	return ((time.Now().UnixNano() / 1000) - b.reset) < 0
}

func (b *Bucket) timeout() int64 {
	return (time.Now().UnixNano() / 1000) - b.reset
}
