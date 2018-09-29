package lru

import (
	"sync"
	"time"

	"github.com/andersfylling/snowflake/v2"
)

type Snowflake = snowflake.Snowflake

func NewCacheItem(content interface{}) *CacheItem {
	return &CacheItem{
		item: content,
	}
}

type CacheItem struct {
	item interface{}

	// unix timestamp when item is considered outdated/dead.
	// update on creation/changes
	death int64
}

func (i *CacheItem) Object() interface{} {
	return i.item
}

func (i *CacheItem) Set(v interface{}) {
	i.item = v
}

func (i *CacheItem) update(lifetime time.Duration) {
	i.death = time.Now().Add(lifetime).UnixNano()
}

func (i *CacheItem) dead(now time.Time) bool {
	return i.death <= now.UnixNano()
}

// olderActivityThan checks if the last time `other` was used/updated is a more recent point in time then `i`
func (i *CacheItem) olderActivityThan(other *CacheItem) bool {
	return i.death < other.death
}

func NewCacheList(size uint, lifetime time.Duration, updateLifetime bool) *CacheList {
	return &CacheList{
		items:                 make(map[Snowflake]*CacheItem, size),
		limit:                 size,
		lifetime:              lifetime,
		updateLifetimeOnUsage: updateLifetime,
	}
}

type CacheList struct {
	sync.RWMutex
	items                 map[Snowflake]*CacheItem
	limit                 uint          // 0 == unlimited
	lifetime              time.Duration // 0 == unlimited
	updateLifetimeOnUsage bool
}

func (list *CacheList) size() uint {
	return uint(len(list.items))
}

func (list *CacheList) First() (item *CacheItem, key Snowflake) {
	for key, item = range list.items {
		return
	}

	return
}

// set adds a new item to the list or returns false if the item already exists
func (list *CacheList) Set(id Snowflake, new *CacheItem) {
	if new.death == 0 || list.updateLifetimeOnUsage {
		new.update(list.lifetime)
	}
	if item, exists := list.items[id]; exists { // check if it points to a diff item
		if item.item != new.item {
			*item = *new
		}
		return
	} else {
		list.items[id] = new
	}

	// if limit is reached, replace the content of the least recently used (lru)
	if list.limit != 0 && list.size() == list.limit {
		lru, lruKey := list.First()
		now := time.Now()
		for key, item := range list.items {
			if item.dead(now) {
				// TODO: can .dead(...) actually be slower?
				lruKey = key
				break
			}
			if item.olderActivityThan(lru) {
				lru = item
				lruKey = key
			}
		}

		delete(list.items, lruKey)
	}
}

// get an item from the list.
func (list *CacheList) Get(id Snowflake) (item *CacheItem, exists bool) {
	item, exists = list.items[id]
	if list.updateLifetimeOnUsage {
		item.update(list.lifetime)
	}
	return
}
