package disgord

import (
	"encoding/json"
	"testing"
)

func TestGuildMemberMarshalling(t *testing.T) {
	timestamp := "2015-04-26T06:26:56.936000+00:00"
	jsonStr := "{\"user\":{},\"nick\":\"NOT API SUPPORT\",\"roles\":[],\"joined_at\":\"" + timestamp + "\",\"deaf\":false,\"mute\":false}"

	var res = []byte(jsonStr)

	guildMember := &GuildMember{}
	err := json.Unmarshal(res, guildMember)
	if err != nil {
		t.Error(err)
	}

	got := guildMember.JoinedAt.String()
	if got != timestamp {
		t.Errorf("Incorrect formatting of JoinedAt timestamp. Got %s, wants %s", got, timestamp)
	}

	data, err := json.Marshal(guildMember)
	if err != nil {
		t.Error(err)
	}

	// match json structures
	if string(res) != string(data) {
		t.Errorf("json data differs. Got %s, wants %s", string(data), string(res))
	}
}
