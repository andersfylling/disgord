package httd

import (
	"testing"
)

func TestLtBucket_AcquireLock(t *testing.T) {
	t.Run("already-locked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		bucket := newLeakyBucket(global)
		bucket.atomicLock.AcquireLock()

		if success := bucket.AcquireLock(); success {
			t.Error("should not be able to acquire lock on locked bucket")
		}
	})
	t.Run("unlocked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		bucket := newLeakyBucket(global)

		if success := bucket.AcquireLock(); !success {
			t.Error("should be able to lock unlocked bucket")
		}
	})
	t.Run("global-unlocked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		global.remaining = 2
		if !global.active() {
			t.Fatal("incorrectly configured global bucket")
		}
		bucket := newLeakyBucket(global)

		if success := bucket.AcquireLock(); !success {
			t.Error("should be able to lock unlocked bucket")
		}
	})
	t.Run("global-locked", func(t *testing.T) {
		global := newLeakyBucket(nil)
		global.atomicLock.AcquireLock()
		global.remaining = 2
		if !global.active() {
			t.Fatal("incorrectly configured global bucket")
		}

		bucket := newLeakyBucket(global)

		if success := bucket.AcquireLock(); success {
			t.Error("should be able to lock when global is locked bucket")
		}
	})
}
