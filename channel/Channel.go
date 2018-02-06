package channel

import (
	"fmt"

	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/user"
	"github.com/andersfylling/snowflake"
)

type Channel struct {
	ID                   snowflake.ID           `json:"id"`
	Type                 uint                   `json:"type"`
	GuildID              snowflake.ID           `json:"guild_id,omitempty"`
	Position             uint                   `json:"position,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	Topic                string                 `json:"topic,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
	LastMessageID        *snowflake.ID          `json:"last_message_id,omitempty"`
	Bitrate              uint                   `json:"bitrate,omitempty"`
	UserLimit            uint                   `json:"user_limit,omitempty"`
	Recipients           []*user.User           `json:"recipient,omitempty"`
	Icon                 *string                `json:"icon,omitempty"`
	OwnerID              snowflake.ID           `json:owner_id,omitempty`
	ApplicationID        snowflake.ID           `json:"applicaiton_id,omitempty"`
	ParentID             snowflake.ID           `json:"parent_id,omitempty"`
	LastPingTimestamp    discord.Timestamp      `json:"last_ping_timestamp,omitempty"`

	// Messages used for caching only. is always empty when fresh from the discord API
	Messages []*Message `json:"-"` // should prolly set a cache limit of 100
}

func NewChannel() *Channel {
	return &Channel{}
}

func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%d>", c.ID)
}

func (c *Channel) Compare(other *Channel) bool {
	return (c == nil && other == nil) || (other != nil && c.ID == other.ID)
}

func (c *Channel) Clear() {
	c.LastMessageID = nil
	// c.Icon = nil // Do I really want to clear this?
	for _, pmo := range c.PermissionOverwrites {
		pmo.Clear()
		pmo = nil
	}
	c.PermissionOverwrites = nil

	//for _,
}
