package disgord

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func jsonbytes(format string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(format, args...))
}

func deadlockTest(t *testing.T, c Cache, evt string, data []byte) {
	// all locks should have been released
	t.Run("deadlock", func(t *testing.T) {
		done := make(chan struct{})
		go func() {
			if _, err := cacheDispatcher(c, evt, data); err != nil {
				t.Error("failed to create channel from event data", err)
			}
			close(done)
		}()

		select {
		case <-time.After(1*time.Second):
			t.Fatal("deadlock detected")
		case <-done:
		}
	})
}

func TestBasicCache_Channels(t *testing.T) {
	cache := NewBasicCache()

	id := Snowflake(1)
	topic := "test"
	name := "anders"

	t.Run("create", func(t *testing.T) {
		evt, err := cacheDispatcher(cache, EvtChannelCreate, jsonbytes(`{"id":%d,"topic":"%s"}`, id, topic))
		if err != nil {
			t.Fatal("failed to create channel from event data", err)
		}

		holder, ok := evt.(*ChannelCreate)
		if !ok {
			t.Fatal("unable to cast event to ChannelCreate type")
		}

		if holder == nil {
			t.Fatal("holder is nil")
		}

		channel := holder.Channel
		if channel == nil {
			t.Fatal("channel is nil")
		}

		if channel.ID != id {
			t.Errorf("channel id should be %d, got %d", id, channel.ID)
		}
		if channel.Topic != topic {
			t.Errorf("channel topic should be %s, got %s", topic, channel.Topic)
		}
	})

	deadlockTest(t, cache, EvtChannelCreate, jsonbytes(`{"id":%d,"topic":"%s"}`, id*2, "abdsa"))

	t.Run("duplicate-create", func(t *testing.T) {
		// should just use the latest event
		// perhaps we were kicked from the guild and re-added. We might not have
		// properly deleted the channel, so the latest event should be utilized.

		currentChannel, err := cache.GetChannel(id)
		if err != nil {
			t.Fatal("cache has no channel")
		}

		// don't overwrite just given fields, completely delete old reference
		evt, err := cacheDispatcher(cache, EvtChannelCreate, jsonbytes(`{"id":%d,"name":"%s"}`, id, name))
		if err != nil {
			t.Fatal("failed to update channel from event data", err)
		}

		holder, ok := evt.(*ChannelCreate)
		if !ok {
			t.Fatal("unable to cast event to ChannelCreate type")
		}
		if holder == nil {
			t.Fatal("holder is nil")
		}

		channel := holder.Channel
		if channel == nil {
			t.Fatal("channel is nil")
		}

		if channel.Name != name {
			t.Errorf("channel name is %s, expects %s", channel.Name, name)
		}

		// make sure old data does not exist
		if channel.Topic != "" {
			t.Errorf("topic should be empty, but got %s", channel.Topic)
		}

		// make sure the old reference was not updated
		if currentChannel.Name == channel.Name {
			t.Errorf("old reference was updated, cache is not read only!")
		}
	})

	t.Run("update-existing", func(t *testing.T) {
		// update a channel already in cache
		oldChannel, err := cache.GetChannel(id)
		if err != nil {
			t.Fatal("cache has no channel")
		}

		updatedName := "test v2"

		// should only update given fields
		evt, err := cacheDispatcher(cache, EvtChannelUpdate, jsonbytes(`{"id":%d,"name":"%s"}`, id, updatedName))
		if err != nil {
			t.Fatal("failed to create channel from event data", err)
		}

		holder := evt.(*ChannelUpdate)
		channel := holder.Channel

		if channel.Topic != oldChannel.Topic {
			t.Error("topic was overwritten")
		}
		if channel.Name != updatedName {
			t.Errorf("name is %s, expected %s", channel.Name, updatedName)
		}
		if oldChannel.Name == updatedName {
			t.Error("cache is not read only")
		}
	})

	deadlockTest(t, cache, EvtChannelUpdate, jsonbytes(`{"id":%d,"topic":"%s"}`, id*2, "dsffddsfdf"))

	t.Run("update-unknown-channel", func(t *testing.T) {
		// if the channel does not exist, we should just create it
		unknownID := id * 23
		oldChannel, err := cache.GetChannel(unknownID)
		if !errors.Is(err, CacheMissErr) {
			t.Fatal("should have been a cache miss error")
		}
		if oldChannel != nil {
			t.Fatal("returned object should be nil")
		}

		evt, err := cacheDispatcher(cache, EvtChannelUpdate, jsonbytes(`{"id":%d,"topic":%s,name":"%s"}`, unknownID, topic, name))
		if err != nil {
			t.Fatal("failed to create channel from event data", err)
		}

		holder := evt.(*ChannelUpdate)
		channel := holder.Channel

		if channel.ID != id {
			t.Errorf("channel id should be %d, got %d", unknownID, channel.ID)
		}
		if channel.Topic != topic {
			t.Errorf("channel topic should be %s, got %s", topic, channel.Topic)
		}
		if channel.Name != name {
			t.Errorf("channel name should be %s, got %s", name, channel.Name)
		}
	})

	t.Run("pin update", func(t *testing.T) {
		channel, err := cache.GetChannel(id)
		if err != nil {
			t.Fatal("cache has no channel")
		}

		now := Time{
			Time: time.Now(),
		}

		t.Run("new", func(t *testing.T) {
			data, err := now.MarshalJSON()
			if err != nil {
				t.Fatal("unable to marshal pin timestamp")
			}

			evt, err := cacheDispatcher(cache, EvtChannelPinsUpdate, jsonbytes(`{"channel_id":%d,"last_pin_timestamp":%s}`, id, data))
			if err != nil {
				t.Fatal("failed to create event struct", err)
			}

			holder := evt.(*ChannelPinsUpdate)
			if holder.LastPinTimestamp.Second() != now.Second() {
				t.Errorf("incorrect time. Got %d, wants %d", holder.LastPinTimestamp.Second(), now.Second())
			}

			if !channel.LastPinTimestamp.IsZero() {
				t.Error("cache is not read-only")
			}

			channelNow, _ := cache.GetChannel(id)
			if channelNow.LastPinTimestamp.IsZero() {
				t.Error("last ping timestamp was not updated")
			}
		})

		t.Run("outdated", func(t *testing.T) {
			now.Add(-10*time.Second)
			data, err := now.MarshalJSON()
			if err != nil {
				t.Fatal("unable to marshal pin timestamp")
			}

			evt, err := cacheDispatcher(cache, EvtChannelPinsUpdate, jsonbytes(`{"channel_id":%d,"last_pin_timestamp":%s}`, id, data))
			if err != nil {
				t.Fatal("failed to create event struct", err)
			}

			holder := evt.(*ChannelPinsUpdate)
			if holder.LastPinTimestamp.Second() == now.Second() {
				t.Error("timestamp was updated")
			}

			if !channel.LastPinTimestamp.IsZero() {
				t.Error("cache is not read-only")
			}

			channelNow, _ := cache.GetChannel(id)
			if channelNow.LastPinTimestamp.IsZero() {
				t.Error("last ping timestamp was not updated")
			}
		})

	})

	t.Run("delete", func(t *testing.T) {
		channel, err := cache.GetChannel(id)
		if err != nil {
			t.Fatal("cache has no channel")
		}
		if channel == nil {
			t.Fatal("returned channel should not be nil")
		}

		evt, err := cacheDispatcher(cache, EvtChannelDelete, jsonbytes(`{"id":%d}`, id))
		if err != nil {
			t.Fatal("failed to create event struct", err)
		}

		holder := evt.(*ChannelDelete)
		if holder.Channel.ID != id {
			t.Errorf("expected channel id to be %d, got %d", id, holder.Channel.ID)
		}

		channel, err = cache.GetChannel(id)
		if !errors.Is(err, CacheMissErr) {
			t.Fatal("should have been a cache miss error")
		}
		if channel != nil {
			t.Fatal("returned object should be nil")
		}
	})
}