package disgord

import (
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/json"
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
		case <-time.After(1 * time.Second):
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

	t.Run("get", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			cache := NewBasicCache()
			cache.Channels.Store[id] = &Channel{ID: id}

			channel, err := cache.GetChannel(id)
			if err != nil {
				t.Error("cache has no channel")
			}
			if channel == nil {
				t.Error("channel is nil")
			}
		})
		t.Run("get unknown", func(t *testing.T) {
			cache := NewBasicCache()

			channel, err := cache.GetChannel(id)
			if err == nil {
				t.Error("should return error when channel is unknown")
			}
			if channel != nil {
				t.Error("channel should be nil")
			}
			if !errors.Is(err, CacheMissErr) {
				t.Error("expected error to be a cache miss err")
			}
		})
	})

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
			now.Add(-10 * time.Second)
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

func TestBasicCache_TypingStart(t *testing.T) {
	cache := NewBasicCache()

	userID := Snowflake(123)
	channelID := Snowflake(348765348)
	evtData := jsonbytes(`{"user_id":%d,"channel_id":"%s"}`, userID, channelID)

	t.Run("event", func(t *testing.T) {
		evt, err := cacheDispatcher(cache, EvtTypingStart, evtData)
		if err != nil {
			t.Fatal("failed to create event struct", err)
		}

		typingStart := evt.(*TypingStart)

		if typingStart.UserID != userID {
			t.Errorf("incorrect user id. Got %d, wants %d", typingStart.UserID, userID)
		}
		if typingStart.ChannelID != channelID {
			t.Errorf("incorrect channel id. Got %d, wants %d", typingStart.ChannelID, channelID)
		}
	})

	deadlockTest(t, cache, EvtTypingStart, evtData)
}

func TestBasicCache_Ready(t *testing.T) {
	cache := NewBasicCache()

	guildIDsToGuilds := func(ids []Snowflake) (container []*GuildUnavailable) {
		for _, id := range ids {
			container = append(container, &GuildUnavailable{ID: id, Unavailable: true})
		}
		return
	}

	guilds := guildIDsToGuilds([]Snowflake{3, 4, 6, 7, 3})
	guildsJson, err := json.Marshal(guilds)
	if err != nil {
		t.Fatal("unable to marshal unavail guilds")
	}

	evtData := jsonbytes(`{"v":8,"user":%s,"guilds":%s,"session_id":"gf7k4gfe78g"}`, `{"id":234}`, guildsJson)

	t.Run("event", func(t *testing.T) {
		evt, err := cacheDispatcher(cache, EvtReady, evtData)
		if err != nil {
			t.Fatal("failed to create event struct", err)
		}

		ready := evt.(*Ready)

		if ready.User.ID != 234 {
			t.Errorf("incorrect user id. Got %d, wants %d", ready.User.ID, 234)
		}

		if len(guilds) != len(ready.Guilds) {
			t.Error("incorrect number of guilds")
		}

		if len(cache.Guilds.Store) != len(guilds) {
			t.Errorf("cache has incorrect number of guilds pre-allocated. Got %d, wants %d", len(cache.Guilds.Store), len(guilds))
		}

		for _, sourceGuild := range guilds {
			if _, ok := cache.Guilds.Store[sourceGuild.ID]; !ok {
				t.Errorf("store is missing guild ID %d", sourceGuild.ID)
			}
		}

		if cache.CurrentUser.ID != 234 {
			t.Error("current user id was not updated")
		}
	})

	deadlockTest(t, cache, EvtReady, evtData)
}

