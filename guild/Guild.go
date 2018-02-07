package guild

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/emoji"
	"github.com/andersfylling/disgord/lvl"
	"github.com/andersfylling/disgord/voice"
	"github.com/andersfylling/snowflake"
)

func NewGuildFromJSON(data []byte) *Guild {
	guild := &Guild{}
	err := json.Unmarshal(data, guild)
	if err != nil {
		panic(err)
	}

	return guild
}

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
type guildJSON struct {
	ID                          snowflake.ID                   `json:"id"`
	ApplicationID               *snowflake.ID                  `json:"application_id"` //   |?
	Name                        string                         `json:"name"`
	Icon                        *string                        `json:"icon"`            //  |?, icon hash
	Splash                      *string                        `json:"splash"`          //  |?, image hash
	Owner                       bool                           `json:"owner,omitempty"` // ?|
	OwnerID                     snowflake.ID                   `json:"owner_id"`
	Permissions                 uint64                         `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
	Region                      string                         `json:"region"`
	AfkChannelID                snowflake.ID                   `json:"afk_channel_id"`
	AfkTimeout                  uint                           `json:"afk_timeout"`
	EmbedEnabled                bool                           `json:"embed_enabled"`
	EmbedChannelID              snowflake.ID                   `json:"embed_channel_id"`
	VerificationLevel           lvl.Verification               `json:"verification_level"`
	DefaultMessageNotifications lvl.DefaultMessageNotification `json:"default_message_notifications"`
	ExplicitContentFilter       lvl.ExplicitContentFilter      `json:"explicit_content_filter"`
	MFALevel                    lvl.MFA                        `json:"mfa_level"`
	WidgetEnabled               bool                           `json:"widget_enabled"`    //   |
	WidgetChannelID             snowflake.ID                   `json:"widget_channel_id"` //   |
	Roles                       []*discord.Role                `json:"roles"`
	Emojis                      []*emoji.Emoji                 `json:"emojis"`
	Features                    []string                       `json:"features"`
	SystemChannelID             *snowflake.ID                  `json:"system_channel_id,omitempty"` //   |?

	// JoinedAt must be a pointer, as we can't hide non-nil structs
	JoinedAt    *discord.Timestamp  `json:"joined_at,omitempty"`    // ?*|
	Large       bool                `json:"large,omitempty"`        // ?*|
	Unavailable bool                `json:"unavailable"`            // ?*|
	MemberCount uint                `json:"member_count,omitempty"` // ?*|
	VoiceStates []*voice.State      `json:"voice_states,omitempty"` // ?*|
	Members     []*Member           `json:"members,omitempty"`      // ?*|
	Channels    []*channel.Channel  `json:"channels,omitempty"`     // ?*|
	Presences   []*discord.Presence `json:"presences,omitempty"`    // ?*|
}

type GuildInterface interface {
	ID() snowflake.ID
	Channel(ID snowflake.ID)
	Channels() []*channel.Channel
}

type GuildIDer interface {
	ID() snowflake.ID
}

type Guild struct {
	d guildJSON // struct data

	sync.RWMutex
}

func (g *Guild) ID() snowflake.ID {
	g.Lock()
	defer g.Unlock()

	return g.d.ID
}

// Channels
func (g *Guild) Channels() []*channel.Channel {
	return g.d.Channels
}

// Compare two guild objects
func (g *Guild) Compare(other *Guild) bool {
	// TODO: this is shit..
	g.Lock()
	defer g.Unlock()

	return (g == nil && other == nil) || (other != nil && g.d.ID == other.d.ID)
}

func (g *Guild) UnmarshalJSON(data []byte) (err error) {
	g.Lock()
	defer g.Unlock()

	return json.Unmarshal(data, &g.d)
}

func (g *Guild) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	var err error
	if g.d.Unavailable {
		guildUnavailable := struct {
			ID          snowflake.ID `json:"id"`
			Unavailable bool         `json:"unavailable"` // ?*|
		}{
			ID:          g.d.ID,
			Unavailable: true,
		}
		jsonData, err = json.Marshal(&guildUnavailable)
		if err != nil {
			return []byte(""), nil
		}
	} else {
		jsonData, err = json.Marshal(g.d)
		if err != nil {
			return []byte(""), nil
		}
	}

	return jsonData, nil
}

// sortChannels Only while in lock
func (g *Guild) sortChannels() {
	sort.Slice(g.d.Channels, func(i, j int) bool {
		return g.d.Channels[i].ID < g.d.Channels[j].ID
	})
}

func (g *Guild) AddChannel(c *channel.Channel) error {
	g.Lock()
	defer g.Unlock()

	g.d.Channels = append(g.d.Channels, c)
	g.sortChannels()

	return nil
}

func (g *Guild) DeleteChannel(c *channel.Channel) error {
	return g.DeleteChannelByID(c.ID)
}
func (g *Guild) DeleteChannelByID(ID snowflake.ID) error {
	g.Lock()
	defer g.Unlock()

	index := -1
	for i, c := range g.d.Channels {
		if c.ID == ID {
			index = i
		}
	}

	if index == -1 {
		return errors.New("channel with ID{" + ID.String() + "} does not exist in cache")
	}

	// delete the entry
	g.d.Channels[index] = g.d.Channels[len(g.d.Channels)-1]
	g.d.Channels[len(g.d.Channels)-1] = nil
	g.d.Channels = g.d.Channels[:len(g.d.Channels)-1]

	g.sortChannels()

	return nil
}

func (g *Guild) AddMember(member *Member) error {
	g.Lock()
	defer g.Unlock()

	// TODO: implement sorting for faster searching later
	g.d.Members = append(g.d.Members, member)

	return nil
}

func (g *Guild) AddRole(role *discord.Role) error {
	g.Lock()
	defer g.Unlock()

	// TODO: implement sorting for faster searching later
	g.d.Roles = append(g.d.Roles, role)

	return nil
}

// Member return a member by his/her userid
func (g *Guild) Member(id snowflake.ID) (*Member, error) {
	g.RLock()
	defer g.RUnlock()

	for _, member := range g.d.Members {
		if member.User.ID == id {
			return member, nil
		}
	}

	return nil, errors.New("member not found in guild")
}

// MemberByName retrieve a slice of members with same username or nickname
func (g *Guild) MemberByName(name string) ([]*Member, error) {
	g.RLock()
	defer g.RUnlock()

	var members []*Member
	for _, member := range g.d.Members {
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
func (g *Guild) Role(id snowflake.ID) (*discord.Role, error) {
	g.RLock()
	defer g.RUnlock()

	for _, role := range g.d.Roles {
		if role.ID == id {
			return role, nil
		}
	}

	return nil, errors.New("role not found in guild")
}

// RoleByTitle retrieves a slice of roles with same name
func (g *Guild) RoleByName(name string) ([]*discord.Role, error) {
	g.RLock()
	defer g.RUnlock()

	var roles []*discord.Role
	for _, role := range g.d.Roles {
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
	g.RLock()
	defer g.RUnlock()

	for _, channel := range g.d.Channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, errors.New("channel not found in guild")
}

// Update update the reference content
func (g *Guild) Update(new *Guild) {

}

// Clear all the pointers
func (g *Guild) Clear() {
	g.Lock() // what if another process tries to read this, but awais while locked for clearing?
	defer g.Unlock()

	g.d.ApplicationID = nil
	//g.Icon = nil // should this be cleared?
	//g.Splash = nil // should this be cleared?

	for _, r := range g.d.Roles {
		r.Clear()
		r = nil
	}
	g.d.Roles = nil

	for _, e := range g.d.Emojis {
		e.Clear()
		e = nil
	}
	g.d.Emojis = nil

	g.d.SystemChannelID = nil
	g.d.JoinedAt = nil

	for _, vst := range g.d.VoiceStates {
		vst.Clear()
		vst = nil
	}
	g.d.VoiceStates = nil

	deletedUsers := []snowflake.ID{}
	for _, m := range g.d.Members {
		deletedUsers = append(deletedUsers, m.Clear())
		m = nil
	}
	g.d.Members = nil

	for _, c := range g.d.Channels {
		c.Clear()
		c = nil
	}
	g.d.Channels = nil

	for _, p := range g.d.Presences {
		p.Clear()
		p = nil
	}
	g.d.Presences = nil

}
