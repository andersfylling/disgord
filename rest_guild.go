package disgord

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
)

func ratelimitGuild(id Snowflake) string {
	return "g:" + id.String()
}
func ratelimitGuildAuditLogs(id Snowflake) string {
	return ratelimitGuild(id) + ":a-l"
}
func ratelimitGuildEmbed(id Snowflake) string {
	return ratelimitGuild(id) + ":e"
}
func ratelimitGuildVanityURL(id Snowflake) string {
	return ratelimitGuild(id) + ":vurl"
}
func ratelimitGuildChannels(id Snowflake) string {
	return ratelimitGuild(id) + ":c"
}
func ratelimitGuildMembers(id Snowflake) string {
	return ratelimitGuild(id) + ":m"
}
func ratelimitGuildBans(id Snowflake) string {
	return ratelimitGuild(id) + ":b"
}
func ratelimitGuildRoles(id Snowflake) string {
	return ratelimitGuild(id) + ":r"
}
func ratelimitGuildRegions(id Snowflake) string {
	return ratelimitGuild(id) + ":regions"
}
func ratelimitGuildIntegrations(id Snowflake) string {
	return ratelimitGuild(id) + ":i"
}
func ratelimitGuildInvites(id Snowflake) string {
	return ratelimitGuild(id) + ":inv"
}
func ratelimitGuildPrune(id Snowflake) string {
	return ratelimitGuild(id) + ":p"
}
func ratelimitGuildWebhooks(id Snowflake) string {
	return ratelimitGuild(id) + ":w"
}

// CreateGuildParams ...
// https://discordapp.com/developers/docs/resources/guild#create-guild-json-params
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

