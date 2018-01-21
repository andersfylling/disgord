package guild

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestGuildMemberMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/member1.json")
	testutil.Check(err, t)

	v := Member{}
	err = testutil.ValidateJSONMarshalling(data, &v)
	testutil.Check(err, t)
}
