package disgord

import (
	"errors"
	"time"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/cache/lfu"
	"github.com/andersfylling/disgord/cache/lru"
)

// cache keys to redirect to the related cache system
const (
	UserCache = iota
	ChannelCache
	GuildCache
	VoiceStateCache
)

// the different cache replacement algorithms
const (
	CacheAlgLRU  = "lru"
	CacheAlgLFU  = "lfu"
	CacheAlgTLRU = "tlru"
)

// Cacher gives basic cache interaction options, and won't require changes when adding more cache systems
type Cacher interface {
	Update(key int, v interface{}) (err error)
	Get(key int, id Snowflake, args ...interface{}) (v interface{}, err error)
}

func newErrorCacheItemNotFound(id Snowflake) *ErrorCacheItemNotFound {
	return &ErrorCacheItemNotFound{
		info: "item with id{" + id.String() + "} was not found in cache",
	}
}

// ErrorCacheItemNotFound requested item was not found in cache
type ErrorCacheItemNotFound struct {
	info string
}

// Error ...
func (e *ErrorCacheItemNotFound) Error() string {
	return e.info
}

func newErrorUsingDeactivatedCache(cacheName string) *ErrorUsingDeactivatedCache {
	return &ErrorUsingDeactivatedCache{
		info: "cannot use deactivated cache: " + cacheName,
	}
}

// ErrorUsingDeactivatedCache the cache system you are trying to access has been disabled in the CacheConfig
type ErrorUsingDeactivatedCache struct {
	info string
}

func (e *ErrorUsingDeactivatedCache) Error() string {
	return e.info
}

func constructSpecificCacher(alg string, limit uint, lifetime time.Duration) (cacher interfaces.CacheAlger, err error) {
	switch alg {
	case CacheAlgTLRU:
		//cacher = tlru.NewCacheList(limit, lifetime)
		err = errors.New("TLRU is missing schedulerer for clearing dead/timed out objects and is therefore deactivated")
	case CacheAlgLRU:
		cacher = lru.NewCacheList(limit)
	case CacheAlgLFU:
		cacher = lfu.NewCacheList(limit)
	default:
		err = errors.New("unsupported caching algorithm")
	}

	return
}

func newCache(conf *CacheConfig) (*Cache, error) {
	userCacher, err := createUserCacher(conf)
	if err != nil {
		return nil, err
	}

	voiceStateCacher, err := createVoiceStateCacher(conf)
	if err != nil {
		return nil, err
	}

	channelCacher, err := createChannelCacher(conf)
	if err != nil {
		return nil, err
	}

	return &Cache{
		immutable:   conf.Immutable,
		conf:        conf,
		users:       userCacher,
		voiceStates: voiceStateCacher,
		channels:    channelCacher,
	}, nil
}

// CacheConfig allows for tweaking the cache system on a personal need
type CacheConfig struct {
	Immutable bool // Must be immutable to support concurrent access and long-running tasks(!)

	DisableUserCaching bool
	UserCacheLimitMiB  uint
	UserCacheLifetime  time.Duration
	UserCacheAlgorithm string

	DisableVoiceStateCaching bool
	//VoiceStateCacheLimitMiB              uint
	VoiceStateCacheLifetime  time.Duration
	VoiceStateCacheAlgorithm string

	DisableChannelCaching bool
	ChannelCacheLimitMiB  uint
	ChannelCacheLifetime  time.Duration
	ChannelCacheAlgorithm string

	DisableGuildCaching bool
	GuildCacheLimitMiB  uint
	GuildCacheLifetime  time.Duration
	GuildCacheAlgorithm string
}

// Cache is the actual cache. It holds the different systems which can be tweaked using the CacheConfig.
type Cache struct {
	conf        *CacheConfig
	immutable   bool
	users       interfaces.CacheAlger
	voiceStates interfaces.CacheAlger
	channels    interfaces.CacheAlger
	guilds      interfaces.CacheAlger
}

// Updates does the same as Update. But allows for a slice of entries instead.
func (c *Cache) Updates(key int, vs []interface{}) (err error) {
	for _, v := range vs {
		err = c.Update(key, v)
		if err != nil {
			return
		}
	}

	return
}

// Update updates a item in the cache given the key identifier and the new content.
// It also checks if the given structs implements the required interfaces (See below).
func (c *Cache) Update(key int, v interface{}) (err error) {
	if v == nil {
		err = errors.New("object was nil")
		return
	}

	_, implementsDeepCopier := v.(DeepCopier)
	_, implementsCacheCopier := v.(cacheCopier)
	if !implementsCacheCopier && !implementsDeepCopier && c.immutable {
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
	case ChannelCache:
		if channel, isChannel := v.(*Channel); isChannel {
			c.SetChannel(channel)
		} else {
			err = errors.New("can only save *Channel structures to channel cache")
		}
	default:
		err = errors.New("caching for given type is not yet implemented")
	}

	return
}

// Get retrieve a item in the cache, or get an error when not found or if the cache system is disabled
// in your CacheConfig configuration.
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
	case ChannelCache:
		v, err = c.GetChannel(id)
	default:
		err = errors.New("caching for given type is not yet implemented")
	}

	return
}
