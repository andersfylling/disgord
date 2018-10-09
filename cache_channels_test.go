package disgord

import "testing"

func TestCache_ChannelCreate(t *testing.T) {
	t.Run("immutable", func(t *testing.T) {
		cache, _ := NewCache(&CacheConfig{
			Immutable:                true,
			ChannelCacheAlgorithm:    CacheAlg_LRU,
			DisableGuildCaching:      true,
			DisableUserCaching:       true,
			DisableVoiceStateCaching: true,
		})

		c1 := NewChannel()
		c1.ID = Snowflake(1234123)
		cache.SetChannel(c1)

		c1.ID = Snowflake(4537345435)

		if c1_r, err := cache.GetChannel(Snowflake(1234123)); c1_r.ID != Snowflake(1234123) || err != nil {
			t.Error(err)
			t.Error("error with retrieving channel")
		}
	})

	t.Run("mutable", func(t *testing.T) {
		cache, _ := NewCache(&CacheConfig{
			ChannelCacheAlgorithm:    CacheAlg_LRU,
			DisableGuildCaching:      true,
			DisableUserCaching:       true,
			DisableVoiceStateCaching: true,
		})

		c1 := NewChannel()
		c1.ID = Snowflake(1234123)

		test := "test"
		c1.Icon = &test
		cache.SetChannel(c1)

		changed := "changed"
		c1.Icon = &changed

		if c1_r, _ := cache.GetChannel(c1.ID); *c1_r.Icon != "changed" {
			t.Error("channel was not affected by external changes")
		}
	})
}
