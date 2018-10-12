package endpoint

import "fmt"

// GuildEmojis /guilds/{guild.id}/emojis
func GuildEmojis(id fmt.Stringer) string {
	return Guild(id) + emojis
}

// GuildEmoji /guilds/{guild.id}/emojis/{emoji.id}
func GuildEmoji(guildID, emojiID fmt.Stringer) string {
	return GuildEmojis(guildID) + "/" + emojiID.String()
}
