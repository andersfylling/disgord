package disgord

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGuildMarshal(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/guild1.json")
	check(err, t)

	v := Guild{}
	err = validateJSONMarshalling(data, &v)
	check(err, t)
}

func TestGuildMarshalUnavailable(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/guildUnavailable1.json")
	check(err, t)

	v := Guild{}
	err = validateJSONMarshalling(data, &v)
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
	err = validateJSONMarshalling(data, &ban)
	check(err, t)
}

// --------
func TestGuildEmbed(t *testing.T) {
	res := []byte("{\"enabled\":true,\"channel_id\":\"41771983444115456\"}")

	// convert to struct
	guildEmbed := GuildEmbed{}
	err := json.Unmarshal(res, &guildEmbed)
	if err != nil {
		t.Error(err)
	}

	// back to json
	data, err := json.Marshal(&guildEmbed)
	if err != nil {
		t.Error(err)
	}

	// match
	if string(res) != string(data) {
		t.Errorf("json data differs. Got %s, wants %s", string(data), string(res))
	}
}

// -------------
func TestGuildMemberMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/member1.json")
	check(err, t)

	v := Member{}
	err = validateJSONMarshalling(data, &v)
	check(err, t)
}
