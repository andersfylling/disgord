package disgord

import (
	"encoding/json"
	"testing"
)

func TestBanObject(t *testing.T) {
	jsonStr := "{\"reason\":\"mentioning b1nzy\",\"user\":{\"username\":\"Mason\",\"discriminator\":\"9999\",\"id\":\"53908099506183680\",\"avatar\":\"a_bab14f271d565501444b2ca3be944b25\"}}"

	var res = []byte(jsonStr)

	ban := &Ban{}
	err := json.Unmarshal(res, ban)
	if err != nil {
		t.Error(err)
	}

	data, err := json.Marshal(ban)
	if err != nil {
		t.Error(err)
	}

	// match json structures
	if string(res) != string(data) {
		t.Errorf("json data differs. \nGot   %s, \nwants %s", string(data), string(res))
	}
}
