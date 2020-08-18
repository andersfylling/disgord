// +build !integration

package crs

import (
	"testing"
)

func BenchmarkCRS(b *testing.B) {
	// lfu
	lfuCache := NewLFU(0)
	benchmarkCacheSet(b, "lfu-unlimited", lfuCache)
	benchmarkCacheSet(b, "lfu-limited", NewLFU(10000))
	benchmarkCacheUpdate(b, "lfu", lfuCache)
	benchmarkCacheGet(b, "lfu", lfuCache) // should output 1
}

func benchmarkCacheSet(b *testing.B, name string, cache *LFU) {
	b.Run("set-"+name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := Snowflake(uint64(i))
			cache.Set(id, cache.CreateCacheableItem(&randomStruct{ID: id}))
		}
	})
}

func benchmarkCacheUpdate(b *testing.B, name string, cache *LFU) {
	b.Run("update-"+name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := Snowflake(uint64(i))
			cache.Set(id, cache.CreateCacheableItem(&randomStruct{ID: id}))
		}
	})
}

func benchmarkCacheGet(b *testing.B, name string, cache *LFU) {
	once := false
	b.Run("get-"+name, func(b *testing.B) {
		for i := 0; i < b.N && i < int(cache.Size()); i++ {
			id := Snowflake(uint64(i))
			cache.Get(id)
		}
		if !once {
			// if it prints out more than once, cause of a race condition, it's not a big deal. it's just annoying
			once = true
			b.Logf("efficiency: %f", cache.Efficiency())
		}
	})
}
