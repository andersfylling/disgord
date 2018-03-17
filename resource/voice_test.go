package resource

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/testutil"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/voice/state1.json")
	testutil.Check(err, t)

	state := VoiceState{}
	err = testutil.ValidateJSONMarshalling(data, &state)
	testutil.Check(err, t)
}
