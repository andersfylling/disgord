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
	UpdateRegisters(key string, res *http.Response, responseBody []byte)
	WaitTime(req *Request) time.Duration
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

// NewDiscordTimeDiff constructor
func NewDiscordTimeDiff() *DiscordTimeDiff {
	return &DiscordTimeDiff{
		Local:   time.Now(),
		Discord: time.Now(),
	}
}

// DiscordTimeDiff allows for synchronizing the timestamp found in http responses with out local timestamp to avoid
// desyncs
type DiscordTimeDiff struct {
	sync.RWMutex
	Local   time.Time
	Discord time.Time
	offset  time.Duration
}

// Update updates the current time difference on a new response
func (d *DiscordTimeDiff) Update(now time.Time, discord time.Time) {
	d.Lock()
	defer d.Unlock()

	d.Local = now
	d.Discord = discord
	d.calculateOffset()
}

func (d *DiscordTimeDiff) calculateOffset() {
	if d.Local.After(d.Discord) {
		d.offset = d.Local.Sub(d.Discord) * -1
	} else {
		d.offset = d.Discord.Sub(d.Local)
	}
}

// Now calculates the current timestamp with the http response timestamp in mind
func (d *DiscordTimeDiff) Now() (t time.Time) {
	d.RLock()
	defer d.RUnlock()

	t = time.Now().Add(d.offset)
	return
}

// NewRateLimit creates a new rate limit manager
func NewRateLimit() *RateLimit {
	return &RateLimit{
		buckets:  make(map[string]*Bucket),
		global:   &Bucket{},
		TimeDiff: NewDiscordTimeDiff(),
	}
}

// RateLimit ...
type RateLimit struct {
	buckets  map[string]*Bucket
	global   *Bucket
	TimeDiff *DiscordTimeDiff

	mu sync.RWMutex
}

// Bucket returns a bucket given the key (or ID) for a rate limit bucket. If
// no bucket exists for the key, one will be created.
func (r *RateLimit) Bucket(key string) *Bucket {
	var bucket *Bucket
	var exists bool

	r.mu.Lock()
	if bucket, exists = r.buckets[key]; !exists {
		r.buckets[key] = &Bucket{
			endpoint: key,
			reset:    r.TimeDiff.Now().UnixNano() / int64(time.Millisecond),
		}
		bucket = r.buckets[key]
	}
	r.mu.Unlock()

	return bucket
}

// RateLimitTimeout returns the time left before the rate limit for a given key
// is reset. This takes the global rate limit into account.
func (r *RateLimit) RateLimitTimeout(key string) int64 {
	global := r.global.timeout(r.TimeDiff.Now())

	bucket := r.Bucket(key)
	unique := bucket.timeout(r.TimeDiff.Now())

	if global > unique {
		return global
	}
	return unique
}

// RateLimited checks if the given key is rate limited. This takes the global
// rate limiter into account.
func (r *RateLimit) RateLimited(key string) bool {
	if r.global.limited(r.TimeDiff.Now()) {
		return true
	}

	bucket := r.Bucket(key)
	return bucket.limited(r.TimeDiff.Now())
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

// UpdateRegisters updates the relevant buckets and time desync between the
// client and the Discord servers.
func (r *RateLimit) UpdateRegisters(key string, resp *http.Response, content []byte) {
	// update time difference
	if discordTime, err := HeaderToTime(&resp.Header); err == nil {
		r.TimeDiff.Update(time.Now(), discordTime)
	}

	// update bucket
	info, err := ExtractRateLimitInfo(resp, content)
	if err != nil {
		return // TODO: logging
	}

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
	bucket.update(info, r.TimeDiff.Now())
	bucket.mu.Unlock()
}

// ---------------------

// Bucket holds the rate limit info for a given key or endpoint
type Bucket struct {
	endpoint  string // endpoint where rate limit is applied. endpoint = key
	limit     uint64 // total allowed requests before rate limit
	remaining uint64 // remaining requests
	reset     int64  // unix milliseconds, even tho discord prefers seconds. global uses milliseconds however.

	mu sync.RWMutex
}

func (b *Bucket) update(info *RateLimitInfo, now time.Time) {
	b.limit = uint64(info.Limit)
	b.remaining = uint64(info.Remaining)
	b.reset = info.Reset

	retryAt := info.RetryAfter + (now.UnixNano() / int64(time.Millisecond))
	if b.reset < retryAt {
		b.reset = retryAt
	}
}

func (b *Bucket) limited(now time.Time) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.reset > (now.UnixNano()/int64(time.Millisecond)) && b.remaining == 0
}

func (b *Bucket) timeout(now time.Time) int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	nowMilli := now.UnixNano() / int64(time.Millisecond)
	var timeout int64
	if b.reset > nowMilli && b.remaining == 0 { // will b.reset > nowMilli if remaining == 0?
		timeout = b.reset - nowMilli
	}

	return timeout
}
