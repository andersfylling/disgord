package disgord

import (
	"errors"
	"time"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/cache/lru"
)

// cache keys
const (
	UserCache = iota
	ChannelCache
	GuildCache

	CacheAlg_LRU = "lru"
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
	const userWeight = 1 // MiB. TODO: what is the actual max size?
	limit := conf.UserCacheLimitMiB / userWeight

	var userCacheAlg interfaces.CacheAlger
	switch conf.UserCacheAlgorithm {
	case CacheAlg_LRU:
		fallthrough
	default:
		userCacheAlg = lru.NewCacheList(limit, conf.UserCacheLifetime, conf.UserCacheUpdateLifetimeOnUsage)
	}

	return &Cache{
		conf:  conf,
		users: userCacheAlg,
	}
}

type CacheConfig struct {
	Immutable bool

	UserCacheLimitMiB              uint
	UserCacheLifetime              time.Duration
	UserCacheUpdateLifetimeOnUsage bool
	UserCacheAlgorithm             string
}

type Cache struct {
	conf  *CacheConfig
	users interfaces.CacheAlger
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

	_, implementsDeepCopier := v.(DeepCopier)
	_, implementsCacheCopier := v.(cacheCopier)
	if !implementsCacheCopier && !implementsDeepCopier && c.conf.Immutable {
		err = errors.New("object does not implement DeepCopier & cacheCopier and must do so when config.Immutable is set")
		return
	}

	switch key {
	case UserCache:
		if user, isUser := v.(*User); isUser {
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

	c.users.Lock()
	defer c.users.Unlock()
	if item, exists := c.users.Get(new.ID); exists {
		if c.conf.Immutable {
			new.copyOverToCache(item.Object())
		} else {
			item.Set(new)
		}
		c.users.Set(new.ID, item)
	} else {
		var content interface{}
		if c.conf.Immutable {
			content = new.DeepCopy()
		} else {
			content = new
		}
		c.users.Set(new.ID, lru.NewCacheItem(content))
	}
}

func (c *Cache) GetUser(id Snowflake) (user *User, err error) {
	c.users.RLock()
	defer c.users.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.users.Get(id); !exists {
		err = NewErrorCacheItemNotFound(id)
		return
	}

	if c.conf.Immutable {
		user = result.Object().(*User).DeepCopy().(*User)
	} else {
		user = result.Object().(*User)
	}

	return
}
