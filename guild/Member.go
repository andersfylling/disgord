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

func (m *Member) Clear() snowflake.ID {
	// do i want to delete user?.. what if there is a PM?
	// Check for user id in DM's
	// or.. since the user object is sent on channel_create events, the user can be reintialized when needed.
	// but should be properly removed from other arrays.
	m.User.Clear()
	id := m.User.ID()
	m.User = nil

	// use this ID to check in other places. To avoid pointing to random memory spaces
	return id
}