// CreateGuild [REST] Create a new guild. Returns a guild object on success. Fires a Guild Create Gateway event.
//  Method                  POST
//  Endpoint                /guilds
//  Rate limiter            /guilds
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild
//  Reviewed                2018-08-16
//  Comment                 This endpoint. can be used only by bots in less than 10 guilds. Creating channel
//                          categories from this endpoint. is not supported.
func CreateGuild(client httd.Poster, params *CreateGuildParams) (ret *Guild, err error) {
	// TODO: check if bot
	// TODO-2: is bot in less than 10 guilds?
	var body []byte
	_, body, err = client.Post(&httd.Request{
		Ratelimiter: endpoint.Guilds(),
		Endpoint:    endpoint.Guilds(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	err = unmarshal(body, &ret)
	return ret, err
}

// GetGuild [REST] Returns the guild object for the given id.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}
//  Rate limiter            /guilds/{guild.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild
//  Reviewed                2018-08-17
//  Comment                 -
func (c *client) GetGuild(id Snowflake, flags ...Flag) (guild *Guild, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.Guild(id),
	}, flags)
	r.factory = func() interface{} {
		return &Guild{}
	}
	r.CacheRegistry = GuildCache
	r.ID = id
	r.preUpdateCache = func(x interface{}) {
		if x == nil {
			return
		}
		if guild, ok := x.(*Guild); ok {
			for i := range guild.Roles {
				guild.Roles[i].guildID = id
			}
		}
	}

	return getGuild(r.Execute)
}

// UpdateGuildParams https://discordapp.com/developers/docs/resources/guild#modify-guild-json-params
// TODO: support nullable Icon, anything else?
type UpdateGuildParams struct {
	Name                    string                        `json:"name,omitempty"`
	Region                  string                        `json:"region,omitempty"`
	VerificationLvl         int                           `json:"verification_level,omitempty"`
	DefaultMsgNotifications DefaultMessageNotificationLvl `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter   ExplicitContentFilterLvl      `json:"explicit_content_filter,omitempty"`
	AFKChannelID            Snowflake                     `json:"afk_channel_id,omitempty"`
	AFKTimeout              int                           `json:"afk_timeout,omitempty"`
	Icon                    string                        `json:"icon,omitempty"`
	OwnerID                 Snowflake                     `json:"owner_id,omitempty"`
	Splash                  string                        `json:"splash,omitempty"`
	SystemChannelID         Snowflake                     `json:"system_channel_id,omitempty"`
}

// ModifyGuild [REST] Modify a guild's settings. Requires the 'MANAGE_GUILD' permission. Returns the updated guild
// object on success. Fires a Guild Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}
//  Rate limiter            /guilds/{guild.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild
//  Reviewed                2018-08-17
//  Comment                 All parameters to this endpoint. are optional
func ModifyGuild(client httd.Patcher, id Snowflake, params *UpdateGuildParams) (ret *Guild, err error) {
	var body []byte
	_, body, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.Guild(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	if err = unmarshal(body, &ret); err != nil {
		return nil, err
	}

	// add guild id to roles
	for _, role := range ret.Roles {
		role.guildID = id
	}

	return ret, nil
}

// DeleteGuild [REST] Delete a guild permanently. User must be owner. Returns 204 No Content on success.
// Fires a Guild Delete Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}
//  Rate limiter            /guilds/{guild.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#delete-guild
//  Reviewed                2018-08-17
//  Comment                 -
func DeleteGuild(client httd.Deleter, id Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Delete(&httd.Request{
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.Guild(id),
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

// GetGuildChannels [REST] Returns a list of guild channel objects.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/channels
//  Rate limiter            /guilds/{guild.id}/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-channels
//  Reviewed                2018-08-17
//  Comment                 -
func GetGuildChannels(client httd.Getter, id Snowflake) (ret []*Channel, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: ratelimitGuildChannels(id),
		Endpoint:    endpoint.GuildChannels(id),
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// CreateGuildChannelParams https://discordapp.com/developers/docs/resources/guild#create-guild-channel-json-params
type CreateGuildChannelParams struct {
	Name                 string                `json:"name"` // required
	Type                 uint                  `json:"type,omitempty"`
	Topic                string                `json:"topic,omitempty"`
	Bitrate              uint                  `json:"bitrate,omitempty"`
	UserLimit            uint                  `json:"user_limit,omitempty"`
	RateLimitPerUser     uint                  `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             Snowflake             `json:"parent_id,omitempty"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
}

// CreateGuildChannel [REST] Create a new channel object for the guild. Requires the 'MANAGE_CHANNELS' permission.
// Returns the new channel object on success. Fires a Channel Create Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/channels
//  Rate limiter            /guilds/{guild.id}/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild-channel
//  Reviewed                2018-08-17
//  Comment                 All parameters for this endpoint. are optional excluding 'name'
func CreateGuildChannel(client httd.Poster, id Snowflake, params *CreateGuildChannelParams) (ret *Channel, err error) {
	var body []byte
	_, body, err = client.Post(&httd.Request{
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.GuildChannels(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// UpdateGuildChannelPositionsParams ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-channel-positions-json-params
type UpdateGuildChannelPositionsParams struct {
	ID       Snowflake `json:"id"`
	Position int       `json:"position"`
}

// ModifyGuildChannelPositions [REST] Modify the positions of a set of channel objects for the guild.
// Requires 'MANAGE_CHANNELS' permission. Returns a 204 empty response on success. Fires multiple Channel Update
// Gateway events.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/channels
//  Rate limiter            /guilds/{guild.id}/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-channel-positions
//  Reviewed                2018-08-17
//  Comment                 Only channels to be modified are required, with the minimum being a swap
//                          between at least two channels.
func ModifyGuildChannelPositions(client httd.Patcher, id Snowflake, params []UpdateGuildChannelPositionsParams) (ret *Guild, err error) {
	var resp *http.Response
	resp, _, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuildChannels(id),
		Endpoint:    endpoint.GuildChannels(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
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

// GetGuildMember [REST] Returns a guild member object for the specified user.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/members/{user.id}
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-member
//  Reviewed                2018-08-17
//  Comment                 -
func (c *client) GetGuildMember(guildID, userID Snowflake, flags ...Flag) (ret *Member, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMember(guildID, userID),
	}, flags)
	r.CacheRegistry = GuildMembersCache
	r.ID = userID
	r.factory = func() interface{} {
		return &Member{}
	}

	return getMember(r.Execute)
}

type GetGuildMembersParams struct {
	After Snowflake `urlparam:"after,omitempty"`
	Limit int       `urlparam:"limit,omitempty"`
}

var _ URLQueryStringer = (*GetGuildMembersParams)(nil)

func (g *GetGuildMembersParams) FindErrors() error {
	if g.Limit > 1000 || g.Limit < 0 {
		return errors.New("limit value should be less than or equal to 1000, and non-negative")
	}
	return nil
}

// GetGuildMembers [REST] Returns a list of guild member objects that are members of the guild. The `after` param
// refers to the highest snowflake.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/members
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-members
//  Reviewed                2018-08-17
//  Comment                 All parameters to this endpoint. are optional
//  Comment#2               "List Guild Members"
//  Comment#3               https://discordapp.com/developers/docs/resources/guild#list-guild-members-query-string-params
func (c *client) GetGuildMembers(guildID Snowflake, params *GetGuildMembersParams, flags ...Flag) (ret []*Member, err error) {
	if params == nil {
		params = &GetGuildMembersParams{}
	}
	if err = params.FindErrors(); err != nil {
		return nil, err
	}

	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMembers(guildID) + params.URLQueryString(),
	}, flags)
	r.CacheRegistry = GuildMembersCache
	r.factory = func() interface{} {
		tmp := make([]*Member, 0)
		return &tmp
	}

	return getMembers(r.Execute)
}

// AddGuildMemberParams ...
// https://discordapp.com/developers/docs/resources/guild#add-guild-member-json-params
type AddGuildMemberParams struct {
	AccessToken string      `json:"access_token"`
	Nick        string      `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles"`
	Mute        bool        `json:"mute"`
	Deaf        bool        `json:"deaf"`
}

// AddGuildMember [REST] Adds a user to the guild, provided you have a valid oauth2 access token for the user with
// the guilds.join scope. Returns a 201 Created with the guild member as the body, or 204 No Content if the user is
// already a member of the guild. Fires a Guild Member Add Gateway event. Requires the bot to have the
// CREATE_INSTANT_INVITE permission.
//  Method                  PUT
//  Endpoint                /guilds/{guild.id}/members/{user.id}
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#add-guild-member
//  Reviewed                2018-08-18
//  Comment                 All parameters to this endpoint. except for access_token are optional.
func AddGuildMember(client httd.Puter, guildID, userID Snowflake, params *AddGuildMemberParams) (ret *Member, err error) {
	var resp *http.Response
	var body []byte
	resp, body, err = client.Put(&httd.Request{
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMember(guildID, userID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
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

// UpdateGuildMemberParams ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-member-json-params
type UpdateGuildMemberParams struct {
	data map[string]interface{}
}

func (m *UpdateGuildMemberParams) init() {
	if m.data != nil {
		return
	}

	m.data = map[string]interface{}{}
}

// SetNick set new nickname for user. Requires permission MANAGE_NICKNAMES
func (m *UpdateGuildMemberParams) SetNick(name string) error {
	if err := ValidateUsername(name); err != nil {
		return err
	}

	m.init()
	m.data["nick"] = name
	return nil
}

// RemoveNick removes nickname for user. Requires permission MANAGE_NICKNAMES
func (m *UpdateGuildMemberParams) RemoveNick() {
	m.init()
	m.data["nick"] = nil
}

// SetRoles updates the member with new roles. Requires permissions MANAGE_ROLES
func (m *UpdateGuildMemberParams) SetRoles(roles []Snowflake) {
	m.init()
	m.data["roles"] = roles
}

// SetMute mutes a member. Requires permission MUTE_MEMBERS
func (m *UpdateGuildMemberParams) SetMute(yes bool) {
	m.init()
	m.data["mute"] = yes
}

// SetDeaf deafens a member. Requires permission DEAFEN_MEMBERS
func (m *UpdateGuildMemberParams) SetDeaf(yes bool) {
	m.init()
	m.data["deaf"] = yes
}

// SetChannelID moves a member from one channel to another. Requires permission MOVE_MEMBERS
func (m *UpdateGuildMemberParams) SetChannelID(id Snowflake) error {
	if id.Empty() {
		return errors.New("empty snowflake")
	}

	m.init()
	m.data["channel_id"] = id
	return nil
}

func (m *UpdateGuildMemberParams) MarshalJSON() ([]byte, error) {
	if len(m.data) == 0 {
		return []byte(`{}`), nil
	}

	return httd.Marshal(m.data)
}

var _ json.Marshaler = (*UpdateGuildMemberParams)(nil)

// ModifyGuildMember [REST] Modify attributes of a guild member. Returns a 204 empty response on success.
// Fires a Guild Member Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/members/{user.id}
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-member
//  Reviewed                2018-08-17
//  Comment                 All parameters to this endpoint. are optional. When moving members to channels,
//                          the API user must have permissions to both connect to the channel and have the
//                          MOVE_MEMBERS permission.
func ModifyGuildMember(client httd.Patcher, guildID, userID Snowflake, params *UpdateGuildMemberParams) (err error) {
	var resp *http.Response
	resp, _, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMember(guildID, userID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "could not change attributes of member. Does the member exist, and do you have permissions?"
		err = errors.New(msg)
	}

	return
}

// AddGuildMemberRole [REST] Adds a role to a guild member. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
//  Method                  PUT
//  Endpoint                /guilds/{guild.id}/members/{user.id}/roles/{role.id}
//  Rate limiter            /guilds/{guild.id}/members TODO: I don't know if this is correct
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#add-guild-member-role
//  Reviewed                2018-08-18
//  Comment                 -
func AddGuildMemberRole(client httd.Puter, guildID, userID, roleID Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Put(&httd.Request{
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMemberRole(guildID, userID, roleID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not add role to user. Do you have the MANAGE_ROLES permission?"
		err = errors.New(msg)
	}

	return
}

// RemoveGuildMemberRole [REST] Removes a role from a guild member. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/members/{user.id}/roles/{role.id}
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#remove-guild-member-role
//  Reviewed                2018-08-18
//  Comment                 -
func RemoveGuildMemberRole(client httd.Deleter, guildID, userID, roleID Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Delete(&httd.Request{
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMemberRole(guildID, userID, roleID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove role from user. Do you have the MANAGE_ROLES permission?"
		err = errors.New(msg)
	}

	return
}

// RemoveGuildMember [REST] Remove a member from a guild. Requires 'KICK_MEMBERS' permission.
// Returns a 204 empty response on success. Fires a Guild Member Remove Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/members/{user.id}
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#remove-guild-member
//  Reviewed                2018-08-18
//  Comment                 -
func RemoveGuildMember(client httd.Deleter, guildID, userID Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Delete(&httd.Request{
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMember(guildID, userID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove user from guild. Do you have the KICK_MEMBERS permission?"
		err = errors.New(msg)
	}

	return
}

// GetGuildBans [REST] Returns a list of ban objects for the users banned from this guild. Requires the 'BAN_MEMBERS' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/bans
//  Rate limiter            /guilds/{guild.id}/bans
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-bans
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) GetGuildBans(id Snowflake, flags ...Flag) (bans []*Ban, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildBans(id),
		Endpoint:    endpoint.GuildBans(id),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Ban, 0)
		return &tmp
	}

	var vs interface{}
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if cons, ok := vs.(*[]*Ban); ok {
		return *cons, nil
	}
	return nil, errors.New("unable to cast guild slice")
}

// GetGuildBan [REST] Returns a ban object for the given user or a 404 not found if the ban cannot be found.
// Requires the 'BAN_MEMBERS' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/bans/{user.id}
//  Rate limiter            /guilds/{guild.id}/bans
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-ban
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) GetGuildBan(guildID, userID Snowflake, flags ...Flag) (ret *Ban, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildBans(guildID),
		Endpoint:    endpoint.GuildBan(guildID, userID),
	}, flags)
	r.factory = func() interface{} {
		return &Ban{User: c.pool.user.Get().(*User)}
	}

	return getBan(r.Execute)
}

// BanMemberParams ...
// https://discordapp.com/developers/docs/resources/guild#create-guild-ban-query-string-params
type BanMemberParams struct {
	DeleteMessageDays int    `urlparam:"delete_message_days,omitempty"` // number of days to delete messages for (0-7)
	Reason            string `urlparam:"reason,omitempty"`              // reason for being banned
}

var _ URLQueryStringer = (*BanMemberParams)(nil)

func (b *BanMemberParams) FindErrors() error {
	if !(0 <= b.DeleteMessageDays && b.DeleteMessageDays <= 7) {
		return errors.New("DeleteMessageDays must be a value in the range of [0, 7], got " + strconv.Itoa(b.DeleteMessageDays))
	}
	return nil
}

// BanMember [REST] Create a guild ban, and optionally delete previous messages sent by the banned user. Requires
// the 'BAN_MEMBERS' permission. Returns a 204 empty response on success. Fires a Guild Ban Add Gateway event.
//  Method                  PUT
//  Endpoint                /guilds/{guild.id}/bans/{user.id}
//  Rate limiter            /guilds/{guild.id}/bans
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild-ban
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) BanMember(guildID, userID Snowflake, params *BanMemberParams, flags ...Flag) (err error) {
	if params == nil {
		return errors.New("params was nil")
	}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPut,
		Ratelimiter: ratelimitGuildBans(guildID),
		Endpoint:    endpoint.GuildBan(guildID, userID) + params.URLQueryString(),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return
}

// UnbanMember [REST] Remove the ban for a user. Requires the 'BAN_MEMBERS' permissions.
// Returns a 204 empty response on success. Fires a Guild Ban Remove Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/bans/{user.id}
//  Rate limiter            /guilds/{guild.id}/bans
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#remove-guild-ban
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) UnbanMember(guildID, userID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitGuildBans(guildID),
		Endpoint:    endpoint.GuildBan(guildID, userID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return
}

// PruneMembersParams will delete members, this is the same as kicking.
// https://discordapp.com/developers/docs/resources/guild#get-guild-prune-count-query-string-params
type pruneMembersParams struct {
	// Days number of days to count prune for (1 or more)
	Days int `urlparam:"days"`

	// ComputePruneCount whether 'pruned' is returned, discouraged for large guilds
	ComputePruneCount bool `urlparam:"compute_prune_count"`
}

func (d *pruneMembersParams) FindErrors() (err error) {
	if d.Days < 1 {
		err = errors.New("days must be at least 1, got " + strconv.Itoa(d.Days))
	}
	return
}

var _ URLQueryStringer = (*pruneMembersParams)(nil)

// GuildPruneCount ...
type guildPruneCount struct {
	Pruned int `json:"pruned"`
}

// EstimatePruneMembersCount [REST] Returns an object with one 'pruned' key indicating the number of members that would be
// removed in a prune operation. Requires the 'KICK_MEMBERS' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/prune
//  Rate limiter            /guilds/{guild.id}/prune
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-prune-count
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) EstimatePruneMembersCount(id Snowflake, days int, flags ...Flag) (estimate int, err error) {
	if id.Empty() {
		return 0, errors.New("guildID can not be " + id.String())
	}
	params := pruneMembersParams{Days: days}
	if err = params.FindErrors(); err != nil {
		return 0, err
	}

	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildPrune(id),
		Endpoint:    endpoint.GuildPrune(id) + params.URLQueryString(),
	}, flags)
	r.factory = func() interface{} {
		return &guildPruneCount{}
	}

	var v interface{}
	if v, err = r.Execute(); err != nil {
		return 0, err
	}

	if v == nil {
		return 0, nil
	}

	return v.(*guildPruneCount).Pruned, nil
}

// PruneMembers [REST] Kicks members from N day back. Requires the 'KICK_MEMBERS' permission.
// The estimate of kicked people is not returned. Use EstimatePruneMembersCount before calling PruneMembers
// if you need it.
// Fires multiple Guild Member Remove Gateway events.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/prune
//  Rate limiter            /guilds/{guild.id}/prune
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#begin-guild-prune
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) PruneMembers(id Snowflake, days int, flags ...Flag) (err error) {
	params := pruneMembersParams{Days: days}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: ratelimitGuildPrune(id),
		Endpoint:    endpoint.GuildPrune(id) + params.URLQueryString(),
	}, flags)

	_, err = r.Execute()
	return err
}

