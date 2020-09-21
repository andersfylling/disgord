package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andersfylling/disgord/internal/constant"

	"github.com/andersfylling/disgord"
)

type keys struct {
	GuildAdmin   disgord.Snowflake
	GuildDefault disgord.Snowflake
}

const randomBase64Emoji = "data:image/png;base64,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"

func notARateLimitIssue(err error) bool {
	return !strings.Contains(err.Error(), "You are being rate limited.")
}

func setupKeys() *keys {
	keys := &keys{}

	str1 := os.Getenv(constant.DisgordTestGuildDefault)
	g1, err := disgord.GetSnowflake(str1)
	if err != nil {
		panic("missing default guild id")
	}
	keys.GuildDefault = g1

	str2 := os.Getenv(constant.DisgordTestGuildAdmin)
	g2, err := disgord.GetSnowflake(str2)
	if err != nil {
		panic("missing admin guild id")
	}
	keys.GuildAdmin = g2

	return keys
}

func main() {
	c := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
	})
	keys := setupKeys()

	// -------------------
	// AUDIT-LOGS
	// -------------------
	logs, err := c.Guild(keys.GuildAdmin).GetAuditLogs().Execute()
	if err != nil {
		panic(err)
	} else if logs == nil {
		panic("did not get a datastructure from GetGuildAuditLogs()")
	}

	// -------------------
	// CHANNELS
	// -------------------
	channelID := disgord.Snowflake(486833611564253186)
	func() {
		channel, err := c.Channel(channelID).Get()
		if err != nil {
			panic(err)
		} else if channel == nil {
			panic("channel was nil")
		} else if channel.ID != channelID {
			panic("incorrect channel id")
		}
	}()

	// create
	func() {
		channel, err := c.Guild(keys.GuildAdmin).CreateChannel("test", nil)
		if err != nil {
			panic("cannot create channel, therefore skipped")
		} else if channel.ID.IsZero() {
			panic("channel ID of created channel was empty")
		}

		channelID = channel.ID
	}()

	// modify
	func() {
		channel, err := c.Channel(channelID).Update().SetName("hello").Execute()
		if err != nil {
			panic(err)
		} else if channel == nil {
			panic("channel was nil")
		}
	}()

	// delete
	func() {
		channel, err := c.Channel(channelID).Delete()
		if err != nil {
			panic(err)
		} else if channel.ID != channelID {
			panic("incorrect channel id")
		}

		_, err = c.Channel(channelID).Get()
		if err == nil {
			panic("able to retrieve deleted channel")
		}
	}()

	// -------------------
	// GUILDS
	// -------------------

	// -------------------
	// USERS
	// -------------------

	// TestGetCurrentUser
	func() {
		_, err = c.CurrentUser().Get(disgord.IgnoreCache)
		if err != nil {
			panic(err)
		}
	}()
	// TestGetUser
	func() {
		const userID = 140413331470024704
		user, err := c.User(userID).Get(disgord.IgnoreCache)
		if err != nil {
			panic(err)
		} else if user.ID != userID {
			panic("user ID missmatch")
		}
	}()
	// TestModifyCurrentUser
	func() {
		// this has been verified to work
		// however, you cannot change username often so this is
		// can give an error

		// TODO: rewrite; this was moved from the disgord root pkg, but have not been rewritten to work here yet
		// var originalUsername string
		// t.Run("getting original username", func(t *testing.T) {
		// 	user, err := GetCurrentUser(client)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		//
		// 	originalUsername = user.Username
		// })
		//
		// t.Run("changing username", func(t *testing.T) {
		// 	if originalUsername == "" {
		// 		panic()
		// 		return
		// 	}
		// 	params := &ModifyCurrentUserParams{}
		// 	params.SetUsername("sldfhksghs")
		// 	_, err := ModifyCurrentUser(client, params)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// })
		//
		// t.Run("resetting username", func(t *testing.T) {
		// 	if originalUsername == "" {
		// 		panic()
		// 		return
		// 	}
		// 	params := &ModifyCurrentUserParams{}
		// 	params.SetUsername(originalUsername)
		// 	_, err := ModifyCurrentUser(client, params)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// })
	}()

	// TestLeaveGuild
	func() {
		// Nope. Not gonna automate this.
	}()
	// TestUserDMs
	func() {
		// TODO
	}()

	// TestCreateDM
	func() {
		channel, err := c.User(228846961774559232).CreateDM()
		if err != nil {
			panic(err)
		} else if channel == nil {
			panic("channel was nil")
		}
	}()
	// TestCreateGroupDM
	func() {
		// TODO
	}()
	// TestGetUserConnections
	func() {
		// Missing OAuth2
	}()

	// -------------------
	// EMOJIS
	// -------------------

	// TestListGuildEmojis
	func() {
		emojis, err := c.Guild(keys.GuildDefault).GetEmojis()
		if err != nil && !notARateLimitIssue(err) {
			panic("rate limited")
		}
		if err != nil && notARateLimitIssue(err) {
			panic(err)
		}

		if len(emojis) != 1 && notARateLimitIssue(err) {
			panic("expected guild to have one emoji")
		}
	}()

	// TestGetGuildEmoji
	func() {
		emojiIDStr := os.Getenv(constant.DisgordTestGuildDefaultEmojiSnowflake)
		//fmt.Println(emojiIDStr)
		emojiID, err := disgord.GetSnowflake(emojiIDStr)
		if err != nil {
			panic("snowflake makes no sense")
			return
		}

		emoji, err := c.Guild(keys.GuildDefault).Emoji(emojiID).Get()
		if err != nil && !notARateLimitIssue(err) {
			panic("rate limited")
		} else if err != nil && notARateLimitIssue(err) {
			panic(err)
		} else if emoji == nil {
			panic("emoji was nil...")
		} else if emoji != nil && emoji.ID != emojiID && notARateLimitIssue(err) {
			panic("emoji ID missmatch")
		}
	}()

	// TestCreateAndDeleteGuildEmoji
	func() {
		var emoji *disgord.Emoji
		var err error

		// create emoji
		func() {
			emoji, err = c.Guild(keys.GuildDefault).CreateEmoji(&disgord.CreateGuildEmojiParams{
				Name:  "testing4324",
				Image: randomBase64Emoji,
			})
			if err != nil && !notARateLimitIssue(err) {
				panic("rate limited")
			}
			if err != nil && notARateLimitIssue(err) {
				panic(err)
			}

			if emoji.ID.IsZero() && notARateLimitIssue(err) {
				panic("emoji ID missing")
			}
		}()

		// delete created emoji
		func() {
			err := c.Guild(keys.GuildDefault).Emoji(emoji.ID).Delete()
			if err != nil && !notARateLimitIssue(err) {
				panic("rate limited")
			}
			if err != nil && notARateLimitIssue(err) {
				panic(err)
			}
		}()
	}()

	// TestModifyGuildEmoji
	func() {
		var emoji *disgord.Emoji
		var err error
		newName := "asldjjkasd" // "super-emoji-ok" <- causes regex issue

		// create emoji
		func() {
			emoji, err = c.Guild(keys.GuildDefault).CreateEmoji(&disgord.CreateGuildEmojiParams{Name: "test6547465", Image: randomBase64Emoji})
			if err != nil && !notARateLimitIssue(err) {
				panic("rate limited")
			} else if err != nil && notARateLimitIssue(err) {
				panic(err)
			} else if emoji.ID.IsZero() && notARateLimitIssue(err) {
				panic("emoji ID missing")
			}
		}()

		// modify emoji
		func() {
			_, err = c.Guild(keys.GuildDefault).Emoji(emoji.ID).Update().SetName(newName).Execute()
			if err != nil && !notARateLimitIssue(err) {
				panic("rate limited")
			} else if err != nil && notARateLimitIssue(err) {
				panic(err)
			}
		}()

		// delete created emoji
		func() {
			time.Sleep(1 * time.Second) // just ensure that this get's run
			err = c.Guild(keys.GuildDefault).Emoji(emoji.ID).Delete()
			if err != nil && !notARateLimitIssue(err) {
				panic("rate limited")
			} else if err != nil && notARateLimitIssue(err) {
				panic(err)
			}
		}()
	}()

	// TestValidEmojiName
	func() {
		var emoji *disgord.Emoji
		var err error

		illegalNames := []string{
			"testing-this-thing-here",
		}

		var mustDelete = false
		// create emoji
		func() {

			emoji, err = c.Guild(keys.GuildAdmin).CreateEmoji(&disgord.CreateGuildEmojiParams{
				Name:  illegalNames[0],
				Image: randomBase64Emoji,
			})
			if err != nil && !notARateLimitIssue(err) {
				panic("rate limited")
			} else if err != nil {
				panic(err)
			}

			if err == nil {
				fmt.Println("discord does accept '-' in emoji names now. Please update validEmojiName()")
				mustDelete = true
			} else {
				panic(err)
			}
		}()

		// delete created emoji
		func() {
			if !mustDelete {
				panic("no new emoji created")
			}
			err = c.Guild(keys.GuildDefault).Emoji(emoji.ID).Delete()
			if err != nil && !notARateLimitIssue(err) {
				panic("rate limited")
			}
			if err != nil && notARateLimitIssue(err) {
				panic(err)
			}
		}()
	}()

	// -------------------
	// INVITES
	// -------------------

	// TestGetInvite
	func() {
		// TODO: invite codes....
		// will only return 404 without an invite code
		//
		// inviteCode := ""
		// _, err = GetInvite(client, inviteCode, false)
		// if err != nil {
		// 	panic(err)
		// }
		//
		// _, err = GetInvite(client, inviteCode, true)
		// if err != nil {
		// 	panic(err)
		// }
	}()

	// TestDeleteInvite
	func() {
		// TODO: invite codes....
		// will only return 404 without an invite code
		//
		// inviteCode := ""
		// _, err = DeleteInvite(client, inviteCode)
		// if err != nil {
		// 	panic(err)
		// }
	}()

	// -------------------
	// VOICES
	// -------------------

	// TestListVoiceRegions
	func() {
		list, err := c.GetVoiceRegions()
		if err != nil {
			panic(err)
		} else if len(list) == 0 {
			panic("expected at least one voice region")
		}
	}()
}
