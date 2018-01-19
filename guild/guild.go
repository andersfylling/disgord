package guild

import (
	"encoding/json"

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
