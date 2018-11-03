package cache

import (
	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/disgord/cache/lfu"
	"testing"
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
		b.Logf("efficiency: %f", cache.Efficiency())
	})
}
