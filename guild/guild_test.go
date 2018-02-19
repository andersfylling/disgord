package guild

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/testutil"
	"github.com/andersfylling/snowflake"
)

func TestGuildMarshal(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guild1.json")
	testutil.Check(err, t)

	v := Guild{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}

func TestGuildMarshalUnavailable(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/guildUnavailable1.json")
	testutil.Check(err, t)

	v := Guild{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}

func TestGuildChannelSorting(t *testing.T) {
	g := &Guild{}
	total := 1000
	for i := total; i > 0; i-- {
		c := &channel.Channel{ID: snowflake.NewID(uint64(i))}
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
