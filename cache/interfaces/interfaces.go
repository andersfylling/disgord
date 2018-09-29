package interfaces

import (
	"github.com/andersfylling/snowflake/v2"
)

type Snowflake = snowflake.Snowflake

type CacheableItem interface {
	Object() interface{}
	Set(v interface{})
}

type CacheAlger interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()

	Get(id Snowflake) (item CacheableItem, exists bool)
	Set(id Snowflake, new CacheableItem)
	CreateCacheableItem(content interface{}) CacheableItem
	Efficiency() float64
}
