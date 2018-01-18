package disgord

import (
	"github.com/andersfylling/snowflake"
)

// GuildMember ...
type GuildMember struct {
	GuildID  snowflake.ID     `json:"guild_id,string,omitempty"`
	User     *User            `json:"user"`
	Nick     string           `json:"nick,omitempty"` // ?|
	Roles    []snowflake.ID   `json:"roles,string"`
	JoinedAt DiscordTimestamp `json:"joined_at,omitempty"`
	Deaf     bool             `json:"deaf"`
	Mute     bool             `json:"mute"`
}
