package guild_test

import (
	"encoding/json"
	"testing"

	"github.com/andersfylling/disgord/guild"
)

func TestGuildEmbed(t *testing.T) {
	res := []byte("{\"enabled\":true,\"channel_id\":\"41771983444115456\"}")

	// convert to struct
	guildEmbed := guild.Embed{}
	err := json.Unmarshal(res, &guildEmbed)
	if err != nil {
		t.Error(err)
	}

	// back to json
	data, err := json.Marshal(&guildEmbed)
	if err != nil {
		t.Error(err)
	}

	// match
	if string(res) != string(data) {
		t.Errorf("json data differs. Got %s, wants %s", string(data), string(res))
	}
}
