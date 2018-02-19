package guild

import (
	"errors"
	"sync"

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

	sync.RWMutex `json:"-"`
}

func (m *Member) Clear() snowflake.ID {
	// do i want to delete user?.. what if there is a PM?
	// Check for user id in DM's
	// or.. since the user object is sent on channel_create events, the user can be reintialized when needed.
	// but should be properly removed from other arrays.
	m.User.Clear()
	id := m.User.ID
	m.User = nil

	// use this ID to check in other places. To avoid pointing to random memory spaces
	return id
}

func (m *Member) Update(new *Member) (err error) {
	if m.User.ID != new.User.ID || m.GuildID != new.GuildID {
		err = errors.New("cannot update user when the new struct has a different ID")
		return
	}
	// make sure that new is not the same pointer!
	if m == new {
		err = errors.New("cannot update user when the new struct points to the same memory space")
		return
	}

	m.Lock()
	new.RLock()
	m.Nick = new.Nick
	m.Roles = new.Roles
	m.JoinedAt = new.JoinedAt
	m.Deaf = new.Deaf
	m.Mute = new.Mute
	new.RUnlock()
	m.Unlock()

	return
}
