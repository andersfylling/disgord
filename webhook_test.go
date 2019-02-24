package disgord

import (
	"encoding/json"
	"testing"

	"github.com/andersfylling/disgord/httd"
)

func getJSONMap(v interface{}) (map[string]*json.RawMessage, error) {
	data, err := httd.Marshal(v)
	if err != nil {
		return nil, err
	}

	var objmap map[string]*json.RawMessage
	err = httd.Unmarshal(data, &objmap)
	return objmap, err
}

func notContain(t *testing.T, list map[string]*json.RawMessage, key string) {
	if _, exists := list[key]; exists {
		t.Error(key + " should not be set")
	}
}

func contain(t *testing.T, list map[string]*json.RawMessage, key string) {
	if _, exists := list[key]; !exists {
		t.Error(key + " should be set")
	}
}

func TestModifyWebhook(t *testing.T) {
	t.Run("params", func(t *testing.T) {
		var params *UpdateWebhookParams
		var partial map[string]*json.RawMessage
		var err error

		// 1
		params = &UpdateWebhookParams{}
		partial, err = getJSONMap(params)
		if err != nil {
			t.Fatal(err)
		}

		notContain(t, partial, "channel_id")
		notContain(t, partial, "avatar")
		notContain(t, partial, "name")

		// 2
		params = &UpdateWebhookParams{}
		params.SetChannelID(45363)
		partial, err = getJSONMap(params)
		if err != nil {
			t.Fatal(err)
		}

		contain(t, partial, "channel_id")
		notContain(t, partial, "avatar")
		notContain(t, partial, "name")

		// 3
		params = &UpdateWebhookParams{}
		params.SetName("shfisudhf")
		partial, err = getJSONMap(params)
		if err != nil {
			t.Fatal(err)
		}

		notContain(t, partial, "channel_id")
		notContain(t, partial, "avatar")
		contain(t, partial, "name")

		// 4
		params = &UpdateWebhookParams{}
		params.SetAvatar("hfjhsdfklsahkfjsdhfksdhf")
		partial, err = getJSONMap(params)
		if err != nil {
			t.Fatal(err)
		}

		notContain(t, partial, "channel_id")
		contain(t, partial, "avatar")
		notContain(t, partial, "name")

		// 5
		params = &UpdateWebhookParams{}
		params.UseDefaultAvatar()
		partial, err = getJSONMap(params)
		if err != nil {
			t.Fatal(err)
		}

		notContain(t, partial, "channel_id")
		contain(t, partial, "avatar")
		notContain(t, partial, "name")
	})
}
