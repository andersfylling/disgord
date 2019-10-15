package util

import (
	"sync/atomic"
)

type AtomicLock int32

func (l *AtomicLock) AcquireLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(l), 0, 1)
}

func (l *AtomicLock) IsLocked() bool {
	return atomic.LoadInt32((*int32)(l)) == 1
}

func (l *AtomicLock) Unlock() {
	atomic.CompareAndSwapInt32((*int32)(l), 1, 0)
}
