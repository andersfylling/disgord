package guild

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/emoji"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/disgord/voice"
	"github.com/andersfylling/snowflake"
)

func NewGuild() *Guild {
	return &Guild{}
}

func NewGuildFromJSON(data []byte) *Guild {
	guild := &Guild{}
	err := json.Unmarshal(data, guild)
	if err != nil {
		panic(err)
	}

	return guild
}

func NewGuildFromUnavailable(gu *Unavailable) *Guild {
	g := &Guild{
		ID:          gu.ID,
		Unavailable: gu.Unavailable,
	}

	return g
}

func NewGuildUnavailable(ID snowflake.ID) *Unavailable {
	gu := &Unavailable{
		ID:          ID,
		Unavailable: true,
	}

	return gu
}

type Unavailable struct {
	ID           snowflake.ID `json:"id"`
	Unavailable  bool         `json:"unavailable"` // ?*|
	sync.RWMutex `json:"-"`
}

type GuildInterface interface {
	Channel(ID snowflake.ID)
}

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
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
	Emojis                      []*emoji.Emoji                        `json:"emojis"`
	Features                    []string                              `json:"features"`
	SystemChannelID             *snowflake.ID                         `json:"system_channel_id,omitempty"` //   |?

	// JoinedAt must be a pointer, as we can't hide non-nil structs
	JoinedAt       *discord.Timestamp `json:"joined_at,omitempty"`    // ?*|
	Large          bool               `json:"large,omitempty"`        // ?*|
	Unavailable    bool               `json:"unavailable"`            // ?*|
	MemberCount    uint               `json:"member_count,omitempty"` // ?*|
	VoiceStates    []*voice.State     `json:"voice_states,omitempty"` // ?*|
	Members        []*Member          `json:"members,omitempty"`      // ?*|
	Channels       []*channel.Channel `json:"channels,omitempty"`     // ?*|
	Presences      []*user.Presence   `json:"presences,omitempty"`    // ?*|
	PresencesMutex sync.RWMutex       `json:"-"`

	sync.RWMutex `json:"-"`
}

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

func (g *Guild) AddChannel(c *channel.Channel) error {
	g.Channels = append(g.Channels, c)
	g.sortChannels()

	return nil
}

func (g *Guild) DeleteChannel(c *channel.Channel) error {
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

func (g *Guild) Channel(id snowflake.ID) (*channel.Channel, error) {
	for _, channel := range g.Channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, errors.New("channel not found in guild")
}

func (g *Guild) UpdatePresence(p *user.Presence) {
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

	g.Lock()
	new.RLock()

	// if it's a unavailable guild object, don't update the remaining fields
	if new.Unavailable {
		g.Unavailable = true
	} else {
		// normal update
		// TODO
	}

	g.Unlock()
	new.RUnlock()
}

// Clear all the pointers
func (g *Guild) Clear() {
	g.Lock() // what if another process tries to read this, but awais while locked for clearing?
	defer g.Unlock()

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

	deletedUsers := []snowflake.ID{}
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
