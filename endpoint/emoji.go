package endpoint

import "fmt"

func GuildEmojis(id fmt.Stringer) string {
	return Guild(id) + emojis
}

func GuildEmoji(guildID, emojiID fmt.Stringer) string {
	return GuildEmojis(guildID) + "/" + emojiID.String()
}
