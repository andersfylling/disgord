package httd

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// http rate limit identifiers
const (
	XRateLimitLimit      = "X-RateLimit-Limit"
	XRateLimitRemaining  = "X-RateLimit-Remaining"
	XRateLimitReset      = "X-RateLimit-Reset" // is converted from seconds to milliseconds!
	XRateLimitGlobal     = "X-RateLimit-Global"
	RateLimitRetryAfter  = "Retry-After"
	GlobalRateLimiterKey = ""
)

// RateLimiter is the interface for the ratelimit manager
type RateLimiter interface {
	Bucket(key string) *Bucket
	RateLimitTimeout(key string) int64
	RateLimited(key string) bool
	UpdateRegisters(key string, adjuster RateLimitAdjuster, res *http.Response, responseBody []byte)
	WaitTime(req *Request) time.Duration
	RequestPermit(key string) (timeout time.Duration, err error)
}

type ratelimitBody struct {
	Message    string `json:"message"`
	RetryAfter int64  `json:"retry_after"`
	Global     bool   `json:"global"`
	Empty      bool   `json:"-"`
}

// RateLimitInfo is populated by Discord http responses in order to obtain rate limits
type RateLimitInfo struct {
	Message    string `json:"message"`
	RetryAfter int64  `json:"retry_after"`
	Global     bool   `json:"global"`
	Limit      int    `json:"-"`
	Remaining  int    `json:"-"`
	Reset      int64  `json:"-"`
	Empty      bool   `json:"-"`
}

// RateLimited check if a response was rate limited
func RateLimited(resp *http.Response) bool {
	return resp.StatusCode == http.StatusTooManyRequests
}

// GlobalRateLimit assumes that there will always be a header entry when a global rate limit kicks in
func GlobalRateLimit(resp *http.Response) bool {
	return resp.Header.Get(XRateLimitGlobal) == "true"
}

// ExtractRateLimitInfo uses the RateLimitInfo struct to obtain rate limits from the Discord response
func ExtractRateLimitInfo(resp *http.Response, body []byte) (info *RateLimitInfo, err error) {
	info = &RateLimitInfo{}

	// extract header information
	limitStr := resp.Header.Get(XRateLimitLimit)
	remainingStr := resp.Header.Get(XRateLimitRemaining)
	resetStr := resp.Header.Get(XRateLimitReset)
	retryAfterStr := resp.Header.Get(RateLimitRetryAfter)

	// convert types
	if limitStr != "" {
		info.Limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return
		}
	}
	if remainingStr != "" {
		info.Remaining, err = strconv.Atoi(remainingStr)
		if err != nil {
			return
		}
	}
	if resetStr != "" {
		info.Reset, err = strconv.ParseInt(resetStr, 10, 64)
		if err != nil {
			return
		}
		info.Reset *= 1000 // second => milliseconds
	}
	if retryAfterStr != "" {
		info.RetryAfter, err = strconv.ParseInt(retryAfterStr, 10, 64)
		if err != nil {
			return
		}
	}

	// the body only contains information when a rate limit is exceeded
	if RateLimited(resp) && len(body) > 0 {
		err = Unmarshal(body, &info)
	}
	if !info.Global && GlobalRateLimit(resp) {
		info.Global = true
	}
	return
}

// HeaderToTime takes the response header from Discord and extracts the
// timestamp. Useful for detecting time desync between discord and client
func HeaderToTime(header *http.Header) (t time.Time, err error) {
	// date: Fri, 14 Sep 2018 19:04:24 GMT
	dateStr := header.Get("date")
	if dateStr == "" {
		err = errors.New("missing header field 'date'")
		return
	}

	t, err = time.Parse(time.RFC1123, dateStr)
	return
}

// NewRateLimit creates a new rate limit manager
func NewRateLimit() *RateLimit {
	return &RateLimit{
		buckets: make(map[string]*Bucket),
		global:  &Bucket{},
	}
}

// RateLimit ...
type RateLimit struct {
	buckets map[string]*Bucket
	global  *Bucket

	mu sync.RWMutex
}

// Bucket returns a bucket given the key (or ID) for a rate limit bucket. If
// no bucket exists for the key, one will be created.
func (r *RateLimit) Bucket(key string) *Bucket {
	var bucket *Bucket
	var exists bool

	r.mu.RLock()
	bucket, exists = r.buckets[key]
	r.mu.RUnlock()

	if !exists {
		r.mu.Lock()
		if bucket, exists = r.buckets[key]; !exists {
			r.buckets[key] = &Bucket{
				endpoint:        key,
				limit:           1,
				remaining:       1,
				shortestTimeout: int(time.Second.Nanoseconds() / int64(time.Millisecond)),
				reset:           time.Now().UnixNano() / int64(time.Millisecond),
			}
			bucket = r.buckets[key]
		}
		r.mu.Unlock()
	}

	return bucket
}

// RateLimitTimeout returns the time left before the rate limit for a given key
// is reset. This takes the global rate limit into account.
func (r *RateLimit) RateLimitTimeout(key string) int64 {
	now := time.Now()
	global := r.global.timeout(now)

	bucket := r.Bucket(key)
	unique := bucket.timeout(now)

	if global > unique {
		return global
	}
	return unique
}

