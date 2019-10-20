package disgord

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"

	"github.com/andersfylling/disgord/internal/constant"
)

// NOTE! Credit for defining the Permission consts in a clean way goes to DiscordGo.
// This is pretty much a copy from their project. I would have made it a dependency if
// the consts were in a isolated sub-pkg. Note that in respect to their license, DisGord
// has no affiliation with DiscordGo.
//
// Source code reference:
//  https://github.com/bwmarrin/discordgo/blob/8325a6bf6dd6c91ed4040a1617b07287b8fb0eba/structs.go#L854

type PermissionBit = uint64
type PermissionBits = PermissionBit

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

// NewGuildFromJSON ...
func NewGuildFromJSON(data []byte) (guild *Guild) {
	guild = NewGuild()
	err := unmarshal(data, guild)
	if err != nil {
		panic(err)
	}

	return guild
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
	Lockable    `json:"-"`
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

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
// reviewed: 2018-08-25
type Guild struct {
	Lockable `json:"-"`

	ID                          Snowflake                     `json:"id"`
	ApplicationID               Snowflake                     `json:"application_id"` //   |?
	Name                        string                        `json:"name"`
	Icon                        string                        `json:"icon"`            //  |?, icon hash
	Splash                      string                        `json:"splash"`          //  |?, image hash
	Owner                       bool                          `json:"owner,omitempty"` // ?|
	OwnerID                     Snowflake                     `json:"owner_id"`
	Permissions                 PermissionBits                `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
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

	//highestSnowflakeAmoungMembers Snowflake
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
	for i := range g.Emojis {
		g.Emojis[i].guildID = g.ID
	}
	for i := range g.Channels {
		g.Channels[i].GuildID = g.ID
	}
}

func (g *Guild) copyOverToCache(other interface{}) (err error) {
	guild := other.(*Guild)

	if constant.LockedMethods {
		g.RLock()
		guild.Lock()
	}

	//guild.ID = g.ID
	if g.Name != "" {
		guild.Name = g.Name
	}
	guild.Owner = g.Owner
	// Use ownerID to check if you are the owner of the guild(!)
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
	g.Icon = g.Icon

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

	if constant.LockedMethods {
		g.RUnlock()
		guild.Unlock()
	}

	return
}

// GetMemberWithHighestSnowflake finds the member with the highest snowflake value.
func (g *Guild) GetMemberWithHighestSnowflake() *Member {
	g.RLock()
	defer g.RUnlock()

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

// func (g *Guild) UnmarshalJSON(data []byte) (err error) {
// 	return json.Unmarshal(data, &g.guildJSON)
// }

// MarshalJSON see interface json.Marshaler
// TODO: fix copying of mutex lock
func (g *Guild) MarshalJSON() ([]byte, error) {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

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
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

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
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

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

func (g *Guild) addMember(member *Member) error {
	if member == nil {
		return errors.New("member was nil")
	}
	// TODO: implement sorting for faster searching later
	g.Members = append(g.Members, member)

	return nil
}

// AddMembers adds multiple members to the Guild object. Note that this method does not interact with Discord.
func (g *Guild) AddMembers(members []*Member) {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

	for _, member := range members {
		g.addMember(member)
	}
}

// AddMember adds a member to the Guild object. Note that this method does not interact with Discord.
func (g *Guild) AddMember(member *Member) error {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

	return g.addMember(member)
}

// LoadAllMembers fetches all the members for this guild from the Discord REST API
func (g *Guild) LoadAllMembers(s Session) (err error) {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}
	// TODO: what if members have already been loaded? use Guild.MembersCount?

	members, err := s.GetMembers(g.ID, nil)
	if err != nil {
		return err
	}

	for i := range members {
		_ = g.addMember(members[i])
	}

	return nil
}

// GetMembersCountEstimate estimates the number of members in a guild without fetching everyone.
// There is no proper way to get this number, so a invite is created and the estimate
// is read from there. The invite is then deleted again.
func (g *Guild) GetMembersCountEstimate(s Session) (estimate int, err error) {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

	var channelID Snowflake
	if len(g.Channels) == 0 {
		channels, err := s.GetGuildChannels(g.ID)
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

	invite, err := s.CreateChannelInvites(channelID, &CreateChannelInvitesParams{
		MaxAge: 1,
	})
	if err != nil {
		return 0, err
	}
	_ = s.DeleteFromDiscord(invite) // delete if possible

	return invite.ApproximateMemberCount, nil
}

// AddRole adds a role to the Guild object. Note that this does not interact with Discord.
func (g *Guild) AddRole(role *Role) error {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

	// TODO: implement sorting for faster searching later
	role.guildID = g.ID
	g.Roles = append(g.Roles, role)

	return nil
}

// Member return a member by his/her userid
func (g *Guild) Member(id Snowflake) (*Member, error) {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

	for _, member := range g.Members {
		if member.User.ID == id {
			return member, nil
		}
	}

	return nil, errors.New("member not found in guild")
}

// MembersByName retrieve a slice of members with same username or nickname
func (g *Guild) MembersByName(name string) (members []*Member) {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

	for _, member := range g.Members {
		if member.Nick == name || member.User.Username == name {
			members = append(members, member)
		}
	}

	return
}

// Role retrieve a role based on role id
func (g *Guild) Role(id Snowflake) (role *Role, err error) {
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

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
	if constant.LockedMethods {
		g.Lock()
		defer g.Unlock()
	}

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
	if constant.LockedMethods {
		g.RLock()
		defer g.RUnlock()
	}

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
	if constant.LockedMethods {
		g.RLock()
		defer g.RUnlock()
	}

	for _, channel := range g.Channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, errors.New("channel not found in guild")
}

// Emoji get a guild emoji by it's ID
func (g *Guild) Emoji(id Snowflake) (emoji *Emoji, err error) {
	if constant.LockedMethods {
		g.RLock()
		defer g.RUnlock()
	}

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

	if constant.LockedMethods {
		g.RLock()
		guild.Lock()
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
	g.Icon = g.Icon

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

	if constant.LockedMethods {
		g.RUnlock()
		guild.Unlock()
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

// Ban https://discordapp.com/developers/docs/resources/guild#ban-object
type Ban struct {
	Lockable `json:"-"`

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

	if constant.LockedMethods {
		b.RLock()
		ban.Lock()
	}

	ban.Reason = b.Reason

	if b.User != nil {
		ban.User = b.User.DeepCopy().(*User)
	}

	if constant.LockedMethods {
		b.RUnlock()
		ban.Unlock()
	}

	return
}

// ------------

// GuildEmbed https://discordapp.com/developers/docs/resources/guild#guild-embed-object
type GuildEmbed struct {
	Lockable `json:"-"`

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

	if constant.LockedMethods {
		e.RLock()
		embed.Lock()
	}

	embed.Enabled = e.Enabled
	embed.ChannelID = e.ChannelID

	if constant.LockedMethods {
		e.RUnlock()
		embed.Unlock()
	}

	return
}

// -------

// Integration https://discordapp.com/developers/docs/resources/guild#integration-object
type Integration struct {
	Lockable `json:"-"`

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

	if constant.LockedMethods {
		i.RLock()
		integration.Lock()
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

	if constant.LockedMethods {
		i.RUnlock()
		integration.Unlock()
	}

	return
}

// IntegrationAccount https://discordapp.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	Lockable `json:"-"`

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

	if constant.LockedMethods {
		i.RLock()
		account.Lock()
	}

	account.ID = i.ID
	account.Name = i.Name

	if constant.LockedMethods {
		i.RUnlock()
		account.Unlock()
	}

	return
}

// -------

// Member https://discordapp.com/developers/docs/resources/guild#guild-member-object
type Member struct {
	Lockable `json:"-"`

	GuildID  Snowflake   `json:"guild_id,omitempty"`
	User     *User       `json:"user"`
	Nick     string      `json:"nick,omitempty"` // ?|
	Roles    []Snowflake `json:"roles"`
	JoinedAt Time        `json:"joined_at,omitempty"`

	// voice
	Deaf bool `json:"deaf"`
	Mute bool `json:"mute"`

	// used for caching
	userID Snowflake
}

var _ Reseter = (*Member)(nil)

func (m *Member) String() string {
	usrname := m.Nick
	if m.User != nil {
		usrname = m.User.Username
	}
	id := m.userID
	if m.userID.IsZero() && m.User != nil {
		id = m.User.ID
	}
	return "member{user:" + usrname + ", nick:" + m.Nick + ", ID:" + id.String() + "}"
}

type nickUpdater interface {
	UpdateGuildMember(guildID, userID Snowflake, flags ...Flag) *updateGuildMemberBuilder
}

func (m *Member) UpdateNick(client nickUpdater, nickname string, flags ...Flag) error {
	return client.UpdateGuildMember(m.GuildID, m.userID, flags...).SetNick(nickname).Execute()
}

func (m *Member) GetPermissions(s Session) (p uint64, err error) {
	uID := m.userID
	if uID.IsZero() {
		usr, err := m.GetUser(s)
		if err != nil {
			return 0, err
		}
		uID = usr.ID
	}
	return s.GetMemberPermissions(m.GuildID, uID)
}

// GetUser tries to ensure that you get a user object and not a nil. The user can be nil if the guild
// was fetched from the cache.
func (m *Member) GetUser(session Session) (usr *User, err error) {
	if m.User != nil {
		return m.User, nil
	}

	return session.GetUser(m.userID)
}

// Mention creates a string which is parsed into a member mention on Discord GUI's
func (m *Member) Mention() string {
	var id Snowflake
	if !m.userID.IsZero() {
		id = m.userID
	} else if m.User != nil {
		id = m.User.ID
	} else {
		fmt.Println("ERRPR: unable to fetch user id. please create a issue at github.com/andersfylling/disgord")
		return "*" + m.Nick + "*"
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

	if constant.LockedMethods {
		m.RLock()
		member.Lock()
	}

	member.GuildID = m.GuildID
	member.Nick = m.Nick
	member.Roles = m.Roles
	member.JoinedAt = m.JoinedAt
	member.Deaf = m.Deaf
	member.Mute = m.Mute
	member.userID = m.userID

	if m.User != nil {
		member.User = m.User.DeepCopy().(*User)
	}

	if constant.LockedMethods {
		m.RUnlock()
		member.Unlock()
	}

	return
}

//////////////////////////////////////////////////////
//
// REST Methods
//
//////////////////////////////////////////////////////

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

// CreateGuild [REST] Add a new guild. Returns a guild object on success. Fires a Guild Add Gateway event.
//  Method                  POST
//  Endpoint                /guilds
//  Rate limiter            /guilds
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild
//  Reviewed                2018-08-16
//  Comment                 This endpoint. can be used only by bots in less than 10 guilds. Creating channel
//                          categories from this endpoint. is not supported.
//							The params argument is optional.
func (c *Client) CreateGuild(guildName string, params *CreateGuildParams, flags ...Flag) (ret *Guild, err error) {
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
		Method:      httd.MethodPost,
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
func (c *Client) GetGuild(id Snowflake, flags ...Flag) (guild *Guild, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Guild(id),
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
func (c *Client) UpdateGuild(id Snowflake, flags ...Flag) (builder *updateGuildBuilder) {
	builder = &updateGuildBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Guild{}
	}
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
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
func (c *Client) DeleteGuild(id Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.Guild(id),
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
func (c *Client) GetGuildChannels(guildID Snowflake, flags ...Flag) (ret []*Channel, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildChannels(guildID),
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
	Position             int                   `json:"position"` // can not omitempty in case position is 0
}

// CreateGuildChannel [REST] Add a new channel object for the guild. Requires the 'MANAGE_CHANNELS' permission.
// Returns the new channel object on success. Fires a Channel Add Gateway event.
//  Method                  POST
//  Endpoint                /guilds/{guild.id}/channels
//  Rate limiter            /guilds/{guild.id}/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild-channel
//  Reviewed                2018-08-17
//  Comment                 All parameters for this endpoint. are optional excluding 'name'
func (c *Client) CreateGuildChannel(guildID Snowflake, channelName string, params *CreateGuildChannelParams, flags ...Flag) (ret *Channel, err error) {
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
		Method:      httd.MethodPost,
		Endpoint:    endpoint.GuildChannels(guildID),
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

// UpdateGuildChannelPositions [REST] Modify the positions of a set of channel objects for the guild.
// Requires 'MANAGE_CHANNELS' permission. Returns a 204 empty response on success. Fires multiple Channel Update
// Gateway events.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/channels
//  Rate limiter            /guilds/{guild.id}/channels
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-channel-positions
//  Reviewed                2018-08-17
//  Comment                 Only channels to be modified are required, with the minimum being a swap
//                          between at least two channels.
func (c *Client) UpdateGuildChannelPositions(guildID Snowflake, params []UpdateGuildChannelPositionsParams, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
		Endpoint:    endpoint.GuildChannels(guildID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
	}, flags)
	r.expectsStatusCode = http.StatusNoContent
	// TODO: update ordering of guild channels in cache

	_, err = r.Execute()
	return err
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
// https://discordapp.com/developers/docs/resources/guild#modify-guild-role-positions-json-params
type UpdateGuildRolePositionsParams struct {
	ID       Snowflake `json:"id"`
	Position int       `json:"position"`
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
func (c *Client) UpdateGuildRolePositions(guildID Snowflake, params []UpdateGuildRolePositionsParams, flags ...Flag) (roles []*Role, err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
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

// UpdateGuildRolesPos "ensures" you have all the roles required,
//func (c *Client) UpdateGuildRolesPos(guildID Snowflake, rs []*Role, flags ...Flag) (updated []*Role, err error) {
//	var current []*Role
//	if current, err = c.GetGuildRoles(guildID, flags...); err != nil {
//		return nil, err
//	}
//	SortRoles(current)
//
//	// add the changes / additions
//	rsOld := rs
//	for i := range current {
//		var exist bool
//		for j := range rs {
//			if current[i].ID == rs[j].ID {
//				exist = true
//				_ = rs[j].CopyOverTo(current[i])
//				break
//			}
//		}
//		if exist {
//			rs[i] = rs[len(rs)-1]
//			rs[len(rs)-1] = nil
//			rs = rs[:len(rs)-1]
//		}
//	}
//	current = append(current, rs...)
//	SortRoles(current)
//
//	// TODO: verify order
//	for i := range rsOld {
//
//	}
//
//	params := make([]UpdateGuildRolePositionsParams, 0, len(current))
//	for i := range current {
//		params = append(params, UpdateGuildRolePositionsParams{
//			ID:       current[i].ID,
//			Position: current[i].Position,
//		})
//	}
//	return c.UpdateGuildRolePositions(guildID, params, flags...)
//}

// GetMember [REST] Returns a guild member object for the specified user.
//  Method                  GET
//  Endpoint                /guilds/{guild.id}/members/{user.id}
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-member
//  Reviewed                2018-08-17
//  Comment                 -
func (c *Client) GetMember(guildID, userID Snowflake, flags ...Flag) (ret *Member, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildMember(guildID, userID),
	}, flags)
	r.CacheRegistry = GuildMembersCache
	r.ID = userID
	r.factory = func() interface{} {
		return &Member{}
	}

	return getMember(r.Execute)
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
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#get-guild-members
//  Reviewed                2018-08-17
//  Comment                 All parameters to this endpoint. are optional
//  Comment#2               "List Guild Members"
//  Comment#3               https://discordapp.com/developers/docs/resources/guild#list-guild-members-query-string-params
func (c *Client) getGuildMembers(guildID Snowflake, params *getGuildMembersParams, flags ...Flag) (ret []*Member, err error) {
	if params == nil {
		params = &getGuildMembersParams{}
	}
	if err = params.FindErrors(); err != nil {
		return nil, err
	}

	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildMembers(guildID) + params.URLQueryString(),
	}, flags)
	r.CacheRegistry = GuildMembersCache
	r.checkCache = func() (v interface{}, err error) {
		return c.cache.GetGuildMembersAfter(guildID, params.After, params.Limit)
	}
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
// If you specify a limit of +1,000 DisGord will run N requests until that amount is met, or until you run
// out of members to fetch.
type GetMembersParams struct {
	After Snowflake `urlparam:"after,omitempty"`
	Limit uint32    `urlparam:"limit,omitempty"` // 0 will fetch everyone
}

// GetMembers uses the GetGuildMembers endpoint iteratively until the your restriction/query params are met.
func (c *Client) GetMembers(guildID Snowflake, params *GetMembersParams, flags ...Flag) (members []*Member, err error) {
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

	var ms []*Member
	for {
		ms, err = c.getGuildMembers(guildID, &p, flags...)
		members = append(members, ms...)
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
func (c *Client) AddGuildMember(guildID, userID Snowflake, accessToken string, params *AddGuildMemberParams, flags ...Flag) (member *Member, err error) {
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
		Method:      httd.MethodPut,
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

// UpdateGuildMember [REST] Modify attributes of a guild member. Returns a 204 empty response on success.
// Fires a Guild Member Update Gateway event.
//  Method                  PATCH
//  Endpoint                /guilds/{guild.id}/members/{user.id}
//  Rate limiter            /guilds/{guild.id}/members
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#modify-guild-member
//  Reviewed                2018-08-17
//  Comment                 All parameters to this endpoint. are optional. When moving members to channels,
//                          the API user must have permissions to both connect to the channel and have the
//                          MOVE_MEMBERS permission.
func (c *Client) UpdateGuildMember(guildID, userID Snowflake, flags ...Flag) (builder *updateGuildMemberBuilder) {
	builder = &updateGuildMemberBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Member{}
	}
	builder.r.flags = flags
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
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
func (c *Client) AddGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodPut,
		Endpoint: endpoint.GuildMemberRole(guildID, userID, roleID),
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
func (c *Client) RemoveGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildMemberRole(guildID, userID, roleID),
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
func (c *Client) KickMember(guildID, userID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildMember(guildID, userID),
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
func (c *Client) GetGuildBans(id Snowflake, flags ...Flag) (bans []*Ban, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildBans(id),
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
func (c *Client) GetGuildBan(guildID, userID Snowflake, flags ...Flag) (ret *Ban, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildBan(guildID, userID),
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

// BanMember [REST] Add a guild ban, and optionally delete previous messages sent by the banned user. Requires
// the 'BAN_MEMBERS' permission. Returns a 204 empty response on success. Fires a Guild Ban Add Gateway event.
//  Method                  PUT
//  Endpoint                /guilds/{guild.id}/bans/{user.id}
//  Rate limiter            /guilds/{guild.id}/bans
//  Discord documentation   https://discordapp.com/developers/docs/resources/guild#create-guild-ban
//  Reviewed                2018-08-18
//  Comment                 -
func (c *Client) BanMember(guildID, userID Snowflake, params *BanMemberParams, flags ...Flag) (err error) {
	if params == nil {
		return errors.New("params was nil")
	}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodPut,
		Endpoint: endpoint.GuildBan(guildID, userID) + params.URLQueryString(),
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
func (c *Client) UnbanMember(guildID, userID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildBan(guildID, userID),
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
func (c *Client) EstimatePruneMembersCount(id Snowflake, days int, flags ...Flag) (estimate int, err error) {
	if id.IsZero() {
		return 0, errors.New("guildID can not be " + id.String())
	}
	params := pruneMembersParams{Days: days}
	if err = params.FindErrors(); err != nil {
		return 0, err
	}

	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildPrune(id) + params.URLQueryString(),
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
func (c *Client) PruneMembers(id Snowflake, days int, flags ...Flag) (err error) {
	params := pruneMembersParams{Days: days}
	if err = params.FindErrors(); err != nil {
		return err
	}

	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodPost,
		Endpoint: endpoint.GuildPrune(id) + params.URLQueryString(),
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
func (c *Client) GetGuildVoiceRegions(id Snowflake, flags ...Flag) (ret []*VoiceRegion, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildRegions(id),
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
func (c *Client) GetGuildInvites(id Snowflake, flags ...Flag) (ret []*Invite, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildInvites(id),
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
func (c *Client) GetGuildIntegrations(id Snowflake, flags ...Flag) (ret []*Integration, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildIntegrations(id),
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
func (c *Client) CreateGuildIntegration(guildID Snowflake, params *CreateGuildIntegrationParams, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPost,
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
func (c *Client) UpdateGuildIntegration(guildID, integrationID Snowflake, params *UpdateGuildIntegrationParams, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
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
func (c *Client) DeleteGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.GuildIntegration(guildID, integrationID),
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
func (c *Client) SyncGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodPost,
		Endpoint: endpoint.GuildIntegrationSync(guildID, integrationID),
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
func (c *Client) SetCurrentUserNick(id Snowflake, nick string, flags ...Flag) (newNick string, err error) {
	params := &updateCurrentUserNickParams{
		Nick: nick,
	}

	r := c.newRESTRequest(&httd.Request{
		Method:      httd.MethodPatch,
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
func (c *Client) GetGuildEmbed(guildID Snowflake, flags ...Flag) (embed *GuildEmbed, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildEmbed(guildID),
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
func (c *Client) UpdateGuildEmbed(guildID Snowflake, flags ...Flag) (builder *updateGuildEmbedBuilder) {
	builder = &updateGuildEmbedBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &GuildEmbed{}
	}
	builder.r.flags = flags
	builder.r.setup(c.cache, c.req, &httd.Request{
		Method:      httd.MethodPatch,
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
func (c *Client) GetGuildVanityURL(guildID Snowflake, flags ...Flag) (ret *PartialInvite, err error) {
	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.GuildVanityURL(guildID),
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

// KickFromVoice kicks member out of voice channel. Assuming he/she/it is in one.
func (b *updateGuildMemberBuilder) KickFromVoice() *updateGuildMemberBuilder {
	b.r.param("channel_id", 0)
	return b
}

// DeleteNick removes nickname for user. Requires permission MANAGE_NICKNAMES
func (b *updateGuildMemberBuilder) DeleteNick() *updateGuildMemberBuilder {
	b.r.param("nick", "")
	return b
}
