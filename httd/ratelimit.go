package httd

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andersfylling/snowflake"
)

const (
	XRateLimitLimit     = "X-RateLimit-Limit"
	XRateLimitRemaining = "X-RateLimit-Remaining"
	XRateLimitReset     = "X-RateLimit-Reset"
	XRateLimitGlobal    = "X-RateLimit-Global"
	RateLimitRetryAfter = "Retry-After"
)

// const
var majorEndpointPrefixes = []string{
	"/channels/",
	"/guilds/",
	"/webhooks/",
}

// TODO: fix ratelimiting logic
func RatelimitChannel(id snowflake.ID) string {
	return "c:" + id.String()
}

func RatelimitGuild(id snowflake.ID) string {
	return "g:" + id.String()
}

func RatelimitWebsocket(id snowflake.ID) string {
	return "w:" + id.String()
}

func RatelimitUsers() string {
	return "u"
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
		global: &Bucket{
			global:true,

		},
	}
}

// RateLimit
// TODO: a bucket is created for every request. Might want to delete them after a while. seriously.
// `/users/1` has the same ratelimiter as `/users/2`
// but any major endpoint prefix does not: `/channels/1` != `/channels/2`
type RateLimit struct {
	buckets map[string]*Bucket
	global *Bucket
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

type ratelimitBody struct {
	Message string `json:"message"`
	RetryAfter int64 `json:"retry_after"`
	Global bool `json:"global"`
}

func (r *RateLimit) HandleResponse(key string, res *http.Response) {
	var err error
	var global bool
	var limit uint64
	var remaining uint64
	var reset int64
	var body *ratelimitBody
	var noBody bool

	// read body as well
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		noBody = true
	}

	// global?
	if res.Header.Get(XRateLimitGlobal) == "true" || (!noBody && body.Global) {
		global = true
	}

	// max number of request before reset
	if res.Header.Get(XRateLimitLimit) != "" || (!noBody && body.Global) {
		limit, err = strconv.ParseUint(res.Header.Get(XRateLimitLimit), 10, 64)
		if err != nil {
			// TODO: logging
		}
	}

	// remaining requests before reset
	remainingStr := res.Header.Get(XRateLimitRemaining)
	if remainingStr != "" {
		remaining, err = strconv.ParseUint(remainingStr, 10, 64)
		if err != nil {
			// TODO: logging
		}
	}

	// reset unix timestamp
	resetStr := res.Header.Get(XRateLimitReset)
	if resetStr != "" {
		// here we get a unix timestamp in seconds, which we convert to milliseconds
		reset, err = strconv.ParseInt(remainingStr, 10, 64)
		if err == nil {
			reset *= 1000 // => milliseconds
		} else {
			// TODO: logging
		}
	} else if res.Header.Get(RateLimitRetryAfter) != "" || (!noBody && body.RetryAfter > 0) {
		// here we are given a delay in millisecond, which we need to convert into a timestamp
		if res.Header.Get(RateLimitRetryAfter) != "" {
			reset, err = strconv.ParseInt(res.Header.Get(RateLimitRetryAfter), 10, 64)
			if err != nil {
				reset = 0
			}
		} else if !noBody && body.RetryAfter > 0 {
			reset = body.RetryAfter
		}

		// convert diff to timestamp
		reset += time.Now().UnixNano() / 1000
	}

	if global {
		r.global.mu.Lock()
		defer r.global.mu.Unlock()

		if limit != 0 {
			r.global.limit = limit
		}
		if remaining != 0 {
			r.global.remaining = remaining
		}
		if reset != 0 {
			r.global.reset = reset
		}
	} else {
		bucket := r.Bucket(key)
		bucket.mu.Lock()
		defer bucket.mu.Unlock()

		if limit != 0 {
			bucket.limit = limit
		}
		if remaining != 0 {
			bucket.remaining = remaining
		}
		if reset != 0 {
			bucket.reset = reset
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
