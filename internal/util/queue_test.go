// +build !integration

package util

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewThreadSafeQueue()
	q.Push(true)

	l := q.Len()
	if l != 1 {
		t.Errorf("mutexQueue should contain 1 entry but holds %d", l)
	}

	_, _ = q.Pop()
	l = q.Len()
	if l != 0 {
		t.Errorf("mutexQueue should contain no entries but holds %d", l)
	}

	if _, err := q.Pop(); err == nil {
		t.Error("expected pop to fail on empty mutexQueue")
	}
}

func TestQueueOrder(t *testing.T) {
	q := NewThreadSafeQueue()

	iterate := func(cb func(int)) {
		for i := 0; i < 10; i++ {
			cb(i)
		}
	}

	// add
	iterate(func(i int) {
		q.Push(i)
	})

	// check order
	iterate(func(i int) {
		j, err := q.Pop()
		if err != nil {
			t.Error(err)
		}

		if i != j.(int) {
			t.Errorf("expected ordered outputs. Got %d, wants %d", j, i)
		}
	})

	q.Push(1, 2, 3)
	q.Push(4)
	q.Push(5)

	for i := 1; i <= 5; i++ {
		j, err := q.Pop()
		if err != nil {
			t.Error(err)
		}

		if i != j.(int) {
			t.Errorf("expected ordered outputs. Got %d, wants %d", j, i)
		}
	}
}
