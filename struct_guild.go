package disgord

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
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

func NewGuildFromJSON(data []byte) (guild *Guild) {
	guild = NewGuild()
	err := json.Unmarshal(data, guild)
	if err != nil {
		panic(err)
	}

	return guild
}

func NewPartialGuild(ID Snowflake) (guild *Guild) {
	guild = NewGuild()
	guild.ID = ID
	guild.Unavailable = true

	return
}

func NewGuildFromUnavailable(gu *GuildUnavailable) *Guild {
	return NewPartialGuild(gu.ID)
}

func NewGuildUnavailable(ID Snowflake) *GuildUnavailable {
	gu := &GuildUnavailable{
		ID:          ID,
		Unavailable: true,
	}

	return gu
}

type GuildUnavailable struct {
	ID           Snowflake `json:"id"`
	Unavailable  bool      `json:"unavailable"` // ?*|
	sync.RWMutex `json:"-"`
}

type GuildInterface interface {
	Channel(ID Snowflake)
}

// if loading is deactivated, then check state, then do a request.
// if loading is activated, check state only.
// type Members interface {
// 	Member(userID snowflake.Snowflake) *Member
// 	MembersWithName( /*username*/ name string) map[snowflake.Snowflake]*Member
// 	MemberByUsername( /*username#discriminator*/ username string) *Member
// 	MemberByAlias(alias string) *Member
// 	EverythingInMemory() bool
// }

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
// TODO: lazyload everything
// reviewed: 2018-08-25
type PartialGuild = Guild
type Guild struct {
	sync.RWMutex `json:"-"`

	ID                          Snowflake                     `json:"id"`
	ApplicationID               *Snowflake                    `json:"application_id"` //   |?
	Name                        string                        `json:"name"`
	Icon                        *string                       `json:"icon"`            //  |?, icon hash
	Splash                      *string                       `json:"splash"`          //  |?, image hash
	Owner                       bool                          `json:"owner,omitempty"` // ?|
	OwnerID                     Snowflake                     `json:"owner_id"`
	Permissions                 uint64                        `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
	Region                      string                        `json:"region"`
	AfkChannelID                *Snowflake                    `json:"afk_channel_id"` // |?
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
	WidgetChannelID             Snowflake                     `json:"widget_channel_id,omit_empty"` //   |
	SystemChannelID             *Snowflake                    `json:"system_channel_id,omitempty"`  //   |?

	// JoinedAt must be a pointer, as we can't hide non-nil structs
	JoinedAt    *Timestamp      `json:"joined_at,omitempty"`    // ?*|
	Large       bool            `json:"large,omitempty"`        // ?*|
	Unavailable bool            `json:"unavailable"`            // ?*| omitempty?
	MemberCount uint            `json:"member_count,omitempty"` // ?*|
	VoiceStates []*VoiceState   `json:"voice_states,omitempty"` // ?*|
	Members     []*Member       `json:"members,omitempty"`      // ?*|
	Channels    []*Channel      `json:"channels,omitempty"`     // ?*|
	Presences   []*UserPresence `json:"presences,omitempty"`    // ?*|

	//highestSnowflakeAmoungMembers Snowflake
}

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

//func (g *Guild) EverythingInMemory() bool {
//	return false
//}

// Compare two guild objects
// TODO: remove?
//func (g *Guild) Compare(other *Guild) bool {
//	// TODO: this is shit..
//	return (g == nil && other == nil) || (other != nil && g.ID == other.ID)
//}

// func (g *Guild) UnmarshalJSON(data []byte) (err error) {
// 	return json.Unmarshal(data, &g.guildJSON)
// }

// TODO: fix copying of mutex lock
func (g *Guild) MarshalJSON() ([]byte, error) {
	g.Lock()
	defer g.Unlock()
	// TODO: check for dead locks

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
	g.Lock()
	defer g.Unlock()

	g.Channels = append(g.Channels, c)
	g.sortChannels()

	return nil
}

func (g *Guild) DeleteChannel(c *Channel) error {
	return g.DeleteChannelByID(c.ID)
}
func (g *Guild) DeleteChannelByID(ID Snowflake) error {
	g.Lock()
	defer g.Unlock()

	index := -1
	for i, c := range g.Channels {
		if c.ID == ID {
			index = i
		}
	}

	if index == -1 {
		return errors.New("channel with Snowflake{" + ID.String() + "} does not exist in cache")
	}

	// delete the entry
	g.Channels[index] = g.Channels[len(g.Channels)-1]
	g.Channels[len(g.Channels)-1] = nil
	g.Channels = g.Channels[:len(g.Channels)-1]

	g.sortChannels()

	return nil
}

func (g *Guild) AddMember(member *Member) error {
	g.Lock()
	defer g.Unlock()

	// TODO: implement sorting for faster searching later
	g.Members = append(g.Members, member)

	return nil
}

func (g *Guild) AddRole(role *Role) error {
	g.Lock()
	defer g.Unlock()

	// TODO: implement sorting for faster searching later
	role.guildID = g.ID
	g.Roles = append(g.Roles, role)

	return nil
}

// Member return a member by his/her userid
func (g *Guild) Member(id Snowflake) (*Member, error) {
	g.RLock()
	defer g.RUnlock()

	for _, member := range g.Members {
		if member.User.ID == id {
			return member, nil
		}
	}

	return nil, errors.New("member not found in guild")
}

// MembersByName retrieve a slice of members with same username or nickname
func (g *Guild) MembersByName(name string) (members []*Member) {
	g.RLock()
	defer g.RUnlock()

	for _, member := range g.Members {
		if member.Nick == name || member.User.Username == name {
			members = append(members, member)
		}
	}

	return
}

// Role retrieve a role based on role id
func (g *Guild) Role(id Snowflake) (role *Role, err error) {
	g.RLock()
	defer g.RUnlock()

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

func (g *Guild) DeleteRoleByID(ID Snowflake) {
	g.Lock()
	defer g.Unlock()

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
	g.RLock()
	defer g.RUnlock()

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

func (g *Guild) Channel(id Snowflake) (*Channel, error) {
	g.RLock()
	defer g.RUnlock()

	for _, channel := range g.Channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, errors.New("channel not found in guild")
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
// 	g.Lock()
// 	g.Presences = append(g.Presences, p) // TODO: update the user pointer?
// 	g.Unlock()
// }

// Clear all the pointers
// func (g *Guild) Clear() {
// 	g.Lock() // what if another process tries to read this, but awais while locked for clearing?
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

func (g *Guild) DeepCopy() (copy interface{}) {
	copy = NewGuild()
	g.CopyOverTo(copy)

	return
}

func (g *Guild) CopyOverTo(other interface{}) (err error) {
	var guild *Guild
	var valid bool
	if guild, valid = other.(*Guild); !valid {
		err = NewErrorUnsupportedType("argument given is not a *Guild type")
		return
	}

	g.RLock()
	guild.Lock()

	// TODO-guild: handle string pointers
	guild.ID = g.ID
	guild.Name = g.Name
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
	guild.Large = g.Large
	guild.Unavailable = g.Unavailable
	guild.MemberCount = g.MemberCount

	// pointers
	if g.ApplicationID != nil {
		id := *g.ApplicationID
		guild.ApplicationID = &id
	}
	if g.Icon != nil {
		icon := *g.Icon
		guild.Icon = &icon
	}
	if g.Splash != nil {
		splash := *g.Splash
		guild.Splash = &splash
	}
	if g.AfkChannelID != nil {
		channel := *g.AfkChannelID
		guild.AfkChannelID = &channel
	}
	if g.JoinedAt != nil {
		joined := *g.JoinedAt
		guild.JoinedAt = &joined
	}

	for _, roleP := range g.Roles {
		guild.Roles = append(guild.Roles, roleP.DeepCopy().(*Role))
	}
	for _, emojiP := range g.Emojis {
		guild.Emojis = append(guild.Emojis, emojiP.DeepCopy().(*Emoji))
	}
	guild.Features = g.Features

	for _, vsP := range g.VoiceStates {
		guild.VoiceStates = append(guild.VoiceStates, vsP.DeepCopy().(*VoiceState))
	}
	for _, memberP := range g.Members {
		guild.Members = append(guild.Members, memberP.DeepCopy().(*Member))
	}
	for _, channelP := range g.Channels {
		guild.Channels = append(guild.Channels, channelP.DeepCopy().(*Channel))
	}
	for _, presenceP := range g.Presences {
		guild.Presences = append(guild.Presences, presenceP.DeepCopy().(*UserPresence))
	}

	g.RUnlock()
	guild.Unlock()

	return
}

// saveToDiscord creates a new Guild if ID is empty or updates an existing one
func (g *Guild) saveToDiscord(session Session) (err error) {
	return errors.New("not implemented")
}
func (g *Guild) deleteFromDiscord(session Session) (err error) {
	return errors.New("not implemented")
}

// --------------
// Ban https://discordapp.com/developers/docs/resources/guild#ban-object
type Ban struct {
	Reason string `json:"reason"`
	User   *User  `json:"user"`
}

// ------------
// GuildEmbed https://discordapp.com/developers/docs/resources/guild#guild-embed-object
type GuildEmbed struct {
	Enabled   bool      `json:"enabled"`
	ChannelID Snowflake `json:"channel_id"`
}

// -------
// Integration https://discordapp.com/developers/docs/resources/guild#integration-object
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

// IntegrationAccount https://discordapp.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id"`   // id of the account
	Name string `json:"name"` // name of the account
}

