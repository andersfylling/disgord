package disgord

import (
	"github.com/andersfylling/disgord/httd"
	"io/ioutil"
	"testing"
)

func TestStateMarshalling(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/voice/state1.json")
	check(err, t)

	state := VoiceState{}
	err = httd.Unmarshal(data, &state)
	check(err, t)
}

func TestVoice_InterfaceImplementations(t *testing.T) {
	t.Run("VoiceState", func(t *testing.T) {
		var u interface{} = &VoiceState{}
		t.Run("DeepCopier", func(t *testing.T) {
			if _, ok := u.(DeepCopier); !ok {
				t.Error("does not implement DeepCopier")
			}
		})

		t.Run("Copier", func(t *testing.T) {
			if _, ok := u.(Copier); !ok {
				t.Error("does not implement Copier")
			}
		})
	})
	t.Run("VoiceRegion", func(t *testing.T) {
		var u interface{} = &VoiceRegion{}
		t.Run("DeepCopier", func(t *testing.T) {
			if _, ok := u.(DeepCopier); !ok {
				t.Error("does not implement DeepCopier")
			}
		})

		t.Run("Copier", func(t *testing.T) {
			if _, ok := u.(Copier); !ok {
				t.Error("does not implement Copier")
			}
		})
	})
}
