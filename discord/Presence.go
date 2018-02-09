package discord

import "github.com/andersfylling/snowflake"
import "github.com/andersfylling/disgord/user"

type Presence struct {
	User    *user.User     `json:"user"`
	Roles   []snowflake.ID `json:"roles"`
	Game    *Activity      `json:"activty"`
	GuildID snowflake.ID   `json:"guild_id"`
	Nick    *string        `json:"nick"`
	Status  string         `json:"status"`
}

func NewPresence() *Presence {
	return &Presence{}
}

func (p *Presence) Update(status string) {
	// Update the presence.
	// talk to the discord api
}

func (p *Presence) String() string {
	return p.Status
}

func (p *Presence) Clear() {
	p.Game = nil
}
