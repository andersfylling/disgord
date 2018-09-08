package event

import (
	"context"

	"github.com/andersfylling/disgord/resource"
	. "github.com/andersfylling/snowflake"
)

// KeyGuildEmojisUpdate Sent when a guild's emojis have been updated.
//                   Fields:
//                   * GuildID int64 or discord.Snowflake
const KeyGuildEmojisUpdate = "GUILD_EMOJI_UPDATE"

// GuildEmojisUpdate	guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID Snowflake         `json:"guild_id"`
	Emojis  []*resource.Emoji `json:"emojis"`
	Ctx     context.Context   `json:"-"`
}
