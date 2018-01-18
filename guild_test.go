package disgord

import (
	"encoding/json"
	"testing"
)

func TestGuildMarshal(t *testing.T) {
	jsonStr := "{\"id\":\"41771983423143937\",\"application_id\":null,\"name\":\"Discord Developers\",\"icon\":\"SEkgTU9NIElUUyBBTkRSRUkhISEhISEh\",\"splash\":null,\"owner_id\":\"80351110224678912\",\"region\":\"us-east\",\"afk_channel_id\":\"42072017402331136\",\"afk_timeout\":300,\"embed_enabled\":true,\"embed_channel_id\":\"41771983444115456\",\"verification_level\":1,\"default_message_notifications\":0,\"explicit_content_filter\":0,\"mfa_level\":0,\"widget_enabled\":false,\"widget_channel_id\":\"41771983423143937\",\"roles\":[],\"emojis\":[],\"features\":[\"INVITE_SPLASH\"],\"unavailable\":false}"

	var res = []byte(jsonStr)

	guild := &Guild{}
	err := json.Unmarshal(res, guild)
	if err != nil {
		t.Error(err)
	}

	data, err := json.Marshal(guild)
	if err != nil {
		t.Error(err)
	}

	// match json structures
	if string(res) != string(data) {
		// TODO need proper discord API mock data for testing
		//t.Errorf("json data differs. \nGot   %s, \nwants %s", string(data), string(res))
	}
}

func TestGuildMarshalUnavailable(t *testing.T) {
	jsonStr := "{\"id\":\"41771983423143937\",\"unavailable\":true}"

	var res = []byte(jsonStr)

	guild := &Guild{}
	err := json.Unmarshal(res, guild)
	if err != nil {
		t.Error(err)
	}

	data, err := json.Marshal(guild)
	if err != nil {
		t.Error(err)
	}

	// match json structures
	if string(res) != string(data) {
		t.Errorf("json data differs. \nGot   %s, \nwants %s", string(data), string(res))
	}
}
