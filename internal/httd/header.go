package httd

import (
	"encoding/json"
	"errors"
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

// NormalizeDiscordHeaders cleans up the fields and remove fields that are considered noise.
//
//  eg. the header from the following response...
//    < HTTP/1.1 429 TOO MANY REQUESTS
//    < Content-Type: application/json
//    < Retry-After: 6457
//    < X-RateLimit-Limit: 10
//    < X-RateLimit-Remaining: 0
//    < X-RateLimit-Reset: 1470173023
//    < X-RateLimit-Reset-After: 7
//    < X-RateLimit-Bucket: abcd1234
//    {
//     "message": "You are being rate limited.",
//     "retry_after": 6457,
//     "global": false
//    }
//
//  is transformed into
//    < HTTP/1.1 429 TOO MANY REQUESTS
//    < Content-Type: application/json
//    < Retry-After: 6457
//    < X-RateLimit-Limit: 10
//    < X-RateLimit-Remaining: 0
//    < X-RateLimit-Bucket: abcd1234
func NormalizeDiscordHeader(statusCode int, header http.Header, body []byte) (h http.Header, err error) {
	// don't care about 2 different time delay estimates for the ltBucket reset.
	// So lets take Retry-After and X-RateLimit-Reset-After to set the reset
	var delay int64 = -1
	if retryAfter := header.Get(RateLimitRetryAfter); retryAfter != "" {
		// The number of milliseconds to wait before submitting another request.
		delay, _ = strconv.ParseInt(retryAfter, 10, 64)

		// we dont care about reset-after; let's stick to http standards
		// header.Set(XRateLimitResetAfter, strconv.FormatInt(delay, 10))
	} else if retry := header.Get(XRateLimitResetAfter); retry != "" {
		// Total time (in seconds with ms precision) of when the current rate limit bucket will reset.
		delayF, _ := strconv.ParseFloat(retry, 64)
		delayF *= 1000 // seconds => milliseconds
		delay = int64(delayF)
		header.Set(RateLimitRetryAfter, strconv.FormatInt(delay, 10))
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
		if delay == -1 && rateLimitBodyInfo.RetryAfter > 0 {
			delay = rateLimitBodyInfo.RetryAfter
			header.Set(RateLimitRetryAfter, strconv.FormatInt(delay, 10))
		}
	}

	// if delay is still not been set, try to extract it from XRateLimitReset date
	if delay == -1 {
		// Epoch time (seconds since 00:00:00 UTC on January 1, 1970) at which the rate limit resets
		// in our case Disgord requests milliseconds so we should not covert it

		// if reset := header.Get(XRateLimitReset); reset != "" {
		// 	epoch, _ := strconv.ParseFloat(reset, 64)
		// 	epoch *= 1000 // seconds => milliseconds
		// 	header.Set(XRateLimitReset, strconv.FormatInt(int64(epoch), 10))
		// }

		// if delay is not set, but Reset was defined; populate retry-after and reset-after
		if reset := header.Get(XRateLimitReset); reset != "" {
			if epoch, err := strconv.ParseFloat(reset, 64); err == nil && epoch > 0 {
				timestamp, err := HeaderToTime(header)
				if err != nil {
					// does add an delay, but there is no reason
					// to go insane if timestamp could not be handled
					timestamp = time.Now()
				}
				delayDuration := time.Unix(0, int64(time.Duration(int64(epoch))*time.Millisecond)).Sub(timestamp)
				delay = delayDuration.Milliseconds()

				// since the header date is in seconds, we add a 1s penalty to avoid 429
				delay += 1000

				header.Set(XRateLimitResetAfter, strconv.FormatInt(delay, 10))
				header.Set(RateLimitRetryAfter, strconv.FormatInt(delay, 10))
			}
		}
	}
	header.Set(DisgordNormalizedHeader, "true")

	// consider the following fields useless/noise
	header.Del(XRateLimitReset)
	header.Del(XRateLimitResetAfter)

	return header, nil
}
