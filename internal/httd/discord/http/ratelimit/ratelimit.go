package ratelimit

// DiscordRateLimitCompliant is a way to test that the injected RoundTripper implementation support discord rate limits.
// TODO: is there a way to use a predefined request to verify this instead?
type DiscordRateLimitCompliant interface {
	DiscordRateLimitCompliant()
}