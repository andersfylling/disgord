package event

import (
	"context"

	. "github.com/andersfylling/snowflake"
)

// WebhooksUpdate Sent when a guild channel's webhook is created, updated, or
//                deleted.
//                Fields:
//                * GuildID   int64 or discord.Snowflake
//                * ChannelID int64 or discord.Snowflake
const KeyWebhooksUpdate = "WEBHOOK_UPDATE"

// WebhooksUpdate guild channel webhook was created, update, or deleted
type WebhooksUpdate struct {
	GuildID   Snowflake       `json:"guild_id"`
	ChannelID Snowflake       `json:"channel_id"`
	Ctx       context.Context `json:"-"`
}
