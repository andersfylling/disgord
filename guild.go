package disgord

import (
	"encoding/json"

	"github.com/andersfylling/snowflake"
)

// DefaultMessageNotificationLevel ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type DefaultMessageNotificationLevel uint

func (dmnl *DefaultMessageNotificationLevel) AllMessages() bool {
	return *dmnl == 0
}
func (dmnl *DefaultMessageNotificationLevel) OnlyMentions() bool {
	return *dmnl == 1
}

// ExplicitContentFilterLevel ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
type ExplicitContentFilterLevel uint

func (ecfl *ExplicitContentFilterLevel) Disabled() bool {
	return *ecfl == 0
}
func (ecfl *ExplicitContentFilterLevel) MembersWithoutRoles() bool {
	return *ecfl == 1
}
func (ecfl *ExplicitContentFilterLevel) AllMembers() bool {
	return *ecfl == 2
}

// MFALevel ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-mfa-level
type MFALevel uint

func (mfal *MFALevel) None() bool {
	return *mfal == 0
}
func (mfal *MFALevel) Elevated() bool {
	return *mfal == 1
}

// VerificationLevel ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-verification-level
type VerificationLevel uint

// None unrestricted
func (vl *VerificationLevel) None() bool {
	return *vl == 0
}

// Low must have verified email on account
func (vl *VerificationLevel) Low() bool {
	return *vl == 1
}

// Medium must be registered on Discord for longer than 5 minutes
func (vl *VerificationLevel) Medium() bool {
	return *vl == 2
}

// High (╯°□°）╯︵ ┻━┻ - must be a member of the server for longer than 10 minutes
func (vl *VerificationLevel) High() bool {
	return *vl == 3
}

// VeryHigh ┻━┻ミヽ(ಠ益ಠ)ﾉ彡┻━┻ - must have a verified phone number
func (vl *VerificationLevel) VeryHigh() bool {
	return *vl == 4
}

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
type Guild struct {
	ID                          snowflake.ID     `json:"id,string"`
	ApplicationID               *snowflake.ID    `json:"application_id"` //   |?
	Name                        string           `json:"name"`
	Icon                        *string          `json:"icon"`            //  |?, icon hash
	Splash                      *string          `json:"splash"`          //  |?, image hash
	Owner                       bool             `json:"owner,omitempty"` // ?|
	OwnerID                     snowflake.ID     `json:"owner_id,string"`
	Permissions                 uint64           `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
	Region                      string           `json:"region"`
	AfkChannelID                snowflake.ID     `json:"afk_channel_id,string"`
	AfkTimeout                  uint             `json:"afk_timeout"`
	EmbedEnabled                bool             `json:"embed_enabled"`
	EmbedChannelID              snowflake.ID     `json:"embed_channel_id,string"`
	VerificationLevel           uint             `json:"verification_level"`
	DefaultMessageNotifications uint             `json:"default_message_notifications"`
	ExplicitContentFilter       uint             `json:"explicit_content_filter"`
	MFALevel                    uint             `json:"mfa_level"`
	WidgetEnabled               bool             `json:"widget_enabled"`           //   |
	WidgetChannelID             snowflake.ID     `json:"widget_channel_id,string"` //   |
	Roles                       []*Role          `json:"roles"`
	Emojis                      []*Emoji         `json:"emojis"`
	Features                    []string         `json:"features"`
	SystemChannelID             *snowflake.ID    `json:"system_channel_id,string,omitempty"` //   |?
	JoinedAt                    DiscordTimestamp `json:"joined_at,omitempty"`                // ?*|
	Large                       bool             `json:"large,omitempty"`                    // ?*|
	Unavailable                 bool             `json:"unavailable"`                        // ?*|
	MemberCount                 uint             `json:"member_count,omitempty"`             // ?*|
	VoiceStates                 []*VoiceState    `json:"voice_states,omitempty"`             // ?*|
	Members                     []*GuildMember   `json:"members,omitempty"`                  // ?*|
	Channels                    []*Channel       `json:"channels,omitempty"`                 // ?*|
	Presences                   []*Presence      `json:"presences,omitempty"`                // ?*|
}
type GuildUnavailable struct {
	ID          snowflake.ID `json:"id,string"`
	Unavailable bool         `json:"unavailable"` // ?*|
}

// Compare two guild objects
func (guild *Guild) Compare(g *Guild) bool {
	return (guild == nil && g == nil) || (g != nil && guild.ID == g.ID)
}

func (guild *Guild) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	var err error
	if guild.Unavailable {
		guildUnavailable := GuildUnavailable{ID: guild.ID, Unavailable: true}
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
