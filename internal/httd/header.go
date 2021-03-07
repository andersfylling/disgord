package httd

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/andersfylling/disgord/json"
)

// http rate limit identifiers
const (
	XAuditLogReason         = "X-Audit-Log-Reason"
	XRateLimitPrecision     = "X-RateLimit-Precision"
	XRateLimitBucket        = "X-RateLimit-Bucket"
	XRateLimitLimit         = "X-RateLimit-Limit"
	XRateLimitRemaining     = "X-RateLimit-Remaining"
	XRateLimitReset         = "X-RateLimit-Reset"
	XRateLimitResetAfter    = "X-RateLimit-Reset-After"
	XRateLimitGlobal        = "X-RateLimit-Global"
	RateLimitRetryAfter     = "Retry-After"
	DisgordNormalizedHeader = "X-Disgord-Normalized-Kufdsfksduhf-S47yf"
	XDisgordNow             = "X-Disgord-Now-fsagkhf"
)

// HeaderToTime takes the response header from Discord and extracts the
// timestamp. Useful for detecting time desync between discord and client
func HeaderToTime(header http.Header) (t time.Time, err error) {
	// date: Fri, 14 Sep 2018 19:04:24 GMT
	dateStr := header.Get("date")
	if dateStr == "" {
		err = errors.New("missing header field 'date'")
		return
	}

	t, err = time.Parse(time.RFC1123, dateStr)
	return
}

type RateLimitResponseStructure struct {
	Message    string  `json:"message"`     // A message saying you are being rate limited.
	RetryAfter float64 `json:"retry_after"` // The number of seconds to wait before submitting another request.
	Global     bool    `json:"global"`      // A value indicating if you are being globally rate limited or not
}

// NormalizeDiscordHeader overrides header fields with body content and make sure every header field
// uses milliseconds and not seconds. Regards rate limits only.
func NormalizeDiscordHeader(statusCode int, header http.Header, body []byte) (h http.Header, err error) {
	secondsToMilli := func(s float64) int64 {
		s *= 1000
		return int64(s)
	}

	var now time.Time
	if field := header.Get(XDisgordNow); field != "" {
		n, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			now = time.Now()
		} else {
			now = time.Unix(0, int64(time.Duration(n)*time.Millisecond))
		}
	} else {
		now = time.Now()
	}

	// don't care about 2 different time delay estimates for the ltBucket reset.
	// So lets take Retry-After and X-RateLimit-Reset-After to set the reset
	var delay int64
	if retry := header.Get(XRateLimitResetAfter); delay == 0 && retry != "" {
		delayF, _ := strconv.ParseFloat(retry, 64)
		delay = secondsToMilli(delayF)
	}

	// sometimes the body might be populated too
	if delay == 0 && statusCode == http.StatusTooManyRequests && body != nil {
		var rateLimitBodyInfo *RateLimitResponseStructure
		if err = json.Unmarshal(body, &rateLimitBodyInfo); err != nil {
			return nil, err
		}
		if rateLimitBodyInfo.Global {
			header.Set(XRateLimitGlobal, "true")
		}
		if delay == 0 && rateLimitBodyInfo.RetryAfter > 0 {
			delay = secondsToMilli(rateLimitBodyInfo.RetryAfter)
		}
	}

	// convert reset to store milliseconds and not seconds
	// if there is no content, we create a reset unix using the delay
	if reset := header.Get(XRateLimitReset); reset != "" {
		if delay == 0 {
			epoch, _ := strconv.ParseFloat(reset, 64)
			epochMilli := secondsToMilli(epoch)
			header.Set(XRateLimitReset, strconv.FormatInt(epochMilli, 10))
		} else {
			epochNow := now.UnixNano() / int64(time.Millisecond)
			header.Set(XRateLimitReset, strconv.FormatInt(epochNow+delay, 10))
		}
	} else if delay > 0 {
		timestamp, err := HeaderToTime(header)
		if err != nil {
			// does add an delay, but there is no reason
			// to go insane if timestamp could not be handled
			timestamp = now
		}

		reset := timestamp.Add(time.Duration(delay) * time.Millisecond)
		ms := reset.UnixNano() / int64(time.Millisecond)
		header.Set(XRateLimitReset, strconv.FormatInt(ms, 10))
	}

	header.Set(DisgordNormalizedHeader, "true")
	return header, nil
}
