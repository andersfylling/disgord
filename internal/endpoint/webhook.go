package endpoint

import "fmt"

// Slack /slack suffix
func Slack() string {
	return slack
}

// GitHub /github suffix
func GitHub() string {
	return github
}

// Webhook /webhooks/{webhook.id}
func Webhook(id fmt.Stringer) string {
	return webhooks + "/" + id.String()
}

// WebhookToken /webhooks/{webhook.id}/{webhook.token}
func WebhookToken(id fmt.Stringer, token string) string {
	return Webhook(id) + "/" + token
}

// ChannelWebhooks /channels/{channel.id}/webhooks
func ChannelWebhooks(id fmt.Stringer) string {
	return Channel(id) + webhooks
}

// GuildWebhooks /guilds/{guild.id}/webhooks
func GuildWebhooks(id fmt.Stringer) string {
	return Guild(id) + webhooks
}
