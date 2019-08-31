package metrics

import (
	"sync"
	"time"
)

type TimestampLogger interface {
	LogTimestamp()
}

func NewTimestampLogger(limit uint) TimestampLogger {
	return &timestampStack{
		tms:   make([]time.Time, limit),
		limit: limit,
	}
}

type timestampStack struct {
	sync.RWMutex
	tms   []time.Time
	index uint
	limit uint
}

var _ TimestampLogger = (*timestampStack)(nil)

func (s *timestampStack) shift() {
	offset := uint(float32(s.index) * 0.2)
	for i := uint(0); i < s.index; i++ {
		s.tms[i] = s.tms[offset]
		offset++
	}
	s.tms = s.tms[:s.index-offset]
	s.index -= offset
}

func (s *timestampStack) LogTimestamp() {
	s.Lock()
	if s.index == s.limit {
		s.shift()
	}
	s.tms[s.index] = time.Now()
	s.index++
	s.Unlock()
}
