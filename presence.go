package disgord

import "github.com/andersfylling/snowflake"

type Presence struct {
	User    *User          `json:"user"`
	Roles   []snowflake.ID `json:"roles"`
	Game    *Activity      `json:"activty"`
	GuildID snowflake.ID   `json:"guild_id"`
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
