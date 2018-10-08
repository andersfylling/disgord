package disgord

import (
	"io/ioutil"
	"testing"
)

func TestChannel_InterfaceImplementations(t *testing.T) {
	var c interface{} = &Channel{}

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := c.(DeepCopier); !ok {
			t.Error("Channel does not implement DeepCopier")
		}

		test := NewChannel()
		test.Icon = "sdljfdsjf"
		test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
			Type: "first",
		})
		test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
			Type: "second",
		})

		copy := test.DeepCopy().(*Channel)
		test.Icon = "sfkjdsf"
		if copy.Icon != "sdljfdsjf" {
			t.Error("deep copy failed")
		}

		test.PermissionOverwrites = append(test.PermissionOverwrites, PermissionOverwrite{
			Type: "third",
		})
		if len(copy.PermissionOverwrites) != 2 {
			t.Error("deep copy failed")
		}
	})

	t.Run("Copier", func(t *testing.T) {
		if _, ok := c.(Copier); !ok {
			t.Error("Channel does not implement Copier")
		}
	})

	t.Run("discordSaver", func(t *testing.T) {
		if _, ok := c.(discordSaver); !ok {
			t.Error("Channel does not implement discordSaver")
		}
	})

	t.Run("discordDeleter", func(t *testing.T) {
		if _, ok := c.(discordDeleter); !ok {
			t.Error("Channel does not implement discordDeleter")
		}
	})
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