// -------

// Member https://discordapp.com/developers/docs/resources/guild#guild-member-object
type Member struct {
	sync.RWMutex `json:"-"`

	GuildID  Snowflake   `json:"guild_id,omitempty"`
	User     *User       `json:"user"`
	Nick     string      `json:"nick,omitempty"` // ?|
	Roles    []Snowflake `json:"roles"`
	JoinedAt Timestamp   `json:"joined_at,omitempty"`
	Deaf     bool        `json:"deaf"`
	Mute     bool        `json:"mute"`
}

func (m *Member) String() string {
	return "member{user:" + m.User.Username + ", nick:" + m.Nick + ", ID:" + m.User.ID.String() + "}"
}
func (m *Member) Mention() string {
	return m.User.Mention()
}
func (m *Member) DeepCopy() (copy interface{}) {
	copy = &Member{}
	m.CopyOverTo(copy)

	return
}
func (m *Member) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var member *Member
	if member, ok = other.(*Member); !ok {
		err = NewErrorUnsupportedType("given interface{} was not of type *Member")
		return
	}

	m.RLock()
	member.Lock()

	member.GuildID = m.GuildID
	member.User = m.User.DeepCopy().(*User)
	member.Nick = m.Nick
	member.Roles = m.Roles
	member.JoinedAt = m.JoinedAt
	member.Deaf = m.Deaf
	member.Mute = m.Mute

	m.RUnlock()
	member.Unlock()

	return
}

