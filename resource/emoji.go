package resource

import (
	"github.com/andersfylling/disgord/request"
	"github.com/andersfylling/snowflake"
)

// Emoji
type Emoji struct {
	ID            snowflake.ID   `json:"id"`
	Name          string         `json:"name"`
	Roles         []snowflake.ID `json:"roles,omitempty"`
	User          *User          `json:"user,omitempty"` // the user who created the emoji
	RequireColons bool           `json:"require_colons,omitempty"`
	Managed       bool           `json:"managed,omitempty"`
	Animated      bool           `json:"animated,omitempty"`
}
type PartialEmoji = Emoji

// Mention
// TODO: review
func (e *Emoji) Mention() string {
	return "<" + e.Name + ":" + e.ID.String() + ">"
}

// MentionAnimated add the animation prefix if a animated emoji
// TODO: review
func (e *Emoji) MentionAnimated() string {
	prefix := ""
	if e.Animated {
		prefix = "a:"
	}

	return "<" + prefix + e.Name + ":" + e.ID.String() + ">"
}

func (e *Emoji) Clear() {
	// obviously don't delete the user ...
}

// endpoints
//
// https://discordapp.com/developers/docs/resources/emoji#emoji-resource
// Routes for controlling emojis do not follow the normal rate limit conventions.
// These routes are specifically limited on a per-guild basis to prevent abuse.
// This means that the quota returned by our APIs may be inaccurate,
// and you may encounter 429s.

// ReqListGuildEmojis [GET] Returns a list of emoji objects for the given guild.
// Endpoint                 /guilds/{guild.id}/emojis
// Rate limiter [MAJOR]     /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/emoji#list-guild-emojis
// Reviewed                 2018-03-17
// Comment                  -
func ReqListGuildEmojis(requester request.DiscordGetter, guildID snowflake.ID) (emoji *Emoji, err error) {
	rateLimitKey := "/guilds/" + guildID.String()
	path := rateLimitKey + "/emojis"

	_, err = requester.Get(rateLimitKey, path, emoji)

	return emoji, err
}

// ReqGetGuildEmoji [GET] Returns an emoji object for the given guild and emoji IDs.
// Endpoint               /guilds/{guild.id}/emojis/{emoji.id}
// Rate limiter [MAJOR]   /guilds/{guild.id}
// Discord documentation  https://discordapp.com/developers/docs/resources/emoji#get-guild-emoji
// Reviewed               2018-03-17
// Comment                -
func ReqGetGuildEmoji(requester request.DiscordGetter, guildID, emojiID snowflake.ID) (emoji *Emoji, err error) {
	rateLimitKey := "/guilds/" + guildID.String()
	path := rateLimitKey + "/emojis/" + emojiID.String()

	_, err = requester.Get(rateLimitKey, path, emoji)

	return emoji, err
}

// ReqCreateGuildEmoji [POST] Create a new emoji for the guild. Requires the
//                            'MANAGE_EMOJIS' permission. Returns the new emoji
//                            object on success. Fires a Guild Emojis Update Gateway event.
// Endpoint                   /guilds/{guild.id}/emojis
// Rate limiter [MAJOR]       /guilds/{guild.id}
// Discord documentation      https://discordapp.com/developers/docs/resources/emoji#create-guild-emoji
// Reviewed                   2018-03-17
// Comment                    "Emojis and animated emojis have a maximum file size of 256kb.
//                            Attempting to upload an emoji larger than this limit will fail
//                            and return 400 Bad Request and an error message, but not a JSON
//                            status code." - Discord docs
func ReqCreateGuildEmoji(requester request.DiscordPoster, guildID snowflake.ID) (emoji *Emoji, err error) {
	rateLimitKey := "/guilds/" + guildID.String()
	path := rateLimitKey + "/emojis"

	_, err = requester.Post(rateLimitKey, path, emoji, nil)

	return emoji, err
}

// ReqModifyGuildEmoji [PATCH]  Modify the given emoji. Requires the 'MANAGE_EMOJIS'
//                              permission. Returns the updated emoji object on success.
//                              Fires a Guild Emojis Update Gateway event.
// Endpoint                     /guilds/{guild.id}/emojis/{emoji.id}
// Rate limiter [MAJOR]         /guilds/{guild.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/emoji#modify-guild-emoji
// Reviewed                     2018-03-17
// Comment                      -
func ReqModifyGuildEmoji(requester request.DiscordPatcher, guildID, emojiID snowflake.ID) (emoji *Emoji, err error) {
	rateLimitKey := "/guilds/" + guildID.String()
	path := rateLimitKey + "/emojis/" + emojiID.String()

	_, err = requester.Patch(rateLimitKey, path, emoji, nil)

	return emoji, err
}

// ReqDeleteGuildEmoji [DELETE] Delete the given emoji. Requires the
//                              'MANAGE_EMOJIS' permission. Returns 204
//                              No Content on success. Fires a Guild Emojis
//                              Update Gateway event.
// Endpoint                     /guilds/{guild.id}/emojis/{emoji.id}
// Rate limiter [MAJOR]         /guilds/{guild.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/emoji#delete-guild-emoji
// Reviewed                     2018-03-17
// Comment                      -
func ReqDeleteGuildEmoji(requester request.DiscordDeleter, guildID, emojiID snowflake.ID) (err error) {
	rateLimitKey := "/guilds/" + guildID.String()
	path := rateLimitKey + "/emojis/" + emojiID.String()

	_, err = requester.Delete(rateLimitKey, path)

	return err
}
