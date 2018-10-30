package disgord

import (
	"github.com/andersfylling/disgord/httd"
	"testing"
)

const channelID = 486833611564253186

func TestGetChannel(t *testing.T) {
	client, _, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	channel, err := GetChannel(client, channelID)
	if err != nil {
		t.Error(err)
		return
	}

	if channel == nil {
		t.Error("channel was nil")
		return
	}

	if channel.ID != channelID {
		t.Error("incorrect channel id")
	}
}

func TestCreateModifyDeleteChannel(t *testing.T) {
	client, keys, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	var channelID Snowflake

	t.Run("create", func(t *testing.T) {
		channel, err := CreateGuildChannel(client, keys.GuildAdmin, &CreateGuildChannelParams{
			Name: "test",
		})
		if err != nil {
			t.Skip("cannot create channel, therefore skipped")
			return
		}

		channelID = channel.ID
	})

	if channelID.Empty() {
		return
	}

	t.Run("modify", func(t *testing.T) {
		changes := NewModifyTextChannelParams()
		changes.SetName("hello")
		channel, err := ModifyChannel(client, channelID, changes)
		if err != nil {
			t.Error(err)
		}
		if channel == nil {
			t.Error("channel was nil")
		}
	})

	t.Run("delete", func(t *testing.T) {
		channel, err := DeleteChannel(client, channelID)
		if err != nil {
			t.Error(err)
		}
		if channel.ID != channelID {
			t.Error("incorrect channel id")
		}

		_, err = GetChannel(client, channelID)
		if err == nil {
			t.Error("able to retrieve deleted channel")
		}
	})
}

func TestModifyChannelParams(t *testing.T) {
	t.Run("type-all", func(t *testing.T) {
		params := ModifyChannelParams{}
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
		params := NewModifyVoiceChannelParams()
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
		params := ModifyChannelParams{}
		data, err := httd.Marshal(params)
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != "{}" {
			t.Error("expected an empty json object")
		}
	})
}
