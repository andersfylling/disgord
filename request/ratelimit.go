package request

type RateLimiter interface{}

func NewRateLimit() *RateLimit {
	return &RateLimit{}
}

type RateLimit struct{}
