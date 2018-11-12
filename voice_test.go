package disgord

import (
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/ratelimit"
	"io/ioutil"
	"net/http"
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

func TestListVoiceRegions(t *testing.T) {
	client, _, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	builder := &listVoiceRegionsBuilder{}
	builder.setup(client, &httd.Request{
		Method:      http.MethodGet,
		Ratelimiter: ratelimit.VoiceRegions(),
		Endpoint:    endpoint.VoiceRegions(),
	}, nil)

	list, err := builder.Execute()
	if err != nil {
		t.Error(err)
	}

	if len(list) == 0 {
		t.Error("expected at least one voice region")
	}
}
