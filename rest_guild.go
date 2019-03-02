package disgord

import (
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
	Name                    string                        `json:"name"` // required
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
//							The params argument is optional.
func (c *client) CreateGuild(guildName string, params *CreateGuildParams, flags ...Flag) (ret *Guild, err error) {
	// TODO: check if bot
	// TODO-2: is bot in less than 10 guilds?

	if guildName == "" {
		return nil, errors.New("guild name is required")
	}
	if l := len(guildName); !(2 <= l && l <= 100) {
		return nil, errors.New("guild name must be 2 or more characters and no more than 100 characters")
	}

	if params == nil {
		params = &CreateGuildParams{}
	}
	params.Name = guildName

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: endpoint.Guilds(),
		Endpoint:    endpoint.Guilds(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Guild{}
	}
	r.CacheRegistry = GuildCache

	return getGuild(r.Execute)
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

	return getGuild(r.Execute)
}

// ModifyGuild [REST] Modify a guild's settings. Requires the 'MANAGE_GUILD' permission. Returns the updated guild
// object on success. Fires a Guild Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}
//  Rate limiter            /guilds/{guild.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild
//  Reviewed                2018-08-17
//  Comment                 All parameters to this endpoint. are optional
func (c *client) UpdateGuild(id Snowflake, flags ...Flag) (builder *updateGuildBuilder) {
	builder = &updateGuildBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Guild{}
	}
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.Guild(id),
		ContentType: httd.ContentTypeJSON,
	}, nil)
	builder.r.cacheRegistry = GuildCache
	builder.r.cacheItemID = id
	builder.r.flags = flags

	return builder
}

// DeleteGuild [REST] Delete a guild permanently. User must be owner. Returns 204 No Content on success.
// Fires a Guild Delete Gateway event.
//  Method                  DELETE
//  Endpoint                /guilds/{guild.id}
//  Rate limiter            /guilds/{guild.id}
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#delete-guild
//  Reviewed                2018-08-17
//  Comment                 -
func (c *client) DeleteGuild(id Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.Guild(id),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return
}

