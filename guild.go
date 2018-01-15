package disgord

import (
	"time"

	"github.com/andersfylling/snowflake"
)

type Guild struct {
	// Unique ID of the guild
	ID snowflake.ID `json:"id"`
	// Name of the guild
	Name string `json:"name"`

	// Icon is an image hash
	Icon string `json:"icon"`
	// Splash is an image hash
	Splash string `json:"splash"`

	// OwnerID is the unique user ID of the guild's owner
	OwnerID                     snowflake.ID   `json:"owner_id"`
	Region                      string         `json:"region"`
	AfkChannelID                snowflake.ID   `json:"afk_channel_id"`
	EmbedChannelID              snowflake.ID   `json:"embed_channel_id"`
	JoinedAt                    time.Time      `json:"joined_at"`
	AfkTimeout                  uint           `json:"afk_timeout"`
	MemberCount                 uint           `json:"member_count"`
	VerificationLevel           uint           `json:"verification_level"`
	EmbedEnabled                bool           `json:"embed_enabled"`
	Large                       bool           `json:"large"` // ??
	DefaultMessageNotifications int            `json:"default_message_notifications"`
	Roles                       []*Role        `json:"roles"`
	Emojis                      []*Emoji       `json:"emojis"`
	Members                     []*GuildMember `json:"members"`
	Presences                   []*Presence    `json:"presences"`
	Channels                    []*Channel     `json:"channels"`
	VoiceStates                 []*VoiceState  `json:"voice_states"`
	Unavailable                 bool           `json:"unavailable"`
}

func (guild *Guild) Compare(g *Guild) bool {
	return g != nil && guild.ID == g.ID
}
