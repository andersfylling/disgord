package disgord

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/andersfylling/snowflake/v3"
)

func TestPrepareBox(t *testing.T) {
	injectRandomEvents(t, func(name string, evt interface{}) error {
		//prepareBox(name, evt) // removed
		executeInternalUpdater(evt)
		return nil
	})
}

type mockCacheEvent struct{}

func (m *mockCacheEvent) Update(key cacheRegistry, v interface{}) (err error) {
	return nil
}
func (m *mockCacheEvent) Get(key cacheRegistry, id Snowflake, args ...interface{}) (v interface{}, err error) {
	return nil, nil
}
func (m *mockCacheEvent) UpdateGuildRole(guildID Snowflake, role *Role, messages json.RawMessage) bool {
	return false
}
func (m *mockCacheEvent) DeleteChannel(channelID snowflake.ID)                                      {}
func (m *mockCacheEvent) DeleteGuildChannel(guildID snowflake.ID, channelID snowflake.ID)           {}
func (m *mockCacheEvent) AddGuildChannel(guildID snowflake.ID, channelID snowflake.ID)              {}
func (m *mockCacheEvent) UpdateChannelPin(channelID snowflake.ID, lastPinTimestamp Timestamp)       {}
func (m *mockCacheEvent) DeleteGuild(guildID snowflake.ID)                                          {}
func (m *mockCacheEvent) DeleteGuildRole(guildID snowflake.ID, roleID snowflake.ID)                 {}
func (m *mockCacheEvent) AddGuildRole(GuildID Snowflake, role *Role)                                {}
func (m *mockCacheEvent) UpdateChannelLastMessageID(channelID snowflake.ID, messageID snowflake.ID) {}
func (m *mockCacheEvent) AddGuildMember(guildID snowflake.ID, member *Member)                       {}
func (m *mockCacheEvent) RemoveGuildMember(guildID snowflake.ID, memberID snowflake.ID)             {}
func (m *mockCacheEvent) UpdateMemberAndUser(guildID, userID snowflake.ID, data json.RawMessage)    {}
func (m *mockCacheEvent) SetGuildEmojis(guildID Snowflake, emojis []*Emoji)                         {}
func (m *mockCacheEvent) Updates(key cacheRegistry, vs []interface{}) error {
	return nil
}

func TestCacheEvent(t *testing.T) {
	cache := &mockCacheEvent{}
	injectRandomEvents(t, func(name string, evt interface{}) error {
		return cacheEvent(cache, name, evt, nil)
	})
}

func TestChannelCreate_UnmarshalJSON(t *testing.T) {
	channel := &Channel{}
	evt := &ChannelCreate{}

	data, err := ioutil.ReadFile("testdata/channel/channel_create.json")
	check(err, t)

	err = unmarshal(data, channel)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Channel.Name != channel.Name {
		t.Error("different names")
	}

	if evt.Channel.ID != channel.ID {
		t.Error("different ID")
	}
}

func TestChannelUpdate_UnmarshalJSON(t *testing.T) {
	channel := &Channel{}
	evt := &ChannelUpdate{}

	data, err := ioutil.ReadFile("testdata/channel/update_topic.json")
	check(err, t)

	err = unmarshal(data, channel)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Channel.Name != channel.Name {
		t.Error("different names")
	}

	if evt.Channel.ID != channel.ID {
		t.Error("different ID")
	}
}

func TestChannelDelete_UnmarshalJSON(t *testing.T) {
	channel := &Channel{}
	evt := &ChannelDelete{}

	data, err := ioutil.ReadFile("testdata/channel/delete.json")
	check(err, t)

	err = unmarshal(data, channel)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Channel.Name != channel.Name {
		t.Error("different names")
	}

	if evt.Channel.ID != channel.ID {
		t.Error("different ID")
	}
}

func TestMessageCreate_UnmarshalJSON(t *testing.T) {
	message := &Message{}
	evt := &MessageCreate{}

	data, err := ioutil.ReadFile("testdata/channel/message_create_guild_invite.json")
	check(err, t)

	err = unmarshal(data, message)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Message.Content != message.Content {
		t.Error("different content")
	}

	if evt.Message.ID != message.ID {
		t.Error("different ID")
	}
}

func TestMessageUpdate_UnmarshalJSON(t *testing.T) {
	message := &Message{}
	evt := &MessageUpdate{}

	data, err := ioutil.ReadFile("testdata/channel/message_update.json")
	check(err, t)

	err = unmarshal(data, message)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Message.Content != message.Content {
		t.Error("different content")
	}

	if evt.Message.ID != message.ID {
		t.Error("different ID")
	}
}

func TestMessageDelete_UnmarshalJSON(t *testing.T) {
	message := &Message{}
	evt := &MessageDelete{}

	data, err := ioutil.ReadFile("testdata/channel/message_delete.json")
	check(err, t)

	err = unmarshal(data, message)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.MessageID != message.ID {
		t.Error("different ID")
	}

	if evt.ChannelID != message.ChannelID {
		t.Error("different channel ID")
	}

	if evt.GuildID.Empty() {
		t.Error("expected guild id to be set")
	}
}

func TestGuildCreate_UnmarshalJSON(t *testing.T) {
	guild := &Guild{}
	evt := &GuildCreate{}

	data, err := ioutil.ReadFile("testdata/guild/create.json")
	check(err, t)

	err = unmarshal(data, guild)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Guild.ID != guild.ID {
		t.Error("different ID")
	}

	if evt.Guild.MemberCount != guild.MemberCount {
		t.Error("different member count")
	}
}

func TestGuildUpdate_UnmarshalJSON(t *testing.T) {
	guild := &Guild{}
	evt := &GuildUpdate{}

	data, err := ioutil.ReadFile("testdata/guild/update.json")
	check(err, t)

	err = unmarshal(data, guild)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Guild.ID != guild.ID {
		t.Error("different ID")
	}

	if evt.Guild.MemberCount != guild.MemberCount {
		t.Error("different member count")
	}
}

func TestGuildDelete_UnmarshalJSON(t *testing.T) {
	guild := &Guild{}
	evt := &GuildDelete{}

	data, err := ioutil.ReadFile("testdata/guild/delete_by_kick.json")
	check(err, t)

	err = unmarshal(data, guild)
	if err != nil {
		t.Error(err)
	}

	err = unmarshal(data, evt)
	if err != nil {
		t.Error(err)
	}

	if evt.UnavailableGuild.ID != guild.ID {
		t.Error("different ID")
	}
}
