package httd

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func createErrMsg(t string, a, b int) error {
	return fmt.Errorf("%s is incorrect. got %d, wants %d", t, a, b)
}

func verifyRateLimitInfo(t *testing.T, info *RateLimitInfo, limit, remaining int, reset int64) {
	if info.Limit != limit {
		t.Error(createErrMsg("limit", info.Limit, limit))
	}
	if info.Remaining != remaining {
		t.Error(createErrMsg("remaining", info.Remaining, remaining))
	}
	if info.Reset != reset*1000 {
		t.Error(createErrMsg("reset", int(info.Reset), int(reset)*1000))
	}
}

func TestDiscordTimeDiff(t *testing.T) {
	t.Run("calculating 0 offset", func(t *testing.T) {
		diff := NewDiscordTimeDiff()
		now := time.Now()
		diff.Update(now, now)

		if diff.offset.Nanoseconds()/1000 != 0 {
			t.Errorf("offset incorrect. Wants 0, got %s", diff.offset.String())
		}
	})

	t.Run("offset, discord behind", func(t *testing.T) {
		diff := NewDiscordTimeDiff()
		var behind int64
		var discord time.Time
		var now time.Time

		now = time.Now()
		behind = now.Unix() - 200
		discord = time.Unix(behind, 0)
		diff.Update(now, discord)

		if diff.offset.Seconds() > -200 || diff.offset.Seconds() < -205 {
			t.Errorf("offset incorrect. Wants [-200s, -205s), got %s", diff.offset.String())
		}
	})

	t.Run("offset, local behind", func(t *testing.T) {
		diff := NewDiscordTimeDiff()
		var behind int64
		var local time.Time
		var now time.Time

		now = time.Now()
		behind = now.Unix() - 200
		local = time.Unix(behind, 0)
		diff.Update(local, now)

		if diff.offset.Seconds() < 200 || diff.offset.Seconds() > 205 {
			t.Errorf("offset incorrect. Wants [200s, 205s), got %s", diff.offset.String())
		}
	})
}

func TestExtractRateLimitInfo(t *testing.T) {
	limit := 2
	remaining := 4
	reset := int64(456344)

	resp := &http.Response{
		Header: make(http.Header, 3),
	}
	resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
	resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
	resp.Header.Set(XRateLimitReset, strconv.FormatInt(reset, 10))
	resp.Header.Set("date", time.Now().Format(time.RFC1123))

	info, err := ExtractRateLimitInfo(resp, []byte(""))
	if err != nil {
		t.Error(err)
	}

	verifyRateLimitInfo(t, info, limit, remaining, reset)
	if info.Global {
		t.Error("rate limit is registered as global even though it is not")
	}
}

func TestExtractRateLimitInfoGlobal(t *testing.T) {
	limit := 2
	remaining := 4
	reset := int64(time.Now().Nanosecond() / 100) // just a large epoch ms

	resp := &http.Response{
		Header: make(http.Header, 4),
	}
	resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
	resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
	resp.Header.Set(XRateLimitReset, strconv.FormatInt(reset, 10))
	resp.Header.Set(XRateLimitGlobal, "true")
	resp.Header.Set("date", time.Now().Format(time.RFC1123))

	info, err := ExtractRateLimitInfo(resp, []byte(""))
	if err != nil {
		t.Error(err)
	}

	verifyRateLimitInfo(t, info, limit, remaining, reset)
	if !info.Global {
		t.Error("rate limit is not registered as global even though it is global")
	}
}

func TestExtractRateLimitGlobal(t *testing.T) {
	limit := 2
	remaining := 0
	reset := int64(235345325435) // just a large epoch ms

	resp := &http.Response{
		Header:     make(http.Header, 4),
		StatusCode: http.StatusTooManyRequests,
	}
	resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
	resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
	resp.Header.Set(XRateLimitReset, strconv.FormatInt(reset, 10))
	resp.Header.Set(XRateLimitGlobal, "true")
	resp.Header.Set("date", time.Now().Format(time.RFC1123))

	rl := NewRateLimit()
	rl.UpdateRegisters("something", nil, resp, []byte(""))

	if !rl.RateLimited("random") {
		t.Error("was not rate limited on a global scale")
	}
}
