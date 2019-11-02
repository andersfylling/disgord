package util

import (
	"errors"
	"sync"
)

type Queue interface {
	Pop() (interface{}, error)
	Push(...interface{})
	Len() uint
}

func NewQueue() Queue {
	return &queue{}
}

type queue struct {
	sync.RWMutex
	items []interface{}
}

var _ Queue = (*queue)(nil)

func (q *queue) Pop() (interface{}, error) {
	q.Lock()
	defer q.Unlock()

	if len(q.items) == 0 {
		return nil, errors.New("queue is empty")
	}

	// get first
	x := q.items[0]

	// shift
	for i := 0; i < len(q.items)-1; i++ {
		q.items[i] = q.items[i+1]
	}
	q.items = q.items[:len(q.items)-1]

	return x, nil
}

func (q *queue) Push(x ...interface{}) {
	q.Lock()
	q.items = append(q.items, x...)
	q.Unlock()
}

func (q *queue) Len() uint {
	q.Lock()
	defer q.Unlock()
	return uint(len(q.items))
}
