package disgord

import (
	"io/ioutil"
	"testing"
)

func TestChannel_DeepCopy(t *testing.T) {
	test := NewChannel()
	icon1 := "sdljfdsjf"
	test.Icon = &icon1
	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: "first",
	})
	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: "second",
	})

	copy := test.DeepCopy().(*Channel)
	icon2 := "sfkjdsf"
	test.Icon = &icon2
	if *copy.Icon != icon1 {
		t.Error("deep copy failed")
	}

	test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
		Type: "third",
	})
	if len(copy.PermissionOverwrites) != 2 {
		t.Error("deep copy failed")
	}
}

func verifyChannelUnmarshal(t *testing.T, data []byte) {
	v := Channel{}
	err := validateJSONMarshalling(data, &v)
	check(err, t)
}

func checkForChannelUnmarshalErr(t *testing.T, data []byte) {
	v := Channel{}
	if err := unmarshal(data, &v); err != nil {
		t.Error(err)
	}
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
