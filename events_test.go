// +build !integration

package disgord

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/json"
)

func TestPrepareBox(t *testing.T) {
	injectRandomEvents(t, func(name string, evt interface{}) error {
		//prepareBox(name, evt) // removed
		executeInternalUpdater(evt)
		return nil
	})
}

func TestChannelCreate_UnmarshalJSON(t *testing.T) {
	channel := &Channel{}
	evt := &ChannelCreate{}

	data, err := ioutil.ReadFile("testdata/channel/channel_create.json")
	check(err, t)

	if err = json.Unmarshal(data, channel); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(channel)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

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

	if err = json.Unmarshal(data, channel); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(channel)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

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

	if err = json.Unmarshal(data, channel); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(channel)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

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

	if err = json.Unmarshal(data, message); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(message)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

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

	if err = json.Unmarshal(data, message); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(message)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

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

	if err = json.Unmarshal(data, message); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(message)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

	if evt.MessageID != message.ID {
		t.Error("different ID")
	}

	if evt.ChannelID != message.ChannelID {
		t.Error("different channel ID")
	}

	if evt.GuildID.IsZero() {
		t.Error("expected guild id to be set")
	}
}

func TestGuildCreate_UnmarshalJSON(t *testing.T) {
	guild := &Guild{}
	evt := &GuildCreate{}

	data, err := ioutil.ReadFile("testdata/guild/create.json")
	check(err, t)

	if err = json.Unmarshal(data, guild); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(guild)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

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

	if err = json.Unmarshal(data, guild); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(guild)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

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

	if err = json.Unmarshal(data, guild); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(guild)

	if err = json.Unmarshal(data, evt); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(evt)

	if evt.UnavailableGuild.ID != guild.ID {
		t.Error("different ID")
	}
}
