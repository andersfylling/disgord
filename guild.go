package disgord

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/andersfylling/disgord/json"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

// NOTE! Credit for defining the Permission consts in a clean way goes to DiscordGo.
// This is pretty much a copy from their project. I would have made it a dependency if
// the consts were in a isolated sub-pkg. Note that in respect to their license, Disgord
// has no affiliation with DiscordGo.
//
// Source code reference:
//  https://github.com/bwmarrin/discordgo/blob/8325a6bf6dd6c91ed4040a1617b07287b8fb0eba/structs.go#L854

// PermissionBit is used to define the permission bit(s) which are set.
type PermissionBit uint64

var _ json.Unmarshaler = (*PermissionBit)(nil)
var _ json.Marshaler = (*PermissionBit)(nil)

func (b *PermissionBit) MarshalJSON() ([]byte, error) {
	str := strconv.FormatUint(uint64(*b), 10)
	return []byte(strconv.Quote(str)), nil
}

func (b *PermissionBit) UnmarshalJSON(bytes []byte) error {
	sb := string(bytes)
	str, err := strconv.Unquote(sb)
	if err != nil {
		return fmt.Errorf("PermissionBit#UnmarshalJSON - unable to unquote bytes{%s}: %w", sb, err)
	}

	v, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("PermissionBit#UnmarshalJSON - parsing string to uint64 failed: %w", err)
	}

	*b = PermissionBit(v)
	return nil
}

// Contains is used to check if the permission bits contains the bits specified.
func (b PermissionBit) Contains(Bits PermissionBit) bool {
	return (b & Bits) == Bits
}

// Constants for the different bit offsets of text channel permissions
const (
	PermissionReadMessages PermissionBit = 1 << (iota + 10)
	PermissionSendMessages
	PermissionSendTTSMessages
	PermissionManageMessages
	PermissionEmbedLinks
	PermissionAttachFiles
	PermissionReadMessageHistory
	PermissionMentionEveryone
	PermissionUseExternalEmojis
	PermissionViewGuildInsights
)

// Constants for the different bit offsets of voice permissions
const (
	PermissionVoiceConnect PermissionBit = 1 << (iota + 20)
	PermissionVoiceSpeak
	PermissionVoiceMuteMembers
	PermissionVoiceDeafenMembers
	PermissionVoiceMoveMembers
	PermissionVoiceUseVAD
	PermissionVoicePrioritySpeaker PermissionBit = 1 << (iota + 2)
	PermissionStream
)

// Constants for general management.
const (
	PermissionChangeNickname PermissionBit = 1 << (iota + 26)
	PermissionManageNicknames
	PermissionManageRoles
	PermissionManageWebhooks
	PermissionManageEmojis
	_
	_
	_
	PermissionManageThreads
	PermissionCreatePublicThreads
	PermissionCreatePrivateThreads
	_
	PermissionSendMessagesInThreads
)

// Constants for the different bit offsets of general permissions
const (
	PermissionCreateInstantInvite PermissionBit = 1 << iota
	PermissionKickMembers
	PermissionBanMembers
	PermissionAdministrator
	PermissionManageChannels
	PermissionManageServer
	PermissionAddReactions
	PermissionViewAuditLogs

	PermissionTextAll = PermissionReadMessages |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone
	PermissionAllVoice = PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD
	PermissionChannelAll = PermissionTextAll |
		PermissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	PermissionAll = PermissionChannelAll |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator
)

// GuildUnavailable is a partial Guild object.
type GuildUnavailable struct {
	ID          Snowflake `json:"id"`
	Unavailable bool      `json:"unavailable"` // ?*|
}

// Guild Guilds in Discord represent an isolated collection of Users and Channels,
//  and are often referred to as "servers" in the UI.
// https://discord.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
// reviewed: 2018-08-25
type Guild struct {
	ID                          Snowflake                     `json:"id"`
	ApplicationID               Snowflake                     `json:"application_id"` //   |?
	Name                        string                        `json:"name"`
	Icon                        string                        `json:"icon"`            //  |?, icon hash
	Splash                      string                        `json:"splash"`          //  |?, image hash
	Owner                       bool                          `json:"owner,omitempty"` // ?|
	OwnerID                     Snowflake                     `json:"owner_id"`
	Permissions                 PermissionBit                 `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
	Region                      string                        `json:"region"`
	AfkChannelID                Snowflake                     `json:"afk_channel_id"` // |?
	AfkTimeout                  uint                          `json:"afk_timeout"`
	VerificationLevel           VerificationLvl               `json:"verification_level"`
	DefaultMessageNotifications DefaultMessageNotificationLvl `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLvl      `json:"explicit_content_filter"`
	Roles                       []*Role                       `json:"roles"`
	Emojis                      []*Emoji                      `json:"emojis"`
	Features                    []string                      `json:"features"`
	MFALevel                    MFALvl                        `json:"mfa_level"`
	WidgetEnabled               bool                          `json:"widget_enabled,omit_empty"`    //   |
	WidgetChannelID             Snowflake                     `json:"widget_channel_id,omit_empty"` //   |?
	SystemChannelID             Snowflake                     `json:"system_channel_id,omitempty"`  //   |?
	DiscoverySplash             string                        `json:"discovery_splash,omitempty"`
	VanityUrl                   string                        `json:"vanity_url_code,omitempty"`
	Description                 string                        `json:"description,omitempty"`
	Banner                      string                        `json:"banner,omitempty"`
	PremiumTier                 PremiumTier                   `json:"premium_tier"`
	PremiumSubscriptionCount    uint                          `json:"premium_subscription_count,omitempty"`

	// JoinedAt must be a pointer, as we can't hide non-nil structs
	JoinedAt    *Time           `json:"joined_at,omitempty"`    // ?*|
	Large       bool            `json:"large,omitempty"`        // ?*|
	Unavailable bool            `json:"unavailable"`            // ?*| omitempty?
	MemberCount uint            `json:"member_count,omitempty"` // ?*|
	VoiceStates []*VoiceState   `json:"voice_states,omitempty"` // ?*|
	Members     []*Member       `json:"members,omitempty"`      // ?*|
	Channels    []*Channel      `json:"channels,omitempty"`     // ?*|
	Presences   []*UserPresence `json:"presences,omitempty"`    // ?*|
}

