// +build !integration

package disgord

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/json"
)

func TestChannel_DeepCopy(t *testing.T) {
	test := NewChannel()
	icon1 := "sdljfdsjf"
	test.Icon = icon1
	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: "first",
	})
	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: "second",
	})

	copy := test.DeepCopy().(*Channel)
	icon2 := "sfkjdsf"
	test.Icon = icon2
	if copy.Icon != icon1 {
		t.Error("deep copy failed")
	}

	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: "third",
	})
	if len(copy.PermissionOverwrites) != 2 {
		t.Error("deep copy failed")
	}
}

func checkForChannelUnmarshalErr(t *testing.T, data []byte) {
	v := Channel{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Error(err)
	}
	executeInternalUpdater(v)
}

func TestChannel_UnmarshalJSON(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		data, err := ioutil.ReadFile("testdata/channel/channel_create.json")
		check(err, t)
		checkForChannelUnmarshalErr(t, data)
	})
	t.Run("update", func(t *testing.T) {
		files := []string{
			"testdata/channel/update_name.json",
			"testdata/channel/update_nsfw1.json",
			"testdata/channel/update_nsfw2.json",
			"testdata/channel/update_ratelimit.json",
			"testdata/channel/update_ratelimit_removed.json",
			"testdata/channel/update_topic.json",
			"testdata/channel/update_topic_removed.json",
		}
		for _, file := range files {
			data, err := ioutil.ReadFile(file)
			check(err, t)
			checkForChannelUnmarshalErr(t, data)
		}
	})
	t.Run("delete", func(t *testing.T) {
		data, err := ioutil.ReadFile("testdata/channel/delete.json")
		check(err, t)
		checkForChannelUnmarshalErr(t, data)
	})
}

func TestChannel_saveToDiscord(t *testing.T) {

}
func TestChannel_JSONIconNull(t *testing.T) {
	// check if null's in json are parsed as an empty string
	data := []byte(`{"id":"324234235","type":1,"icon":null}`)
	var c *struct {
		ID   Snowflake `json:"id"`
		Type int       `json:"type"`
		Icon string    `json:"icon"`
	}
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	executeInternalUpdater(c)

	if c.Icon != "" {
		t.Error(c.Icon, "was not empty")
	}
}
