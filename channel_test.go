package disgord

import (
	"io/ioutil"
	"testing"

	"github.com/andersfylling/disgord/httd"
)

func TestChannel_InterfaceImplementations(t *testing.T) {
	var c interface{} = &Channel{}

	t.Run("DeepCopier", func(t *testing.T) {
		if _, ok := c.(DeepCopier); !ok {
			t.Error("Channel does not implement DeepCopier")
		}

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

func TestChannel_saveToDiscord(t *testing.T) {

}

func TestModifyChannelParams(t *testing.T) {
	t.Run("type-all", func(t *testing.T) {
		params := UpdateChannelParams{}
		var err error

		params.SetName("test")
		params.SetPosition(2)
		params.SetPermissionOverwrites([]PermissionOverwrite{})

		if _, exists := params.data["name"]; !exists {
			t.Error("missing name key")
		}
		if _, exists := params.data["position"]; !exists {
			t.Error("missing position key")
		}
		if _, exists := params.data["permission_overwrites"]; !exists {
			t.Error("missing permission_overwrites key")
		}

		// invalid content
		err = params.SetName("a")
		if err == nil {
			t.Error("expected insert to fail")
		}
		err = params.SetParentID(342)
		if err == nil {
			t.Error("expected insert to fail")
		}
		err = params.RemoveParentID()
		if err == nil {
			t.Error("expected insert to fail")
		}
	})
	t.Run("type-text", func(t *testing.T) {
		params := NewModifyTextChannelParams()
		var err error

		err = params.SetBitrate(9000)
		if err == nil {
			t.Error("expected change to voice only attribute to fail for text channel")
		}

		err = params.SetUserLimit(4)
		if err == nil {
			t.Error("expected change to voice only attribute to fail for text channel")
		}

		params.SetRateLimitPerUser(0)
		if _, exists := params.data["rate_limit_per_user"]; !exists {
			t.Error("missing rate_limit_per_user key")
		}
	})
	t.Run("type-voice", func(t *testing.T) {
		params := NewUpdateVoiceChannelParams()
		var err error

		err = params.SetBitrate(9000)
		if err != nil {
			t.Error(err)
		}

		err = params.SetUserLimit(4)
		if err != nil {
			t.Error(err)
		}

		err = params.SetTopic("test")
		if err == nil {
			t.Error("expected change to voice only attribute to fail for text channel")
		}
	})
	t.Run("empty-marshal", func(t *testing.T) {
		params := UpdateChannelParams{}
		data, err := httd.Marshal(params)
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != "{}" {
			t.Error("expected an empty json object")
		}
	})
}
