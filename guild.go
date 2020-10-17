package disgord

import (
	"context"
	"errors"
	"fmt"
	"github.com/andersfylling/disgord/json"
	"math"
	"net/http"
	"sort"
	"strconv"

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

// NewGuild ...
func NewGuild() *Guild {
	return &Guild{
		Roles:       []*Role{},
		Emojis:      []*Emoji{},
		Features:    []string{},
		VoiceStates: []*VoiceState{},
		Members:     []*Member{},
		Channels:    []*Channel{},
		Presences:   []*UserPresence{},
	}
}

// NewPartialGuild ...
func NewPartialGuild(ID Snowflake) (guild *Guild) {
	guild = NewGuild()
	guild.ID = ID
	guild.Unavailable = true

	return
}

// NewGuildFromUnavailable converts a unavailable guild object into a normal Guild object
func NewGuildFromUnavailable(gu *GuildUnavailable) *Guild {
	return NewPartialGuild(gu.ID)
}

// NewGuildUnavailable ...
func NewGuildUnavailable(ID Snowflake) *GuildUnavailable {
	gu := &GuildUnavailable{
		ID:          ID,
		Unavailable: true,
	}

	return gu
}

// GuildUnavailable is a partial Guild object.
type GuildUnavailable struct {
	ID          Snowflake `json:"id"`
	Unavailable bool      `json:"unavailable"` // ?*|
}

//type GuildInterface interface {
//	Channel(ID Snowflake)
//}

// if loading is deactivated, then check state, then do a request.
// if loading is activated, check state only.
// type Members interface {
// 	Member(userID snowflake.Snowflake) *Member
// 	MembersWithName( /*username*/ name string) map[snowflake.Snowflake]*Member
// 	MemberByUsername( /*username#discriminator*/ username string) *Member
// 	MemberByAlias(alias string) *Member
// 	EverythingInMemory() bool
// }

// PartialGuild see Guild
type PartialGuild = Guild // TODO: find the actual data struct for partial guild

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
	EmbedEnabled                bool                          `json:"embed_enabled,omit_empty"`
	EmbedChannelID              Snowflake                     `json:"embed_channel_id,omit_empty"`
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

	// JoinedAt must be a pointer, as we can't hide non-nil structs
	JoinedAt    *Time           `json:"joined_at,omitempty"`    // ?*|
	Large       bool            `json:"large,omitempty"`        // ?*|
	Unavailable bool            `json:"unavailable"`            // ?*| omitempty?
	MemberCount uint            `json:"member_count,omitempty"` // ?*|
	VoiceStates []*VoiceState   `json:"voice_states,omitempty"` // ?*|
	Members     []*Member       `json:"members,omitempty"`      // ?*|
	Channels    []*Channel      `json:"channels,omitempty"`     // ?*|
	Presences   []*UserPresence `json:"presences,omitempty"`    // ?*|

	//highestSnowflakeAmongMembers Snowflake
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

// MarshalJSON see interface json.Marshaler
// TODO: fix copying of mutex lock
func (g *Guild) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	var err error
	if g.Unavailable {
		guildUnavailable := NewGuildUnavailable(g.ID)
		jsonData, err = json.Marshal(guildUnavailable)
		if err != nil {
			return []byte(""), nil
		}
	} else {
		jsonData, err = json.Marshal(Guild(*g))
		if err != nil {
			return []byte(""), nil
		}
	}

	return jsonData, nil
}

