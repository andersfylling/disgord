// +build !integration

package httd

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"testing"
// 	"time"
// )
//
// func createErrMsg(t string, a, b int) error {
// 	return fmt.Errorf("%s is incorrect. got %d, wants %d", t, a, b)
// }
//
// func TestExtractRateLimitInfo(t *testing.T) {
// 	limit := 2
// 	remaining := 4
// 	reset := int64(456344)
//
// 	resp := &http.Response{
// 		Header: make(http.Header, 3),
// 	}
// 	resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
// 	resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
// 	resp.Header.Set(XRateLimitReset, strconv.FormatInt(reset, 10))
// 	resp.Header.Set("date", time.Now().Format(time.RFC1123))
//
// 	info, err := ExtractRateLimitInfo(resp, []byte(""))
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	verifyRateLimitInfo(t, info, limit, remaining, reset)
// 	if info.Global {
// 		t.Error("rate limit is registered as global even though it is not")
// 	}
// }
//
// func TestExtractRateLimitInfoGlobal(t *testing.T) {
// 	limit := 2
// 	remaining := 4
// 	reset := int64(time.Now().Nanosecond() / 100) // just a large epoch ms
//
// 	resp := &http.Response{
// 		Header: make(http.Header, 4),
// 	}
// 	resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
// 	resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
// 	resp.Header.Set(XRateLimitReset, strconv.FormatInt(reset, 10))
// 	resp.Header.Set(XRateLimitGlobal, "true")
// 	resp.Header.Set("date", time.Now().Format(time.RFC1123))
//
// 	info, err := ExtractRateLimitInfo(resp, []byte(""))
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	verifyRateLimitInfo(t, info, limit, remaining, reset)
// 	if !info.Global {
// 		t.Error("rate limit is not registered as global even though it is global")
// 	}
// }
//
// func TestExtractRateLimitGlobal(t *testing.T) {
// 	limit := 2
// 	remaining := 0
// 	reset := int64(235345325435) // just a large epoch ms
//
// 	resp := &http.Response{
// 		Header:     make(http.Header, 4),
// 		StatusCode: http.StatusTooManyRequests,
// 	}
// 	resp.Header.Set(XRateLimitLimit, strconv.Itoa(limit))
// 	resp.Header.Set(XRateLimitRemaining, strconv.Itoa(remaining))
// 	resp.Header.Set(XRateLimitReset, strconv.FormatInt(reset, 10))
// 	resp.Header.Set(XRateLimitGlobal, "true")
// 	resp.Header.Set("date", time.Now().Format(time.RFC1123))
//
// 	rl := NewRateLimit()
// 	rl.UpdateRegisters("something", resp, []byte(""))
//
// 	if !rl.RateLimited("random") {
// 		t.Error("was not rate limited on a global scale")
// 	}
// }
