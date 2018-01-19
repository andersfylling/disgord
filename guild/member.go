package guild

import (
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

// Member ...
type Member struct {
	GuildID  snowflake.ID      `json:"guild_id,omitempty"`
	User     *user.User        `json:"user"`
	Nick     string            `json:"nick,omitempty"` // ?|
	Roles    []snowflake.ID    `json:"roles"`
	JoinedAt discord.Timestamp `json:"joined_at,omitempty"`
	Deaf     bool              `json:"deaf"`
	Mute     bool              `json:"mute"`
}
