package disgord

import (
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
		newName := "hello"
		channel, err := ModifyChannel(client, channelID, &ModifyChannelParams{
			Name: &newName,
		})
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
