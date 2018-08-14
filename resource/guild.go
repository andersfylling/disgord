package resource

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"

	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/snowflake"
)

// consts inspired by: https://github.com/bwmarrin/discordgo/blob/master/structs.go

// Constants for the different bit offsets of text channel permissions
const (
	ReadMessagesPermission = 1 << (iota + 10)
	SendMessagesPermission
	SendTTSMessagesPermission
	ManageMessagesPermission
	EmbedLinksPermission
	AttachFilesPermission
	ReadMessageHistoryPermission
	MentionEveryonePermission
	UseExternalEmojisPermission
)

// Constants for the different bit offsets of voice permissions
const (
	VoiceConnectPermission = 1 << (iota + 20)
	VoiceSpeakPermission
	VoiceMuteMembersPermission
	VoiceDeafenMembersPermission
	VoiceMoveMembersPermission
	VoiceUseVADPermission
)

// Constants for general management.
const (
	ChangeNicknamePermission = 1 << (iota + 26)
	ManageNicknamesPermission
	ManageRolesPermission
	ManageWebhooksPermission
	ManageEmojisPermission
)

// Constants for the different bit offsets of general permissions
const (
	CreateInstantInvitePermission = 1 << iota
	KickMembersPermission
	BanMembersPermission
	AdministratorPermission
	ManageChannelsPermission
	ManageServerPermission
	AddReactionsPermission
	ViewAuditLogsPermission

	AllTextPermission = ReadMessagesPermission |
		SendMessagesPermission |
		SendTTSMessagesPermission |
		ManageMessagesPermission |
		EmbedLinksPermission |
		AttachFilesPermission |
		ReadMessageHistoryPermission |
		MentionEveryonePermission
	AllVoicePermission = VoiceConnectPermission |
		VoiceSpeakPermission |
		VoiceMuteMembersPermission |
		VoiceDeafenMembersPermission |
		VoiceMoveMembersPermission |
		VoiceUseVADPermission
	AllChannelPermission = AllTextPermission |
		AllVoicePermission |
		CreateInstantInvitePermission |
		ManageRolesPermission |
		ManageChannelsPermission |
		AddReactionsPermission |
		ViewAuditLogsPermission
	AllPermission = AllChannelPermission |
		KickMembersPermission |
		BanMembersPermission |
		ManageServerPermission |
		AdministratorPermission
)

func NewGuild() *Guild {
	return &Guild{}
}

func NewGuildFromJSON(data []byte) *Guild {
	guild := &Guild{
		Roles:       []*Role{},
		Emojis:      []*Emoji{},
		Features:    []string{},
		VoiceStates: []*VoiceState{},
		Members:     []*Member{},
		Channels:    []*Channel{},
		Presences:   []*UserPresence{},
	}
	err := json.Unmarshal(data, guild)
	if err != nil {
		panic(err)
	}

	return guild
}

func NewPartialGuild(ID snowflake.ID) *Guild {
	return &Guild{
		ID:          ID,
		Unavailable: true,
		Roles:       []*Role{},
		Emojis:      []*Emoji{},
		Features:    []string{},
		VoiceStates: []*VoiceState{},
		Members:     []*Member{},
		Channels:    []*Channel{},
		Presences:   []*UserPresence{},
	}
}

func NewGuildFromUnavailable(gu *GuildUnavailable) *Guild {
	return NewPartialGuild(gu.ID)
}

func NewGuildUnavailable(ID snowflake.ID) *GuildUnavailable {
	gu := &GuildUnavailable{
		ID:          ID,
		Unavailable: true,
	}

	return gu
}

type GuildUnavailable struct {
	ID           snowflake.ID `json:"id"`
	Unavailable  bool         `json:"unavailable"` // ?*|
	sync.RWMutex `json:"-"`
}

type GuildInterface interface {
	Channel(ID snowflake.ID)
}

