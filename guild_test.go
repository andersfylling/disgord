// +build !integration

package disgord

import (
	"github.com/andersfylling/disgord/json"
	"io/ioutil"
	"testing"
)

func TestGuildMarshal(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/guild1.json")
	check(err, t)

	v := Guild{}
	err = json.Unmarshal(data, &v)
	check(err, t)
}

func TestGuildMarshalUnavailable(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/guildUnavailable1.json")
	check(err, t)

	v := Guild{}
	err = json.Unmarshal(data, &v)
	check(err, t)
}

func TestGuild_ChannelSorting(t *testing.T) {
	g := &Guild{}
	total := 1000
	for i := total; i > 0; i-- {
		s := NewSnowflake(uint64(i))
		c := &Channel{ID: s}
		g.AddChannel(c)
	}

	chans := g.Channels
	for i := 1; i <= total; i++ {
		if chans[i-1].ID != NewSnowflake(uint64(i)) {
			t.Error("wrong order")
			break
		}
	}
}

// ---------
func TestGuildBanObject(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/ban1.json")
	check(err, t)

	ban := Ban{}
	err = json.Unmarshal(data, &ban)
	check(err, t)
}

// --------
func TestGuildEmbed(t *testing.T) {
	res := []byte("{\"enabled\":true,\"channel_id\":\"41771983444115456\"}")
	expects := []byte("{\"enabled\":true,\"channel_id\":41771983444115456}")

	// convert to struct
	guildEmbed := GuildEmbed{}
	if err := json.Unmarshal(res, &guildEmbed); err != nil {
		t.Error(err)
	}

	// back to json
	data, err := json.Marshal(&guildEmbed)
	if err != nil {
		t.Error(err)
	}

	// match
	if string(expects) != string(data) {
		t.Errorf("json data differs. Got %s, wants %s", string(data), string(expects))
	}
}

// -------------

func TestGuild_sortChannels(t *testing.T) {
	snowflakes := []Snowflake{
		NewSnowflake(6),
		NewSnowflake(65),
		NewSnowflake(324),
		NewSnowflake(5435),
		NewSnowflake(63453),
		NewSnowflake(111111111),
	}

	guild := NewGuild()

	for i := range snowflakes {
		channel := NewChannel()
		channel.ID = snowflakes[len(snowflakes)-1-i] // reverse

		guild.Channels = append(guild.Channels, channel)
	}

	guild.sortChannels()
	for i, c := range guild.Channels {
		if snowflakes[i] != c.ID {
			t.Error("Channels in guild did not sort correctly")
		}
	}
}

func TestGuild_AddChannel(t *testing.T) {
	snowflakes := []Snowflake{
		NewSnowflake(6),
		NewSnowflake(65),
		NewSnowflake(324),
		NewSnowflake(5435),
		NewSnowflake(63453),
		NewSnowflake(111111111),
	}

	guild := NewGuild()

	for i := range snowflakes {
		channel := NewChannel()
		channel.ID = snowflakes[len(snowflakes)-1-i] // reverse

		guild.AddChannel(channel)
	}

	for i, c := range guild.Channels {
		if snowflakes[i] != c.ID {
			t.Error("Channels in guild did not sort correctly")
		}
	}
}

func TestGuild_DeleteChannel(t *testing.T) {
	snowflakes := []Snowflake{
		NewSnowflake(6),
		NewSnowflake(65),
		NewSnowflake(324),
		NewSnowflake(5435),
		NewSnowflake(63453),
		NewSnowflake(111111111),
	}

	guild := NewGuild()

	for i := range snowflakes {
		channel := NewChannel()
		channel.ID = snowflakes[len(snowflakes)-1-i] // reverse

		guild.AddChannel(channel)
	}

	id := snowflakes[3]
	channel := NewChannel()
	channel.ID = id
	guild.DeleteChannel(channel)
	_, err := guild.Channel(id)
	if err == nil {
		t.Error("no error given when requesting a deleted channel")
	}
}

func TestPermissionBit(t *testing.T) {
	// test permission bit checking
	testBits := PermissionSendMessages | PermissionReadMessages
	if testBits.Contains(PermissionAdministrator) {
		t.Fatal("does not have administrator")
	}
	if !testBits.Contains(PermissionSendMessages) {
		t.Fatal("does have send messages")
	}
	if !testBits.Contains(PermissionReadMessages) {
		t.Fatal("does have read messages")
	}

	// Test json marshal/unmarshal
	b, err := defaultMarshaler(testBits)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &testBits); err != nil {
		t.Fatal(err)
	}
	executeInternalUpdater(testBits)

	if !testBits.Contains(PermissionReadMessages) {
		t.Fatal("does have read messages")
	}
}
