// +build !integration

package gateway

import (
	"testing"
	"time"
)

func TestRlBucket(t *testing.T) {
	t.Run("blocked", func(t *testing.T) {
		b1 := newRatelimitBucket(1, 0)
		if b1.Blocked() {
			t.Error("expected first request to not be blocked")
		}

		b1.entries[0].unix = time.Now().UnixNano() + 1*int64(time.Second)
		if !b1.Blocked() {
			t.Error("should be blocked")
		}

		b1.entries[0].unix = time.Now().UnixNano() - 1
		if b1.Blocked() {
			t.Error("should not be blocked")
		}

	})
	t.Run("insert", func(t *testing.T) {
		b1 := newRatelimitBucket(1, 0)
		if b1.entries[0].unix > 0 {
			t.Error("expects unix to be 0")
		}
		b1.Insert("test")
		if b1.entries[0].cmd != "test" {
			t.Error("expected first entry to have cmd `test`")
		}
		b1.Insert("lol")
		if b1.entries[0].cmd != "lol" {
			t.Error("expected first entry to have cmd `lol`")
		}
		if len(b1.entries) > 1 {
			t.Error("length should be 1")
		}

		b2 := newRatelimitBucket(5, 0)
		cmds := []string{
			"1", "2", "3", "4", "5",
		}

		for i := len(cmds) - 1; i >= 0; i-- {
			b2.Insert(cmds[i])
		}
		for i := range b2.entries {
			if b2.entries[i].cmd != cmds[i] {
				t.Errorf("b2.entries is not correctly ordered. At index %d, got %s, wants %s", i, b2.entries[i].cmd, cmds[i])
			}
		}

	})
}