// if loading is deactivated, then check state, then do a request.
// if loading is activated, check state only.
// type Members interface {
// 	Member(userID snowflake.ID) *Member
// 	MembersWithName( /*username*/ name string) map[snowflake.ID]*Member
// 	MemberByUsername( /*username#discriminator*/ username string) *Member
// 	MemberByAlias(alias string) *Member
// 	EverythingInMemory() bool
// }

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
// TODO: lazyload everything
type PartialGuild = Guild
type Guild struct {
	ID                          snowflake.ID                          `json:"id"`
	ApplicationID               *snowflake.ID                         `json:"application_id"` //   |?
	Name                        string                                `json:"name"`
	Icon                        *string                               `json:"icon"`            //  |?, icon hash
	Splash                      *string                               `json:"splash"`          //  |?, image hash
	Owner                       bool                                  `json:"owner,omitempty"` // ?|
	OwnerID                     snowflake.ID                          `json:"owner_id"`
	Permissions                 uint64                                `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
	Region                      string                                `json:"region"`
	AfkChannelID                snowflake.ID                          `json:"afk_channel_id"`
	AfkTimeout                  uint                                  `json:"afk_timeout"`
	EmbedEnabled                bool                                  `json:"embed_enabled"`
	EmbedChannelID              snowflake.ID                          `json:"embed_channel_id"`
	VerificationLevel           discord.VerificationLvl               `json:"verification_level"`
	DefaultMessageNotifications discord.DefaultMessageNotificationLvl `json:"default_message_notifications"`
	ExplicitContentFilter       discord.ExplicitContentFilterLvl      `json:"explicit_content_filter"`
	MFALevel                    discord.MFALvl                        `json:"mfa_level"`
	WidgetEnabled               bool                                  `json:"widget_enabled"`    //   |
	WidgetChannelID             snowflake.ID                          `json:"widget_channel_id"` //   |
	Roles                       []*Role                               `json:"roles"`
	Emojis                      []*Emoji                              `json:"emojis"`
	Features                    []string                              `json:"features"`
	SystemChannelID             *snowflake.ID                         `json:"system_channel_id,omitempty"` //   |?

	// JoinedAt must be a pointer, as we can't hide non-nil structs
	JoinedAt       *discord.Timestamp `json:"joined_at,omitempty"`    // ?*|
	Large          bool               `json:"large,omitempty"`        // ?*|
	Unavailable    bool               `json:"unavailable"`            // ?*|
	MemberCount    uint               `json:"member_count,omitempty"` // ?*|
	VoiceStates    []*VoiceState      `json:"voice_states,omitempty"` // ?*|
	Members        []*Member          `json:"members,omitempty"`      // ?*|
	Channels       []*Channel         `json:"channels,omitempty"`     // ?*|
	Presences      []*UserPresence    `json:"presences,omitempty"`    // ?*|
	PresencesMutex sync.RWMutex       `json:"-"`

	mu sync.RWMutex `json:"-"`
}

//func (g *Guild) EverythingInMemory() bool {
//	return false
//}

// Compare two guild objects
func (g *Guild) Compare(other *Guild) bool {
	// TODO: this is shit..
	return (g == nil && other == nil) || (other != nil && g.ID == other.ID)
}

// func (g *Guild) UnmarshalJSON(data []byte) (err error) {
// 	return json.Unmarshal(data, &g.guildJSON)
// }

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

func (g *Guild) AddChannel(c *Channel) error {
	g.Channels = append(g.Channels, c)
	g.sortChannels()

	return nil
}

func (g *Guild) DeleteChannel(c *Channel) error {
	return g.DeleteChannelByID(c.ID)
}
func (g *Guild) DeleteChannelByID(ID snowflake.ID) error {
	index := -1
	for i, c := range g.Channels {
		if c.ID == ID {
			index = i
		}
	}

	if index == -1 {
		return errors.New("channel with ID{" + ID.String() + "} does not exist in cache")
	}

	// delete the entry
	g.Channels[index] = g.Channels[len(g.Channels)-1]
	g.Channels[len(g.Channels)-1] = nil
	g.Channels = g.Channels[:len(g.Channels)-1]

	g.sortChannels()

	return nil
}

func (g *Guild) AddMember(member *Member) error {
	// TODO: implement sorting for faster searching later
	g.Members = append(g.Members, member)

	return nil
}

func (g *Guild) AddRole(role *Role) error {
	// TODO: implement sorting for faster searching later
	g.Roles = append(g.Roles, role)

	return nil
}

// Member return a member by his/her userid
func (g *Guild) Member(id snowflake.ID) (*Member, error) {
	for _, member := range g.Members {
		if member.User.ID == id {
			return member, nil
		}
	}

	return nil, errors.New("member not found in guild")
}

// MemberByName retrieve a slice of members with same username or nickname
func (g *Guild) MemberByName(name string) ([]*Member, error) {
	var members []*Member
	for _, member := range g.Members {
		if member.Nick == name || member.User.Username == name {
			members = append(members, member)
		}
	}

	if len(members) == 0 {
		return nil, errors.New("no members with that nick or username was found in guild")
	}

	return members, nil
}

// Role retrieve a role based on role id
func (g *Guild) Role(id snowflake.ID) (*Role, error) {
	for _, role := range g.Roles {
		if role.ID == id {
			return role, nil
		}
	}

	return nil, errors.New("role not found in guild")
}

func (g *Guild) UpdateRole(r *Role) {
	for _, role := range g.Roles {
		if role.ID == r.ID {
			*role = *r
			break
		}
	}
}
func (g *Guild) DeleteRoleByID(ID snowflake.ID) {
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

// RoleByTitle retrieves a slice of roles with same name
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

func (g *Guild) Channel(id snowflake.ID) (*Channel, error) {
	for _, channel := range g.Channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, errors.New("channel not found in guild")
}

func (g *Guild) UpdatePresence(p *UserPresence) {
	g.PresencesMutex.RLock()
	index := -1
	for i, presence := range g.Presences {
		if presence.User.ID == p.User.ID {
			index = i
			break
		}
	}
	g.PresencesMutex.RUnlock()

	if index != -1 {
		// update
		return
	}

	// otherwise add
	g.PresencesMutex.Lock()
	g.Presences = append(g.Presences, p) // TODO: update the user pointer?
	g.PresencesMutex.Unlock()
}

// Update update the reference content
func (g *Guild) Update(new *Guild) {
	// must have same ID
	if g.ID != new.ID {
		return
	}

	// must not be the same pointer as it causes a deadlock
	if g == new {
		return
	}

	g.mu.Lock()
	new.mu.RLock()

	// if it's a unavailable guild object, don't update the remaining fields
	if new.Unavailable {
		g.Unavailable = true
	} else {
		// normal update
		// TODO
	}

	g.mu.Unlock()
	new.mu.RUnlock()
}

// Clear all the pointers
func (g *Guild) Clear() {
	g.mu.Lock() // what if another process tries to read this, but awais while locked for clearing?
	defer g.mu.Unlock()

	g.ApplicationID = nil
	//g.Icon = nil // should this be cleared?
	//g.Splash = nil // should this be cleared?

	for _, r := range g.Roles {
		r.Clear()
		r = nil
	}
	g.Roles = nil

	for _, e := range g.Emojis {
		e.Clear()
		e = nil
	}
	g.Emojis = nil

	g.SystemChannelID = nil
	g.JoinedAt = nil

	for _, vst := range g.VoiceStates {
		vst.Clear()
		vst = nil
	}
	g.VoiceStates = nil

	var deletedUsers []snowflake.ID
	for _, m := range g.Members {
		deletedUsers = append(deletedUsers, m.Clear())
		m = nil
	}
	g.Members = nil

	for _, c := range g.Channels {
		c.Clear()
		c = nil
	}
	g.Channels = nil

	for _, p := range g.Presences {
		p.Clear()
		p = nil
	}
	g.Presences = nil

}

func (g *Guild) DeepCopy() *Guild {
	guild := NewGuild()

	g.mu.RLock()

	// TODO-guild: handle string pointers
	guild.ID = g.ID
	guild.ApplicationID = g.ApplicationID
	guild.Name = g.Name
	guild.Icon = g.Icon
	guild.Splash = g.Splash
	guild.Owner = g.Owner
	guild.OwnerID = g.OwnerID
	guild.Permissions = g.Permissions
	guild.Region = g.Region
	guild.AfkChannelID = g.AfkChannelID
	guild.AfkTimeout = g.AfkTimeout
	guild.EmbedEnabled = g.EmbedEnabled
	guild.EmbedChannelID = g.EmbedChannelID
	guild.VerificationLevel = g.VerificationLevel
	guild.DefaultMessageNotifications = g.DefaultMessageNotifications
	guild.ExplicitContentFilter = g.ExplicitContentFilter
	guild.MFALevel = g.MFALevel
	guild.WidgetEnabled = g.WidgetEnabled
	guild.WidgetChannelID = g.WidgetChannelID
	guild.SystemChannelID = g.SystemChannelID
	guild.JoinedAt = g.JoinedAt
	guild.Large = g.Large
	guild.Unavailable = g.Unavailable
	guild.MemberCount = g.MemberCount
	guild.PresencesMutex = g.PresencesMutex

	// handle deep copy of slices
	//TODO-guild: implement deep copying for fields
	guild.Roles = g.Roles
	guild.Emojis = g.Emojis
	guild.Features = g.Features

	guild.VoiceStates = g.VoiceStates
	guild.Members = g.Members
	guild.Channels = g.Channels
	guild.Presences = g.Presences

	g.mu.RUnlock()

	return guild
}

//--------------
type GuildBan struct {
	Reason *string `json:"reason"`
	User   *User   `json:"user"`
}

//------------
type GuildEmbed struct {
	Enabled   bool         `json:"enabled"`
	ChannelID snowflake.ID `json:"channel_id"`
}

// -------
type GuildIntegration struct {
}

type GuildIntegrationAccount struct {
}

// -------

// Member ...
type Member struct {
	GuildID  snowflake.ID      `json:"guild_id,omitempty"`
	User     *User             `json:"user"`
	Nick     string            `json:"nick,omitempty"` // ?|
	Roles    []snowflake.ID    `json:"roles"`
	JoinedAt discord.Timestamp `json:"joined_at,omitempty"`
	Deaf     bool              `json:"deaf"`
	Mute     bool              `json:"mute"`

	sync.RWMutex `json:"-"`
}

func (m *Member) Clear() snowflake.ID {
	// do i want to delete user?.. what if there is a PM?
	// Check for user id in DM's
	// or.. since the user object is sent on channel_create events, the user can be reintialized when needed.
	// but should be properly removed from other arrays.
	m.User.Clear()
	id := m.User.ID
	m.User = nil

	// use this ID to check in other places. To avoid pointing to random memory spaces
	return id
}

func (m *Member) Update(new *Member) (err error) {
	if m.User.ID != new.User.ID || m.GuildID != new.GuildID {
		err = errors.New("cannot update user when the new struct has a different ID")
		return
	}
	// make sure that new is not the same pointer!
	if m == new {
		err = errors.New("cannot update user when the new struct points to the same memory space")
		return
	}

	m.Lock()
	new.RLock()
	m.Nick = new.Nick
	m.Roles = new.Roles
	m.JoinedAt = new.JoinedAt
	m.Deaf = new.Deaf
	m.Mute = new.Mute
	new.RUnlock()
	m.Unlock()

	return
}

// --------------

type Role struct {
	ID          snowflake.ID `json:"id"`
	Name        string       `json:"name"`
	Managed     bool         `json:"managed"`
	Mentionable bool         `json:"mentionable"`
	Hoist       bool         `json:"hoist"`
	Color       int          `json:"color"`
	Position    int          `json:"position"`
	Permissions uint64       `json:"permissions"`
}

func NewRole() *Role {
	return &Role{}
}

func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

func (r *Role) Clear() {

}

const (
	EndpointGuild = "/guilds/"
)
