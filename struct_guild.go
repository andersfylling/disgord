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
	ID           Snowflake `json:"id"`
	Unavailable  bool      `json:"unavailable"` // ?*|
	sync.RWMutex `json:"-"`
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
type PartialGuild = Guild

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
// reviewed: 2018-08-25
type Guild struct {
	sync.RWMutex `json:"-"`

	ID                          Snowflake                     `json:"id"`
	ApplicationID               Snowflake                     `json:"application_id"` //   |?
	Name                        string                        `json:"name"`
	Icon                        *string                       `json:"icon"`            //  |?, icon hash
	Splash                      *string                       `json:"splash"`          //  |?, image hash
	Owner                       bool                          `json:"owner,omitempty"` // ?|
	OwnerID                     Snowflake                     `json:"owner_id"`
	Permissions                 uint64                        `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
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

func (g *Guild) copyOverToCache(other interface{}) (err error) {
	guild := other.(*Guild)

	g.RLock()
	guild.Lock()

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

	// pointers
	if !g.ApplicationID.Empty() {
		guild.ApplicationID = g.ApplicationID
	}
	if g.Splash != nil {
		splash := *g.Splash
		guild.Splash = &splash
	}
	if g.Icon != nil {
		icon := *g.Icon
		guild.Icon = &icon
	}
	if !g.AfkChannelID.Empty() {
		guild.AfkChannelID = g.AfkChannelID
	}
	if !g.SystemChannelID.Empty() {
		guild.SystemChannelID = g.SystemChannelID
	}
	if g.JoinedAt != nil {
		joined := *g.JoinedAt
		guild.JoinedAt = &joined
	}

	g.RUnlock()
	guild.Unlock()

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

// AddChannel adds a channel to the Guild object. Note that this method does not interact with Discord.
func (g *Guild) AddChannel(c *Channel) error {
	g.Lock()
	defer g.Unlock()

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
	g.Lock()
	defer g.Unlock()

	for _, member := range members {
		g.addMember(member)
	}
}

// AddMember adds a member to the Guild object. Note that this method does not interact with Discord.
func (g *Guild) AddMember(member *Member) error {
	g.Lock()
	defer g.Unlock()

	return g.addMember(member)
}

// LoadAllMembers fetches all teh members for this guild from the Discord REST API
func (g *Guild) LoadAllMembers(session Session) (err error) {
	g.Lock()
	defer g.Unlock()

	// TODO-1: check cache
	// TODO-2: what if members have already been loaded? use Guild.MembersCount?

	var lastCount = 1000
	var failsafe bool
	// TODO-3: failsafe is set when the number of users returned is less
	// than 1,000 two times
	highestSnowflake := NewSnowflake(0)

	for {
		if lastCount == 0 || failsafe {
			break
		}
		var members []*Member
		members, err = session.GetGuildMembers(g.ID, highestSnowflake, 1000)
		if err != nil {
			return
		}

		for _, member := range members {
			g.addMember(member)

			s := member.User.ID
			if s > highestSnowflake {
				highestSnowflake = s
			}
		}

		lastCount = len(members)
	}

	return nil
}

// AddRole adds a role to the Guild object. Note that this does not interact with Discord.
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

// DeleteRoleByID remove a role from the guild struct
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

// RoleByName retrieves a slice of roles with same name
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

// Channel get a guild channel given it's ID
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

// Emoji get a guild emoji by it's ID
func (g *Guild) Emoji(id Snowflake) (emoji *Emoji, err error) {
	g.RLock()
	defer g.RUnlock()

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

	g.RLock()
	guild.Lock()

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

	// pointers
	if !g.ApplicationID.Empty() {
		guild.ApplicationID = g.ApplicationID
	}
	if g.Splash != nil {
		splash := *g.Splash
		guild.Splash = &splash
	}
	if g.Icon != nil {
		icon := *g.Icon
		guild.Icon = &icon
	}
	if !g.AfkChannelID.Empty() {
		guild.AfkChannelID = g.AfkChannelID
	}
	if !g.SystemChannelID.Empty() {
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
	sync.RWMutex `json:"-"`

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

	b.RLock()
	ban.Lock()

	ban.Reason = b.Reason

	if b.User != nil {
		ban.User = b.User.DeepCopy().(*User)
	}

	b.RUnlock()
	ban.Unlock()

	return
}

// ------------

// GuildEmbed https://discordapp.com/developers/docs/resources/guild#guild-embed-object
type GuildEmbed struct {
	sync.RWMutex `json:"-"`

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

	e.RLock()
	embed.Lock()

	embed.Enabled = e.Enabled
	embed.ChannelID = e.ChannelID

	e.RUnlock()
	embed.Unlock()

	return
}

// -------

// Integration https://discordapp.com/developers/docs/resources/guild#integration-object
type Integration struct {
	sync.RWMutex `json:"-"`

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

	i.RLock()
	integration.Lock()

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

	i.RUnlock()
	integration.Unlock()

	return
}

// IntegrationAccount https://discordapp.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	sync.RWMutex `json:"-"`

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

	i.RLock()
	account.Lock()

	account.ID = i.ID
	account.Name = i.Name

	i.RUnlock()
	account.Unlock()
	return
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

	// used for caching
	userID Snowflake
}

func (m *Member) String() string {
	return "member{user:" + m.User.Username + ", nick:" + m.Nick + ", ID:" + m.User.ID.String() + "}"
}

// Mention creates a string which is parsed into a member mention on Discord GUI's
func (m *Member) Mention() string {
	return m.User.MentionNickname()
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

// GuildPruneCount ...
type GuildPruneCount struct {
	Pruned int `json:"pruned"`
}
