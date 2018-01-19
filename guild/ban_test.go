package guild_test

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/testutil"
)

func TestBanObject(t *testing.T) {
	data, err := ioutil.ReadFile("examples/ban1.json")
	testutil.Check(err, t)

	ban := guild.Ban{}
	err = testutil.ValidateJSONMarshalling(data, &ban)
	testutil.Check(err, t)
}