var _ Reseter = (*Guild)(nil)
var _ fmt.Stringer = (*Guild)(nil)
var _ Copier = (*Guild)(nil)
var _ DeepCopier = (*Guild)(nil)
var _ internalUpdater = (*Guild)(nil)

func (g *Guild) String() string {
	return g.Name + "{" + g.ID.String() + "}"
}

func (g *Guild) updateInternals() {
	for i := range g.Roles {
		g.Roles[i].guildID = g.ID
	}
	for i := range g.Channels {
		g.Channels[i].GuildID = g.ID
	}
	for i := range g.Members {
		g.Members[i].GuildID = g.ID
		g.Members[i].updateInternals()
	}
}

// GetMemberWithHighestSnowflake finds the member with the highest snowflake value.
func (g *Guild) GetMemberWithHighestSnowflake() *Member {
	if len(g.Members) == 0 {
		return nil
	}

	highest := g.Members[0]
	for _, member := range g.Members {
		if member.User.ID > highest.User.ID {
			highest = member
		}
	}

	return highest
}

// sortChannels Only while in lock
func (g *Guild) sortChannels() {
	Sort(g.Channels, SortByID)
}

// AddChannel adds a channel to the Guild object. Note that this method does not interact with Discord.
func (g *Guild) AddChannel(c *Channel) error {
	g.Channels = append(g.Channels, c)
	g.sortChannels()

	return nil
}

// DeleteChannel removes a channel from the Guild object. Note that this method does not interact with Discord.
func (g *Guild) DeleteChannel(c *Channel) error {
	return g.DeleteChannelByID(c.ID)
}

// DeleteChannelByID removes a channel from the Guild object. Note that this method does not interact with Discord.
func (g *Guild) DeleteChannelByID(ID Snowflake) error {
	index := -1
	for i, c := range g.Channels {
		if c.ID == ID {
			index = i
		}
	}

	if index == -1 {
		return errors.New("channel with Snowflake{" + ID.String() + "} does not exist in cacheLink")
	}

	// delete the entry
	copy(g.Channels[index:], g.Channels[index+1:])
	g.Channels[len(g.Channels)-1] = nil // or the zero value of T
	g.Channels = g.Channels[:len(g.Channels)-1]

	return nil
}

func (g *Guild) hasMember(id Snowflake) bool {
	for i := len(g.Members) - 1; i >= 0; i-- {
		if g.Members[i].UserID == id {
			return true
		}
	}

	return false
}

func (g *Guild) addMembers(members ...*Member) {
	// TODO: implement sorting for faster searching later
	g.Members = append(g.Members, members...)
}

// AddMembers adds multiple members to the Guild object. Note that this method does not interact with Discord.
func (g *Guild) AddMembers(members []*Member) {
	// Reduces allocations
	membersToAdd := members[:0]

	for _, member := range members {
		// TODO: Check for UserID.IsZero()
		if !g.hasMember(member.UserID) {
			membersToAdd = append(membersToAdd, member)
		}
	}

	g.addMembers(membersToAdd...)
}

// AddMember adds a member to the Guild object. Note that this method does not interact with Discord.
func (g *Guild) AddMember(member *Member) error {
	// TODO: Check for UserID.IsZero()
	if !g.hasMember(member.UserID) {
		g.addMembers(member)
	}

	return nil
}

// GetMembersCountEstimate estimates the number of members in a guild without fetching everyone.
// There is no proper way to get this number, so a invite is created and the estimate
// is read from there. The invite is then deleted again.
func (g *Guild) GetMembersCountEstimate(ctx context.Context, s Session) (estimate int, err error) {
	var channelID Snowflake
	if len(g.Channels) == 0 {
		channels, err := s.Guild(g.ID).WithContext(ctx).GetChannels()
		if err != nil {
			return 0, err
		}

		for i := range channels {
			channelID = channels[i].ID

			// prefer the main channel
			if channelID == g.ID {
				break
			}
		}

		// TODO: update g.Channels
	}
	if channelID.IsZero() {
		return 0, MissingChannelIDErr
	}

	invite, err := s.Channel(channelID).WithContext(ctx).CreateInvite().
		SetMaxAge(1).
		Execute()
	if err != nil {
		return 0, err
	}
	_, _ = s.Invite(invite.Code).WithContext(ctx).Delete() // delete if possible

	return invite.ApproximateMemberCount, nil
}

// AddRole adds a role to the Guild object. Note that this does not interact with Discord.
func (g *Guild) AddRole(role *Role) error {
	// TODO: implement sorting for faster searching later
	role.guildID = g.ID
	g.Roles = append(g.Roles, role)

	return nil
}

// Member return a member by his/her userid
func (g *Guild) Member(id Snowflake) (*Member, error) {
	for _, member := range g.Members {
		if member.UserID == id {
			return member, nil
		}
	}

	return nil, errors.New("member not found in guild")
}

// MembersByName retrieve a slice of members with same username or nickname
func (g *Guild) MembersByName(name string) (members []*Member) {
	for _, member := range g.Members {
		if member.Nick == name || member.User.Username == name {
			members = append(members, member)
		}
	}

	return
}

// Role retrieve a role based on role id
func (g *Guild) Role(id Snowflake) (*Role, error) {
	for _, role := range g.Roles {
		if role.ID == id {
			return role, nil
		}
	}
	return nil, errors.New("role not found in guild")
}

// DeleteRoleByID remove a role from the guild struct
func (g *Guild) DeleteRoleByID(ID Snowflake) {
	index := -1
	for i, r := range g.Roles {
		if r.ID == ID {
			index = i
			break
		}
	}

	if index != -1 {
		// delete the entry
		g.Roles[index] = g.Roles[len(g.Roles)-1]
		g.Roles[len(g.Roles)-1] = nil
		g.Roles = g.Roles[:len(g.Roles)-1]
	}
}

