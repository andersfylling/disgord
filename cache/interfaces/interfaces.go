package interfaces

import (
	"github.com/andersfylling/snowflake/v3"
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

	// Size returns the number of actual elements in the list
	Size() uint

	// Cap returns the capacity of the cache. Note that 0 means there is no limit
	Cap() uint

	// ClearSoft calls ClearTables and creates a new list slice/map.
	ClearSoft()

	// ClearHard forces instant GC on every element in the list + tables
	ClearHard()

	// ClearTableNils removes nil objects from the tracking table
	ClearTableNils()

	// ClearTables clears tracking tables, such that the list can be overwritten with new content. But does not clear the list.
	ClearTables()
}
