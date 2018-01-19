package guild_test

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/guild"
	"github.com/andersfylling/disgord/testutil"
)

func TestGuildMemberMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("examples/member1.json")
	testutil.Check(err, t)

	v := guild.Member{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}
