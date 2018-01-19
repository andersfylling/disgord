package voice_test

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
	"github.com/andersfylling/disgord/voice"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("examples/state1.json")
	testutil.Check(err, t)

	state := voice.State{}
	err = testutil.ValidateJSONMarshalling(data, &state)
	testutil.Check(err, t)
}
