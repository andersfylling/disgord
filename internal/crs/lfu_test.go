// +build !integration

package crs

import "testing"

type randomStruct struct {
	ID Snowflake
}

func TestCacheList(t *testing.T) {
	t.Run("size limit", func(t *testing.T) {
		limit := uint(10)
		list := NewLFU(limit)
		if list.size != 0 {
			t.Error("size if not 0")
		}
		for i := 1; i < 256; i++ {
			usr := &randomStruct{}
			usr.ID = Snowflake(i)

			item := newLFUItem(usr)
			list.Set(usr.ID, item)
		}

		if list.size > limit {
			t.Errorf("list has a greater size than expected limit. Got %d, wants %d", list.size, limit)
		}
	})
	t.Run("replaces only LFU", func(t *testing.T) {
		ids := []Snowflake{1, 3, 5, 7, 9}
		list := NewLFU(uint(len(ids)))
		for i := 0; i < 256; i++ {
			usr := &randomStruct{}
			usr.ID = Snowflake(i)
			item := newLFUItem(usr)

			for _, id := range ids {
				if usr.ID == id {
					item.counter += 4
					break
				} else if usr.ID < id || usr.ID > ids[len(ids)-1] {
					break
				}
			}

			list.Set(usr.ID, item)
		}

		for i, item := range list.items {
			// except index 0. This is because th last content created must be placed in the cache, and then overwrite an
			// content with a count of 4. Since we always start at index 0 when we look for items with a smaller counter
			// and every entry has a count of 4, the first entry will be replace as there is no better alternative.
			if item.counter < 4 && i != 0 {
				t.Errorf("expected lfu counter for index %d to be higher. Got %d, wants above %d", i, item.counter, 4)
			}
		}
	})
}
