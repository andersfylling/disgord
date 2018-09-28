package lru

import (
	"testing"
	"time"
)

type randomStruct struct {
	ID Snowflake
}

func TestCacheItem(t *testing.T) {
	type smth struct {
		val string
	}
	lifetime := time.Duration(1) * time.Hour

	i := smth{"test"}
	item := NewCacheItem(&i)

	old := item.death
	item.update(lifetime)
	if old == item.death {
		t.Error("update method does not change death timestamp")
	}

	if item.dead(time.Now()) {
		t.Error("item is considered dead an hour before it's time")
	}
	item.update(lifetime * -1)
	if !item.dead(time.Now()) {
		t.Error("item was expected to be dead an hour ago")
	}
}

func TestCacheList(t *testing.T) {
	t.Run("size limit", func(t *testing.T) {
		limit := uint(10)
		list := NewCacheList(limit, time.Duration(1)*time.Hour, false)
		if list.size() != 0 {
			t.Error("size if not 0")
		}
		for i := 1; i < 256; i++ {
			usr := &randomStruct{}
			usr.ID = Snowflake(i)

			item := NewCacheItem(usr)
			list.Set(usr.ID, item)
		}

		if list.size() > limit {
			t.Errorf("list has a greater size than expected limit. Got %d, wants %d", list.size(), limit)
		}
	})
	t.Run("replaces only LRU", func(t *testing.T) {
		limit := uint(10)
		list := NewCacheList(limit, time.Duration(1)*time.Hour, false)
		for i := 1; i < 256; i++ {
			usr := &randomStruct{}
			usr.ID = Snowflake(i)

			item := NewCacheItem(usr)
			list.Set(usr.ID, item)
		}

		for key := range list.items {
			if key < 245 {
				t.Errorf("expected key to be higher. Got %d, wants above %d", key, 245)
			}
		}
	})
}
