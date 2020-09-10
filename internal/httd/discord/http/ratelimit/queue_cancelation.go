package ratelimit

import (
	"context"
	"sync"

	"go.uber.org/atomic"
)

func deadContext(ctx context.Context) bool {
	return ctx.Err() != nil || ctx.Done() == nil
}

type DoneChan chan interface{}

type queueItem struct {
	ctx context.Context
	notify chan (chan interface{})
}

type queueWithCancellation struct {
	mu sync.Mutex
	items []*queueItem
	change chan int
	dead atomic.Bool
}

func (queue *queueWithCancellation) insert(item *queueItem) {
	queue.mu.Lock()
	if queue.dead.Load() {
		queue.mu.Unlock()
		close(item.notify)
		return
	}
	defer func() {
		queue.mu.Unlock()
		queue.change <- 1
	}()

	if cap(queue.items) == len(queue.items) {
		length := len(queue.items)
		tmp := queue.items
		queue.items = make([]*queueItem, length, length*2)
		copy(queue.items, tmp)
	}
	queue.items = append(queue.items, item)
}

func (queue *queueWithCancellation) close() {
	queue.dead.Store(true)
	close(queue.change)
}

func (queue *queueWithCancellation) watch() {
	var debt uint
	var lastChan chan interface{}

	popBack := func() *queueItem {
		queue.mu.Lock()
		data := queue.items[len(queue.items) - 1]
		queue.items = queue.items[:len(queue.items) - 1]
		queue.mu.Unlock()

		return data
	}

	popHandle := func() bool {
		data := popBack()
		if deadContext(data.ctx) {
			return false
		}

		// close(lastChan) // nop
		lastChan = make(chan interface{})
		data.notify <- lastChan
		return true
	}

	reactOnNew := true
	for {
		select {
		case _, open := <-queue.change:
			// assume every change is just a increment of one
			if !open {
				queue.mu.Lock()
				for i := 0; i < len(queue.items); i++ {
					close(queue.items[i].notify)
				}
				queue.mu.Unlock()
				break
			}

			if reactOnNew {
				if popHandle() {
					reactOnNew = false
				}
			} else {
				debt++
			}
		case <-lastChan:
			if debt == 0 {
				lastChan = make(chan interface{})
				reactOnNew = true
			} else {
				popHandle()
				debt--
			}
		}
	}
}