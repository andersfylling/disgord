package disgord

import (
	"errors"
	"sync"
)

// cache keys
const (
	UserCache = iota
	ChannelCache
	GuildCache
)

type Cacher interface {
	Update(key int, v interface{}) (err error)
	Get(key int, id Snowflake) (v interface{}, err error)
}

func NewErrorCacheItemNotFound(id Snowflake) *ErrorCacheItemNotFound {
	return &ErrorCacheItemNotFound{
		info: "item with id{" + id.String() + "} was not found in cache",
	}
}

type ErrorCacheItemNotFound struct {
	info string
}

func (e *ErrorCacheItemNotFound) Error() string {
	return e.info
}

func NewCache(conf *CacheConfig) *Cache {
	return &Cache{
		conf:  conf,
		users: make(map[Snowflake]*User),
	}
}

type CacheConfig struct {
	Immutable bool
}

type Cache struct {
	conf       *CacheConfig
	users      map[Snowflake]*User
	usersMutex sync.RWMutex
}

func (c *Cache) Updates(key int, vs []interface{}) (err error) {
	for _, v := range vs {
		err = c.Update(key, v)
		if err != nil {
			return
		}
	}

	return
}

func (c *Cache) Update(key int, v interface{}) (err error) {
	if v == nil {
		err = errors.New("object was nil")
		return
	}

	var obj interface{}
	if c.conf.Immutable {
		if copyable, ok := v.(DeepCopier); ok {
			obj = copyable.DeepCopy()
		} else {
			err = errors.New("object does not implement DeepCopier and must do so when config.Immutable is set")
		}
	} else {
		obj = v
	}

	switch key {
	case UserCache:
		if user, isUser := obj.(*User); isUser {
			c.SetUser(user)
		} else {
			err = errors.New("can only save *User structures to user cache")
		}
	default:
		err = errors.New("caching for given type is not yet implemented")
	}

	return
}

func (c *Cache) Get(key int, id Snowflake) (v interface{}, err error) {
	switch key {
	case UserCache:
		v, err = c.GetUser(id)
	default:
		err = errors.New("caching for given type is not yet implemented")
	}

	return
}

func (c *Cache) SetUser(new *User) {
	if new == nil {
		return
	}

	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()

	if user, exists := c.users[new.ID]; exists {
		new.CopyOverTo(user)
	} else {
		c.users[new.ID] = new
	}
}

func (c *Cache) GetUser(id Snowflake) (user *User, err error) {
	c.usersMutex.RLock()
	defer c.usersMutex.RUnlock()

	var exists bool
	var result *User
	if result, exists = c.users[id]; !exists {
		err = NewErrorCacheItemNotFound(id)
		return
	}

	if c.conf.Immutable {
		user = result.DeepCopy().(*User)
	} else {
		user = result
	}

	return
}
