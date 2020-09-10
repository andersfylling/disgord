package ratelimit

import (
	"context"
	"sync"

	"go.uber.org/atomic"
)

// YOU MUST INSTANTIATE homeChan!
type passingGame struct {
	homeChan chan *ball
	once sync.Once
	done atomic.Bool
	queue []*player
	sync.Mutex
}

func (p *passingGame) add(player *player) {
	p.Lock()
	p.queue = append(p.queue, player)
	p.Unlock()

	p.once.Do(func() {
		p.homeChan <- &ball{p.homeChan}
	})
}

func (p *passingGame) run(ctx context.Context) {
	if !p.done.CAS(false, true) {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.homeChan:

		}
	}
}

type player struct {
	c chan<- *ball
	ctx context.Context
}

func (p *player) skip() bool {
	return p.ctx.Err() != nil
}

type ball struct {
	home chan *ball
}

func (b *ball) next() {
	b.home <- b
}