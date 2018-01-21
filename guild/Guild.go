package guild

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/emoji"
	"github.com/andersfylling/disgord/lvl"
	"github.com/andersfylling/disgord/voice"
	"github.com/andersfylling/snowflake"
)

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
type Guild struct {
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

	sync.RWMutex `json:"-"`
}

// Compare two guild objects
func (guild *Guild) Compare(g *Guild) bool {
	return (guild == nil && g == nil) || (g != nil && guild.ID == g.ID)
}

func (guild *Guild) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	var err error
	if guild.Unavailable {
		guildUnavailable := struct {
			ID          snowflake.ID `json:"id"`
			Unavailable bool         `json:"unavailable"` // ?*|
		}{
			ID:          guild.ID,
			Unavailable: true,
		}
		jsonData, err = json.Marshal(&guildUnavailable)
		if err != nil {
			return []byte(""), nil
		}
	} else {
		g := Guild(*guild) // avoid stack overflow by recursive call of Marshal
		jsonData, err = json.Marshal(g)
		if err != nil {
			return []byte(""), nil
		}
	}

	return jsonData, nil
}

func (guild *Guild) AddChannel(channel *channel.Channel) error {
	guild.Lock()
	guild.Unlock()

	// TODO: implement sorting for faster searching later
	guild.Channels = append(guild.Channels, channel)

	return nil
}

func (guild *Guild) AddMember(member *Member) error {
	guild.Lock()
	guild.Unlock()

	// TODO: implement sorting for faster searching later
	guild.Members = append(guild.Members, member)

	return nil
}

func (guild *Guild) AddRole(role *discord.Role) error {
	guild.Lock()
	guild.Unlock()

	// TODO: implement sorting for faster searching later
	guild.Roles = append(guild.Roles, role)

	return nil
}

// Member return a member by his/her userid
func (guild *Guild) Member(id snowflake.ID) (*Member, error) {
	guild.RLock()
	defer guild.RUnlock()

	for _, member := range guild.Members {
		if member.User.ID == id {
			return member, nil
		}
	}

	return nil, errors.New("member not found in guild")
}

// MemberByName retrieve a slice of members with same username or nickname
func (guild *Guild) MemberByName(name string) ([]*Member, error) {
	guild.RLock()
	defer guild.RUnlock()

	var members []*Member
	for _, member := range guild.Members {
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
func (guild *Guild) Role(id snowflake.ID) (*discord.Role, error) {
	guild.RLock()
	defer guild.RUnlock()

	for _, role := range guild.Roles {
		if role.ID == id {
			return role, nil
		}
	}

	return nil, errors.New("role not found in guild")
}

// RoleByTitle retrieves a slice of roles with same name
func (guild *Guild) RoleByName(name string) ([]*discord.Role, error) {
	guild.RLock()
	defer guild.RUnlock()

	var roles []*discord.Role
	for _, role := range guild.Roles {
		if role.Name == name {
			roles = append(roles, role)
		}
	}

	if len(roles) == 0 {
		return nil, errors.New("no roles were found in guild")
	}

	return roles, nil
}

func (guild *Guild) Channel(id snowflake.ID) (*channel.Channel, error) {
	guild.RLock()
	defer guild.RUnlock()

	for _, channel := range guild.Channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, errors.New("channel not found in guild")
}
