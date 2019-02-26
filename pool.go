package disgord

import "sync"

type Pool interface {
	Put(x Reseter)
	Get() (x Reseter)
}

type pool struct {
	pool sync.Pool

	// New specifies a function to generate a
	// value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func() Reseter
}

// Put resets the object before it is put back into the pool. We reset it here
// to quickly detect if there are other owners than the pool as it is inserted.
func (p *pool) Put(x Reseter) {
	x.Reset()
	p.pool.Put(x)
}

// Get selects an arbitrary item from the Pool, removes it from the
// Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to Put and
// the values returned by Get.
//
// This assumes that p.New is always set.
func (p *pool) Get() (x Reseter) {
	var ok bool
	if x, ok = p.pool.Get().(Reseter); x == nil || !ok {
		x = p.New()
	}

	return
}

//////////////////////////////////////////////////////
//
// Resource Pools
//
//////////////////////////////////////////////////////

func newPools() *pools {
	p := &pools{}

	p.channel = &pool{
		New: func() Reseter {
			return &Channel{}
		},
	}
	p.user = &pool{
		New: func() Reseter {
			return &User{}
		},
	}
	p.message = &pool{
		New: func() Reseter {
			return &Message{
				Author: p.user.Get().(*User),
			}
		},
	}
	p.emoji = &pool{
		New: func() Reseter {
			return &Emoji{}
		},
	}

	//p.msgCreate = &pool{
	//	New: func() Reseter {
	//		return &MessageCreate{
	//			Message: p.message.Get().(*Message),
	//		}
	//	},
	//}

	return p
}

type pools struct {
	channel Pool
	message Pool
	user    Pool
	emoji   Pool

	// events
	//msgCreate Pool // this is actually slower
}

func (p *pools) ChannelPool() Pool {
	return p.channel
}
func (p *pools) MessagePool() Pool {
	return p.message
}
func (p *pools) UserPool() Pool {
	return p.user
}
func (p *pools) EmojiPool() Pool {
	return p.emoji
}
