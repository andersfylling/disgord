package rest

import (
	"os"
	"testing"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/endpoint"
	"github.com/andersfylling/disgord/rest/httd"
	. "github.com/andersfylling/snowflake"
)

const randomBase64Emoji = "data:image/png;base64,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"

func TestListGuildEmojis(t *testing.T) {
	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	guildIDStr := os.Getenv(constant.DisgordTestGuildDefault)
	guildID, err := GetSnowflake(guildIDStr)
	if err != nil {
		t.Skip()
		return
	}

	emojis, err := ListGuildEmojis(client, guildID)
	if err != nil {
		t.Error(err)
	}

	if len(emojis) != 1 {
		t.Error("expected guild to have one emoji")
	}
}

func TestGetGuildEmoji(t *testing.T) {
	guildIDStr := os.Getenv(constant.DisgordTestGuildDefault)
	guildID, err := GetSnowflake(guildIDStr)
	if err != nil {
		t.Skip()
		return
	}
	emojiIDStr := os.Getenv(constant.DisgordTestGuildDefaultEmojiSnowflake)
	emojiID, err := GetSnowflake(emojiIDStr)
	if err != nil {
		t.Skip()
		return
	}

	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	emoji, err := GetGuildEmoji(client, guildID, emojiID)
	if err != nil {
		t.Error(err)
	}

	if emoji.ID != emojiID {
		t.Error("emoji ID missmatch")
	}
}

func TestCreateAndDeleteGuildEmoji(t *testing.T) {
	var emoji *resource.Emoji
	var err error

	guildIDStr := os.Getenv(constant.DisgordTestGuildAdmin)
	guildID, err := GetSnowflake(guildIDStr)
	if err != nil {
		t.Skip()
		return
	}

	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	t.Run("create emoji", func(t *testing.T) {
		params := &CreateGuildEmojiParams{
			Name:  "testing4324",
			Image: randomBase64Emoji,
		}
		emoji, err = CreateGuildEmoji(client, guildID, params)
		if err != nil {
			t.Error(err)
		}

		if emoji.ID.Empty() {
			t.Error("emoji ID missing")
		}
	})

	t.Run("verify creation", func(t *testing.T) {
		_, err = GetGuildEmoji(client, guildID, emoji.ID)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("delete created emoji", func(t *testing.T) {
		err = DeleteGuildEmoji(client, guildID, emoji.ID)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("verify deletion", func(t *testing.T) {
		_, err = GetGuildEmoji(client, guildID, emoji.ID)
		if err == nil {
			t.Error("emoji was found when it should have been deleted")
		}
	})
}

func TestModifyGuildEmoji(t *testing.T) {
	var emoji *resource.Emoji
	var err error
	newName := "asldjjkasd" // "super-emoji-ok" <- causes regex issue

	guildIDStr := os.Getenv(constant.DisgordTestGuildAdmin)
	guildID, err := GetSnowflake(guildIDStr)
	if err != nil {
		t.Skip()
		return
	}

	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	t.Run("create emoji", func(t *testing.T) {
		params := &CreateGuildEmojiParams{
			Name:  "test6547465",
			Image: randomBase64Emoji,
		}
		emoji, err = CreateGuildEmoji(client, guildID, params)
		if err != nil {
			t.Error(err)
		}

		if emoji.ID.Empty() {
			t.Error("emoji ID missing")
		}
	})

	t.Run("verify creation", func(t *testing.T) {
		_, err = GetGuildEmoji(client, guildID, emoji.ID)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("modify emoji", func(t *testing.T) {
		params := &ModifyGuildEmojiParams{
			Name: newName,
		}
		_, err = ModifyGuildEmoji(client, guildID, emoji.ID, params)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("verify modification", func(t *testing.T) {
		emoji, err = GetGuildEmoji(client, guildID, emoji.ID)
		if err != nil {
			t.Error(err)
		}

		if emoji.Name != newName {
			t.Error("emoji name was not updated")
		}
	})

	t.Run("delete created emoji", func(t *testing.T) {
		err = DeleteGuildEmoji(client, guildID, emoji.ID)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("verify deletion", func(t *testing.T) {
		_, err = GetGuildEmoji(client, guildID, emoji.ID)
		if err == nil {
			t.Error("emoji was found when it should have been deleted")
		}
	})
}

func TestValidEmojiName(t *testing.T) {
	var emoji *resource.Emoji
	var err error

	illegalNames := []string{
		"testing-this-thing-here",
	}

	guildIDStr := os.Getenv(constant.DisgordTestGuildAdmin)
	guildID, err := GetSnowflake(guildIDStr)
	if err != nil {
		t.Skip()
		return
	}

	client, err := createTestRequester()
	if err != nil {
		t.Skip()
		return
	}

	var mustDelete bool = false
	t.Run("create emoji", func(t *testing.T) {
		params := &CreateGuildEmojiParams{
			Name:  illegalNames[0],
			Image: randomBase64Emoji,
		}

		_, body, err := client.Post(&httd.Request{
			Ratelimiter: httd.RatelimitGuild(guildID),
			Endpoint:    endpoint.GuildEmojis(guildID),
			JSONParams:  params,
		})
		if err != nil {
			return
		}

		err = unmarshal(body, emoji)
		if err == nil {
			t.Error("discord does accept '-' in emoji names now. Please update validEmojiName()")
			mustDelete = true
		} else {
			t.Log(err)
		}
	})

	t.Run("verify creation", func(t *testing.T) {
		if !mustDelete {
			t.Skip()
			return
		}
		_, err = GetGuildEmoji(client, guildID, emoji.ID)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("delete created emoji", func(t *testing.T) {
		if !mustDelete {
			t.Skip()
			return
		}
		err = DeleteGuildEmoji(client, guildID, emoji.ID)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("verify deletion", func(t *testing.T) {
		if !mustDelete {
			t.Skip()
			return
		}
		_, err = GetGuildEmoji(client, guildID, emoji.ID)
		if err == nil {
			t.Error("emoji was found when it should have been deleted")
		}
	})
}