// TODO
// func (m *Member) Clear() Snowflake {
// 	// do i want to delete user?.. what if there is a PM?
// 	// Check for user id in DM's
// 	// or.. since the user object is sent on channel_create events, the user can be reintialized when needed.
// 	// but should be properly removed from other arrays.
// 	m.User.Clear()
// 	id := m.User.ID
// 	m.User = nil
//
// 	// use this Snowflake to check in other places. To avoid pointing to random memory spaces
// 	return id
// }

// TODO
// func (m *Member) Update(new *Member) (err error) {
// 	if m.User.ID != new.User.ID || m.GuildID != new.GuildID {
// 		err = errors.New("cannot update user when the new struct has a different Snowflake")
// 		return
// 	}
// 	// make sure that new is not the same pointer!
// 	if m == new {
// 		err = errors.New("cannot update user when the new struct points to the same memory space")
// 		return
// 	}
//
// 	m.Lock()
// 	new.RLock()
// 	m.Nick = new.Nick
// 	m.Roles = new.Roles
// 	m.JoinedAt = new.JoinedAt
// 	m.Deaf = new.Deaf
// 	m.Mute = new.Mute
// 	new.RUnlock()
// 	m.Unlock()
//
// 	return
// }

type GuildPruneCount struct {
	Pruned int `json:"pruned"`
}