// sortChannels Only while in lock
func (g *Guild) sortChannels() {
	sort.Slice(g.Channels, func(i, j int) bool {
		return g.Channels[i].ID < g.Channels[j].ID
	})
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
		return 0, errors.New("unable to decide which channel to create invite for")
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
		if member.User.ID == id {
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
func (g *Guild) Role(id Snowflake) (role *Role, err error) {
	for _, role = range g.Roles {
		if role.ID == id {
			return
		}
	}

	err = errors.New("role not found in guild")
	return
}

// TODO
//func (g *Guild) UpdateRole(r *Role) {
//	for _, role := range g.Roles {
//		if role.ID == r.ID {
//			*role = *r
//			break
//		}
//	}
//}

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

// TODO
// func (g *Guild) UpdatePresence(p *UserPresence) {
// 	g.RLock()
// 	index := -1
// 	for i, presence := range g.Presences {
// 		if presence.User.ID == p.User.ID {
// 			index = i
// 			break
// 		}
// 	}
// 	g.RUnlock()
//
// 	if index != -1 {
// 		// update
// 		return
// 	}
//
// 	// otherwise add
// 	g.AcquireLock()
// 	g.Presences = append(g.Presences, p) // TODO: update the user pointer?
// 	g.Unlock()
// }

// Clear all the pointers
// func (g *Guild) Clear() {
// 	g.AcquireLock() // what if another process tries to read this, but awais while locked for clearing?
// 	defer g.Unlock()
//
// 	//g.Icon = nil // should this be cleared?
// 	//g.Splash = nil // should this be cleared?
//
// 	for _, r := range g.Roles {
// 		r.Clear()
// 		r = nil
// 	}
// 	g.Roles = nil
//
// 	for _, e := range g.Emojis {
// 		e.Clear()
// 		e = nil
// 	}
// 	g.Emojis = nil
//
// 	for _, vst := range g.VoiceStates {
// 		vst.Clear()
// 		vst = nil
// 	}
// 	g.VoiceStates = nil
//
// 	var deletedUsers []Snowflake
// 	for _, m := range g.Members {
// 		deletedUsers = append(deletedUsers, m.Clear())
// 		m = nil
// 	}
// 	g.Members = nil
//
// 	for _, c := range g.Channels {
// 		c.Clear()
// 		c = nil
// 	}
// 	g.Channels = nil
//
// 	for _, p := range g.Presences {
// 		p.Clear()
// 		p = nil
// 	}
// 	g.Presences = nil
//
// }

// DeepCopy see interface at struct.go#DeepCopier
func (g *Guild) DeepCopy() (copy interface{}) {
	copy = NewGuild()
	g.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (g *Guild) CopyOverTo(other interface{}) (err error) {
	var guild *Guild
	var valid bool
	if guild, valid = other.(*Guild); !valid {
		err = newErrorUnsupportedType("argument given is not a *Guild type")
		return
	}

	guild.ID = g.ID
	guild.Name = g.Name
	guild.Owner = g.Owner
	guild.OwnerID = g.OwnerID
	guild.Permissions = g.Permissions
	guild.Region = g.Region
	guild.AfkTimeout = g.AfkTimeout
	guild.EmbedEnabled = g.EmbedEnabled
	guild.EmbedChannelID = g.EmbedChannelID
	guild.VerificationLevel = g.VerificationLevel
	guild.DefaultMessageNotifications = g.DefaultMessageNotifications
	guild.ExplicitContentFilter = g.ExplicitContentFilter
	guild.Features = g.Features
	guild.MFALevel = g.MFALevel
	guild.WidgetEnabled = g.WidgetEnabled
	guild.WidgetChannelID = g.WidgetChannelID
	guild.SystemChannelID = g.SystemChannelID
	guild.Large = g.Large
	guild.Unavailable = g.Unavailable
	guild.MemberCount = g.MemberCount
	guild.Splash = g.Splash
	guild.Icon = g.Icon

	// pointers
	if !g.ApplicationID.IsZero() {
		guild.ApplicationID = g.ApplicationID
	}
	if !g.AfkChannelID.IsZero() {
		guild.AfkChannelID = g.AfkChannelID
	}
	if !g.SystemChannelID.IsZero() {
		guild.SystemChannelID = g.SystemChannelID
	}
	if g.JoinedAt != nil {
		joined := *g.JoinedAt
		guild.JoinedAt = &joined
	}

	for _, roleP := range g.Roles {
		if roleP == nil {
			continue
		}
		guild.Roles = append(guild.Roles, roleP.DeepCopy().(*Role))
	}
	for _, emojiP := range g.Emojis {
		if emojiP == nil {
			continue
		}
		guild.Emojis = append(guild.Emojis, emojiP.DeepCopy().(*Emoji))
	}

	for _, vsP := range g.VoiceStates {
		if vsP == nil {
			continue
		}
		guild.VoiceStates = append(guild.VoiceStates, vsP.DeepCopy().(*VoiceState))
	}
	for _, memberP := range g.Members {
		if memberP == nil {
			continue
		}
		guild.Members = append(guild.Members, memberP.DeepCopy().(*Member))
	}
	for _, channelP := range g.Channels {
		if channelP == nil {
			continue
		}
		guild.Channels = append(guild.Channels, channelP.DeepCopy().(*Channel))
	}
	for _, presenceP := range g.Presences {
		if presenceP == nil {
			continue
		}
		guild.Presences = append(guild.Presences, presenceP.DeepCopy().(*UserPresence))
	}

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

// DeepCopy see interface at struct.go#DeepCopier
func (b *Ban) DeepCopy() (copy interface{}) {
	copy = &Ban{}
	b.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (b *Ban) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var ban *Ban
	if ban, ok = other.(*Ban); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *Ban")
		return
	}

	ban.Reason = b.Reason

	if b.User != nil {
		ban.User = b.User.DeepCopy().(*User)
	}

	return
}

// ------------

// GuildEmbed https://discord.com/developers/docs/resources/guild#guild-embed-object
type GuildEmbed struct {
	Enabled   bool      `json:"enabled"`
	ChannelID Snowflake `json:"channel_id"`
}

// DeepCopy see interface at struct.go#DeepCopier
func (e *GuildEmbed) DeepCopy() (copy interface{}) {
	copy = &GuildEmbed{}
	e.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (e *GuildEmbed) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var embed *GuildEmbed
	if embed, ok = other.(*GuildEmbed); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *GuildEmbed")
		return
	}

	embed.Enabled = e.Enabled
	embed.ChannelID = e.ChannelID

	return
}

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

// DeepCopy see interface at struct.go#DeepCopier
func (i *Integration) DeepCopy() (copy interface{}) {
	copy = &Integration{}
	i.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (i *Integration) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var integration *Integration
	if integration, ok = other.(*Integration); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *Integration")
		return
	}

	integration.ID = i.ID
	integration.Name = i.Name
	integration.Type = i.Type
	integration.Enabled = i.Enabled
	integration.Syncing = i.Syncing
	integration.RoleID = i.RoleID
	integration.ExpireBehavior = i.ExpireBehavior
	integration.ExpireGracePeriod = i.ExpireGracePeriod

	if i.User != nil {
		integration.User = i.User.DeepCopy().(*User)
	}
	if i.Account != nil {
		integration.Account = i.Account.DeepCopy().(*IntegrationAccount)
	}

	return
}

// IntegrationAccount https://discord.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id"`   // id of the account
	Name string `json:"name"` // name of the account
}

