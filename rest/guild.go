package rest

import (
	"errors"
	"net/http"
	"strconv"

	. "github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/rest/httd"
	. "github.com/andersfylling/disgord/snowflake"
)

// CreateGuildParams https://discordapp.com/developers/docs/resources/guild#create-guild-json-params
// example partial channel object:
// {
//    "name": "naming-things-is-hard",
//    "type": 0
// }
type CreateGuildParams struct {
	Name                    string                        `json:"name"`
	Region                  string                        `json:"region"`
	Icon                    string                        `json:"icon"`
	VerificationLvl         int                           `json:"verification_level"`
	DefaultMsgNotifications DefaultMessageNotificationLvl `json:"default_message_notifications"`
	ExplicitContentFilter   ExplicitContentFilterLvl      `json:"explicit_content_filter"`
	Roles                   []*Role                       `json:"roles"`
	Channels                []*PartialChannel             `json:"channels"`
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

	err = unmarshal(body, &ret)
	return
}

// GetGuild [GET]           Returns the guild object for the given id.
// Endpoint                 /guilds/{guild.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild
// Reviewed                 2018-08-17
// Comment                  -
func GetGuild(client httd.Getter, guildID Snowflake) (ret *Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyGuildParams https://discordapp.com/developers/docs/resources/guild#modify-guild-json-params
type ModifyGuildParams struct {
	Name                    string                        `json:"name,omitempty"`
	Region                  string                        `json:"region,omitempty"`
	VerificationLvl         int                           `json:"verification_level,omitempty"`
	DefaultMsgNotifications DefaultMessageNotificationLvl `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter   ExplicitContentFilterLvl      `json:"explicit_content_filter,omitempty"`
	AFKChannelID            Snowflake                  `json:"afk_channel_id,omitempty"`
	AFKTimeout              int                           `json:"afk_timeout,omitempty"`
	Icon                    string                        `json:"icon,omitempty"`
	OwnerID                 Snowflake                  `json:"owner_id,omitempty"`
	Splash                  string                        `json:"splash,omitempty"`
	SystemChannelID         Snowflake                  `json:"system_channel_id,omitempty"`
}

// ModifyGuild [PATCH]      Modify a guild's settings. Requires the 'MANAGE_GUILD' permission. Returns the updated
//                          guild object on success. Fires a Guild Update Gateway event.
// Endpoint                 /guilds/{guild.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#modify-guild
// Reviewed                 2018-08-17
// Comment                  All parameters to this endpoint are optional
func ModifyGuild(client httd.Patcher, guildID Snowflake, params *ModifyGuildParams) (ret *Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String(),
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// DeleteGuild [DELETE]     Delete a guild permanently. User must be owner. Returns 204 No Content on success.
//                          Fires a Guild Delete Gateway event.
// Endpoint                 /guilds/{guild.id}
// Rate limiter             /guilds/{guild.id}
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#delete-guild
// Reviewed                 2018-08-17
// Comment                  -
func DeleteGuild(client httd.Deleter, guildID Snowflake) (err error) {
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
// Rate limiter             /guilds/{guild.id}/channels
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-channels
// Reviewed                 2018-08-17
// Comment                  -
func GetGuildChannels(client httd.Getter, guildID Snowflake) (ret *[]Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildChannels(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/channels",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
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
	ParentID             Snowflake          `json:"parent_id,omitempty"`             // ?|
	NSFW                 bool                  `json:"nsfw,omitempty"`                  // ?|
}

// CreateGuildChannel [POST]    Create a new channel object for the guild. Requires the 'MANAGE_CHANNELS' permission.
//                              Returns the new channel object on success. Fires a Channel Create Gateway event.
// Endpoint                     /guilds/{guild.id}/channels
// Rate limiter                 /guilds/{guild.id}/channels
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#create-guild-channel
// Reviewed                     2018-08-17
// Comment                      All parameters for this endpoint are optional excluding 'name'
func CreateGuildChannel(client httd.Poster, guildID Snowflake, params *CreateGuildChannelParams) (ret *Channel, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuild(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/channels",
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyGuildChannelPositionsParams https://discordapp.com/developers/docs/resources/guild#modify-guild-channel-positions-json-params
type ModifyGuildChannelPositionsParams struct {
	ID       Snowflake `json:"id"`
	Position int          `json:"position"`
}

// ModifyGuildChannelPositions [PATCH]  Modify the positions of a set of channel objects for the guild. Requires
//                                      'MANAGE_CHANNELS' permission. Returns a 204 empty response on success.
//                                      Fires multiple Channel Update Gateway events.
// Endpoint                             /guilds/{guild.id}/channels
// Rate limiter                         /guilds/{guild.id}/channels
// Discord documentation                https://discordapp.com/developers/docs/resources/guild#modify-guild-channel-positions
// Reviewed                             2018-08-17
// Comment                              Only channels to be modified are required, with the minimum being a swap
//                                      between at least two channels.
func ModifyGuildChannelPositions(client httd.Patcher, guildID Snowflake, params *ModifyGuildChannelPositionsParams) (ret *Guild, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildChannels(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/channels",
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildMember [GET]     Returns a guild member object for the specified user.
// Endpoint                 /guilds/{guild.id}/members/{user.id}
// Rate limiter             /guilds/{guild.id}/members
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-member
// Reviewed                 2018-08-17
// Comment                  -
func GetGuildMember(client httd.Getter, guildID, userID Snowflake) (ret *Member, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildMembers [GET]    Returns a list of guild member objects that are members of the guild.
// Endpoint                 /guilds/{guild.id}/members
// Rate limiter             /guilds/{guild.id}/members
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-members
// Reviewed                 2018-08-17
// Comment                  All parameters to this endpoint are optional
// Comment#2                "List Guild Members"
// Comment#3                https://discordapp.com/developers/docs/resources/guild#list-guild-members-query-string-params
func GetGuildMembers(client httd.Getter, guildID, after Snowflake, limit int) (ret []*Member, err error) {
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
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members" + query,
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// AddGuildMemberParams https://discordapp.com/developers/docs/resources/guild#add-guild-member-json-params
type AddGuildMemberParams struct {
	AccessToken string         `json:"access_token"`
	Nick        string         `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles"`
	Mute        bool           `json:"mute"`
	Deaf        bool           `json:"deaf"`
}

// AddGuildMember [PUT]     Adds a user to the guild, provided you have a valid oauth2 access token for the user
//                          with the guilds.join scope. Returns a 201 Created with the guild member as the body,
//                          or 204 No Content if the user is already a member of the guild. Fires a Guild Member Add
//                          Gateway event. Requires the bot to have the CREATE_INSTANT_INVITE permission.
// Endpoint                 /guilds/{guild.id}/members/{user.id}
// Rate limiter             /guilds/{guild.id}/members
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#add-guild-member
// Reviewed                 2018-08-18
// Comment                  All parameters to this endpoint except for access_token are optional.
func AddGuildMember(client httd.Puter, guildID, userID Snowflake, params *AddGuildMemberParams) (ret *Member, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/members/" + userID.String(),
		JSONParams:  params,
	}
	resp, body, err := client.Put(details)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusCreated {
		err = unmarshal(body, &ret)
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
	Nick      string         `json:"nick,omitempty"`       // :MANAGE_NICKNAMES
	Roles     []Snowflake `json:"roles,omitempty"`      // :MANAGE_ROLES
	Mute      bool           `json:"mute,omitempty"`       // :MUTE_MEMBERS
	Deaf      bool           `json:"deaf,omitempty"`       // :DEAFEN_MEMBERS
	ChannelID Snowflake   `json:"channel_id,omitempty"` // :MOVE_MEMBERS
}

// ModifyGuildMember [PATCH]    Modify attributes of a guild member. Returns a 204 empty response on success.
//                              Fires a Guild Member Update Gateway event.
// Endpoint                     /guilds/{guild.id}/members/{user.id}
// Rate limiter                 /guilds/{guild.id}/members
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#modify-guild-member
// Reviewed                     2018-08-17
// Comment                      All parameters to this endpoint are optional. When moving members to channels,
//                              the API user must have permissions to both connect to the channel and have the
//                              MOVE_MEMBERS permission.
func ModifyGuildMember(client httd.Patcher, guildID, userID Snowflake, params *ModifyGuildMemberParams) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
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
// Rate limiter                     /guilds/{guild.id}/members TODO: I don't know if this is correct
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#modify-current-user-nick
// Reviewed                         2018-08-18
// Comment                          -
func ModifyCurrentUserNick(client httd.Patcher, guildID Snowflake, params *ModifyCurrentUserNickParams) (nickname string, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
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

	err = unmarshal(body, nickname)
	return
}

// AddGuildMemberRole [PUT] Adds a role to a guild member. Requires the 'MANAGE_ROLES' permission. Returns a 204
//                          empty response on success. Fires a Guild Member Update Gateway event.
// Endpoint                 /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// Rate limiter             /guilds/{guild.id}/members TODO: I don't know if this is correct
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#add-guild-member-role
// Reviewed                 2018-08-18
// Comment                  -
func AddGuildMemberRole(client httd.Puter, guildID, userID, roleID Snowflake) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
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
// Rate limiter                     /guilds/{guild.id}/members
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#remove-guild-member-role
// Reviewed                         2018-08-18
// Comment                          -
func RemoveGuildMemberRole(client httd.Deleter, guildID, userID, roleID Snowflake) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
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
// Rate limiter                 /guilds/{guild.id}/members
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#remove-guild-member
// Reviewed                     2018-08-18
// Comment                      -
func RemoveGuildMember(client httd.Deleter, guildID, userID Snowflake) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildMembers(guildID),
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
// Rate limiter             /guilds/{guild.id}/bans
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-bans
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildBans(client httd.Getter, guildID Snowflake) (ret []*Ban, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildBans(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/bans",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildBan [GET]        Returns a ban object for the given user or a 404 not found if the ban cannot be found.
//                          Requires the 'BAN_MEMBERS' permission.
// Endpoint                 /guilds/{guild.id}/bans/{user.id}
// Rate limiter             /guilds/{guild.id}/bans
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-ban
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildBan(client httd.Getter, guildID, userID Snowflake) (ret *Ban, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildBans(guildID),
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

	err = unmarshal(body, &ret)
	return
}

// CreateGuildBanParams https://discordapp.com/developers/docs/resources/guild#create-guild-ban-query-string-params
type CreateGuildBanParams struct {
	DeleteMessageDays int    `urlparam:"delete_message_days"` // number of days to delete messages for (0-7)
	Reason            string `urlparam:"reason"`              // reason for being banned
}

// getQueryString this ins't really pretty, but it works.
func (params *CreateGuildBanParams) getQueryString() string {
	separator := "?"
	query := ""

	if params.DeleteMessageDays > 0 {
		query += separator + "delete_message_days=" + strconv.Itoa(params.DeleteMessageDays)
		separator = "&"
	}

	if params.Reason != "" {
		query += separator + "reason=" + params.Reason
	}

	return query
}

// CreateGuildBan [PUT]     Create a guild ban, and optionally delete previous messages sent by the banned user.
//                          Requires the 'BAN_MEMBERS' permission. Returns a 204 empty response on success.
//                          Fires a Guild Ban Add Gateway event.
// Endpoint                 /guilds/{guild.id}/bans/{user.id}
// Rate limiter             /guilds/{guild.id}/bans
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#create-guild-ban
// Reviewed                 2018-08-18
// Comment                  -
func CreateGuildBan(client httd.Puter, guildID, userID Snowflake, params *CreateGuildBanParams) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildBans(guildID),
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
// Rate limiter             /guilds/{guild.id}/bans
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#remove-guild-ban
// Reviewed                 2018-08-18
// Comment                  -
func RemoveGuildBan(client httd.Deleter, guildID, userID Snowflake) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildBans(guildID),
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
// Rate limiter             /guilds/{guild.id}/roles
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-roles
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildRoles(client httd.Getter, guildID Snowflake) (ret []*Role, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildRoles(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// CreateGuildRoleParams https://discordapp.com/developers/docs/resources/guild#create-guild-role-json-params
type CreateGuildRoleParams struct {
	Name        string `json:"name,omitempty"`
	Permissions int    `json:"permissions,omitempty"`
	Color       int    `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`
}

// CreateGuildRole [POST]   Create a new role for the guild. Requires the 'MANAGE_ROLES' permission.
//                          Returns the new role object on success. Fires a Guild Role Create Gateway event.
// Endpoint                 /guilds/{guild.id}/roles
// Rate limiter             /guilds/{guild.id}/roles
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#create-guild-role
// Reviewed                 2018-08-18
// Comment                  All JSON params are optional.
func CreateGuildRole(client httd.Poster, guildID Snowflake, params *CreateGuildRoleParams) (ret *Role, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildRoles(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles",
		JSONParams:  params,
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyGuildRolePositionsParams https://discordapp.com/developers/docs/resources/guild#modify-guild-role-positions-json-params
type ModifyGuildRolePositionsParams struct {
	ID       Snowflake `json:"id"`
	Position int          `json:"position"`
}

// ModifyGuildRolePositions [PATCH] Modify the positions of a set of role objects for the guild. Requires the
//                                  'MANAGE_ROLES' permission. Returns a list of all of the guild's role objects
//                                  on success. Fires multiple Guild Role Update Gateway events.
// Endpoint                         /guilds/{guild.id}/roles
// Rate limiter                     /guilds/{guild.id}/roles
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#modify-guild-role-positions
// Reviewed                         2018-08-18
// Comment                          -
func ModifyGuildRolePositions(client httd.Patcher, guildID Snowflake, params *ModifyGuildRolePositionsParams) (ret []*Role, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildRoles(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles",
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

type ModifyGuildRoleParams struct {
	Name        string `json:"name,omitempty"`
	Permissions int    `json:"permissions,omitempty"`
	Color       int    `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`
}

// ModifyGuildRole [PATCH]  Modify a guild role. Requires the 'MANAGE_ROLES' permission. Returns the updated
//                          role on success. Fires a Guild Role Update Gateway event.
// Endpoint                 /guilds/{guild.id}/roles/{role.id}
// Rate limiter             /guilds/{guild.id}/roles
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#modify-guild-role
// Reviewed                 2018-08-18
// Comment                  -
func ModifyGuildRole(client httd.Patcher, guildID, roleID Snowflake, params *ModifyGuildRoleParams) (ret []*Role, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildRoles(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/roles/" + roleID.String(),
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// DeleteGuildRole [DELETE] Delete a guild role. Requires the 'MANAGE_ROLES' permission. Returns a 204 empty
//                          response on success. Fires a Guild Role Delete Gateway event.
// Endpoint                 /guilds/{guild.id}/roles/{role.id}
// Rate limiter             /guilds/{guild.id}/roles
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#delete-guild-role
// Reviewed                 2018-08-18
// Comment                  -
func DeleteGuildRole(client httd.Deleter, guildID, integrationID Snowflake) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildRoles(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/integrations/" + integrationID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove role. Do you have the MANAGE_ROLES permission?"
		err = errors.New(msg)
	}

	return
}

// GetGuildPruneCountParams https://discordapp.com/developers/docs/resources/guild#get-guild-prune-count-query-string-params
type GuildPruneParams struct {
	Days int `urlparam:"days"` // number of days to count prune for (1 or more)
}

// getQueryString this ins't really pretty, but it works.
func (params *GuildPruneParams) getQueryString() string {
	separator := "?"
	query := ""

	if params.Days > 0 {
		query += separator + "days=" + strconv.Itoa(params.Days)
	}

	return query
}

// GetGuildPruneCount [GET] Returns an object with one 'pruned' key indicating the number of members that would
//                          be removed in a prune operation. Requires the 'KICK_MEMBERS' permission.
// Endpoint                 /guilds/{guild.id}/prune
// Rate limiter             /guilds/{guild.id}/prune
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-prune-count
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildPruneCount(client httd.Getter, guildID Snowflake, params *GuildPruneParams) (ret *GuildPruneCount, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildPrune(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/prune" + params.getQueryString(),
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// BeginGuildPrune [POST]   Begin a prune operation. Requires the 'KICK_MEMBERS' permission. Returns an object
//                          with one 'pruned' key indicating the number of members that were removed in the
//                          prune operation. Fires multiple Guild Member Remove Gateway events.
// Endpoint                 /guilds/{guild.id}/prune
// Rate limiter             /guilds/{guild.id}/prune
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#begin-guild-prune
// Reviewed                 2018-08-18
// Comment                  -
func BeginGuildPrune(client httd.Poster, guildID Snowflake, params *GuildPruneParams) (ret *GuildPruneCount, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildPrune(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/prune" + params.getQueryString(),
	}
	_, body, err := client.Post(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildVoiceRegions [GET]   Returns a list of voice region objects for the guild. Unlike the similar /voice
//                              route, this returns VIP servers when the guild is VIP-enabled.
// Endpoint                     /guilds/{guild.id}/regions
// Rate limiter                 /guilds/{guild.id}/regions
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#get-guild-voice-regions
// Reviewed                     2018-08-18
// Comment                      -
func GetGuildVoiceRegions(client httd.Getter, guildID Snowflake) (ret []*VoiceRegion, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildRegions(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/regions",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildInvites [GET]    Returns a list of invite objects (with invite metadata) for the guild.
//                          Requires the 'MANAGE_GUILD' permission.
// Endpoint                 /guilds/{guild.id}/invites
// Rate limiter             /guilds/{guild.id}/invites
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-invites
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildInvites(client httd.Getter, guildID Snowflake) (ret []*Invite, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildInvites(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/invites",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildIntegrations [GET]   Returns a list of integration objects for the guild.
//                              Requires the 'MANAGE_GUILD' permission.
// Endpoint                     /guilds/{guild.id}/integrations
// Rate limiter                 /guilds/{guild.id}/integrations
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#get-guild-integrations
// Reviewed                     2018-08-18
// Comment                      -
func GetGuildIntegrations(client httd.Getter, guildID Snowflake) (ret []*Integration, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildIntegrations(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/integrations",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// CreateGuildIntegrationParams https://discordapp.com/developers/docs/resources/guild#create-guild-integration-json-params
type CreateGuildIntegrationParams struct {
	Type string       `json:"type"`
	ID   Snowflake `json:"id"`
}

// CreateGuildIntegration [POST]    Attach an integration object from the current user to the guild. Requires
//                                  the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
//                                  Fires a Guild Integrations Update Gateway event.
// Endpoint                         /guilds/{guild.id}/integrations
// Rate limiter                     /guilds/{guild.id}/integrations
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#create-guild-integration
// Reviewed                         2018-08-18
// Comment                          -
func CreateGuildIntegration(client httd.Poster, guildID Snowflake, params *CreateGuildIntegrationParams) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildIntegrations(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/integrations",
		JSONParams:  params,
	}
	resp, _, err := client.Post(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not create the integration object. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}

	return
}

// ModifyGuildIntegrationParams https://discordapp.com/developers/docs/resources/guild#modify-guild-integration-json-params
type ModifyGuildIntegrationParams struct {
	ExpireBehavior    int  `json:"expire_behavior"`
	ExpireGracePeriod int  `json:"expire_grace_period"`
	EnableEmoticons   bool `json:"enable_emoticons"`
}

// ModifyGuildIntegration [POST]    Modify the behavior and settings of a integration object for the guild.
//                                  Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
//                                  Fires a Guild Integrations Update Gateway event.
// Endpoint                         /guilds/{guild.id}/integrations/{integration.id}
// Rate limiter                     /guilds/{guild.id}/integrations
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#modify-guild-integration
// Reviewed                         2018-08-18
// Comment                          -
func ModifyGuildIntegration(client httd.Patcher, guildID, integrationID Snowflake, params *ModifyGuildIntegrationParams) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildIntegrations(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/integrations/" + integrationID.String(),
		JSONParams:  params,
	}
	resp, _, err := client.Patch(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not modify the integration object. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}

	return
}

// DeleteGuildIntegration [DELETE]  Delete the attached integration object for the guild.
//                                  Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
//                                  Fires a Guild Integrations Update Gateway event.
// Endpoint                         /guilds/{guild.id}/integrations/{integration.id}
// Rate limiter                     /guilds/{guild.id}/integrations
// Discord documentation            https://discordapp.com/developers/docs/resources/guild#delete-guild-integration
// Reviewed                         2018-08-18
// Comment                          -
func DeleteGuildIntegration(client httd.Deleter, guildID, integrationID Snowflake) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildIntegrations(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/integrations/" + integrationID.String(),
	}
	resp, _, err := client.Delete(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove the integration object for the guild. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}

	return
}

// SyncGuildIntegration [POST]  Sync an integration. Requires the 'MANAGE_GUILD' permission.
//                              Returns a 204 empty response on success.
// Endpoint                     /guilds/{guild.id}/integrations/{integration.id}/sync
// Rate limiter                 /guilds/{guild.id}/integrations TODO: is this correct?
// Discord documentation        https://discordapp.com/developers/docs/resources/guild#sync-guild-integration
// Reviewed                     2018-08-18
// Comment                      -
func SyncGuildIntegration(client httd.Poster, guildID, integrationID Snowflake) (err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildIntegrations(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/integrations/" + integrationID.String() + "/sync",
	}
	resp, _, err := client.Post(details)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "could not sync guild integrations. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}
	return
}

// GetGuildEmbed [GET]      Returns the guild embed object. Requires the 'MANAGE_GUILD' permission.
// Endpoint                 /guilds/{guild.id}/embed
// Rate limiter             /guilds/{guild.id}/embed
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-embed
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildEmbed(client httd.Getter, guildID Snowflake) (ret *GuildEmbed, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildEmbed(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/embed",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// ModifyGuildEmbed [PATCH] Modify a guild embed object for the guild. All attributes may be passed in
//                          with JSON and modified. Requires the 'MANAGE_GUILD' permission.
//                          Returns the updated guild embed object.
// Endpoint                 /guilds/{guild.id}/embed
// Rate limiter             /guilds/{guild.id}/embed
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#modify-guild-embed
// Reviewed                 2018-08-18
// Comment                  -
func ModifyGuildEmbed(client httd.Patcher, guildID Snowflake, params *GuildEmbed) (ret *GuildEmbed, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildEmbed(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/embed",
		JSONParams:  params,
	}
	_, body, err := client.Patch(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildVanityURL [GET]  Returns a partial invite object for guilds with that feature enabled.
//                          Requires the 'MANAGE_GUILD' permission.
// Endpoint                 /guilds/{guild.id}/vanity-url
// Rate limiter             /guilds/{guild.id}/vanity-url
// Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-vanity-url
// Reviewed                 2018-08-18
// Comment                  -
func GetGuildVanityURL(client httd.Getter, guildID Snowflake) (ret *PartialInvite, err error) {
	details := &httd.Request{
		Ratelimiter: httd.RatelimitGuildVanityURL(guildID),
		Endpoint:    "/guilds/" + guildID.String() + "/vanity-url",
	}
	_, body, err := client.Get(details)
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}
