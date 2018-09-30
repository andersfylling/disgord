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
	Delete(id Snowflake)
	CreateCacheableItem(content interface{}) CacheableItem
	RefreshAfterDiscordUpdate(item CacheableItem)
	Efficiency() float64
}
