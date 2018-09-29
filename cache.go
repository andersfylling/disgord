package disgord

import (
	"errors"
	"time"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/cache/lfu"
	"github.com/andersfylling/disgord/cache/lru"
	"github.com/andersfylling/disgord/cache/tlru"
)

// cache keys
const (
	UserCache = iota
	ChannelCache
	GuildCache
	VoiceStateCache

	CacheAlg_LRU  = "lru"
	CacheAlg_LFU  = "lfu"
	CacheAlg_TLRU = "tlru"
)

type Cacher interface {
	Update(key int, v interface{}) (err error)
	Get(key int, id Snowflake, args ...interface{}) (v interface{}, err error)
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

func NewErrorUsingDeactivatedCache(cacheName string) *ErrorUsingDeactivatedCache {
	return &ErrorUsingDeactivatedCache{
		info: "cannot use deactivated cache: " + cacheName,
	}
}

type ErrorUsingDeactivatedCache struct {
	info string
}

func (e *ErrorUsingDeactivatedCache) Error() string {
	return e.info
}

func constructSpecificCacher(alg string, limit uint, lifetime time.Duration) (cacher interfaces.CacheAlger, err error) {
	switch alg {
	case CacheAlg_TLRU:
		cacher = tlru.NewCacheList(limit, lifetime)
	case CacheAlg_LRU:
		cacher = lru.NewCacheList(limit)
	case CacheAlg_LFU:
		cacher = lfu.NewCacheList(limit)
	default:
		err = errors.New("unsupported caching algorithm")
	}

	return
}

func createUserCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if !conf.UserCaching {
		return nil, nil
	}

	const userWeight = 1 // MiB. TODO: what is the actual max size?
	limit := conf.UserCacheLimitMiB / userWeight

	cacher, err = constructSpecificCacher(conf.UserCacheAlgorithm, limit, conf.UserCacheLifetime)
	return
}

func createVoiceStateCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if !conf.VoiceStateCaching {
		return nil, nil
	}

	cacher, err = constructSpecificCacher(conf.VoiceStateCacheAlgorithm, 0, conf.VoiceStateCacheLifetime)
	return
}

func NewCache(conf *CacheConfig) (*Cache, error) {

	userCacher, err := createUserCacher(conf)
	if err != nil {
		return nil, err
	}

	voiceStateCacher, err := createVoiceStateCacher(conf)
	if err != nil {
		return nil, err
	}

	return &Cache{
		conf:        conf,
		users:       userCacher,
		voiceStates: voiceStateCacher,
	}, nil
}

type CacheConfig struct {
	Immutable bool

	UserCaching        bool
	UserCacheLimitMiB  uint
	UserCacheLifetime  time.Duration
	UserCacheAlgorithm string

	VoiceStateCaching bool
	//VoiceStateCacheLimitMiB              uint
	VoiceStateCacheLifetime  time.Duration
	VoiceStateCacheAlgorithm string
}

type Cache struct {
	conf        *CacheConfig
	users       interfaces.CacheAlger
	voiceStates interfaces.CacheAlger
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
	case VoiceStateCache:
		if state, isVS := v.(*VoiceState); isVS {
			c.SetVoiceState(state)
		} else {
			err = errors.New("can only save *VoiceState structures to voice state cache")
		}
	default:
		err = errors.New("caching for given type is not yet implemented")
	}

	return
}

func (c *Cache) Get(key int, id Snowflake, args ...interface{}) (v interface{}, err error) {
	switch key {
	case UserCache:
		v, err = c.GetUser(id)
	case VoiceStateCache:
		if len(args) > 0 {
			if params, ok := args[0].(*guildVoiceStateCacheParams); ok {
				v, err = c.GetVoiceState(id, params)
			} else {
				err = errors.New("voice state cache extraction requires an addition argument of type *guildVoiceStateCacheParams")
			}
		} else {
			err = errors.New("voice state cache extraction requires an addition argument for filtering")
		}
	default:
		err = errors.New("caching for given type is not yet implemented")
	}

	return
}
