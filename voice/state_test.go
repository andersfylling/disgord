package voice

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/state1.json")
	testutil.Check(err, t)

	state := State{}
	err = testutil.ValidateJSONMarshalling(data, &state)
	testutil.Check(err, t)
}