// RateLimited checks if the given key is rate limited. This takes the global
// rate limiter into account.
func (r *RateLimit) RateLimited(key string) bool {
	now := time.Now()
	if r.global.limited(now) {
		return true
	}

	bucket := r.Bucket(key)
	return bucket.limited(now)
}

// WaitTime get's the remaining time before another request can be made.
// returns a time.Duration of milliseconds.
func (r *RateLimit) WaitTime(req *Request) time.Duration {
	timeout := int64(0)
	if r.RateLimited(req.Ratelimiter) {
		timeout = r.RateLimitTimeout(req.Ratelimiter) // number of milliseconds
	}

	// Duration requires nano seconds argument, so multiply with millisecond
	return time.Duration(timeout) * time.Millisecond
}

func (r *RateLimit) RequestPermit(key string) (timeout time.Duration, err error) {
	now := time.Now()

	r.global.mu.Lock()
	if r.global._limited(now) {
		// TODO: makes no sense to check limited first...
		timeout, err = r.global._requestPermit(key, now)
	}
	r.global.mu.Unlock()
	if err != nil || timeout > 0 {
		return
	}

	bucket := r.Bucket(key)
	bucket.mu.Lock()
	timeout, err = bucket._requestPermit(key, now)
	bucket.mu.Unlock()

	return
}

func adjustReset(timeout int64, adjuster RateLimitAdjuster) (newTimeout int64) {
	if adjuster != nil {
		d := time.Duration(timeout) * time.Millisecond
		d = adjuster(d)
		timeout = d.Nanoseconds() / int64(time.Millisecond)
	}

	return timeout
}

// UpdateRegisters updates the relevant buckets and time desync between the
// client and the Discord servers.
func (r *RateLimit) UpdateRegisters(key string, adjuster RateLimitAdjuster, resp *http.Response, content []byte) {
	now := time.Now()
	// update time difference
	var discordTime time.Time
	var err error
	if discordTime, err = HeaderToTime(&resp.Header); err != nil {
		discordTime = now
	}

	// update bucket
	info, err := ExtractRateLimitInfo(resp, content)
	if err != nil {
		return // TODO: logging
	}

	// adjust rate limit if desired (however, respect global rate limits)
	// In DisGord the Reset value is in milliseconds, not seconds.
	timeout := info.Reset - (int64(discordTime.UnixNano()) / int64(time.Millisecond))
	if !info.Global && timeout > 0 && adjuster != nil {
		timeout = adjustReset(timeout, adjuster)
	}
	if info.RetryAfter > 0 {
		timeout = info.RetryAfter
	}
	info.Reset = (now.UnixNano() / int64(time.Millisecond)) + timeout

	// select bucket
	// TODO: what if "key" is an endpoint with a global rate limiter only?
	var bucket *Bucket
	if info.Global {
		bucket = r.global
	} else {
		bucket = r.Bucket(key)
	}

	// update
	bucket.mu.Lock()
	bucket.update(info, now)
	if bucket.longestTimeout < int(timeout) {
		bucket.longestTimeout = int(timeout)
	} else {
		bucket.shortestTimeout = int(timeout)
	}
	bucket.mu.Unlock()
}

// ---------------------

// Bucket holds the rate limit info for a given key or endpoint
type Bucket struct {
	endpoint  string // endpoint where rate limit is applied. endpoint = key
	limit     int    // total allowed requests before rate limit
	remaining int    // remaining requests
	reset     int64  // unix milliseconds, even tho discord prefers seconds. global uses milliseconds however.

	// milliseconds
	longestTimeout  int // store the longest timeout to simulate a reset correctly
	shortestTimeout int

	mu sync.RWMutex
}

func (b *Bucket) update(info *RateLimitInfo, now time.Time) {
	b.limit = info.Limit
	b.remaining = info.Remaining
	b.reset = info.Reset
}

func (b *Bucket) _limited(now time.Time) bool {
	return b.reset > (now.UnixNano()/int64(time.Millisecond)) && b.remaining == 0
}

func (b *Bucket) limited(now time.Time) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b._limited(now)
}

func (b *Bucket) _timeout(now time.Time) int64 {
	nowMilli := now.UnixNano() / int64(time.Millisecond)
	var timeout int64
	if b.reset > nowMilli && b.remaining == 0 { // will b.reset > nowMilli if remaining == 0?
		timeout = b.reset - nowMilli
	}

	return timeout
}

func (b *Bucket) timeout(now time.Time) int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b._timeout(now)
}

func (b *Bucket) _requestPermit(key string, now time.Time) (timeout time.Duration, err error) {
	// make sure the restrictions are valid
	nowMilli := now.UnixNano() / int64(time.Millisecond)
	if b.reset <= nowMilli {
		longestTimeout := int64(b.longestTimeout)
		if longestTimeout == 0 {
			longestTimeout = time.Hour.Nanoseconds()
		}
		b.reset = nowMilli + longestTimeout
		b.remaining = b.limit
		if b.remaining == 0 {
			b.remaining++ // so we can do one request to get the new rate limits
		}
	}

	// see if we can execute a request right now, or for how long we need to wait
	b.remaining--
	if b.remaining < 0 {
		b.remaining++
		timeout = time.Duration(b.shortestTimeout) * time.Millisecond
	}

	return
}
