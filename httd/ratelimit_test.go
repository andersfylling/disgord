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
	remaining := 4
	reset := int64(time.Now().UnixNano()) // just a large epoch ms

	resp := &http.Response{
		Header:     make(http.Header, 4),
		StatusCode: http.StatusTooManyRequests,
	}
	resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
	resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
	resp.Header.Set(XRateLimitReset, strconv.FormatInt(reset, 10))
	resp.Header.Set(XRateLimitGlobal, "true")

	rl := NewRateLimit()
	rl.UpdateRegisters("something", resp, []byte(""))

	if !rl.RateLimited("random") {
		t.Error("was not rate limited on a global scale")
	}
}
