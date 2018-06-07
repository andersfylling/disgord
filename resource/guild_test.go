package resource

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
	"github.com/andersfylling/snowflake"
)

func TestGuildMarshal(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/guild1.json")
	testutil.Check(err, t)

	v := Guild{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}

func TestGuildMarshalUnavailable(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/guildUnavailable1.json")
	testutil.Check(err, t)

	v := Guild{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}

func TestGuild_ChannelSorting(t *testing.T) {
	g := &Guild{}
	total := 1000
	for i := total; i > 0; i-- {
		c := &Channel{ID: snowflake.NewID(uint64(i))}
		g.AddChannel(c)
	}

	chans := g.Channels
	for i := 1; i <= total; i++ {
		if chans[i-1].ID != snowflake.NewID(uint64(i)) {
			t.Error("wrong order")
			break
		}
	}
}

// ---------
func TestGuildBanObject(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild/ban1.json")
	testutil.Check(err, t)

	ban := GuildBan{}
	err = testutil.ValidateJSONMarshalling(data, &ban)
	testutil.Check(err, t)
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
	testutil.Check(err, t)

	v := Member{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}
