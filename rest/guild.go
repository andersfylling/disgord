package rest

import (
	"github.com/andersfylling/disgord/rest/httd"
	. "github.com/andersfylling/disgord/resource"
	"encoding/json"
	"github.com/andersfylling/snowflake"
	"net/http"
	"errors"
)

// CreateGuildParams https://discordapp.com/developers/docs/resources/guild#create-guild-json-params
// example partial channel object:
// {
//    "name": "naming-things-is-hard",
//    "type": 0
// }
type CreateGuildParams struct {
	Name string `json:"name"`
	Region string `json:"region"`
	Icon string `json:"icon"`
	VerificationLvl int `json:"verification_level"`
	DefaultMsgNotifications DefaultMessageNotificationLvl `json:"default_message_notifications"`
	ExplicitContentFilter ExplicitContentFilterLvl `json:"explicit_content_filter"`
	Roles []*Role `json:"roles"`
	Channels []*PartialChannel `json:"channels"`
}

// CreateGuild [POST]       Create a new guild. Returns a guild object on success. Fires a Guild Create Gateway event.
// Endpoint                 /guilds
// Rate limiter             /guilds
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#create-guild
// Reviewed                 2018-08-16
// Comment                  This endpoint can be used only by bots in less than 10 guilds. Creating channel
//                          categories from this endpoint is not supported.
func CreateGuild(client httd.Poster, params *CreateGuildParams) (ret *Guild, err error) {
	// TODO: check if bot
	// TODO-2: is bot in less than 10 guilds?
	details := &httd.Request{
		Ratelimiter: "/guilds",
		Endpoint:    "/guilds",
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}

// GetGuild [GET]           Returns the guild object for the given id.
// Endpoint                 /guilds/{guild.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild
// Reviewed                 2018-08-17
// Comment                  -
func GetGuild(client httd.Getter, guildID snowflake.ID) (ret *Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}


// ModifyGuildParams https://discordapp.com/developers/docs/resources/guild#modify-guild-json-params
type ModifyGuildParams struct {
	Name string `json:"name"`
	Region string `json:"region"`
	VerificationLvl int `json:"verification_level"`
	DefaultMsgNotifications DefaultMessageNotificationLvl `json:"default_message_notifications"`
	ExplicitContentFilter ExplicitContentFilterLvl `json:"explicit_content_filter"`
	AFKChannelID snowflake.ID `json:"afk_channel_id"`
	AFKTimeout int `json:"afk_timeout"`
	Icon string `json:"icon"`
	OwnerID snowflake.ID `json:"owner_id"`
	Splash string `json:"splash"`
	SystemChannelID snowflake.ID `json:"system_channel_id"`
}

// ModifyGuild [PATCH]      Modify a guild's settings. Requires the 'MANAGE_GUILD' permission. Returns the updated
//                          guild object on success. Fires a Guild Update Gateway event.
// Endpoint                 /guilds/{guild.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#modify-guild
// Reviewed                 2018-08-17
// Comment                  All parameters to this endpoint are optional
func ModifyGuild(client httd.Patcher, guildID snowflake.ID, params *ModifyGuildParams) (ret *Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String(),
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}

// DeleteGuild [DELETE]     Delete a guild permanently. User must be owner. Returns 204 No Content on success.
//                          Fires a Guild Delete Gateway event.
// Endpoint                 /guilds/{guild.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#delete-guild
// Reviewed                 2018-08-17
// Comment                  -
func DeleteGuild(client httd.Deleter, guildID snowflake.ID) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "unexpected http response code. Got " + resp.Status + ", wants " + http.StatusText(http.StatusNoContent)
		err = errors.New(msg)
	}
	return
}