package httd

import (
	"errors"
	"github.com/andersfylling/disgord/json"
	"net/http"
	"strconv"
	"time"
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
	Message    string `json:"message"`     // A message saying you are being rate limited.
	RetryAfter int64  `json:"retry_after"` // The number of milliseconds to wait before submitting another request.
	Global     bool   `json:"global"`      // A value indicating if you are being globally rate limited or not
}

// NormalizeDiscordHeader overrides header fields with body content and make sure every header field
// uses milliseconds and not seconds. Regards rate limits only.
func NormalizeDiscordHeader(statusCode int, header http.Header, body []byte) (h http.Header, err error) {
	// don't care about 2 different time delay estimates for the ltBucket reset.
	// So lets take Retry-After and X-RateLimit-Reset-After to set the reset
	var delay int64
	if retryAfter := header.Get(RateLimitRetryAfter); retryAfter != "" {
		delay, _ = strconv.ParseInt(retryAfter, 10, 64)
	}
	if retry := header.Get(XRateLimitResetAfter); delay == 0 && retry != "" {
		delayF, _ := strconv.ParseFloat(retry, 64)
		delayF *= 1000 // seconds => milliseconds
		delay = int64(delayF)
	}

	// sometimes the body might be populated too
	if statusCode == http.StatusTooManyRequests && body != nil {
		var rateLimitBodyInfo *RateLimitResponseStructure
		if err = json.Unmarshal(body, &rateLimitBodyInfo); err != nil {
			return nil, err
		}
		if rateLimitBodyInfo.Global {
			header.Set(XRateLimitGlobal, "true")
		}
		if delay == 0 && rateLimitBodyInfo.RetryAfter > 0 {
			delay = rateLimitBodyInfo.RetryAfter
		}
	}

	// convert Reset to store milliseconds and not seconds
	// if there is no content, we create a Reset unix using the delay
	if reset := header.Get(XRateLimitReset); reset != "" {
		epoch, _ := strconv.ParseFloat(reset, 64)
		epoch *= 1000 // seconds => milliseconds
		header.Set(XRateLimitReset, strconv.FormatInt(int64(epoch), 10))
	} else if delay > 0 {
		timestamp, err := HeaderToTime(header)
		if err != nil {
			// does add an delay, but there is no reason
			// to go insane if timestamp could not be handled
			timestamp = time.Now()
		}

		reset := timestamp.Add(time.Duration(delay) * time.Millisecond)
		ms := reset.UnixNano() / int64(time.Millisecond)
		header.Set(XRateLimitReset, strconv.FormatInt(ms, 10))
	}

	header.Set(DisgordNormalizedHeader, "true")
	return header, nil
}
