package guild

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
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
