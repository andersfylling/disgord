package ratelimit

import (
	"errors"
	"net/http"
	"strings"
)

const DiscordAPIURLPrefix = "https://discord.com/api/v"

type RateLimit struct {
	Original http.RoundTripper
	queue queueWithCancellation
}

var _ http.RoundTripper = (*RateLimit)(nil)
var _ DiscordRateLimitCompliant = (*RateLimit)(nil)

func (r *RateLimit) DiscordRateLimitCompliant() {}

func (r *RateLimit) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if !r.isDiscordAPIRequest(req) {
		return r.Original.RoundTrip(req)
	}

	return r.rateLimit(req, func() (*http.Response, error) {
		return r.Original.RoundTrip(req)
	})
}

func (r *RateLimit) isDiscordAPIRequest(req *http.Request) bool {
	return strings.HasPrefix(req.URL.String(), DiscordAPIURLPrefix)
}

func (r *RateLimit) rateLimit(req *http.Request, cb func() (*http.Response, error)) (*http.Response, error) {
	ctx := req.Context()


	wait := make(chan (chan interface{}))
	r.queue.insert(&queueItem{
		ctx:    ctx,
		notify: wait,
	})
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case done, open := <-wait:
		defer func() {
			if open && done != nil {
				close(done)
			}
		}()
		if !open {
			return nil, errors.New("rate limit access chanel closed")
		}

		resp, err := cb()
		if err != nil {
			return nil, err
		}
		resp.Header = NormalizeDiscordHeader(resp.StatusCode, resp.Header)

		// TODO: update bucket info

		return resp, nil
	}
}
