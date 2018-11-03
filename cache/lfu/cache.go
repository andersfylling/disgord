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
	id      Snowflake
	item    interface{}
	counter uint64
}

// Object ...
func (i *CacheItem) Object() interface{} {
	i.increment()
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
	list := &CacheList{
		items:    make([]CacheItem, size),
		table:    make(map[Snowflake]int, size),
		nilTable: make([]int, size),
		limit:    size,
	}

	for i := 0; i < int(size); i++ {
		list.nilTable[i] = i
	}

	return list
}

type CacheList struct {
	sync.RWMutex
	items    []CacheItem
	table    map[Snowflake]int
	nilTable []int
	limit    uint // 0 == unlimited
	size     uint

	misses uint64 // opposite of cache hits
	hits   uint64
}

// Set set adds a new item to the list or returns false if the item already exists
func (list *CacheList) Set(id Snowflake, newItemI interfaces.CacheableItem) {
	newItem := newItemI.(*CacheItem)
	if key, exists := list.table[id]; exists && key != -1 {
		list.items[key].item = newItem.item
		return
	}

	if list.limit > 0 && list.size >= list.limit {
		// if limit is reached, replace the content of the least recently counter (lru)
		list.removeLFU()
	}

	if len(list.nilTable) > 0 {
		key := list.nilTable[len(list.nilTable)-1]
		list.nilTable = list.nilTable[:len(list.nilTable)-1]
		list.items[key] = *newItem
		list.table[id] = key
		list.size++
	} else {
		key := len(list.items)
		list.items = append(list.items, *newItem)
		list.table[id] = key
		list.size++
	}
}

func (list *CacheList) removeLFU() {
	lfuKey := 0
	lfu := list.items[lfuKey]
	var i int
	for i = range list.items {
		if list.items[i].counter < lfu.counter {
			// TODO: create a link to lowest counter for later?
			lfu = list.items[i]
			lfuKey = i
		}
	}

	list.deleteUnsafe(lfuKey, lfu.id)
}

// RefreshAfterDiscordUpdate ...
func (list *CacheList) RefreshAfterDiscordUpdate(itemI interfaces.CacheableItem) {
	item := itemI.(*CacheItem)
	item.increment()
}

// Get get an item from the list.
func (list *CacheList) Get(id Snowflake) (ret interfaces.CacheableItem, exists bool) {
	if key, exists := list.table[id]; exists && key != -1 {
		ret = &list.items[key]
		list.items[key].increment()
		list.hits++
	} else {
		list.misses++
	}
	return
}

func (list *CacheList) deleteUnsafe(key int, id Snowflake) {
	list.table[id] = -1
	list.nilTable = append(list.nilTable, key)
	list.size--
	//list.items[key] = list.items[len(list.items)-1]
	//list.items = list.items[:len(list.items)-1]
}

// Delete ...
func (list *CacheList) Delete(id Snowflake) {
	if key, exists := list.table[id]; exists && key != -1 {
		list.deleteUnsafe(key, id)
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
