package channel

import (
	"fmt"

	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

type Channel struct {
	ID                   snowflake.ID           `json:"id"`
	GuildID              snowflake.ID           `json:"guild_id"`
	Name                 string                 `json:"name"`
	Topic                string                 `json:"topic"`
	Type                 uint                   `json:"type"`
	LastMessageID        snowflake.ID           `json:"last_message_id"`
	NSFW                 bool                   `json:"nsfw"`
	Position             uint                   `json:"position"`
	Bitrate              int                    `json:"bitrate"`
	Recipients           []*user.User           `json:"recipient"`
	Messages             []*Message             `json:"-"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`
}

func NewChannel() *Channel {
	return &Channel{}
}

func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%d>", c.ID)
}

func (channel *Channel) Compare(c *Channel) bool {
	return c != nil && channel.ID == c.ID
}
