package endpoint

import . "github.com/andersfylling/snowflake"

func GuildEmojis(id Snowflake) string {
	return Guild(id) + emojis
}

func GuildEmoji(guildID, emojiID Snowflake) string {
	return GuildEmojis(guildID) + "/" + emojiID.String()
}
