package disgord

import (
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/json"
	"strings"
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

func deadlockGetTest(t *testing.T, cb func()) {
	// all locks should have been released
	t.Run("deadlock", func(t *testing.T) {
		done := make(chan struct{})
		go func() {
			cb()
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
			if !errors.Is(err, ErrCacheMiss) {
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
		if !errors.Is(err, ErrCacheMiss) {
			t.Fatal("should have been a cache miss error")
		}
		if oldChannel != nil {
			t.Fatal("returned object should be nil")
		}

		evt, err := cacheDispatcher(cache, EvtChannelUpdate, jsonbytes(`{"id":%d,"topic":"%s","name":"%s"}`, unknownID, topic, name))
		if err != nil {
			t.Fatal("failed to create channel from event data", err)
		}

		holder := evt.(*ChannelUpdate)
		channel := holder.Channel

		if channel.ID != unknownID {
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
			now.Time = now.Add(-10 * time.Second)
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
				t.Error("incorrect timestamp")
			}

			if !channel.LastPinTimestamp.IsZero() {
				t.Error("cache is not read-only")
			}

			channelNow, _ := cache.GetChannel(id)
			if channelNow.LastPinTimestamp.Second() == now.Second() {
				t.Error("last ping timestamp was updated")
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
		if !errors.Is(err, ErrCacheMiss) {
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

		if len(cache.Guilds.Store) != 0 {
			t.Errorf("cache pre-allocated the guilds, but is it really needed?")
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
		if errors.Is(err, ErrCacheMiss) {
			t.Fatal("DM channel was not created for message")
		}

		if channel.Type != ChannelTypeDM {
			t.Errorf("channel was created with incorrect type. Got %d, wants %d", channel.Type, ChannelTypeDM)
		}
	})
}

func TestBasicCache_UserUpdate(t *testing.T) {
	cache := NewBasicCache()
	cache.CurrentUser = &User{ID: 1, Username: "anders", Bot: true}

	evt, err := cacheDispatcher(cache, EvtUserUpdate, jsonbytes(`{"id":1,"username":"test"}`))
	if err != nil {
		t.Fatal("failed to create event struct", err)
	}

	usr := evt.(*UserUpdate)
	if usr.ID != 1 {
		t.Fatalf("incorrect user id. Got %d", usr.ID)
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
			cache.Guilds.Store[1] = &guildCacheContainer{Guild: &Guild{ID: 1}}

			guild, err := cache.GetGuild(1)
			if err != nil {
				t.Error("cache has no such guild")
			}
			if guild == nil {
				t.Error("guild is nil")
			}
		})
		t.Run("get unknown", func(t *testing.T) {
			cache := NewBasicCache()

			guild, err := cache.GetChannel(1)
			if err == nil {
				t.Error("should return error when guild is unknown")
			}
			if guild != nil {
				t.Error("guild should be nil")
			}
			if !errors.Is(err, ErrCacheMiss) {
				t.Error("expected error to be a cache miss err")
			}
		})
	})

	t.Run("get complex", func(t *testing.T) {
		cache := NewBasicCache()
		cache.Guilds.Store[1] = &guildCacheContainer{
			Guild:      &Guild{ID: 1},
			ChannelIDs: []Snowflake{1, 4},
			Members: map[Snowflake]*Member{
				3:  {UserID: 3, Nick: "andy"},
				56: {UserID: 56},
				34: {UserID: 34},
			},
		}
		cache.Users.Store[3] = &User{ID: 3, Username: "anders"}
		cache.Users.Store[56] = &User{ID: 56, Username: "test"}
		cache.Users.Store[34] = &User{ID: 34, Username: "botlol"}
		cache.Channels.Store[1] = &Channel{ID: 1, Name: "channel#1"}
		cache.Channels.Store[4] = &Channel{ID: 4, Name: "fourth"}

		guild, err := cache.GetGuild(1)
		if err != nil {
			t.Error("cache has no such guild")
		}
		if guild == nil {
			t.Fatal("guild is nil")
		}

		if len(guild.Channels) != 2 {
			t.Error("incorrect number of channels")
		}
		if len(guild.Members) != 3 {
			t.Error("incorrect number of members")
		}
		if guild.Members[0].User == nil {
			t.Error("member is missing user object")
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("kicked", func(t *testing.T) {
			cache := NewBasicCache()
			cache.Guilds.Store[1] = &guildCacheContainer{Guild: &Guild{ID: 1}}

			evt, err := cacheDispatcher(cache, EvtGuildDelete, jsonbytes(`{"id":%d}`, 1))
			if err != nil {
				t.Fatal("failed to create event struct", err)
			}

			guildEvt := evt.(*GuildDelete).UnavailableGuild
			if guildEvt.ID != 1 {
				t.Error("incorrect guild id")
			}

			guild, err := cache.GetGuild(1)
			if !errors.Is(err, ErrCacheMiss) {
				t.Error("expected cache miss err")
			}
			if guild != nil {
				t.Error("guild should be nil")
			}
		})
		t.Run("deleted", func(t *testing.T) {
			cache := NewBasicCache()
			cache.Guilds.Store[1] = &guildCacheContainer{Guild: &Guild{ID: 1}}

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
			if !errors.Is(err, ErrCacheMiss) {
				t.Error("expected cache miss err")
			}
			if guild != nil {
				t.Error("guild should be nil")
			}
		})
	})

	cache := NewBasicCache()

	id := Snowflake(10)
	name := "test guild"

	memberID := Snowflake(1)
	nick := "andy"
	username := "anders"

	channelID := Snowflake(352)
	channelName := "sidgfs asd"

	t.Run("create-with-members-and-channels", func(t *testing.T) {
		memberData := jsonbytes(`{"nick":"%s","user":{"id":%d,"username":"%s"}}`, nick, memberID, username)
		channelData := jsonbytes(`{"id":%d,"name":"%s"}`, channelID, channelName)

		data := jsonbytes(`{"id":%d,"name":"%s","members":[%s],"channels":[%s]}`, id, name, memberData, channelData)

		evt, err := cacheDispatcher(cache, EvtGuildCreate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		holder, ok := evt.(*GuildCreate)
		if !ok {
			t.Fatal("unable to cast event to GuildCreate type")
		}

		if holder == nil {
			t.Fatal("holder is nil")
		}

		guild := holder.Guild
		if guild == nil {
			t.Fatal("guild is nil")
		}

		if guild.ID != id {
			t.Errorf("channel id should be %d, got %d", id, guild.ID)
		}
		if guild.Name != name {
			t.Errorf("channel topic should be %s, got %s", name, guild.Name)
		}

		container, ok := cache.Guilds.Store[id]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		if len(container.ChannelIDs) != 1 {
			t.Fatal("missing channel id")
		}
		if cid := container.ChannelIDs[0]; cid.IsZero() {
			t.Fatal("channel id was stored as 0")
		} else {
			channel, ok := cache.Channels.Store[cid]
			if !ok {
				t.Fatal("channel was not saved to cache")
			}

			if channel.ID != channelID {
				t.Error("incorrect channel id")
			}
			if channel.Name != channelName {
				t.Error("incorrect channel name")
			}
		}

		if member, ok := container.Members[memberID]; !ok {
			t.Error("member was not stored")
		} else {
			if member == nil {
				t.Fatal("member is nil")
			}
			if member.UserID != memberID {
				t.Error("member is missing user id")
			}
			if member.Nick != nick {
				t.Error("incorrect nickname")
			}
		}

		if user, ok := cache.Users.Store[memberID]; !ok {
			t.Error("user was not stored")
		} else {
			if user == nil {
				t.Fatal("user is nil")
			}
			if user.ID != memberID {
				t.Error("user is missing user id")
			}
			if user.Username != username {
				t.Error("incorrect username")
			}
		}
	})

	deadlockTest(t, cache, EvtGuildCreate, jsonbytes(`{"id":%d,"name":"%s"}`, id*2, "abdsa"))

	t.Run("create-without-members-and-channels", func(t *testing.T) {
		data := jsonbytes(`{"id":%d,"name":"%s","members":[],"channels":[]}`, id, name)

		cache := NewBasicCache()
		evt, err := cacheDispatcher(cache, EvtGuildCreate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		guild := evt.(*GuildCreate).Guild
		if guild.ID != id {
			t.Errorf("channel id should be %d, got %d", id, guild.ID)
		}
		if guild.Name != name {
			t.Errorf("channel topic should be %s, got %s", name, guild.Name)
		}

		container, ok := cache.Guilds.Store[id]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		if container.Members == nil {
			t.Error("members map was not initialized")
		}
	})

	t.Run("update-without-members-and-channels", func(t *testing.T) {
		data := jsonbytes(`{"id":%d,"name":"%s","members":[],"channels":[]}`, id, name)

		cache := NewBasicCache()
		evt, err := cacheDispatcher(cache, EvtGuildUpdate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		guild := evt.(*GuildUpdate).Guild
		if guild.ID != id {
			t.Errorf("channel id should be %d, got %d", id, guild.ID)
		}
		if guild.Name != name {
			t.Errorf("channel topic should be %s, got %s", name, guild.Name)
		}

		container, ok := cache.Guilds.Store[id]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		if container.Members == nil {
			t.Error("members map was not initialized")
		}
	})

	t.Run("update-with-members-and-channels", func(t *testing.T) {
		// these should not be stored as they will not be in a guild update
		memberData := jsonbytes(`{"nick":"test","user":{"id":345,"username":"fdhhghgj"}}`)
		channelData := jsonbytes(`{"id":345,"name":"fgdfhjlll"}`)

		newGuildName := name + " v2"
		data := jsonbytes(`{"id":%d,"name":"%s","members":[%s],"channels":[%s]}`, id, newGuildName, memberData, channelData)

		evt, err := cacheDispatcher(cache, EvtGuildUpdate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		holder, ok := evt.(*GuildUpdate)
		if !ok {
			t.Fatal("unable to cast event to GuildCreate type")
		}

		if holder == nil {
			t.Fatal("holder is nil")
		}

		guild := holder.Guild
		if guild == nil {
			t.Fatal("guild is nil")
		}

		if guild.ID != id {
			t.Errorf("channel id should be %d, got %d", id, guild.ID)
		}
		if guild.Name != newGuildName {
			t.Errorf("channel topic should be %s, got %s", newGuildName, guild.Name)
		}

		// check that nothing else changed!
		container, ok := cache.Guilds.Store[id]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		if len(container.ChannelIDs) != 1 {
			t.Fatal("missing channel id")
		}
		if cid := container.ChannelIDs[0]; cid.IsZero() {
			t.Fatal("channel id was stored as 0")
		} else {
			channel, ok := cache.Channels.Store[cid]
			if !ok {
				t.Fatal("channel was not saved to cache")
			}

			if channel.ID != channelID {
				t.Error("incorrect channel id")
			}
			if channel.Name != channelName {
				t.Error("incorrect channel name")
			}
		}

		if len(container.Members) != 1 {
			t.Error("incorrect number of members")
		}

		if member, ok := container.Members[memberID]; !ok {
			t.Error("member was not stored")
		} else {
			if member == nil {
				t.Fatal("member is nil")
			}
			if member.UserID != memberID {
				t.Error("member is missing user id")
			}
			if member.Nick != nick {
				t.Error("incorrect nickname")
			}
		}

		if user, ok := cache.Users.Store[memberID]; !ok {
			t.Error("user was not stored")
		} else {
			if user == nil {
				t.Fatal("user is nil")
			}
			if user.ID != memberID {
				t.Error("user is missing user id")
			}
			if user.Username != username {
				t.Error("incorrect username")
			}
		}
	})

	t.Run("update-on-unknown-guild", func(t *testing.T) {
		cache := NewBasicCache()

		memberData := jsonbytes(`{"nick":"%s","user":{"id":%d,"username":"%s"}}`, nick, memberID, username)
		channelData := jsonbytes(`{"id":%d,"name":"%s"}`, channelID, channelName)

		data := jsonbytes(`{"id":%d,"name":"%s","members":[%s],"channels":[%s]}`, id, name, memberData, channelData)

		evt, err := cacheDispatcher(cache, EvtGuildUpdate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		holder, ok := evt.(*GuildUpdate)
		if !ok {
			t.Fatal("unable to cast event to GuildCreate type")
		}

		if holder == nil {
			t.Fatal("holder is nil")
		}

		guild := holder.Guild
		if guild == nil {
			t.Fatal("guild is nil")
		}

		if guild.ID != id {
			t.Errorf("channel id should be %d, got %d", id, guild.ID)
		}
		if guild.Name != name {
			t.Errorf("channel topic should be %s, got %s", name, guild.Name)
		}

		// check that nothing else changed!
		container, ok := cache.Guilds.Store[id]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		if len(container.ChannelIDs) != 1 {
			t.Fatal("missing channel id")
		}
		if cid := container.ChannelIDs[0]; cid.IsZero() {
			t.Fatal("channel id was stored as 0")
		} else {
			channel, ok := cache.Channels.Store[cid]
			if !ok {
				t.Fatal("channel was not saved to cache")
			}

			if channel.ID != channelID {
				t.Error("incorrect channel id")
			}
			if channel.Name != channelName {
				t.Error("incorrect channel name")
			}
		}

		if len(container.Members) != 1 {
			t.Error("incorrect number of members")
		}

		if member, ok := container.Members[memberID]; !ok {
			t.Error("member was not stored")
		} else {
			if member == nil {
				t.Fatal("member is nil")
			}
			if member.UserID != memberID {
				t.Error("member is missing user id")
			}
			if member.Nick != nick {
				t.Error("incorrect nickname")
			}
		}

		if user, ok := cache.Users.Store[memberID]; !ok {
			t.Error("user was not stored")
		} else {
			if user == nil {
				t.Fatal("user is nil")
			}
			if user.ID != memberID {
				t.Error("user is missing user id")
			}
			if user.Username != username {
				t.Error("incorrect username")
			}
		}
	})

	t.Run("channel-create", func(t *testing.T) {
		const guildID = 1
		const channelID = 20

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{Guild: &Guild{ID: guildID}}

		data := []byte(fmt.Sprintf(`{"id":%d,"guild_id":%d}`, channelID, guildID))
		_, err := cacheDispatcher(cache, EvtChannelCreate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		channelIDs := cache.Guilds.Store[guildID].ChannelIDs
		if len(channelIDs) != 1 {
			t.Fatal("channel was not linked to guild")
		}
		if channelIDs[0] != Snowflake(channelID) {
			t.Error("wrong channel id")
		}
	})

	t.Run("channel-update", func(t *testing.T) {
		const guildID = 1
		const channelID = 20

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{Guild: &Guild{ID: guildID}}

		data := []byte(fmt.Sprintf(`{"id":%d,"guild_id":%d}`, channelID, guildID))
		_, err := cacheDispatcher(cache, EvtChannelUpdate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		channelIDs := cache.Guilds.Store[guildID].ChannelIDs
		if len(channelIDs) != 1 {
			t.Fatal("channel was not linked to guild")
		}
		if channelIDs[0] != Snowflake(channelID) {
			t.Error("wrong channel id")
		}
	})

	t.Run("channel-delete", func(t *testing.T) {
		const guildID = 1
		const channelID = 20

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild:      &Guild{ID: guildID},
			ChannelIDs: []Snowflake{channelID},
		}

		data := []byte(fmt.Sprintf(`{"id":%d,"guild_id":%d}`, channelID, guildID))
		_, err := cacheDispatcher(cache, EvtChannelDelete, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		channelIDs := cache.Guilds.Store[guildID].ChannelIDs
		if len(channelIDs) != 0 {
			t.Error("there should be no channels")
		}
	})
}

func TestBasicCache_GuildMembers(t *testing.T) {
	id := Snowflake(2523)

	t.Run("members chunk", func(t *testing.T) {
		cache := NewBasicCache()
		cache.Guilds.Store[id] = &guildCacheContainer{
			Guild:   &Guild{ID: id},
			Members: map[Snowflake]*Member{},
		}

		memberJson := func(id Snowflake, nick, username string) []byte {
			return jsonbytes(`{"user":{"id":%d,"username":"%s"},"nick":"%s"}`, id, username, nick)
		}

		type ref struct {
			id       Snowflake
			nick     string
			username string
		}

		memberRefs := []ref{
			{45, "test", "amazontester"},
			{245436, "", "okay"},
			{2345, "andy", "anders"},
		}

		memberJsons := []string{}
		for i := range memberRefs {
			m := memberRefs[i]
			j := memberJson(m.id, m.nick, m.username)
			memberJsons = append(memberJsons, string(j))
		}

		joinedMembers := strings.Join(memberJsons, ",")
		data := jsonbytes(`{"guild_id":%d,"members":[%s]}`, id, joinedMembers)

		evt, err := cacheDispatcher(cache, EvtGuildMembersChunk, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		holder, ok := evt.(*GuildMembersChunk)
		if !ok {
			t.Fatal("unable to cast event to GuildMembersChunk type")
		}

		if holder == nil {
			t.Fatal("holder is nil")
		}

		guildID := holder.GuildID
		if guildID != id {
			t.Fatal("guild id is incorrect")
		}
		if len(holder.Members) != len(memberRefs) {
			t.Fatal("incorrect number of members")
		}

		container, ok := cache.Guilds.Store[guildID]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		for i := range memberRefs {
			m := memberRefs[i]
			ref := fmt.Sprintf("member %s: ", m.id)
			if member, ok := container.Members[m.id]; !ok {
				t.Errorf(ref+"was not stored", m.id)
			} else {
				if member == nil {
					t.Fatalf(ref+"is nil", m.id)
				}
				if member.UserID != m.id {
					t.Errorf(ref + "is missing user id")
				}
				if member.Nick != m.nick {
					t.Errorf(ref + "incorrect nickname")
				}
			}
		}

		for i := range memberRefs {
			m := memberRefs[i]
			ref := fmt.Sprintf("user %s: ", m.id)
			if user, ok := cache.Users.Store[m.id]; !ok {
				t.Errorf(ref+"was not stored", m.id)
			} else {
				if user == nil {
					t.Fatalf(ref+"is nil", m.id)
				}
				if user.ID != m.id {
					t.Errorf(ref + "is missing user id")
				}
				if user.Username != m.username {
					t.Errorf(ref + "incorrect username")
				}
			}
		}

		deadlockTest(t, cache, EvtGuildMembersChunk, data)
	})

	t.Run("member add", func(t *testing.T) {
		cache := NewBasicCache()
		cache.Guilds.Store[id] = &guildCacheContainer{
			Guild:   &Guild{ID: id},
			Members: map[Snowflake]*Member{},
		}

		memberJson := func(id Snowflake, nick, username string) []byte {
			return jsonbytes(`"user":{"id":%d,"username":"%s"},"nick":"%s"`, id, username, nick)
		}

		memberID := Snowflake(345)
		nick := "sjghsfg"
		username := "dfs"

		data := jsonbytes(`{"guild_id":%d,%s}`, id, memberJson(memberID, nick, username))

		evt, err := cacheDispatcher(cache, EvtGuildMemberAdd, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		if _, ok := evt.(*GuildMemberAdd); !ok {
			t.Fatal("unable to cast event to GuildMemberAdd type")
		}

		container, ok := cache.Guilds.Store[id]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		if member, ok := container.Members[memberID]; !ok {
			t.Error("member was not stored")
		} else {
			if member == nil {
				t.Fatal("member is nil")
			}
			if member.UserID != memberID {
				t.Error("wrong user id")
			}
			if member.Nick != nick {
				t.Error("incorrect nickname")
			}
		}

		if user, ok := cache.Users.Store[memberID]; !ok {
			t.Error("user was not stored")
		} else {
			if user == nil {
				t.Fatal("user is nil")
			}
			if user.ID != memberID {
				t.Error("wrong user id")
			}
			if user.Username != username {
				t.Error("incorrect username")
			}
		}

		if container.Guild.MemberCount != 1 {
			t.Error("incorrect member count")
		}

		deadlockTest(t, cache, EvtGuildMemberAdd, data)
	})

	t.Run("member update", func(t *testing.T) {
		// I'm uncertain if this event contains decent user data, so I'm not caching the user field
		cache := NewBasicCache()
		cache.Guilds.Store[id] = &guildCacheContainer{
			Guild:   &Guild{ID: id},
			Members: map[Snowflake]*Member{},
		}

		memberJson := func(id Snowflake, nick string) []byte {
			return jsonbytes(`"user":{"id":%d},"nick":"%s"`, id, nick)
		}
		payloadJson := func(guildID, memberID Snowflake, nick string) []byte {
			return jsonbytes(`{"guild_id":%d,%s}`, guildID, memberJson(memberID, nick))
		}

		memberID := Snowflake(23425)
		nick := "sjghsfg"
		data := payloadJson(id, memberID, nick)

		t.Run("unknown member", func(t *testing.T) {
			// when the member does not exist we should create it
			// the create event could have been missed..
			evt, err := cacheDispatcher(cache, EvtGuildMemberUpdate, data)
			if err != nil {
				t.Fatal("failed to create event", err)
			}

			if _, ok := evt.(*GuildMemberUpdate); !ok {
				t.Fatal("unable to cast event to GuildMemberAdd type")
			}

			container, ok := cache.Guilds.Store[id]
			if !ok || container.Guild == nil {
				t.Error("guild was not cached")
			}

			if member, ok := container.Members[memberID]; !ok {
				t.Error("member was not stored")
			} else {
				if member == nil {
					t.Fatal("member is nil")
				}
				if member.UserID != memberID {
					t.Error("wrong user id")
				}
				if member.Nick != nick {
					t.Error("incorrect nickname")
				}
			}

			if container.Guild.MemberCount != 1 {
				t.Error("incorrect member count")
			}
		})

		t.Run("existing member", func(t *testing.T) {
			// when the member does not exist we should create it
			// the create event could have been missed..
			updatedNick := "andy bandy"
			data = payloadJson(id, memberID, updatedNick)
			evt, err := cacheDispatcher(cache, EvtGuildMemberUpdate, data)
			if err != nil {
				t.Fatal("failed to create event", err)
			}

			if _, ok := evt.(*GuildMemberUpdate); !ok {
				t.Fatal("unable to cast event to GuildMemberAdd type")
			}

			container, ok := cache.Guilds.Store[id]
			if !ok || container.Guild == nil {
				t.Error("guild was not in cached")
			}

			if member, ok := container.Members[memberID]; !ok {
				t.Error("member was not stored")
			} else {
				if member == nil {
					t.Fatal("member is nil")
				}
				if member.UserID != memberID {
					t.Error("wrong user id")
				}
				if member.Nick != updatedNick {
					t.Error("incorrect nickname")
				}
			}

			if container.Guild.MemberCount != 1 {
				t.Error("incorrect member count")
			}
		})

		deadlockTest(t, cache, EvtGuildMemberAdd, data)
	})

	t.Run("member remove", func(t *testing.T) {
		memberID := Snowflake(345)

		cache := NewBasicCache()
		cache.Guilds.Store[id] = &guildCacheContainer{
			Guild: &Guild{ID: id, MemberCount: 1},
			Members: map[Snowflake]*Member{
				memberID: {UserID: memberID, Nick: "test"},
			},
		}

		data := jsonbytes(`{"user":{"id":%d},"guild_id":%d}`, memberID, id)

		evt, err := cacheDispatcher(cache, EvtGuildMemberRemove, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		if _, ok := evt.(*GuildMemberRemove); !ok {
			t.Fatal("unable to cast event to GuildMemberRemove type")
		}

		container, ok := cache.Guilds.Store[id]
		if !ok || container.Guild == nil {
			t.Error("guild was not cached")
		}

		if _, ok := container.Members[memberID]; ok {
			t.Error("member was not deleted")
		}

		if container.Guild.MemberCount != 0 {
			t.Error("incorrect member count")
		}

		deadlockTest(t, cache, EvtGuildMemberRemove, data)
	})

}

func TestBasicCache_GuildRoles(t *testing.T) {
	t.Run("create without guild", func(t *testing.T) {
		cache := NewBasicCache()

		guildID := Snowflake(3546)
		roleID := Snowflake(5)
		name := "test"
		position := 3
		data := jsonbytes(`{"guild_id":%d,"role":{"id":%d,"name":"%s","position":%d}}`, guildID, roleID, name, position)

		evt, err := cacheDispatcher(cache, EvtGuildRoleCreate, data)
		if err != nil {
			t.Fatal("failed to create event", err)
		}

		roleCreate := evt.(*GuildRoleCreate)
		role := roleCreate.Role

		if roleCreate.GuildID != guildID {
			t.Fatal("incorrect guild id")
		}
		if role.ID != roleID {
			t.Fatal("incorrect role id")
		}

		if len(cache.Guilds.Store) != 0 {
			t.Error("a guild object was created, expected none to be created")
		}

		deadlockTest(t, cache, EvtGuildRoleCreate, data)
	})
	t.Run("create", func(t *testing.T) {
		guildID := Snowflake(3546)

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild: &Guild{ID: guildID},
		}

		roleID := Snowflake(5)
		name := "test"
		position := 3
		data := jsonbytes(`{"guild_id":%d,"role":{"id":%d,"name":"%s","position":%d}}`, guildID, roleID, name, position)

		if _, err := cacheDispatcher(cache, EvtGuildRoleCreate, data); err != nil {
			t.Fatal("failed to create event", err)
		}

		if len(cache.Guilds.Store) != 1 {
			t.Fatal("missing guild")
		}

		roles := cache.Guilds.Store[guildID].Guild.Roles
		if len(roles) != 1 {
			t.Fatal("role was not cached")
		}

		role := roles[0]
		if role.ID != roleID {
			t.Error("incorrect role id")
		}
		if role.Name != name {
			t.Error("incorrect role name")
		}
		if role.Position != position {
			t.Error("incorrect role position")
		}

		deadlockTest(t, cache, EvtGuildRoleCreate, data)
	})
	t.Run("update", func(t *testing.T) {
		guildID := Snowflake(3546)

		roleID := Snowflake(436345)
		roleName := "testing"

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild: &Guild{
				ID: guildID,
				Roles: []*Role{
					{ID: roleID, Name: roleName},
				},
			},
		}

		updatedRoleName := "test x2222"
		data := jsonbytes(`{"guild_id":%d,"role":{"id":%d,"name":"%s"}}`, guildID, roleID, updatedRoleName)

		if _, err := cacheDispatcher(cache, EvtGuildRoleUpdate, data); err != nil {
			t.Fatal("failed to create event", err)
		}

		if len(cache.Guilds.Store) != 1 {
			t.Fatal("missing guild")
		}

		roles := cache.Guilds.Store[guildID].Guild.Roles
		if len(roles) != 1 {
			t.Fatal("missing role from cache")
		}

		role := roles[0]
		if role.ID != roleID {
			t.Error("incorrect role id")
		}
		if role.Name != updatedRoleName {
			t.Error("incorrect role name")
		}

		deadlockTest(t, cache, EvtGuildRoleCreate, data)

		// what happens when we only update the position?
		// the name should still exist
		position := 35
		data = jsonbytes(`{"guild_id":%d,"role":{"id":%d,"position":%d}}`, guildID, roleID, position)
		if _, err := cacheDispatcher(cache, EvtGuildRoleUpdate, data); err != nil {
			t.Fatal("failed to create event", err)
		}

		role = roles[0]
		if role.ID != roleID {
			t.Error("incorrect role id")
		}
		if role.Name != updatedRoleName {
			t.Error("incorrect role name")
		}
		if role.Position != position {
			t.Error("incorrect role position")
		}

	})
	t.Run("delete", func(t *testing.T) {
		guildID := Snowflake(3546)
		roleID := Snowflake(5)

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild: &Guild{
				ID: guildID,
				Roles: []*Role{
					{ID: roleID, Name: "test"},
				},
			},
		}

		data := jsonbytes(`{"guild_id":%d,"role_id":%d}`, guildID, roleID)

		if _, err := cacheDispatcher(cache, EvtGuildRoleDelete, data); err != nil {
			t.Fatal("failed to create event", err)
		}

		if len(cache.Guilds.Store) != 1 {
			t.Fatal("missing guild")
		}

		roles := cache.Guilds.Store[guildID].Guild.Roles
		if len(roles) != 0 {
			t.Fatal("role was not deleted")
		}

		deadlockTest(t, cache, EvtGuildRoleDelete, data)
	})
	t.Run("delete without guild", func(t *testing.T) {
		cache := NewBasicCache()

		guildID := Snowflake(3546)
		roleID := Snowflake(5)
		data := jsonbytes(`{"guild_id":%d,"role_id":%d}`, guildID, roleID)

		if _, err := cacheDispatcher(cache, EvtGuildRoleDelete, data); err != nil {
			t.Fatal("failed to create event", err)
		}

		if len(cache.Guilds.Store) != 0 {
			t.Fatal("a guild was created")
		}

		deadlockTest(t, cache, EvtGuildRoleDelete, data)
	})
}

func TestBasicCache_GetGuildEmoji(t *testing.T) {
	t.Run("unknown guild", func(t *testing.T) {
		cache := NewBasicCache()

		emoji, err := cache.GetGuildEmoji(0, 0)
		if err == nil {
			t.Fatal("there should be an error..")
		}

		if emoji != nil {
			t.Error("emoji should be nil")
		}
		if !errors.Is(err, ErrCacheMiss) {
			t.Error("errpr type should have been CacheMissErr")
		}

		deadlockGetTest(t, func() {
			_, _ = cache.GetGuildEmoji(0, 0)
		})
	})
	t.Run("unknown emoji", func(t *testing.T) {
		guildID := Snowflake(2)

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild: &Guild{ID: guildID},
		}

		emoji, err := cache.GetGuildEmoji(guildID, 0)
		if err == nil {
			t.Fatal("there should be an error..")
		}

		if emoji != nil {
			t.Error("emoji should be nil")
		}
		if !errors.Is(err, ErrCacheMiss) {
			t.Error("errpr type should have been CacheMissErr")
		}

		deadlockGetTest(t, func() {
			_, _ = cache.GetGuildEmoji(guildID, 0)
		})
	})
	t.Run("existing emoji", func(t *testing.T) {
		guildID := Snowflake(2)
		emojiID := Snowflake(34)

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild: &Guild{
				ID: guildID,
				Emojis: []*Emoji{
					{ID: emojiID},
				},
			},
		}

		emoji, err := cache.GetGuildEmoji(guildID, emojiID)
		if err != nil {
			t.Fatal("this should succeed")
		}
		if emoji == nil {
			t.Fatal("emoji should not be nil")
		}
		if emoji.ID != emojiID {
			t.Error("incorrect emoji id")
		}

		t.Run("read only", func(t *testing.T) {
			cachedEmoji := cache.Guilds.Store[guildID].Guild.Emojis[0]
			if emoji == cachedEmoji {
				t.Error("emoji address is shared")
			}
		})

		deadlockGetTest(t, func() {
			_, _ = cache.GetGuildEmoji(guildID, emojiID)
		})
	})
}

func TestBasicCache_GetGuildEmojis(t *testing.T) {
	t.Run("unknown guild", func(t *testing.T) {
		cache := NewBasicCache()

		emoji, err := cache.GetGuildEmojis(0)
		if err == nil {
			t.Fatal("there should be an error..")
		}

		if emoji != nil {
			t.Error("emoji should be nil")
		}
		if !errors.Is(err, ErrCacheMiss) {
			t.Error("errpr type should have been CacheMissErr")
		}

		deadlockGetTest(t, func() {
			_, _ = cache.GetGuildEmojis(0)
		})
	})
	t.Run("no emojis", func(t *testing.T) {
		guildID := Snowflake(2)

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild: &Guild{ID: guildID},
		}

		emojis, err := cache.GetGuildEmojis(guildID)
		if err != nil {
			t.Fatal("got err")
		}

		if len(emojis) != 0 {
			t.Error("emojis should be nil")
		}

		deadlockGetTest(t, func() {
			_, _ = cache.GetGuildEmojis(guildID)
		})
	})
	t.Run("existing emojis", func(t *testing.T) {
		guildID := Snowflake(2)
		emojiID := Snowflake(34)

		cache := NewBasicCache()
		cache.Guilds.Store[guildID] = &guildCacheContainer{
			Guild: &Guild{
				ID: guildID,
				Emojis: []*Emoji{
					{ID: emojiID},
				},
			},
		}

		emojis, err := cache.GetGuildEmojis(guildID)
		if err != nil {
			t.Fatal("this should succeed")
		}
		if len(emojis) == 0 {
			t.Fatal("emojis should not be nil")
		}
		if emojis[0].ID != emojiID {
			t.Error("incorrect emoji id")
		}

		t.Run("read only", func(t *testing.T) {
			cachedEmoji := cache.Guilds.Store[guildID].Guild.Emojis[0]
			if emojis[0] == cachedEmoji {
				t.Error("emoji address is shared")
			}
		})

		deadlockGetTest(t, func() {
			_, _ = cache.GetGuildEmojis(guildID)
		})
	})
}
