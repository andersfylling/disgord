package disgord

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/httd"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/voice/state1.json")
	check(err, t)

	state := VoiceState{}
	err = httd.Unmarshal(data, &state)
	check(err, t)
}
