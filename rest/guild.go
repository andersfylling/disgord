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
// Comment#2                "List Guild Members"
// Comment#3                https://discordapp.com/developers/docs/resources/guild#list-guild-members-query-string-params
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

// AddGuildMemberParams https://discordapp.com/developers/docs/resources/guild#add-guild-member-json-params
type AddGuildMemberParams struct {
	AccessToken string `json:"access_token"`
	Nick string `json:"nick,omitempty"`
	Roles []snowflake.ID `json:"roles"`
	Mute bool `json:"mute"`
	Deaf bool `json:"deaf"`
}

// AddGuildMember [PUT]     Adds a user to the guild, provided you have a valid oauth2 access token for the user
//                          with the guilds.join scope. Returns a 201 Created with the guild member as the body,
//                          or 204 No Content if the user is already a member of the guild. Fires a Guild Member Add
//                          Gateway event. Requires the bot to have the CREATE_INSTANT_INVITE permission.
// Endpoint                 /guilds/{guild.id}/members/{user.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#add-guild-member
// Reviewed                 2018-08-18
// Comment                  All parameters to this endpoint except for access_token are optional.
func AddGuildMember(client httd.Puter, guildID, userID snowflake.ID, params *AddGuildMemberParams) (ret *Member, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String(),
		JSONParams:  params,
	}
	resp, body, err := client.Put(details)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusCreated {
		err = json.Unmarshal(body, &ret)
		return
	}

	if resp.StatusCode == http.StatusNoContent {
		msg := "User{id:" + userID.String() + "} already exists in guild{id:" + guildID.String() + "}"
		err = errors.New(msg)
	}

	return
}


// ModifyGuildMemberParams https://discordapp.com/developers/docs/resources/guild#modify-guild-member-json-params
type ModifyGuildMemberParams struct {
	Nick string `json:"nick,omitempty"` // :MANAGE_NICKNAMES
	Roles []snowflake.ID `json:"roles,omitempty"` // :MANAGE_ROLES
	Mute bool `json:"mute,omitempty"`             // :MUTE_MEMBERS
	Deaf bool `json:"deaf,omitempty"`             // :DEAFEN_MEMBERS
	ChannelID snowflake.ID `json:"channel_id,omitempty"` // :MOVE_MEMBERS
}

// ModifyGuildMember [PATCH]    Modify attributes of a guild member. Returns a 204 empty response on success.
//                              Fires a Guild Member Update Gateway event.
// Endpoint                     /guilds/{guild.id}/members/{user.id}
// Rate limiter                 /guilds/{guild.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#modify-guild-member
// Reviewed                     2018-08-17
// Comment                      All parameters to this endpoint are optional. When moving members to channels,
//                              the API user must have permissions to both connect to the channel and have the
//                              MOVE_MEMBERS permission.
func ModifyGuildMember(client httd.Patcher, guildID, userID snowflake.ID, params *ModifyGuildMemberParams) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String(),
		JSONParams:  params,
	}
	resp, _, err := client.Patch(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not change attributes of member. Does the member exist, and do you have permissions?"
		err = errors.New(msg)
	}

	return
}


// ModifyGuildMemberParams https://discordapp.com/developers/docs/resources/guild#modify-guild-member-json-params
type ModifyCurrentUserNickParams struct {
	Nick string `json:"nick"` // :CHANGE_NICKNAME
}

// ModifyCurrentUserNick [PATCH]    Modifies the nickname of the current user in a guild. Returns a 200 with the
//                                  nickname on success. Fires a Guild Member Update Gateway event.
// Endpoint                         /guilds/{guild.id}/members/@me/nick
// Rate limiter                     /guilds/{guild.id}
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#modify-current-user-nick
// Reviewed                         2018-08-18
// Comment                          -
func ModifyCurrentUserNick(client httd.Patcher, guildID snowflake.ID, params *ModifyCurrentUserNickParams) (nickname string, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/@me/nick",
		JSONParams:  params,
	}
	resp, body, err := client.Patch(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		msg := "Could not change nickname. Do you have the CHANGE_NICKNAME permission?"
		err = errors.New(msg)
		return
	}

	err = json.Unmarshal(body, nickname)
	return
}


// AddGuildMemberRole [PUT] Adds a role to a guild member. Requires the 'MANAGE_ROLES' permission. Returns a 204
//                          empty response on success. Fires a Guild Member Update Gateway event.
// Endpoint                 /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#add-guild-member-role
// Reviewed                 2018-08-18
// Comment                  -
func AddGuildMemberRole(client httd.Puter, guildID, userID, roleID snowflake.ID) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String(),
	}
	resp, _, err := client.Put(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not add role to user. Do you have the MANAGE_ROLES permission?"
		err = errors.New(msg)
	}

	return
}