// RoleByName retrieves a slice of roles with same name
func (g *Guild) RoleByName(name string) ([]*Role, error) {
	var roles []*Role
	for _, role := range g.Roles {
		if role.Name == name {
			roles = append(roles, role)
		}
	}

	if len(roles) == 0 {
		return nil, errors.New("no roles were found in guild")
	}

	return roles, nil
}

// Channel get a guild channel given it's ID
func (g *Guild) Channel(id Snowflake) (*Channel, error) {
	for _, channel := range g.Channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, errors.New("channel not found in guild")
}

// Emoji get a guild emoji by it's ID
func (g *Guild) Emoji(id Snowflake) (emoji *Emoji, err error) {
	for _, emoji = range g.Emojis {
		if emoji.ID == id {
			return
		}
	}

	err = errors.New("emoji not found in guild")
	return
}

// --------------

// PartialBan is used by audit logs
type PartialBan struct {
	Reason                 string
	BannedUserID           Snowflake
	ModeratorResponsibleID Snowflake
}

var _ fmt.Stringer = (*PartialBan)(nil)

func (p *PartialBan) String() string {
	return fmt.Sprintf("mod{%d} banned member{%d}, reason: %s.", p.ModeratorResponsibleID, p.BannedUserID, p.Reason)
}

// Ban https://discord.com/developers/docs/resources/guild#ban-object
type Ban struct {
	Reason string `json:"reason"`
	User   *User  `json:"user"`
}

var _ Copier = (*Ban)(nil)
var _ DeepCopier = (*Ban)(nil)

// ------------

// GuildEmbed https://discord.com/developers/docs/resources/guild#guild-embed-object
type GuildWidget struct {
	Enabled   bool      `json:"enabled"`
	ChannelID Snowflake `json:"channel_id"`
}

var _ Copier = (*GuildWidget)(nil)
var _ DeepCopier = (*GuildWidget)(nil)

// -------

// Integration https://discord.com/developers/docs/resources/guild#integration-object
type Integration struct {
	ID                Snowflake           `json:"id"`
	Name              string              `json:"name"`
	Type              string              `json:"type"`
	Enabled           bool                `json:"enabled"`
	Syncing           bool                `json:"syncing"`
	RoleID            Snowflake           `json:"role_id"`
	ExpireBehavior    int                 `json:"expire_behavior"`
	ExpireGracePeriod int                 `json:"expire_grace_period"`
	User              *User               `json:"user"`
	Account           *IntegrationAccount `json:"account"`
}

var _ Copier = (*Integration)(nil)
var _ DeepCopier = (*Integration)(nil)

// IntegrationAccount https://discord.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id"`   // id of the account
	Name string `json:"name"` // name of the account
}

var _ Copier = (*IntegrationAccount)(nil)
var _ DeepCopier = (*IntegrationAccount)(nil)

// -------

// Member https://discord.com/developers/docs/resources/guild#guild-member-object
type Member struct {
	GuildID      Snowflake   `json:"guild_id,omitempty"`
	User         *User       `json:"user"`
	Nick         string      `json:"nick,omitempty"`
	Roles        []Snowflake `json:"roles"`
	JoinedAt     Time        `json:"joined_at,omitempty"`
	PremiumSince Time        `json:"premium_since,omitempty"`
	Deaf         bool        `json:"deaf"`
	Mute         bool        `json:"mute"`
	Pending      bool        `json:"pending"`

	// custom
	UserID Snowflake `json:"-"`
}

var _ Reseter = (*Member)(nil)
var _ fmt.Stringer = (*Member)(nil)
var _ internalUpdater = (*Member)(nil)
var _ Mentioner = (*Member)(nil)
var _ Copier = (*Member)(nil)
var _ DeepCopier = (*Member)(nil)

func (m *Member) updateInternals() {
	if m.User != nil {
		m.UserID = m.User.ID
	}
}

func (m *Member) String() string {
	username := m.Nick
	if m.User != nil {
		username = m.User.Username
	}
	id := m.UserID
	if m.UserID.IsZero() && m.User != nil {
		id = m.User.ID
	}
	return "member{user:" + username + ", nick:" + m.Nick + ", ID:" + id.String() + "}"
}

type GuildQueryBuilderCaller interface {
	Guild(id Snowflake) GuildQueryBuilder
}

func (m *Member) UpdateNick(ctx context.Context, client GuildQueryBuilderCaller, nickname string) error {
	builder := client.Guild(m.GuildID).Member(m.UserID).WithContext(ctx).UpdateBuilder()
	return builder.
		SetNick(nickname).
		Execute()
}

// GetPermissions populates a uint64 with all the permission flags
func (m *Member) GetPermissions(ctx context.Context, s GuildQueryBuilderCaller) (permissions PermissionBit, err error) {
	// TODO: Don't deep copy channels for this in the future!
	roles, err := s.Guild(m.GuildID).WithContext(ctx).GetRoles()
	if err != nil {
		return 0, err
	}

	unprocessedRoles := len(m.Roles)
	for _, roleInfo := range roles {
		for _, roleId := range m.Roles {
			if roleInfo.ID == roleId {
				permissions |= (PermissionBit)(roleInfo.Permissions)
				unprocessedRoles--
				break
			}
		}

		if unprocessedRoles == 0 {
			break
		}
	}

	return permissions, nil
}

// GetUser tries to ensure that you get a user object and not a nil. The user can be nil if the guild
// was fetched from the cache.
func (m *Member) GetUser(ctx context.Context, session Session) (usr *User, err error) {
	if m.User != nil {
		return m.User, nil
	}

	return session.User(m.UserID).WithContext(ctx).Get()
}

