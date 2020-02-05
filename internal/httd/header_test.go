package httd

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func equalDiscordHeaders(a, b http.Header) bool {
	if a.Get(RateLimitRetryAfter) != b.Get(RateLimitRetryAfter) {
		return false
	}
	if a.Get(XRateLimitLimit) != b.Get(XRateLimitLimit) {
		return false
	}
	if a.Get(XRateLimitRemaining) != b.Get(XRateLimitRemaining) {
		return false
	}
	if a.Get(XRateLimitReset) != b.Get(XRateLimitReset) {
		return false
	}
	if a.Get(XRateLimitResetAfter) != b.Get(XRateLimitResetAfter) {
		return false
	}
	if a.Get(XRateLimitBucket) != b.Get(XRateLimitBucket) {
		return false
	}
	if a.Get(XRateLimitGlobal) != b.Get(XRateLimitGlobal) {
		return false
	}
	return true
}

func TestNormalizeDiscordHeader(t *testing.T) {
	// < HTTP/1.1 429 TOO MANY REQUESTS
	// < Content-Type: application/json
	// < Retry-After: 6457
	// < X-RateLimit-Limit: 10
	// < X-RateLimit-Remaining: 0
	// < X-RateLimit-Bucket: abcd1234
	normalized := http.Header{}
	normalized.Set(RateLimitRetryAfter, "6457")
	normalized.Set(XRateLimitLimit, "10")
	normalized.Set(XRateLimitRemaining, "0")
	normalized.Set(XRateLimitBucket, "abcd1234")

	// < HTTP/1.1 429 TOO MANY REQUESTS
	// < Content-Type: application/json
	// < X-RateLimit-Limit: 10
	// < X-RateLimit-Remaining: 1
	// < X-RateLimit-Bucket: abcd1234
	normalized2 := http.Header{}
	normalized2.Set(XRateLimitLimit, "10")
	normalized2.Set(XRateLimitRemaining, "0")
	normalized2.Set(XRateLimitBucket, "abcd1234")

	now := time.Now()
	headerDate := now.Format(time.RFC1123)
	epochMSNow := now.UnixNano() / int64(time.Millisecond)
	retryAfterMS := int64(6457)

	t.Run("a", func(t *testing.T) {
		// < HTTP/1.1 429 TOO MANY REQUESTS
		// < Content-Type: application/json
		// < Retry-After: 6457
		// < X-RateLimit-Limit: 10
		// < X-RateLimit-Remaining: 0
		// < X-RateLimit-Reset: ${time.Now().Milliseconds()+6457ms}
		// < X-RateLimit-Reset-After: 6.457
		// < X-RateLimit-Bucket: abcd1234
		header := http.Header{}
		header.Set("date", headerDate)
		header.Set(RateLimitRetryAfter, "6457")
		header.Set(XRateLimitLimit, "10")
		header.Set(XRateLimitRemaining, "0")
		header.Set(XRateLimitReset, strconv.FormatInt(epochMSNow+retryAfterMS, 10))
		header.Set(XRateLimitResetAfter, "6.457")
		header.Set(XRateLimitBucket, "abcd1234")

		transformed, err := NormalizeDiscordHeader(http.StatusTooManyRequests, header, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !equalDiscordHeaders(transformed, normalized) {
			t.Fatal(fmt.Sprintf("expected \n%#v\n Got \n%#v", transformed, normalized))
		}
	})

	t.Run("b", func(t *testing.T) {
		// < HTTP/1.1 429 TOO MANY REQUESTS
		// < Content-Type: application/json
		// < X-RateLimit-Limit: 10
		// < X-RateLimit-Remaining: 0
		// < X-RateLimit-Reset: ${time.Now().Milliseconds()+6457ms}
		// < X-RateLimit-Reset-After: 6.457
		// < X-RateLimit-Bucket: abcd1234
		header := http.Header{}
		header.Set("date", headerDate)
		header.Set(XRateLimitLimit, "10")
		header.Set(XRateLimitRemaining, "0")
		header.Set(XRateLimitReset, strconv.FormatInt(epochMSNow+retryAfterMS, 10))
		header.Set(XRateLimitResetAfter, "6.457")
		header.Set(XRateLimitBucket, "abcd1234")

		transformed, err := NormalizeDiscordHeader(http.StatusTooManyRequests, header, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !equalDiscordHeaders(transformed, normalized) {
			t.Fatal(fmt.Sprintf("expected \n%#v\n Got \n%#v", transformed, normalized))
		}
	})

	t.Run("c.1", func(t *testing.T) {
		// < HTTP/1.1 429 TOO MANY REQUESTS
		// < Content-Type: application/json
		// < X-RateLimit-Limit: 10
		// < X-RateLimit-Remaining: 0
		// < X-RateLimit-Reset: ${time.Now().Milliseconds()+6457ms} + desync
		// < X-RateLimit-Bucket: abcd1234
		header := http.Header{}
		header.Set("date", now.Add(234*time.Millisecond).Format(time.RFC1123))
		header.Set(XRateLimitLimit, "10")
		header.Set(XRateLimitRemaining, "0")
		header.Set(XRateLimitReset, strconv.FormatInt(epochMSNow+retryAfterMS, 10))
		header.Set(XRateLimitBucket, "abcd1234")

		transformed, err := NormalizeDiscordHeader(http.StatusTooManyRequests, header, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !equalDiscordHeaders(transformed, normalized) {
			t.Fatal(fmt.Sprintf("expected \n%#v\n Got \n%#v", transformed, normalized))
		}
	})

	t.Run("c.2", func(t *testing.T) {
		// < HTTP/1.1 429 TOO MANY REQUESTS
		// < Content-Type: application/json
		// < X-RateLimit-Limit: 10
		// < X-RateLimit-Remaining: 0
		// < X-RateLimit-Reset: ${time.Now().Milliseconds()+6457ms} + (-desync)
		// < X-RateLimit-Bucket: abcd1234
		header := http.Header{}
		header.Set("date", now.Add(-1*234*time.Millisecond).Format(time.RFC1123))
		header.Set(XRateLimitLimit, "10")
		header.Set(XRateLimitRemaining, "0")
		header.Set(XRateLimitReset, strconv.FormatInt(epochMSNow+retryAfterMS, 10))
		header.Set(XRateLimitBucket, "abcd1234")

		transformed, err := NormalizeDiscordHeader(http.StatusTooManyRequests, header, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !equalDiscordHeaders(transformed, normalized) {
			t.Fatal(fmt.Sprintf("expected \n%#v\n Got \n%#v", transformed, normalized))
		}
	})

	t.Run("d", func(t *testing.T) {
		// < HTTP/1.1 429 TOO MANY REQUESTS
		// < Content-Type: application/json
		// < X-RateLimit-Limit: 10
		// < X-RateLimit-Remaining: 0
		// < X-RateLimit-Reset-After: 6.457
		// < X-RateLimit-Bucket: abcd1234
		header := http.Header{}
		header.Set("date", headerDate)
		header.Set(XRateLimitLimit, "10")
		header.Set(XRateLimitRemaining, "0")
		header.Set(XRateLimitResetAfter, "6.457")
		header.Set(XRateLimitBucket, "abcd1234")

		transformed, err := NormalizeDiscordHeader(http.StatusTooManyRequests, header, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !equalDiscordHeaders(transformed, normalized) {
			t.Fatal(fmt.Sprintf("expected \n%#v\n Got \n%#v", transformed, normalized))
		}
	})

	t.Run("e", func(t *testing.T) {
		// < HTTP/1.1 429 TOO MANY REQUESTS
		// < Content-Type: application/json
		// < Retry-After: 6457
		// < X-RateLimit-Limit: 10
		// < X-RateLimit-Remaining: 1
		// < X-RateLimit-Reset: ${time.Now().Milliseconds()+6457ms}
		// < X-RateLimit-Reset-After: 6.457
		// < X-RateLimit-Bucket: abcd1234
		header := http.Header{}
		header.Set("date", headerDate)
		header.Set(RateLimitRetryAfter, "6457")
		header.Set(XRateLimitLimit, "10")
		header.Set(XRateLimitRemaining, "0")
		header.Set(XRateLimitReset, strconv.FormatInt(epochMSNow+retryAfterMS, 10))
		header.Set(XRateLimitResetAfter, "6.457")
		header.Set(XRateLimitBucket, "abcd1234")

		transformed, err := NormalizeDiscordHeader(http.StatusTooManyRequests, header, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !equalDiscordHeaders(transformed, normalized) {
			t.Fatal(fmt.Sprintf("expected \n%#v\n Got \n%#v", transformed, normalized))
		}
	})

}