// RemoveGuildMemberRole [DELETE]   Removes a role from a guild member. Requires the 'MANAGE_ROLES' permission. Returns
//                                  a 204 empty response on success. Fires a Guild Member Update Gateway event.
// Endpoint                         /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// Rate limiter                     /guilds/{guild.id}
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#remove-guild-member-role
// Reviewed                         2018-08-18
// Comment                          -
func RemoveGuildMemberRole(client httd.Deleter, guildID, userID, roleID snowflake.ID) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove role from user. Do you have the MANAGE_ROLES permission?"
		err = errors.New(msg)
	}

	return
}

// RemoveGuildMember [DELETE]   Remove a member from a guild. Requires 'KICK_MEMBERS' permission. Returns a 204
//                              empty response on success. Fires a Guild Member Remove Gateway event.
// Endpoint                     /guilds/{guild.id}/members/{user.id}
// Rate limiter                 /guilds/{guild.id}
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#remove-guild-member
// Reviewed                     2018-08-18
// Comment                      -
func RemoveGuildMember(client httd.Deleter, guildID, userID snowflake.ID) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove user from guild. Do you have the KICK_MEMBERS permission?"
		err = errors.New(msg)
	}

	return
}



// GetGuildBans [GET]       Returns a list of ban objects for the users banned from this guild.
//                          Requires the 'BAN_MEMBERS' permission.
// Endpoint                 /guilds/{guild.id}/bans
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-bans
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildBans(client httd.Getter, guildID snowflake.ID) (ret []*Ban, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/bans",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}

// GetGuildBan [GET]        Returns a ban object for the given user or a 404 not found if the ban cannot be found.
//                          Requires the 'BAN_MEMBERS' permission.
// Endpoint                 /guilds/{guild.id}/bans/{user.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-ban
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildBan(client httd.Getter, guildID, userID snowflake.ID) (ret *Ban, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/bans/" + userID.String(),
	}
	resp, body, err := client.Get(details)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		msg := "given user is not registered as banned"
		err = errors.New(msg)
	}

	err = json.Unmarshal(body, &ret)
	return
}

// CreateGuildBanParams https://discordapp.com/developers/docs/resources/guild#create-guild-ban-query-string-params
type CreateGuildBanParams struct {
	DeleteMessageDays int `urlparam:"delete_message_days"` // number of days to delete messages for (0-7)
	Reason string `urlparam:"reason"` // reason for being banned
}

// getQueryString this ins't really pretty, but it works.
func (params *CreateGuildBanParams) getQueryString() string {
	seperator := "?"
	query := ""

	if params.DeleteMessageDays > 0 {
		query += seperator + "delete_message_days=" + strconv.Itoa(params.DeleteMessageDays)
		seperator = "&"
	}

	if params.Reason != "" {
		query += seperator + "reason=" + params.Reason
	}

	return query
}

// CreateGuildBan [PUT]     Create a guild ban, and optionally delete previous messages sent by the banned user.
//                          Requires the 'BAN_MEMBERS' permission. Returns a 204 empty response on success.
//                          Fires a Guild Ban Add Gateway event.
// Endpoint                 /guilds/{guild.id}/bans/{user.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#create-guild-ban
// Reviewed                 2018-08-18
// Comment                  -
func CreateGuildBan(client httd.Puter, guildID, userID snowflake.ID, params *CreateGuildBanParams) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/bans/" + userID.String() + params.getQueryString(),
	}
	resp, _, err := client.Put(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "could not ban member"
		err = errors.New(msg)
	}

	return
}

// RemoveGuildBan [DELETE]  Remove the ban for a user. Requires the 'BAN_MEMBERS' permissions.
//                          Returns a 204 empty response on success. Fires a Guild Ban Remove Gateway event.
// Endpoint                 /guilds/{guild.id}/bans/{user.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#remove-guild-ban
// Reviewed                 2018-08-18
// Comment                  -
func RemoveGuildBan(client httd.Deleter, guildID, userID snowflake.ID) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/bans/" + userID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove ban on user. Do you have the BAN_MEMBERS permission?"
		err = errors.New(msg)
	}

	return
}

// GetGuildRoles [GET]      Returns a list of role objects for the guild.
// Endpoint                 /guilds/{guild.id}/roles
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-roles
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildRoles(client httd.Getter, guildID snowflake.ID) (ret []*Role, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}


// CreateGuildRoleParams https://discordapp.com/developers/docs/resources/guild#create-guild-role-json-params
type CreateGuildRoleParams struct {
	Name string `json:"name,omitempty"`
	Permissions int `json:"permissions,omitempty"`
	Color int `json:"color,omitempty"`
	Hoist bool `json:"hoist,omitempty"`
	Mentionable bool `json:"mentionable,omitempty"`
}

// CreateGuildRole [POST]   Create a new role for the guild. Requires the 'MANAGE_ROLES' permission.
//                          Returns the new role object on success. Fires a Guild Role Create Gateway event.
// Endpoint                 /guilds/{guild.id}/roles
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#create-guild-role
// Reviewed                 2018-08-18
// Comment                  All JSON params are optional.
func CreateGuildRole(client httd.Poster, guildID snowflake.ID, params *CreateGuildRoleParams) (ret *Role, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles",
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ret)
	return
}