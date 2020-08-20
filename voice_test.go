// +build !integration

package disgord

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/json"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/voice/state1.json")
	check(err, t)

	state := VoiceState{}
	err = json.Unmarshal(data, &state)
	check(err, t)
}