// GetGuildVoiceRegions [REST] Returns a list of voice region objects for the guild. Unlike the similar /voice route,
// this returns VIP servers when the guild is VIP-enabled.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/regions
//  Rate limiter            /guilds/{guild.id}/regions
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-voice-regions
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) GetGuildVoiceRegions(id Snowflake, flags ...Flag) (ret []*VoiceRegion, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildRegions(id),
		Endpoint:    endpoint.GuildRegions(id),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*VoiceRegion, 0)
		return &tmp
	}

	return getVoiceRegions(r.Execute)
}

// GetGuildInvites [REST] Returns a list of invite objects (with invite metadata) for the guild.
// Requires the 'MANAGE_GUILD' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/invites
//  Rate limiter            /guilds/{guild.id}/invites
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-invites
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) GetGuildInvites(id Snowflake, flags ...Flag) (ret []*Invite, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildInvites(id),
		Endpoint:    endpoint.GuildInvites(id),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Invite, 0)
		return &tmp
	}

	return getInvites(r.Execute)
}

// GetGuildIntegrations [REST] Returns a list of integration objects for the guild.
// Requires the 'MANAGE_GUILD' permission.
//  Method                  GET
//  Endpoint                 /guilds/{guild.id}/integrations
//  Rate limiter             /guilds/{guild.id}/integrations
//  Discord documentation    https://discordapp.com/developers/docs/resources/guild#get-guild-integrations
//  Reviewed                 2018-08-18
//  Comment                  -
func (c *client) GetGuildIntegrations(id Snowflake, flags ...Flag) (ret []*Integration, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildIntegrations(id),
		Endpoint:    endpoint.GuildIntegrations(id),
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Integration, 0)
		return &tmp
	}

	return getIntegrations(r.Execute)
}

