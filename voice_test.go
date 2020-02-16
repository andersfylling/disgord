// +build !integration

package disgord

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/internal/util"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/voice/state1.json")
	check(err, t)

	state := VoiceState{}
	err = util.Unmarshal(data, &state)
	check(err, t)
}