// DeepCopy see interface at struct.go#DeepCopier
func (i *IntegrationAccount) DeepCopy() (copy interface{}) {
	copy = &IntegrationAccount{}
	i.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (i *IntegrationAccount) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var account *IntegrationAccount
	if account, ok = other.(*IntegrationAccount); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *IntegrationAccount")
		return
	}

	account.ID = i.ID
	account.Name = i.Name
	return
}

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

	// custom
	UserID Snowflake `json:"-"`
}

var _ Reseter = (*Member)(nil)
var _ fmt.Stringer = (*Member)(nil)
var _ internalUpdater = (*Member)(nil)
var _ Mentioner = (*Member)(nil)

func (m *Member) updateInternals() {
	if m.User != nil {
		m.UserID = m.User.ID
	}
}

func (m *Member) String() string {
	usrname := m.Nick
	if m.User != nil {
		usrname = m.User.Username
	}
	id := m.UserID
	if m.UserID.IsZero() && m.User != nil {
		id = m.User.ID
	}
	return "member{user:" + usrname + ", nick:" + m.Nick + ", ID:" + id.String() + "}"
}

type GuildQueryBuilderCaller interface {
	Guild(id Snowflake) GuildQueryBuilder
}

func (m *Member) UpdateNick(ctx context.Context, client GuildQueryBuilderCaller, nickname string, flags ...Flag) error {
	return client.
		Guild(m.GuildID).
		Member(m.UserID).
		WithContext(ctx).
		Update(flags...).
		SetNick(nickname).
		Execute()
}

// GetPermissions populates a uint64 with all the permission flags
func (m *Member) GetPermissions(ctx context.Context, s GuildQueryBuilderCaller, flags ...Flag) (permissions PermissionBit, err error) {
	// TODO: Don't deep copy channels for this in the future!
	roles, err := s.Guild(m.GuildID).WithContext(ctx).GetRoles(flags...)
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

// DeepCopy see interface at struct.go#DeepCopier
func (m *Member) DeepCopy() (copy interface{}) {
	copy = &Member{}
	m.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (m *Member) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var member *Member
	if member, ok = other.(*Member); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *Member")
		return
	}

	member.GuildID = m.GuildID
	member.Nick = m.Nick
	member.Roles = m.Roles
	member.JoinedAt = m.JoinedAt
	member.Deaf = m.Deaf
	member.Mute = m.Mute
	member.UserID = m.UserID

	if m.User != nil {
		member.User = m.User.DeepCopy().(*User)
	}
	return
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

// CreateGuildParams ...
// https://discord.com/developers/docs/resources/guild#create-guild-json-params
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
//  Discord documentation   https://discord.com/developers/docs/resources/guild#create-guild
//  Reviewed                2018-08-16
//  Comment                 This endpoint. can be used only by bots in less than 10 Guilds. Creating channel
//                          categories from this endpoint. is not supported.
//							The params argument is optional.
func (c clientQueryBuilder) CreateGuild(guildName string, params *CreateGuildParams, flags ...Flag) (ret *Guild, err error) {
	// TODO: check if bot
	// TODO-2: is bot in less than 10 Guilds?

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

	r := c.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         c.ctx,
		Endpoint:    endpoint.Guilds(),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Guild{}
	}

	return getGuild(r.Execute)
}

