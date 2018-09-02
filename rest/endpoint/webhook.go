package endpoint

import . "github.com/andersfylling/snowflake"

// Slack /slack suffix
func Slack() string {
	return slack
}

// GitHub /github suffix
func GitHub() string {
	return github
}

// Webhook /webhooks/{webhook.id}
func Webhook(id Snowflake) string {
	return webhooks + "/" + id.String()
}

// WebhookToken /webhooks/{webhook.id}/{webhook.token}
func WebhookToken(id Snowflake, token string) string {
	return Webhook(id) + "/" + token
}

// ChannelWebhooks /channels/{channel.id}/webhooks
func ChannelWebhooks(id Snowflake) string {
	return Channel(id) + webhooks
}

// GuildWebhooks /guilds/{guild.id}/webhooks
func GuildWebhooks(id Snowflake) string {
	return Guild(id) + webhooks
}
