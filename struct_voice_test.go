package disgord

import (
	"io/ioutil"
	"testing"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/voice/state1.json")
	check(err, t)

	state := VoiceState{}
	err = validateJSONMarshalling(data, &state)
	check(err, t)
}