// GuildQueryBuilder defines the exposed functions from the guild query builder.
type GuildQueryBuilder interface {
	WithContext(ctx context.Context) GuildQueryBuilder

	// TODO: Add more guild attribute things. Waiting for caching changes before then.
	Get(flags ...Flag) (guild *Guild, err error)
	// TODO: For GetChannels, it might sense to have the option for a function to filter before each channel ends up deep copied.
	// TODO-2: This could be much more performant in guilds with a large number of channels.
	GetChannels(flags ...Flag) ([]*Channel, error)
	// TODO: For GetMembers, it might sense to have the option for a function to filter before each member ends up deep copied.
	// TODO-2: This could be much more performant in larger guilds where this is needed.
	GetMembers(params *GetMembersParams, flags ...Flag) ([]*Member, error)
	Update(flags ...Flag) UpdateGuildBuilder
	Delete(flags ...Flag) error

	CreateChannel(name string, params *CreateGuildChannelParams, flags ...Flag) (*Channel, error)
	UpdateChannelPositions(params []UpdateGuildChannelPositionsParams, flags ...Flag) error
	CreateMember(userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (*Member, error)
	Member(userID Snowflake) GuildMemberQueryBuilder

	KickVoiceParticipant(userID Snowflake) error
	SetCurrentUserNick(nick string, flags ...Flag) (newNick string, err error)
	GetBans(flags ...Flag) ([]*Ban, error)
	GetBan(userID Snowflake, flags ...Flag) (*Ban, error)
	UnbanUser(userID Snowflake, reason string, flags ...Flag) error
	// TODO: For GetRoles, it might sense to have the option for a function to filter before each role ends up deep copied.
	// TODO-2: This could be much more performant in larger guilds where this is needed.
	// TODO-3: Add GetRole.
	GetRoles(flags ...Flag) ([]*Role, error)
	GetMemberPermissions(userID Snowflake, flags ...Flag) (permissions PermissionBit, err error)
	UpdateRolePositions(params []UpdateGuildRolePositionsParams, flags ...Flag) ([]*Role, error)
	CreateRole(params *CreateGuildRoleParams, flags ...Flag) (*Role, error)
	Role(roleID Snowflake) GuildRoleQueryBuilder

	EstimatePruneMembersCount(days int, flags ...Flag) (estimate int, err error)
	PruneMembers(days int, reason string, flags ...Flag) error
	GetVoiceRegions(flags ...Flag) ([]*VoiceRegion, error)
	GetInvites(flags ...Flag) ([]*Invite, error)

	GetIntegrations(flags ...Flag) ([]*Integration, error)
	CreateIntegration(params *CreateGuildIntegrationParams, flags ...Flag) error
	UpdateIntegration(integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error
	DeleteIntegration(integrationID Snowflake, flags ...Flag) error
	SyncIntegration(integrationID Snowflake, flags ...Flag) error

	GetEmbed(flags ...Flag) (*GuildEmbed, error)
	UpdateEmbed(flags ...Flag) UpdateGuildEmbedBuilder
	GetVanityURL(flags ...Flag) (*PartialInvite, error)
	GetAuditLogs(flags ...Flag) GuildAuditLogsBuilder
	VoiceConnect(channelID Snowflake) (ret VoiceConnection, err error)

	// TODO: For GetEmojis, it might sense to have the option for a function to filter before each emoji ends up deep copied.
	// TODO-2: This could be much more performant in guilds with a large number of channels.
	GetEmojis(flags ...Flag) ([]*Emoji, error)
	CreateEmoji(params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error)
	Emoji(emojiID Snowflake) GuildEmojiQueryBuilder

	GetWebhooks(flags ...Flag) (ret []*Webhook, err error)
}

// Guild is used to create a guild query builder.
func (c clientQueryBuilder) Guild(id Snowflake) GuildQueryBuilder {
	return &guildQueryBuilder{client: c.client, gid: id}
}

// The default guild query builder.
type guildQueryBuilder struct {
	ctx    context.Context
	client *Client
	gid    Snowflake
}

func (g guildQueryBuilder) WithContext(ctx context.Context) GuildQueryBuilder {
	g.ctx = ctx
	return &g
}

// Get is used to get the Guild struct containing all information from it.
// Note that it's significantly quicker in most instances where you have the cache enabled (as is by default) to get the individual parts you need.
func (g guildQueryBuilder) Get(flags ...Flag) (guild *Guild, err error) {
	if guild, _ = g.client.cache.GetGuild(g.gid); guild != nil {
		return guild, nil
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Guild(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		return &Guild{}
	}

	return getGuild(r.Execute)
}

// Update is used to create a guild update builder.
func (g guildQueryBuilder) Update(flags ...Flag) UpdateGuildBuilder {
	builder := &updateGuildBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Guild{}
	}
	builder.r.setup(g.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.Guild(g.gid),
		ContentType: httd.ContentTypeJSON,
	}, nil)
	builder.r.flags = flags

	return builder
}

// Delete is used to delete a guild.
func (g guildQueryBuilder) Delete(flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.Guild(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// GetChannels is used to get a guilds channels.
func (g guildQueryBuilder) GetChannels(flags ...Flag) ([]*Channel, error) {
	if channels, _ := g.client.cache.GetGuildChannels(g.gid); channels != nil {
		return channels, nil
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildChannels(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Channel, 0)
		return &tmp
	}

	return getChannels(r.Execute)
}

// CreateChannel Create a new channel object for the guild. Requires the 'MANAGE_CHANNELS' permission.
// Returns the new channel object on success. Fires a Channel Create Gateway event.
func (g guildQueryBuilder) CreateChannel(name string, params *CreateGuildChannelParams, flags ...Flag) (*Channel, error) {
	if name == "" && (params == nil || params.Name == "") {
		return nil, errors.New("channel name is required")
	}
	if l := len(name); !(2 <= l && l <= 100) {
		return nil, errors.New("channel name must be 2 or more characters and no more than 100 characters")
	}

	if params == nil {
		params = &CreateGuildChannelParams{}
	}
	if name != "" && params.Name == "" {
		params.Name = name
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildChannels(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      params.Reason,
	}, flags)
	r.factory = func() interface{} {
		return &Channel{}
	}

	return getChannel(r.Execute)
}

// UpdateChannelPositions Modify the positions of a set of channel objects for the guild.
// Requires 'MANAGE_CHANNELS' permission. Returns a 204 empty response on success. Fires multiple Channel Update
// Gateway events.
func (g guildQueryBuilder) UpdateChannelPositions(params []UpdateGuildChannelPositionsParams, flags ...Flag) error {
	var reason string
	for i := range params {
		if params[i].Reason != "" {
			reason = params[i].Reason
			break
		}
	}
	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildChannels(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      reason,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// GetMembers uses the GetGuildMembers endpoint iteratively until your query params are met.
func (g guildQueryBuilder) GetMembers(params *GetMembersParams, flags ...Flag) ([]*Member, error) {
	if params == nil {
		params = &GetMembersParams{
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

	p := getGuildMembersParams{
		After: params.After,
	}
	if params.Limit == 0 || params.Limit > 1000 {
		p.Limit = 1000
	} else {
		p.Limit = int(params.Limit)
	}

	members := make([]*Member, 0)
	var ms []*Member
	var err error
	for {
		ms, err = g.getGuildMembers(&p, flags...)
		if ms != nil {
			members = append(members, ms...)
		}
		if err != nil {
			return members, err
		}

		// stop if we're on the last page/block of members
		if len(ms) < 1000 {
			break
		}

		// set limit such that we don't retrieve redundant members
		max := params.Limit << 1
		max = max >> 1
		lim := int(max) - len(members)
		if lim < 1000 {
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
func (g guildQueryBuilder) CreateMember(userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (*Member, error) {
	if accessToken == "" && (params == nil || params.AccessToken == "") {
		return nil, errors.New("access token is required")
	}

	if params == nil {
		params = &AddGuildMemberParams{}
	}
	if accessToken != "" && params.AccessToken == "" {
		params.AccessToken = accessToken
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPut,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildMember(g.gid, userID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.factory = func() interface{} {
		return &Member{
			GuildID: g.gid,
			UserID:  userID,
		}
	}
	r.expectsStatusCode = http.StatusCreated

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
func (g guildQueryBuilder) SetCurrentUserNick(nick string, flags ...Flag) (newNick string, err error) {
	params := &updateCurrentUserNickParams{
		Nick: nick,
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildMembersMeNick(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusOK
	r.factory = func() interface{} {
		return &nickNameResponse{}
	}

	return getNickName(r.Execute)
}

// GetBans returns an array of ban objects for the Users banned from this guild. Requires the 'BAN_MEMBERS' permission.
func (g guildQueryBuilder) GetBans(flags ...Flag) ([]*Ban, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildBans(g.gid),
		Ctx:      g.ctx,
	}, flags)
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
func (g guildQueryBuilder) GetBan(userID Snowflake, flags ...Flag) (*Ban, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildBan(g.gid, userID),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		return &Ban{User: g.client.pool.user.Get().(*User)}
	}

	return getBan(r.Execute)
}

// UnbanMember Remove the ban for a user. Requires the 'BAN_MEMBERS' permissions.
// Returns a 204 empty response on success. Fires a Guild Ban Remove Gateway event.
func (g guildQueryBuilder) UnbanUser(userID Snowflake, reason string, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildBan(g.gid, userID),
		Reason:   reason,
		Ctx:      g.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// GetRoles Returns a list of role objects for the guild.
func (g guildQueryBuilder) GetRoles(flags ...Flag) ([]*Role, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: "/guilds/" + g.gid.String() + "/roles",
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Role, 0)
		return &tmp
	}

	return getRoles(r.Execute)
}

// GetMemberPermissions is used to return the members permissions.
func (g guildQueryBuilder) GetMemberPermissions(userID Snowflake, flags ...Flag) (permissions PermissionBit, err error) {
	member, err := g.Member(userID).WithContext(g.ctx).Get(flags...)
	if err != nil {
		return 0, err
	}
	return member.GetPermissions(g.ctx, g.client, flags...)
}

// CreateGuildRoleParams ...
// https://discord.com/developers/docs/resources/guild#create-guild-role-json-params
type CreateGuildRoleParams struct {
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
func (g guildQueryBuilder) CreateRole(params *CreateGuildRoleParams, flags ...Flag) (*Role, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildRoles(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      params.Reason,
	}, flags)
	r.factory = func() interface{} {
		return &Role{}
	}

	return getRole(r.Execute)
}

// UpdateRolePositions Modify the positions of a set of role objects for the guild.
// Requires the 'MANAGE_ROLES' permission. Returns a list of all of the guild's role objects on success.
// Fires multiple Guild Role Update Gateway events.
func (g guildQueryBuilder) UpdateRolePositions(params []UpdateGuildRolePositionsParams, flags ...Flag) ([]*Role, error) {
	var reason string
	for i := range params {
		if params[i].Reason != "" {
			reason = params[i].Reason
			break
		}
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildRoles(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      reason,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Role, 0)
		return &tmp
	}

	return getRoles(r.Execute)
}

// EstimatePruneMembersCount Returns an object with one 'pruned' key indicating the number of members that would be
// removed in a prune operation. Requires the 'KICK_MEMBERS' permission.
func (g guildQueryBuilder) EstimatePruneMembersCount(days int, flags ...Flag) (estimate int, err error) {
	if g.gid.IsZero() {
		return 0, errors.New("guildID can not be " + g.gid.String())
	}
	params := pruneMembersParams{Days: days}
	if err = params.FindErrors(); err != nil {
		return 0, err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildPrune(g.gid) + params.URLQueryString(),
		Ctx:      g.ctx,
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

// PruneMembers Kicks members from N day back. Requires the 'KICK_MEMBERS' permission.
// The estimate of kicked people is not returned. Use EstimatePruneMembersCount before calling PruneMembers
// if you need it. Fires multiple Guild Member Remove Gateway events.
func (g guildQueryBuilder) PruneMembers(days int, reason string, flags ...Flag) (err error) {
	params := pruneMembersParams{Days: days}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodPost,
		Endpoint: endpoint.GuildPrune(g.gid) + params.URLQueryString(),
		Ctx:      g.ctx,
		Reason:   reason,
	}, flags)

	_, err = r.Execute()
	return err
}

// GetVoiceRegions Returns a list of voice region objects for the guild. Unlike the similar /voice route,
// this returns VIP servers when the guild is VIP-enabled.
func (g guildQueryBuilder) GetVoiceRegions(flags ...Flag) ([]*VoiceRegion, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildRegions(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*VoiceRegion, 0)
		return &tmp
	}

	return getVoiceRegions(r.Execute)
}

// GetInvites Returns a list of invite objects (with invite metadata) for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (g guildQueryBuilder) GetInvites(flags ...Flag) ([]*Invite, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildInvites(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Invite, 0)
		return &tmp
	}

	return getInvites(r.Execute)
}

// GetIntegrations Returns a list of integration objects for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (g guildQueryBuilder) GetIntegrations(flags ...Flag) ([]*Integration, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildIntegrations(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Integration, 0)
		return &tmp
	}

	return getIntegrations(r.Execute)
}

// CreateIntegration attaches an integration object from the current user to the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
func (g guildQueryBuilder) CreateIntegration(params *CreateGuildIntegrationParams, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildIntegrations(g.gid),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// UpdateIntegration Modify the behavior and settings of a integration object for the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
func (g guildQueryBuilder) UpdateIntegration(integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildIntegration(g.gid, integrationID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// DeleteIntegration Delete the attached integration object for the guild.
// Requires the 'MANAGE_GUILD' permission. Returns a 204 empty response on success.
// Fires a Guild Integrations Update Gateway event.
func (g guildQueryBuilder) DeleteIntegration(integrationID Snowflake, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodDelete,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildIntegration(g.gid, integrationID),
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// SyncIntegration Sync an integration. Requires the 'MANAGE_GUILD' permission.
// Returns a 204 empty response on success.
func (g guildQueryBuilder) SyncIntegration(integrationID Snowflake, flags ...Flag) error {
	r := g.client.newRESTRequest(&httd.Request{
		Method:   httd.MethodPost,
		Endpoint: endpoint.GuildIntegrationSync(g.gid, integrationID),
		Ctx:      g.ctx,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent

	_, err := r.Execute()
	return err
}

// GetEmbed Returns the guild embed object. Requires the 'MANAGE_GUILD' permission.
func (g guildQueryBuilder) GetEmbed(flags ...Flag) (*GuildEmbed, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmbed(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		return &GuildEmbed{}
	}

	return getGuildEmbed(r.Execute)
}

// UpdateEmbed Modify a guild embed object for the guild. All attributes may be passed in with JSON and
// modified. Requires the 'MANAGE_GUILD' permission. Returns the updated guild embed object.
func (g guildQueryBuilder) UpdateEmbed(flags ...Flag) UpdateGuildEmbedBuilder {
	builder := &updateGuildEmbedBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &GuildEmbed{}
	}
	builder.r.flags = flags
	builder.r.setup(g.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildEmbed(g.gid),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

// GetVanityURL Returns a partial invite object for Guilds with that feature enabled.
// Requires the 'MANAGE_GUILD' permission.
func (g guildQueryBuilder) GetVanityURL(flags ...Flag) (*PartialInvite, error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildVanityURL(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		return &PartialInvite{}
	}

	return getPartialInvite(r.Execute)
}

// GetAuditLogs Returns an audit log object for the guild. Requires the 'VIEW_AUDIT_LOG' permission.
// Note that this request will _always_ send a REST request, regardless of you calling IgnoreCache or not.
func (g guildQueryBuilder) GetAuditLogs(flags ...Flag) GuildAuditLogsBuilder {
	builder := &guildAuditLogsBuilder{}
	builder.r.itemFactory = auditLogFactory
	builder.r.flags = flags
	builder.r.IgnoreCache().setup(g.client.req, &httd.Request{
		Ctx:      g.ctx,
		Method:   httd.MethodGet,
		Endpoint: endpoint.GuildAuditLogs(g.gid),
	}, nil)

	return builder
}

// VoiceConnect is used to handle making a voice connection.
func (g guildQueryBuilder) VoiceConnect(channelID Snowflake) (VoiceConnection, error) {
	return g.client.VoiceConnectOptions(g.gid, channelID, true, false)
}

// GetEmojis Returns a list of emoji objects for the given guild.
func (g guildQueryBuilder) GetEmojis(flags ...Flag) ([]*Emoji, error) {
	if emojis, _ := g.client.cache.GetGuildEmojis(g.gid); emojis != nil {
		return emojis, nil
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmojis(g.gid),
		Ctx:      g.ctx,
	}, flags)
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

// CreateGuildEmojiParams JSON params for func CreateGuildEmoji
type CreateGuildEmojiParams struct {
	Name  string      `json:"name"`  // required
	Image string      `json:"image"` // required
	Roles []Snowflake `json:"roles"` // optional

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// CreateEmoji Create a new emoji for the guild. Requires the 'MANAGE_EMOJIS' permission.
// Returns the new emoji object on success. Fires a Guild Emojis Update Gateway event.
func (g guildQueryBuilder) CreateEmoji(params *CreateGuildEmojiParams, flags ...Flag) (*Emoji, error) {
	if g.gid.IsZero() {
		return nil, errors.New("guildID must be set, was " + g.gid.String())
	}

	if params == nil {
		return nil, errors.New("params object can not be nil")
	}
	if !validEmojiName(params.Name) {
		return nil, errors.New("invalid emoji name")
	}
	if !validAvatarPrefix(params.Image) {
		return nil, errors.New("image string must be base64 encoded with base64 prefix")
	}

	r := g.client.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
		Ctx:         g.ctx,
		Endpoint:    endpoint.GuildEmojis(g.gid),
		ContentType: httd.ContentTypeJSON,
		Body:        params,
		Reason:      params.Reason,
	}, flags)
	r.pool = g.client.pool.emoji
	r.factory = func() interface{} {
		return &Emoji{}
	}

	return getEmoji(r.Execute)
}

// KickVoiceParticipant is used to kick someone from voice.
func (g guildQueryBuilder) KickVoiceParticipant(userID Snowflake) error {
	return g.Member(userID).WithContext(g.ctx).Update().KickFromVoice().Execute()
}

// GetWebhooks Returns a list of guild webhook objects. Requires the 'MANAGE_WEBHOOKS' permission.
func (g guildQueryBuilder) GetWebhooks(flags ...Flag) (ret []*Webhook, err error) {
	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildWebhooks(g.gid),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Webhook, 0)
		return &tmp
	}

	return getWebhooks(r.Execute)
}

// CreateGuildChannelParams https://discord.com/developers/docs/resources/guild#create-guild-channel-json-params
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
	Position             int                   `json:"position"` // can not omitempty in case position is 0

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

// UpdateGuildChannelPositionsParams ...
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions-json-params
type UpdateGuildChannelPositionsParams struct {
	ID       Snowflake `json:"id"`
	Position int       `json:"position"`

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	// just reuse the string. Go will optimize it to point to the same memory anyways
	// TODO: improve this?
	Reason string `json:"-"`
}

func NewUpdateGuildRolePositionsParams(rs []*Role) (p []UpdateGuildRolePositionsParams) {
	p = make([]UpdateGuildRolePositionsParams, 0, len(rs))
	for i := range rs {
		p = append(p, UpdateGuildRolePositionsParams{
			ID:       rs[i].ID,
			Position: rs[i].Position,
		})
	}

	return p
}

// UpdateGuildRolePositionsParams ...
// https://discord.com/developers/docs/resources/guild#modify-guild-role-positions-json-params
type UpdateGuildRolePositionsParams struct {
	ID       Snowflake `json:"id"`
	Position int       `json:"position"`

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string `json:"-"`
}

type getGuildMembersParams struct {
	After Snowflake `urlparam:"after,omitempty"`
	Limit int       `urlparam:"limit,omitempty"` // 1 is default. even if 0 is supplied.
}

var _ URLQueryStringer = (*getGuildMembersParams)(nil)

func (g *getGuildMembersParams) FindErrors() error {
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
func (g guildQueryBuilder) getGuildMembers(params *getGuildMembersParams, flags ...Flag) (ret []*Member, err error) {
	if params == nil {
		params = &getGuildMembersParams{}
	}
	if err = params.FindErrors(); err != nil {
		return nil, err
	}

	r := g.client.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildMembers(g.gid) + params.URLQueryString(),
		Ctx:      g.ctx,
	}, flags)
	r.factory = func() interface{} {
		tmp := make([]*Member, 0)
		return &tmp
	}

	return getMembers(r.Execute)
}

// GetMembersParams if Limit is 0, every member is fetched. This does not follow the Discord API where a 0
// is converted into a 1. 0 = every member. The rest is exactly the same, you should be able to do everything
// the Discord docs says with the addition that you can bypass a limit of 1,000.
//
// If you specify a limit of +1,000 Disgord will run N requests until that amount is met, or until you run
// out of members to fetch.
type GetMembersParams struct {
	After Snowflake `urlparam:"after,omitempty"`
	Limit uint32    `urlparam:"limit,omitempty"` // 0 will fetch everyone
}

// AddGuildMemberParams ...
// https://discord.com/developers/docs/resources/guild#add-guild-member-json-params
type AddGuildMemberParams struct {
	AccessToken string      `json:"access_token"` // required
	Nick        string      `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Mute        bool        `json:"mute,omitempty"`
	Deaf        bool        `json:"deaf,omitempty"`
}

// BanMemberParams ...
// https://discord.com/developers/docs/resources/guild#create-guild-ban-query-string-params
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

// PruneMembersParams will delete members, this is the same as kicking.
// https://discord.com/developers/docs/resources/guild#get-guild-prune-count-query-string-params
type pruneMembersParams struct {
	// Days number of days to count prune for (1 or more)
	Days int `urlparam:"days"`

	// ComputePruneCount whether 'pruned' is returned, discouraged for large Guilds
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

// CreateGuildIntegrationParams ...
// https://discord.com/developers/docs/resources/guild#create-guild-integration-json-params
type CreateGuildIntegrationParams struct {
	Type string    `json:"type"`
	ID   Snowflake `json:"id"`
}

// UpdateGuildIntegrationParams ...
// https://discord.com/developers/docs/resources/guild#modify-guild-integration-json-params
// TODO: currently unsure which are required/optional params
type UpdateGuildIntegrationParams struct {
	ExpireBehavior    int  `json:"expire_behavior"`
	ExpireGracePeriod int  `json:"expire_grace_period"`
	EnableEmoticons   bool `json:"enable_emoticons"`
}

// updateCurrentUserNickParams ...
// https://discord.com/developers/docs/resources/guild#modify-guild-member-json-params
type updateCurrentUserNickParams struct {
	Nick string `json:"nick"` // :CHANGE_NICKNAME
}

type nickNameResponse struct {
	Nickname string `json:"nickname"`
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// updateGuildBuilder https://discord.com/developers/docs/resources/guild#modify-guild-json-params
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
// https://discord.com/developers/docs/resources/guild#modify-guild-member-json-params
//generate-rest-params: nick:string, roles:[]Snowflake, mute:bool, deaf:bool, channel_id:Snowflake,
//generate-rest-basic-execute: err:error,
type updateGuildMemberBuilder struct {
	r RESTBuilder
}

// KickFromVoice kicks member out of voice channel. Assuming they are in one.
func (b *updateGuildMemberBuilder) KickFromVoice() UpdateGuildMemberBuilder {
	b.r.param("channel_id", 0)
	return b
}

// DeleteNick removes nickname for user. Requires permission MANAGE_NICKNAMES
func (b *updateGuildMemberBuilder) DeleteNick() UpdateGuildMemberBuilder {
	b.r.param("nick", "")
	return b
}
