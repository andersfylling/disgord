package disgord

import (
	"errors"
	"net/http"
	"strings"

	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
)

// endpoints
//
// https://discordapp.com/developers/docs/resources/emoji#emoji-resource
// Routes for controlling emojis do not follow the normal rate limit conventions.
// These routes are specifically limited on a per-guild basis to prevent abuse.
// This means that the quota returned by our APIs may be inaccurate,
// and you may encounter 429s.

// [REST] Returns a list of emoji objects for the given guild.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#list-guild-emojis
//  Reviewed                2018-06-10
//  Comment                 -
func ListGuildEmojis(client httd.Getter, id Snowflake) (ret []*Emoji, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.GuildEmojis(id),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// [REST] Returns an emoji object for the given guild and emoji IDs.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#get-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func GetGuildEmoji(client httd.Getter, guildID, emojiID Snowflake) (ret *Emoji, err error) {
	_, body, err := client.Get(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

func validEmojiName(name string) bool {
	// TODO: what is the allowed format?
	// a test showed that using "-" caused regex issues
	return !strings.Contains(name, "-")
}

// CreateGuildEmojiParams JSON params for func CreateGuildEmoji
type CreateGuildEmojiParams struct {
	Name  string      `json:"name"`
	Image string      `json:"image"`
	Roles []Snowflake `json:"roles"`
}

// [REST] Create a new emoji for the guild. Requires the 'MANAGE_EMOJIS' permission. Returns the new emoji object
// on success. Fires a Guild Emojis Update Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/emojis
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#create-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 Emojis and animated emojis have a maximum file size of 256kb. Attempting to upload
//                          an emoji larger than this limit will fail and return 400 Bad Request and an
//                          error message, but not a JSON status code.
func CreateGuildEmoji(client httd.Poster, guildID Snowflake, params *CreateGuildEmojiParams) (ret *Emoji, err error) {
	if !validEmojiName(params.Name) {
		err = errors.New("emoji name contains illegal characters. Did not send request")
		return
	}
	_, body, err := client.Post(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmojis(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyGuildEmojiParams JSON params for func ModifyGuildEmoji
type ModifyGuildEmojiParams struct {
	Name  string      `json:"name"`
	Roles []Snowflake `json:"roles"`
}

// [REST] Modify the given emoji. Requires the 'MANAGE_EMOJIS' permission. Returns the updated emoji object on success.
// Fires a Guild Emojis Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#modify-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func ModifyGuildEmoji(client httd.Patcher, guildID, emojiID Snowflake, params *ModifyGuildEmojiParams) (ret *Emoji, err error) {
	if !validEmojiName(params.Name) {
		err = errors.New("emoji name contains illegal characters. Did not send request")
		return
	}
	_, body, err := client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// [REST] Delete the given emoji. Requires the 'MANAGE_EMOJIS' permission. Returns 204 No Content on success.
// Fires a Guild Emojis Update Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/emojis/{emoji.id}
//  Rate limiter [MAJOR]    /guilds/{guild.id} // TODO: no idea if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/emoji#delete-guild-emoji
//  Reviewed                2018-06-10
//  Comment                 -
func DeleteGuildEmoji(client httd.Deleter, guildID, emojiID Snowflake) (err error) {
	resp, _, err := client.Delete(&httd.Request{
		Ratelimiter: ratelimitGuild(guildID),
		Endpoint:    endpoint.GuildEmoji(guildID, emojiID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}
