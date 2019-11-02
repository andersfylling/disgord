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

//////////////////////////////////////////////////////
//
// QUEUE: Thread Safe
//
//////////////////////////////////////////////////////

func NewThreadSafeQueue() Queue {
	return &mutexQueue{}
}

type mutexQueue struct {
	sync.RWMutex
	queue
}

var _ Queue = (*mutexQueue)(nil)

func (q *mutexQueue) Pop() (interface{}, error) {
	q.Lock()
	defer q.Unlock()
	return q.queue.Pop()
}

func (q *mutexQueue) Push(x ...interface{}) {
	q.Lock()
	q.queue.Push(x...)
	q.Unlock()
}

func (q *mutexQueue) Len() uint {
	q.Lock()
	defer q.Unlock()
	return q.queue.Len()
}

//////////////////////////////////////////////////////
//
// QUEUE: Unsafe
//
//////////////////////////////////////////////////////

func NewThreadUnsafeQueue() Queue {
	return &queue{}
}

type queue struct {
	items []interface{}
}

var _ Queue = (*queue)(nil)

func (q *queue) Pop() (interface{}, error) {
	if len(q.items) == 0 {
		return nil, errors.New("mutexQueue is empty")
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
	q.items = append(q.items, x...)
}

func (q *queue) Len() uint {
	return uint(len(q.items))
}