func TestBasicCache_Message(t *testing.T) {
	cache := NewBasicCache()

	evtData := jsonbytes(`{"id":1,"content":"testing","guild_id":2,"channel_id":3}`)

	t.Run("create", func(t *testing.T) {
		evt, err := cacheDispatcher(cache, EvtMessageCreate, evtData)
		if err != nil {
			t.Fatal("failed to create event struct", err)
		}

		msg := evt.(*MessageCreate).Message
		if msg.ID != 1 {
			t.Error("incorrect message id")
		}

		// should not create a DM channel
		if len(cache.Channels.Store) > 0 {
			t.Error("channel was created")
		}
	})

	deadlockTest(t, cache, EvtMessageCreate, evtData)

	t.Run("create DM", func(t *testing.T) {
		// if guild id is not set, it's a DM message
		// TODO: group DM
		evt, err := cacheDispatcher(cache, EvtMessageCreate, jsonbytes(`{"id":1,"content":"testing","channel_id":3}`))
		if err != nil {
			t.Fatal("failed to create event struct", err)
		}

		msg := evt.(*MessageCreate).Message
		if msg.ID != 1 {
			t.Error("incorrect message id")
		}

		if len(cache.Channels.Store) == 0 {
			t.Error("missing DM channel")
		}

		channel, err := cache.GetChannel(3)
		if errors.Is(err, CacheMissErr) {
			t.Fatal("DM channel was not created for message")
		}

		if channel.Type != ChannelTypeDM {
			t.Errorf("channel was created with incorrect type. Got %d, wants %d", channel.Type, ChannelTypeDM)
		}
	})
}

func TestBasicCache_UserUpdate(t *testing.T) {
	cache := NewBasicCache()
	cache.Users.Store[1] = &User{ID: 1, Username: "anders", Bot: true}

	evt, err := cacheDispatcher(cache, EvtUserUpdate, jsonbytes(`{"id":1,"username":"test"}`))
	if err != nil {
		t.Fatal("failed to create event struct", err)
	}

	usr := evt.(*UserUpdate)
	if usr.ID != 1 {
		t.Fatal("incorrect user id")
	}

	if usr.Username != "test" {
		t.Error("username was not updated")
	}
	if !usr.Bot {
		t.Error("bot value was overwritten")
	}

	if cache.CurrentUser.Username != "test" {
		t.Error("current users username was not updated")
	}
}

func TestBasicCache_Guilds(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			cache := NewBasicCache()
			cache.Guilds.Store[1] = &Guild{ID: 1}

			channel, err := cache.GetGuild(1)
			if err != nil {
				t.Error("cache has no such guild")
			}
			if channel == nil {
				t.Error("guild is nil")
			}
		})
		t.Run("get unknown", func(t *testing.T) {
			cache := NewBasicCache()

			channel, err := cache.GetChannel(1)
			if err == nil {
				t.Error("should return error when guild is unknown")
			}
			if channel != nil {
				t.Error("guild should be nil")
			}
			if !errors.Is(err, CacheMissErr) {
				t.Error("expected error to be a cache miss err")
			}
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("kicked", func(t *testing.T) {
			cache := NewBasicCache()
			cache.Guilds.Store[1] = &Guild{ID: 1}

			evt, err := cacheDispatcher(cache, EvtGuildDelete, jsonbytes(`{"id":%d}`, 1))
			if err != nil {
				t.Fatal("failed to create event struct", err)
			}

			guildEvt := evt.(*GuildDelete).UnavailableGuild
			if guildEvt.ID != 1 {
				t.Error("incorrect guild id")
			}

			guild, err := cache.GetGuild(1)
			if !errors.Is(err, CacheMissErr) {
				t.Error("expected cache miss err")
			}
			if guild != nil {
				t.Error("guild should be nil")
			}
		})
		t.Run("deleted", func(t *testing.T) {
			cache := NewBasicCache()
			cache.Guilds.Store[1] = &Guild{ID: 1}

			evt, err := cacheDispatcher(cache, EvtGuildDelete, jsonbytes(`{"id":%d,"unavailable":true}`, 1))
			if err != nil {
				t.Fatal("failed to create event struct", err)
			}

			guildEvt := evt.(*GuildDelete).UnavailableGuild
			if guildEvt.ID != 1 {
				t.Error("incorrect guild id")
			}
			if !guildEvt.Unavailable {
				t.Error("should have been unavail")
			}

			guild, err := cache.GetGuild(1)
			if !errors.Is(err, CacheMissErr) {
				t.Error("expected cache miss err")
			}
			if guild != nil {
				t.Error("guild should be nil")
			}
		})
	})
}
