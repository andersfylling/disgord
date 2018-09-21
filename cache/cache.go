package cahce

import (
	. "github.com/andersfylling/snowflake"
	"errors"
	"sync"
)

type DeepCopy interface {
	DeepCopy() (v interface{})
}

// Rewriter in stead of creating a new copy, this allows rewriting object already in memory
// to reduce memory usage, speed up updating object in cache, and reduce GC usage
type Rewriter interface {
	// Rewrite handles mutex locking internally
	Rewrite(new interface{}) (err error)
}

type Cacher interface {
	Set(id Snowflake, new interface{}) (err error)
	Add(id Snowflake, new interface{}) (err error)
	Get(id Snowflake) (v interface{}, err error)
}

type Config struct {
	UseCopies bool
}

type Client struct {
	config *Config
	Items map[Snowflake]interface{}
	mu sync.RWMutex
}

func (c *Client) Set(id Snowflake, new interface{}) (err error) {
	if old, exists := c.Get(id); exists {
		if c.config.UseCopies {
			old.(Rewriter).Rewrite(new)
		} else {
			c.mu.Lock()
			c.Items[id] = new
			c.mu.Unlock()
		}
	} else {
		err = c.Add(id, new)
	}

	return
}

func (c *Client) Add(id Snowflake, new interface{}) (err error) {
	if _, exists := c.Get(id); exists {
		err = errors.New("an item already exist with given key: " + id.String())
	} else {
		c.mu.Lock()
		c.Items[id] = new
		c.mu.Unlock()
	}

	return
}

func (c *Client) Get(id Snowflake) (v interface{}, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.config.UseCopies {
		var p interface{}
		p, exists = c.Items[id]
		v = p.(DeepCopy).DeepCopy()
	} else {
		v, exists = c.Items[id]
	}

	return
}