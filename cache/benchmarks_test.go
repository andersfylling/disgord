package cache

import (
	"testing"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/cache/lfu"
	"github.com/andersfylling/disgord/cache/lru"
)

type Snowflake = interfaces.Snowflake

type randomStruct struct {
	ID Snowflake
}

func BenchmarkAllCacheReplacementStrategies(b *testing.B) {

	// lfu
	lfuCache := lfu.NewCacheList(0)
	benchmarkCacheSet(b, "lfu-unlimited", lfuCache)
	benchmarkCacheSet(b, "lfu-limited", lfu.NewCacheList(10000))
	benchmarkCacheUpdate(b, "lfu", lfuCache)
	benchmarkCacheGet(b, "lfu", lfuCache) // should output 1

	// lru
	lruCache := lru.NewCacheList(0)
	benchmarkCacheSet(b, "lru-unlimited", lruCache)
	benchmarkCacheSet(b, "lru-limited", lru.NewCacheList(10000))
	benchmarkCacheUpdate(b, "lru", lruCache)
	benchmarkCacheGet(b, "lru", lruCache) // should output 1
}

func benchmarkCacheSet(b *testing.B, name string, cache interfaces.CacheAlger) {
	b.Run("set-"+name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := Snowflake(uint64(i))
			cache.Set(id, cache.CreateCacheableItem(&randomStruct{ID: id}))
		}
	})
}

func benchmarkCacheUpdate(b *testing.B, name string, cache interfaces.CacheAlger) {
	b.Run("update-"+name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := Snowflake(uint64(i))
			cache.Set(id, cache.CreateCacheableItem(&randomStruct{ID: id}))
		}
	})
}

func benchmarkCacheGet(b *testing.B, name string, cache interfaces.CacheAlger) {
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
