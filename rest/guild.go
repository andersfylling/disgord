package rest

import (
	"github.com/andersfylling/disgord/rest/httd"
	. "github.com/andersfylling/disgord/resource"
	"encoding/json"
	"github.com/andersfylling/snowflake"
	"net/http"
	"errors"
	"strconv"
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
	Name string `json:"name,omitempty"`
	Region string `json:"region,omitempty"`
	VerificationLvl int `json:"verification_level,omitempty"`
	DefaultMsgNotifications DefaultMessageNotificationLvl `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter ExplicitContentFilterLvl `json:"explicit_content_filter,omitempty"`
	AFKChannelID snowflake.ID `json:"afk_channel_id,omitempty"`
	AFKTimeout int `json:"afk_timeout,omitempty"`
	Icon string `json:"icon,omitempty"`
	OwnerID snowflake.ID `json:"owner_id,omitempty"`
	Splash string `json:"splash,omitempty"`
	SystemChannelID snowflake.ID `json:"system_channel_id,omitempty"`
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

// GetGuildChannels [GET]   Returns a list of guild channel objects.
// Endpoint                 /guilds/{guild.id}/channels
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-channels
// Reviewed                 2018-08-17
// Comment                  -
func GetGuildChannels(client httd.Getter, guildID snowflake.ID)  (ret *[]Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/channels",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}


// CreateGuildChannelParams https://discordapp.com/developers/docs/resources/guild#create-guild-channel-json-params
type CreateGuildChannelParams struct {
	Name                 string                `json:"name"`                            //  |
	Type                 uint                  `json:"type,omitempty"`                  // ?|
	Topic                string                `json:"topic,omitempty"`                 // ?|
	Bitrate              uint                  `json:"bitrate,omitempty"`               // ?|
	UserLimit            uint                  `json:"user_limit,omitempty"`            // ?|
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"` // ?|
	ParentID             snowflake.ID          `json:"parent_id,omitempty"`             // ?|
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
}

// CreateGuildChannel [POST]    Create a new channel object for the guild. Requires the 'MANAGE_CHANNELS' permission.
//                              Returns the new channel object on success. Fires a Channel Create Gateway event.
// Endpoint                     /guilds/{guild.id}/channels
// Rate limiter                 /guilds/{guild.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#create-guild-channel
// Reviewed                     2018-08-17
// Comment                      All parameters for this endpoint are optional excluding 'name'
func CreateGuildChannel(client httd.Poster, guildID snowflake.ID, params *CreateGuildChannelParams) (ret *Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/channels",
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}



// ModifyGuildChannelPositionsParams https://discordapp.com/developers/docs/resources/guild#modify-guild-channel-positions-json-params
type ModifyGuildChannelPositionsParams struct {
	ID snowflake.ID `json:"id"`
	Position int `json:"position"`
}

// ModifyGuildChannelPositions [PATCH]  Modify the positions of a set of channel objects for the guild. Requires
//                                      'MANAGE_CHANNELS' permission. Returns a 204 empty response on success.
//                                      Fires multiple Channel Update Gateway events.
// Endpoint                             /guilds/{guild.id}/channels
// Rate limiter                         /guilds/{guild.id}
// Discord documentation                https://discordapp.com/developers/docs/resources/guild#modify-guild-channel-positions
// Reviewed                             2018-08-17
// Comment                              Only channels to be modified are required, with the minimum being a swap
//                                      between at least two channels.
func ModifyGuildChannelPositions(client httd.Patcher, guildID snowflake.ID, params *ModifyGuildChannelPositionsParams) (ret *Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/channels",
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}

// GetGuildMember [GET]     Returns a guild member object for the specified user.
// Endpoint                 /guilds/{guild.id}/members/{user.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-member
// Reviewed                 2018-08-17
// Comment                  -
func GetGuildMember(client httd.Getter, guildID, userID snowflake.ID)  (ret *Member, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}



// GetGuildMembers [GET]    Returns a list of guild member objects that are members of the guild.
// Endpoint                 /guilds/{guild.id}/members
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-members
// Reviewed                 2018-08-17
// Comment                  All parameters to this endpoint are optional
// Comment#2                https://discordapp.com/developers/docs/resources/guild#list-guild-members-query-string-params
func GetGuildMembers(client httd.Getter, guildID, after snowflake.ID, limit int)  (ret []*Member, err error) {
	// TODO: convert after and limit to a query struct
	// omg i hate myself. use reflection to convert a query struct to string(?). it's at least better.
	query := ""
	if limit > 0 || !after.Empty() {
		query += "?"

		if !after.Empty() {
			query += "after=" + after.String()

			if limit > 0 {
				query += "&"
			}
		}

		if limit > 0 {
			query += "limit=" + strconv.Itoa(limit)
		}
	}

	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members" + query,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}