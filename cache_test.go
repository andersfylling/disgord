package disgord

import "testing"

func TestCache_ChannelCreate(t *testing.T) {
	t.Run("immutable", func(t *testing.T) {
		cache, _ := newCache(&CacheConfig{
			DisableGuildCaching:      true,
			DisableUserCaching:       true,
			DisableVoiceStateCaching: true,
		})

		c1 := NewChannel()
		c1.ID = Snowflake(1234123)
		cache.SetChannel(c1)

		c1.ID = Snowflake(4537345435)

		if r, err := cache.GetChannel(Snowflake(1234123)); r.ID != Snowflake(1234123) || err != nil {
			t.Error(err)
			t.Error("error with retrieving channel")
		}
	})

	t.Run("mutable", func(t *testing.T) {
		cache, _ := newCache(&CacheConfig{
			Mutable:                  true,
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

		if r, _ := cache.GetChannel(c1.ID); *r.Icon != "changed" {
			t.Error("channel was not affected by external changes")
		}
	})
}