// CreateGuildIntegrationParams ...
// https://discordapp.com/developers/docs/resources/guild#create-guild-integration-json-params
type CreateGuildIntegrationParams struct {
	Type string    `json:"type"`
	ID   Snowflake `json:"id"`
}

// CreateGuildIntegration [REST] Attach an integration object from the current user to the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/integrations
//  Rate limiter            /guilds/{guild.id}/integrations
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild-integration
//  Reviewed                2018-08-18
//  Comment                 -
func CreateGuildIntegration(client httd.Poster, guildID Snowflake, params *CreateGuildIntegrationParams) (err error) {
	var resp *http.Response
	resp, _, err = client.Post(&httd.Request{
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegrations(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not create the integration object. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}

	return
}

// UpdateGuildIntegrationParams ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-integration-json-params
type UpdateGuildIntegrationParams struct {
	ExpireBehavior    int  `json:"expire_behavior"`
	ExpireGracePeriod int  `json:"expire_grace_period"`
	EnableEmoticons   bool `json:"enable_emoticons"`
}

// UpdateGuildIntegration [REST] Modify the behavior and settings of a integration object for the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/integrations/{integration.id}
//  Rate limiter            /guilds/{guild.id}/integrations
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-integration
//  Reviewed                2018-08-18
//  Comment                 -
func UpdateGuildIntegration(client httd.Patcher, guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams) (err error) {
	var resp *http.Response
	resp, _, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegration(guildID, integrationID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not modify the integration object. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}

	return
}

// DeleteGuildIntegration [REST] Delete the attached integration object for the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}/integrations/{integration.id}
//  Rate limiter            /guilds/{guild.id}/integrations
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#delete-guild-integration
//  Reviewed                2018-08-18
//  Comment                 -
func DeleteGuildIntegration(client httd.Deleter, guildID, integrationID Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Delete(&httd.Request{
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegration(guildID, integrationID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "Could not remove the integration object for the guild. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}

	return
}

// SyncGuildIntegration [REST] Sync an integration. Requires the 'MANAGE_GUILD' permission.
// Returns a 204 empty response on success.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/integrations/{integration.id}/sync
//  Rate limiter            /guilds/{guild.id}/integrations
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#sync-guild-integration
//  Reviewed                2018-08-18
//  Comment                 -
func SyncGuildIntegration(client httd.Poster, guildID, integrationID Snowflake) (err error) {
	var resp *http.Response
	resp, _, err = client.Post(&httd.Request{
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegrationSync(guildID, integrationID),
	})
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		msg := "could not sync guild integrations. Do you have the MANAGE_GUILD permission?"
		err = errors.New(msg)
	}
	return
}

// GetGuildEmbed [REST] Returns the guild embed object. Requires the 'MANAGE_GUILD' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/embed
//  Rate limiter            /guilds/{guild.id}/embed
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-embed
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) GetGuildEmbed(guildID Snowflake, flags ...Flag) (embed *GuildEmbed, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildEmbed(guildID),
		Endpoint:    endpoint.GuildEmbed(guildID),
	}, flags)
	r.factory = func() interface{} {
		return &GuildEmbed{}
	}

	return getGuildEmbed(r.Execute)
}

// ModifyGuildEmbed [REST] Modify a guild embed object for the guild. All attributes may be passed in with JSON and
// modified. Requires the 'MANAGE_GUILD' permission. Returns the updated guild embed object.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/embed
//  Rate limiter            /guilds/{guild.id}/embed
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-embed
//  Reviewed                2018-08-18
//  Comment                 -
func ModifyGuildEmbed(client httd.Patcher, guildID Snowflake, params *GuildEmbed) (ret *GuildEmbed, err error) {
	var body []byte
	_, body, err = client.Patch(&httd.Request{
		Ratelimiter: ratelimitGuildEmbed(guildID),
		Endpoint:    endpoint.GuildEmbed(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &ret)
	return
}

// GetGuildVanityURL [REST] Returns a partial invite object for guilds with that feature enabled.
// Requires the 'MANAGE_GUILD' permission.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/vanity-url
//  Rate limiter            /guilds/{guild.id}/vanity-url
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-vanity-url
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) GetGuildVanityURL(guildID Snowflake, flags ...Flag) (ret *PartialInvite, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildVanityURL(guildID),
		Endpoint:    endpoint.GuildVanityURL(guildID),
	}, flags)
	r.factory = func() interface{} {
		return &PartialInvite{}
	}

	return getPartialInvite(r.Execute)
}
