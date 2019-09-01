// lfu (least frequently counter) will overwrite cached items that have been counter the least when the cache limit is reached.
package lfu

import (
	"runtime"
	"sync"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/depalias"
)

type Snowflake = depalias.Snowflake

// NewCacheItem ...
func NewCacheItem(content interface{}) *CacheItem {
	return &CacheItem{
		content: content,
	}
}

// CacheItem ...
type CacheItem struct {
	id      Snowflake
	content interface{}
	counter uint64
}

// Object ...
func (i *CacheItem) Object() interface{} {
	return i.content
}

// Set ...
func (i *CacheItem) Set(v interface{}) {
	i.content = v
}

func (i *CacheItem) increment() {
	i.counter++
}

// NewCacheList ...
func NewCacheList(size uint) *CacheList {
	list := &CacheList{
		limit: size,
	}

	list.ClearSoft()

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

func (list *CacheList) Size() uint {
	return list.size
}

func (list *CacheList) Cap() uint {
	return list.limit
}

func (list *CacheList) ClearSoft() {
	for i := range list.items {
		// TODO: is this needed?
		list.items[i].content = nil
	}
	list.items = make([]CacheItem, list.limit)
	list.ClearTables()
}

func (list *CacheList) ClearHard() {
	list.ClearSoft()
	runtime.GC()
}

func (list *CacheList) ClearTableNils() {
	size := 0
	for key := range list.table {
		if list.table[key] != -1 {
			size++
		}
	}
	// TODO: create a tmp slice which holds only valid entries, and loop through those instead of re-looping?
	newTable := make(map[Snowflake]int, list.limit)
	for key := range list.table {
		if list.table[key] != -1 {
			newTable[key] = list.table[key]
		}
	}
	list.table = newTable
}

func (list *CacheList) ClearTables() {
	list.table = make(map[Snowflake]int)
	list.nilTable = make([]int, list.limit)

	for i := 0; i < int(list.limit); i++ {
		list.nilTable[i] = i
	}
}

// Set set adds a new content to the list or returns false if the content already exists
func (list *CacheList) Set(id Snowflake, newItemI interfaces.CacheableItem) {
	newItem := newItemI.(*CacheItem)
	newItem.id = id
	if key, exists := list.table[id]; exists && key != -1 {
		list.items[key].content = newItem.content
		return
	}

	if list.limit > 0 && list.size >= list.limit {
		// if limit is reached, replace the content of the least recently counter (lru)
		list.removeLFU()
	}

	var key int
	if len(list.nilTable) > 0 {
		key = list.nilTable[len(list.nilTable)-1]
		list.nilTable = list.nilTable[:len(list.nilTable)-1]
		list.items[key] = *newItem
	} else {
		key = len(list.items)
		list.items = append(list.items, *newItem)
	}
	list.table[id] = key
	list.size++
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

// Get get an content from the list.
func (list *CacheList) Get(id Snowflake) (ret interfaces.CacheableItem, exists bool) {
	key, exists := list.table[id]
	exists = exists && key != -1

	if exists {
		ret = &list.items[key]
		list.items[key].increment()
		list.hits++
	} else {
		list.misses++
	}

	return ret, exists
}

func (list *CacheList) deleteUnsafe(key int, id Snowflake) {
	list.table[id] = -1
	list.items[key].content = nil // prepare for GC
	list.nilTable = append(list.nilTable, key)
	list.size--
}

// Delete ...
func (list *CacheList) Delete(id Snowflake) {
	if key, exists := list.table[id]; exists && key != -1 {
		list.deleteUnsafe(key, id)
	}
}

func (list *CacheList) Foreach(cb func(interface{})) {
	for k := range list.items {
		cb(list.items[k])
	}
}

func (list *CacheList) ListIDs() (ids []Snowflake) {
	ids = make([]Snowflake, 0, len(list.items))
	for i := range list.items {
		ids = append(ids, list.items[i].id)
	}

	return ids
}

// CreateCacheableItem ...
func (list *CacheList) CreateCacheableItem(content interface{}) interfaces.CacheableItem {
	return NewCacheItem(content)
}

// Efficiency ...
func (list *CacheList) Efficiency() float64 {
	if list.hits == 0 {
		return 0.0
	}
	return float64(list.hits) / float64(list.misses+list.hits)
}

var _ interfaces.CacheAlger = (*CacheList)(nil)
var _ interfaces.CacheableItem = (*CacheItem)(nil)