// Mention creates a string which is parsed into a member mention on Discord GUI's
func (m *Member) Mention() string {
	var id Snowflake
	if !m.UserID.IsZero() {
		id = m.UserID
	} else if m.User != nil {
		id = m.User.ID
	}

	return "<@!" + id.String() + ">"
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

// CreateGuild ...
// https://discord.com/developers/docs/resources/guild#create-guild-json-params
// example partial channel object:
// {
//    "name": "naming-things-is-hard",
//    "type": 0
// }
type CreateGuild struct {
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
//  Discord documentation   https://discord.com/developers/docs/resources/guild#create-guild
//  Reviewed                2018-08-16
//  Comment                 This endpoint. can be used only by bots in less than 10 Guilds. Creating channel
//                          categories from this endpoint. is not supported.
//							The params argument is optional.
func (c clientQueryBuilder) CreateGuild(guildName string, params *CreateGuild) (ret *Guild, err error) {
	// TODO: check if bot
	// TODO-2: is bot in less than 10 Guilds?

	if guildName == "" {
		return nil, MissingGuildNameErr
	}
	if l := len(guildName); !(2 <= l && l <= 100) {
		return nil, fmt.Errorf("guild name must be 2 or more characters and no more than 100 characters: %w", IllegalValueErr)
	}

	if params == nil {
		params = &CreateGuild{}
	}
	params.Name = guildName

	r := c.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ctx:         c.ctx,
		Endpoint:    endpoint.Guilds(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, c.flags)
	r.factory = func() interface{} {
		return &Guild{}
	}

	return getGuild(r.Execute)
}

// GuildQueryBuilder defines the exposed functions from the guild query builder.
type GuildQueryBuilder interface {
	WithContext(ctx context.Context) GuildQueryBuilder
	WithFlags(flags ...Flag) GuildQueryBuilder

	// Get
	// TODO: Add more guild attribute things. Waiting for caching changes before then.
	Get() (guild *Guild, err error)
	Update(params *UpdateGuild) (*Guild, error)
	Delete() error

	// Leave leaves the given guild
	Leave() error

	// GetChannels
	// TODO: For GetChannels, it might sense to have the option for a function to filter before each channel ends up deep copied.
	// TODO-2: This could be much more performant in guilds with a large number of channels.
	GetChannels() ([]*Channel, error)

	// GetMembers
	// TODO: For GetMembers, it might sense to have the option for a function to filter before each member ends up deep copied.
	// TODO-2: This could be much more performant in larger guilds where this is needed.
	GetMembers(params *GetMembers) ([]*Member, error)
	UpdateBuilder() UpdateGuildBuilder

	CreateChannel(name string, params *CreateGuildChannel) (*Channel, error)
	UpdateChannelPositions(params []UpdateGuildChannelPositions) error
	CreateMember(userID Snowflake, accessToken string, params *AddGuildMember) (*Member, error)
	Member(userID Snowflake) GuildMemberQueryBuilder

	// Deprecated: use DisconnectVoiceParticipant
	KickVoiceParticipant(userID Snowflake) error

	DisconnectVoiceParticipant(userID Snowflake) error
	SetCurrentUserNick(nick string) (newNick string, err error)
	GetBans() ([]*Ban, error)
	GetBan(userID Snowflake) (*Ban, error)
	UnbanUser(userID Snowflake, reason string) error

	// GetRoles
	// TODO: For GetRoles, it might sense to have the option for a function to filter before each role ends up deep copied.
	// TODO-2: This could be much more performant in larger guilds where this is needed.
	// TODO-3: Add GetRole.
	GetRoles() ([]*Role, error)
	UpdateRolePositions(params []UpdateGuildRolePositions) ([]*Role, error)
	CreateRole(params *CreateGuildRole) (*Role, error)
	Role(roleID Snowflake) GuildRoleQueryBuilder

	EstimatePruneMembersCount(days int) (estimate int, err error)
	PruneMembers(days int, reason string) error
	GetVoiceRegions() ([]*VoiceRegion, error)
	GetInvites() ([]*Invite, error)

	GetIntegrations() ([]*Integration, error)
	CreateIntegration(params *CreateGuildIntegration) error
	UpdateIntegration(integrationID Snowflake, params *UpdateGuildIntegration) error
	DeleteIntegration(integrationID Snowflake) error
	SyncIntegration(integrationID Snowflake) error

	GetWidget() (*GuildWidget, error)
	UpdateWidget(params *UpdateGuildWidget) (*GuildWidget, error)
	GetVanityURL() (*PartialInvite, error)
	GetAuditLogs(logs *GetAuditLogs) (*AuditLog, error)

	VoiceChannel(channelID Snowflake) VoiceChannelQueryBuilder

	// GetEmojis
	// TODO: For GetEmojis, it might sense to have the option for a function to filter before each emoji ends up deep copied.
	// TODO-2: This could be much more performant in guilds with a large number of channels.
	GetEmojis() ([]*Emoji, error)
	CreateEmoji(params *CreateGuildEmoji) (*Emoji, error)
	Emoji(emojiID Snowflake) GuildEmojiQueryBuilder

	GetWebhooks() (ret []*Webhook, err error)

	// GetActiveThreads Returns all active threads in the guild, including public and private threads. Threads are ordered
	// by their id, in descending order.
	GetActiveThreads() (*ActiveGuildThreads, error)

	// Deprecated: use UpdateEmbed
	UpdateEmbedBuilder() UpdateGuildEmbedBuilder
	// Deprecated: use GetWidget
	GetEmbed() (*GuildEmbed, error)
}

// Guild is used to create a guild query builder.
func (c clientQueryBuilder) Guild(id Snowflake) GuildQueryBuilder {
	return &guildQueryBuilder{client: c.client, gid: id}
}

// The default guild query builder.
type guildQueryBuilder struct {
	ctx    context.Context
	flags  Flag
	client *Client
	gid    Snowflake
}

func (g *guildQueryBuilder) validate() error {
	if g.client == nil {
		return MissingClientInstanceErr
	}
	if g.gid.IsZero() {
		return MissingGuildIDErr
	}
	return nil
}

func (g guildQueryBuilder) WithContext(ctx context.Context) GuildQueryBuilder {
	g.ctx = ctx
	return &g
}

func (g guildQueryBuilder) WithFlags(flags ...Flag) GuildQueryBuilder {
	g.flags = mergeFlags(flags)
	return &g
}

// Get is used to get the Guild struct containing all information from it.
// Note that it's significantly quicker in most instances where you have the cache enabled (as is by default) to get the individual parts you need.
func (g guildQueryBuilder) Get() (guild *Guild, err error) {
	if !ignoreCache(g.flags) {
		if guild, _ = g.client.cache.GetGuild(g.gid); guild != nil {
			return guild, nil
		}
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Guild(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		return &Guild{}
	}

	return getGuild(r.Execute)
}

// Update update a guild
func (g guildQueryBuilder) Update(params *UpdateGuild) (*Guild, error) {
	if params == nil {
		return nil, MissingRESTParamsErr
	}
	if err := g.validate(); err != nil {
		return nil, err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.Guild(g.gid),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
		Reason:      params.AuditLogReason,
	}, g.flags)
	r.factory = func() interface{} {
		return &Guild{}
	}

	return getGuild(r.Execute)
}

type UpdateGuild struct {
	Name                        *string                        `json:"name,omitempty"`
	Region                      *string                        `json:"region,omitempty"`
	VerificationLvl             *VerificationLvl               `json:"verification_lvl,omitempty"`
	DefaultMessageNotifications *DefaultMessageNotificationLvl `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *ExplicitContentFilterLvl      `json:"explicit_content_filter,omitempty"`
	AFKChannelID                *Snowflake                     `json:"afk_channel_id,omitempty"`
	Icon                        *string                        `json:"icon,omitempty"`
	OwnerID                     *Snowflake                     `json:"owner_id,omitempty"`
	Splash                      *string                        `json:"splash,omitempty"`
	DiscoverySplash             *string                        `json:"discovery_splash,omitempty"`
	Banner                      *string                        `json:"banner,omitempty"`
	SystemChannelID             *Snowflake                     `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *uint                          `json:"system_channel_flags,omitempty"`
	RulesChannelID              *Snowflake                     `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *Snowflake                     `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *string                        `json:"preferred_locale,omitempty"`
	Features                    *[]string                      `json:"features,omitempty"`
	Description                 *string                        `json:"description,omitempty"`

	AuditLogReason string `json:"-"`
}

// Delete is used to delete a guild.
func (g guildQueryBuilder) Delete() error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.Guild(g.gid),
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// Leave https://discord.com/developers/docs/resources/user#leave-guild
func (g guildQueryBuilder) Leave() error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.UserMeGuild(g.gid),
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// GetChannels is used to get a guilds channels.
func (g guildQueryBuilder) GetChannels() ([]*Channel, error) {
	if channels, _ := g.client.cache.GetGuildChannels(g.gid); channels != nil {
		return channels, nil
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildChannels(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Channel, 0)
		return &tmp
	}

	return getChannels(r.Execute)
}

// CreateChannel Create a new channel object for the guild. Requires the 'MANAGE_CHANNELS' permission.
// Returns the new channel object on success. Fires a Channel Create Gateway event.
func (g guildQueryBuilder) CreateChannel(name string, params *CreateGuildChannel) (*Channel, error) {
	if name == "" && (params == nil || params.Name == "") {
		return nil, MissingChannelNameErr
	}
	if l := len(name); !(2 <= l && l <= 100) {
		return nil, fmt.Errorf("channel name must be 2 or more characters and no more than 100 characters: %w", IllegalValueErr)
	}

	if params == nil {
		params = &CreateGuildChannel{}
	}
	if name != "" && params.Name == "" {
		params.Name = name
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildChannels(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      params.Reason,
	}, g.flags)
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

// UpdateChannelPositions Modify the positions of a set of channel objects for the guild.
// Requires 'MANAGE_CHANNELS' permission. Returns a 204 empty response on success. Fires multiple Channel Update
// Gateway events.
func (g guildQueryBuilder) UpdateChannelPositions(params []UpdateGuildChannelPositions) error {
	var reason string
	for i := range params {
		if params[i].Reason != "" {
			reason = params[i].Reason
			break
		}
	}
	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildChannels(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      reason,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// GetMembers uses the GetGuildMembers endpoint iteratively until your query params are met.
func (g guildQueryBuilder) GetMembers(params *GetMembers) ([]*Member, error) {
	const QueryLimit uint32 = 1000

	if params == nil {
		params = &GetMembers{
			Limit: math.MaxUint32,
		}
	}

	highestSnowflake := func(ms []*Member) (highest Snowflake) {
		for i := range ms {
			if ms[i].User != nil && ms[i].User.ID > highest {
				highest = ms[i].User.ID
			}
		}
		return highest
	}

	p := getGuildMembers{
		After: params.After,
	}
	if params.Limit == 0 || params.Limit > QueryLimit {
		p.Limit = int(QueryLimit)
	} else {
		p.Limit = int(params.Limit)
	}

	members := make([]*Member, 0)
	var ms []*Member
	var err error
	for {
		ms, err = g.getGuildMembers(&p)
		if ms != nil {
			members = append(members, ms...)
		}
		if err != nil {
			return members, err
		}

		// stop if we're on the last page/block of members
		if len(ms) < int(QueryLimit) {
			break
		}

		// set limit such that we don't retrieve redundant members
		max := params.Limit << 1
		max = max >> 1
		lim := int(max) - len(members)
		if lim < int(QueryLimit) {
			if lim <= 0 {
				// should never be less than 0
				break
			}
			p.Limit = lim
		}

		params.After = highestSnowflake(ms)
	}

	return members, err
}

// CreateMember Adds a user to the guild, provided you have a valid oauth2 access token for the user with
// the Guilds.join scope. Returns a 201 Created with the guild member as the body, or 204 No Content if the user is
// already a member of the guild. Fires a Guild Member Add Gateway event. Requires the bot to have the
// CREATE_INSTANT_INVITE permission.
func (g guildQueryBuilder) CreateMember(userID Snowflake, accessToken string, params *AddGuildMember) (*Member, error) {
	if accessToken == "" && (params == nil || params.AccessToken == "") {
		return nil, errors.New("access token is required")
	}

	if params == nil {
		params = &AddGuildMember{}
	}
	if accessToken != "" && params.AccessToken == "" {
		params.AccessToken = accessToken
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPut,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildMember(g.gid, userID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, g.flags)
	r.factory = func() interface{} {
		return &Member{
			GuildID: g.gid,
			UserID:  userID,
		}
	}

	// TODO: update guild cache
	var member *Member
	var err error
	if member, err = getMember(r.Execute); err != nil {
		if errRest, ok := err.(*httd.ErrREST); ok && errRest.HTTPCode == http.StatusNoContent {
			errRest.Msg = "member{" + userID.String() + "} is already in Guild{" + g.gid.String() + "}"
		}
	}

	if member != nil {
		member.GuildID = g.gid
	}

	return member, err
}

// SetCurrentUserNick Modifies the nickname of the current user in a guild. Returns a 200
// with the nickname on success. Fires a Guild Member Update Gateway event.
func (g guildQueryBuilder) SetCurrentUserNick(nick string) (newNick string, err error) {
	params := &updateCurrentUserNick{
		Nick: nick,
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildMembersMeNick(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, g.flags)
	r.factory = func() interface{} {
		return &nickNameResponse{}
	}

	return getNickName(r.Execute)
}

// GetBans returns an array of ban objects for the Users banned from this guild. Requires the 'BAN_MEMBERS' permission.
func (g guildQueryBuilder) GetBans() ([]*Ban, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildBans(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Ban, 0)
		return &tmp
	}

	var vs interface{}
	var err error
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if cons, ok := vs.(*[]*Ban); ok {
		return *cons, nil
	}
	return nil, errors.New("unable to cast guild slice")
}

// GetBan Returns a ban object for the given user or a 404 not found if the ban cannot be found.
// Requires the 'BAN_MEMBERS' permission.
func (g guildQueryBuilder) GetBan(userID Snowflake) (*Ban, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildBan(g.gid, userID),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		return &Ban{User: g.client.pool.user.Get().(*User)}
	}

	return getBan(r.Execute)
}

// UnbanUser Remove the ban for a user. Requires the 'BAN_MEMBERS' permissions.
// Returns a 204 empty response on success. Fires a Guild Ban Remove Gateway event.
func (g guildQueryBuilder) UnbanUser(userID Snowflake, reason string) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.GuildBan(g.gid, userID),
		Reason:   reason,
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// GetRoles Returns a list of role objects for the guild.
func (g guildQueryBuilder) GetRoles() ([]*Role, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: "/guilds/" + g.gid.String() + "/roles",
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Role, 0)
		return &tmp
	}

	return getRoles(r.Execute)
}

// CreateGuildRole ...
// https://discord.com/developers/docs/resources/guild#create-guild-role-json-params
type CreateGuildRole struct {
	Name        string `json:"name,omitempty"`
	Permissions uint64 `json:"permissions,omitempty"`
	Color       uint   `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// CreateRole Create a new role for the guild. Requires the 'MANAGE_ROLES' permission.
// Returns the new role object on success. Fires a Guild Role Create Gateway event.
func (g guildQueryBuilder) CreateRole(params *CreateGuildRole) (*Role, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildRoles(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      params.Reason,
	}, g.flags)
	r.factory = func() interface{} {
		return &Role{}
	}

	return getRole(r.Execute)
}

// UpdateRolePositions Modify the positions of a set of role objects for the guild.
// Requires the 'MANAGE_ROLES' permission. Returns a list of all of the guild's role objects on success.
// Fires multiple Guild Role Update Gateway events.
func (g guildQueryBuilder) UpdateRolePositions(params []UpdateGuildRolePositions) ([]*Role, error) {
	var reason string
	for i := range params {
		if params[i].Reason != "" {
			reason = params[i].Reason
			break
		}
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildRoles(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      reason,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Role, 0)
		return &tmp
	}

	return getRoles(r.Execute)
}

// EstimatePruneMembersCount Returns an object with one 'pruned' key indicating the number of members that would be
// removed in a prune operation. Requires the 'KICK_MEMBERS' permission.
func (g guildQueryBuilder) EstimatePruneMembersCount(days int) (estimate int, err error) {
	if g.gid.IsZero() {
		return 0, MissingGuildIDErr
	}
	params := pruneMembers{Days: days}
	if err = params.FindErrors(); err != nil {
		return 0, err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildPrune(g.gid) + params.URLQueryString(),
		Ctx:      g.ctx,
	}, g.flags)
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

// PruneMembers Kicks members from N day back. Requires the 'KICK_MEMBERS' permission.
// The estimate of kicked people is not returned. Use EstimatePruneMembersCount before calling PruneMembers
// if you need it. Fires multiple Guild Member Remove Gateway events.
func (g guildQueryBuilder) PruneMembers(days int, reason string) (err error) {
	params := pruneMembers{Days: days}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodPost,
		Endpoint: endpoint.GuildPrune(g.gid) + params.URLQueryString(),
		Ctx:      g.ctx,
		Reason:   reason,
	}, g.flags)

	_, err = r.Execute()
	return err
}

// GetVoiceRegions Returns a list of voice region objects for the guild. Unlike the similar /voice route,
// this returns VIP servers when the guild is VIP-enabled.
func (g guildQueryBuilder) GetVoiceRegions() ([]*VoiceRegion, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildRegions(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*VoiceRegion, 0)
		return &tmp
	}

	return getVoiceRegions(r.Execute)
}

// GetInvites Returns a list of invite objects (with invite metadata) for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (g guildQueryBuilder) GetInvites() ([]*Invite, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildInvites(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Invite, 0)
		return &tmp
	}

	return getInvites(r.Execute)
}

// GetIntegrations Returns a list of integration objects for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (g guildQueryBuilder) GetIntegrations() ([]*Integration, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildIntegrations(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Integration, 0)
		return &tmp
	}

	return getIntegrations(r.Execute)
}

// CreateIntegration attaches an integration object from the current user to the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
func (g guildQueryBuilder) CreateIntegration(params *CreateGuildIntegration) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildIntegrations(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// UpdateIntegration Modify the behavior and settings of a integration object for the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
func (g guildQueryBuilder) UpdateIntegration(integrationID Snowflake, params *UpdateGuildIntegration) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildIntegration(g.gid, integrationID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// DeleteIntegration Delete the attached integration object for the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
func (g guildQueryBuilder) DeleteIntegration(integrationID Snowflake) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodDelete,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildIntegration(g.gid, integrationID),
		ContentType: httd.ContentTypeJSON,
	}, g.flags)

	_, err := r.Execute()
	return err
}

// SyncIntegration Sync an integration. Requires the 'MANAGE_GUILD' permission.
// Returns a 204 empty response on success.
func (g guildQueryBuilder) SyncIntegration(integrationID Snowflake) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   http.MethodPost,
		Endpoint: endpoint.GuildIntegrationSync(g.gid, integrationID),
		Ctx:      g.ctx,
	}, g.flags)

	_, err := r.Execute()
	return err
}

func (g guildQueryBuilder) GetEmbed() (*GuildEmbed, error) {
	return g.GetWidget()
}

func (g guildQueryBuilder) GetWidget() (*GuildWidget, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmbed(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		return &GuildWidget{}
	}

	return getGuildWidget(r.Execute)
}

func (g guildQueryBuilder) UpdateWidget(params *UpdateGuildWidget) (*GuildWidget, error) {
	if params == nil {
		return nil, MissingRESTParamsErr
	}
	if err := g.validate(); err != nil {
		return nil, err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.Guild(g.gid),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
		Reason:      params.AuditLogReason,
	}, g.flags)
	r.factory = func() interface{} {
		return &GuildWidget{}
	}

	return getGuildWidget(r.Execute)
}

type UpdateGuildWidget struct {
	Enabled   *bool      `json:"enabled,omitempty"`
	ChannelID *Snowflake `json:"channel_id,omitempty"`

	AuditLogReason string `json:"-"`
}

// GetVanityURL Returns a partial invite object for Guilds with that feature enabled.
// Requires the 'MANAGE_GUILD' permission.
func (g guildQueryBuilder) GetVanityURL() (*PartialInvite, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildVanityURL(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		return &PartialInvite{}
	}

	return getPartialInvite(r.Execute)
}

// GetAuditLogs Returns an audit log object for the guild. Requires the 'VIEW_AUDIT_LOG' permission.
// Note that this request will _always_ send a REST request, regardless of you calling IgnoreCache or not.
func (g guildQueryBuilder) GetAuditLogs(params *GetAuditLogs) (*AuditLog, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildAuditLogs(g.gid) + params.URLQueryString(),
		Method:   http.MethodGet,
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = auditLogFactory

	logs, err := r.Execute()
	if err != nil {
		return nil, err
	}

	return logs.(*AuditLog), nil
}

type GetAuditLogs struct {
	UserID     Snowflake `urlparam:"user_id"`
	ActionType int       `urlparam:"action_type"`
	Before     Snowflake `urlparam:"before,omitempty"`
	Limit      int       `urlparam:"limit,omitempty"`
}

var _ URLQueryStringer = (*GetAuditLogs)(nil)

// GetEmojis Returns a list of emoji objects for the given guild.
func (g guildQueryBuilder) GetEmojis() ([]*Emoji, error) {
	if emojis, _ := g.client.cache.GetGuildEmojis(g.gid); emojis != nil {
		return emojis, nil
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmojis(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Emoji, 0)
		return &tmp
	}

	var vs interface{}
	var err error
	if vs, err = r.Execute(); err != nil {
		return nil, err
	}

	if ems, ok := vs.(*[]*Emoji); ok {
		return *ems, nil
	}
	return vs.([]*Emoji), nil
}

// CreateGuildEmoji JSON params for func CreateGuildEmoji
type CreateGuildEmoji struct {
	Name  string      `json:"name"`  // required
	Image string      `json:"image"` // required
	Roles []Snowflake `json:"roles"` // optional

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// CreateEmoji Create a new emoji for the guild. Requires the 'MANAGE_EMOJIS' permission.
// Returns the new emoji object on success. Fires a Guild Emojis Update Gateway event.
func (g guildQueryBuilder) CreateEmoji(params *CreateGuildEmoji) (*Emoji, error) {
	if g.gid.IsZero() {
		return nil, MissingGuildIDErr
	}

	if params == nil {
		return nil, errors.New("params object can not be nil")
	}
	if !validEmojiName(params.Name) {
		return nil, fmt.Errorf("invalid emoji name: %w", IllegalValueErr)
	}
	if !validAvatarPrefix(params.Image) {
		return nil, errors.New("image string must be base64 encoded with base64 prefix")
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildEmojis(g.gid),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
		Reason:      params.Reason,
	}, g.flags)
	r.pool = g.client.pool.emoji
	r.factory = func() interface{} {
		return &Emoji{}
	}

	return getEmoji(r.Execute)
}

// DisconnectVoiceParticipant is used to kick someone from voice.
func (g guildQueryBuilder) DisconnectVoiceParticipant(userID Snowflake) error {
	builder := g.Member(userID).WithContext(g.ctx).UpdateBuilder()
	return builder.
		KickFromVoice().
		Execute()
}

func (g guildQueryBuilder) KickVoiceParticipant(userID Snowflake) error {
	return g.DisconnectVoiceParticipant(userID)
}

// GetWebhooks Returns a list of guild webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
func (g guildQueryBuilder) GetWebhooks() (ret []*Webhook, err error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildWebhooks(g.gid),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Webhook, 0)
		return &tmp
	}

	return getWebhooks(r.Execute)
}

// CreateGuildChannel https://discord.com/developers/docs/resources/guild#create-guild-channel-json-params
type CreateGuildChannel struct {
	Name                 string                `json:"name"` // required
	Type                 ChannelType           `json:"type,omitempty"`
	Topic                string                `json:"topic,omitempty"`
	Bitrate              uint                  `json:"bitrate,omitempty"`
	UserLimit            uint                  `json:"user_limit,omitempty"`
	RateLimitPerUser     uint                  `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             Snowflake             `json:"parent_id,omitempty"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
	Position             int                   `json:"position"` // can not omitempty in case position is 0

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// UpdateGuildChannelPositions
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions-json-params
type UpdateGuildChannelPositions struct {
	ID       Snowflake `json:"id"`
	Position int       `json:"position"`

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	// just reuse the string. Go will optimize it to point to the same memory anyways
	// TODO: improve this?
	Reason string `json:"-"`
}

func NewUpdateGuildRolePositions(rs []*Role) (p []*UpdateGuildRolePositions) {
	p = make([]*UpdateGuildRolePositions, 0, len(rs))
	for i := range rs {
		p = append(p, &UpdateGuildRolePositions{
			ID:       rs[i].ID,
			Position: rs[i].Position,
		})
	}

	return p
}

// UpdateGuildRolePositions
// https://discord.com/developers/docs/resources/guild#modify-guild-role-positions-json-params
type UpdateGuildRolePositions struct {
	ID       Snowflake `json:"id"`
	Position int       `json:"position"`

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

type getGuildMembers struct {
	After Snowflake `urlparam:"after,omitempty"`
	Limit int       `urlparam:"limit,omitempty"` // 1 is default. even if 0 is supplied.
}

var _ URLQueryStringer = (*getGuildMembers)(nil)

func (g *getGuildMembers) FindErrors() error {
	if g.Limit > 1000 || g.Limit < 1 {
		return errors.New("limit value should be less than or equal to 1000, and 1 or more")
	}
	return nil
}

// GetGuildMembers [REST] Returns a list of guild member objects that are members of the guild. The `after` param
// refers to the highest snowflake.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/members
//  Discord documentation   https://discord.com/developers/docs/resources/guild#get-guild-members
//  Reviewed                2018-08-17
//  Comment                 All parameters to this endpoint. are optional
//  Comment#2               "List Guild Members"
//  Comment#3               https://discord.com/developers/docs/resources/guild#list-guild-members-query-string-params
func (g guildQueryBuilder) getGuildMembers(params *getGuildMembers) (ret []*Member, err error) {
	if params == nil {
		params = &getGuildMembers{}
	}
	if err = params.FindErrors(); err != nil {
		return nil, err
	}

	if !ignoreCache(g.flags) {
		p := &GetMembers{After: params.After, Limit: uint32(params.Limit)}
		members, err := g.client.cache.GetMembers(g.gid, p)
		if err == nil && len(members) > 0 {
			return members, nil
		}
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildMembers(g.gid) + params.URLQueryString(),
		Ctx:      g.ctx,
	}, g.flags)
	r.factory = func() interface{} {
		tmp := make([]*Member, 0)
		return &tmp
	}

	return getMembers(r.Execute)
}

// GetMembers if Limit is 0, every member is fetched. This does not follow the Discord API where a 0
// is converted into a 1. 0 = every member. The rest is exactly the same, you should be able to do everything
// the Discord docs says with the addition that you can bypass a limit of 1,000.
//
// If you specify a limit of +1,000 Disgord will run N requests until that amount is met, or until you run
// out of members to fetch.
type GetMembers struct {
	After Snowflake `urlparam:"after,omitempty"`
	Limit uint32    `urlparam:"limit,omitempty"` // 0 will fetch everyone
}

// AddGuildMember ...
// https://discord.com/developers/docs/resources/guild#add-guild-member-json-params
type AddGuildMember struct {
	AccessToken string      `json:"access_token"` // required
	Nick        string      `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Mute        bool        `json:"mute,omitempty"`
	Deaf        bool        `json:"deaf,omitempty"`
}

// BanMember ...
// https://discord.com/developers/docs/resources/guild#create-guild-ban-query-string-params
type BanMember struct {
	DeleteMessageDays int    `urlparam:"delete_message_days,omitempty"` // number of days to delete messages for (0-7)
	Reason            string `urlparam:"reason,omitempty"`              // reason for being banned
}

var _ URLQueryStringer = (*BanMember)(nil)

func (b *BanMember) FindErrors() error {
	if !(0 <= b.DeleteMessageDays && b.DeleteMessageDays <= 7) {
		return errors.New("DeleteMessageDays must be a value in the range of [0, 7], got " + strconv.Itoa(b.DeleteMessageDays))
	}
	return nil
}

// PruneMembers will delete members, this is the same as kicking.
// https://discord.com/developers/docs/resources/guild#get-guild-prune-count-query-string-params
type pruneMembers struct {
	// Days number of days to count prune for (1 or more)
	Days int `urlparam:"days"`

	// ComputePruneCount whether 'pruned' is returned, discouraged for large Guilds
	ComputePruneCount bool `urlparam:"compute_prune_count"`
}

var _ URLQueryStringer = (*pruneMembers)(nil)

func (d *pruneMembers) FindErrors() (err error) {
	if d.Days < 1 {
		err = errors.New("days must be at least 1, got " + strconv.Itoa(d.Days))
	}
	return
}

// GuildPruneCount ...
type guildPruneCount struct {
	Pruned int `json:"pruned"`
}

// CreateGuildIntegration ...
// https://discord.com/developers/docs/resources/guild#create-guild-integration-json-params
type CreateGuildIntegration struct {
	Type string    `json:"type"`
	ID   Snowflake `json:"id"`
}

// UpdateGuildIntegration ...
// https://discord.com/developers/docs/resources/guild#modify-guild-integration-json-params
// TODO: currently unsure which are required/optional params
type UpdateGuildIntegration struct {
	ExpireBehavior    int  `json:"expire_behavior"`
	ExpireGracePeriod int  `json:"expire_grace_period"`
	EnableEmoticons   bool `json:"enable_emoticons"`
}

// updateCurrentUserNick ...
// https://discord.com/developers/docs/resources/guild#modify-guild-member-json-params
type updateCurrentUserNick struct {
	Nick string `json:"nick"` // :CHANGE_NICKNAME
}

type nickNameResponse struct {
	Nickname string `json:"nickname"`
}

// ActiveGuildThreads https://discord.com/developers/docs/resources/guild#list-active-threads-response-body
type ActiveGuildThreads struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
}

// GetActiveThreads https://discord.com/developers/docs/resources/guild#list-active-threads
func (g guildQueryBuilder) GetActiveThreads() (*ActiveGuildThreads, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      http.MethodGet,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildThreadsActive(g.gid),
		ContentType: httd.ContentTypeJSON,
	}, g.flags)
	r.factory = func() interface{} {
		return &ActiveGuildThreads{
			Threads: make([]*Channel, 0),
			Members: make([]*ThreadMember, 0),
		}
	}

	return getActiveGuildThreads(r.Execute)
}
