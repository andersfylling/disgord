// lfu (least frequently counter) will overwrite cached items that have been counter the least when the cache limit is reached.
package lfu

import (
	"sync"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/snowflake/v2"
)

type Snowflake = snowflake.Snowflake

// NewCacheItem ...
func NewCacheItem(content interface{}) *CacheItem {
	return &CacheItem{
		item: content,
	}
}

// CacheItem ...
type CacheItem struct {
	item    interface{}
	counter uint64
}

// Object ...
func (i *CacheItem) Object() interface{} {
	return i.item
}

// Set ...
func (i *CacheItem) Set(v interface{}) {
	i.item = v
}

func (i *CacheItem) increment() {
	i.counter++
}

// NewCacheList ...
func NewCacheList(size uint) *CacheList {
	return &CacheList{
		items: make(map[Snowflake]*CacheItem, size),
		limit: size,
	}
}

type CacheList struct {
	sync.RWMutex
	items map[Snowflake]*CacheItem
	limit uint // 0 == unlimited

	misses uint64 // opposite of cache hits
	hits   uint64
}

func (list *CacheList) size() uint {
	return uint(len(list.items))
}

// First ...
func (list *CacheList) First() (item *CacheItem, key Snowflake) {
	for key, item = range list.items {
		return
	}

	return
}

// Set set adds a new item to the list or returns false if the item already exists
func (list *CacheList) Set(id Snowflake, newItemI interfaces.CacheableItem) {
	newItem := newItemI.(*CacheItem)
	if item, exists := list.items[id]; exists { // check if it points to a diff item
		if item.item != newItem.item || item != newItem {
			*item = *newItem
		}
		return
	} else {
		list.items[id] = newItem
	}

	if list.limit == 0 || list.size() <= list.limit {
		return
	}
	// if limit is reached, replace the content of the least recently counter (lru)
	list.removeLFU(id)
}

func (list *CacheList) removeLFU(exception Snowflake) {
	lfu, lfuKey := list.First()
	for key, item := range list.items {
		if key != exception && item.counter < lfu.counter {
			// TODO: create an lru map, for later?
			lfu = item
			lfuKey = key
		}
	}

	delete(list.items, lfuKey)
}

// RefreshAfterDiscordUpdate ...
func (list *CacheList) RefreshAfterDiscordUpdate(itemI interfaces.CacheableItem) {
	item := itemI.(*CacheItem)
	item.increment()
}

// Get get an item from the list.
func (list *CacheList) Get(id Snowflake) (ret interfaces.CacheableItem, exists bool) {
	var item *CacheItem
	if item, exists = list.items[id]; exists {
		ret = item
		item.increment()
		list.hits++
	} else {
		list.misses++
	}
	return
}

// Delete ...
func (list *CacheList) Delete(id Snowflake) {
	if _, exists := list.items[id]; exists {
		delete(list.items, id)
	}
}

// CreateCacheableItem ...
func (list *CacheList) CreateCacheableItem(content interface{}) interfaces.CacheableItem {
	return NewCacheItem(content)
}

// Efficiency ...
func (list *CacheList) Efficiency() float64 {
	return float64(list.hits) / float64(list.misses+list.hits)
}

var _ interfaces.CacheAlger = (*CacheList)(nil)
var _ interfaces.CacheableItem = (*CacheItem)(nil)