// GetGuildChannels [REST] Returns a list of guild channel objects.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/channels
//  Rate limiter            /guilds/{guild.id}/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-channels
//  Reviewed                2018-08-17
//  Comment                 -
func (c *client) GetGuildChannels(id Snowflake, flags ...Flag) (ret []*Channel, err error) {
	r := c.newRESTRequest(&httd.Request{
		Ratelimiter: ratelimitGuildChannels(id),
		Endpoint:    endpoint.GuildChannels(id),
	}, flags)
	r.CacheRegistry = ChannelCache
	r.factory = func() interface{} {
		tmp := make([]*Channel, 0)
		return &tmp
	}
	// TODO: update guild cache

	return getChannels(r.Execute)
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
func (c *client) CreateGuildChannel(id Snowflake, channelName string, params *CreateGuildChannelParams, flags ...Flag) (ret *Channel, err error) {
	if channelName == "" && (params == nil || params.Name == "") {
		return nil, errors.New("channel name is required")
	}
	if l := len(channelName); !(2 <= l && l <= 100) {
		return nil, errors.New("channel name must be 2 or more characters and no more than 100 characters")
	}

	if params == nil {
		params = &CreateGuildChannelParams{}
	}
	if channelName != "" && params.Name == "" {
		params.Name = channelName
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: ratelimitGuild(id),
		Endpoint:    endpoint.GuildChannels(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Channel{}
	}
	r.CacheRegistry = ChannelCache
	// TODO: update guild cache

	return getChannel(r.Execute)
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
func (c *client) UpdateGuildChannelPositions(id Snowflake, params []UpdateGuildChannelPositionsParams, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuildChannels(id),
		Endpoint:    endpoint.GuildChannels(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent
	// TODO: update ordering of guild channels in cache

	_, err = r.Execute()
	return err
}

// UpdateGuildRolePositionsParams ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-role-positions-json-params
type UpdateGuildRolePositionsParams struct {
	ID       Snowflake `json:"id"`
	Position uint      `json:"position"`
}

// UpdateGuildRolePositions [REST] Modify the positions of a set of role objects for the guild.
// Requires the 'MANAGE_ROLES' permission. Returns a list of all of the guild's role objects on success.
// Fires multiple Guild Role Update Gateway events.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/roles
//  Rate limiter            /guilds/{guild.id}/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-role-positions
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) UpdateGuildRolePositions(guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) (roles []*Role, err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuildRoles(guildID),
		Endpoint:    endpoint.GuildRoles(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Role, 0)
		return &tmp
	}
	// TODO: update ordering of guild roles in cache

	return getRoles(r.Execute)
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
	AccessToken string      `json:"access_token"` // required
	Nick        string      `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Mute        bool        `json:"mute,omitempty"`
	Deaf        bool        `json:"deaf,omitempty"`
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
func (c *client) AddGuildMember(guildID, userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (member *Member, err error) {
	if accessToken == "" && (params == nil || params.AccessToken == "") {
		return nil, errors.New("access token is required")
	}

	if params == nil {
		params = &AddGuildMemberParams{}
	}
	if accessToken != "" && params.AccessToken == "" {
		params.AccessToken = accessToken
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPut,
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMember(guildID, userID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Member{}
	}
	r.expectsStatusCode = http.StatusCreated

	// TODO: update guild cache
	if member, err = getMember(r.Execute); err != nil {
		if errRest, ok := err.(*httd.ErrREST); ok && errRest.HTTPCode == http.StatusNoContent {
			errRest.Msg = "member{" + userID.String() + "} is already in Guild{" + guildID.String() + "}"
		}
	}

	return member, err
}

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
func (c *client) UpdateGuildMember(guildID, userID Snowflake, flags ...Flag) (builder *updateGuildMemberBuilder) {
	builder = &updateGuildMemberBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Member{}
	}
	builder.r.flags = flags
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMember(guildID, userID),
		ContentType: httd.ContentTypeJSON,
	}, func(resp *http.Response, body []byte, err error) error {
		if resp.StatusCode != http.StatusNoContent {
			msg := "could not change attributes of member. Does the member exist, and do you have permissions?"
			return errors.New(msg)
		}
		return nil
	})

	// TODO: cache member changes
	return builder
}

// AddGuildMemberRole [REST] Adds a role to a guild member. Requires the 'MANAGE_ROLES' permission.
// Returns a 204 empty response on success. Fires a Guild Member Update Gateway event.
//  Method                  PUT
//  Endpoint                /guilds/{guild.id}/members/{user.id}/roles/{role.id}
//  Rate limiter            /guilds/{guild.id}/members/roles
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#add-guild-member-role
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) AddGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPut,
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMemberRole(guildID, userID, roleID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
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
func (c *client) RemoveGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMemberRole(guildID, userID, roleID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
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
func (c *client) KickMember(guildID, userID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitGuildMembers(guildID),
		Endpoint:    endpoint.GuildMember(guildID, userID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
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

var _ URLQueryStringer = (*pruneMembersParams)(nil)

func (d *pruneMembersParams) FindErrors() (err error) {
	if d.Days < 1 {
		err = errors.New("days must be at least 1, got " + strconv.Itoa(d.Days))
	}
	return
}

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
//  Method                   GET
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
func (c *client) CreateGuildIntegration(guildID Snowflake, params *CreateGuildIntegrationParams, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegrations(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// UpdateGuildIntegrationParams ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-integration-json-params
// TODO: currently unsure which are required/optional params
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
func (c *client) UpdateGuildIntegration(guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegration(guildID, integrationID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
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
func (c *client) DeleteGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegration(guildID, integrationID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// SyncGuildIntegration [REST] Sync an integration. Requires the 'MANAGE_GUILD' permission.
// Returns a 204 empty response on success.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/integrations/{integration.id}/sync
//  Rate limiter            /guilds/{guild.id}/integrations
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#sync-guild-integration
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) SyncGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ratelimiter: ratelimitGuildIntegrations(guildID),
		Endpoint:    endpoint.GuildIntegrationSync(guildID, integrationID),
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err = r.Execute()
	return err
}

// updateCurrentUserNickParams ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-member-json-params
type updateCurrentUserNickParams struct {
	Nick string `json:"nick"` // :CHANGE_NICKNAME
}

type nickNameResponse struct {
	Nickname string `json:"nickname"`
}

// SetCurrentUserNick [REST] Modifies the nickname of the current user in a guild. Returns a 200
// with the nickname on success. Fires a Guild Member Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/members/@me/nick
//  Rate limiter            /guilds/{guild.id}/members/@me/nick
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-current-user-nick
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) SetCurrentUserNick(id Snowflake, nick string, flags ...Flag) (newNick string, err error) {
	params := &updateCurrentUserNickParams{
		Nick: nick,
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuildMembers(id),
		Endpoint:    endpoint.GuildMembersMeNick(id),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusOK
	r.factory = func() interface{} {
		return &nickNameResponse{}
	}

	return getNickName(r.Execute)
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

// UpdateGuildEmbed [REST] Modify a guild embed object for the guild. All attributes may be passed in with JSON and
// modified. Requires the 'MANAGE_GUILD' permission. Returns the updated guild embed object.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/embed
//  Rate limiter            /guilds/{guild.id}/embed
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-embed
//  Reviewed                2018-08-18
//  Comment                 -
func (c *client) UpdateGuildEmbed(guildID Snowflake, flags ...Flag) (builder *updateGuildEmbedBuilder) {
	builder = &updateGuildEmbedBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &GuildEmbed{}
	}
	builder.r.flags = flags
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      http.MethodPatch,
		Ratelimiter: ratelimitGuildEmbed(guildID),
		Endpoint:    endpoint.GuildEmbed(guildID),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
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

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// updateGuildBuilder https://discordapp.com/developers/docs/resources/guild#modify-guild-json-params
//generate-rest-params: name:string, region:string, verification_level:int, default_message_notifications:DefaultMessageNotificationLvl, explicit_content_filter:ExplicitContentFilterLvl, afk_channel_id:Snowflake, afk_timeout:int, icon:string, owner_id:Snowflake, splash:string, system_channel_id:Snowflake,
//generate-rest-basic-execute: guild:*Guild,
type updateGuildBuilder struct {
	r RESTBuilder
}

//generate-rest-params: enabled:bool, channel_id:Snowflake,
//generate-rest-basic-execute: embed:*GuildEmbed,
type updateGuildEmbedBuilder struct {
	r RESTBuilder
}

// updateGuildMemberBuilder ...
// https://discordapp.com/developers/docs/resources/guild#modify-guild-member-json-params
//generate-rest-params: nick:string, roles:[]Snowflake, mute:bool, deaf:bool, channel_id:Snowflake,
//generate-rest-basic-execute: err:error,
type updateGuildMemberBuilder struct {
	r RESTBuilder
}

// RemoveNick removes nickname for user. Requires permission MANAGE_NICKNAMES
func (b *updateGuildMemberBuilder) SetDefaultNick() *updateGuildMemberBuilder {
	b.r.param("nick", nil)
	return b
}
