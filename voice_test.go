// +build !integration

package disgord

import (
	"io/ioutil"
	"testing"
)

func TestStateMarshalling(t *testing.T) {
	unmarshal := createUnmarshalUpdater(defaultUnmarshaler)

	data, err := ioutil.ReadFile("testdata/voice/state1.json")
	check(err, t)

	state := VoiceState{}
	err = unmarshal(data, &state)
	check(err, t)
}
