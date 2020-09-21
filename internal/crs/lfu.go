package crs

import (
	"runtime"
	"sync"
)

// New ...
func NewLFU(size uint) *LFU {
	list := &LFU{
		limit: size,
		table: make(map[Snowflake]int),
	}

	list.ClearSoft()

	return list
}

func SetLimit(v interface{}, limit uint) {
	switch t := v.(type) {
	case *LFU:
		t.limit = limit
		if t.table == nil {
			// this is hacky, should be initialised somewhere else!
			t.table = make(map[Snowflake]int)
		}
	default:
		panic("unsupported cache replacement system")
	}
}

type LFU struct {
	sync.RWMutex
	items    []LFUItem
	table    map[Snowflake]int
	nilTable []int
	limit    uint // 0 == unlimited
	size     uint

	misses uint64 // opposite of cache hits
	hits   uint64
}

func (list *LFU) Size() uint {
	return list.size
}

func (list *LFU) Cap() uint {
	return list.limit
}

func (list *LFU) ClearSoft() {
	for i := range list.items {
		// TODO: is this needed?
		list.items[i].Val = nil
	}
	list.items = make([]LFUItem, list.limit)
	list.ClearTables()
}

func (list *LFU) ClearHard() {
	list.ClearSoft()
	runtime.GC()
}

func (list *LFU) ClearTableNils() {
	size := 0
	for key := range list.table {
		if list.table[key] != -1 {
			size++
		}
	}

	newTable := make(map[Snowflake]int, list.limit)
	for key := range list.table {
		if list.table[key] != -1 {
			newTable[key] = list.table[key]
		}
	}
	list.table = newTable
}

func (list *LFU) ClearTables() {
	list.table = make(map[Snowflake]int)
	list.nilTable = make([]int, list.limit)

	for i := 0; i < int(list.limit); i++ {
		list.nilTable[i] = i
	}
}

// Set set adds a new content to the list or returns false if the content already exists
func (list *LFU) Set(id Snowflake, newItem *LFUItem) {
	newItem.ID = id
	if key, exists := list.table[id]; exists && key != -1 {
		list.items[key].Val = newItem.Val
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

func (list *LFU) removeLFU() {
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

	list.deleteUnsafe(lfuKey, lfu.ID)
}

// RefreshAfterDiscordUpdate ...
func (list *LFU) RefreshAfterDiscordUpdate(item *LFUItem) {
	item.increment()
}

// Get get an content from the list.
func (list *LFU) Get(id Snowflake) (ret *LFUItem, exists bool) {
	var key int
	if key, exists = list.table[id]; exists && key != -1 {
		ret = &list.items[key]
		ret.increment()
		list.hits++
	} else {
		exists = false // if key == -1, exists might still be true
		list.misses++
	}
	return
}

func (list *LFU) deleteUnsafe(key int, id Snowflake) {
	list.table[id] = -1
	list.items[key].Val = nil // prepare for GC
	list.nilTable = append(list.nilTable, key)
	list.size--
}

// Delete ...
func (list *LFU) Delete(id Snowflake) {
	if key, exists := list.table[id]; exists && key != -1 {
		list.deleteUnsafe(key, id)
	}
}

// CreateCacheableItem ...
func (list *LFU) CreateCacheableItem(content interface{}) *LFUItem {
	return newLFUItem(content)
}

// Efficiency ...
func (list *LFU) Efficiency() float64 {
	if list.hits == 0 {
		return 0.0
	}
	return float64(list.hits) / float64(list.misses+list.hits)
}
