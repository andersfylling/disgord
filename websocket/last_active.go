package websocket

import (
	"sync"
	"time"
)

// this is poisen..
type lastActivity struct {
	mu sync.Mutex
	tm time.Time
}

func (l *lastActivity) Update() {
	l.mu.Lock()
	l.tm = time.Now()
	l.mu.Unlock()
}

func (l *lastActivity) Time() time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tm
}

func (l *lastActivity) OlderThan(d time.Duration) (older bool) {
	l.mu.Lock()
	now := time.Now()
	older = now.Sub(l.tm) > d
	l.mu.Unlock()
	return older
}
