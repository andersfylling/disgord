package interfaces

import (
	"github.com/andersfylling/snowflake/v2"
)

// Snowflake ...
type Snowflake = snowflake.Snowflake

// CacheableItem an item that can be cached
type CacheableItem interface {
	Object() interface{}
	Set(v interface{})
}

// CacheAlger a cache replacement algorithm interface
type CacheAlger interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()

	Get(id Snowflake) (item CacheableItem, exists bool)
	Set(id Snowflake, new CacheableItem)
	Delete(id Snowflake)
	CreateCacheableItem(content interface{}) CacheableItem
	RefreshAfterDiscordUpdate(item CacheableItem)
	Efficiency() float64
}
