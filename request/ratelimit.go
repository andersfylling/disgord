package request

import "net/http"

// const
var majorEndpointPrefixes = []string{
	"/channels/",
	"/guilds/",
	"/webhooks/",
}

type RateLimiter interface{}

func NewRateLimit() *RateLimit {
	return &RateLimit{}
}

type RateLimit struct{}

func (r *RateLimit) HandleResponse(res *http.Response) {
}
